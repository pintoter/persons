package db

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pintoter/persons/internal/entity"
	"github.com/pintoter/persons/internal/service"
)

func getPersonBuilder(id int) (string, []interface{}, error) {
	builder := sq.Select("id", "name", "surname", "patronymic", "age", "gender", "nationalize").
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

func getPersonsBuilder(data *service.GetFilters) (string, []interface{}, error) {
	builder := sq.Select("id", "name", "surname", "patronymic", "age", "gender", "nationalize").
		From(persons).
		OrderBy("id ASC").
		PlaceholderFormat(sq.Dollar)

	if data.Name != nil {
		builder = builder.Where(sq.Eq{"name": *data.Name})
	}

	if data.Surname != nil {
		builder = builder.Where(sq.Eq{"surname": *data.Surname})
	}

	if data.Patronymic != nil {
		builder = builder.Where(sq.Eq{"patronymic": *data.Patronymic})
	}

	if data.Age != nil {
		builder = builder.Where(sq.Eq{"age": *data.Age})
	}

	if data.Gender != nil {
		builder = builder.Where(sq.Eq{"gender": *data.Gender})
	}

	if data.Nationalize != nil {
		builder = builder.Where(sq.Eq{"nationalize": *data.Nationalize})
	}

	builder = builder.Limit(uint64(data.Limit)).Offset(uint64(data.Offset))

	return builder.ToSql()
}

func (r *DBRepo) GetPersons(ctx context.Context, data *service.GetFilters) ([]entity.Person, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := getPersonsBuilder(data)
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
