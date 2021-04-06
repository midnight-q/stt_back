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

func CurrentUserFind(filter types.CurrentUserFilter) (result []types.CurrentUser, totalRecords int, err error) {

	foundIds := []int{}
	dbmodelData := []dbmodels.CurrentUser{}
	limit := filter.PerPage
	offset := filter.GetOffset()

	filterIds := filter.GetIds()
	filterExceptIds := filter.GetExceptIds()

	var count int

	criteria := core.Db.Where(dbmodels.CurrentUser{})

	//CurrentUser.FindFilterCode remove this line for disable generator functionality

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
	//            if core.Db.NewScope(&dbmodels.CurrentUser{}).HasColumn(field) {
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

	q := criteria.Model(dbmodels.CurrentUser{}).Count(&count)

	if q.Error != nil {
		log.Println("FindCurrentUser > Ошибка получения данных:", q.Error)
		return result, 0, nil
	}

	// order global criteria
	if len(filter.Order) > 0 {
		for index, Field := range filter.Order {
			if core.Db.NewScope(&dbmodels.CurrentUser{}).HasColumn(Field) {
				criteria = criteria.Order("\"" + strings.ToLower(Field) + "\"" + " " + filter.OrderDirection[index])
			} else {
				err = errors.NewErrorWithCode("Ordering by unknown Field", errors.ErrorCodeNotValid, Field)
				return
			}
		}
	}

	q = criteria.Limit(limit).Offset(offset).Find(&dbmodelData)

	if q.Error != nil {
		log.Println("FindCurrentUser > Ошибка получения данных2:", q.Error)
		return []types.CurrentUser{}, 0, nil
	}

	// подготовка id для получения связанных сущностей
	for _, item := range dbmodelData {
		foundIds = append(foundIds, item.ID)
	}

	// получение связнаных сущностей

	//формирование результатов
	for _, item := range dbmodelData {
		result = append(result, AssignCurrentUserTypeFromDb(item))
	}

	return result, count, nil
}

func CurrentUserMultiCreate(filter types.CurrentUserFilter) (data []types.CurrentUser, err error) {

	typeModelList, err := filter.GetCurrentUserModelList()

	if err != nil {
		return
	}

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetCurrentUserModel(typeModel)
		item, e := CurrentUserCreate(filter, tx)

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

func CurrentUserCreate(filter types.CurrentUserFilter, query *gorm.DB) (data types.CurrentUser, err error) {

	typeModel := filter.GetCurrentUserModel()
	dbModel := AssignCurrentUserDbFromType(typeModel)
	dbModel.ID = 0

	dbModel.Validate()

	if !dbModel.IsValid() {
		fmt.Println("CurrentUserCreate > Create CurrentUser error:", dbModel)
		return types.CurrentUser{}, dbModel.GetValidationError()
	}

	query = query.Create(&dbModel)

	if query.Error != nil {
		fmt.Println("CurrentUserCreate > Create CurrentUser error:", query.Error)
		return types.CurrentUser{}, errors.NewErrorWithCode("cant create CurrentUser", errors.ErrorCodeSqlError, "")
	}

	return AssignCurrentUserTypeFromDb(dbModel), nil
}

func CurrentUserRead(filter types.CurrentUserFilter) (data types.CurrentUser, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1
	filter.ClearIds()
	filter.AddId(filter.GetCurrentId())

	findData, _, err := CurrentUserFind(filter)

	if len(findData) > 0 {
		return findData[0], nil
	}

	return types.CurrentUser{}, errors.NewErrorWithCode("Not found", errors.ErrorCodeNotFound, "")
}

func CurrentUserMultiUpdate(filter types.CurrentUserFilter) (data []types.CurrentUser, err error) {

	typeModelList, err := filter.GetCurrentUserModelList()

	if err != nil {
		return
	}

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetCurrentUserModel(typeModel)
		filter.ClearIds()
		filter.SetCurrentId(typeModel.Id)

		item, e := CurrentUserUpdate(filter, tx)

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

func CurrentUserUpdate(filter types.CurrentUserFilter, query *gorm.DB) (data types.CurrentUser, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	existsModel, err := CurrentUserRead(filter)

	if existsModel.Id < 1 || err != nil {
		err = errors.NewErrorWithCode("CurrentUser not found in db with id: "+strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "Id")
		return
	}

	newModel := filter.GetCurrentUserModel()

	updateModel := AssignCurrentUserDbFromType(newModel)
	updateModel.ID = existsModel.Id

	//updateModel.Some = newModel.Some

	//updateModel.Field remove this line for disable generator functionality

	updateModel.Validate()

	if !updateModel.IsValid() {
		err = updateModel.GetValidationError()
		return
	}

	q := query.Model(dbmodels.CurrentUser{}).Save(&updateModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	data = AssignCurrentUserTypeFromDb(updateModel)
	return
}

func CurrentUserMultiDelete(filter types.CurrentUserFilter) (isOk bool, err error) {

	typeModelList, err := filter.GetCurrentUserModelList()

	if err != nil {
		return
	}

	isOk = true

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetCurrentUserModel(typeModel)
		filter.ClearIds()
		filter.SetCurrentId(typeModel.Id)

		_, e := CurrentUserDelete(filter, tx)

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

func CurrentUserDelete(filter types.CurrentUserFilter, query *gorm.DB) (isOk bool, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	existsModel, err := CurrentUserRead(filter)

	if existsModel.Id < 1 || err != nil {

		if err != nil {
			err = errors.NewErrorWithCode("CurrentUser not found in db with id: "+strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "")
		}
		return
	}

	dbModel := AssignCurrentUserDbFromType(existsModel)
	q := query.Model(dbmodels.CurrentUser{}).Where(dbmodels.CurrentUser{ID: dbModel.ID}).Delete(&dbModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	isOk = true
	return
}

func CurrentUserFindOrCreate(filter types.CurrentUserFilter) (data types.CurrentUser, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	findOrCreateModel := AssignCurrentUserDbFromType(filter.GetCurrentUserModel())
	//findOrCreateModel.Field remove this line for disable generator functionality

	findOrCreateModel.Validate()

	if !findOrCreateModel.IsValid() {
		err = findOrCreateModel.GetValidationError()
		return
	}

	q := core.Db.Model(dbmodels.CurrentUser{}).Where(dbmodels.CurrentUser{ID: findOrCreateModel.ID}).FirstOrCreate(&findOrCreateModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	data = AssignCurrentUserTypeFromDb(findOrCreateModel)
	return
}

func CurrentUserUpdateOrCreate(filter types.CurrentUserFilter) (data types.CurrentUser, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	updateOrCreateModel := AssignCurrentUserDbFromType(filter.GetCurrentUserModel())
	//updateOrCreateModel.Field remove this line for disable generator functionality

	updateOrCreateModel.Validate()

	if !updateOrCreateModel.IsValid() {
		err = updateOrCreateModel.GetValidationError()
		return
	}

	//please uncomment and set criteria
	//q := core.Db.Model(dbmodels.CurrentUser{}).Where(dbmodels.CurrentUser{ID: updateOrCreateModel.ID}).Assign(dbmodels.CurrentUser{/*PLEASE SET CRITERIA*/}).FirstOrCreate(&updateOrCreateModel)

	//if q.Error != nil {
	//    err = q.Error
	//    return
	//}

	data = AssignCurrentUserTypeFromDb(updateOrCreateModel)
	return
}

// add all assign functions

func AssignCurrentUserTypeFromDb(dbCurrentUser dbmodels.CurrentUser) types.CurrentUser {

	//AssignCurrentUserTypeFromDb predefine remove this line for disable generator functionality

	return types.CurrentUser{
		Id: dbCurrentUser.ID,
		//AssignCurrentUserTypeFromDb.Field remove this line for disable generator functionality
	}
}

func AssignCurrentUserDbFromType(typeModel types.CurrentUser) dbmodels.CurrentUser {

	//AssignCurrentUserDbFromType predefine remove this line for disable generator functionality

	return dbmodels.CurrentUser{
		ID: typeModel.Id,
		//AssignCurrentUserDbFromType.Field remove this line for disable generator functionality
	}
}
