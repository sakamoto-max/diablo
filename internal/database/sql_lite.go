package database

import (
	"database/sql"
	"fmt"

	_ "github.com/ncruces/go-sqlite3/driver"
)

func New() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:demo.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite conn : %w", err)
	}

	return db, nil
}
