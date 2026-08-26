package password

import (
	"errors"
	"fmt"

	myErrs "github.com/sakamoto-max/diablo/internal/pkg/myerrors"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passInBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", myErrs.Wrap(fmt.Errorf("failed to hash the password : %w", err), myErrs.Internal)
	}

	return string(passInBytes), nil
}

func ComparePassword(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return nil
	}

	return myErrs.Wrap(errors.New("passwords do not match"), myErrs.BadRequest)
}
