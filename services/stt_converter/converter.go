package stt_converter

import (
	"fmt"
	"stt_back/common"
)

type Result struct {
	Text      string
	RawResult string
}

func ConvertSpeechToText(data []byte) (res Result, err error) {
	// TODO: Implement this
	res.Text = common.RandomString(20)
	res.RawResult = fmt.Sprintf(`{"result":"%s"}`, res.Text)

	return
}
