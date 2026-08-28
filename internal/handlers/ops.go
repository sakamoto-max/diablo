package handlers

import (
	"net/http"
)

type Ops struct {
	// service *services.User
	// myerrors.HttpErr
}

func (h *Ops) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
