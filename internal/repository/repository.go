package repository

import (
	"database/sql"
)

type Db struct {
	Suite SuiteIface
}

func New(db *sql.DB) *Db {
	return &Db{
		Suite: &Suite{db: db},
	}
}


