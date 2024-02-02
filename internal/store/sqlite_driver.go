package store

import (
	"log"

	"github.com/lthnh15032001/ngrok-impl/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSqliteDB() (*gorm.DB, error) {
	//  testing purpose
	db, err := gorm.Open(sqlite.Open("./iot.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	log.Printf("Successfuly Connect Sqlite")
	if err := db.AutoMigrate(
		&models.TunnelAgentModel{},
		&models.UserModel{},
	); err != nil {
		return db, err
	}
	return db, nil
}
