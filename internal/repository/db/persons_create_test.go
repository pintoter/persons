package db

import (
	"context"
	"errors"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pintoter/persons/internal/entity"
	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		person entity.Person
	}

	type mockBehavior func(args args)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantId       int
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedExecInPerson := "INSERT INTO person (name,surname,patronymic,age,gender) VALUES ($1,$2,$3,$4,$5) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExecInPerson)).
					WithArgs(
						args.person.Name,
						args.person.Surname,
						args.person.Patronymic,
						args.person.Age,
						args.person.Gender,
					).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				expectedExecInNationality := "INSERT INTO person_nationality (person_id,nationalize,probability) VALUES ($1,$2,$3),($4,$5,$6)"

				mock.ExpectExec(regexp.QuoteMeta(expectedExecInNationality)).
					WithArgs(1, args.person.Nationalize[0].Country, args.person.Nationalize[0].Probability, 1, args.person.Nationalize[1].Country, args.person.Nationalize[1].Probability).
					WillReturnResult(sqlmock.NewResult(0, 2))

				mock.ExpectCommit()
			},
			args: args{
				person: entity.Person{
					Name:       "name",
					Surname:    "surname",
					Patronymic: "patronymic",
					Age:        18,
					Gender:     "male",
					Nationalize: []entity.Nationality{
						{
							Country:     "RU",
							Probability: 0.02,
						},
						{
							Country:     "GE",
							Probability: 0.2,
						},
					},
				},
			},
			wantId: 1,
		},
		{
			name: "Success_WithEmptyPatronymic",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO person (name,surname,patronymic,age,gender) VALUES ($1,$2,$3,$4,$5) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(
						args.person.Name,
						args.person.Surname,
						args.person.Patronymic,
						args.person.Age,
						args.person.Gender,
					).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				expectedExecInNationality := "INSERT INTO person_nationality (person_id,nationalize,probability) VALUES ($1,$2,$3)"
				mock.ExpectExec(regexp.QuoteMeta(expectedExecInNationality)).
					WithArgs(1, args.person.Nationalize[0].Country, args.person.Nationalize[0].Probability, 1, args.person.Nationalize[1].Country, args.person.Nationalize[1].Probability).
					WillReturnResult(sqlmock.NewResult(0, 2))
				mock.ExpectCommit()
			},
			args: args{
				person: entity.Person{
					Name:    "name",
					Surname: "surname",
					Age:     18,
					Gender:  "male",
					Nationalize: []entity.Nationality{
						{
							Country:     "RU",
							Probability: 0.02,
						},
						{
							Country:     "GE",
							Probability: 0.2,
						},
					},
				},
			},
			wantId: 1,
		},
		{
			name: "Failed_EmptyName",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO person (name,surname,patronymic,age,gender) VALUES ($1,$2,$3,$4,$5) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(
						args.person.Name,
						args.person.Surname,
						args.person.Patronymic,
						args.person.Age,
						args.person.Gender,
					).WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
			args: args{
				person: entity.Person{
					Name:    "name",
					Surname: "surname",
					Age:     18,
					Gender:  "male",
					Nationalize: []entity.Nationality{
						{
							Country:     "RU",
							Probability: 0.02,
						},
						{
							Country:     "GE",
							Probability: 0.2,
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			gotId, err := r.Create(context.Background(), tt.args.person)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantId, gotId)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
