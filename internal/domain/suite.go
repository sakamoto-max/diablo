package domain

import (
	"time"

	"github.com/sakamoto-max/diablo/internal/dto"
)

type Suite struct {
	Id    string
	Name  string
	Files []File
}

type File struct {
	Id          string
	Name        string
	Path        string
	Data        []byte
	LastUpdated time.Time
	LastSynced  time.Time
	IsDir       bool
}

type LastSyncedData struct {
	SuiteName  string
	SuiteId    string
	LastSynced time.Time
	Events     []dto.Event
}
