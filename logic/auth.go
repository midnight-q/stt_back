package logic

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"stt_back/common"
	"stt_back/core"
	"stt_back/dbmodels"
	"stt_back/settings"
	"stt_back/types"
)

func AuthFind(filter types.AuthFilter) (result []types.Auth, totalRecords int, err error) {
	return
}

func AuthCreate(filter types.AuthFilter, query *gorm.DB) (data types.Auth, err error) {

	typeModel := filter.GetAuthModel()
	dbAuth := AssignAuthDbFromType(typeModel)
	dbAuth.Validate()

	if dbAuth.IsValid() {

		dbUser := dbmodels.User{}

		query = query.Where(dbmodels.User{Email: dbAuth.Email}).Find(&dbUser)

		if dbUser.ID < 1 {
			return types.Auth{}, errors.New("cant create Auth")
		}
		hashErr := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(typeModel.Password+settings.PASSWORD_SALT))

		if hashErr == nil {

			token := common.RandomString(32)
			dbAuth.Token = token

			core.Db.Model(dbmodels.Auth{}).Where(dbmodels.Auth{Token: token}).First(&dbAuth)
			dbAuth.IsActive = true
			dbAuth.UserId = dbUser.ID
			dbAuth.Email = dbUser.Email
			passLen := "empty"
			if len(dbUser.Password) > 0 {
				passLen = "used"
			}
			dbAuth.Password = passLen
			q := core.Db.Save(&dbAuth)

			if q.Error != nil {

				log.Println("GetNewToken > Ошибка сохранения данных:", q.Error)

			} else {

				typeModel.Token = dbAuth.Token
				typeModel.Password = "******"

				return typeModel, nil
			}
		}

		fmt.Println("AuthCreate > Create Auth error:", hashErr)
		return types.Auth{}, errors.New("cant create Auth")

	} else {

		fmt.Println("AuthCreate > Invalid data:", dbAuth)
		return types.Auth{}, dbAuth.GetValidationError()
	}
}

func AuthMultiCreate(filter types.AuthFilter) (data types.Auth, err error) {
	return
}

func AuthRead(filter types.AuthFilter) (data types.Auth, err error) {
	return
}

func AuthUpdate(filter types.AuthFilter, query *gorm.DB) (data types.Auth, err error) {
	return
}

func AuthMultiUpdate(filter types.AuthFilter) (data types.Auth, err error) {
	return
}

func AuthMultiDelete(filter types.AuthFilter) (isOk bool, err error) {
	return
}

func AuthDelete(filter types.AuthFilter, query *gorm.DB) (isOk bool, err error) {

	dbAuth := dbmodels.Auth{}

	q := query.Model(dbmodels.Auth{}).Where(dbmodels.Auth{Token: filter.Token}).Find(&dbAuth)

	dbAuth.IsActive = false
	q = core.Db.Model(dbmodels.Auth{}).Save(&dbAuth)

	if q.Error != nil {
		err = q.Error
		return
	}

	isOk = true
	return
}

func AuthFindOrCreate(filter types.AuthFilter) (data types.Auth, err error) {
	return
}

func AuthUpdateOrCreate(filter types.AuthFilter) (data types.Auth, err error) {
	return
}

func AssignAuthTypeFromDb(dbAuth dbmodels.Auth) types.Auth {

	//AssignAuthTypeFromDb predefine remove this line for disable generator functionality

	return types.Auth{
		Email:    dbAuth.Email,
		Password: "*******",
		Token:    dbAuth.Token,
		UserId:   dbAuth.UserId,
		//AssignAuthTypeFromDb.Field remove this line for disable generator functionality
	}
}

func AssignAuthDbFromType(typeModel types.Auth) dbmodels.Auth {

	//AssignAuthDbFromType predefine remove this line for disable generator functionality
	token := common.RandomString(32)

	return dbmodels.Auth{
		Token:    token,
		UserId:   typeModel.UserId,
		Email:    typeModel.Email,
		Password: typeModel.Password,
		//AssignAuthDbFromType.Field remove this line for disable generator functionality
	}
}
