package dto

import (
	"github.com/go-playground/validator/v10"
)

type RegisterUser struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (r *RegisterUser) Validate() error {
	newValidator := validator.New(validator.WithRequiredStructEnabled())
	return newValidator.Struct(r)
}

type UserIp struct {
	IP string `json:"ip"`
}
