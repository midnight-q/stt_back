package webapp

import (
    "stt_back/core"
    "stt_back/logic"
    "net/http"
    "stt_back/mdl"
    "stt_back/types"
    "stt_back/settings"
)

    

func RoleResourceFind(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetRoleResourceFilter(httpRequest, settings.FunctionTypeFind)
	if err != nil {
		ErrResponse(w, err, http.StatusBadRequest, requestDto)
		return
	}

    if !requestDto.IsAuthorized() {
		AuthErr(w, requestDto)
		return
	}

	if !requestDto.IsValid() {
		Bad(w, requestDto, requestDto.GetValidationError())
		return
	}

    
    // Получаем список
    data, totalRecords, err := logic.RoleResourceFind(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseFind{
        data,
        totalRecords,
    })

    return
}

    
func RoleResourceMultiCreate(w http.ResponseWriter, httpRequest *http.Request) {


    requestDto, err := types.GetRoleResourceFilter(httpRequest, settings.FunctionTypeMultiCreate)
	if err != nil {
		ErrResponse(w, err, http.StatusBadRequest, requestDto)
		return
	}

    if !requestDto.IsAuthorized() {
		AuthErr(w, requestDto)
		return
	}

	if !requestDto.IsValid() {
		Bad(w, requestDto, requestDto.GetValidationError())
		return
	}


    data, err := logic.RoleResourceMultiCreate(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseCreate{
        data,
    })

    return
}

func RoleResourceCreate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetRoleResourceFilter(httpRequest, settings.FunctionTypeCreate)
	if err != nil {
		ErrResponse(w, err, http.StatusBadRequest, requestDto)
		return
	}

    if !requestDto.IsAuthorized() {
		AuthErr(w, requestDto)
		return
	}

	if !requestDto.IsValid() {
		Bad(w, requestDto, requestDto.GetValidationError())
		return
	}


    data, err := logic.RoleResourceCreate(requestDto, core.Db)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseCreate{
        data,
    })

    return
}

    

func RoleResourceRead(w http.ResponseWriter, httpRequest *http.Request) {


    requestDto, err := types.GetRoleResourceFilter(httpRequest, settings.FunctionTypeRead)
	if err != nil {
		ErrResponse(w, err, http.StatusBadRequest, requestDto)
		return
	}

    if !requestDto.IsAuthorized() {
		AuthErr(w, requestDto)
		return
	}

    requestDto.PerPage = 1
    requestDto.CurrentPage = 1

	if !requestDto.IsValid() {
		Bad(w, requestDto, requestDto.GetValidationError())
		return
	}


    data, err := logic.RoleResourceRead(requestDto)

    // Создаём структуру ответа
    if err != nil {
        code := http.StatusBadRequest
        if err.Error() == "Not found" {
            code = http.StatusNotFound
        }
        ErrResponse(w, err, code, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseRead{
        data,
    })

    return
}

    


func RoleResourceMultiUpdate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetRoleResourceFilter(httpRequest, settings.FunctionTypeMultiUpdate)
	if err != nil {
		ErrResponse(w, err, http.StatusBadRequest, requestDto)
		return
	}

    if !requestDto.IsAuthorized() {
		AuthErr(w, requestDto)
		return
	}

	if !requestDto.IsValid() {
		Bad(w, requestDto, requestDto.GetValidationError())
		return
	}


    data, err := logic.RoleResourceMultiUpdate(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseUpdate{
        data,
    })

    return
}

func RoleResourceUpdate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetRoleResourceFilter(httpRequest, settings.FunctionTypeUpdate)
	if err != nil {
		ErrResponse(w, err, http.StatusBadRequest, requestDto)
		return
	}

    if !requestDto.IsAuthorized() {
		AuthErr(w, requestDto)
		return
	}

	if !requestDto.IsValid() {
		Bad(w, requestDto, requestDto.GetValidationError())
		return
	}


    data, err := logic.RoleResourceUpdate(requestDto, core.Db)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseUpdate{
        data,
    })

    return
}

    

func RoleResourceMultiDelete(w http.ResponseWriter, httpRequest *http.Request) {


    requestDto, err := types.GetRoleResourceFilter(httpRequest, settings.FunctionTypeMultiDelete)
	if err != nil {
		ErrResponse(w, err, http.StatusBadRequest, requestDto)
		return
	}

    if !requestDto.IsAuthorized() {
		AuthErr(w, requestDto)
		return
	}

	if !requestDto.IsValid() {
		Bad(w, requestDto, requestDto.GetValidationError())
		return
	}


    isOk, err := logic.RoleResourceMultiDelete(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseDelete{
        isOk,
    })

    return
}

func RoleResourceDelete(w http.ResponseWriter, httpRequest *http.Request) {


    requestDto, err := types.GetRoleResourceFilter(httpRequest, settings.FunctionTypeDelete)
	if err != nil {
		ErrResponse(w, err, http.StatusBadRequest, requestDto)
		return
	}

    if !requestDto.IsAuthorized() {
		AuthErr(w, requestDto)
		return
	}

	if !requestDto.IsValid() {
		Bad(w, requestDto, requestDto.GetValidationError())
		return
	}


    isOk, err := logic.RoleResourceDelete(requestDto, core.Db)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseDelete{
        isOk,
    })

    return
}

    

func RoleResourceFindOrCreate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetRoleResourceFilter(httpRequest, settings.FunctionTypeFindOrCreate)
	if err != nil {
		ErrResponse(w, err, http.StatusBadRequest, requestDto)
		return
	}

    if !requestDto.IsAuthorized() {
		AuthErr(w, requestDto)
		return
	}

	if !requestDto.IsValid() {
		Bad(w, requestDto, requestDto.GetValidationError())
		return
	}


    data, err := logic.RoleResourceFindOrCreate(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseFindOrCreate{
        data,
    })

    return
}

    

func RoleResourceUpdateOrCreate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetRoleResourceFilter(httpRequest, settings.FunctionTypeUpdateOrCreate)
	if err != nil {
		ErrResponse(w, err, http.StatusBadRequest, requestDto)
		return
	}

    if !requestDto.IsAuthorized() {
		AuthErr(w, requestDto)
		return
	}

	if !requestDto.IsValid() {
		Bad(w, requestDto, requestDto.GetValidationError())
		return
	}


    data, err := logic.RoleResourceUpdateOrCreate(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseUpdateOrCreate{
        data,
    })

    return
}

