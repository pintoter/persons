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

func TestNote_GetNoteById(t *testing.T) {
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

	id := 1
	userId := 1
	notes := []entity.Note{
		{
			ID:          id,
			UserId:      userId,
			Title:       "Test title",
			Description: "Test description",
			Date:        time.Now().Round(time.Second),
			Status:      entity.StatusDone,
		},
	}

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantNote     entity.Note
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "date", "status"}).
					AddRow(notes[0].ID, notes[0].UserId, notes[0].Title, notes[0].Description, notes[0].Date, notes[0].Status)

				expectedQuery := "SELECT id, user_id, title, description, date, status FROM notes WHERE user_id = $1 AND id = $2"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.userId, args.id).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			args:     args{id: id, userId: userId},
			wantNote: notes[0],
		},
		{
			name: "Failed_NotFound",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "date", "status"})

				expectedQuery := "SELECT id, user_id, title, description, date, status FROM notes WHERE user_id = $1 AND id = $2"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.userId, args.id).WillReturnError(errors.New("failed test")).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			args:     args{id: id, userId: userId},
			wantNote: entity.Note{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			gotNote, err := r.GetNoteById(context.Background(), tt.args.id, tt.args.userId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantNote, gotNote)
		})
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestNote_GetNoteByTitle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		title  string
		userId int
	}

	type mockBehavior func(args args)

	id := 1
	title := "Test title"
	userId := 1
	notes := []entity.Note{
		{
			ID:          id,
			UserId:      userId,
			Title:       title,
			Description: "Test description",
			Date:        time.Now().Round(time.Second),
			Status:      entity.StatusDone,
		},
	}

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantNote     entity.Note
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "date", "status"}).
					AddRow(notes[0].ID, notes[0].UserId, notes[0].Title, notes[0].Description, notes[0].Date, notes[0].Status)

				expectedQuery := "SELECT id, user_id, title, description, date, status FROM notes WHERE user_id = $1 AND title = $2"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.userId, args.title).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			args:     args{title: title, userId: userId},
			wantNote: notes[0],
		},
		{
			name: "Failed_NotFound",
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "date", "status"})

				expectedQuery := "SELECT id, user_id, title, description, date, status FROM notes WHERE user_id = $1 AND title = $2"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.userId, args.title).WillReturnError(errors.New("failed test")).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			args:     args{title: title, userId: userId},
			wantNote: entity.Note{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			gotNote, err := r.GetNoteByTitle(context.Background(), tt.args.title, tt.args.userId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantNote, gotNote)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetNotes(t *testing.T) {
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

	notes := []entity.Note{
		{
			ID:          1,
			UserId:      1,
			Title:       "Test title 1",
			Description: "Test description 1",
			Date:        time.Now().Round(time.Second),
			Status:      entity.StatusDone,
		},
		{
			ID:          2,
			UserId:      1,
			Title:       "Test title 2",
			Description: "Test description 2",
			Date:        time.Now().Round(time.Second),
			Status:      entity.StatusDone,
		},
		{
			ID:          3,
			UserId:      1,
			Title:       "Test title 3",
			Description: "Test description 3",
			Date:        time.Now().Round(time.Second),
			Status:      entity.StatusDone,
		},
		{
			ID:          4,
			UserId:      1,
			Title:       "Test title 4",
			Description: "Test description 4",
			Date:        time.Now().Round(time.Second),
			Status:      entity.StatusDone,
		},
	}

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantNotes    []entity.Note
		wantErr      bool
	}{
		{
			name: "Success",
			args: args{userId: 1},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "date", "status"}).
					AddRow(notes[0].ID, notes[0].UserId, notes[0].Title, notes[0].Description, notes[0].Date, notes[0].Status).
					AddRow(notes[1].ID, notes[1].UserId, notes[1].Title, notes[1].Description, notes[1].Date, notes[1].Status).
					AddRow(notes[2].ID, notes[2].UserId, notes[2].Title, notes[2].Description, notes[2].Date, notes[2].Status).
					AddRow(notes[3].ID, notes[3].UserId, notes[3].Title, notes[3].Description, notes[3].Date, notes[3].Status)

				expectedQuery := "SELECT id, user_id, title, description, date, status FROM notes WHERE user_id = $1 ORDER BY id ASC"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.userId).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			wantNotes: []entity.Note{notes[0], notes[1], notes[2], notes[3]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			got, err := r.GetNotes(context.Background(), tt.args.userId)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantNotes, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetNotesExtended(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := New(db)

	type args struct {
		limit  int
		offset int
		status string
		date   time.Time
		userId int
	}

	type mockBehavior func(args args)

	dateFormatted, _ := time.Parse("2006-01-02", "2020-04-18")

	notes := []entity.Note{
		{
			ID:          1,
			UserId:      1,
			Title:       "Test title 1",
			Description: "Test description 1",
			Date:        time.Time{},
			Status:      entity.StatusNotDone,
		},
		{
			ID:          2,
			UserId:      1,
			Title:       "Test title 2",
			Description: "Test description 2",
			Date:        time.Time{},
			Status:      entity.StatusDone,
		},
		{
			ID:          3,
			UserId:      1,
			Title:       "Test title 3",
			Description: "Test description 3",
			Date:        time.Time{},
			Status:      entity.StatusNotDone,
		},
		{
			ID:          4,
			UserId:      1,
			Title:       "Test title 4",
			Description: "Test description 4",
			Date:        time.Time{},
			Status:      entity.StatusDone,
		},
	}

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantNotes    []entity.Note
		wantErr      bool
	}{
		{
			name: "SuccessWithoutStatusAndDate",
			args: args{
				limit:  5,
				offset: 5,
				userId: 1,
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "date", "status"}).
					AddRow(notes[0].ID, notes[0].UserId, notes[0].Title, notes[0].Description, notes[0].Date, notes[0].Status).
					AddRow(notes[1].ID, notes[1].UserId, notes[1].Title, notes[1].Description, notes[1].Date, notes[1].Status).
					AddRow(notes[2].ID, notes[2].UserId, notes[2].Title, notes[2].Description, notes[2].Date, notes[2].Status).
					AddRow(notes[3].ID, notes[3].UserId, notes[3].Title, notes[3].Description, notes[3].Date, notes[3].Status)

				expectedQuery := "SELECT id, user_id, title, description, date, status FROM notes WHERE user_id = $1 ORDER BY id ASC LIMIT 5 OFFSET 5"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.userId).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			wantNotes: []entity.Note{notes[0], notes[1], notes[2], notes[3]},
			wantErr:   false,
		},
		{
			name: "SuccessWithoutDate",
			args: args{
				limit:  5,
				offset: 5,
				status: entity.StatusDone,
				userId: 1,
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "date", "status"}).
					AddRow(notes[1].ID, notes[1].UserId, notes[1].Title, notes[1].Description, notes[1].Date, notes[1].Status).
					AddRow(notes[3].ID, notes[3].UserId, notes[3].Title, notes[3].Description, notes[3].Date, notes[3].Status)

				expectedQuery := "SELECT id, user_id, title, description, date, status FROM notes WHERE user_id = $1 AND status = $2 ORDER BY id ASC LIMIT 5 OFFSET 5"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.userId, args.status).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			wantNotes: []entity.Note{notes[1], notes[3]},
			wantErr:   false,
		},
		{
			name: "SuccessWithStatusAndDate",
			args: args{
				limit:  5,
				offset: 5,
				date:   dateFormatted,
				status: entity.StatusDone,
				userId: 1,
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "date", "status"}).
					AddRow(notes[1].ID, notes[1].UserId, notes[1].Title, notes[1].Description, dateFormatted, notes[1].Status).
					AddRow(notes[3].ID, notes[1].UserId, notes[3].Title, notes[3].Description, dateFormatted, notes[3].Status)

				expectedQuery := "SELECT id, user_id, title, description, date, status FROM notes WHERE user_id = $1 AND status = $2 AND date = $3 ORDER BY id ASC LIMIT 5 OFFSET 5"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.userId, args.status, args.date).WillReturnRows(rows)

				mock.ExpectCommit()
			},
			wantNotes: []entity.Note{
				{
					ID:          notes[1].ID,
					UserId:      notes[1].UserId,
					Title:       notes[1].Title,
					Description: notes[1].Description,
					Date:        dateFormatted,
					Status:      notes[1].Status,
				},
				{
					ID:          notes[3].ID,
					UserId:      notes[3].UserId,
					Title:       notes[3].Title,
					Description: notes[3].Description,
					Date:        dateFormatted,
					Status:      notes[3].Status,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			got, err := r.GetNotesExtended(context.Background(), tt.args.limit, tt.args.offset, tt.args.status, tt.args.date, tt.args.userId)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantNotes, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
