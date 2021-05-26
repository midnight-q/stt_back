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

var validModelCheckConverterLog = types.CheckConverterLog{
		Id:   1,
		//Name: "Some Name",
	}

var updateModelCheckConverterLog = types.CheckConverterLog{
		Id:   1,
		//Name: "Some Another Name",
	}

var idsForRemoveCheckConverterLog = []int{}

func validateFieldsCheckConverterLog(t *testing.T, testModel types.CheckConverterLog, validModelCheckConverterLog types.CheckConverterLog, response *httptest.ResponseRecorder) {

	if testModel.Id < 1 {
		t.Error("Fail test creating new CheckConverterLog", "expect id > 0", "got", testModel.Id, "response:", response.Body)
	}

	//if testModel.Name != validModelCheckConverterLog.Name {
	//	t.Error("Fail test creating new CheckConverterLog", "expect Name =", validModelCheckConverterLog.Name, "got", testModel.Name, "response:", response.Body)
	//}
}

var createdModelCheckConverterLog types.CheckConverterLog

var testCreateFuncCheckConverterLog = func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {
	response := httptest.NewRecorder()
	CheckConverterLogCreate(response, tt.Request)
	return response, nil
}

var createAdminRequestCheckConverterLog = tests.GetCreateAdminRequest(settings.CheckConverterLogRoute, validModelCheckConverterLog)

var testCasesCheckConverterLog = []tests.WebTest{
    {
		Name:         "Find CheckConverterLogs as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			req := tests.GetCreateAdminRequest(settings.CheckConverterLogRoute, validModelCheckConverterLog)
			webtest := tests.WebTest{Request: req}
			webtest.Name = "Creating before find"

			_, err := testCreateFuncCheckConverterLog(webtest)

			if err != nil {
				return nil, err
			}

			route := settings.CheckConverterLogRoute
			request := tests.GetFindAdminRequest(route, 1, 1)

			return tests.SendRequest(route, request, CheckConverterLogFind, http.MethodGet), nil
		},
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {

			list, total := getCheckConverterLogParsedList(response)

			if total < 1 {
				t.Error("Error in find CheckConverterLog. Total rows must be > 0, got", total)
				return
			}

			for _, item := range list {
				validateFieldsCheckConverterLog(t, item, validModelCheckConverterLog, response)
			}
		},
	},
    {
		Name:         "Create new CheckConverterLog as admin",
		Request:      createAdminRequestCheckConverterLog,
		ResponseCode: 201,
		TestFunc:     testCreateFuncCheckConverterLog,
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {

			createdModelCheckConverterLog = getCheckConverterLogParsedModel(response)

			idsForRemoveCheckConverterLog = append(idsForRemoveCheckConverterLog, createdModelCheckConverterLog.Id)
			validateFieldsCheckConverterLog(t, createdModelCheckConverterLog, validModelCheckConverterLog, response)
		},
	},
	{
		Name:            "Create new CheckConverterLog as non authorized user",
		Request:         tests.GetCreateNonAuthorizedUserRequest(settings.CheckConverterLogRoute, validModelCheckConverterLog),
		ResponseCode:    403,
		TestFunc:        testCreateFuncCheckConverterLog,
	},
    {
		Name:         "Read CheckConverterLog as non authorized user",
		ResponseCode: 403,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {
			request := tests.GetReadNonAuthorizedUserRequest(settings.CheckConverterLogRoute + "/" + strconv.Itoa(updateModelCheckConverterLog.Id))
			return tests.SendRequest(settings.CheckConverterLogRoute + "/{id}", request, CheckConverterLogRead, http.MethodGet), nil
		},
	},
	{
		Name:         "Read CheckConverterLog as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			tt.Request = tests.GetCreateAdminRequest(settings.CheckConverterLogRoute, updateModelCheckConverterLog)
			responseCreate, err := testCreateFuncCheckConverterLog(tt)

			if err != nil {
				return nil, err
			}

			responseData, err := getCheckConverterLogResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			idsForRemoveCheckConverterLog = append(idsForRemoveCheckConverterLog, responseData.Id)

			id := strconv.Itoa(responseData.Id)
			request := tests.GetReadAdminRequest(settings.CheckConverterLogRoute + "/" + id)

			return tests.SendRequest(settings.CheckConverterLogRoute + "/{id}", request, CheckConverterLogRead, http.MethodGet), nil

		},
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {
			createdModelCheckConverterLog = getCheckConverterLogParsedModel(response)
			validateFieldsCheckConverterLog(t, createdModelCheckConverterLog, updateModelCheckConverterLog, response)
		},
	},
    {
		Name:         "Update CheckConverterLog as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			tt.Request = tests.GetCreateAdminRequest(settings.CheckConverterLogRoute, validModelCheckConverterLog)
			responseCreate, err := testCreateFuncCheckConverterLog(tt)

			if err != nil {
				return nil, err
			}

			responseData, err := getCheckConverterLogResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			idsForRemoveCheckConverterLog = append(idsForRemoveCheckConverterLog, responseData.Id)

			updateModelCheckConverterLog.Id = responseData.Id

			id := strconv.Itoa(responseData.Id)
			request := tests.GetUpdateAdminRequest(settings.CheckConverterLogRoute + "/" + id, updateModelCheckConverterLog)

			return tests.SendRequest(settings.CheckConverterLogRoute + "/{id}", request, CheckConverterLogUpdate, http.MethodPut), nil
		},
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {
			model := getCheckConverterLogParsedModel(response)
			validateFieldsCheckConverterLog(t, model, updateModelCheckConverterLog, response)

			id := strconv.Itoa(model.Id)
			request := tests.GetReadAdminRequest(settings.CheckConverterLogRoute + "/" + id)
			readResponse := tests.SendRequest(settings.CheckConverterLogRoute + "/{id}", request, CheckConverterLogRead, http.MethodGet)

			model = getCheckConverterLogParsedModel(readResponse)
			validateFieldsCheckConverterLog(t, model, updateModelCheckConverterLog, readResponse)
		},
	},
    {
		Name: "Delete CheckConverterLog as unauthorized user",
		//Request: inside delete func,
		ResponseCode: 403,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			responseCreate := httptest.NewRecorder()
			CheckConverterLogCreate(responseCreate, tests.GetCreateAdminRequest(settings.CheckConverterLogRoute, validModelCheckConverterLog))

			responseData, err := getCheckConverterLogResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			idsForRemoveCheckConverterLog = append(idsForRemoveCheckConverterLog, responseData.Id)

			id := strconv.Itoa(responseData.Id)
			tt.Request = tests.GetDeleteNonAuthorizedUserRequest(settings.CheckConverterLogRoute + "/" + id)

			return tests.SendRequest(settings.CheckConverterLogRoute+"/{id}", tt.Request, CheckConverterLogDelete, http.MethodDelete), nil
		},
	},
	{
		Name: "Delete CheckConverterLog as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			// create CheckConverterLog for next delete
			responseCreate := httptest.NewRecorder()
			CheckConverterLogCreate(responseCreate, tests.GetCreateAdminRequest(settings.CheckConverterLogRoute, validModelCheckConverterLog))

			responseData, err := getCheckConverterLogResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			id := strconv.Itoa(responseData.Id)
			req := tests.GetDeleteAdminRequest(settings.CheckConverterLogRoute + "/" + id)

			return tests.SendRequest(settings.CheckConverterLogRoute+"/{id}", req, CheckConverterLogDelete, http.MethodDelete), nil
		},
	},
	{
		Name: "Delete CheckConverterLog as admin two times",
		ResponseCode: 400,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			// create CheckConverterLog for next delete
			responseCreate := httptest.NewRecorder()
			CheckConverterLogCreate(responseCreate, tests.GetCreateAdminRequest(settings.CheckConverterLogRoute, validModelCheckConverterLog))

			responseData, err := getCheckConverterLogResponseModel(tt, responseCreate)
			if err != nil {
				return nil, err
			}

			id := strconv.Itoa(responseData.Id)
			req := tests.GetDeleteAdminRequest(settings.CheckConverterLogRoute + "/" + id)
			tests.SendRequest(settings.CheckConverterLogRoute+"/{id}", req, CheckConverterLogDelete, http.MethodDelete)

			return tests.SendRequest(settings.CheckConverterLogRoute+"/{id}", req, CheckConverterLogDelete, http.MethodDelete), nil
		},
	},
    

    

}

func getCheckConverterLogResponseModel(tt tests.WebTest, response *httptest.ResponseRecorder) (types.CheckConverterLog, error) {

	model := getCheckConverterLogParsedModel(response)

	if model.Id < 1 {
		return types.CheckConverterLog{}, errors.New("Test " + tt.Name + " fail")
	}

	return model, nil
}


func getCheckConverterLogParsedList(response *httptest.ResponseRecorder) (list []types.CheckConverterLog, total int) {

	responseData := struct{
		List []types.CheckConverterLog
		Total int
	} {
		List: []types.CheckConverterLog{},
		Total: 1,
	}

	json.Unmarshal(response.Body.Bytes(), &responseData)

	return responseData.List, responseData.Total
}

func getCheckConverterLogParsedModel(response *httptest.ResponseRecorder) types.CheckConverterLog {

	responseData := struct{
		Model types.CheckConverterLog
	} {
		Model:types.CheckConverterLog{},
	}
	json.Unmarshal(response.Body.Bytes(), &responseData)

	return responseData.Model
}


func TestCheckConverterLog(t *testing.T) {

	tmpDb := core.Db
	core.Db = tmpDb.Begin()

	for _, tt := range testCasesCheckConverterLog {
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

