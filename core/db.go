package core

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"stt_back/settings"
)

const DbConnectString = "host='" + settings.DbHost +
	"' port='" + settings.DbPort +
	"' user='" + settings.DbUser +
	"' password='" + settings.DbPass +
	"' dbname='" + settings.DbName +
	"' sslmode='disable'"

var Db, DbErr = gorm.Open("postgres", DbConnectString)

func EnableSqlLog() {
	Db.LogMode(true)
}

func DisableSqlLog() {
	Db.LogMode(false)
}

func GetTableName(dbmodel interface{}) string {
	return Db.NewScope(dbmodel).TableName()
}
