package db

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

func deleteQuery(table string, id int) (string, []interface{}, error) {
	builder := sq.Delete(table).
		PlaceholderFormat(sq.Dollar)

	if table == personTable {
		builder = builder.Where(sq.Eq{"id": id})
	}

	if table == nationalityTable {
		builder = builder.Where(sq.Eq{"person_id": id})
	}

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

	query, args, err := deleteQuery(nationalityTable, id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	query, args, err = deleteQuery(personTable, id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}
