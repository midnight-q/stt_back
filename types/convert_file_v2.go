package types

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"stt_back/errors"
	"stt_back/settings"
)

type ConvertFileV2 struct {
	Id             int
	TimeEstimate   int
	ConverterLogId int
	//ConvertFileV2 remove this line for disable generator functionality

	validator
}

func (convertFileV2 *ConvertFileV2) Validate() {
	//Validate remove this line for disable generator functionality
}

type ConvertFileV2Filter struct {
	model        ConvertFileV2
	list         []ConvertFileV2
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
	//ConvertFileV2Filter remove this line for disable generator functionality

	AbstractFilter
}

func GetConvertFileV2Filter(request *http.Request, functionType string) (filter ConvertFileV2Filter, err error) {

	filter.request = request
	if filter.IsFormDataContentType() {
		filter.DataUrl = request.FormValue("DataUrl")
		filter.File, filter.Header, err = request.FormFile("data")
		if len(filter.DataUrl) < 1 && err != nil {
			fmt.Println("sdsdsd", err)
			err = errors.New("Not found data or DataUrl")
			return ConvertFileV2Filter{}, err
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
			return ConvertFileV2Filter{}, err
		}
	} else {
		err = errors.New("Unsupported content type")
		return ConvertFileV2Filter{}, err
	}

	//GetConvertFileV2Filter remove this line for disable generator functionality

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

func (filter *ConvertFileV2Filter) GetConvertFileV2Model() ConvertFileV2 {

	filter.model.Validate()

	return filter.model
}

func (filter *ConvertFileV2Filter) GetConvertFileV2ModelList() (data []ConvertFileV2, err error) {

	for k, _ := range filter.list {
		filter.list[k].Validate()

		if !filter.list[k].IsValid() {
			err = filter.list[k].GetValidationError()
			break
		}
	}

	return filter.list, nil
}

func (filter *ConvertFileV2Filter) SetConvertFileV2Model(typeModel ConvertFileV2) {

	filter.model = typeModel
}
