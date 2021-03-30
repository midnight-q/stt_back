package types

import (
    "net/http"
    "stt_back/settings"
)

type Resource struct {
    Id   int
    Name string
	Code string
	TypeId int
	//Resource remove this line for disable generator functionality
    
    validator
}

func (resource *Resource) Validate()  {
    //Validate remove this line for disable generator functionality
}

type ResourceFilter struct {
    model Resource
    list []Resource
    //ResourceFilter remove this line for disable generator functionality

    AbstractFilter
}

func GetResourceFilter(request *http.Request, functionType string) (filter ResourceFilter, err error) {

    filter.request = request
	filter.rawRequestBody, err = GetRawBodyContent(request)
    if err != nil {
        return filter, err
    }
    //filter.TestFilter, _ = strconv.Atoi(request.FormValue("TestFilter"))

    //GetResourceFilter remove this line for disable generator functionality

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


func (filter *ResourceFilter) GetResourceModel() Resource {

    filter.model.Validate()

    return  filter.model
}

func (filter *ResourceFilter) GetResourceModelList() (data []Resource, err error) {

    for k, _ := range filter.list {
        filter.list[k].Validate()

        if ! filter.list[k].IsValid() {
            err = filter.list[k].GetValidationError()
            break
        }
    }

    return  filter.list, nil
}

func (filter *ResourceFilter) SetResourceModel(typeModel Resource) {

    filter.model = typeModel
}
