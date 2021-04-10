package logic

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"path/filepath"
	"stt_back/common"
	"stt_back/core"
	"stt_back/errors"
	"stt_back/services/file_storage"
	"stt_back/services/stt_converter"
	"stt_back/types"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/wav"
	"github.com/jinzhu/gorm"
	"github.com/orcaman/writerseeker"
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

	inputFile, err := convertInputFile(filter.File, filter.Header)
	if err != nil {
		return types.ConvertFile{}, err
	}

	result, err := stt_converter.ConvertSpeechToText(inputFile)
	if err != nil {
		return types.ConvertFile{}, err
	}

	converterParams := stt_converter.Params{
		TimeFrame:         filter.TimeFrame,
		IsShowEmotion:     filter.IsShowEmotion,
		IsShowSpeaker:     filter.IsShowSpeaker,
		IsShowTag:         filter.IsShowTag,
		IsShowPunctuation: filter.IsShowPunctuation,
	}
	resultTextPath := stt_converter.ConvertDataToText(result.Data, converterParams)
	resultHtmlPath := stt_converter.ConvertDataToHtml(result.Data, converterParams)
	resultFilePdfPath := stt_converter.ConvertDataToPdf(result.Data, converterParams)
	resultFileDocPath := stt_converter.ConvertDataToDoc(result.Data, converterParams) // TODO: Implement this

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
		FilePath:          filePath,
		ResultTextPath:    resultTextPath,
		ResultHtmlPath:    resultHtmlPath,
		ResultFileDocPath: resultFileDocPath,
		ResultFilePdfPath: resultFilePdfPath,
		RawResult:         result.RawResult,
		UserId:            filter.UserId,
		SourceFilePath:    sourceFilePath,
	})
	_, err = ConverterLogCreate(f, core.Db)

	data = types.ConvertFile{
		ResultTextPath:    resultTextPath,
		ResultHtmlPath:    resultHtmlPath,
		ResultFileDocPath: resultFileDocPath,
		ResultFilePdfPath: resultFilePdfPath,
		Data:              result.Data,
		SourceFilePath:    sourceFilePath,
	}
	return
}

func getUnsupportedFileFormatError() error {
	return errors.NewErrorWithCode("Unsupported file format", errors.ErrorCodeUnsupportedFileFormat, "InputFilename")
}

func convertInputFile(file multipart.File, header *multipart.FileHeader) (res []byte, err error) {
	fileFormat := common.GetFileFormatFromName(header.Filename)

	switch fileFormat {
	case ".wav":
		streamer, format, err := wav.Decode(file)
		if err != nil {
			return []byte{}, err
		}
		defer streamer.Close()
		r := beep.Resample(3, format.SampleRate, beep.SampleRate(8000), streamer)
		buf := writerseeker.WriterSeeker{}
		err = wav.Encode(&buf, r, beep.Format{
			SampleRate:  8000,
			NumChannels: 1,
			Precision:   2,
		})
		if err != nil {
			return []byte{}, err
		}
		res, err = ioutil.ReadAll(buf.Reader())
		if err != nil {
			return []byte{}, err
		}

		break
	case ".mp3":
		streamer, format, err := mp3.Decode(file)
		if err != nil {
			return []byte{}, err
		}
		defer streamer.Close()
		r := beep.Resample(3, format.SampleRate, beep.SampleRate(8000), streamer)
		buf := writerseeker.WriterSeeker{}
		err = wav.Encode(&buf, r, beep.Format{
			SampleRate:  8000,
			NumChannels: 1,
			Precision:   2,
		})
		if err != nil {
			return []byte{}, err
		}
		res, err = ioutil.ReadAll(buf.Reader())
		if err != nil {
			return []byte{}, err
		}
		break
	default:
		return []byte{}, getUnsupportedFileFormatError()
	}

	return
}
