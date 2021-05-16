package logic

import (
	"strconv"
	"stt_back/types"

	"fmt"
	"log"
	"strings"
	"stt_back/core"
	"stt_back/dbmodels"
	"stt_back/errors"

	"github.com/jinzhu/gorm"
)

func ConverterLogFind(filter types.ConverterLogFilter) (result []types.ConverterLog, totalRecords int, err error) {

	foundIds := []int{}
	dbmodelData := []dbmodels.ConverterLog{}
	limit := filter.PerPage
	offset := filter.GetOffset()

	filterIds := filter.GetIds()
	filterExceptIds := filter.GetExceptIds()

	var count int

	criteria := core.Db.Where(dbmodels.ConverterLog{})

	//ConverterLog.FindFilterCode remove this line for disable generator functionality

	if len(filterIds) > 0 {
		criteria = criteria.Where("id in (?)", filterIds)
	}

	if len(filterExceptIds) > 0 {
		criteria = criteria.Where("id not in (?)", filterExceptIds)
	}

	if filter.UserId > 0 {
		criteria = criteria.Where(dbmodels.ConverterLog{UserId: filter.UserId})
	}

	//if len(filter.Search) > 0 {
	//
	//    s := ("%" + filter.Search + "%")
	//
	//    if len(filter.SearchBy) > 0 {
	//
	//        for _, field := range filter.SearchBy {
	//
	//            if core.Db.NewScope(&dbmodels.ConverterLog{}).HasColumn(field) {
	//                criteria = criteria.Or("`"+field+"`"+" ilike ?", s)
	//            } else {
	//                err = errors.NewErrorWithCode("Search by unknown field", errors.ErrorCodeNotValid ,field)
	//                return
	//            }
	//        }
	//    } else {
	//      criteria = criteria.Where("name ilike ? or code ilike ?", ("%" + filter.Search + "%"), ("%" + filter.Search + "%"))
	//    }
	//}

	q := criteria.Model(dbmodels.ConverterLog{}).Count(&count)

	if q.Error != nil {
		log.Println("FindConverterLog > Ошибка получения данных:", q.Error)
		return result, 0, nil
	}

	// order global criteria
	if len(filter.Order) > 0 {
		for index, Field := range filter.Order {
			if core.Db.NewScope(&dbmodels.ConverterLog{}).HasColumn(Field) {
				criteria = criteria.Order("\"" + strings.ToLower(Field) + "\"" + " " + filter.OrderDirection[index])
			} else {
				err = errors.NewErrorWithCode("Ordering by unknown Field", errors.ErrorCodeNotValid, Field)
				return
			}
		}
	}

	q = criteria.Limit(limit).Offset(offset).Find(&dbmodelData)

	if q.Error != nil {
		log.Println("FindConverterLog > Ошибка получения данных2:", q.Error)
		return []types.ConverterLog{}, 0, nil
	}

	// подготовка id для получения связанных сущностей
	for _, item := range dbmodelData {
		foundIds = append(foundIds, item.ID)
	}

	// получение связнаных сущностей

	//формирование результатов
	for _, item := range dbmodelData {
		result = append(result, AssignConverterLogTypeFromDb(item))
	}

	return result, count, nil
}

func ConverterLogMultiCreate(filter types.ConverterLogFilter) (data []types.ConverterLog, err error) {

	typeModelList, err := filter.GetConverterLogModelList()

	if err != nil {
		return
	}

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetConverterLogModel(typeModel)
		item, e := ConverterLogCreate(filter, tx)

		if e != nil {
			err = e
			data = nil
			break
		}

		data = append(data, item)
	}

	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}

	return
}

func ConverterLogCreate(filter types.ConverterLogFilter, query *gorm.DB) (data types.ConverterLog, err error) {

	typeModel := filter.GetConverterLogModel()
	dbModel := AssignConverterLogDbFromType(typeModel)
	dbModel.ID = 0

	dbModel.Validate()

	if !dbModel.IsValid() {
		fmt.Println("ConverterLogCreate > Create ConverterLog error:", dbModel)
		return types.ConverterLog{}, dbModel.GetValidationError()
	}

	query = query.Create(&dbModel)

	if query.Error != nil {
		fmt.Println("ConverterLogCreate > Create ConverterLog error:", query.Error)
		return types.ConverterLog{}, errors.NewErrorWithCode("cant create ConverterLog", errors.ErrorCodeSqlError, "")
	}

	return AssignConverterLogTypeFromDb(dbModel), nil
}

func ConverterLogRead(filter types.ConverterLogFilter) (data types.ConverterLog, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1
	filter.ClearIds()
	filter.AddId(filter.GetCurrentId())

	findData, _, err := ConverterLogFind(filter)

	if len(findData) > 0 {
		return findData[0], nil
	}

	return types.ConverterLog{}, errors.NewErrorWithCode("Not found", errors.ErrorCodeNotFound, "")
}

func ConverterLogMultiUpdate(filter types.ConverterLogFilter) (data []types.ConverterLog, err error) {

	typeModelList, err := filter.GetConverterLogModelList()

	if err != nil {
		return
	}

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetConverterLogModel(typeModel)
		filter.ClearIds()
		filter.SetCurrentId(typeModel.Id)

		item, e := ConverterLogUpdate(filter, tx)

		if e != nil {
			err = e
			data = nil
			break
		}

		data = append(data, item)
	}

	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}

	return data, nil
}

func ConverterLogUpdate(filter types.ConverterLogFilter, query *gorm.DB) (data types.ConverterLog, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	existsModel, err := ConverterLogRead(filter)

	if existsModel.Id < 1 || err != nil {
		err = errors.NewErrorWithCode("ConverterLog not found in db with id: "+strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "Id")
		return
	}

	newModel := filter.GetConverterLogModel()

	updateModel := AssignConverterLogDbFromType(newModel)
	updateModel.ID = existsModel.Id

	//updateModel.Some = newModel.Some

	updateModel.FilePath = newModel.FilePath
	updateModel.ResultTextPath = newModel.ResultTextPath
	updateModel.ResultFilePath = newModel.ResultFilePath
	updateModel.ResultFormat = newModel.ResultFormat
	updateModel.RawResult = newModel.RawResult
	updateModel.ResultHtmlPath = newModel.ResultHtmlPath
	updateModel.ResultFileDocPath = newModel.ResultFileDocPath
	updateModel.ResultFilePdfPath = newModel.ResultFilePdfPath
	updateModel.UserId = newModel.UserId
	updateModel.SourceFilePath = newModel.SourceFilePath
	updateModel.RecordNumber = newModel.RecordNumber
	//updateModel.Field remove this line for disable generator functionality

	updateModel.Validate()

	if !updateModel.IsValid() {
		err = updateModel.GetValidationError()
		return
	}

	q := query.Model(dbmodels.ConverterLog{}).Save(&updateModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	data = AssignConverterLogTypeFromDb(updateModel)
	return
}

func ConverterLogMultiDelete(filter types.ConverterLogFilter) (isOk bool, err error) {

	typeModelList, err := filter.GetConverterLogModelList()

	if err != nil {
		return
	}

	isOk = true

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetConverterLogModel(typeModel)
		filter.ClearIds()
		filter.SetCurrentId(typeModel.Id)

		_, e := ConverterLogDelete(filter, tx)

		if e != nil {
			err = e
			isOk = false
			break
		}
	}

	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}

	return isOk, err
}

func ConverterLogDelete(filter types.ConverterLogFilter, query *gorm.DB) (isOk bool, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	existsModel, err := ConverterLogRead(filter)

	if existsModel.Id < 1 || err != nil {

		if err != nil {
			err = errors.NewErrorWithCode("ConverterLog not found in db with id: "+strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "")
		}
		return
	}

	dbModel := AssignConverterLogDbFromType(existsModel)
	q := query.Model(dbmodels.ConverterLog{}).Where(dbmodels.ConverterLog{ID: dbModel.ID}).Delete(&dbModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	isOk = true
	return
}

func ConverterLogFindOrCreate(filter types.ConverterLogFilter) (data types.ConverterLog, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	findOrCreateModel := AssignConverterLogDbFromType(filter.GetConverterLogModel())
	//findOrCreateModel.Field remove this line for disable generator functionality

	findOrCreateModel.Validate()

	if !findOrCreateModel.IsValid() {
		err = findOrCreateModel.GetValidationError()
		return
	}

	q := core.Db.Model(dbmodels.ConverterLog{}).Where(dbmodels.ConverterLog{ID: findOrCreateModel.ID}).FirstOrCreate(&findOrCreateModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	data = AssignConverterLogTypeFromDb(findOrCreateModel)
	return
}

func ConverterLogUpdateOrCreate(filter types.ConverterLogFilter) (data types.ConverterLog, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	updateOrCreateModel := AssignConverterLogDbFromType(filter.GetConverterLogModel())
	//updateOrCreateModel.Field remove this line for disable generator functionality

	updateOrCreateModel.Validate()

	if !updateOrCreateModel.IsValid() {
		err = updateOrCreateModel.GetValidationError()
		return
	}

	//please uncomment and set criteria
	//q := core.Db.Model(dbmodels.ConverterLog{}).Where(dbmodels.ConverterLog{ID: updateOrCreateModel.ID}).Assign(dbmodels.ConverterLog{/*PLEASE SET CRITERIA*/}).FirstOrCreate(&updateOrCreateModel)

	//if q.Error != nil {
	//    err = q.Error
	//    return
	//}

	data = AssignConverterLogTypeFromDb(updateOrCreateModel)
	return
}

// add all assign functions

func AssignConverterLogTypeFromDb(dbConverterLog dbmodels.ConverterLog) types.ConverterLog {

	//AssignConverterLogTypeFromDb predefine remove this line for disable generator functionality

	return types.ConverterLog{
		Id:                dbConverterLog.ID,
		FilePath:          dbConverterLog.FilePath,
		ResultTextPath:        dbConverterLog.ResultTextPath,
		ResultFilePath:    dbConverterLog.ResultFilePath,
		ResultFormat:      dbConverterLog.ResultFormat,
		RawResult:         dbConverterLog.RawResult,
		ResultHtmlPath:        dbConverterLog.ResultHtmlPath,
		ResultFileDocPath: dbConverterLog.ResultFileDocPath,
		ResultFilePdfPath: dbConverterLog.ResultFilePdfPath,
		UserId:            dbConverterLog.UserId,
		CreatedAt:         dbConverterLog.CreatedAt,
		SourceFilePath: dbConverterLog.SourceFilePath,
		RecordNumber: dbConverterLog.RecordNumber,
		//AssignConverterLogTypeFromDb.Field remove this line for disable generator functionality
	}
}

func AssignConverterLogDbFromType(typeModel types.ConverterLog) dbmodels.ConverterLog {

	//AssignConverterLogDbFromType predefine remove this line for disable generator functionality

	return dbmodels.ConverterLog{
		ID:                typeModel.Id,
		FilePath:          typeModel.FilePath,
		ResultTextPath:    typeModel.ResultTextPath,
		ResultFilePath:    typeModel.ResultFilePath,
		ResultFormat:      typeModel.ResultFormat,
		RawResult:         typeModel.RawResult,
		ResultHtmlPath:    typeModel.ResultHtmlPath,
		ResultFileDocPath: typeModel.ResultFileDocPath,
		ResultFilePdfPath: typeModel.ResultFilePdfPath,
		UserId:            typeModel.UserId,
		SourceFilePath:    typeModel.SourceFilePath,
		RecordNumber: typeModel.RecordNumber,
		//AssignConverterLogDbFromType.Field remove this line for disable generator functionality
	}
}
