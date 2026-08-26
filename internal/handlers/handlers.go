package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sakamoto-max/diablo/internal/services"
)

type Handlers struct {
	User       *User
	FileSystem *FileSystem
}

func NewHandlers(service *services.Service) *Handlers {
	return &Handlers{
		User: &User{service: service.User},
		FileSystem: &FileSystem{service: service.FileSystem},
	}
}

func respWriter(w http.ResponseWriter, resp any, statusCode int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
