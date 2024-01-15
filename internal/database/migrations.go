package migrations

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	sourceURL = "file://migrations"
)

type Config interface {
	GetDSN() string
}

func Do(cfg Config) error {
	m, err := migrate.New(sourceURL, cfg.GetDSN())
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	defer func() {
		m.Close()
	}()

	return nil
}
