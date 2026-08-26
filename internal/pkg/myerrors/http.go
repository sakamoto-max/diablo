package myerrors

import (
	"encoding/json"
	"log"
	"net/http"
)

type HttpErr struct{}

func (h *HttpErr) WriteError(w http.ResponseWriter, err error) {
	unWrappedErr := unWrap(err)

	if unWrappedErr.HttpCode() == http.StatusInternalServerError {
		log.Println("internal server error : %w", err)
		unWrappedErr.message = "internal server error"
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(unWrappedErr.HttpCode())
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}
