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
	Find(dest interface{}, conds ...interface{}) *gorm.DB // ✅ Added Find method
	Where(query interface{}, args ...interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}) *gorm.DB
}

// DatabaseImpl is a wrapper around *gorm.DB to implement Database interface
type DatabaseImpl struct {
	DB *gorm.DB
}

func (d *DatabaseImpl) Create(value interface{}) *gorm.DB {
	return d.DB.Create(value)
}

func (d *DatabaseImpl) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	return d.DB.Find(dest, conds...)
}

func (d *DatabaseImpl) Where(query interface{}, args ...interface{}) *gorm.DB {
	return d.DB.Where(query, args...)
}
func (d *DatabaseImpl) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return d.DB.First(dest, conds...)
}
func (d *DatabaseImpl) Save(value interface{}) *gorm.DB {
	return d.DB.Save(value)
}
func (d *DatabaseImpl) Model(value interface{}) *gorm.DB {
	return d.DB.Model(value)
}
func (d *DatabaseImpl) Delete(value interface{}) *gorm.DB {
	return d.DB.Model(value)
}

var gormDB *gorm.DB

func OpenDatabase() error {
	var err error
	dbConfig := config.GetDatabaseConfig()

	// Construct the PostgreSQL connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)

	gormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
		return err
	}

	// Automatically migrate the Book model
	err = gormDB.AutoMigrate(&models.Book{})
	if err != nil {
		log.Fatalf("failed to migrate database schema: %v", err)
		return err
	}

	return nil
}

func CloseDatabase() error {
	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB returns the wrapped database instance
func GetDB() Database {
	return &DatabaseImpl{DB: gormDB} // ✅ Returns the wrapped struct instead of raw *gorm.DB
}
