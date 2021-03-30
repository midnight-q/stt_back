package types

import (
    "net/http"
    "stt_back/settings"
)

type CurrentUser struct {
    Id   int
    //CurrentUser remove this line for disable generator functionality
    
    validator
}

func (currentUser *CurrentUser) Validate()  {
    //Validate remove this line for disable generator functionality
}

type CurrentUserFilter struct {
    model CurrentUser
    list []CurrentUser
    //CurrentUserFilter remove this line for disable generator functionality

    AbstractFilter
}

func GetCurrentUserFilter(request *http.Request, functionType string) (filter CurrentUserFilter, err error) {

    filter.request = request
	filter.rawRequestBody, err = GetRawBodyContent(request)
    if err != nil {
        return filter, err
    }
    //filter.TestFilter, _ = strconv.Atoi(request.FormValue("TestFilter"))

    //GetCurrentUserFilter remove this line for disable generator functionality

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


func (filter *CurrentUserFilter) GetCurrentUserModel() CurrentUser {

    filter.model.Validate()

    return  filter.model
}

func (filter *CurrentUserFilter) GetCurrentUserModelList() (data []CurrentUser, err error) {

    for k, _ := range filter.list {
        filter.list[k].Validate()

        if ! filter.list[k].IsValid() {
            err = filter.list[k].GetValidationError()
            break
        }
    }

    return  filter.list, nil
}

func (filter *CurrentUserFilter) SetCurrentUserModel(typeModel CurrentUser) {

    filter.model = typeModel
}
