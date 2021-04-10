package logic

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
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
	file := ioutil.NopCloser(f)
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
