package dbrepo

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pintoter/persons/internal/entity"
)

func getPersonBuilder(id int) (string, []interface{}, error) {
	builder := sq.Select("id", "user_id", "title", "description", "date", "status").
		From(persons).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *DBRepo) GetPerson(ctx context.Context, id int) (entity.Person, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return entity.Person{}, err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := getPersonBuilder(id)
	if err != nil {
		return entity.Person{}, err
	}

	var person entity.Person
	err = tx.QueryRowContext(ctx, query, args...).
		Scan(
			&person.ID,
			&person.Name,
			&person.Surname,
			&person.Patronymic,
			&person.Age,
			&person.Gender,
			&person.Nationalize,
		)
	if err != nil {
		return entity.Person{}, err
	}

	return person, tx.Commit()
}

/*-----------------------------
					GET ALL NOTES
 ----------------------------- */

func getNotesBuilder(limit, offset int) (string, []interface{}, error) {
	builder := sq.Select("id", "user_id", "title", "description", "date", "status").
		From(persons).
		OrderBy("id ASC").
		Where(sq.Eq{"user_id": userId}).
		PlaceholderFormat(sq.Dollar)

	if status != "" {
		builder = builder.Where(sq.Eq{"status": status})
	}

	if !date.Equal(time.Time{}) {
		builder = builder.Where(sq.Eq{"date": date})
	}

	if limit != 0 || offset != 0 {
		builder = builder.Limit(uint64(limit)).Offset(uint64(offset))
	}

	return builder.ToSql()
}

func (r *DBRepo) GetPersons(ctx context.Context, limit, offset int) ([]entity.Person, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := getNotesBuilder(0, 0)
	if err != nil {
		return nil, err
	}

	var persons []entity.Person
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var person entity.Person
		if err := rows.Scan(
			&person.ID,
			&person.Name,
			&person.Surname,
			&person.Patronymic,
			&person.Age,
			&person.Gender,
			&person.Nationalize,
		); err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}

	return persons, tx.Commit()
}
