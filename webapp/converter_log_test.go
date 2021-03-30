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

var validModelConverterLog = types.ConverterLog{
		Id:   1,
		//Name: "Some Name",
	}

var updateModelConverterLog = types.ConverterLog{
		Id:   1,
		//Name: "Some Another Name",
	}

var idsForRemoveConverterLog = []int{}

func validateFieldsConverterLog(t *testing.T, testModel types.ConverterLog, validModelConverterLog types.ConverterLog, response *httptest.ResponseRecorder) {

	if testModel.Id < 1 {
		t.Error("Fail test creating new ConverterLog", "expect id > 0", "got", testModel.Id, "response:", response.Body)
	}

	//if testModel.Name != validModelConverterLog.Name {
	//	t.Error("Fail test creating new ConverterLog", "expect Name =", validModelConverterLog.Name, "got", testModel.Name, "response:", response.Body)
	//}
}

var createdModelConverterLog types.ConverterLog

var testCreateFuncConverterLog = func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {
	response := httptest.NewRecorder()
	ConverterLogCreate(response, tt.Request)
	return response, nil
}

var createAdminRequestConverterLog = tests.GetCreateAdminRequest(settings.ConverterLogRoute, validModelConverterLog)

var testCasesConverterLog = []tests.WebTest{
    {
		Name:         "Find ConverterLogs as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			req := tests.GetCreateAdminRequest(settings.ConverterLogRoute, validModelConverterLog)
			webtest := tests.WebTest{Request: req}
			webtest.Name = "Creating before find"

			_, err := testCreateFuncConverterLog(webtest)

			if err != nil {
				return nil, err
			}

			route := settings.ConverterLogRoute
			request := tests.GetFindAdminRequest(route, 1, 1)

			return tests.SendRequest(route, request, ConverterLogFind, http.MethodGet), nil
		},
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {

			list, total := getConverterLogParsedList(response)

			if total < 1 {
				t.Error("Error in find ConverterLog. Total rows must be > 0, got", total)
				return
			}

			for _, item := range list {
				validateFieldsConverterLog(t, item, validModelConverterLog, response)
			}
		},
	},
    {
		Name:         "Create new ConverterLog as admin",
		Request:      createAdminRequestConverterLog,
		ResponseCode: 201,
		TestFunc:     testCreateFuncConverterLog,
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {

			createdModelConverterLog = getConverterLogParsedModel(response)

			idsForRemoveConverterLog = append(idsForRemoveConverterLog, createdModelConverterLog.Id)
			validateFieldsConverterLog(t, createdModelConverterLog, validModelConverterLog, response)
		},
	},
	{
		Name:            "Create new ConverterLog as non authorized user",
		Request:         tests.GetCreateNonAuthorizedUserRequest(settings.ConverterLogRoute, validModelConverterLog),
		ResponseCode:    403,
		TestFunc:        testCreateFuncConverterLog,
	},
    {
		Name:         "Read ConverterLog as non authorized user",
		ResponseCode: 403,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {
			request := tests.GetReadNonAuthorizedUserRequest(settings.ConverterLogRoute + "/" + strconv.Itoa(updateModelConverterLog.Id))
			return tests.SendRequest(settings.ConverterLogRoute + "/{id}", request, ConverterLogRead, http.MethodGet), nil
		},
	},
	{
		Name:         "Read ConverterLog as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			tt.Request = tests.GetCreateAdminRequest(settings.ConverterLogRoute, updateModelConverterLog)
			responseCreate, err := testCreateFuncConverterLog(tt)

			if err != nil {
				return nil, err
			}

			responseData, err := getConverterLogResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			idsForRemoveConverterLog = append(idsForRemoveConverterLog, responseData.Id)

			id := strconv.Itoa(responseData.Id)
			request := tests.GetReadAdminRequest(settings.ConverterLogRoute + "/" + id)

			return tests.SendRequest(settings.ConverterLogRoute + "/{id}", request, ConverterLogRead, http.MethodGet), nil

		},
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {
			createdModelConverterLog = getConverterLogParsedModel(response)
			validateFieldsConverterLog(t, createdModelConverterLog, updateModelConverterLog, response)
		},
	},
    {
		Name:         "Update ConverterLog as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			tt.Request = tests.GetCreateAdminRequest(settings.ConverterLogRoute, validModelConverterLog)
			responseCreate, err := testCreateFuncConverterLog(tt)

			if err != nil {
				return nil, err
			}

			responseData, err := getConverterLogResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			idsForRemoveConverterLog = append(idsForRemoveConverterLog, responseData.Id)

			updateModelConverterLog.Id = responseData.Id

			id := strconv.Itoa(responseData.Id)
			request := tests.GetUpdateAdminRequest(settings.ConverterLogRoute + "/" + id, updateModelConverterLog)

			return tests.SendRequest(settings.ConverterLogRoute + "/{id}", request, ConverterLogUpdate, http.MethodPut), nil
		},
		ResultValidator: func(t *testing.T, response *httptest.ResponseRecorder) {
			model := getConverterLogParsedModel(response)
			validateFieldsConverterLog(t, model, updateModelConverterLog, response)

			id := strconv.Itoa(model.Id)
			request := tests.GetReadAdminRequest(settings.ConverterLogRoute + "/" + id)
			readResponse := tests.SendRequest(settings.ConverterLogRoute + "/{id}", request, ConverterLogRead, http.MethodGet)

			model = getConverterLogParsedModel(readResponse)
			validateFieldsConverterLog(t, model, updateModelConverterLog, readResponse)
		},
	},
    {
		Name: "Delete ConverterLog as unauthorized user",
		//Request: inside delete func,
		ResponseCode: 403,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			responseCreate := httptest.NewRecorder()
			ConverterLogCreate(responseCreate, tests.GetCreateAdminRequest(settings.ConverterLogRoute, validModelConverterLog))

			responseData, err := getConverterLogResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			idsForRemoveConverterLog = append(idsForRemoveConverterLog, responseData.Id)

			id := strconv.Itoa(responseData.Id)
			tt.Request = tests.GetDeleteNonAuthorizedUserRequest(settings.ConverterLogRoute + "/" + id)

			return tests.SendRequest(settings.ConverterLogRoute+"/{id}", tt.Request, ConverterLogDelete, http.MethodDelete), nil
		},
	},
	{
		Name: "Delete ConverterLog as admin",
		ResponseCode: 200,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			// create ConverterLog for next delete
			responseCreate := httptest.NewRecorder()
			ConverterLogCreate(responseCreate, tests.GetCreateAdminRequest(settings.ConverterLogRoute, validModelConverterLog))

			responseData, err := getConverterLogResponseModel(tt, responseCreate)

			if err != nil {
				return nil, err
			}

			id := strconv.Itoa(responseData.Id)
			req := tests.GetDeleteAdminRequest(settings.ConverterLogRoute + "/" + id)

			return tests.SendRequest(settings.ConverterLogRoute+"/{id}", req, ConverterLogDelete, http.MethodDelete), nil
		},
	},
	{
		Name: "Delete ConverterLog as admin two times",
		ResponseCode: 400,
		TestFunc: func(tt tests.WebTest) (*httptest.ResponseRecorder, error) {

			// create ConverterLog for next delete
			responseCreate := httptest.NewRecorder()
			ConverterLogCreate(responseCreate, tests.GetCreateAdminRequest(settings.ConverterLogRoute, validModelConverterLog))

			responseData, err := getConverterLogResponseModel(tt, responseCreate)
			if err != nil {
				return nil, err
			}

			id := strconv.Itoa(responseData.Id)
			req := tests.GetDeleteAdminRequest(settings.ConverterLogRoute + "/" + id)
			tests.SendRequest(settings.ConverterLogRoute+"/{id}", req, ConverterLogDelete, http.MethodDelete)

			return tests.SendRequest(settings.ConverterLogRoute+"/{id}", req, ConverterLogDelete, http.MethodDelete), nil
		},
	},
    

    

}

func getConverterLogResponseModel(tt tests.WebTest, response *httptest.ResponseRecorder) (types.ConverterLog, error) {

	model := getConverterLogParsedModel(response)

	if model.Id < 1 {
		return types.ConverterLog{}, errors.New("Test " + tt.Name + " fail")
	}

	return model, nil
}


func getConverterLogParsedList(response *httptest.ResponseRecorder) (list []types.ConverterLog, total int) {

	responseData := struct{
		List []types.ConverterLog
		Total int
	} {
		List: []types.ConverterLog{},
		Total: 1,
	}

	json.Unmarshal(response.Body.Bytes(), &responseData)

	return responseData.List, responseData.Total
}

func getConverterLogParsedModel(response *httptest.ResponseRecorder) types.ConverterLog {

	responseData := struct{
		Model types.ConverterLog
	} {
		Model:types.ConverterLog{},
	}
	json.Unmarshal(response.Body.Bytes(), &responseData)

	return responseData.Model
}


func TestConverterLog(t *testing.T) {

	tmpDb := core.Db
	core.Db = tmpDb.Begin()

	for _, tt := range testCasesConverterLog {
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

