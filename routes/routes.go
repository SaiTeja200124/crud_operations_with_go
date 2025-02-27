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
	r.Get("/books", handler.GetAll)
	r.Get("/books/{query}", handler.Get)
	r.Put("/books/{bookID}", handler.Update)
	r.Delete("/books/{bookID}", handler.Delete)

	return r
}
