package dbrepo

import (
	"context"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pintoter/todo-list/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestUpdateNote(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		id          int
		userId      int
		title       string
		description string
		status      string
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
				id:          1,
				userId:      1,
				title:       "Test title NEW",
				description: "Test description NEW",
				status:      entity.StatusDone,
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "UPDATE notes SET title = $1, description = $2, status = $3 WHERE id = $4 AND user_id = $5"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.title, args.description, args.status, args.id, args.userId).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Failed",
			args: args{
				id:          100,
				userId:      1,
				title:       "Test title NEW",
				description: "Test description NEW",
				status:      entity.StatusDone,
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				expectedQuery := "UPDATE notes SET title = $1, description = $2, status = $3 WHERE id = $4 AND user_id = $5"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.title, args.description, args.status, args.id, args.userId).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			err := r.UpdateNote(context.Background(), tt.args.id, tt.args.userId, tt.args.title, tt.args.description, tt.args.status)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
