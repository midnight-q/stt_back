package types

import (
	"net/http"
	"strconv"
	"stt_back/settings"
	"time"
)

type ConverterLog struct {
	Id                int
	FilePath          string
	ResultTextPath    string
	ResultFilePath    string
	ResultFormat      string
	RawResult         string
	ResultHtmlPath    string
	ResultFileDocPath string
	ResultFilePdfPath string
	UserId            int
	CreatedAt         time.Time
	SourceFilePath    string
	RecordNumber      int
	TimeFrame         int
	IsShowEmotion     bool
	IsShowSpeaker     bool
	IsShowTag         bool
	IsShowPunctuation bool
	NamedEntityTypes  string
	Status            string
	Token             string
	ResultText        string
	ErrorString       string
	Duration int
	//ConverterLog remove this line for disable generator functionality

	validator
}

func (converterLog *ConverterLog) Validate() {
	//Validate remove this line for disable generator functionality
}

type ConverterLogFilter struct {
	model ConverterLog
	list  []ConverterLog
	//ConverterLogFilter remove this line for disable generator functionality
	UserId int

	AbstractFilter
}

func GetConverterLogFilter(request *http.Request, functionType string) (filter ConverterLogFilter, err error) {

	filter.request = request
	filter.rawRequestBody, err = GetRawBodyContent(request)
	if err != nil {
		return filter, err
	}
	//filter.TestFilter, _ = strconv.Atoi(request.FormValue("TestFilter"))
	filter.UserId, _ = strconv.Atoi(request.FormValue("UserId"))

	//GetConverterLogFilter remove this line for disable generator functionality

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

func (filter *ConverterLogFilter) GetConverterLogModel() ConverterLog {

	filter.model.Validate()

	return filter.model
}

func (filter *ConverterLogFilter) GetConverterLogModelList() (data []ConverterLog, err error) {

	for k, _ := range filter.list {
		filter.list[k].Validate()

		if !filter.list[k].IsValid() {
			err = filter.list[k].GetValidationError()
			break
		}
	}

	return filter.list, nil
}

func (filter *ConverterLogFilter) SetConverterLogModel(typeModel ConverterLog) {

	filter.model = typeModel
}
