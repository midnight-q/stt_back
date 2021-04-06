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

func RoleFind(filter types.RoleFilter) (result []types.Role, totalRecords int, err error) {

	foundIds := []int{}
	dbmodelData := []dbmodels.Role{}
	limit := filter.PerPage
	offset := filter.GetOffset()

	filterIds := filter.GetIds()
	filterExceptIds := filter.GetExceptIds()

	var count int

	criteria := core.Db.Where(dbmodels.Role{})

	//Role.FindFilterCode remove this line for disable generator functionality

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
	//            if core.Db.NewScope(&dbmodels.Role{}).HasColumn(field) {
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

	q := criteria.Model(dbmodels.Role{}).Count(&count)

	if q.Error != nil {
		log.Println("FindRole > Ошибка получения данных:", q.Error)
		return result, 0, nil
	}

	// order global criteria
	if len(filter.Order) > 0 {
		for index, Field := range filter.Order {
			if core.Db.NewScope(&dbmodels.Role{}).HasColumn(Field) {
				criteria = criteria.Order("\"" + strings.ToLower(Field) + "\"" + " " + filter.OrderDirection[index])
			} else {
				err = errors.NewErrorWithCode("Ordering by unknown Field", errors.ErrorCodeNotValid, Field)
				return
			}
		}
	}

	q = criteria.Limit(limit).Offset(offset).Find(&dbmodelData)

	if q.Error != nil {
		log.Println("FindRole > Ошибка получения данных2:", q.Error)
		return []types.Role{}, 0, nil
	}

	// подготовка id для получения связанных сущностей
	for _, item := range dbmodelData {
		foundIds = append(foundIds, item.ID)
	}

	// получение связнаных сущностей

	//формирование результатов
	for _, item := range dbmodelData {
		result = append(result, AssignRoleTypeFromDb(item))
	}

	return result, count, nil
}

func RoleMultiCreate(filter types.RoleFilter) (data []types.Role, err error) {

	typeModelList, err := filter.GetRoleModelList()

	if err != nil {
		return
	}

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetRoleModel(typeModel)
		item, e := RoleCreate(filter, tx)

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

func RoleCreate(filter types.RoleFilter, query *gorm.DB) (data types.Role, err error) {

	typeModel := filter.GetRoleModel()
	dbModel := AssignRoleDbFromType(typeModel)
	dbModel.ID = 0

	dbModel.Validate()

	if !dbModel.IsValid() {
		fmt.Println("RoleCreate > Create Role error:", dbModel)
		return types.Role{}, dbModel.GetValidationError()
	}

	query = query.Create(&dbModel)

	if query.Error != nil {
		fmt.Println("RoleCreate > Create Role error:", query.Error)
		return types.Role{}, errors.NewErrorWithCode("cant create Role", errors.ErrorCodeSqlError, "")
	}

	return AssignRoleTypeFromDb(dbModel), nil
}

func RoleRead(filter types.RoleFilter) (data types.Role, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1
	filter.ClearIds()
	filter.AddId(filter.GetCurrentId())

	findData, _, err := RoleFind(filter)

	if len(findData) > 0 {
		return findData[0], nil
	}

	return types.Role{}, errors.NewErrorWithCode("Not found", errors.ErrorCodeNotFound, "")
}

func RoleMultiUpdate(filter types.RoleFilter) (data []types.Role, err error) {

	typeModelList, err := filter.GetRoleModelList()

	if err != nil {
		return
	}

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetRoleModel(typeModel)
		filter.ClearIds()
		filter.SetCurrentId(typeModel.Id)

		item, e := RoleUpdate(filter, tx)

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

func RoleUpdate(filter types.RoleFilter, query *gorm.DB) (data types.Role, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	existsModel, err := RoleRead(filter)

	if existsModel.Id < 1 || err != nil {
		err = errors.NewErrorWithCode("Role not found in db with id: "+strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "Id")
		return
	}

	newModel := filter.GetRoleModel()

	updateModel := AssignRoleDbFromType(newModel)
	updateModel.ID = existsModel.Id

	//updateModel.Some = newModel.Some

	updateModel.Name = newModel.Name
	updateModel.Description = newModel.Description
	//updateModel.Field remove this line for disable generator functionality

	updateModel.Validate()

	if !updateModel.IsValid() {
		err = updateModel.GetValidationError()
		return
	}

	q := query.Model(dbmodels.Role{}).Save(&updateModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	data = AssignRoleTypeFromDb(updateModel)
	return
}

func RoleMultiDelete(filter types.RoleFilter) (isOk bool, err error) {

	typeModelList, err := filter.GetRoleModelList()

	if err != nil {
		return
	}

	isOk = true

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetRoleModel(typeModel)
		filter.ClearIds()
		filter.SetCurrentId(typeModel.Id)

		_, e := RoleDelete(filter, tx)

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

func RoleDelete(filter types.RoleFilter, query *gorm.DB) (isOk bool, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	existsModel, err := RoleRead(filter)

	if existsModel.Id < 1 || err != nil {

		if err != nil {
			err = errors.NewErrorWithCode("Role not found in db with id: "+strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "")
		}
		return
	}

	dbModel := AssignRoleDbFromType(existsModel)
	q := query.Model(dbmodels.Role{}).Where(dbmodels.Role{ID: dbModel.ID}).Delete(&dbModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	isOk = true
	return
}

func RoleFindOrCreate(filter types.RoleFilter) (data types.Role, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	findOrCreateModel := AssignRoleDbFromType(filter.GetRoleModel())
	//findOrCreateModel.Field remove this line for disable generator functionality

	findOrCreateModel.Validate()

	if !findOrCreateModel.IsValid() {
		err = findOrCreateModel.GetValidationError()
		return
	}

	q := core.Db.Model(dbmodels.Role{}).Where(dbmodels.Role{ID: findOrCreateModel.ID}).FirstOrCreate(&findOrCreateModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	data = AssignRoleTypeFromDb(findOrCreateModel)
	return
}

func RoleUpdateOrCreate(filter types.RoleFilter) (data types.Role, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	updateOrCreateModel := AssignRoleDbFromType(filter.GetRoleModel())
	//updateOrCreateModel.Field remove this line for disable generator functionality

	updateOrCreateModel.Validate()

	if !updateOrCreateModel.IsValid() {
		err = updateOrCreateModel.GetValidationError()
		return
	}

	//please uncomment and set criteria
	//q := core.Db.Model(dbmodels.Role{}).Where(dbmodels.Role{ID: updateOrCreateModel.ID}).Assign(dbmodels.Role{/*PLEASE SET CRITERIA*/}).FirstOrCreate(&updateOrCreateModel)

	//if q.Error != nil {
	//    err = q.Error
	//    return
	//}

	data = AssignRoleTypeFromDb(updateOrCreateModel)
	return
}

// add all assign functions

func AssignRoleTypeFromDb(dbRole dbmodels.Role) types.Role {

	//AssignRoleTypeFromDb predefine remove this line for disable generator functionality

	return types.Role{
		Id:          dbRole.ID,
		Name:        dbRole.Name,
		Description: dbRole.Description,
		//AssignRoleTypeFromDb.Field remove this line for disable generator functionality
	}
}

func AssignRoleDbFromType(typeModel types.Role) dbmodels.Role {

	//AssignRoleDbFromType predefine remove this line for disable generator functionality

	return dbmodels.Role{
		ID:          typeModel.Id,
		Name:        typeModel.Name,
		Description: typeModel.Description,
		//AssignRoleDbFromType.Field remove this line for disable generator functionality
	}
}
