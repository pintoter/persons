package dbrepo

import (
	"context"
	"errors"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		id     int
		userId int
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
				id:     1,
				userId: 1,
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "DELETE FROM notes WHERE id = $1 AND user_id = $2"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.id, args.userId).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Failed",
			args: args{
				id:     100,
				userId: 1,
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "DELETE FROM notes WHERE id = $1 AND user_id = $2"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.id, args.userId).
					WillReturnError(errors.New("new error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			err := r.DeleteNoteById(context.Background(), tt.args.id, tt.args.userId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteNotes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		userId int
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
				userId: 1,
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "DELETE FROM notes WHERE user_id = $1"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.userId).
					WillReturnResult(sqlmock.NewResult(0, 5))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Failed",
			args: args{
				userId: 1,
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "DELETE FROM notes WHERE user_id = $1"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.userId).
					WillReturnError(errors.New("new error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			err := r.DeleteNotes(context.Background(), tt.args.userId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
