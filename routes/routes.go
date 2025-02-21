// routes/routes.go
package routes

import (
	"connection_to_pg/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// SetupRoutes initializes the router with all routes
func SetupRoutes(handler *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Post("/books", handler.Create)
	r.Get("/books", handlers.GetAll)
	r.Get("/books/{query}", handlers.Get)
	r.Put("/books/{bookID}", handlers.Update)
	r.Delete("/books/{bookID}", handlers.DeleteBook)

	return r
}
