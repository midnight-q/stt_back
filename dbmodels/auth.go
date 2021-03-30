package dbmodels

import (
    "stt_back/common"
    "time"
    "stt_back/errors"
)

type Auth struct {

    ID        int       `gorm:"primary_key"`
    Email     string
    Password  string
    Token     string
    IsActive bool
    UserId   int
    //Auth remove this line for disable generator functionality

    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time `sql:"index" json:"-"`

    validator
}

func (auth *Auth) Validate() {
        if len(auth.Email) < 3 || ! common.ValidateEmail(auth.Email)  {
        auth.AddValidationError("User email not valid", errors.ErrorCodeFieldLengthTooShort, "Email")
    }

    if len(auth.Password) < 1 {
        auth.AddValidationError("User password is empty", errors.ErrorCodeNotEmpty, "Password")
    }

    //Validate remove this line for disable generator functionality
}

