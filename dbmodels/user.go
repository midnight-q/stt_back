package dbmodels

import (
	"stt_back/common"
	"stt_back/errors"
	"time"
)

type User struct {
	ID          int    `gorm:"primary_key"`
	Email       string `gorm:"type:varchar(100);unique_index"`
	FirstName   string
	IsActive    bool
	LastName    string
	MobilePhone string
	Password    string
	//User remove this line for disable generator functionality

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index" json:"-"`

	validator
}

func (user *User) Validate() {

	if len(user.FirstName) < 1 {
		user.AddValidationError("User first name is empty", errors.ErrorCodeFieldLengthTooShort, "FirstName")
	}

	if len(user.LastName) < 1 {
		user.AddValidationError("User last name is empty", errors.ErrorCodeFieldLengthTooShort, "LastName")
	}

	if len(user.Email) < 3 || !common.ValidateEmail(user.Email) {
		user.AddValidationError("User email not valid", errors.ErrorCodeNotValid, "Email")
	}

	if len(user.MobilePhone) > 3 && !common.ValidateMobile(user.MobilePhone) {
		user.AddValidationError("User mobile phone should be valid or empty. Format +0123456789... ", errors.ErrorCodeNotValid, "MobilePhone")
	}

	//Validate remove this line for disable generator functionality

}
