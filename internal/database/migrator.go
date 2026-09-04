package database

import (
	"database/sql"
	"embed"
	"fmt"
)

//go:embed migrations/up.sql
var upMigrationFiles embed.FS

//go:embed migrations/down.sql
var downMigrationFiles embed.FS

func CreateTables(db *sql.DB) error {
	dataInBytes, err := upMigrationFiles.ReadFile("migrations/up.sql")
	if err != nil {
		return fmt.Errorf("failed to read file : %w", err)
	}

	_, err = db.Exec(string(dataInBytes))
	if err != nil {
		return fmt.Errorf("failed to execute sql : %w", err)
	}

	return nil
}

func DropTables(db *sql.DB) error {
	dataInBytes, err := downMigrationFiles.ReadFile("migrations/down.sql")
	if err != nil {
		return fmt.Errorf("failed to read file : %w", err)
	}

	_, err = db.Exec(string(dataInBytes))
	if err != nil {
		return fmt.Errorf("failed to execute sql : %w", err)
	}
	return nil
}
