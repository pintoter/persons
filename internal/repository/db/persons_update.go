package db

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pintoter/persons/internal/service"
	"github.com/pintoter/persons/pkg/logger"
)

func updateBuilder(id int, data *service.UpdateParams) (string, []interface{}, error) {
	builder := sq.Update(personTable).
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
	return builder.ToSql()
}

func (r *DBRepo) Update(ctx context.Context, id int, params *service.UpdateParams) error {
	logMethod := "repository.Update"
	logger.DebugKV(ctx, "update start", "layer", logMethod, "params", params)
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := updateBuilder(id, params)
	logger.DebugKV(ctx, "update builder", "layer", logMethod, "query", query, "err", err)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}
