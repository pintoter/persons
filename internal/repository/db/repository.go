package db

import (
	"database/sql"
)

const (
	persons = "persons"
)

type DBRepo struct {
	db *sql.DB
}

func New(db *sql.DB) *DBRepo {
	return &DBRepo{
		db: db,
	}
}

// Create address of const value for test Get and Update methods
func GetAddress[T any](x T) *T {
	return &x
}
