package main

import (
	"connection_to_pg/db"
	"connection_to_pg/handlers"
	"connection_to_pg/routes"
	"log"
	"net/http"
)

func main() {
	// Open database connection
	err := db.OpenDatabase()
	if err != nil {
		log.Fatalf("error opening database connection: %v", err)
	}
	defer db.CloseDatabase()

	// Get the actual database instance
	database := db.GetDB()

	// Create a handler with the database dependency
	handler := &handlers.Handler{DB: database}

	// Setup router with the handler instance
	r := routes.SetupRoutes(handler) // Load routes from separate file

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
