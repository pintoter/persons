package db

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pintoter/persons/internal/service"
)

func updateBuilder(id int, data *service.UpdateParams) (string, []interface{}, error) {
	builder := sq.Update(persons).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	if data.Name != nil {
		builder = builder.Set("name", *data.Name)
	}

	if data.Surname != nil {
		builder = builder.Set("surname", *data.Surname)
	}

	if data.Patronymic != nil {
		builder = builder.Set("patronymic", *data.Patronymic)
	}

	if data.Age != nil {
		builder = builder.Set("age", *data.Age)
	}

	if data.Gender != nil {
		builder = builder.Set("gender", *data.Gender)
	}

	if data.Nationalize != nil {
		builder = builder.Set("nationalize", *data.Nationalize)
	}

	return builder.ToSql()
}

func (r *DBRepo) Update(ctx context.Context, id int, params *service.UpdateParams) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := updateBuilder(id, params)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}
