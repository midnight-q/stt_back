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

var validModelConvertFile = types.ConvertFile{
		Id:   1,
		//Name: "Some Name",
	}

var updateModelConvertFile = types.ConvertFile{
		Id:   1,
		//Name: "Some Another Name",
	}

var idsForRemoveConvertFile = []int{}

func validateFieldsConvertFile(t *testing.T, testModel types.ConvertFile, validModelConvertFile types.ConvertFile, response *httptest.ResponseRecorder) {

	if testModel.Id < 1 {
		t.Error("Fail test creating new ConvertFile", "expect id > 0", "got", testModel.Id, "response:", response.Body)
	}

	//if testModel.Name != validModelConvertFile.Name {
	//	t.Error("Fail test creating new ConvertFile", "expect Name =", validModelConvertFile.Name, "got", testModel.Name, "response:", response.Body)
	//}
}

var createdModelConvertFile types.ConvertFile

var testCreateFuncConvertFile = func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {
	response := httptest.NewRecorder()
	ConvertFileCreate(response, tt.Request)
	return response, nil
}

var createAdminRequestConvertFile = tests.GetCreateAdminRequest(settings.ConvertFileRoute, validModelConvertFile)

var testCasesConvertFile = []tests.WebTest{
    {
		Name:         "Find ConvertFiles as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			req := tests.GetCreateAdminRequest(settings.ConvertFileRoute, validModelConvertFile)
			webtest := tests.WebTest{Request: req}
			webtest.Name = "Creating before find"

			_, err := testCreateFuncConvertFile(webtest)

			if err != nil {
				return nil, err
			}

			route := settings.ConvertFileRoute
			request := tests.GetFindAdminRequest(route, 1, 1)

			return tests.SendRequest(route, request, ConvertFileFind, http.MethodGet), nil
		},
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {

			list, total := getConvertFileParsedList(response)

			if total < 1 {
				t.Error("Error in find ConvertFile. Total rows must be > 0, got", total)
				return
			}

			for _, item := range list {
				validateFieldsConvertFile(t, item, validModelConvertFile, response)
			}
		},
	},
    {
		Name:         "Create new ConvertFile as admin",
		Request:      createAdminRequestConvertFile,
		ResponseCode: 201,
		TestFunc:     testCreateFuncConvertFile,
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {

			createdModelConvertFile = getConvertFileParsedModel(response)

			idsForRemoveConvertFile = append(idsForRemoveConvertFile, createdModelConvertFile.Id)
			validateFieldsConvertFile(t, createdModelConvertFile, validModelConvertFile, response)
		},
	},
	{
		Name:            "Create new ConvertFile as non authorized user",
		Request:         tests.GetCreateNonAuthorizedUserRequest(settings.ConvertFileRoute, validModelConvertFile),
		ResponseCode:    403,
		TestFunc:        testCreateFuncConvertFile,
	},
    {
		Name:         "Read ConvertFile as non authorized user",
		ResponseCode: 403,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {
			request := tests.GetReadNonAuthorizedUserRequest(settings.ConvertFileRoute + "/" + strconv.Itoa(updateModelConvertFile.Id))
			return tests.SendRequest(settings.ConvertFileRoute + "/{id}", request, ConvertFileRead, http.MethodGet), nil
		},
	},
	{
		Name:         "Read ConvertFile as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			tt.Request = tests.GetCreateAdminRequest(settings.ConvertFileRoute, updateModelConvertFile)
			responseCreate, err := testCreateFuncConvertFile(tt)

			if err != nil {
				return nil, err
			}

			responseData, err := getConvertFileResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			idsForRemoveConvertFile = append(idsForRemoveConvertFile, responseData.Id)

			id := strconv.Itoa(responseData.Id)
			request := tests.GetReadAdminRequest(settings.ConvertFileRoute + "/" + id)

			return tests.SendRequest(settings.ConvertFileRoute + "/{id}", request, ConvertFileRead, http.MethodGet), nil

		},
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {
			createdModelConvertFile = getConvertFileParsedModel(response)
			validateFieldsConvertFile(t, createdModelConvertFile, updateModelConvertFile, response)
		},
	},
    {
		Name:         "Update ConvertFile as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			tt.Request = tests.GetCreateAdminRequest(settings.ConvertFileRoute, validModelConvertFile)
			responseCreate, err := testCreateFuncConvertFile(tt)

			if err != nil {
				return nil, err
			}

			responseData, err := getConvertFileResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			idsForRemoveConvertFile = append(idsForRemoveConvertFile, responseData.Id)

			updateModelConvertFile.Id = responseData.Id

			id := strconv.Itoa(responseData.Id)
			request := tests.GetUpdateAdminRequest(settings.ConvertFileRoute + "/" + id, updateModelConvertFile)

			return tests.SendRequest(settings.ConvertFileRoute + "/{id}", request, ConvertFileUpdate, http.MethodPut), nil
		},
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {
			model := getConvertFileParsedModel(response)
			validateFieldsConvertFile(t, model, updateModelConvertFile, response)

			id := strconv.Itoa(model.Id)
			request := tests.GetReadAdminRequest(settings.ConvertFileRoute + "/" + id)
			readResponse := tests.SendRequest(settings.ConvertFileRoute + "/{id}", request, ConvertFileRead, http.MethodGet)

			model = getConvertFileParsedModel(readResponse)
			validateFieldsConvertFile(t, model, updateModelConvertFile, readResponse)
		},
	},
    {
		Name: "Delete ConvertFile as unauthorized user",
		//Request: inside delete func,
		ResponseCode: 403,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			responseCreate := httptest.NewRecorder()
			ConvertFileCreate(responseCreate, tests.GetCreateAdminRequest(settings.ConvertFileRoute, validModelConvertFile))

			responseData, err := getConvertFileResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			idsForRemoveConvertFile = append(idsForRemoveConvertFile, responseData.Id)

			id := strconv.Itoa(responseData.Id)
			tt.Request = tests.GetDeleteNonAuthorizedUserRequest(settings.ConvertFileRoute + "/" + id)

			return tests.SendRequest(settings.ConvertFileRoute+"/{id}", tt.Request, ConvertFileDelete, http.MethodDelete), nil
		},
	},
	{
		Name: "Delete ConvertFile as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			// create ConvertFile for next delete
			responseCreate := httptest.NewRecorder()
			ConvertFileCreate(responseCreate, tests.GetCreateAdminRequest(settings.ConvertFileRoute, validModelConvertFile))

			responseData, err := getConvertFileResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			id := strconv.Itoa(responseData.Id)
			req := tests.GetDeleteAdminRequest(settings.ConvertFileRoute + "/" + id)

			return tests.SendRequest(settings.ConvertFileRoute+"/{id}", req, ConvertFileDelete, http.MethodDelete), nil
		},
	},
	{
		Name: "Delete ConvertFile as admin two times",
		ResponseCode: 400,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			// create ConvertFile for next delete
			responseCreate := httptest.NewRecorder()
			ConvertFileCreate(responseCreate, tests.GetCreateAdminRequest(settings.ConvertFileRoute, validModelConvertFile))

			responseData, err := getConvertFileResponseModel(tt, responseCreate)
			if err != nil {
				return nil, err
			}

			id := strconv.Itoa(responseData.Id)
			req := tests.GetDeleteAdminRequest(settings.ConvertFileRoute + "/" + id)
			tests.SendRequest(settings.ConvertFileRoute+"/{id}", req, ConvertFileDelete, http.MethodDelete)

			return tests.SendRequest(settings.ConvertFileRoute+"/{id}", req, ConvertFileDelete, http.MethodDelete), nil
		},
	},
    

    

}

func getConvertFileResponseModel(tt tests.WebTest, response *httptest.ResponseRecorder) (types.ConvertFile, error) {

	model := getConvertFileParsedModel(response)

	if model.Id < 1 {
		return types.ConvertFile{}, errors.New("Test " + tt.Name + " fail")
	}

	return model, nil
}


func getConvertFileParsedList(response *httptest.ResponseRecorder) (list []types.ConvertFile, total int) {

	responseData := struct{
		List []types.ConvertFile
		Total int
	} {
		List: []types.ConvertFile{},
		Total: 1,
	}

	json.Unmarshal(response.Body.Bytes(), &responseData)

	return responseData.List, responseData.Total
}

func getConvertFileParsedModel(response *httptest.ResponseRecorder) types.ConvertFile {

	responseData := struct{
		Model types.ConvertFile
	} {
		Model:types.ConvertFile{},
	}
	json.Unmarshal(response.Body.Bytes(), &responseData)

	return responseData.Model
}


func TestConvertFile(t *testing.T) {

	tmpDb := core.Db
	core.Db = tmpDb.Begin()

	for _, tt := range testCasesConvertFile {
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

