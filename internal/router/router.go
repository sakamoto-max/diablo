package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/sakamoto-max/diablo/internal/handlers"
)

func NewRouter(handlers *handlers.Handlers) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/ping", handlers.Ops.Ping)
	r.Post("/new", handlers.FileSystem.CreateNewFileSystem)
	r.Post("/sync", handlers.FileSystem.Sync)
	r.Get("/suite", handlers.FileSystem.Suite)
	r.Post("/alive", handlers.FileSystem.Alive)

	return r
}
