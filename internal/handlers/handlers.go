package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sakamoto-max/diablo/internal/services"
)

type Handlers struct {
	Ops        *Ops
	FileSystem *FileSystem
}

func NewHandlers(service *services.Service) *Handlers {
	return &Handlers{
		Ops:        &Ops{},
		FileSystem: &FileSystem{service: service.FileSystem},
	}
}

func respWriter(w http.ResponseWriter, resp any, statusCode int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
