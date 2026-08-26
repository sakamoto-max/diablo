package dto

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ValidationErr struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
}

func ValidationErrWriter(w http.ResponseWriter, err error) {

	allErrs := extrackValidationErrs(err)
	
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(allErrs)
}

func extrackValidationErrs(err error) []ValidationErr {
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		panic(fmt.Errorf("failed to get all the validation errs : %w", err))
	}

	var allErrs []ValidationErr

	for _, err := range validationErrs {
		field := err.Field()
		tag := err.Tag()

		allErrs = append(allErrs, ValidationErr{
			Field: field,
			Tag:   tag,
		})
	}

	return allErrs
}
