package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sakamoto-max/diablo/internal/domain"
	"github.com/sakamoto-max/diablo/internal/dto"
	myErrs "github.com/sakamoto-max/diablo/internal/pkg/myerrors"
)

type User struct {
	pool *pgxpool.Pool
}

type UserIface interface {
	Register(ctx context.Context, user dto.RegisterUser) (*domain.User, error)
}

func (u *User) Register(ctx context.Context, user dto.RegisterUser) (*domain.User, error) {

	query := `
		INSERT INTO USERS (
			NAME,
			EMAIL,
			PASSWORD
		) VALUES (
			@name,
			@email,
			@password
		)
		RETURNING ID, CREATED_AT
	`

	var id string
	var createdAt time.Time

	err := u.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	}).Scan(&id, &createdAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
				return nil, myErrs.Wrap(errors.New("email already exists"), myErrs.AlreadyExists)
			}

			return nil, myErrs.Wrap(fmt.Errorf("failed to register user : %w", err), myErrs.Internal)
		}
	}

	return &domain.User{
		Id:        id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: createdAt,
	}, nil

}
