package dbrepo

import "database/sql"

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
