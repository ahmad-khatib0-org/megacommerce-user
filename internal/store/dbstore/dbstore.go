// Package dbstore store/retrieve data in the database
package dbstore

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStore struct {
	db *pgxpool.Pool
}

func NewDBStore(pool *pgxpool.Pool) *DBStore {
	return &DBStore{db: pool}
}
