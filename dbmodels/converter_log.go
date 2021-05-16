package dbmodels

import (
	"time"
)

type ConverterLog struct {
	ID                int `gorm:"primary_key"`
	FilePath          string
	ResultTextPath    string
	ResultFilePath    string
	ResultFormat      string
	RawResult         string
	ResultHtmlPath    string
	ResultFileDocPath string
	ResultFilePdfPath string
	UserId            int
	SourceFilePath    string
	RecordNumber int
	//ConverterLog remove this line for disable generator functionality

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index" json:"-"`

	validator
}

func (converterLog *ConverterLog) Validate() {
	//Validate remove this line for disable generator functionality
}
