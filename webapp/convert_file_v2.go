package webapp

import (
	"net/http"
	"stt_back/core"
	"stt_back/logic"
	"stt_back/mdl"
	"stt_back/settings"
	"stt_back/types"
)

func ConvertFileV2Create(w http.ResponseWriter, httpRequest *http.Request) {

    requestDto, err := types.GetConvertFileV2Filter(httpRequest, settings.FunctionTypeCreate)
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


    data, err := logic.ConvertFileV2Create(requestDto, core.Db)

    if err != nil {
        ErrResponse(w, err, http.StatusBadRequest, requestDto)
        return
    }

    ValidResponse(w, mdl.ResponseCreate{
        data,
    })

    return
}
