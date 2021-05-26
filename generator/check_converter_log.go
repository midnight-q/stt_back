package generator

import (
    "stt_back/types"
    "math/rand"
    "strings"
)

func GenCheckConverterLog() types.CheckConverterLog {

	return types.CheckConverterLog{
		Id:   rand.Intn(100500),
		ConverterLog: strings.Title(Babbler2.Babble()),
		TimeEstimate: rand.Intn(100500),
		//CheckConverterLog remove this line for disable generator functionality
	}
}

func GenListCheckConverterLog() (list []types.CheckConverterLog) {

	for i:=0; i<rand.Intn(5) + 2; i++{
		list = append(list, GenCheckConverterLog())
	}

	return
}
