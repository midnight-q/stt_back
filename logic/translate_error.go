package logic

    import (
        "stt_back/types"
        "strconv"

        "log"
        "fmt"
        "stt_back/core"
        "stt_back/errors"
        "stt_back/dbmodels"
		"strings"

        "github.com/jinzhu/gorm"
    )

func TranslateErrorFind(filter types.TranslateErrorFilter)  (result []types.TranslateError, totalRecords int, err error) {

    foundIds 	:= []int{}
    dbmodelData	:= []dbmodels.TranslateError{}
    limit       := filter.PerPage
    offset      := filter.GetOffset()

    filterIds 	:= filter.GetIds()
    filterExceptIds 	:= filter.GetExceptIds()

    var count int

    criteria := core.Db.Where(dbmodels.TranslateError{})

	if len(filter.ErrorCodes) > 0 {
        criteria = criteria.Where("code in (?)", filter.ErrorCodes)
    }
	//TranslateError.FindFilterCode remove this line for disable generator functionality

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
    //            if core.Db.NewScope(&dbmodels.TranslateError{}).HasColumn(field) {
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

    q := criteria.Model(dbmodels.TranslateError{}).Count(&count)

    if q.Error != nil {
       log.Println("FindTranslateError > Ошибка получения данных:", q.Error)
       return result, 0, nil
    }

    // order global criteria
    if len(filter.Order) > 0  {
        for index, Field := range filter.Order {
             if core.Db.NewScope(&dbmodels.TranslateError{}).HasColumn(Field) {
                criteria = criteria.Order("\"" + strings.ToLower(Field) + "\"" + " " + filter.OrderDirection[index])
            } else {
				err = errors.NewErrorWithCode("Ordering by unknown Field", errors.ErrorCodeNotValid ,Field)
                return
            }
        }
    }

    q = criteria.Limit(limit).Offset(offset).Find(&dbmodelData)

    if q.Error != nil {
       log.Println("FindTranslateError > Ошибка получения данных2:", q.Error)
       return []types.TranslateError{}, 0, nil
    }

    // подготовка id для получения связанных сущностей
    for _, item := range dbmodelData {
        foundIds = append(foundIds, item.ID)
    }

    // получение связнаных сущностей

    //формирование результатов
    for _, item := range dbmodelData {
       result = append(result, AssignTranslateErrorTypeFromDb(item))
    }

    return result, count, nil
}


func TranslateErrorMultiCreate(filter types.TranslateErrorFilter)  (data []types.TranslateError, err error) {

    typeModelList, err := filter.GetTranslateErrorModelList()

    if err != nil {
        return
    }

    tx := core.Db.Begin()

    for _, typeModel := range typeModelList {

        filter.SetTranslateErrorModel(typeModel)
        item, e := TranslateErrorCreate(filter, tx)

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

func TranslateErrorCreate(filter types.TranslateErrorFilter, query *gorm.DB)  (data types.TranslateError, err error) {

    typeModel := filter.GetTranslateErrorModel()
    dbModel := AssignTranslateErrorDbFromType(typeModel)
    dbModel.ID = 0

    dbModel.Validate()

    if ! dbModel.IsValid() {
        fmt.Println("TranslateErrorCreate > Create TranslateError error:", dbModel)
        return types.TranslateError{}, dbModel.GetValidationError()
    }

    query = query.Create(&dbModel)

    if query.Error != nil {
        fmt.Println("TranslateErrorCreate > Create TranslateError error:", query.Error)
        return types.TranslateError{}, errors.NewErrorWithCode("cant create TranslateError", errors.ErrorCodeSqlError, "")
    }

    return AssignTranslateErrorTypeFromDb(dbModel), nil
}

func TranslateErrorRead(filter types.TranslateErrorFilter)  (data types.TranslateError, err error) {

    filter.Pagination.CurrentPage = 1
    filter.Pagination.PerPage = 1
    filter.ClearIds()
    filter.AddId(filter.GetCurrentId())

    findData, _, err := TranslateErrorFind(filter)

    if len(findData) > 0 {
        return findData[0], nil
    }

    return types.TranslateError{}, errors.NewErrorWithCode("Not found", errors.ErrorCodeNotFound, "")
}


func TranslateErrorMultiUpdate(filter types.TranslateErrorFilter)  (data []types.TranslateError, err error) {

    typeModelList, err := filter.GetTranslateErrorModelList()

    if err != nil {
        return
    }

    tx := core.Db.Begin()

    for _, typeModel := range typeModelList {

        filter.SetTranslateErrorModel(typeModel)
        filter.ClearIds()
        filter.SetCurrentId(typeModel.Id)

        item, e := TranslateErrorUpdate(filter, tx)

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

func TranslateErrorUpdate(filter types.TranslateErrorFilter, query *gorm.DB)  (data types.TranslateError, err error) {

    filter.Pagination.CurrentPage = 1
    filter.Pagination.PerPage = 1

    existsModel, err := TranslateErrorRead(filter)

    if existsModel.Id < 1 || err != nil {
        err = errors.NewErrorWithCode("TranslateError not found in db with id: " + strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "Id")
        return
    }

    newModel := filter.GetTranslateErrorModel()

    updateModel := AssignTranslateErrorDbFromType(newModel)
    updateModel.ID = existsModel.Id

    //updateModel.Some = newModel.Some

    updateModel.Code = newModel.Code
	updateModel.LanguageCode = newModel.LanguageCode
	updateModel.Translate = newModel.Translate
	//updateModel.Field remove this line for disable generator functionality

    updateModel.Validate()

    if !updateModel.IsValid() {
        err = updateModel.GetValidationError()
        return
    }

    q := query.Model(dbmodels.TranslateError{}).Save(&updateModel)

    if q.Error != nil {
        err = q.Error
        return
    }

    data = AssignTranslateErrorTypeFromDb(updateModel)
    return
}


func TranslateErrorMultiDelete(filter types.TranslateErrorFilter)  (isOk bool, err error) {

    typeModelList, err := filter.GetTranslateErrorModelList()

    if err != nil {
        return
    }

    isOk = true

    tx := core.Db.Begin()

    for _, typeModel := range typeModelList {

        filter.SetTranslateErrorModel(typeModel)
        filter.ClearIds()
        filter.SetCurrentId(typeModel.Id)

        _, e := TranslateErrorDelete(filter, tx)

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

func TranslateErrorDelete(filter types.TranslateErrorFilter, query *gorm.DB)  (isOk bool, err error) {

    filter.Pagination.CurrentPage = 1
    filter.Pagination.PerPage = 1

    existsModel, err := TranslateErrorRead(filter)

    if existsModel.Id < 1 || err != nil {

        if err != nil {
            err = errors.NewErrorWithCode("TranslateError not found in db with id: " + strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "")
        }
        return
    }

    dbModel := AssignTranslateErrorDbFromType(existsModel)
    q := query.Model(dbmodels.TranslateError{}).Where(dbmodels.TranslateError{ID: dbModel.ID}).Delete(&dbModel)

    if q.Error != nil {
        err = q.Error
        return
    }

    isOk = true
    return
}

func TranslateErrorFindOrCreate(filter types.TranslateErrorFilter)  (data types.TranslateError, err error) {

    filter.Pagination.CurrentPage = 1
    filter.Pagination.PerPage = 1

    findOrCreateModel := AssignTranslateErrorDbFromType(filter.GetTranslateErrorModel())
	//findOrCreateModel.Field remove this line for disable generator functionality

    findOrCreateModel.Validate()

    if !findOrCreateModel.IsValid() {
        err = findOrCreateModel.GetValidationError()
        return
    }

    q := core.Db.Model(dbmodels.TranslateError{}).Where(dbmodels.TranslateError{ID: findOrCreateModel.ID}).FirstOrCreate(&findOrCreateModel)

    if q.Error != nil {
        err = q.Error
        return
    }

    data = AssignTranslateErrorTypeFromDb(findOrCreateModel)
    return
}


func TranslateErrorUpdateOrCreate(filter types.TranslateErrorFilter)  (data types.TranslateError, err error) {

    filter.Pagination.CurrentPage = 1
    filter.Pagination.PerPage = 1

    updateOrCreateModel := AssignTranslateErrorDbFromType(filter.GetTranslateErrorModel())
	//updateOrCreateModel.Field remove this line for disable generator functionality

    updateOrCreateModel.Validate()

    if !updateOrCreateModel.IsValid() {
        err = updateOrCreateModel.GetValidationError()
        return
    }

    //please uncomment and set criteria
    //q := core.Db.Model(dbmodels.TranslateError{}).Where(dbmodels.TranslateError{ID: updateOrCreateModel.ID}).Assign(dbmodels.TranslateError{/*PLEASE SET CRITERIA*/}).FirstOrCreate(&updateOrCreateModel)

    //if q.Error != nil {
    //    err = q.Error
    //    return
    //}

    data = AssignTranslateErrorTypeFromDb(updateOrCreateModel)
    return
}


// add all assign functions

func AssignTranslateErrorTypeFromDb(dbTranslateError dbmodels.TranslateError) types.TranslateError {

    //AssignTranslateErrorTypeFromDb predefine remove this line for disable generator functionality

    return types.TranslateError{
        Id: dbTranslateError.ID,
        Code: dbTranslateError.Code,
		LanguageCode: dbTranslateError.LanguageCode,
		Translate: dbTranslateError.Translate,
		//AssignTranslateErrorTypeFromDb.Field remove this line for disable generator functionality
    }
}

func AssignTranslateErrorDbFromType(typeModel types.TranslateError) dbmodels.TranslateError {

    //AssignTranslateErrorDbFromType predefine remove this line for disable generator functionality
    
    return dbmodels.TranslateError{
        ID: typeModel.Id,
        Code: typeModel.Code,
		LanguageCode: typeModel.LanguageCode,
		Translate: typeModel.Translate,
		//AssignTranslateErrorDbFromType.Field remove this line for disable generator functionality
    }
}


