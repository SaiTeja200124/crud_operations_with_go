package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	err := OpenDatabase()
	if err != nil {
		log.Printf("error opening database connection %v", err)
	}
	defer CloseDatabase()

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Post("/books", create)

	r.Get("/books", getAll)

	r.Get("/books/{query}", get)

	r.Put("/books/{bookID}", update)

	r.Delete("/books/{bookID}", deleteBook)

	http.ListenAndServe("localhost:8080", r)
}
