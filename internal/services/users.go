package services

import (
	"context"

	"github.com/sakamoto-max/diablo/internal/domain"
	"github.com/sakamoto-max/diablo/internal/dto"
	"github.com/sakamoto-max/diablo/internal/pkg/password"
	"github.com/sakamoto-max/diablo/internal/pkg/token"
	"github.com/sakamoto-max/diablo/internal/repository"
)

type User struct {
	db *repository.Db
}

func (u *User) Register(ctx context.Context, user dto.RegisterUser) (*domain.User, error) {

	hashedPass, err := password.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPass

	resp, err := u.db.User.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := token.GenerateToken(resp.Id)
	if err != nil {
		return nil, err
	}

	resp.Token = token

	return resp, nil
}
