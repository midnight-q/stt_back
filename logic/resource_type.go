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

func ResourceTypeFind(filter types.ResourceTypeFilter) (result []types.ResourceType, totalRecords int, err error) {

	foundIds := []int{}
	dbmodelData := []dbmodels.ResourceType{}
	limit := filter.PerPage
	offset := filter.GetOffset()

	filterIds := filter.GetIds()
	filterExceptIds := filter.GetExceptIds()

	var count int

	criteria := core.Db.Where(dbmodels.ResourceType{})

	//ResourceType.FindFilterCode remove this line for disable generator functionality

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
	//            if core.Db.NewScope(&dbmodels.ResourceType{}).HasColumn(field) {
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

	q := criteria.Model(dbmodels.ResourceType{}).Count(&count)

	if q.Error != nil {
		log.Println("FindResourceType > Ошибка получения данных:", q.Error)
		return result, 0, nil
	}

	// order global criteria
	if len(filter.Order) > 0 {
		for index, Field := range filter.Order {
			if core.Db.NewScope(&dbmodels.ResourceType{}).HasColumn(Field) {
				criteria = criteria.Order("\"" + strings.ToLower(Field) + "\"" + " " + filter.OrderDirection[index])
			} else {
				err = errors.NewErrorWithCode("Ordering by unknown Field", errors.ErrorCodeNotValid, Field)
				return
			}
		}
	}

	q = criteria.Limit(limit).Offset(offset).Find(&dbmodelData)

	if q.Error != nil {
		log.Println("FindResourceType > Ошибка получения данных2:", q.Error)
		return []types.ResourceType{}, 0, nil
	}

	// подготовка id для получения связанных сущностей
	for _, item := range dbmodelData {
		foundIds = append(foundIds, item.ID)
	}

	// получение связнаных сущностей

	//формирование результатов
	for _, item := range dbmodelData {
		result = append(result, AssignResourceTypeTypeFromDb(item))
	}

	return result, count, nil
}

func ResourceTypeMultiCreate(filter types.ResourceTypeFilter) (data []types.ResourceType, err error) {

	typeModelList, err := filter.GetResourceTypeModelList()

	if err != nil {
		return
	}

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetResourceTypeModel(typeModel)
		item, e := ResourceTypeCreate(filter, tx)

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

func ResourceTypeCreate(filter types.ResourceTypeFilter, query *gorm.DB) (data types.ResourceType, err error) {

	typeModel := filter.GetResourceTypeModel()
	dbModel := AssignResourceTypeDbFromType(typeModel)
	dbModel.ID = 0

	dbModel.Validate()

	if !dbModel.IsValid() {
		fmt.Println("ResourceTypeCreate > Create ResourceType error:", dbModel)
		return types.ResourceType{}, dbModel.GetValidationError()
	}

	query = query.Create(&dbModel)

	if query.Error != nil {
		fmt.Println("ResourceTypeCreate > Create ResourceType error:", query.Error)
		return types.ResourceType{}, errors.NewErrorWithCode("cant create ResourceType", errors.ErrorCodeSqlError, "")
	}

	return AssignResourceTypeTypeFromDb(dbModel), nil
}

func ResourceTypeRead(filter types.ResourceTypeFilter) (data types.ResourceType, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1
	filter.ClearIds()
	filter.AddId(filter.GetCurrentId())

	findData, _, err := ResourceTypeFind(filter)

	if len(findData) > 0 {
		return findData[0], nil
	}

	return types.ResourceType{}, errors.NewErrorWithCode("Not found", errors.ErrorCodeNotFound, "")
}

func ResourceTypeMultiUpdate(filter types.ResourceTypeFilter) (data []types.ResourceType, err error) {

	typeModelList, err := filter.GetResourceTypeModelList()

	if err != nil {
		return
	}

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetResourceTypeModel(typeModel)
		filter.ClearIds()
		filter.SetCurrentId(typeModel.Id)

		item, e := ResourceTypeUpdate(filter, tx)

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

func ResourceTypeUpdate(filter types.ResourceTypeFilter, query *gorm.DB) (data types.ResourceType, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	existsModel, err := ResourceTypeRead(filter)

	if existsModel.Id < 1 || err != nil {
		err = errors.NewErrorWithCode("ResourceType not found in db with id: "+strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "Id")
		return
	}

	newModel := filter.GetResourceTypeModel()

	updateModel := AssignResourceTypeDbFromType(newModel)
	updateModel.ID = existsModel.Id

	//updateModel.Some = newModel.Some

	updateModel.Name = newModel.Name
	//updateModel.Field remove this line for disable generator functionality

	updateModel.Validate()

	if !updateModel.IsValid() {
		err = updateModel.GetValidationError()
		return
	}

	q := query.Model(dbmodels.ResourceType{}).Save(&updateModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	data = AssignResourceTypeTypeFromDb(updateModel)
	return
}

func ResourceTypeMultiDelete(filter types.ResourceTypeFilter) (isOk bool, err error) {

	typeModelList, err := filter.GetResourceTypeModelList()

	if err != nil {
		return
	}

	isOk = true

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetResourceTypeModel(typeModel)
		filter.ClearIds()
		filter.SetCurrentId(typeModel.Id)

		_, e := ResourceTypeDelete(filter, tx)

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

func ResourceTypeDelete(filter types.ResourceTypeFilter, query *gorm.DB) (isOk bool, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	existsModel, err := ResourceTypeRead(filter)

	if existsModel.Id < 1 || err != nil {

		if err != nil {
			err = errors.NewErrorWithCode("ResourceType not found in db with id: "+strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "")
		}
		return
	}

	dbModel := AssignResourceTypeDbFromType(existsModel)
	q := query.Model(dbmodels.ResourceType{}).Where(dbmodels.ResourceType{ID: dbModel.ID}).Delete(&dbModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	isOk = true
	return
}

func ResourceTypeFindOrCreate(filter types.ResourceTypeFilter) (data types.ResourceType, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	findOrCreateModel := AssignResourceTypeDbFromType(filter.GetResourceTypeModel())
	//findOrCreateModel.Field remove this line for disable generator functionality

	findOrCreateModel.Validate()

	if !findOrCreateModel.IsValid() {
		err = findOrCreateModel.GetValidationError()
		return
	}

	q := core.Db.Model(dbmodels.ResourceType{}).Where(dbmodels.ResourceType{ID: findOrCreateModel.ID}).FirstOrCreate(&findOrCreateModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	data = AssignResourceTypeTypeFromDb(findOrCreateModel)
	return
}

func ResourceTypeUpdateOrCreate(filter types.ResourceTypeFilter) (data types.ResourceType, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	updateOrCreateModel := AssignResourceTypeDbFromType(filter.GetResourceTypeModel())
	//updateOrCreateModel.Field remove this line for disable generator functionality

	updateOrCreateModel.Validate()

	if !updateOrCreateModel.IsValid() {
		err = updateOrCreateModel.GetValidationError()
		return
	}

	//please uncomment and set criteria
	//q := core.Db.Model(dbmodels.ResourceType{}).Where(dbmodels.ResourceType{ID: updateOrCreateModel.ID}).Assign(dbmodels.ResourceType{/*PLEASE SET CRITERIA*/}).FirstOrCreate(&updateOrCreateModel)

	//if q.Error != nil {
	//    err = q.Error
	//    return
	//}

	data = AssignResourceTypeTypeFromDb(updateOrCreateModel)
	return
}

// add all assign functions

func AssignResourceTypeTypeFromDb(dbResourceType dbmodels.ResourceType) types.ResourceType {

	//AssignResourceTypeTypeFromDb predefine remove this line for disable generator functionality

	return types.ResourceType{
		Id:   dbResourceType.ID,
		Name: dbResourceType.Name,
		//AssignResourceTypeTypeFromDb.Field remove this line for disable generator functionality
	}
}

func AssignResourceTypeDbFromType(typeModel types.ResourceType) dbmodels.ResourceType {

	//AssignResourceTypeDbFromType predefine remove this line for disable generator functionality

	return dbmodels.ResourceType{
		ID:   typeModel.Id,
		Name: typeModel.Name,
		//AssignResourceTypeDbFromType.Field remove this line for disable generator functionality
	}
}
