package dbrepo

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

func updateBuilder(id int, title, description, status string) (string, []interface{}, error) {
	builder := sq.Update(persons).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	if title != "" {
		builder = builder.Set("title", title)
	}

	if description != "" {
		builder = builder.Set("description", description)
	}

	if status != "" {
		builder = builder.Set("status", status)
	}

	return builder.ToSql()
}

func (r *DBRepo) Update(ctx context.Context, id int) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := updateBuilder(id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}
