package dbmodels

import (
    "time"
    
)

type Region struct {

    ID        int       `gorm:"primary_key"`
    Name int
	Code string
	//Region remove this line for disable generator functionality

    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time `sql:"index" json:"-"`

    validator
}

func (region *Region) Validate() {
    //Validate remove this line for disable generator functionality
}

