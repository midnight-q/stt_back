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

func RoleResourceFind(filter types.RoleResourceFilter) (result []types.RoleResource, totalRecords int, err error) {

	foundIds := []int{}
	dbmodelData := []dbmodels.RoleResource{}
	limit := filter.PerPage
	offset := filter.GetOffset()

	filterIds := filter.GetIds()
	filterExceptIds := filter.GetExceptIds()

	var count int

	criteria := core.Db.Where(dbmodels.RoleResource{})

	//RoleResource.FindFilterCode remove this line for disable generator functionality

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
	//            if core.Db.NewScope(&dbmodels.RoleResource{}).HasColumn(field) {
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

	q := criteria.Model(dbmodels.RoleResource{}).Count(&count)

	if q.Error != nil {
		log.Println("FindRoleResource > Ошибка получения данных:", q.Error)
		return result, 0, nil
	}

	// order global criteria
	if len(filter.Order) > 0 {
		for index, Field := range filter.Order {
			if core.Db.NewScope(&dbmodels.RoleResource{}).HasColumn(Field) {
				criteria = criteria.Order("\"" + strings.ToLower(Field) + "\"" + " " + filter.OrderDirection[index])
			} else {
				err = errors.NewErrorWithCode("Ordering by unknown Field", errors.ErrorCodeNotValid, Field)
				return
			}
		}
	}

	q = criteria.Limit(limit).Offset(offset).Find(&dbmodelData)

	if q.Error != nil {
		log.Println("FindRoleResource > Ошибка получения данных2:", q.Error)
		return []types.RoleResource{}, 0, nil
	}

	// подготовка id для получения связанных сущностей
	for _, item := range dbmodelData {
		foundIds = append(foundIds, item.ID)
	}

	// получение связнаных сущностей

	//формирование результатов
	for _, item := range dbmodelData {
		result = append(result, AssignRoleResourceTypeFromDb(item))
	}

	return result, count, nil
}

func RoleResourceMultiCreate(filter types.RoleResourceFilter) (data []types.RoleResource, err error) {

	typeModelList, err := filter.GetRoleResourceModelList()

	if err != nil {
		return
	}

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetRoleResourceModel(typeModel)
		item, e := RoleResourceCreate(filter, tx)

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

func RoleResourceCreate(filter types.RoleResourceFilter, query *gorm.DB) (data types.RoleResource, err error) {

	typeModel := filter.GetRoleResourceModel()
	dbModel := AssignRoleResourceDbFromType(typeModel)
	dbModel.ID = 0

	dbModel.Validate()

	if !dbModel.IsValid() {
		fmt.Println("RoleResourceCreate > Create RoleResource error:", dbModel)
		return types.RoleResource{}, dbModel.GetValidationError()
	}

	query = query.Create(&dbModel)

	if query.Error != nil {
		fmt.Println("RoleResourceCreate > Create RoleResource error:", query.Error)
		return types.RoleResource{}, errors.NewErrorWithCode("cant create RoleResource", errors.ErrorCodeSqlError, "")
	}

	return AssignRoleResourceTypeFromDb(dbModel), nil
}

func RoleResourceRead(filter types.RoleResourceFilter) (data types.RoleResource, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1
	filter.ClearIds()
	filter.AddId(filter.GetCurrentId())

	findData, _, err := RoleResourceFind(filter)

	if len(findData) > 0 {
		return findData[0], nil
	}

	return types.RoleResource{}, errors.NewErrorWithCode("Not found", errors.ErrorCodeNotFound, "")
}

func RoleResourceMultiUpdate(filter types.RoleResourceFilter) (data []types.RoleResource, err error) {

	typeModelList, err := filter.GetRoleResourceModelList()

	if err != nil {
		return
	}

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetRoleResourceModel(typeModel)
		filter.ClearIds()
		filter.SetCurrentId(typeModel.Id)

		item, e := RoleResourceUpdate(filter, tx)

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

func RoleResourceUpdate(filter types.RoleResourceFilter, query *gorm.DB) (data types.RoleResource, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	existsModel, err := RoleResourceRead(filter)

	if existsModel.Id < 1 || err != nil {
		err = errors.NewErrorWithCode("RoleResource not found in db with id: "+strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "Id")
		return
	}

	newModel := filter.GetRoleResourceModel()

	updateModel := AssignRoleResourceDbFromType(newModel)
	updateModel.ID = existsModel.Id

	//updateModel.Some = newModel.Some

	updateModel.RoleId = newModel.RoleId
	updateModel.ResourceId = newModel.ResourceId
	updateModel.Find = newModel.Find
	updateModel.Read = newModel.Read
	updateModel.Create = newModel.Create
	updateModel.Update = newModel.Update
	updateModel.Delete = newModel.Delete
	updateModel.FindOrCreate = newModel.FindOrCreate
	updateModel.UpdateOrCreate = newModel.UpdateOrCreate
	//updateModel.Field remove this line for disable generator functionality

	updateModel.Validate()

	if !updateModel.IsValid() {
		err = updateModel.GetValidationError()
		return
	}

	q := query.Model(dbmodels.RoleResource{}).Save(&updateModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	data = AssignRoleResourceTypeFromDb(updateModel)
	return
}

func RoleResourceMultiDelete(filter types.RoleResourceFilter) (isOk bool, err error) {

	typeModelList, err := filter.GetRoleResourceModelList()

	if err != nil {
		return
	}

	isOk = true

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetRoleResourceModel(typeModel)
		filter.ClearIds()
		filter.SetCurrentId(typeModel.Id)

		_, e := RoleResourceDelete(filter, tx)

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

func RoleResourceDelete(filter types.RoleResourceFilter, query *gorm.DB) (isOk bool, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	existsModel, err := RoleResourceRead(filter)

	if existsModel.Id < 1 || err != nil {

		if err != nil {
			err = errors.NewErrorWithCode("RoleResource not found in db with id: "+strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "")
		}
		return
	}

	dbModel := AssignRoleResourceDbFromType(existsModel)
	q := query.Model(dbmodels.RoleResource{}).Where(dbmodels.RoleResource{ID: dbModel.ID}).Delete(&dbModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	isOk = true
	return
}

func RoleResourceFindOrCreate(filter types.RoleResourceFilter) (data types.RoleResource, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	findOrCreateModel := AssignRoleResourceDbFromType(filter.GetRoleResourceModel())
	//findOrCreateModel.Field remove this line for disable generator functionality

	findOrCreateModel.Validate()

	if !findOrCreateModel.IsValid() {
		err = findOrCreateModel.GetValidationError()
		return
	}

	q := core.Db.Model(dbmodels.RoleResource{}).Where(dbmodels.RoleResource{ID: findOrCreateModel.ID}).FirstOrCreate(&findOrCreateModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	data = AssignRoleResourceTypeFromDb(findOrCreateModel)
	return
}

func RoleResourceUpdateOrCreate(filter types.RoleResourceFilter) (data types.RoleResource, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	updateOrCreateModel := AssignRoleResourceDbFromType(filter.GetRoleResourceModel())
	//updateOrCreateModel.Field remove this line for disable generator functionality

	updateOrCreateModel.Validate()

	if !updateOrCreateModel.IsValid() {
		err = updateOrCreateModel.GetValidationError()
		return
	}

	//please uncomment and set criteria
	//q := core.Db.Model(dbmodels.RoleResource{}).Where(dbmodels.RoleResource{ID: updateOrCreateModel.ID}).Assign(dbmodels.RoleResource{/*PLEASE SET CRITERIA*/}).FirstOrCreate(&updateOrCreateModel)

	//if q.Error != nil {
	//    err = q.Error
	//    return
	//}

	data = AssignRoleResourceTypeFromDb(updateOrCreateModel)
	return
}

// add all assign functions

func AssignRoleResourceTypeFromDb(dbRoleResource dbmodels.RoleResource) types.RoleResource {

	//AssignRoleResourceTypeFromDb predefine remove this line for disable generator functionality

	return types.RoleResource{
		Id:             dbRoleResource.ID,
		RoleId:         dbRoleResource.RoleId,
		ResourceId:     dbRoleResource.ResourceId,
		Find:           dbRoleResource.Find,
		Read:           dbRoleResource.Read,
		Create:         dbRoleResource.Create,
		Update:         dbRoleResource.Update,
		Delete:         dbRoleResource.Delete,
		FindOrCreate:   dbRoleResource.FindOrCreate,
		UpdateOrCreate: dbRoleResource.UpdateOrCreate,
		//AssignRoleResourceTypeFromDb.Field remove this line for disable generator functionality
	}
}

func AssignRoleResourceDbFromType(typeModel types.RoleResource) dbmodels.RoleResource {

	//AssignRoleResourceDbFromType predefine remove this line for disable generator functionality

	return dbmodels.RoleResource{
		ID:             typeModel.Id,
		RoleId:         typeModel.RoleId,
		ResourceId:     typeModel.ResourceId,
		Find:           typeModel.Find,
		Read:           typeModel.Read,
		Create:         typeModel.Create,
		Update:         typeModel.Update,
		Delete:         typeModel.Delete,
		FindOrCreate:   typeModel.FindOrCreate,
		UpdateOrCreate: typeModel.UpdateOrCreate,
		//AssignRoleResourceDbFromType.Field remove this line for disable generator functionality
	}
}
