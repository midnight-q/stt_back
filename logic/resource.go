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

func ResourceFind(filter types.ResourceFilter)  (result []types.Resource, totalRecords int, err error) {

    foundIds 	:= []int{}
    dbmodelData	:= []dbmodels.Resource{}
    limit       := filter.PerPage
    offset      := filter.GetOffset()

    filterIds 	:= filter.GetIds()
    filterExceptIds 	:= filter.GetExceptIds()

    var count int

    criteria := core.Db.Where(dbmodels.Resource{})

	//Resource.FindFilterCode remove this line for disable generator functionality

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
    //            if core.Db.NewScope(&dbmodels.Resource{}).HasColumn(field) {
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

    q := criteria.Model(dbmodels.Resource{}).Count(&count)

    if q.Error != nil {
       log.Println("FindResource > Ошибка получения данных:", q.Error)
       return result, 0, nil
    }

    // order global criteria
    if len(filter.Order) > 0  {
        for index, Field := range filter.Order {
             if core.Db.NewScope(&dbmodels.Resource{}).HasColumn(Field) {
                criteria = criteria.Order("\"" + strings.ToLower(Field) + "\"" + " " + filter.OrderDirection[index])
            } else {
				err = errors.NewErrorWithCode("Ordering by unknown Field", errors.ErrorCodeNotValid ,Field)
                return
            }
        }
    }

    q = criteria.Limit(limit).Offset(offset).Find(&dbmodelData)

    if q.Error != nil {
       log.Println("FindResource > Ошибка получения данных2:", q.Error)
       return []types.Resource{}, 0, nil
    }

    // подготовка id для получения связанных сущностей
    for _, item := range dbmodelData {
        foundIds = append(foundIds, item.ID)
    }

    // получение связнаных сущностей

    //формирование результатов
    for _, item := range dbmodelData {
       result = append(result, AssignResourceTypeFromDb(item))
    }

    return result, count, nil
}


func ResourceMultiCreate(filter types.ResourceFilter)  (data []types.Resource, err error) {

    typeModelList, err := filter.GetResourceModelList()

    if err != nil {
        return
    }

    tx := core.Db.Begin()

    for _, typeModel := range typeModelList {

        filter.SetResourceModel(typeModel)
        item, e := ResourceCreate(filter, tx)

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

func ResourceCreate(filter types.ResourceFilter, query *gorm.DB)  (data types.Resource, err error) {

    typeModel := filter.GetResourceModel()
    dbModel := AssignResourceDbFromType(typeModel)
    dbModel.ID = 0

    dbModel.Validate()

    if ! dbModel.IsValid() {
        fmt.Println("ResourceCreate > Create Resource error:", dbModel)
        return types.Resource{}, dbModel.GetValidationError()
    }

    query = query.Create(&dbModel)

    if query.Error != nil {
        fmt.Println("ResourceCreate > Create Resource error:", query.Error)
        return types.Resource{}, errors.NewErrorWithCode("cant create Resource", errors.ErrorCodeSqlError, "")
    }

    return AssignResourceTypeFromDb(dbModel), nil
}

func ResourceRead(filter types.ResourceFilter)  (data types.Resource, err error) {

    filter.Pagination.CurrentPage = 1
    filter.Pagination.PerPage = 1
    filter.ClearIds()
    filter.AddId(filter.GetCurrentId())

    findData, _, err := ResourceFind(filter)

    if len(findData) > 0 {
        return findData[0], nil
    }

    return types.Resource{}, errors.NewErrorWithCode("Not found", errors.ErrorCodeNotFound, "")
}


func ResourceMultiUpdate(filter types.ResourceFilter)  (data []types.Resource, err error) {

    typeModelList, err := filter.GetResourceModelList()

    if err != nil {
        return
    }

    tx := core.Db.Begin()

    for _, typeModel := range typeModelList {

        filter.SetResourceModel(typeModel)
        filter.ClearIds()
        filter.SetCurrentId(typeModel.Id)

        item, e := ResourceUpdate(filter, tx)

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

func ResourceUpdate(filter types.ResourceFilter, query *gorm.DB)  (data types.Resource, err error) {

    filter.Pagination.CurrentPage = 1
    filter.Pagination.PerPage = 1

    existsModel, err := ResourceRead(filter)

    if existsModel.Id < 1 || err != nil {
        err = errors.NewErrorWithCode("Resource not found in db with id: " + strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "Id")
        return
    }

    newModel := filter.GetResourceModel()

    updateModel := AssignResourceDbFromType(newModel)
    updateModel.ID = existsModel.Id

    //updateModel.Some = newModel.Some

    updateModel.Name = newModel.Name
	updateModel.Code = newModel.Code
	updateModel.TypeId = newModel.TypeId
	//updateModel.Field remove this line for disable generator functionality

    updateModel.Validate()

    if !updateModel.IsValid() {
        err = updateModel.GetValidationError()
        return
    }

    q := query.Model(dbmodels.Resource{}).Save(&updateModel)

    if q.Error != nil {
        err = q.Error
        return
    }

    data = AssignResourceTypeFromDb(updateModel)
    return
}


func ResourceMultiDelete(filter types.ResourceFilter)  (isOk bool, err error) {

    typeModelList, err := filter.GetResourceModelList()

    if err != nil {
        return
    }

    isOk = true

    tx := core.Db.Begin()

    for _, typeModel := range typeModelList {

        filter.SetResourceModel(typeModel)
        filter.ClearIds()
        filter.SetCurrentId(typeModel.Id)

        _, e := ResourceDelete(filter, tx)

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

func ResourceDelete(filter types.ResourceFilter, query *gorm.DB)  (isOk bool, err error) {

    filter.Pagination.CurrentPage = 1
    filter.Pagination.PerPage = 1

    existsModel, err := ResourceRead(filter)

    if existsModel.Id < 1 || err != nil {

        if err != nil {
            err = errors.NewErrorWithCode("Resource not found in db with id: " + strconv.Itoa(filter.GetCurrentId()), errors.ErrorCodeNotFound, "")
        }
        return
    }

    dbModel := AssignResourceDbFromType(existsModel)
    q := query.Model(dbmodels.Resource{}).Where(dbmodels.Resource{ID: dbModel.ID}).Delete(&dbModel)

    if q.Error != nil {
        err = q.Error
        return
    }

    isOk = true
    return
}

func ResourceFindOrCreate(filter types.ResourceFilter)  (data types.Resource, err error) {

    filter.Pagination.CurrentPage = 1
    filter.Pagination.PerPage = 1

    findOrCreateModel := AssignResourceDbFromType(filter.GetResourceModel())
	//findOrCreateModel.Field remove this line for disable generator functionality

    findOrCreateModel.Validate()

    if !findOrCreateModel.IsValid() {
        err = findOrCreateModel.GetValidationError()
        return
    }

    q := core.Db.Model(dbmodels.Resource{}).Where(dbmodels.Resource{ID: findOrCreateModel.ID}).FirstOrCreate(&findOrCreateModel)

    if q.Error != nil {
        err = q.Error
        return
    }

    data = AssignResourceTypeFromDb(findOrCreateModel)
    return
}


func ResourceUpdateOrCreate(filter types.ResourceFilter)  (data types.Resource, err error) {

    filter.Pagination.CurrentPage = 1
    filter.Pagination.PerPage = 1

    updateOrCreateModel := AssignResourceDbFromType(filter.GetResourceModel())
	//updateOrCreateModel.Field remove this line for disable generator functionality

    updateOrCreateModel.Validate()

    if !updateOrCreateModel.IsValid() {
        err = updateOrCreateModel.GetValidationError()
        return
    }

    //please uncomment and set criteria
    //q := core.Db.Model(dbmodels.Resource{}).Where(dbmodels.Resource{ID: updateOrCreateModel.ID}).Assign(dbmodels.Resource{/*PLEASE SET CRITERIA*/}).FirstOrCreate(&updateOrCreateModel)

    //if q.Error != nil {
    //    err = q.Error
    //    return
    //}

    data = AssignResourceTypeFromDb(updateOrCreateModel)
    return
}


// add all assign functions

func AssignResourceTypeFromDb(dbResource dbmodels.Resource) types.Resource {

    //AssignResourceTypeFromDb predefine remove this line for disable generator functionality

    return types.Resource{
        Id: dbResource.ID,
        Name: dbResource.Name,
		Code: dbResource.Code,
		TypeId: dbResource.TypeId,
		//AssignResourceTypeFromDb.Field remove this line for disable generator functionality
    }
}

func AssignResourceDbFromType(typeModel types.Resource) dbmodels.Resource {

    //AssignResourceDbFromType predefine remove this line for disable generator functionality
    
    return dbmodels.Resource{
        ID: typeModel.Id,
        Name: typeModel.Name,
		Code: typeModel.Code,
		TypeId: typeModel.TypeId,
		//AssignResourceDbFromType.Field remove this line for disable generator functionality
    }
}


