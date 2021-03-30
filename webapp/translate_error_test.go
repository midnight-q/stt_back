package webapp

import (
	"stt_back/settings"
	"stt_back/tests"
	"stt_back/types"
	"stt_back/core"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var validModelTranslateError = types.TranslateError{
		Id:   1,
		//Name: "Some Name",
	}

var updateModelTranslateError = types.TranslateError{
		Id:   1,
		//Name: "Some Another Name",
	}

var idsForRemoveTranslateError = []int{}

func validateFieldsTranslateError(t *testing.T, testModel types.TranslateError, validModelTranslateError types.TranslateError, response *httptest.ResponseRecorder) {

	if testModel.Id < 1 {
		t.Error("Fail test creating new TranslateError", "expect id > 0", "got", testModel.Id, "response:", response.Body)
	}

	//if testModel.Name != validModelTranslateError.Name {
	//	t.Error("Fail test creating new TranslateError", "expect Name =", validModelTranslateError.Name, "got", testModel.Name, "response:", response.Body)
	//}
}

var createdModelTranslateError types.TranslateError

var testCreateFuncTranslateError = func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {
	response := httptest.NewRecorder()
	TranslateErrorCreate(response, tt.Request)
	return response, nil
}

var createAdminRequestTranslateError = tests.GetCreateAdminRequest(settings.TranslateErrorRoute, validModelTranslateError)

var testCasesTranslateError = []tests.WebTest{
    {
		Name:         "Find TranslateErrors as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			req := tests.GetCreateAdminRequest(settings.TranslateErrorRoute, validModelTranslateError)
			webtest := tests.WebTest{Request: req}
			webtest.Name = "Creating before find"

			_, err := testCreateFuncTranslateError(webtest)

			if err != nil {
				return nil, err
			}

			route := settings.TranslateErrorRoute
			request := tests.GetFindAdminRequest(route, 1, 1)

			return tests.SendRequest(route, request, TranslateErrorFind, http.MethodGet), nil
		},
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {

			list, total := getTranslateErrorParsedList(response)

			if total < 1 {
				t.Error("Error in find TranslateError. Total rows must be > 0, got", total)
				return
			}

			for _, item := range list {
				validateFieldsTranslateError(t, item, validModelTranslateError, response)
			}
		},
	},
    {
		Name:         "Create new TranslateError as admin",
		Request:      createAdminRequestTranslateError,
		ResponseCode: 201,
		TestFunc:     testCreateFuncTranslateError,
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {

			createdModelTranslateError = getTranslateErrorParsedModel(response)

			idsForRemoveTranslateError = append(idsForRemoveTranslateError, createdModelTranslateError.Id)
			validateFieldsTranslateError(t, createdModelTranslateError, validModelTranslateError, response)
		},
	},
	{
		Name:            "Create new TranslateError as non authorized user",
		Request:         tests.GetCreateNonAuthorizedUserRequest(settings.TranslateErrorRoute, validModelTranslateError),
		ResponseCode:    403,
		TestFunc:        testCreateFuncTranslateError,
	},
    {
		Name:         "Read TranslateError as non authorized user",
		ResponseCode: 403,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {
			request := tests.GetReadNonAuthorizedUserRequest(settings.TranslateErrorRoute + "/" + strconv.Itoa(updateModelTranslateError.Id))
			return tests.SendRequest(settings.TranslateErrorRoute + "/{id}", request, TranslateErrorRead, http.MethodGet), nil
		},
	},
	{
		Name:         "Read TranslateError as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			tt.Request = tests.GetCreateAdminRequest(settings.TranslateErrorRoute, updateModelTranslateError)
			responseCreate, err := testCreateFuncTranslateError(tt)

			if err != nil {
				return nil, err
			}

			responseData, err := getTranslateErrorResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			idsForRemoveTranslateError = append(idsForRemoveTranslateError, responseData.Id)

			id := strconv.Itoa(responseData.Id)
			request := tests.GetReadAdminRequest(settings.TranslateErrorRoute + "/" + id)

			return tests.SendRequest(settings.TranslateErrorRoute + "/{id}", request, TranslateErrorRead, http.MethodGet), nil

		},
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {
			createdModelTranslateError = getTranslateErrorParsedModel(response)
			validateFieldsTranslateError(t, createdModelTranslateError, updateModelTranslateError, response)
		},
	},
    {
		Name:         "Update TranslateError as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			tt.Request = tests.GetCreateAdminRequest(settings.TranslateErrorRoute, validModelTranslateError)
			responseCreate, err := testCreateFuncTranslateError(tt)

			if err != nil {
				return nil, err
			}

			responseData, err := getTranslateErrorResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			idsForRemoveTranslateError = append(idsForRemoveTranslateError, responseData.Id)

			updateModelTranslateError.Id = responseData.Id

			id := strconv.Itoa(responseData.Id)
			request := tests.GetUpdateAdminRequest(settings.TranslateErrorRoute + "/" + id, updateModelTranslateError)

			return tests.SendRequest(settings.TranslateErrorRoute + "/{id}", request, TranslateErrorUpdate, http.MethodPut), nil
		},
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {
			model := getTranslateErrorParsedModel(response)
			validateFieldsTranslateError(t, model, updateModelTranslateError, response)

			id := strconv.Itoa(model.Id)
			request := tests.GetReadAdminRequest(settings.TranslateErrorRoute + "/" + id)
			readResponse := tests.SendRequest(settings.TranslateErrorRoute + "/{id}", request, TranslateErrorRead, http.MethodGet)

			model = getTranslateErrorParsedModel(readResponse)
			validateFieldsTranslateError(t, model, updateModelTranslateError, readResponse)
		},
	},
    {
		Name: "Delete TranslateError as unauthorized user",
		//Request: inside delete func,
		ResponseCode: 403,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			responseCreate := httptest.NewRecorder()
			TranslateErrorCreate(responseCreate, tests.GetCreateAdminRequest(settings.TranslateErrorRoute, validModelTranslateError))

			responseData, err := getTranslateErrorResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			idsForRemoveTranslateError = append(idsForRemoveTranslateError, responseData.Id)

			id := strconv.Itoa(responseData.Id)
			tt.Request = tests.GetDeleteNonAuthorizedUserRequest(settings.TranslateErrorRoute + "/" + id)

			return tests.SendRequest(settings.TranslateErrorRoute+"/{id}", tt.Request, TranslateErrorDelete, http.MethodDelete), nil
		},
	},
	{
		Name: "Delete TranslateError as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			// create TranslateError for next delete
			responseCreate := httptest.NewRecorder()
			TranslateErrorCreate(responseCreate, tests.GetCreateAdminRequest(settings.TranslateErrorRoute, validModelTranslateError))

			responseData, err := getTranslateErrorResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			id := strconv.Itoa(responseData.Id)
			req := tests.GetDeleteAdminRequest(settings.TranslateErrorRoute + "/" + id)

			return tests.SendRequest(settings.TranslateErrorRoute+"/{id}", req, TranslateErrorDelete, http.MethodDelete), nil
		},
	},
	{
		Name: "Delete TranslateError as admin two times",
		ResponseCode: 400,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			// create TranslateError for next delete
			responseCreate := httptest.NewRecorder()
			TranslateErrorCreate(responseCreate, tests.GetCreateAdminRequest(settings.TranslateErrorRoute, validModelTranslateError))

			responseData, err := getTranslateErrorResponseModel(tt, responseCreate)
			if err != nil {
				return nil, err
			}

			id := strconv.Itoa(responseData.Id)
			req := tests.GetDeleteAdminRequest(settings.TranslateErrorRoute + "/" + id)
			tests.SendRequest(settings.TranslateErrorRoute+"/{id}", req, TranslateErrorDelete, http.MethodDelete)

			return tests.SendRequest(settings.TranslateErrorRoute+"/{id}", req, TranslateErrorDelete, http.MethodDelete), nil
		},
	},
    

    

}

func getTranslateErrorResponseModel(tt tests.WebTest, response *httptest.ResponseRecorder) (types.TranslateError, error) {

	model := getTranslateErrorParsedModel(response)

	if model.Id < 1 {
		return types.TranslateError{}, errors.New("Test " + tt.Name + " fail")
	}

	return model, nil
}


func getTranslateErrorParsedList(response *httptest.ResponseRecorder) (list []types.TranslateError, total int) {

	responseData := struct{
		List []types.TranslateError
		Total int
	} {
		List: []types.TranslateError{},
		Total: 1,
	}

	json.Unmarshal(response.Body.Bytes(), &responseData)

	return responseData.List, responseData.Total
}

func getTranslateErrorParsedModel(response *httptest.ResponseRecorder) types.TranslateError {

	responseData := struct{
		Model types.TranslateError
	} {
		Model:types.TranslateError{},
	}
	json.Unmarshal(response.Body.Bytes(), &responseData)

	return responseData.Model
}


func TestTranslateError(t *testing.T) {

	tmpDb := core.Db
	core.Db = tmpDb.Begin()

	for _, tt := range testCasesTranslateError {
		t.Run(tt.Name, func(t *testing.T) {
			tests.TestFunction(t, tt)
		})
	}

	// clear created data from database
	defer func() {
		core.Db.Rollback()
		core.Db = tmpDb
	}()
}

