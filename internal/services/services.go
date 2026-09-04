package services

import (
	"github.com/sakamoto-max/diablo/internal/repository"
)

type Service struct {
	FileSystem *Synchronizer
}

func NewService(db *repository.Db) *Service {
	return &Service{
		FileSystem: &Synchronizer{db: db},
	}
}
