package webapp

import (
    "stt_back/common"
    "stt_back/errors"
    "stt_back/logic"
    "stt_back/mdl"
    "stt_back/types"
    "net/http"
    "encoding/json"
    "fmt"
)

type FilterInterface interface {
    IsDebug() bool
    GetLanguageId() int
}

func Bad(w http.ResponseWriter, requestDto FilterInterface, err error) {
    ErrResponse(w, err, http.StatusBadRequest, requestDto)
}

func AuthErr(w http.ResponseWriter, requestDto FilterInterface) {

    ErrResponse(w, GetAuthErrTpl(common.MyCaller()), http.StatusForbidden, requestDto)
}

func GetAuthErrTpl(operation string) errors.ErrorWithCode {

    return errors.NewErrorWithCode(
        fmt.Sprintf("Invalid authorize in %s", operation),
        errors.ErrorCodeInvalidAuthorize,
        "Token")
}

func ErrResponse(w http.ResponseWriter, err error, status int, filter FilterInterface) {

    response := types.APIError{}
    response.Error = true

    switch e := err.(type) {
    case errors.ValidatorError:
        for _, errWithCode := range e.Errors() {
            newError := types.Error{
                Field:     errWithCode.GetField(),
                ErrorCode: errWithCode.ErrorCode(),
            }
            if filter.IsDebug() {
                newError.ErrorDebug = errWithCode.Error()
            }
            response.Errors = append(response.Errors, newError)
        }
        break

    case errors.ValidatorErrorInterface:
        for _, errWithCode := range e.Errors() {
            newError := types.Error{
                Field:     errWithCode.GetField(),
                ErrorCode: errWithCode.ErrorCode(),
            }
            if filter.IsDebug() {
                newError.ErrorDebug = errWithCode.Error()
            }
            response.Errors = append(response.Errors, newError)
        }
        break

    case errors.ErrorWithCode:
        newError := types.Error{
            Field:     e.GetField(),
            ErrorCode: e.ErrorCode(),
        }
        if filter.IsDebug() {
            newError.ErrorDebug = e.Error()
        }
        response.Errors = append(response.Errors, newError)
        break

    default:
        newError := types.Error{
        }
        if filter.IsDebug() {
            newError.ErrorDebug = e.Error()
        }
        response.Errors = append(response.Errors, newError)
        break
    }

    var errCodes []int
    for _, e := range response.Errors {
        errCodes = append(errCodes, e.ErrorCode)
    }
    errCodes = common.UniqueIntArray(errCodes)
    f := types.TranslateErrorFilter{}
    f.LanguageId = filter.GetLanguageId()
    if f.LanguageId < 1 {
        f.LanguageId = errors.DefaultErrorLanguageId
    }
    f.ErrorCodes = errCodes
    f.CurrentPage = 1
    f.PerPage = len(errCodes)
    translates, _, err := logic.TranslateErrorFind(f)
    if err != nil {
        fmt.Println("TranslateErrorFind err = ", err)
    }

    for i, err := range response.Errors {
        for _, translate := range translates {
            if err.ErrorCode == translate.Code {
                response.Errors[i].ErrorMessage = translate.Translate
                break
            }
        }
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(response)

    return
}


func ValidResponse (w http.ResponseWriter, data interface{}) {

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    switch data.(type) {
    case mdl.ResponseCreate:
        w.WriteHeader(http.StatusCreated)
        break
    default:
        w.WriteHeader(http.StatusOK)
        break
    }
    json.NewEncoder(w).Encode(data)

    return
}