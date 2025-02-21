package db

import (
	"connection_to_pg/config"
	"connection_to_pg/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database interface for dependency injection
type Database interface {
	Create(value interface{}) *gorm.DB
}

var DB *gorm.DB

func OpenDatabase() error {
	var err error
	dbConfig := config.GetDatabaseConfig()

	// Construct the PostgreSQL connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
		return err
	}

	// Automatically migrate the Book model
	err = DB.AutoMigrate(&models.Book{})
	if err != nil {
		log.Fatalf("failed to migrate database schema: %v", err)
		return err
	}

	return nil
}

func CloseDatabase() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func GetDB() Database {
	return DB
}
