package types

import (
    "net/http"
    "stt_back/settings"
)

type ResourceType struct {
    Id   int
    Name string
	//ResourceType remove this line for disable generator functionality
    
    validator
}

func (resourceType *ResourceType) Validate()  {
    //Validate remove this line for disable generator functionality
}

type ResourceTypeFilter struct {
    model ResourceType
    list []ResourceType
    //ResourceTypeFilter remove this line for disable generator functionality

    AbstractFilter
}

func GetResourceTypeFilter(request *http.Request, functionType string) (filter ResourceTypeFilter, err error) {

    filter.request = request
	filter.rawRequestBody, err = GetRawBodyContent(request)
    if err != nil {
        return filter, err
    }
    //filter.TestFilter, _ = strconv.Atoi(request.FormValue("TestFilter"))

    //GetResourceTypeFilter remove this line for disable generator functionality

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


func (filter *ResourceTypeFilter) GetResourceTypeModel() ResourceType {

    filter.model.Validate()

    return  filter.model
}

func (filter *ResourceTypeFilter) GetResourceTypeModelList() (data []ResourceType, err error) {

    for k, _ := range filter.list {
        filter.list[k].Validate()

        if ! filter.list[k].IsValid() {
            err = filter.list[k].GetValidationError()
            break
        }
    }

    return  filter.list, nil
}

func (filter *ResourceTypeFilter) SetResourceTypeModel(typeModel ResourceType) {

    filter.model = typeModel
}
