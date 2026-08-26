package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sakamoto-max/diablo/internal/config"
	myErrs "github.com/sakamoto-max/diablo/internal/pkg/myerrors"
)

type Claims struct {
	UserId string
	jwt.RegisteredClaims
}

var SECRET_KEY string

func GenerateToken(userId string) (string, error) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "diablo",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7 * 3000)),
		},
	})

	token, err := tkn.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", myErrs.Wrap(fmt.Errorf("failed to generate token : %w", err), myErrs.Internal)
	}

	return token, nil
}

func ValidateToken(token string) (*Claims, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return nil, myErrs.Wrap(fmt.Errorf("failed to parse token : %w", err), myErrs.Unauthorized)
	}

	if !tkn.Valid {
		return nil, myErrs.Wrap(errors.New("token is invalid"), myErrs.Unauthorized)
	}

	return claims, nil
}


func Init(config *config.Config) {
	SECRET_KEY = config.Jwt.SecretKey
}