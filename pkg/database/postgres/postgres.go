package postgres

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	driverName = "postgres"
)

type Config interface {
	GetDSN() string
	GetMaxOpenCons() int
	GetMaxIdleCons() int
	GetConnMaxIdleTime() time.Duration
	GetConnMaxLifetime() time.Duration
}

func New(cfg Config) (*sql.DB, error) {
	db, err := sql.Open(driverName, cfg.GetDSN())
	if err != nil {
		return nil, errors.Wrap(err, "opening db")
	}

	db.SetMaxOpenConns(cfg.GetMaxOpenCons())
	db.SetMaxIdleConns(cfg.GetMaxIdleCons())
	db.SetConnMaxIdleTime(cfg.GetConnMaxIdleTime())
	db.SetConnMaxLifetime(cfg.GetConnMaxLifetime())

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "ping DB")
	}

	return db, nil
}
