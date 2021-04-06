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

func LanguageFind(filter types.LanguageFilter) (result []types.Language, totalRecords int, err error) {

	foundIds := []int{}
	dbmodelData := []dbmodels.Language{}
	limit := filter.PerPage
	offset := filter.GetOffset()

	filterIds := filter.GetIds()
	filterExceptIds := filter.GetExceptIds()

	var count int

	criteria := core.Db.Where(dbmodels.Language{})

	//Language.FindFilterCode remove this line for disable generator functionality

	if len(filterIds) > 0 {
		criteria = criteria.Where("id in (?)", filterIds)
	}

	if len(filterExceptIds) > 0 {
		criteria = criteria.Where("id not in (?)", filterExceptIds)
	}

	//if len(filter.Search) > 0 {
	//
	//    s := ("%" + filter.Search + "%")
	//
	//    if len(filter.SearchBy) > 0 {
	//
	//        for _, field := range filter.SearchBy {
	//
	//            if core.Db.NewScope(&dbmodels.Language{}).HasColumn(field) {
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

	q := criteria.Model(dbmodels.Language{}).Count(&count)

	if q.Error != nil {
		log.Println("FindLanguage > Ошибка получения данных:", q.Error)
		return result, 0, nil
	}

	// order global criteria
	if len(filter.Order) > 0 {
		for index, Field := range filter.Order {
			if core.Db.NewScope(&dbmodels.Language{}).HasColumn(Field) {
				criteria = criteria.Order("\"" + strings.ToLower(Field) + "\"" + " " + filter.OrderDirection[index])
			} else {
				err = errors.NewErrorWithCode("Ordering by unknown Field", errors.ErrorCodeNotValid, Field)
				return
			}
		}
	}

	q = criteria.Limit(limit).Offset(offset).Find(&dbmodelData)

	if q.Error != nil {
		log.Println("FindLanguage > Ошибка получения данных2:", q.Error)
		return []types.Language{}, 0, nil
	}

	// подготовка id для получения связанных сущностей
	for _, item := range dbmodelData {
		foundIds = append(foundIds, item.ID)
	}

	// получение связнаных сущностей

	//формирование результатов
	for _, item := range dbmodelData {
		result = append(result, AssignLanguageTypeFromDb(item))
	}

	return result, count, nil
}

func LanguageMultiCreate(filter types.LanguageFilter) (data []types.Language, err error) {

	typeModelList, err := filter.GetLanguageModelList()

	if err != nil {
		return
	}

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetLanguageModel(typeModel)
		item, e := LanguageCreate(filter, tx)

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

func LanguageCreate(filter types.LanguageFilter, query *gorm.DB) (data types.Language, err error) {

	typeModel := filter.GetLanguageModel()
	dbModel := AssignLanguageDbFromType(typeModel)
	dbModel.ID = 0

	dbModel.Validate()

	if !dbModel.IsValid() {
		fmt.Println("LanguageCreate > Create Language error:", dbModel)
		return types.Language{}, dbModel.GetValidationError()
	}

	query = query.Create(&dbModel)

	if query.Error != nil {
		fmt.Println("LanguageCreate > Create Language error:", query.Error)
		return types.Language{}, errors.NewErrorWithCode("cant create Language", errors.ErrorCodeSqlError, "")
	}

	return AssignLanguageTypeFromDb(dbModel), nil
}

func LanguageRead(filter types.LanguageFilter) (data types.Language, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1
	filter.ClearIds()
	filter.AddId(filter.GetCurrentId())

	findData, _, err := LanguageFind(filter)

	if len(findData) > 0 {
		return findData[0], nil
	}

	return types.Language{}, errors.NewErrorWithCode("Not found", errors.ErrorCodeNotFound, "")
}

func LanguageMultiUpdate(filter types.LanguageFilter) (data []types.Language, err error) {

	typeModelList, err := filter.GetLanguageModelList()

	if err != nil {
		return
	}

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetLanguageModel(typeModel)
		filter.ClearIds()
		filter.SetCurrentId(typeModel.Id)

		item, e := LanguageUpdate(filter, tx)

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

func LanguageUpdate(filter types.LanguageFilter, query *gorm.DB) (data types.Language, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	existsModel, err := LanguageRead(filter)

	if existsModel.Id < 1 || err != nil {
		err = errors.NewErrorWithCode("Language not found in db with id: "+strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "Id")
		return
	}

	newModel := filter.GetLanguageModel()

	updateModel := AssignLanguageDbFromType(newModel)
	updateModel.ID = existsModel.Id

	//updateModel.Some = newModel.Some

	updateModel.Name = newModel.Name
	updateModel.Code = newModel.Code
	//updateModel.Field remove this line for disable generator functionality

	updateModel.Validate()

	if !updateModel.IsValid() {
		err = updateModel.GetValidationError()
		return
	}

	q := query.Model(dbmodels.Language{}).Save(&updateModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	data = AssignLanguageTypeFromDb(updateModel)
	return
}

func LanguageMultiDelete(filter types.LanguageFilter) (isOk bool, err error) {

	typeModelList, err := filter.GetLanguageModelList()

	if err != nil {
		return
	}

	isOk = true

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetLanguageModel(typeModel)
		filter.ClearIds()
		filter.SetCurrentId(typeModel.Id)

		_, e := LanguageDelete(filter, tx)

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

func LanguageDelete(filter types.LanguageFilter, query *gorm.DB) (isOk bool, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	existsModel, err := LanguageRead(filter)

	if existsModel.Id < 1 || err != nil {

		if err != nil {
			err = errors.NewErrorWithCode("Language not found in db with id: "+strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "")
		}
		return
	}

	dbModel := AssignLanguageDbFromType(existsModel)
	q := query.Model(dbmodels.Language{}).Where(dbmodels.Language{ID: dbModel.ID}).Delete(&dbModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	isOk = true
	return
}

func LanguageFindOrCreate(filter types.LanguageFilter) (data types.Language, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	findOrCreateModel := AssignLanguageDbFromType(filter.GetLanguageModel())
	//findOrCreateModel.Field remove this line for disable generator functionality

	findOrCreateModel.Validate()

	if !findOrCreateModel.IsValid() {
		err = findOrCreateModel.GetValidationError()
		return
	}

	q := core.Db.Model(dbmodels.Language{}).Where(dbmodels.Language{ID: findOrCreateModel.ID}).FirstOrCreate(&findOrCreateModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	data = AssignLanguageTypeFromDb(findOrCreateModel)
	return
}

func LanguageUpdateOrCreate(filter types.LanguageFilter) (data types.Language, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	updateOrCreateModel := AssignLanguageDbFromType(filter.GetLanguageModel())
	//updateOrCreateModel.Field remove this line for disable generator functionality

	updateOrCreateModel.Validate()

	if !updateOrCreateModel.IsValid() {
		err = updateOrCreateModel.GetValidationError()
		return
	}

	//please uncomment and set criteria
	//q := core.Db.Model(dbmodels.Language{}).Where(dbmodels.Language{ID: updateOrCreateModel.ID}).Assign(dbmodels.Language{/*PLEASE SET CRITERIA*/}).FirstOrCreate(&updateOrCreateModel)

	//if q.Error != nil {
	//    err = q.Error
	//    return
	//}

	data = AssignLanguageTypeFromDb(updateOrCreateModel)
	return
}

// add all assign functions

func AssignLanguageTypeFromDb(dbLanguage dbmodels.Language) types.Language {

	//AssignLanguageTypeFromDb predefine remove this line for disable generator functionality

	return types.Language{
		Id:   dbLanguage.ID,
		Name: dbLanguage.Name,
		Code: dbLanguage.Code,
		//AssignLanguageTypeFromDb.Field remove this line for disable generator functionality
	}
}

func AssignLanguageDbFromType(typeModel types.Language) dbmodels.Language {

	//AssignLanguageDbFromType predefine remove this line for disable generator functionality

	return dbmodels.Language{
		ID:   typeModel.Id,
		Name: typeModel.Name,
		Code: typeModel.Code,
		//AssignLanguageDbFromType.Field remove this line for disable generator functionality
	}
}
