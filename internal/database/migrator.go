package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"github.com/sakamoto-max/diablo/internal/config"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func Migrate(config *config.Config) error {

	log.Println("starting database migrations")
	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
	config.Pg.User,
	config.Pg.Password,
	config.Pg.Host,
	
	config.Pg.Port,
	config.Pg.Database,
	config.Pg.SSLMode,
)

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to connect to database : %v", err)
	}
	defer conn.Close(context.Background())

	migrator, err := migrate.NewMigrator(context.Background(), conn, "schema_version")
	if err != nil {
		return fmt.Errorf("failed to create migrator : %v", err)
	}

	subTree, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("failed to get the sub tree : %w", err)
	}

	err = migrator.LoadMigrations(subTree)
	if err != nil {
		return fmt.Errorf("failed to load migrations : %w", err)
	}

	currentVersion, err := migrator.GetCurrentVersion(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get current version : %w", err)
	}

	err = migrator.Migrate(context.Background())
	if err != nil {
		return fmt.Errorf("failed to migrate : %w", err)
	}

	if currentVersion == int32(len(migrator.Migrations)) {
		log.Printf("database is up to date : version %v", currentVersion)
	} else {
		log.Printf("database has migrated from %v to %v", currentVersion, len(migrator.Migrations))
	}

	return nil
}
