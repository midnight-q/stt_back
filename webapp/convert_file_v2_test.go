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

var validModelConvertFileV2 = types.ConvertFileV2{
		Id:   1,
		//Name: "Some Name",
	}

var updateModelConvertFileV2 = types.ConvertFileV2{
		Id:   1,
		//Name: "Some Another Name",
	}

var idsForRemoveConvertFileV2 = []int{}

func validateFieldsConvertFileV2(t *testing.T, testModel types.ConvertFileV2, validModelConvertFileV2 types.ConvertFileV2, response *httptest.ResponseRecorder) {

	if testModel.Id < 1 {
		t.Error("Fail test creating new ConvertFileV2", "expect id > 0", "got", testModel.Id, "response:", response.Body)
	}

	//if testModel.Name != validModelConvertFileV2.Name {
	//	t.Error("Fail test creating new ConvertFileV2", "expect Name =", validModelConvertFileV2.Name, "got", testModel.Name, "response:", response.Body)
	//}
}

var createdModelConvertFileV2 types.ConvertFileV2

var testCreateFuncConvertFileV2 = func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {
	response := httptest.NewRecorder()
	ConvertFileV2Create(response, tt.Request)
	return response, nil
}

var createAdminRequestConvertFileV2 = tests.GetCreateAdminRequest(settings.ConvertFileV2Route, validModelConvertFileV2)

var testCasesConvertFileV2 = []tests.WebTest{
    {
		Name:         "Find ConvertFileV2s as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			req := tests.GetCreateAdminRequest(settings.ConvertFileV2Route, validModelConvertFileV2)
			webtest := tests.WebTest{Request: req}
			webtest.Name = "Creating before find"

			_, err := testCreateFuncConvertFileV2(webtest)

			if err != nil {
				return nil, err
			}

			route := settings.ConvertFileV2Route
			request := tests.GetFindAdminRequest(route, 1, 1)

			return tests.SendRequest(route, request, ConvertFileV2Find, http.MethodGet), nil
		},
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {

			list, total := getConvertFileV2ParsedList(response)

			if total < 1 {
				t.Error("Error in find ConvertFileV2. Total rows must be > 0, got", total)
				return
			}

			for _, item := range list {
				validateFieldsConvertFileV2(t, item, validModelConvertFileV2, response)
			}
		},
	},
    {
		Name:         "Create new ConvertFileV2 as admin",
		Request:      createAdminRequestConvertFileV2,
		ResponseCode: 201,
		TestFunc:     testCreateFuncConvertFileV2,
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {

			createdModelConvertFileV2 = getConvertFileV2ParsedModel(response)

			idsForRemoveConvertFileV2 = append(idsForRemoveConvertFileV2, createdModelConvertFileV2.Id)
			validateFieldsConvertFileV2(t, createdModelConvertFileV2, validModelConvertFileV2, response)
		},
	},
	{
		Name:            "Create new ConvertFileV2 as non authorized user",
		Request:         tests.GetCreateNonAuthorizedUserRequest(settings.ConvertFileV2Route, validModelConvertFileV2),
		ResponseCode:    403,
		TestFunc:        testCreateFuncConvertFileV2,
	},
    {
		Name:         "Read ConvertFileV2 as non authorized user",
		ResponseCode: 403,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {
			request := tests.GetReadNonAuthorizedUserRequest(settings.ConvertFileV2Route + "/" + strconv.Itoa(updateModelConvertFileV2.Id))
			return tests.SendRequest(settings.ConvertFileV2Route + "/{id}", request, ConvertFileV2Read, http.MethodGet), nil
		},
	},
	{
		Name:         "Read ConvertFileV2 as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			tt.Request = tests.GetCreateAdminRequest(settings.ConvertFileV2Route, updateModelConvertFileV2)
			responseCreate, err := testCreateFuncConvertFileV2(tt)

			if err != nil {
				return nil, err
			}

			responseData, err := getConvertFileV2ResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			idsForRemoveConvertFileV2 = append(idsForRemoveConvertFileV2, responseData.Id)

			id := strconv.Itoa(responseData.Id)
			request := tests.GetReadAdminRequest(settings.ConvertFileV2Route + "/" + id)

			return tests.SendRequest(settings.ConvertFileV2Route + "/{id}", request, ConvertFileV2Read, http.MethodGet), nil

		},
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {
			createdModelConvertFileV2 = getConvertFileV2ParsedModel(response)
			validateFieldsConvertFileV2(t, createdModelConvertFileV2, updateModelConvertFileV2, response)
		},
	},
    {
		Name:         "Update ConvertFileV2 as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			tt.Request = tests.GetCreateAdminRequest(settings.ConvertFileV2Route, validModelConvertFileV2)
			responseCreate, err := testCreateFuncConvertFileV2(tt)

			if err != nil {
				return nil, err
			}

			responseData, err := getConvertFileV2ResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			idsForRemoveConvertFileV2 = append(idsForRemoveConvertFileV2, responseData.Id)

			updateModelConvertFileV2.Id = responseData.Id

			id := strconv.Itoa(responseData.Id)
			request := tests.GetUpdateAdminRequest(settings.ConvertFileV2Route + "/" + id, updateModelConvertFileV2)

			return tests.SendRequest(settings.ConvertFileV2Route + "/{id}", request, ConvertFileV2Update, http.MethodPut), nil
		},
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {
			model := getConvertFileV2ParsedModel(response)
			validateFieldsConvertFileV2(t, model, updateModelConvertFileV2, response)

			id := strconv.Itoa(model.Id)
			request := tests.GetReadAdminRequest(settings.ConvertFileV2Route + "/" + id)
			readResponse := tests.SendRequest(settings.ConvertFileV2Route + "/{id}", request, ConvertFileV2Read, http.MethodGet)

			model = getConvertFileV2ParsedModel(readResponse)
			validateFieldsConvertFileV2(t, model, updateModelConvertFileV2, readResponse)
		},
	},
    {
		Name: "Delete ConvertFileV2 as unauthorized user",
		//Request: inside delete func,
		ResponseCode: 403,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			responseCreate := httptest.NewRecorder()
			ConvertFileV2Create(responseCreate, tests.GetCreateAdminRequest(settings.ConvertFileV2Route, validModelConvertFileV2))

			responseData, err := getConvertFileV2ResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			idsForRemoveConvertFileV2 = append(idsForRemoveConvertFileV2, responseData.Id)

			id := strconv.Itoa(responseData.Id)
			tt.Request = tests.GetDeleteNonAuthorizedUserRequest(settings.ConvertFileV2Route + "/" + id)

			return tests.SendRequest(settings.ConvertFileV2Route+"/{id}", tt.Request, ConvertFileV2Delete, http.MethodDelete), nil
		},
	},
	{
		Name: "Delete ConvertFileV2 as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			// create ConvertFileV2 for next delete
			responseCreate := httptest.NewRecorder()
			ConvertFileV2Create(responseCreate, tests.GetCreateAdminRequest(settings.ConvertFileV2Route, validModelConvertFileV2))

			responseData, err := getConvertFileV2ResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			id := strconv.Itoa(responseData.Id)
			req := tests.GetDeleteAdminRequest(settings.ConvertFileV2Route + "/" + id)

			return tests.SendRequest(settings.ConvertFileV2Route+"/{id}", req, ConvertFileV2Delete, http.MethodDelete), nil
		},
	},
	{
		Name: "Delete ConvertFileV2 as admin two times",
		ResponseCode: 400,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			// create ConvertFileV2 for next delete
			responseCreate := httptest.NewRecorder()
			ConvertFileV2Create(responseCreate, tests.GetCreateAdminRequest(settings.ConvertFileV2Route, validModelConvertFileV2))

			responseData, err := getConvertFileV2ResponseModel(tt, responseCreate)
			if err != nil {
				return nil, err
			}

			id := strconv.Itoa(responseData.Id)
			req := tests.GetDeleteAdminRequest(settings.ConvertFileV2Route + "/" + id)
			tests.SendRequest(settings.ConvertFileV2Route+"/{id}", req, ConvertFileV2Delete, http.MethodDelete)

			return tests.SendRequest(settings.ConvertFileV2Route+"/{id}", req, ConvertFileV2Delete, http.MethodDelete), nil
		},
	},
    

    

}

func getConvertFileV2ResponseModel(tt tests.WebTest, response *httptest.ResponseRecorder) (types.ConvertFileV2, error) {

	model := getConvertFileV2ParsedModel(response)

	if model.Id < 1 {
		return types.ConvertFileV2{}, errors.New("Test " + tt.Name + " fail")
	}

	return model, nil
}


func getConvertFileV2ParsedList(response *httptest.ResponseRecorder) (list []types.ConvertFileV2, total int) {

	responseData := struct{
		List []types.ConvertFileV2
		Total int
	} {
		List: []types.ConvertFileV2{},
		Total: 1,
	}

	json.Unmarshal(response.Body.Bytes(), &responseData)

	return responseData.List, responseData.Total
}

func getConvertFileV2ParsedModel(response *httptest.ResponseRecorder) types.ConvertFileV2 {

	responseData := struct{
		Model types.ConvertFileV2
	} {
		Model:types.ConvertFileV2{},
	}
	json.Unmarshal(response.Body.Bytes(), &responseData)

	return responseData.Model
}


func TestConvertFileV2(t *testing.T) {

	tmpDb := core.Db
	core.Db = tmpDb.Begin()

	for _, tt := range testCasesConvertFileV2 {
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

