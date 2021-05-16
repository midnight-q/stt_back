package stt_converter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"stt_back/common"
	"stt_back/errors"
	"time"
	"unicode/utf8"
)

type ResultData struct {
	Result struct {
		Diarization []struct {
			EndTime     int    `json:"end_time"`
			SpeakerName string `json:"speaker_name"`
			StartTime   int    `json:"start_time"`
		} `json:"diarization"`
		Ner []struct {
			EndTime       int                `json:"end_time"`
			NamedEntities map[string][][]int `json:"named_entities"`
			Sent          string             `json:"sent"`
			SpeakerName   string             `json:"speaker_name"`
			StartTime     int                `json:"start_time"`
			Text          string             `json:"text"`
		} `json:"ner"`
		Re []struct {
			EndTime       int      `json:"end_time"`
			NamedEntities struct{} `json:"named_entities"`
			Relations     struct{} `json:"relations"`
			Sent          string   `json:"sent"`
			SpeakerName   string   `json:"speaker_name"`
			StartTime     int      `json:"start_time"`
			Text          string   `json:"text"`
		} `json:"re"`
		Stt []struct {
			EndTime   int    `json:"end_time"`
			StartTime int    `json:"start_time"`
			Word      string `json:"word"`
		} `json:"stt"`
		SttDictors []struct {
			EndTime       int `json:"end_time"`
			NamedEntities struct {
			} `json:"named_entities"`
			Sent        string `json:"sent"`
			SpeakerName string `json:"speaker_name"`
			StartTime   int    `json:"start_time"`
			Text        string `json:"text"`
		} `json:"stt.dictors"`
		SttPunct []struct {
			EndTime   int    `json:"end_time"`
			Sent      string `json:"sent"`
			StartTime int    `json:"start_time"`
		} `json:"stt.punct"`
		Toxic []struct {
			EndTime       int      `json:"end_time"`
			NamedEntities struct{} `json:"named_entities"`
			Sent          string   `json:"sent"`
			Sentiment     struct {
				Certainty float64 `json:"certainty"`
				Label     string  `json:"label"`
				Output    int     `json:"output"`
			} `json:"sentiment"`
			SpeakerName string `json:"speaker_name"`
			StartTime   int    `json:"start_time"`
			Text        string `json:"text"`
		} `json:"toxic"`
		Vad [][]int `json:"vad"`
	} `json:"result"`
}

type Data struct {
	TimeStart        int
	TimeEnd          int
	Text             string
	RawText          string
	Speaker          string
	Emotion          string
	Tags             []Tag
	IsTimeFrameLabel bool
}

type Tag struct {
	Name  string
	Start int
	End   int
}

type Result struct {
	RawResult string
	RawData   ResultData
	Data      []Data
}

func ConvertSpeechToText(data []byte, params Params) (res Result, err error) {
	reader := bytes.NewReader(data)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("wav", "file.wav")

	if err != nil {
		log.Fatal(err)
	}

	io.Copy(part, reader)
	writer.Close()
	request, err := http.NewRequest("POST", "https://ai.nsu.ru:7777/a8cb46de23/recognize", body)

	if err != nil {
		return Result{}, err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		return Result{}, err
	}

	if response.StatusCode != 201 {
		rawString, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return Result{}, err
		}
		fmt.Println(string(rawString), response.StatusCode)

		err = errors.New("Error in request to converter: " + string(rawString))
		return Result{}, err
	}

	defer response.Body.Close()

	rawString, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Result{}, err
	}
	var resultData = ResultData{}
	err = json.Unmarshal(rawString, &resultData)
	if err != nil {
		return Result{}, err
	}

	res.RawData = resultData
	res.RawResult = string(rawString)


	rawRes := []Data{}
	for _, ner := range resultData.Result.Ner {
		rawRes = append(rawRes, Data{
			TimeStart: ner.StartTime,
			TimeEnd:   ner.EndTime,
			Text:      ner.Text,
			RawText:   clearString(ner.Sent),
			Speaker:   getSpeakerName(ner.SpeakerName),
			Emotion:   getEmotionFromResult(resultData, ner.StartTime, ner.EndTime),
			Tags:      convertTags(ner.NamedEntities, params),
		})
	}

	timeFrame := params.TimeFrame * 1000

	//Split phrase if collision with timeFrame marker exist
	for _, d := range rawRes {
		if d.TimeStart/timeFrame != d.TimeEnd/timeFrame {
			collisionCount := d.TimeEnd/timeFrame - d.TimeStart/timeFrame
			splittedPhrase := splitPhrase(d, collisionCount, timeFrame, &resultData)
			res.Data = append(res.Data, splittedPhrase...)
		} else {
			res.Data = append(res.Data, d)
		}
	}

	lastPhrase := 0
	for _,d := range res.Data {
		if lastPhrase < d.TimeStart {
			lastPhrase = d.TimeStart
		}
	}
	//Add timeFrame markers
	for i := 0; i < lastPhrase/timeFrame+1; i++ {
		t := time.Time{}
		t = t.Add(time.Duration(params.TimeFrame * i * int(time.Second)))
		res.Data = append(res.Data, Data{
			TimeStart:        timeFrame * i,
			TimeEnd:          timeFrame * i,
			Text:             t.Format("04:05"),
			RawText:          t.Format("04:05"),
			IsTimeFrameLabel: true,
		})
	}

	return
}

func splitPhrase(d Data, count int, frame int, resultData *ResultData) (res []Data) {
	if count < 1 {
		res = append(res, d)
		return
	}

	collisionTime := (d.TimeStart/frame + 1) * frame
	for _, word := range resultData.Result.Stt {
		if word.StartTime < collisionTime {
			continue
		}
		if word.StartTime > d.TimeEnd {
			break
		}
		firstPart := Data{
			TimeStart: d.TimeStart,
			TimeEnd:   word.StartTime,
			Text:      splitText(d.Text, word.Word)[0],
			RawText:   splitText(d.RawText, word.Word)[0],
			Speaker:   d.Speaker,
			Emotion:   d.Emotion,
		}
		tagsFirst, tagsSecond := prepareTags(d.Tags, utf8.RuneCountInString(firstPart.Text))
		firstPart.Tags = tagsFirst

		secondPart := Data{
			TimeStart: word.StartTime,
			TimeEnd:   d.TimeEnd,
			Text:      splitText(d.Text, word.Word)[1],
			RawText:   splitText(d.RawText, word.Word)[1],
			Speaker:   d.Speaker,
			Emotion:   d.Emotion,
			Tags:      tagsSecond,
		}

		res = append(res, firstPart)
		if word.EndTime != d.TimeEnd {
			res = append(res, splitPhrase(secondPart, count-1, frame, resultData)...)
		} else {
			res = append(res, secondPart)
		}
		break
	}
	return
}

func prepareTags(tags []Tag, i int) (firstPart, secondPart []Tag) {
	for _, tag := range tags {
		if tag.End < i {
			firstPart = append(firstPart, tag)
		} else {
			secondPart = append(secondPart, Tag{
				Name:  tag.Name,
				Start: tag.Start - i,
				End:   tag.End - i,
			})
		}
	}
	return
}

func splitText(input, word string) []string {
	index := strings.Index(strings.ToLower(input), word)
	return []string{input[:index], input[index:]}
}

func getSpeakerName(name string) string {
	name = strings.TrimSpace(name)
	if strings.Index(name, "Спикер") > -1 {
		return name
	}
	return "Спикер " + name
}

func clearString(sent string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9а-яА-ЯЁё\\- ]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(sent, "")
	return strings.ToLower(processedString)
}

func convertTags(entities map[string][][]int, convertParams Params) (res []Tag) {
	for name, params := range entities {
		if !common.InArray(name, convertParams.NamedEntityTypes) {
			continue
		}
		for _, param := range params {
			res = append(res, Tag{
				Name:  name,
				Start: param[0],
				End:   param[1],
			})
		}
	}
	return
}

func getEmotionFromResult(data ResultData, start int, end int) string {
	for _, emotion := range data.Result.Toxic {
		if emotion.StartTime == start && emotion.EndTime == end {
			return emotion.Sentiment.Label
		}
	}
	return ""
}
