package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/sakamoto-max/diablo/internal/handlers"
	"github.com/sakamoto-max/diablo/internal/middleware"
)

func NewRouter(handlers *handlers.Handlers, middlewares *middleware.Middlewares) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/register", handlers.User.Register)

	r.With(middlewares.Auth.ValidateToken).Post("/alive", handlers.Alive)
	r.Post("/new", handlers.FileSystem.CreateNewFileSystem)
	r.Post("/sync", handlers.FileSystem.Sync)
	r.Get("/suite", handlers.FileSystem.Suite)
	r.Post("/ping", handlers.FileSystem.Ping)

	return r
}
