package types

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"stt_back/errors"
	"stt_back/services/stt_converter"
	"stt_back/settings"
)

type ConvertFile struct {
	Id                int
	Data              []stt_converter.Data
	ResultHtmlPath    string
	ResultTextPath    string
	ResultFilePdfPath string
	ResultFileDocPath string
	SourceFilePath    string
	//ConvertFile remove this line for disable generator functionality

	validator
}

func (convertFile *ConvertFile) Validate() {
	//Validate remove this line for disable generator functionality
}

type ConvertFileFilter struct {
	model ConvertFile
	list  []ConvertFile

	File         multipart.File
	Header       *multipart.FileHeader
	FileSource   multipart.File
	HeaderSource *multipart.FileHeader

	TimeFrame         int
	IsShowEmotion     bool
	IsShowSpeaker     bool
	IsShowTag         bool
	IsShowPunctuation bool
	UserId            int
	DataUrl           string
	NamedEntityTypes  []string
	//ConvertFileFilter remove this line for disable generator functionality

	AbstractFilter
}

func GetConvertFileFilter(request *http.Request, functionType string) (filter ConvertFileFilter, err error) {

	filter.request = request
	if filter.IsFormDataContentType() {
		filter.DataUrl = request.FormValue("DataUrl")
		filter.File, filter.Header, err = request.FormFile("data")
		if len(filter.DataUrl) < 1 && err != nil {
			fmt.Println("sdsdsd", err)
			err = errors.New("Not found data or DataUrl")
			return ConvertFileFilter{}, err
		}
		filter.FileSource, filter.HeaderSource, _ = request.FormFile("dataSource")
		filter.TimeFrame, _ = strconv.Atoi(request.FormValue("TimeFrame"))
		filter.UserId, _ = strconv.Atoi(request.FormValue("UserId"))
		filter.IsShowEmotion, _ = strconv.ParseBool(request.FormValue("IsShowEmotion"))
		filter.IsShowSpeaker, _ = strconv.ParseBool(request.FormValue("IsShowSpeaker"))
		filter.IsShowTag, _ = strconv.ParseBool(request.FormValue("IsShowTag"))
		filter.IsShowPunctuation, _ = strconv.ParseBool(request.FormValue("IsShowPunctuation"))

		NamedEntityTypesString := request.FormValue("NamedEntityTypes")
		err = json.Unmarshal([]byte(NamedEntityTypesString), &filter.NamedEntityTypes)
		if err != nil {
			fmt.Println("Unmarshal error:", err)
			return ConvertFileFilter{}, err
		}
	} else {
		err = errors.New("Unsupported content type")
		return ConvertFileFilter{}, err
	}

	//filter.TestFilter, _ = strconv.Atoi(request.FormValue("TestFilter"))

	//GetConvertFileFilter remove this line for disable generator functionality

	switch functionType {
	case settings.FunctionTypeMultiCreate, settings.FunctionTypeMultiUpdate, settings.FunctionTypeMultiDelete, settings.FunctionTypeMultiFindOrCreate:
		err = ReadJSON(filter.rawRequestBody, &filter.list)
		if err != nil {
			return
		}
		break
	default:
		err = ReadJSON(filter.rawRequestBody, &filter.model)
		if err != nil {
			return
		}
		break
	}

	filter.AbstractFilter, err = GetAbstractFilter(request, filter.rawRequestBody, functionType)

	return filter, err
}

func (filter *ConvertFileFilter) GetConvertFileModel() ConvertFile {

	filter.model.Validate()

	return filter.model
}

func (filter *ConvertFileFilter) GetConvertFileModelList() (data []ConvertFile, err error) {

	for k, _ := range filter.list {
		filter.list[k].Validate()

		if !filter.list[k].IsValid() {
			err = filter.list[k].GetValidationError()
			break
		}
	}

	return filter.list, nil
}

func (filter *ConvertFileFilter) SetConvertFileModel(typeModel ConvertFile) {

	filter.model = typeModel
}
