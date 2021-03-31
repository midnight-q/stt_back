package logic

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/wav"
	"github.com/jinzhu/gorm"
	"github.com/orcaman/writerseeker"
	"io/ioutil"
	"mime/multipart"
	"stt_back/common"
	"stt_back/core"
	"stt_back/errors"
	"stt_back/services/file_storage"
	"stt_back/services/stt_converter"
	"stt_back/types"
)

func ConvertFileCreate(filter types.ConvertFileFilter, query *gorm.DB) (data types.ConvertFile, err error) {
	inputFile, err := convertInputFile(filter.File, filter.Header)
	if err != nil {
		return types.ConvertFile{}, err
	}

	result, err := stt_converter.ConvertSpeechToText(inputFile)
	if err != nil {
		return types.ConvertFile{}, err
	}

	resultFilePath := createResultFileIfNeed(result, filter.ResultFormat)

	filePath, err := file_storage.CreateFileInLocalStorage(inputFile, ".wav")
	if err != nil {
		return types.ConvertFile{}, err
	}

	f := types.ConverterLogFilter{}
	f.SetConverterLogModel(types.ConverterLog{
		FilePath:       filePath,
		ResultText:     result.Text,
		ResultFilePath: resultFilePath,
		ResultFormat:   filter.ResultFormat,
		RawResult:      result.RawResult,
	})
	_, err = ConverterLogCreate(f, core.Db)

	data = types.ConvertFile{
		ResultFilePath: resultFilePath,
		ResultText:     result.Text,
		ResultFormat:   filter.ResultFormat,
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
			Precision:   1,
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
			Precision:   1,
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

func createResultFileIfNeed(result stt_converter.Result, format string) (path string){
	//TODO: Implement creating file for format
	switch format {
	case "pdf":
		path = "pdf_path"
		break
	case "doc":
		path = "doc_path"
		break
	}

	return
}