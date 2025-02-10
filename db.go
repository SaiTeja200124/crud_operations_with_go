package main

import (
	// "database/sql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

var DB *gorm.DB

func OpenDatabase() error {
	var err error
	// Connecting to PostgreSQL using GORM
	DB, err = gorm.Open("postgres", "user=postgres password=password dbname=gopractice sslmode=disable")
	if err != nil {
		return err
	}

	// Automatically migrate the Book model to keep the schema up-to-date
	DB.AutoMigrate(&Book{})
	return nil
}

func CloseDatabase() error {
	return DB.Close()
}
