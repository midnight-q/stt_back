package types

import (
    "stt_back/settings"
    "stt_back/errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type GoshaFilterIds struct {
	Ids       []int
	ExceptIds []int
	currentId int

	validator
}

func (filter *GoshaFilterIds) GetFirstId() (int, error) {
	for _, id := range filter.Ids {
		return id, nil
	}
	return 0, errors.New("Empty array")
}

func (filter *GoshaFilterIds) GetCurrentId() int {
	return filter.currentId
}

func (filter *GoshaFilterIds) SetCurrentId(id int) int {
	filter.currentId = id
	return filter.currentId
}

func (filter *GoshaFilterIds) GetIds() []int {
	return filter.Ids
}

func (filter *GoshaFilterIds) GetExceptIds() []int {
	return filter.ExceptIds
}

func (filter *GoshaFilterIds) AddId(id int) *GoshaFilterIds {
	filter.Ids = append(filter.Ids, id)
	return filter
}

func (filter *GoshaFilterIds) AddExceptIds(id int) *GoshaFilterIds {
	filter.ExceptIds = append(filter.ExceptIds, id)
	return filter
}

func (filter *GoshaFilterIds) AddIds(ids []int) *GoshaFilterIds {
	for _, id := range ids {
		filter.AddId(id)
	}
	return filter
}

func (filter *GoshaFilterIds) Clear() *GoshaFilterIds {
	filter.Ids = []int{}
	return filter
}

func (filter *GoshaFilterIds) ClearIds() *GoshaFilterIds {
	filter.Ids = []int{}
	return filter
}

func (filter *GoshaFilterIds) ClearExceptId() *GoshaFilterIds {
	filter.ExceptIds = []int{}
	return filter
}

// method find read create update delete
func (filter *GoshaFilterIds) Validate(functionType string) {

	switch functionType {
	case settings.FunctionTypeFind:

		break
	case settings.FunctionTypeCreate:

		break
	case settings.FunctionTypeRead:
		if len(filter.GetIds()) != 1 || filter.GetIds()[0] < 1 {
			filter.AddValidationError("Error parse Id", errors.ErrorCodeParseId, "Id")
		}
		break
	case settings.FunctionTypeUpdate:
		if len(filter.GetIds()) != 1 || filter.GetIds()[0] < 1 {
			filter.AddValidationError("Error parse Id", errors.ErrorCodeParseId, "Id")
		}
		break
	case settings.FunctionTypeDelete:
		if len(filter.GetIds()) != 1 || filter.GetIds()[0] < 1 {
			filter.AddValidationError("Error parse Id", errors.ErrorCodeParseId, "Id")
		}
		break
	case settings.FunctionTypeMultiDelete:
		break
	case settings.FunctionTypeFindOrCreate:
		break
	case settings.FunctionTypeMultiCreate:
		break
	case settings.FunctionTypeMultiUpdate:
		break
	case settings.FunctionTypeUpdateOrCreate:
		break
	default:
		filter.AddValidationError("Unsupported function type: "+functionType, errors.ErrorCodeUnsupportedFunctionType, "")
		break
	}
}

type GoshaSearchFilter struct {
	Search   string
	SearchBy []string
}

type GoshaOrderFilter struct {
	Order          []string
	OrderDirection []string
}

type GoshaDebugFilter struct {
	isDebug bool
}

func (filter *GoshaDebugFilter) SetDebug(isDebug bool) {
	filter.isDebug = isDebug
}

func (filter GoshaDebugFilter) IsDebug() bool {
	return filter.isDebug
}

type AbstractFilter struct {
	request *http.Request
	rawRequestBody []byte

    Regionality
	GoshaSearchFilter
	GoshaOrderFilter
	GoshaFilterIds
	Pagination
	validator
	Authenticator
	GoshaDebugFilter
}

func GetAbstractFilter(request *http.Request, requestBody []byte, functionType string) (filter AbstractFilter, err error) {

	filter.request = request
    filter.rawRequestBody = requestBody
	filter.functionType = functionType
	filter.urlPath = request.URL.Path

	if !isGroupFunctionType(functionType) {
        err = ReadJSON(filter.rawRequestBody, &filter.GoshaFilterIds)
        if err != nil {
            return
        }
    }

	filter.Pagination.CurrentPage, _ = strconv.Atoi(request.FormValue("CurrentPage"))
	filter.Pagination.PerPage, _ = strconv.Atoi(request.FormValue("PerPage"))
	filter.Search = request.FormValue("Search")

	isDebug, _ := strconv.ParseBool(request.FormValue("IsDebug"))
	filter.GoshaDebugFilter.SetDebug(isDebug)

	arr, _ := url.ParseQuery(request.URL.RawQuery)

	dirs := []string{}

	for _, dir := range arr["OrderDirection[]"] {

		if strings.ToLower(dir) == settings.OrderDirectionDesc {
			dirs = append(dirs, settings.OrderDirectionDesc)
		} else {
			dirs = append(dirs, settings.OrderDirectionAsc)
		}
	}

	for index, field := range arr["Order[]"] {

		filter.Order = append(filter.Order, gorm.ToColumnName(field))

		if len(dirs) > index && dirs[index] == "desc" {
			filter.OrderDirection = append(filter.OrderDirection, "desc")
		} else {
			filter.OrderDirection = append(filter.OrderDirection, "asc")
		}
	}

	if len(filter.Order) < 1 && len(filter.OrderDirection) < 1 {
		filter.Order = append(filter.Order, "id")
		filter.OrderDirection = append(filter.OrderDirection, "asc")
	}

	for _, field := range arr["SearchBy[]"] {
		filter.SearchBy = append(filter.SearchBy, gorm.ToColumnName(field))
	}

	filter.SetToken(request)
	filter.SetIp(request)

	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])

	if id > 0 {
		filter.AddId(id)
		filter.SetCurrentId(id)
	}

	for _, field := range arr["Ids[]"] {
		id, _ := strconv.Atoi(field)
		filter.AddId(id)
	}
	for _, field := range arr["ExceptIds[]"] {
		id, _ := strconv.Atoi(field)
		filter.AddExceptIds(id)
	}

	filter.Validate(functionType)

	return filter, err
}

func (filter *AbstractFilter) IsValid() bool {

	return filter.GoshaFilterIds.IsValid() &&
		filter.Pagination.IsValid() &&
		filter.validator.IsValid() &&
		filter.Authenticator.IsValid()
}

func (filter *AbstractFilter) Validate(functionType string) {

	filter.GoshaFilterIds.Validate(functionType)
	filter.Pagination.Validate(functionType)
	filter.validator.Validate(functionType)
	filter.Authenticator.Validate(functionType)
}

func (filter *AbstractFilter) GetValidationError() error {
	return errors.JoinValidatorError(filter.GoshaFilterIds.GetValidationError(),
		filter.Pagination.GetValidationError(),
		filter.validator.GetValidationError(),
		filter.Authenticator.GetValidationError())
}

func (filter *AbstractFilter) ValidatePerPage() {
	if filter.PerPage > filter.GetMaxPerPage() {
		filter.AddValidationError("PerPage more than maximum", errors.ErrorCodeInvalidPerPage, "PerPage")
	}
}

func (filter *AbstractFilter) GetHost() string {
	return filter.request.Host
}

func (filter *AbstractFilter) GetCurrentIp() string {

	ip := filter.request.Header.Get("X-Forwarded-For")
	return ip
}

func (filter *AbstractFilter) GetCurrentUserAgent() string {
	return filter.request.UserAgent()
}

func isGroupFunctionType(functionType string) bool {
    switch functionType {
    case settings.FunctionTypeMultiCreate, settings.FunctionTypeMultiUpdate, settings.FunctionTypeMultiDelete, settings.FunctionTypeMultiFindOrCreate:
        return true
    default:
        return false
    }
}

