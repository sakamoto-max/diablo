package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Db struct {
	User  UserIface
	Suite SuiteIface
}

func NewDb(pool *pgxpool.Pool) *Db {
	return &Db{
		User:  &User{pool: pool},
		Suite: &Suite{db: pool},
	}
}

// tables :

// suites:
// ID
// NAME
// LAST_UPDATED
// LAST_SYNCED

// files:
// ID
// NAME
// PATH
// LAST_UPDATED
// LAST_SYNCED
// IS_DIR
