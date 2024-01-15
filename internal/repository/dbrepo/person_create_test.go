package dbrepo

import (
	"context"
	"errors"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pintoter/todo-list/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestNote_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		note entity.Note
	}

	type mockBehavior func(args args)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		id           int
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO notes (user_id,title,description,date,status) VALUES ($1,$2,$3,$4,$5) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.note.UserId, args.note.Title, args.note.Description, args.note.Date, args.note.Status).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()
			},
			args: args{
				note: entity.Note{
					UserId:      1,
					Title:       "Test title",
					Description: "Test describstion",
					Date:        time.Now().Round(time.Second),
					Status:      entity.StatusDone,
				},
			},
			id: 1,
		},
		{
			name: "Success_WithEmptyDateAndDescription",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO notes (user_id,title,description,date,status) VALUES ($1,$2,$3,$4,$5) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.note.UserId, args.note.Title, args.note.Description, args.note.Date, args.note.Status).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()
			},
			args: args{
				note: entity.Note{
					UserId: 1,
					Title:  "Test title",
					Status: entity.StatusDone,
				},
			},
			id: 1,
		},
		{
			name: "Failed_EmptyTitle",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO notes (user_id,title,description,date,status) VALUES ($1,$2,$3,$4,$5) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.note.UserId, args.note.Title, args.note.Description, args.note.Date, args.note.Status).WillReturnError(errors.New("empty title"))

				mock.ExpectRollback()
			},
			args: args{
				note: entity.Note{
					UserId:      1,
					Title:       "",
					Description: "Test describstion",
					Date:        time.Now().Round(time.Second),
					Status:      entity.StatusDone,
				},
			},
			wantErr: true,
		},
		{
			name: "Failed_EmptyUserId",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO notes (user_id,title,description,date,status) VALUES ($1,$2,$3,$4,$5) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.note.UserId, args.note.Title, args.note.Description, args.note.Date, args.note.Status).WillReturnError(errors.New("empty id"))

				mock.ExpectRollback()
			},
			args: args{
				note: entity.Note{
					Title:       "Title",
					Description: "Test describstion",
					Date:        time.Now().Round(time.Second),
					Status:      "not dont",
				},
			},
			wantErr: true,
		},
		{
			name: "Failed_InvalidStatus",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedExec := "INSERT INTO notes (user_id,title,description,date,status) VALUES ($1,$2,$3,$4,$5) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExec)).
					WithArgs(args.note.UserId, args.note.Title, args.note.Description, args.note.Date, args.note.Status).WillReturnError(errors.New("invalid status"))

				mock.ExpectRollback()
			},
			args: args{
				note: entity.Note{
					Title:       "Title",
					Description: "Test describstion",
					Date:        time.Now().Round(time.Second),
					Status:      "not dont",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			got, err := r.CreateNote(context.Background(), tt.args.note)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.id, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
