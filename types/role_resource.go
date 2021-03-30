package types

import (
    "net/http"
    "stt_back/settings"
)

type RoleResource struct {
    Id   int
    RoleId int
	ResourceId int
	Find bool
	Read bool
	Create bool
	Update bool
	Delete bool
	FindOrCreate bool
	UpdateOrCreate bool
	//RoleResource remove this line for disable generator functionality
    
    validator
}

func (roleResource *RoleResource) Validate()  {
    //Validate remove this line for disable generator functionality
}

type RoleResourceFilter struct {
    model RoleResource
    list []RoleResource
    //RoleResourceFilter remove this line for disable generator functionality

    AbstractFilter
}

func GetRoleResourceFilter(request *http.Request, functionType string) (filter RoleResourceFilter, err error) {

    filter.request = request
	filter.rawRequestBody, err = GetRawBodyContent(request)
    if err != nil {
        return filter, err
    }
    //filter.TestFilter, _ = strconv.Atoi(request.FormValue("TestFilter"))

    //GetRoleResourceFilter remove this line for disable generator functionality

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


func (filter *RoleResourceFilter) GetRoleResourceModel() RoleResource {

    filter.model.Validate()

    return  filter.model
}

func (filter *RoleResourceFilter) GetRoleResourceModelList() (data []RoleResource, err error) {

    for k, _ := range filter.list {
        filter.list[k].Validate()

        if ! filter.list[k].IsValid() {
            err = filter.list[k].GetValidationError()
            break
        }
    }

    return  filter.list, nil
}

func (filter *RoleResourceFilter) SetRoleResourceModel(typeModel RoleResource) {

    filter.model = typeModel
}
