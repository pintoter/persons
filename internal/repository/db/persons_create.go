package db

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pintoter/persons/internal/entity"
	"github.com/pintoter/persons/pkg/logger"
)

func createBuilder(person entity.Person) (string, []interface{}, error) {
	builder := sq.Insert(persons).
		Columns("name", "surname", "patronymic", "age", "gender", "nationalize").
		Values(person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationalize).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *DBRepo) Create(ctx context.Context, person entity.Person) (int, error) {
	logMethod := "repository.Create"
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	logger.DebugKV(ctx, "begin tx", "layer", logMethod, "err", err)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := createBuilder(person)
	logger.DebugKV(ctx, "create builder", "layer", logMethod, "query", query, "args", args, "err", err)
	if err != nil {
		return 0, err
	}

	var id int
	err = tx.QueryRowContext(ctx, query, args...).Scan(&id)
	logger.DebugKV(ctx, "insert in db", "layer", logMethod, "err", err)
	if err != nil {
		return 0, err
	}

	return id, tx.Commit()
}
