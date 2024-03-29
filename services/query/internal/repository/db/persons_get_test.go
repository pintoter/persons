package db

import (
	"context"
	"errors"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pintoter/persons/services/query/internal/entity"
	"github.com/pintoter/persons/services/query/internal/service"
	"github.com/stretchr/testify/assert"
)

func Test_GetPerson(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		id int
	}

	type mockBehavior func(args args)

	id := 1
	persons := []entity.Person{
		{
			ID:         1,
			Name:       "name",
			Surname:    "surname",
			Patronymic: "patronymic",
			Age:        18,
			Gender:     "male",
			Nationalize: []entity.Nationality{
				{
					Country:     "RU",
					Probability: 0.1337,
				},
			},
		},
	}

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantPerson   entity.Person
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "name", "surname", "patronymic", "age", "gender", "nationalize", "probability"}).
					AddRow(
						persons[0].ID,
						persons[0].Name,
						persons[0].Surname,
						persons[0].Patronymic,
						persons[0].Age,
						persons[0].Gender,
						persons[0].Nationalize[0].Country,
						persons[0].Nationalize[0].Probability,
					)

				expectedQuery := `SELECT person.id, person.name, person.surname, person.patronymic, person.age,
				person.gender, n.nationalize, n.probability FROM person JOIN person_nationality n ON n.person_id = person.id WHERE person.id = $1`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.id).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			args:       args{id: id},
			wantPerson: persons[0],
		},
		{
			name: "Failed_NotFound",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "name", "surname", "patronymic", "age", "gender", "nationalize", "probability"})

				expectedQuery := `SELECT person.id, person.name, person.surname, person.patronymic, person.age,
				person.gender, n.nationalize, n.probability FROM person JOIN person_nationality n ON n.person_id = person.id WHERE person.id = $1`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.id).WillReturnError(errors.New("some error")).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			args:       args{id: id},
			wantPerson: entity.Person{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			gotPerson, err := r.GetPerson(context.Background(), tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPerson, gotPerson)
			}
		})
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func Test_GetPersons(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		filters *service.GetFilters
	}

	type mockBehavior func(args args)

	persons := []entity.Person{
		{
			ID:         1,
			Name:       "name",
			Surname:    "surname",
			Patronymic: "patronymic",
			Age:        18,
			Gender:     "male",
			Nationalize: []entity.Nationality{
				{
					Country:     "RU",
					Probability: 0.1337,
				},
			},
		},
		{
			ID:         2,
			Name:       "name",
			Surname:    "surname",
			Patronymic: "patronymic",
			Age:        19,
			Gender:     "male",
			Nationalize: []entity.Nationality{
				{
					Country:     "RU",
					Probability: 0.1337,
				},
			},
		},
		{
			ID:         3,
			Name:       "name",
			Surname:    "surname",
			Patronymic: "patronymic",
			Age:        20,
			Gender:     "male",
			Nationalize: []entity.Nationality{
				{
					Country:     "RU",
					Probability: 0.1337,
				},
			},
		},
	}

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantPersons  []entity.Person
		wantErr      bool
	}{
		{
			name: "Success",
			args: args{
				filters: &service.GetFilters{
					Nationalize: GetAddress[string]("RU"),
					Limit:       5,
					Offset:      0,
				},
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "name", "surname", "patronymic", "age", "gender", "nationalize", "probability"}).
					AddRow(
						persons[0].ID,
						persons[0].Name,
						persons[0].Surname,
						persons[0].Patronymic,
						persons[0].Age,
						persons[0].Gender,
						persons[0].Nationalize[0].Country,
						persons[0].Nationalize[0].Probability,
					).AddRow(
					persons[1].ID,
					persons[1].Name,
					persons[1].Surname,
					persons[1].Patronymic,
					persons[1].Age,
					persons[1].Gender,
					persons[1].Nationalize[0].Country,
					persons[1].Nationalize[0].Probability,
				).AddRow(
					persons[2].ID,
					persons[2].Name,
					persons[2].Surname,
					persons[2].Patronymic,
					persons[2].Age,
					persons[2].Gender,
					persons[2].Nationalize[0].Country,
					persons[2].Nationalize[0].Probability,
				)

				expectedQuery := `SELECT person.id, person.name, person.surname, person.patronymic, person.age, person.gender, n.nationalize, n.probability 
				FROM person JOIN person_nationality n ON n.person_id = person.id WHERE n.nationalize = $1 LIMIT 5 OFFSET 0`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.filters.Nationalize).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			wantPersons: []entity.Person{persons[0], persons[1], persons[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			got, err := r.GetPersons(context.Background(), tt.args.filters)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPersons, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
