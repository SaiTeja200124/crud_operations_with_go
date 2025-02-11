package models

type Book struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author"`
}

type CreateBookBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author"`
}

type DatabaseConfig struct {
	User     string
	Password string
	DBName   string
	SSLMode  string
	Host     string
	Port     string
}
