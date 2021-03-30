package webapp

import (
    "stt_back/core"
    "stt_back/logic"
    "net/http"
    "stt_back/mdl"
    "stt_back/types"
    "stt_back/settings"
)

    

func ResourceTypeFind(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetResourceTypeFilter(httpRequest, settings.FunctionTypeFind)
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
    data, totalRecords, err := logic.ResourceTypeFind(requestDto)

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

    
func ResourceTypeMultiCreate(w http.ResponseWriter, httpRequest *http.Request) {


    requestDto, err := types.GetResourceTypeFilter(httpRequest, settings.FunctionTypeMultiCreate)
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


    data, err := logic.ResourceTypeMultiCreate(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseCreate{
        data,
    })

    return
}

func ResourceTypeCreate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetResourceTypeFilter(httpRequest, settings.FunctionTypeCreate)
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


    data, err := logic.ResourceTypeCreate(requestDto, core.Db)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseCreate{
        data,
    })

    return
}

    

func ResourceTypeRead(w http.ResponseWriter, httpRequest *http.Request) {


    requestDto, err := types.GetResourceTypeFilter(httpRequest, settings.FunctionTypeRead)
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


    data, err := logic.ResourceTypeRead(requestDto)

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

    


func ResourceTypeMultiUpdate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetResourceTypeFilter(httpRequest, settings.FunctionTypeMultiUpdate)
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


    data, err := logic.ResourceTypeMultiUpdate(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseUpdate{
        data,
    })

    return
}

func ResourceTypeUpdate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetResourceTypeFilter(httpRequest, settings.FunctionTypeUpdate)
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


    data, err := logic.ResourceTypeUpdate(requestDto, core.Db)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseUpdate{
        data,
    })

    return
}

    

func ResourceTypeMultiDelete(w http.ResponseWriter, httpRequest *http.Request) {


    requestDto, err := types.GetResourceTypeFilter(httpRequest, settings.FunctionTypeMultiDelete)
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


    isOk, err := logic.ResourceTypeMultiDelete(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseDelete{
        isOk,
    })

    return
}

func ResourceTypeDelete(w http.ResponseWriter, httpRequest *http.Request) {


    requestDto, err := types.GetResourceTypeFilter(httpRequest, settings.FunctionTypeDelete)
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


    isOk, err := logic.ResourceTypeDelete(requestDto, core.Db)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseDelete{
        isOk,
    })

    return
}

    

func ResourceTypeFindOrCreate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetResourceTypeFilter(httpRequest, settings.FunctionTypeFindOrCreate)
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


    data, err := logic.ResourceTypeFindOrCreate(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseFindOrCreate{
        data,
    })

    return
}

    

func ResourceTypeUpdateOrCreate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetResourceTypeFilter(httpRequest, settings.FunctionTypeUpdateOrCreate)
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


    data, err := logic.ResourceTypeUpdateOrCreate(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseUpdateOrCreate{
        data,
    })

    return
}

