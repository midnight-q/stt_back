package logic

import (
	"fmt"
	"stt_back/common"
	"stt_back/core"
	"stt_back/services/stt_converter"
	"stt_back/types"
)

func CheckConverterLogRead(filter types.CheckConverterLogFilter) (data types.CheckConverterLog, err error) {
	f := types.ConverterLogFilter{}
	f.SetCurrentId(filter.GetCurrentId())
	log, err := ConverterLogRead(f)
	if err != nil {
		return types.CheckConverterLog{}, err
	}
	namedEntities := []string{}
	err = common.StringToData(log.NamedEntityTypes, &namedEntities)
	if err != nil {
		fmt.Println("Error common.StringToData:", err)
		return types.CheckConverterLog{}, err
	}

	converterParams := stt_converter.Params{
		TimeFrame:         log.TimeFrame,
		IsShowEmotion:     log.IsShowEmotion,
		IsShowSpeaker:     log.IsShowSpeaker,
		IsShowTag:         log.IsShowTag,
		IsShowPunctuation: log.IsShowPunctuation,
		NamedEntityTypes:  namedEntities,
	}

	result, err := stt_converter.CheckConvertStatus(log.Token, converterParams)

	if err != nil {
		return types.CheckConverterLog{}, err
	}

	data.Status = result.Status

	switch result.Status {
	case "Processing":
		data.TimeEstimate = result.TimeEstimate
		data.ConverterLog = log
		return
	case "Complete":
		log.Status = result.Status

		resultTextPath, resultText := stt_converter.ConvertDataToText(result.Data, converterParams)
		resultHtmlPath, _ := stt_converter.ConvertDataToHtml(result.Data, converterParams)
		resultFilePdfPath := stt_converter.ConvertDataToPdf(result.Data, converterParams)
		resultFileDocPath := stt_converter.ConvertDataToDoc(result.Data, converterParams)

		log.ResultText = resultText
		log.ResultTextPath = resultTextPath
		log.ResultHtmlPath = resultHtmlPath
		log.ResultFilePdfPath = resultFilePdfPath
		log.ResultFileDocPath = resultFileDocPath
		log.RawResult = result.RawResult

		f.SetConverterLogModel(log)
		log, err = ConverterLogUpdate(f, core.Db)
		if err != nil {
			return types.CheckConverterLog{}, err
		}
		data.ConverterLog = log
		return data, nil

    case "Partial":
        log.Status = result.Status

        resultTextPath, resultText := stt_converter.ConvertDataToText(result.Data, converterParams)
        resultHtmlPath, _ := stt_converter.ConvertDataToHtml(result.Data, converterParams)
        resultFilePdfPath := stt_converter.ConvertDataToPdf(result.Data, converterParams)
        resultFileDocPath := stt_converter.ConvertDataToDoc(result.Data, converterParams)

        log.ResultText = resultText
        log.ResultTextPath = resultTextPath
        log.ResultHtmlPath = resultHtmlPath
        log.ResultFilePdfPath = resultFilePdfPath
        log.ResultFileDocPath = resultFileDocPath
        log.RawResult = result.RawResult
        log.ErrorString = result.ErrorString

        f.SetConverterLogModel(log)
        log, err = ConverterLogUpdate(f, core.Db)
        if err != nil {
            return types.CheckConverterLog{}, err
        }
        data.ConverterLog = log
        data.ErrorString = result.ErrorString
        return data, nil
	}


	//CheckConverterLog Read logic code
	return
}
