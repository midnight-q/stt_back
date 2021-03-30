package types

import (
    "net/http"
    "stt_back/settings"
)

type UserRole struct {
    Id   int
    UserId int
	RoleId int
	//UserRole remove this line for disable generator functionality
    
    validator
}

func (userRole *UserRole) Validate()  {
    //Validate remove this line for disable generator functionality
}

type UserRoleFilter struct {
    model UserRole
    list []UserRole
    //UserRoleFilter remove this line for disable generator functionality

    AbstractFilter
}

func GetUserRoleFilter(request *http.Request, functionType string) (filter UserRoleFilter, err error) {

    filter.request = request
	filter.rawRequestBody, err = GetRawBodyContent(request)
    if err != nil {
        return filter, err
    }
    //filter.TestFilter, _ = strconv.Atoi(request.FormValue("TestFilter"))

    //GetUserRoleFilter remove this line for disable generator functionality

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


func (filter *UserRoleFilter) GetUserRoleModel() UserRole {

    filter.model.Validate()

    return  filter.model
}

func (filter *UserRoleFilter) GetUserRoleModelList() (data []UserRole, err error) {

    for k, _ := range filter.list {
        filter.list[k].Validate()

        if ! filter.list[k].IsValid() {
            err = filter.list[k].GetValidationError()
            break
        }
    }

    return  filter.list, nil
}

func (filter *UserRoleFilter) SetUserRoleModel(typeModel UserRole) {

    filter.model = typeModel
}
