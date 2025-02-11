// config/config.go
package config

// DatabaseConfig holds the PostgreSQL credentials
import "connection_to_pg/models"

// GetDatabaseConfig returns the database configuration
func GetDatabaseConfig() models.DatabaseConfig {
	return models.DatabaseConfig{
		User:     "postgres",
		Password: "password",
		DBName:   "gopractice",
		SSLMode:  "disable",
		Host:     "localhost",
		Port:     "5432",
	}
}
