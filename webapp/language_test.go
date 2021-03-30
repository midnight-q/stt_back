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

var validModelLanguage = types.Language{
		Id:   1,
		//Name: "Some Name",
	}

var updateModelLanguage = types.Language{
		Id:   1,
		//Name: "Some Another Name",
	}

var idsForRemoveLanguage = []int{}

func validateFieldsLanguage(t *testing.T, testModel types.Language, validModelLanguage types.Language, response *httptest.ResponseRecorder) {

	if testModel.Id < 1 {
		t.Error("Fail test creating new Language", "expect id > 0", "got", testModel.Id, "response:", response.Body)
	}

	//if testModel.Name != validModelLanguage.Name {
	//	t.Error("Fail test creating new Language", "expect Name =", validModelLanguage.Name, "got", testModel.Name, "response:", response.Body)
	//}
}

var createdModelLanguage types.Language

var testCreateFuncLanguage = func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {
	response := httptest.NewRecorder()
	LanguageCreate(response, tt.Request)
	return response, nil
}

var createAdminRequestLanguage = tests.GetCreateAdminRequest(settings.LanguageRoute, validModelLanguage)

var testCasesLanguage = []tests.WebTest{
    {
		Name:         "Find Languages as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			req := tests.GetCreateAdminRequest(settings.LanguageRoute, validModelLanguage)
			webtest := tests.WebTest{Request: req}
			webtest.Name = "Creating before find"

			_, err := testCreateFuncLanguage(webtest)

			if err != nil {
				return nil, err
			}

			route := settings.LanguageRoute
			request := tests.GetFindAdminRequest(route, 1, 1)

			return tests.SendRequest(route, request, LanguageFind, http.MethodGet), nil
		},
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {

			list, total := getLanguageParsedList(response)

			if total < 1 {
				t.Error("Error in find Language. Total rows must be > 0, got", total)
				return
			}

			for _, item := range list {
				validateFieldsLanguage(t, item, validModelLanguage, response)
			}
		},
	},
    {
		Name:         "Create new Language as admin",
		Request:      createAdminRequestLanguage,
		ResponseCode: 201,
		TestFunc:     testCreateFuncLanguage,
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {

			createdModelLanguage = getLanguageParsedModel(response)

			idsForRemoveLanguage = append(idsForRemoveLanguage, createdModelLanguage.Id)
			validateFieldsLanguage(t, createdModelLanguage, validModelLanguage, response)
		},
	},
	{
		Name:            "Create new Language as non authorized user",
		Request:         tests.GetCreateNonAuthorizedUserRequest(settings.LanguageRoute, validModelLanguage),
		ResponseCode:    403,
		TestFunc:        testCreateFuncLanguage,
	},
    {
		Name:         "Read Language as non authorized user",
		ResponseCode: 403,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {
			request := tests.GetReadNonAuthorizedUserRequest(settings.LanguageRoute + "/" + strconv.Itoa(updateModelLanguage.Id))
			return tests.SendRequest(settings.LanguageRoute + "/{id}", request, LanguageRead, http.MethodGet), nil
		},
	},
	{
		Name:         "Read Language as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			tt.Request = tests.GetCreateAdminRequest(settings.LanguageRoute, updateModelLanguage)
			responseCreate, err := testCreateFuncLanguage(tt)

			if err != nil {
				return nil, err
			}

			responseData, err := getLanguageResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			idsForRemoveLanguage = append(idsForRemoveLanguage, responseData.Id)

			id := strconv.Itoa(responseData.Id)
			request := tests.GetReadAdminRequest(settings.LanguageRoute + "/" + id)

			return tests.SendRequest(settings.LanguageRoute + "/{id}", request, LanguageRead, http.MethodGet), nil

		},
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {
			createdModelLanguage = getLanguageParsedModel(response)
			validateFieldsLanguage(t, createdModelLanguage, updateModelLanguage, response)
		},
	},
    {
		Name:         "Update Language as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			tt.Request = tests.GetCreateAdminRequest(settings.LanguageRoute, validModelLanguage)
			responseCreate, err := testCreateFuncLanguage(tt)

			if err != nil {
				return nil, err
			}

			responseData, err := getLanguageResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			idsForRemoveLanguage = append(idsForRemoveLanguage, responseData.Id)

			updateModelLanguage.Id = responseData.Id

			id := strconv.Itoa(responseData.Id)
			request := tests.GetUpdateAdminRequest(settings.LanguageRoute + "/" + id, updateModelLanguage)

			return tests.SendRequest(settings.LanguageRoute + "/{id}", request, LanguageUpdate, http.MethodPut), nil
		},
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {
			model := getLanguageParsedModel(response)
			validateFieldsLanguage(t, model, updateModelLanguage, response)

			id := strconv.Itoa(model.Id)
			request := tests.GetReadAdminRequest(settings.LanguageRoute + "/" + id)
			readResponse := tests.SendRequest(settings.LanguageRoute + "/{id}", request, LanguageRead, http.MethodGet)

			model = getLanguageParsedModel(readResponse)
			validateFieldsLanguage(t, model, updateModelLanguage, readResponse)
		},
	},
    {
		Name: "Delete Language as unauthorized user",
		//Request: inside delete func,
		ResponseCode: 403,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			responseCreate := httptest.NewRecorder()
			LanguageCreate(responseCreate, tests.GetCreateAdminRequest(settings.LanguageRoute, validModelLanguage))

			responseData, err := getLanguageResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			idsForRemoveLanguage = append(idsForRemoveLanguage, responseData.Id)

			id := strconv.Itoa(responseData.Id)
			tt.Request = tests.GetDeleteNonAuthorizedUserRequest(settings.LanguageRoute + "/" + id)

			return tests.SendRequest(settings.LanguageRoute+"/{id}", tt.Request, LanguageDelete, http.MethodDelete), nil
		},
	},
	{
		Name: "Delete Language as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			// create Language for next delete
			responseCreate := httptest.NewRecorder()
			LanguageCreate(responseCreate, tests.GetCreateAdminRequest(settings.LanguageRoute, validModelLanguage))

			responseData, err := getLanguageResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			id := strconv.Itoa(responseData.Id)
			req := tests.GetDeleteAdminRequest(settings.LanguageRoute + "/" + id)

			return tests.SendRequest(settings.LanguageRoute+"/{id}", req, LanguageDelete, http.MethodDelete), nil
		},
	},
	{
		Name: "Delete Language as admin two times",
		ResponseCode: 400,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			// create Language for next delete
			responseCreate := httptest.NewRecorder()
			LanguageCreate(responseCreate, tests.GetCreateAdminRequest(settings.LanguageRoute, validModelLanguage))

			responseData, err := getLanguageResponseModel(tt, responseCreate)
			if err != nil {
				return nil, err
			}

			id := strconv.Itoa(responseData.Id)
			req := tests.GetDeleteAdminRequest(settings.LanguageRoute + "/" + id)
			tests.SendRequest(settings.LanguageRoute+"/{id}", req, LanguageDelete, http.MethodDelete)

			return tests.SendRequest(settings.LanguageRoute+"/{id}", req, LanguageDelete, http.MethodDelete), nil
		},
	},
    

    

}

func getLanguageResponseModel(tt tests.WebTest, response *httptest.ResponseRecorder) (types.Language, error) {

	model := getLanguageParsedModel(response)

	if model.Id < 1 {
		return types.Language{}, errors.New("Test " + tt.Name + " fail")
	}

	return model, nil
}


func getLanguageParsedList(response *httptest.ResponseRecorder) (list []types.Language, total int) {

	responseData := struct{
		List []types.Language
		Total int
	} {
		List: []types.Language{},
		Total: 1,
	}

	json.Unmarshal(response.Body.Bytes(), &responseData)

	return responseData.List, responseData.Total
}

func getLanguageParsedModel(response *httptest.ResponseRecorder) types.Language {

	responseData := struct{
		Model types.Language
	} {
		Model:types.Language{},
	}
	json.Unmarshal(response.Body.Bytes(), &responseData)

	return responseData.Model
}


func TestLanguage(t *testing.T) {

	tmpDb := core.Db
	core.Db = tmpDb.Begin()

	for _, tt := range testCasesLanguage {
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

