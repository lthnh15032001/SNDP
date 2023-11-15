package store

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSqliteDB() (*gorm.DB, error) {
	//  testing purpose
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	// if err := db.AutoMigrate(
	// 	&models.Policy{},
	// 	&models.ScheduleModel{},
	// ); err != nil {
	// 	return db, err
	// }
	return db, nil
}
