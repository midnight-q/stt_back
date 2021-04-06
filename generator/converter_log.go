package generator

import (
	"math/rand"
	"strings"
	"stt_back/types"
)

func GenConverterLog() types.ConverterLog {

	return types.ConverterLog{
		Id:                rand.Intn(100500),
		FilePath:          strings.Title(Babbler2.Babble()),
		ResultText:        strings.Title(Babbler2.Babble()),
		ResultFilePath:    strings.Title(Babbler2.Babble()),
		ResultFormat:      strings.Title(Babbler2.Babble()),
		RawResult:         strings.Title(Babbler2.Babble()),
		ResultHtml:        strings.Title(Babbler2.Babble()),
		ResultFileDocPath: strings.Title(Babbler2.Babble()),
		ResultFilePdfPath: strings.Title(Babbler2.Babble()),
		UserId:            rand.Intn(100500),
		//ConverterLog remove this line for disable generator functionality
	}
}

func GenListConverterLog() (list []types.ConverterLog) {

	for i := 0; i < rand.Intn(5)+2; i++ {
		list = append(list, GenConverterLog())
	}

	return
}
