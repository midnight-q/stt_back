package migrations

import (
	"stt_back/core"
	"stt_back/dbmodels"
)

func AddRecordNumberForLogs() {
	clearDeletedRecord()
	if IsDataMigrate() {
		return
	}

	userIds := getUserIds()

	for _, userId := range userIds {
		logs, count := getLogsByUserId(userId)
		if count < 1 {
			continue
		}

		for i, log := range logs {
			log.RecordNumber = i + 1
			core.Db.Save(&log)
		}
	}
}

func clearDeletedRecord() {
	core.Db.Unscoped().Model(dbmodels.ConverterLog{}).Where("deleted_at is not null").Delete(&dbmodels.ConverterLog{})
}

func IsDataMigrate() bool {
	count := 0
	core.Db.Unscoped().Model(dbmodels.ConverterLog{}).Where("record_number < 1 or record_number is null").Count(&count)
	return count == 0
}

func getLogsByUserId(userId int) (res []dbmodels.ConverterLog, count int) {
	core.Db.Model(dbmodels.ConverterLog{}).Where(dbmodels.ConverterLog{UserId: userId}).Order("created_at").Find(&res)
	return res, len(res)
}

func getUserIds() (res []int) {
	data := []dbmodels.ConverterLog{}
	core.Db.Unscoped().Model(dbmodels.ConverterLog{}).Select("user_id").Group("user_id").Find(&data)
	for _, d := range data {
		res = append(res, d.UserId)
	}
	return
}
