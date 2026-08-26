package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sakamoto-max/diablo/internal/config"
)

func NewPgPool(config *config.Config) (*pgxpool.Pool, error) {
	pgUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.Pg.User,
		config.Pg.Password,
		config.Pg.Host,
		config.Pg.Port,
		config.Pg.Database,
		config.Pg.SSLMode,
	)
	pgConfig, err := pgxpool.ParseConfig(pgUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres url : %w", err)
	}

	pgConfig.MaxConns = int32(config.Pg.MaxOpenConns)
	pgConfig.MaxConnLifetime = time.Duration(config.Pg.MaxLifetime)
	pgConfig.MaxConnIdleTime = time.Duration(config.Pg.MaxIdleTime)

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool : %w", err)
	}

	return pool, nil
}
