package store

import (
	"fmt"
	"log"

	"github.com/lthnh15032001/ngrok-impl/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMySQLDB(host string, user string, password string, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local", user, password, host, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	log.Printf("Successfuly Connect %s:@tcp(%s:3306)/%s\n", user, host, dbName)
	if err := db.AutoMigrate(
		&models.Policy{},
		&models.ScheduleModel{},
	); err != nil {
		return db, err
	}
	return db, nil
}
