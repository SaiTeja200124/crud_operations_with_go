// db/db.go
package db

import (
	"connection_to_pg/config"
	"connection_to_pg/models"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

var DB *gorm.DB

func OpenDatabase() error {
	var err error
	dbConfig := config.GetDatabaseConfig()

	// Construct the PostgreSQL connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)

	// Connect to PostgreSQL using GORM
	DB, err = gorm.Open("postgres", dsn)
	if err != nil {
		return err
	}

	// Automatically migrate the Book model to keep the schema up-to-date
	DB.AutoMigrate(&models.Book{})
	return nil
}

func CloseDatabase() error {
	return DB.Close()
}
