package types

import (
    "stt_back/settings"
    "stt_back/errors"
)

type validator struct {
    validationError errors.ValidatorError
}

func (val *validator) IsValid() bool {

    return val.validationError.IsEmpty()
}

func (val *validator) GetValidationError() errors.ValidatorErrorInterface {
    return &val.validationError
}

func (val *validator) AddValidationError(err string, code errors.ErrorCode, field string) {
    val.validationError.AddError(errors.NewErrorWithCode(err, code, field))
}

func (val *validator) Validate(functionType string) {

    switch functionType {

    case settings.FunctionTypeFind:
        break

    case settings.FunctionTypeCreate:
        break

    case settings.FunctionTypeMultiCreate:
        break

    case settings.FunctionTypeRead:
        break

    case settings.FunctionTypeUpdate:
        break

    case settings.FunctionTypeMultiUpdate:
        break

    case settings.FunctionTypeDelete:
        break

    case settings.FunctionTypeMultiDelete:
        break

    default:
        val.AddValidationError("Unsupported function type: "+functionType, errors.ErrorCodeUnsupportedFunctionType, "")
        break
    }
}
