package db

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pintoter/persons/internal/entity"
	"github.com/pintoter/persons/internal/service"
	"github.com/pintoter/persons/pkg/logger"
)

func getPersonBuilder(id int) (string, []interface{}, error) {
	builder := sq.Select("person.id", "person.name", "person.surname", "person.patronymic", "person.age", "person.gender", "n.nationalize, n.probability").
		From(personTable).
		Join("person_nationality n ON n.person_id = person.id").
		Where(sq.Eq{"person.id": id}).
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *DBRepo) GetPerson(ctx context.Context, id int) (entity.Person, error) {
	logMethod := "repository.GetPerson"
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	logger.DebugKV(ctx, "begin tx", "layer", logMethod, "err", err)
	if err != nil {
		return entity.Person{}, err
	}
	defer func() { _ = tx.Rollback() }()

	query, args, err := getPersonBuilder(id)
	logger.DebugKV(ctx, "get builder", "layer", logMethod, "query", query, "args", args, "err", err)
	if err != nil {
		return entity.Person{}, err
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return entity.Person{}, err
	}
	defer rows.Close()

	isIdScanned := make(map[int]struct{})
	var person entity.Person
	for rows.Next() {
		var id, age int
		var name, patronymic, surname, gender string
		var nationalize string
		var probability float64

		err = rows.Scan(&id, &name, &surname, &patronymic, &age, &gender, &nationalize, &probability)
		if err != nil {
			logger.DebugKV(ctx, "rows.Scan", "layer", logMethod, "err", err)
			return entity.Person{}, err
		}

		if _, ok := isIdScanned[id]; !ok {
			isIdScanned[id] = struct{}{}
			person.ID = id
			person.Name = name
			person.Surname = surname
			person.Patronymic = patronymic
			person.Age = age
			person.Gender = gender
		}
		person.Nationalize = append(person.Nationalize, entity.Nationality{Country: nationalize, Probability: probability})
	}
	logger.DebugKV(ctx, "get builder", "layer", logMethod, "person", person)

	return person, tx.Commit()
}

func getPersonsBuilder(data *service.GetFilters) (string, []interface{}, error) {
	builder := sq.Select("person.id", "person.name", "person.surname", "person.patronymic", "person.age", "person.gender", "n.nationalize, n.probability").
		From(personTable).
		Join("person_nationality n ON n.person_id = person.id").
		PlaceholderFormat(sq.Dollar)

	if data.Name != nil {
		builder = builder.Where(sq.Eq{"person.name": *data.Name})
	}
	if data.Surname != nil {
		builder = builder.Where(sq.Eq{"person.surname": *data.Surname})
	}
	if data.Patronymic != nil {
		builder = builder.Where(sq.Eq{"person.patronymic": *data.Patronymic})
	}
	if data.Age != nil {
		builder = builder.Where(sq.Eq{"person.age": *data.Age})
	}
	if data.Gender != nil {
		builder = builder.Where(sq.Eq{"person.gender": *data.Gender})
	}
	if data.Nationalize != nil {
		builder = builder.Where(sq.Eq{"n.nationalize": *data.Nationalize})
	}
	builder = builder.Limit(uint64(data.Limit)).Offset(uint64(data.Offset))
	return builder.ToSql()
}

/* FIX ME */
func (r *DBRepo) GetPersons(ctx context.Context, data *service.GetFilters) ([]entity.Person, error) {
	logMethod := "repository.GetPersons"
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

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	isIdScanned := make(map[int]struct{})
	var persons []entity.Person
	var count int
	for rows.Next() {
		var person entity.Person
		var nationalize string
		var probability float64

		err = rows.Scan(&person.ID, &person.Name, &person.Surname, &person.Patronymic, &person.Age, &person.Gender, &nationalize, &probability)
		if err != nil {
			logger.DebugKV(ctx, "rows.Scan", "layer", logMethod, "err", err)
			return nil, err
		}

		if _, ok := isIdScanned[person.ID]; !ok {
			isIdScanned[person.ID] = struct{}{}
			count++
			persons = append(persons, person)
		}
		persons[count-1].Nationalize = append(persons[count-1].Nationalize, entity.Nationality{Country: nationalize, Probability: probability})
		logger.DebugKV(ctx, "rows.Scan", "layer", logMethod, "persons", persons)
	}
	logger.DebugKV(ctx, "result of repo.GetPersons", "layer", logMethod, "persons", persons)

	return persons, tx.Commit()
}
