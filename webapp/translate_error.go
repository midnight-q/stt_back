package webapp

import (
    "stt_back/core"
    "stt_back/logic"
    "net/http"
    "stt_back/mdl"
    "stt_back/types"
    "stt_back/settings"
)

    

func TranslateErrorFind(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetTranslateErrorFilter(httpRequest, settings.FunctionTypeFind)
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
    data, totalRecords, err := logic.TranslateErrorFind(requestDto)

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

    
func TranslateErrorMultiCreate(w http.ResponseWriter, httpRequest *http.Request) {


    requestDto, err := types.GetTranslateErrorFilter(httpRequest, settings.FunctionTypeMultiCreate)
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


    data, err := logic.TranslateErrorMultiCreate(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseCreate{
        data,
    })

    return
}

func TranslateErrorCreate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetTranslateErrorFilter(httpRequest, settings.FunctionTypeCreate)
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


    data, err := logic.TranslateErrorCreate(requestDto, core.Db)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseCreate{
        data,
    })

    return
}

    

func TranslateErrorRead(w http.ResponseWriter, httpRequest *http.Request) {


    requestDto, err := types.GetTranslateErrorFilter(httpRequest, settings.FunctionTypeRead)
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


    data, err := logic.TranslateErrorRead(requestDto)

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

    


func TranslateErrorMultiUpdate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetTranslateErrorFilter(httpRequest, settings.FunctionTypeMultiUpdate)
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


    data, err := logic.TranslateErrorMultiUpdate(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseUpdate{
        data,
    })

    return
}

func TranslateErrorUpdate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetTranslateErrorFilter(httpRequest, settings.FunctionTypeUpdate)
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


    data, err := logic.TranslateErrorUpdate(requestDto, core.Db)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseUpdate{
        data,
    })

    return
}

    

func TranslateErrorMultiDelete(w http.ResponseWriter, httpRequest *http.Request) {


    requestDto, err := types.GetTranslateErrorFilter(httpRequest, settings.FunctionTypeMultiDelete)
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


    isOk, err := logic.TranslateErrorMultiDelete(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseDelete{
        isOk,
    })

    return
}

func TranslateErrorDelete(w http.ResponseWriter, httpRequest *http.Request) {


    requestDto, err := types.GetTranslateErrorFilter(httpRequest, settings.FunctionTypeDelete)
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


    isOk, err := logic.TranslateErrorDelete(requestDto, core.Db)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseDelete{
        isOk,
    })

    return
}

    

func TranslateErrorFindOrCreate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetTranslateErrorFilter(httpRequest, settings.FunctionTypeFindOrCreate)
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


    data, err := logic.TranslateErrorFindOrCreate(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseFindOrCreate{
        data,
    })

    return
}

    

func TranslateErrorUpdateOrCreate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetTranslateErrorFilter(httpRequest, settings.FunctionTypeUpdateOrCreate)
	if err != nil {
		ErrResponse(w, err, http.StatusBadRequest, requestDto)
		return
	}

    

	if !requestDto.IsValid() {
		Bad(w, requestDto, requestDto.GetValidationError())
		return
	}


    data, err := logic.TranslateErrorUpdateOrCreate(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseUpdateOrCreate{
        data,
    })

    return
}

