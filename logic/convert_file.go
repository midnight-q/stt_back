package logic

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"stt_back/common"
	"stt_back/core"
	"stt_back/dbmodels"
	"stt_back/errors"
	"stt_back/services/file_storage"
	"stt_back/services/stt_converter"
	"stt_back/types"

	"github.com/jinzhu/gorm"
)

func ConvertFileCreate(filter types.ConvertFileFilter, query *gorm.DB) (data types.ConvertFile, err error) {
	if filter.TimeFrame < 1 {
		err = errors.New("Not found TimeFrame")
		return types.ConvertFile{}, err
	}
	if filter.UserId < 1 {
		err = errors.New("Not found UserId")
		return types.ConvertFile{}, err
	}

	var inputFile []byte
	if filter.Header != nil {
		ext := filepath.Ext(filter.Header.Filename)
		b, err := ioutil.ReadAll(filter.File)
		if err != nil {
			return types.ConvertFile{}, err
		}
		inputFile, err = convertInputFile(bytes.NewReader(b), ext)
		if err != nil {
			return types.ConvertFile{}, err
		}
	} else if len(filter.DataUrl) > 0 {
		b, ext, err := loadFileFormUrl(filter.DataUrl)
		if err != nil {
			return types.ConvertFile{}, err
		}
		inputFile, err = convertInputFile(bytes.NewReader(b), ext)
		if err != nil {
			return types.ConvertFile{}, err
		}
	} else {
		err = errors.New("Not found file")
		return types.ConvertFile{}, err
	}
	converterParams := stt_converter.Params{
		TimeFrame:         filter.TimeFrame,
		IsShowEmotion:     filter.IsShowEmotion,
		IsShowSpeaker:     filter.IsShowSpeaker,
		IsShowTag:         filter.IsShowTag,
		IsShowPunctuation: filter.IsShowPunctuation,
		NamedEntityTypes:  filter.NamedEntityTypes,
	}

	result, err := stt_converter.ConvertSpeechToText(inputFile, converterParams)
	if err != nil {
		return types.ConvertFile{}, err
	}

	resultTextPath, resultText := stt_converter.ConvertDataToText(result.Data, converterParams)
	resultHtmlPath, _ := stt_converter.ConvertDataToHtml(result.Data, converterParams)
	resultFilePdfPath := stt_converter.ConvertDataToPdf(result.Data, converterParams)
	resultFileDocPath := stt_converter.ConvertDataToDoc(result.Data, converterParams)

	filePath, err := file_storage.CreateFileInLocalStorage(inputFile, ".wav")
	if err != nil {
		return types.ConvertFile{}, err
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

	f := types.ConverterLogFilter{}
	f.SetConverterLogModel(types.ConverterLog{
		ResultText:        resultText,
		FilePath:          filePath,
		ResultTextPath:    resultTextPath,
		ResultHtmlPath:    resultHtmlPath,
		ResultFileDocPath: resultFileDocPath,
		ResultFilePdfPath: resultFilePdfPath,
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
		Status:            "Complete",
	})

	_, err = ConverterLogCreate(f, core.Db)

	data = types.ConvertFile{
		ResultTextPath:    resultTextPath,
		ResultHtmlPath:    resultHtmlPath,
		ResultFileDocPath: resultFileDocPath,
		ResultFilePdfPath: resultFilePdfPath,
		Data:              result.Data,
		SourceFilePath:    sourceFilePath,
		ResultText:        resultText,
	}
	return
}

func getNumberForUser(userId int) int {
	count := 0
	core.Db.Unscoped().Model(dbmodels.ConverterLog{}).Where(dbmodels.ConverterLog{UserId: userId}).Count(&count)

	return count + 1
}

func loadFileFormUrl(url string) (res []byte, ext string, err error) {
	ext = filepath.Ext(url)
	if len(ext) < 4 {
		err = errors.New("Unsupported file format")
		return nil, "", err
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, "", getUnsupportedFileFormatError()
	}
	res, err = ioutil.ReadAll(resp.Body)
	return
}

func getUnsupportedFileFormatError() error {
	return errors.NewErrorWithCode("Unsupported file format", errors.ErrorCodeUnsupportedFileFormat, "InputFilename")
}

func convertInputFile(f *bytes.Reader, fileFormat string) (res []byte, err error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	name := uuid.New().String()
	inputName := name + fileFormat
	outputName := name + "_new.wav"
	err = ioutil.WriteFile(inputName, b, os.ModePerm)
	if err != nil {
		return
	}
	cmd := exec.Command("/usr/bin/ffmpeg", "-i", inputName, "-ac", "1", "-ar", "8000", outputName)
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(cmd.String())
		return nil, err
	}

	b, err = ioutil.ReadFile(outputName)
	if err != nil {
		return
	}
	err = os.Remove(inputName)
	if err != nil {
		return
	}
	err = os.Remove(outputName)
	if err != nil {
		return
	}

	return b, nil
}
