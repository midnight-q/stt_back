package dbmodels

import (
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
