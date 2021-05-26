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
		SourceFilePath: strings.Title(Babbler2.Babble()),
		RecordNumber: rand.Intn(100500),
		TimeFrame: rand.Intn(100500),
		IsShowEmotion: (rand.Intn(100500) % 2 > 0),
		IsShowSpeaker: (rand.Intn(100500) % 2 > 0),
		IsShowTag: (rand.Intn(100500) % 2 > 0),
		IsShowPunctuation: (rand.Intn(100500) % 2 > 0),
		NamedEntityTypes: strings.Title(Babbler2.Babble()),
		Status: strings.Title(Babbler2.Babble()),
		Token: strings.Title(Babbler2.Babble()),
		ResultText: strings.Title(Babbler2.Babble()),
		ErrorString: strings.Title(Babbler2.Babble()),
		//ConverterLog remove this line for disable generator functionality
	}
}

func GenListConverterLog() (list []types.ConverterLog) {

	for i := 0; i < rand.Intn(5)+2; i++ {
		list = append(list, GenConverterLog())
	}

	return
}
