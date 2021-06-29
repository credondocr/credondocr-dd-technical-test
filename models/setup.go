package models

import (
	"credondocr-dd-technical-test/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupModels() *gorm.DB {
	config := config.GetConfig()

	h := config.GetString("database.host")
	u := config.GetString("database.user")
	pw := config.GetString("database.password")
	dbn := config.GetString("database.name")
	p := config.GetString("database.port")

	dsn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable", h, p, u, dbn, pw)

	fmt.Println("conname is\t\t", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	db.AutoMigrate(&Song{})

	// // Initialize value
	// m := Book{Author: "author1", Title: "title1"}

	// db.Create(&m)

	return db
}
