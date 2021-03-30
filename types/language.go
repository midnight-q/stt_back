package types

import (
    "net/http"
    "stt_back/settings"
)

type Language struct {
    Id   int
    Name int
	Code string
	//Language remove this line for disable generator functionality
    
    validator
}

func (language *Language) Validate()  {
    //Validate remove this line for disable generator functionality
}

type LanguageFilter struct {
    model Language
    list []Language
    //LanguageFilter remove this line for disable generator functionality

    AbstractFilter
}

func GetLanguageFilter(request *http.Request, functionType string) (filter LanguageFilter, err error) {

    filter.request = request
	filter.rawRequestBody, err = GetRawBodyContent(request)
    if err != nil {
        return filter, err
    }
    //filter.TestFilter, _ = strconv.Atoi(request.FormValue("TestFilter"))

    //GetLanguageFilter remove this line for disable generator functionality

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


func (filter *LanguageFilter) GetLanguageModel() Language {

    filter.model.Validate()

    return  filter.model
}

func (filter *LanguageFilter) GetLanguageModelList() (data []Language, err error) {

    for k, _ := range filter.list {
        filter.list[k].Validate()

        if ! filter.list[k].IsValid() {
            err = filter.list[k].GetValidationError()
            break
        }
    }

    return  filter.list, nil
}

func (filter *LanguageFilter) SetLanguageModel(typeModel Language) {

    filter.model = typeModel
}
