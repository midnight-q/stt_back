package types

import (
	"stt_back/settings"
	"stt_back/errors"
)

type APIStatus struct {
	Status string
}

type APIError struct {
	Error  bool
	Errors []Error
}

type Error struct {
	ErrorMessage string
	ErrorCode    int
	Field        string `json:"Field,omitempty"`
	ErrorDebug   string `json:"ErrorDebug,omitempty"`
}

type Pagination struct {
	CurrentPage int
	PerPage     int

	validator
}


func (pagination *Pagination) GetOffset() int {
	return (pagination.CurrentPage - 1) * pagination.PerPage
}

func (pagination *Pagination) Validate(functionType string) {

	switch functionType {

	case settings.FunctionTypeFind:

		if pagination.CurrentPage < 1 {
			pagination.AddValidationError("CurrentPage parameter is not set", errors.ErrorCodeInvalidCurrentPage,"CurrentPage")
		}

		if pagination.PerPage < 1 {
			pagination.AddValidationError("PerPage parameter is not set", errors.ErrorCodeInvalidPerPage,"PerPage")
		}

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
		pagination.validator.AddValidationError("Usupported function type: " + functionType, errors.ErrorCodeUnsupportedFunctionType, "")
		break
	}
}
