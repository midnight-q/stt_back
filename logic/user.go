package logic

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"strings"
	"stt_back/core"
	"stt_back/dbmodels"
	"stt_back/errors"
	"stt_back/settings"
	"stt_back/types"
)

func UserFind(filter types.UserFilter) (result []types.User, totalRecords int, err error) {

	foundIds := []int{}
	dbmodelData := []dbmodels.User{}
	limit := filter.PerPage
	offset := filter.GetOffset()

	filterIds := filter.GetIds()
	filterExceptIds := filter.GetExceptIds()

	var count int

	criteria := core.Db.Where(dbmodels.User{})

	//User.FindFilterCode remove this line for disable generator functionality

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
	//            if core.Db.NewScope(&dbmodels.User{}).HasColumn(field) {
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

	q := criteria.Model(dbmodels.User{}).Count(&count)

	if q.Error != nil {
		log.Println("FindUser > Ошибка получения данных:", q.Error)
		return result, 0, nil
	}

	// order global criteria
	if len(filter.Order) > 0 {
		for index, Field := range filter.Order {
			if core.Db.NewScope(&dbmodels.User{}).HasColumn(Field) {
				criteria = criteria.Order("\"" + strings.ToLower(Field) + "\"" + " " + filter.OrderDirection[index])
			} else {
				err = errors.NewErrorWithCode("Ordering by unknown Field", errors.ErrorCodeNotValid, Field)
				return
			}
		}
	}

	q = criteria.Limit(limit).Offset(offset).Find(&dbmodelData)

	if q.Error != nil {
		log.Println("FindUser > Ошибка получения данных2:", q.Error)
		return []types.User{}, 0, nil
	}

	// подготовка id для получения связанных сущностей
	for _, item := range dbmodelData {
		foundIds = append(foundIds, item.ID)
	}

	// получение связнаных сущностей

	//формирование результатов
	for _, item := range dbmodelData {
		result = append(result, AssignUserTypeFromDb(item))
	}

	return result, count, nil
}

func UserMultiCreate(filter types.UserFilter) (data []types.User, err error) {

	typeModelList, err := filter.GetUserModelList()

	if err != nil {
		return
	}

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetUserModel(typeModel)
		item, e := UserCreate(filter, tx)

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

func UserCreate(filter types.UserFilter, query *gorm.DB) (data types.User, err error) {

	typeModel := filter.GetUserModel()
	dbModel := AssignUserDbFromType(typeModel)
	dbModel.ID = 0

	dbModel.Validate()

	if !dbModel.IsValid() {
		fmt.Println("UserCreate > Create User error:", dbModel)
		return types.User{}, dbModel.GetValidationError()
	}

	query = query.Create(&dbModel)

	if query.Error != nil {
		fmt.Println("UserCreate > Create User error:", query.Error)
		return types.User{}, errors.NewErrorWithCode("cant create User", errors.ErrorCodeSqlError, "")
	}

	return AssignUserTypeFromDb(dbModel), nil
}

func UserRead(filter types.UserFilter) (data types.User, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1
	filter.ClearIds()
	filter.AddId(filter.GetCurrentId())

	findData, _, err := UserFind(filter)

	if len(findData) > 0 {
		return findData[0], nil
	}

	return types.User{}, errors.NewErrorWithCode("Not found", errors.ErrorCodeNotFound, "")
}

func UserMultiUpdate(filter types.UserFilter) (data []types.User, err error) {

	typeModelList, err := filter.GetUserModelList()

	if err != nil {
		return
	}

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetUserModel(typeModel)
		filter.ClearIds()
		filter.SetCurrentId(typeModel.Id)

		item, e := UserUpdate(filter, tx)

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

func UserUpdate(filter types.UserFilter, query *gorm.DB) (data types.User, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	existsModel, err := UserRead(filter)

	if existsModel.Id < 1 || err != nil {
		err = errors.NewErrorWithCode("User not found in db with id: "+strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "Id")
		return
	}

	newModel := filter.GetUserModel()

	updateModel := AssignUserDbFromType(newModel)
	updateModel.ID = existsModel.Id

	//updateModel.Some = newModel.Some

	//updateModel.Field remove this line for disable generator functionality

	updateModel.Validate()

	if !updateModel.IsValid() {
		err = updateModel.GetValidationError()
		return
	}

	q := query.Model(dbmodels.User{}).Save(&updateModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	data = AssignUserTypeFromDb(updateModel)
	return
}

func UserMultiDelete(filter types.UserFilter) (isOk bool, err error) {

	typeModelList, err := filter.GetUserModelList()

	if err != nil {
		return
	}

	isOk = true

	tx := core.Db.Begin()

	for _, typeModel := range typeModelList {

		filter.SetUserModel(typeModel)
		filter.ClearIds()
		filter.SetCurrentId(typeModel.Id)

		_, e := UserDelete(filter, tx)

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

func UserDelete(filter types.UserFilter, query *gorm.DB) (isOk bool, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	existsModel, err := UserRead(filter)

	if existsModel.Id < 1 || err != nil {

		if err != nil {
			err = errors.NewErrorWithCode("User not found in db with id: "+strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "")
		}
		return
	}

	dbModel := AssignUserDbFromType(existsModel)
	q := query.Model(dbmodels.User{}).Where(dbmodels.User{ID: dbModel.ID}).Delete(&dbModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	isOk = true
	return
}

func UserFindOrCreate(filter types.UserFilter) (data types.User, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	findOrCreateModel := AssignUserDbFromType(filter.GetUserModel())
	//findOrCreateModel.Field remove this line for disable generator functionality

	findOrCreateModel.Validate()

	if !findOrCreateModel.IsValid() {
		err = findOrCreateModel.GetValidationError()
		return
	}

	q := core.Db.Model(dbmodels.User{}).Where(dbmodels.User{ID: findOrCreateModel.ID}).FirstOrCreate(&findOrCreateModel)

	if q.Error != nil {
		err = q.Error
		return
	}

	data = AssignUserTypeFromDb(findOrCreateModel)
	return
}

func UserUpdateOrCreate(filter types.UserFilter) (data types.User, err error) {

	filter.Pagination.CurrentPage = 1
	filter.Pagination.PerPage = 1

	updateOrCreateModel := AssignUserDbFromType(filter.GetUserModel())
	//updateOrCreateModel.Field remove this line for disable generator functionality

	updateOrCreateModel.Validate()

	if !updateOrCreateModel.IsValid() {
		err = updateOrCreateModel.GetValidationError()
		return
	}

	//please uncomment and set criteria
	//q := core.Db.Model(dbmodels.User{}).Where(dbmodels.User{ID: updateOrCreateModel.ID}).Assign(dbmodels.User{/*PLEASE SET CRITERIA*/}).FirstOrCreate(&updateOrCreateModel)

	//if q.Error != nil {
	//    err = q.Error
	//    return
	//}

	data = AssignUserTypeFromDb(updateOrCreateModel)
	return
}

// add all assign functions

func AssignUserTypeFromDb(dbUser dbmodels.User) types.User {

	//AssignUserTypeFromDb predefine remove this line for disable generator functionality

	return types.User{
		Id:          dbUser.ID,
		Email:       dbUser.Email,
		FirstName:   dbUser.FirstName,
		IsActive:    dbUser.IsActive,
		LastName:    dbUser.LastName,
		MobilePhone: dbUser.MobilePhone,
		Password:    "******",
		//AssignUserTypeFromDb.Field remove this line for disable generator functionality
	}
}

func AssignUserDbFromType(typeModel types.User) dbmodels.User {

	password := []byte(typeModel.Password + settings.PASSWORD_SALT)
	hashedPassword, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	//AssignUserTypeFromDb predefine remove this line for disable generator functionality

	return dbmodels.User{
		ID:          typeModel.Id,
		Email:       typeModel.Email,
		FirstName:   typeModel.FirstName,
		IsActive:    typeModel.IsActive,
		LastName:    typeModel.LastName,
		MobilePhone: typeModel.MobilePhone,
		Password:    string(hashedPassword),
		//AssignUserDbFromType.Field remove this line for disable generator functionality
	}
}
