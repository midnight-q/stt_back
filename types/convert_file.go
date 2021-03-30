package types

import (
    "net/http"
    "stt_back/settings"
)

type ConvertFile struct {
    Id   int
    ResultFilePath string
	ResultText string
	InputData []byte
	InputFilename string
	ResultFormat string
	//ConvertFile remove this line for disable generator functionality
    
    validator
}

func (convertFile *ConvertFile) Validate()  {
    //Validate remove this line for disable generator functionality
}

type ConvertFileFilter struct {
    model ConvertFile
    list []ConvertFile
    //ConvertFileFilter remove this line for disable generator functionality

    AbstractFilter
}

func GetConvertFileFilter(request *http.Request, functionType string) (filter ConvertFileFilter, err error) {

    filter.request = request
	filter.rawRequestBody, err = GetRawBodyContent(request)
    if err != nil {
        return filter, err
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

    return  filter, err
}


func (filter *ConvertFileFilter) GetConvertFileModel() ConvertFile {

    filter.model.Validate()

    return  filter.model
}

func (filter *ConvertFileFilter) GetConvertFileModelList() (data []ConvertFile, err error) {

    for k, _ := range filter.list {
        filter.list[k].Validate()

        if ! filter.list[k].IsValid() {
            err = filter.list[k].GetValidationError()
            break
        }
    }

    return  filter.list, nil
}

func (filter *ConvertFileFilter) SetConvertFileModel(typeModel ConvertFile) {

    filter.model = typeModel
}
