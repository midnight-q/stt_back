package dbmodels

import (
	"time"
)

type ConverterLog struct {
	ID                int `gorm:"primary_key"`
	FilePath          string
	ResultText        string
	ResultFilePath    string
	ResultFormat      string
	RawResult         string
	ResultHtml        string
	ResultFileDocPath string
	ResultFilePdfPath string
	UserId            int
	SourceFilePath string
	//ConverterLog remove this line for disable generator functionality

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index" json:"-"`

	validator
}

func (converterLog *ConverterLog) Validate() {
	//Validate remove this line for disable generator functionality
}
