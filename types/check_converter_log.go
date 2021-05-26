package types

import (
	"net/http"
	"stt_back/settings"
)

type CheckConverterLog struct {
	Id           int
	ConverterLog ConverterLog
	TimeEstimate int
	ErrorString  string
	Status       string

	//CheckConverterLog remove this line for disable generator functionality

	validator
}

func (checkConverterLog *CheckConverterLog) Validate() {
	//Validate remove this line for disable generator functionality
}

type CheckConverterLogFilter struct {
	model CheckConverterLog
	list  []CheckConverterLog
	//CheckConverterLogFilter remove this line for disable generator functionality

	AbstractFilter
}

func GetCheckConverterLogFilter(request *http.Request, functionType string) (filter CheckConverterLogFilter, err error) {

	filter.request = request
	filter.rawRequestBody, err = GetRawBodyContent(request)
	if err != nil {
		return filter, err
	}
	//filter.TestFilter, _ = strconv.Atoi(request.FormValue("TestFilter"))

	//GetCheckConverterLogFilter remove this line for disable generator functionality

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

func (filter *CheckConverterLogFilter) GetCheckConverterLogModel() CheckConverterLog {

	filter.model.Validate()

	return filter.model
}

func (filter *CheckConverterLogFilter) GetCheckConverterLogModelList() (data []CheckConverterLog, err error) {

	for k, _ := range filter.list {
		filter.list[k].Validate()

		if !filter.list[k].IsValid() {
			err = filter.list[k].GetValidationError()
			break
		}
	}

	return filter.list, nil
}

func (filter *CheckConverterLogFilter) SetCheckConverterLogModel(typeModel CheckConverterLog) {

	filter.model = typeModel
}
