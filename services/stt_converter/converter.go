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
			EndTime       int `json:"end_time"`
			NamedEntities struct {
			} `json:"named_entities"`
			Relations struct {
			} `json:"relations"`
			Sent        string `json:"sent"`
			SpeakerName string `json:"speaker_name"`
			StartTime   int    `json:"start_time"`
			Text        string `json:"text"`
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
			EndTime       int `json:"end_time"`
			NamedEntities struct {
			} `json:"named_entities"`
			Sent      string `json:"sent"`
			Sentiment struct {
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

	lastPhrase := 0
	for _, ner := range resultData.Result.Ner {
		if lastPhrase < ner.StartTime {
			lastPhrase = ner.StartTime
		}
		res.Data = append(res.Data, Data{
			TimeStart: ner.StartTime,
			TimeEnd:   ner.EndTime,
			Text:      ner.Text,
			RawText:   clearString(ner.Sent),
			Speaker:   ner.SpeakerName,
			Emotion:   getEmotionFromResult(resultData, ner.StartTime, ner.EndTime),
			Tags:      covertTags(ner.NamedEntities, params),
		})
	}

	for i := 0; i < lastPhrase/(params.TimeFrame*1000); i++ {
		t := time.Time{}
		t = t.Add(time.Duration(params.TimeFrame * i * int(time.Second)))
		res.Data = append(res.Data, Data{
			TimeStart:        params.TimeFrame * i * 1000,
			Text:             t.Format("04:05"),
			RawText:          t.Format("04:05"),
			IsTimeFrameLabel: true,
		})
	}

	return
}

func clearString(sent string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9а-яА-ЯЁё\\- ]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(sent, "")
	return strings.ToLower(processedString)
}

func covertTags(entities map[string][][]int, convertParams Params) (res []Tag) {
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
