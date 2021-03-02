package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

// Database is used to inject orm instance in server.
var Database *gorm.DB

func SetupDatabase() error {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	Database = db
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}
