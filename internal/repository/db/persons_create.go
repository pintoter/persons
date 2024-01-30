package db

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pintoter/persons/internal/entity"
	"github.com/pintoter/persons/pkg/logger"
)

func createPersonBuilder(person entity.Person) (string, []interface{}, error) {
	builder := sq.Insert(personTable).
		Columns("name", "surname", "patronymic", "age", "gender").
		Values(person.Name, person.Surname, person.Patronymic, person.Age, person.Gender).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func createNationalizeBuilder(person entity.Person) (string, []interface{}, error) {
	builder := sq.Insert(nationalityTable).
		Columns("person_id", "nationalize", "probability").
		PlaceholderFormat(sq.Dollar)

	for _, nationalize := range person.Nationalize {
		builder = builder.Values(person.ID, nationalize.Country, nationalize.Probability)
	}

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

	query, args, err := createPersonBuilder(person)
	logger.DebugKV(ctx, "create builder", "layer", logMethod, "query", query, "args", args, "err", err)
	if err != nil {
		return 0, err
	}

	err = tx.QueryRowContext(ctx, query, args...).Scan(&person.ID)
	logger.DebugKV(ctx, "insert in person table", "layer", logMethod, "err", err)
	if err != nil {
		return 0, err
	}

	query, args, err = createNationalizeBuilder(person)
	logger.DebugKV(ctx, "create nationalize builder", "layer", logMethod, "query", query, "args", args, "err", err)
	if err != nil {
		return 0, err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	logger.DebugKV(ctx, "insert in nationalize table", "layer", logMethod, "err", err)
	if err != nil {
		return 0, err
	}
	logger.DebugKV(ctx, "end of creating person", "layer", logMethod)

	return person.ID, tx.Commit()
}
