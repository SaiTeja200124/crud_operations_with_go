// routes/routes.go
package routes

import (
	"connection_to_pg/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// SetupRoutes initializes the router with all routes
func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Post("/books", handlers.Create)
	r.Get("/books", handlers.GetAll)
	r.Get("/books/{query}", handlers.Get)
	r.Put("/books/{bookID}", handlers.Update)
	r.Delete("/books/{bookID}", handlers.DeleteBook)

	return r
}
