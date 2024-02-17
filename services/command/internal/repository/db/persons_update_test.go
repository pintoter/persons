package db

import (
	"context"
	"errors"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pintoter/persons/services/command/internal/service"
	"github.com/stretchr/testify/assert"
)

func Test_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		id     int
		params *service.UpdateParams
	}

	type mockBehavior func(args args)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "Success",
			args: args{
				id: 1,
				params: &service.UpdateParams{
					Name: GetAddress[string]("Vlad"),
				},
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "UPDATE person SET name = $1 WHERE id = $2"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.params.Name, args.id).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Failed",
			args: args{
				id: 100,
				params: &service.UpdateParams{
					Surname: GetAddress[string]("Ivanov"),
				},
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "UPDATE person SET surname = $1 WHERE id = $2"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.params.Surname, args.id).
					WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			err := r.Update(context.Background(), tt.args.id, tt.args.params)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
