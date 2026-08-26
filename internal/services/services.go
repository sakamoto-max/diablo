package services

import (
	"github.com/sakamoto-max/diablo/internal/repository"
)

type Service struct {
	User       *User
	FileSystem *Synchronizer
}

func NewService(db *repository.Db) *Service {
	return &Service{
		User:       &User{db: db},
		FileSystem: &Synchronizer{db: db},
	}
}
