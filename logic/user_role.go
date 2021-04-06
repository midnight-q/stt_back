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

func UserRoleFind(filter types.UserRoleFilter) (result []types.UserRole, totalRecords int, err error) {

	foundIds := []int{}
	dbmodelData := []dbmodels.UserRole{}
	limit := filter.PerPage
	offset := filter.GetOffset()

	filterIds := filter.GetIds()
	filterExceptIds := filter.GetExceptIds()

	var count int

	criteria := core.Db.Where(dbmodels.UserRole{})

	//UserRole.FindFilterCode remove this line for disable generator functionality

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
	//            if core.Db.NewScope(&dbmodels.UserRole{}).HasColumn(field) {
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

	q := criteria.Model(dbmodels.UserRole{}).Count(&count)

	if q.Error != nil {
		log.Println("FindUserRole > Ошибка получения данных:", q.Error)
		return result, 0, nil
	}

	// order global criteria
	if len(filter.Order) > 0 {
		for index, Field := range filter.Order {
			if core.Db.NewScope(&dbmodels.UserRole{}).HasColumn(Field) {
				criteria = criteria.Order("\"" + strings.ToLower(Field) + "\"" + " " + filter.OrderDirection[index])
			} else {
				err = errors.NewErrorWithCode("Ordering by unknown Field", errors.ErrorCodeNotValid, Field)
				return
			}
		}
	}

	q = criteria.Limit(limit).Offset(offset).Find(&dbmodelData)

	if q.Error != nil {
		log.Println("FindUserRole > Ошибка получения данных2:", q.Error)
		return []types.UserRole{}, 0, nil
	}

	// подготовка id для получения связанных сущностей
	for _, item := range dbmodelData {
		foundIds = append(foundIds, item.ID)
	}

	// получение связнаных сущностей

	//формирование результатов
	for _, item := range dbmodelData {
		result = append(result, AssignUserRoleTypeFromDb(item))
	}

	return result, count, nil
}

func UserRoleMultiCreate(filter types.UserRoleFilter) (data []types.UserRole, err error) {

	typeModelList, err := filter.GetUserRoleModelList()

	if err != nil {
		return
	}

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetUserRoleModel(typeModel)
		item, e := UserRoleCreate(filter, tx)

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

func UserRoleCreate(filter types.UserRoleFilter, query *gorm.DB) (data types.UserRole, err error) {

	typeModel := filter.GetUserRoleModel()
	dbModel := AssignUserRoleDbFromType(typeModel)
	dbModel.ID = 0

	dbModel.Validate()

	if !dbModel.IsValid() {
		fmt.Println("UserRoleCreate > Create UserRole error:", dbModel)
		return types.UserRole{}, dbModel.GetValidationError()
	}

	query = query.Create(&dbModel)

	if query.Error != nil {
		fmt.Println("UserRoleCreate > Create UserRole error:", query.Error)
		return types.UserRole{}, errors.NewErrorWithCode("cant create UserRole", errors.ErrorCodeSqlError, "")
	}

	return AssignUserRoleTypeFromDb(dbModel), nil
}

func UserRoleRead(filter types.UserRoleFilter) (data types.UserRole, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1
	filter.ClearIds()
	filter.AddId(filter.GetCurrentId())

	findData, _, err := UserRoleFind(filter)

	if len(findData) > 0 {
		return findData[0], nil
	}

	return types.UserRole{}, errors.NewErrorWithCode("Not found", errors.ErrorCodeNotFound, "")
}

func UserRoleMultiUpdate(filter types.UserRoleFilter) (data []types.UserRole, err error) {

	typeModelList, err := filter.GetUserRoleModelList()

	if err != nil {
		return
	}

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetUserRoleModel(typeModel)
		filter.ClearIds()
		filter.SetCurrentId(typeModel.Id)

		item, e := UserRoleUpdate(filter, tx)

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

func UserRoleUpdate(filter types.UserRoleFilter, query *gorm.DB) (data types.UserRole, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	existsModel, err := UserRoleRead(filter)

	if existsModel.Id < 1 || err != nil {
		err = errors.NewErrorWithCode("UserRole not found in db with id: "+strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "Id")
		return
	}

	newModel := filter.GetUserRoleModel()

	updateModel := AssignUserRoleDbFromType(newModel)
	updateModel.ID = existsModel.Id

	//updateModel.Some = newModel.Some

	updateModel.UserId = newModel.UserId
	updateModel.RoleId = newModel.RoleId
	//updateModel.Field remove this line for disable generator functionality

	updateModel.Validate()

	if !updateModel.IsValid() {
		err = updateModel.GetValidationError()
		return
	}

	q := query.Model(dbmodels.UserRole{}).Save(&updateModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	data = AssignUserRoleTypeFromDb(updateModel)
	return
}

func UserRoleMultiDelete(filter types.UserRoleFilter) (isOk bool, err error) {

	typeModelList, err := filter.GetUserRoleModelList()

	if err != nil {
		return
	}

	isOk = true

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetUserRoleModel(typeModel)
		filter.ClearIds()
		filter.SetCurrentId(typeModel.Id)

		_, e := UserRoleDelete(filter, tx)

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

func UserRoleDelete(filter types.UserRoleFilter, query *gorm.DB) (isOk bool, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	existsModel, err := UserRoleRead(filter)

	if existsModel.Id < 1 || err != nil {

		if err != nil {
			err = errors.NewErrorWithCode("UserRole not found in db with id: "+strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "")
		}
		return
	}

	dbModel := AssignUserRoleDbFromType(existsModel)
	q := query.Model(dbmodels.UserRole{}).Where(dbmodels.UserRole{ID: dbModel.ID}).Delete(&dbModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	isOk = true
	return
}

func UserRoleFindOrCreate(filter types.UserRoleFilter) (data types.UserRole, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	findOrCreateModel := AssignUserRoleDbFromType(filter.GetUserRoleModel())
	//findOrCreateModel.Field remove this line for disable generator functionality

	findOrCreateModel.Validate()

	if !findOrCreateModel.IsValid() {
		err = findOrCreateModel.GetValidationError()
		return
	}

	q := core.Db.Model(dbmodels.UserRole{}).Where(dbmodels.UserRole{ID: findOrCreateModel.ID}).FirstOrCreate(&findOrCreateModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	data = AssignUserRoleTypeFromDb(findOrCreateModel)
	return
}

func UserRoleUpdateOrCreate(filter types.UserRoleFilter) (data types.UserRole, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	updateOrCreateModel := AssignUserRoleDbFromType(filter.GetUserRoleModel())
	//updateOrCreateModel.Field remove this line for disable generator functionality

	updateOrCreateModel.Validate()

	if !updateOrCreateModel.IsValid() {
		err = updateOrCreateModel.GetValidationError()
		return
	}

	//please uncomment and set criteria
	//q := core.Db.Model(dbmodels.UserRole{}).Where(dbmodels.UserRole{ID: updateOrCreateModel.ID}).Assign(dbmodels.UserRole{/*PLEASE SET CRITERIA*/}).FirstOrCreate(&updateOrCreateModel)

	//if q.Error != nil {
	//    err = q.Error
	//    return
	//}

	data = AssignUserRoleTypeFromDb(updateOrCreateModel)
	return
}

// add all assign functions

func AssignUserRoleTypeFromDb(dbUserRole dbmodels.UserRole) types.UserRole {

	//AssignUserRoleTypeFromDb predefine remove this line for disable generator functionality

	return types.UserRole{
		Id:     dbUserRole.ID,
		UserId: dbUserRole.UserId,
		RoleId: dbUserRole.RoleId,
		//AssignUserRoleTypeFromDb.Field remove this line for disable generator functionality
	}
}

func AssignUserRoleDbFromType(typeModel types.UserRole) dbmodels.UserRole {

	//AssignUserRoleDbFromType predefine remove this line for disable generator functionality

	return dbmodels.UserRole{
		ID:     typeModel.Id,
		UserId: typeModel.UserId,
		RoleId: typeModel.RoleId,
		//AssignUserRoleDbFromType.Field remove this line for disable generator functionality
	}
}
