package generator

import (
    "stt_back/types"
    "math/rand"
    "strings"
)

func GenConvertFileV2() types.ConvertFileV2 {

	return types.ConvertFileV2{
		Id:   rand.Intn(100500),
		Data: strings.Title(Babbler2.Babble()),
		ResultHtmlPath: strings.Title(Babbler2.Babble()),
		ResultFilePdfPath: strings.Title(Babbler2.Babble()),
		ResultTextPath: strings.Title(Babbler2.Babble()),
		ResultFileDocPath: strings.Title(Babbler2.Babble()),
		SourceFilePath: strings.Title(Babbler2.Babble()),
		ResultText: strings.Title(Babbler2.Babble()),
		TimeEstimate: rand.Intn(100500),
		//ConvertFileV2 remove this line for disable generator functionality
	}
}

func GenListConvertFileV2() (list []types.ConvertFileV2) {

	for i:=0; i<rand.Intn(5) + 2; i++{
		list = append(list, GenConvertFileV2())
	}

	return
}
