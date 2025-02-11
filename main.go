// main.go
package main

import (
	"connection_to_pg/db"
	"connection_to_pg/routes"
	"log"
	"net/http"
)

func main() {
	err := db.OpenDatabase()
	if err != nil {
		log.Printf("error opening database connection %v", err)
	}
	defer db.CloseDatabase()
	r := routes.SetupRoutes() // Load routes from separate file

	http.ListenAndServe("localhost:8080", r)
}
