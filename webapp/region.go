package webapp

import (
    "stt_back/core"
    "stt_back/logic"
    "net/http"
    "stt_back/mdl"
    "stt_back/types"
    "stt_back/settings"
)

    

func RegionFind(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetRegionFilter(httpRequest, settings.FunctionTypeFind)
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
    data, totalRecords, err := logic.RegionFind(requestDto)

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

    
func RegionMultiCreate(w http.ResponseWriter, httpRequest *http.Request) {


    requestDto, err := types.GetRegionFilter(httpRequest, settings.FunctionTypeMultiCreate)
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


    data, err := logic.RegionMultiCreate(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseCreate{
        data,
    })

    return
}

func RegionCreate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetRegionFilter(httpRequest, settings.FunctionTypeCreate)
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


    data, err := logic.RegionCreate(requestDto, core.Db)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseCreate{
        data,
    })

    return
}

    

func RegionRead(w http.ResponseWriter, httpRequest *http.Request) {


    requestDto, err := types.GetRegionFilter(httpRequest, settings.FunctionTypeRead)
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


    data, err := logic.RegionRead(requestDto)

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

    


func RegionMultiUpdate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetRegionFilter(httpRequest, settings.FunctionTypeMultiUpdate)
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


    data, err := logic.RegionMultiUpdate(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseUpdate{
        data,
    })

    return
}

func RegionUpdate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetRegionFilter(httpRequest, settings.FunctionTypeUpdate)
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


    data, err := logic.RegionUpdate(requestDto, core.Db)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseUpdate{
        data,
    })

    return
}

    

func RegionMultiDelete(w http.ResponseWriter, httpRequest *http.Request) {


    requestDto, err := types.GetRegionFilter(httpRequest, settings.FunctionTypeMultiDelete)
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


    isOk, err := logic.RegionMultiDelete(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseDelete{
        isOk,
    })

    return
}

func RegionDelete(w http.ResponseWriter, httpRequest *http.Request) {


    requestDto, err := types.GetRegionFilter(httpRequest, settings.FunctionTypeDelete)
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


    isOk, err := logic.RegionDelete(requestDto, core.Db)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseDelete{
        isOk,
    })

    return
}

    

func RegionFindOrCreate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetRegionFilter(httpRequest, settings.FunctionTypeFindOrCreate)
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


    data, err := logic.RegionFindOrCreate(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseFindOrCreate{
        data,
    })

    return
}

    

func RegionUpdateOrCreate(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetRegionFilter(httpRequest, settings.FunctionTypeUpdateOrCreate)
	if err != nil {
		ErrResponse(w, err, http.StatusBadRequest, requestDto)
		return
	}

    

	if !requestDto.IsValid() {
		Bad(w, requestDto, requestDto.GetValidationError())
		return
	}


    data, err := logic.RegionUpdateOrCreate(requestDto)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseUpdateOrCreate{
        data,
    })

    return
}

