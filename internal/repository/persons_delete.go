package repository

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

func deleteQuery(id int) (string, []interface{}, error) {
	builder := sq.Delete(persons).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *DBRepo) Delete(ctx context.Context, id int) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := deleteQuery(id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}
