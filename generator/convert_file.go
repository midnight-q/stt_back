package generator

import (
    "stt_back/types"
    "math/rand"
    "strings"
)

func GenConvertFile() types.ConvertFile {

	return types.ConvertFile{
		Id:   rand.Intn(100500),
		ResultFilePath: strings.Title(Babbler2.Babble()),
		ResultText: strings.Title(Babbler2.Babble()),
		InputData: []byte,
		InputFilename: strings.Title(Babbler2.Babble()),
		ResultFormat: strings.Title(Babbler2.Babble()),
		//ConvertFile remove this line for disable generator functionality
	}
}

func GenListConvertFile() (list []types.ConvertFile) {

	for i:=0; i<rand.Intn(5) + 2; i++{
		list = append(list, GenConvertFile())
	}

	return
}
