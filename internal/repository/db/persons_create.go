package db

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pintoter/persons/internal/entity"
	"github.com/pintoter/persons/pkg/logger"
)

func createBuilder(person entity.Person) (string, []interface{}, error) {
	builder := sq.Insert(personTable).
		Columns("name", "surname", "patronymic", "age", "gender").
		Values(person.Name, person.Surname, person.Patronymic, person.Age, person.Gender).
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

	builder := sq.Insert(nationalityTable).Columns("person_id", "nationalize", "probability").PlaceholderFormat(sq.Dollar)
	for _, nationalize := range person.Nationalize {
		builder = builder.Values(id, nationalize.Country, nationalize.Probability)
	}
	query1, args1, err := builder.ToSql()
	logger.DebugKV(ctx, "create nationalize builder", "layer", logMethod, "query", query1, "args", args1, "err", err)
	if err != nil {
		logger.DebugKV(ctx, "insert in db", "layer", logMethod, "err", err)
		return 0, err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	logger.DebugKV(ctx, "insert in db", "layer", logMethod, "err", err)
	if err != nil {
		return 0, err
	}

	logger.DebugKV(ctx, "end of creating person", "layer", logMethod)
	return id, tx.Commit()
}
