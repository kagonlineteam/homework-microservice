package models

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error

	if os.Getenv("GIN_MODE") == "release" {
		DB, err = gorm.Open("postgres", os.Getenv("HOMEWORK_POSTGRES_DSN"))
	} else {
		DB, err = gorm.Open("sqlite3", "debug_database.sqlite3")
	}

	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	DB.LogMode(os.Getenv("GIN_MODE") != "release").AutoMigrate(new(Homework))
}
