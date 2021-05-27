package logic

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"stt_back/common"
	"stt_back/core"
	"stt_back/errors"
	"stt_back/services/file_storage"
	"stt_back/services/stt_converter"
	"stt_back/types"

	"github.com/jinzhu/gorm"
)

func ConvertFileV2Create(filter types.ConvertFileV2Filter, query *gorm.DB) (data types.ConvertFileV2, err error) {

	if filter.TimeFrame < 1 {
		err = errors.New("Not found TimeFrame")
		return types.ConvertFileV2{}, err
	}
	if filter.UserId < 1 {
		err = errors.New("Not found UserId")
		return types.ConvertFileV2{}, err
	}

	var inputFile []byte
	var duration float64
	if filter.Header != nil {
		ext := filepath.Ext(filter.Header.Filename)
		b, err := ioutil.ReadAll(filter.File)
		if err != nil {
			return types.ConvertFileV2{}, err
		}
		inputFile, duration, err = convertInputFile(bytes.NewReader(b), ext)
		if err != nil {
			return types.ConvertFileV2{}, err
		}
	} else if len(filter.DataUrl) > 0 {
		b, ext, err := loadFileFormUrl(filter.DataUrl)
		if err != nil {
			return types.ConvertFileV2{}, err
		}
		inputFile, duration, err = convertInputFile(bytes.NewReader(b), ext)
		if err != nil {
			return types.ConvertFileV2{}, err
		}
	} else {
		err = errors.New("Not found file")
		return types.ConvertFileV2{}, err
	}

	result, err := stt_converter.ConvertSpeechToTextV2(inputFile)
	if err != nil {
		return types.ConvertFileV2{}, err
	}

	sourceFilePath := ""
	if filter.HeaderSource != nil {
		fileExt := filepath.Ext(filter.HeaderSource.Filename)
		data, e := ioutil.ReadAll(filter.FileSource)
		if e == nil {
			path, e := file_storage.CreateFileInLocalStorage(data, fileExt)
			if e != nil {
				fmt.Println("Error save source file:", e)
			}
			sourceFilePath = path
		}
	}

	filePath, err := file_storage.CreateFileInLocalStorage(inputFile, ".wav")
	if err != nil {
		return types.ConvertFileV2{}, err
	}

	f := types.ConverterLogFilter{}
	f.SetConverterLogModel(types.ConverterLog{
		FilePath:          filePath,
		RawResult:         result.RawResult,
		UserId:            filter.UserId,
		SourceFilePath:    sourceFilePath,
		RecordNumber:      getNumberForUser(filter.UserId),
		TimeFrame:         filter.TimeFrame,
		IsShowEmotion:     filter.IsShowEmotion,
		IsShowSpeaker:     filter.IsShowSpeaker,
		IsShowTag:         filter.IsShowTag,
		IsShowPunctuation: filter.IsShowPunctuation,
		NamedEntityTypes:  common.DataToString(filter.NamedEntityTypes),
		Status:            "Processing",
		Token:             result.Token,
		Duration:          int(duration),
	})

	converterLog, err := ConverterLogCreate(f, core.Db)
	if err != nil {
		return
	}

	data.ConverterLogId = converterLog.Id

	//ConvertFileV2 Create logic code
	return
}
