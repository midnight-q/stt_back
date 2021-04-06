package stt_converter

import (
	"fmt"
	"sort"
	"strings"
)

type Params struct {
	TimeFrame         int
	IsShowEmotion     bool
	IsShowSpeaker     bool
	IsShowTag         bool
	IsShowPunctuation bool
}

func ConvertDataToText(data []Data, params Params) (res string) {
	resultArr := []string{}
	sort.Slice(data, func(i, j int) bool {
		return data[i].TimeStart < data[j].TimeStart
	})
	for _, d := range data {
		resultArr = append(resultArr, applyTextTemplate(d, params))
	}
	return strings.Join(resultArr, "\n")
}

func ConvertDataToHtml(data []Data, params Params) (res string) {
	resultArr := []string{}
	sort.Slice(data, func(i, j int) bool {
		return data[i].TimeStart < data[j].TimeStart
	})
	for _, d := range data {
		resultArr = append(resultArr, applyHtmlTemplate(d, params))
	}
	return strings.Join(resultArr, "")
}

func ConvertDataToPdf(data []Data, params Params) (path string) {

	return path
}

func ConvertDataToDoc(data []Data, params Params) (path string) {

	return path
}

func applyTextTemplate(data Data, params Params) string {
	str := data.RawText
	if params.IsShowPunctuation {
		str = data.Text
	}
	if params.IsShowSpeaker || params.IsShowEmotion {
		str = ":" + str
	}
	if params.IsShowEmotion {
		str = "(" + data.Emotion + ")" + str
	}
	if params.IsShowSpeaker {
		str = data.Speaker + str
	}
	return str
}

func applyHtmlTemplate(data Data, params Params) string {
	str := fmt.Sprintf("<div>%s</div>", data.RawText)
	if params.IsShowPunctuation {
		str = fmt.Sprintf("<div>%s</div>", data.Text)

	}
	if params.IsShowEmotion {
		str = fmt.Sprintf("<div>%s</div>", data.Emotion) + str
	}
	if params.IsShowSpeaker {
		str = fmt.Sprintf("<div>%s</div>", data.Speaker) + str
	}
	return fmt.Sprintf("<p>%s</p>", str)
}
