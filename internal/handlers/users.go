package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sakamoto-max/diablo/internal/dto"
	"github.com/sakamoto-max/diablo/internal/middleware"
	"github.com/sakamoto-max/diablo/internal/pkg/myerrors"
	"github.com/sakamoto-max/diablo/internal/services"
)

type User struct {
	service *services.User
	myerrors.HttpErr
}

func (u *User) Register(w http.ResponseWriter, r *http.Request) {

	var input dto.RegisterUser

	json.NewDecoder(r.Body).Decode(&input)

	err := input.Validate()
	if err != nil {
		dto.ValidationErrWriter(w, err)
		return
	}

	resp, err := u.service.Register(r.Context(), input)
	if err != nil {
		u.WriteError(w, err)
		return
	}

	respWriter(w, resp, http.StatusCreated)
}

func (h *Handlers) Alive(w http.ResponseWriter, r *http.Request) {

	userId := middleware.GetUserId(r.Context())
	if userId == "" {
		// todo
	}

	w.WriteHeader(http.StatusOK)
}
