package repository

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pintoter/persons/internal/entity"
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
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := createBuilder(person)
	if err != nil {
		return 0, err
	}

	var id int
	err = tx.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, tx.Commit()
}
