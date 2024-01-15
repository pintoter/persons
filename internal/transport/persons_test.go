package transport

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/pintoter/todo-list/internal/entity"
	"github.com/pintoter/todo-list/internal/service"
	mock_service "github.com/pintoter/todo-list/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateNoteHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockINotesRepository, note entity.Note)

	dateFormatted, _ := time.Parse(dateFormat, "2020-01-20")

	tests := []struct {
		name                 string
		inputBody            string
		inputNote            entity.Note
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			inputBody: `{
					"date": "2020-01-20",
					"description": "Test description",
					"status": "not_done",
					"title": "Test one"
				}`,
			inputNote: entity.Note{
				Title:       "Test one",
				Description: "Test description",
				Status:      "not_done",
				Date:        dateFormatted,
			},
			mockBehavior: func(s *mock_service.MockINotesRepository, note entity.Note) {
				s.EXPECT().GetByTitle(gomock.Any(), note.Title).Return(entity.Note{}, sql.ErrNoRows)
				s.EXPECT().Create(gomock.Any(), note).Return(1, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(successCUDResponse{Message: "note created successfully"}, "", "    ")
				return string(resp)
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			noterepo := mock_service.NewMockINotesRepository(c)
			tt.mockBehavior(noterepo, tt.inputNote)

			service := &service.Service{
				IRepository: noterepo,
			}

			handler := NewHandler(service)

			// Init endpoint
			mux := mux.NewRouter()
			mux.HandleFunc("/note", handler.createNoteHandler).Methods(http.MethodPost)

			// Create request
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/note", bytes.NewBufferString(tt.inputBody))

			mux.ServeHTTP(w, r)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}

func TestGetNoteHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockINotesRepository, id int)

	dateFormatted, _ := time.Parse(dateFormat, "2020-01-20")

	tests := []struct {
		name                 string
		inputId              int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "Ok",
			inputId: 1,
			mockBehavior: func(s *mock_service.MockINotesRepository, id int) {
				s.EXPECT().GetById(gomock.Any(), id).Return(entity.Note{
					Title:       "Test one",
					Description: "Test description",
					Status:      "not_done",
					Date:        dateFormatted,
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(getNoteResponse{Note: entity.Note{
					Title:       "Test one",
					Description: "Test description",
					Status:      "not_done",
					Date:        dateFormatted,
				}}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name:    "FailedNotExist",
			inputId: 1,
			mockBehavior: func(s *mock_service.MockINotesRepository, id int) {
				s.EXPECT().GetById(gomock.Any(), id).Return(entity.Note{}, entity.ErrNoteNotExists)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(errorResponse{
					Err: entity.ErrNoteNotExists.Error(),
				}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name:    "FailedWithErr",
			inputId: 1,
			mockBehavior: func(s *mock_service.MockINotesRepository, id int) {
				s.EXPECT().GetById(gomock.Any(), id).Return(entity.Note{}, errors.New("any error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(errorResponse{
					Err: "any error",
				}, "", "    ")
				return string(resp)
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			noterepo := mock_service.NewMockINotesRepository(c)
			tt.mockBehavior(noterepo, tt.inputId)

			service := &service.Service{
				IRepository: noterepo,
			}

			handler := &Handler{
				service: service,
			}

			mux := mux.NewRouter()
			mux.HandleFunc("/note/{id}", handler.getNoteHandler).Methods(http.MethodGet)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/note/1", nil)

			mux.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestGetNotesHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockINotesRepository)

	dateFormatted, _ := time.Parse(dateFormat, "2020-01-20")

	tests := []struct {
		name                 string
		inputId              int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "Ok",
			inputId: 1,
			mockBehavior: func(s *mock_service.MockINotesRepository) {
				s.EXPECT().GetNotes(gomock.Any()).Return([]entity.Note{
					{
						ID:          1,
						Title:       "Test one",
						Description: "Test description",
						Status:      "not_done",
						Date:        dateFormatted,
					},
					{
						ID:          2,
						Title:       "Test two",
						Description: "Test description",
						Status:      "not_done",
						Date:        dateFormatted,
					},
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(getNotesResponse{Notes: []entity.Note{
					{
						ID:          1,
						Title:       "Test one",
						Description: "Test description",
						Status:      "not_done",
						Date:        dateFormatted,
					},
					{
						ID:          2,
						Title:       "Test two",
						Description: "Test description",
						Status:      "not_done",
						Date:        dateFormatted,
					},
				}}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name:    "FailedWithErr",
			inputId: 1,
			mockBehavior: func(s *mock_service.MockINotesRepository) {
				s.EXPECT().GetNotes(gomock.Any()).Return(nil, errors.New("any error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(errorResponse{
					Err: "any error",
				}, "", "    ")
				return string(resp)
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			noterepo := mock_service.NewMockINotesRepository(c)
			tt.mockBehavior(noterepo)

			service := &service.Service{
				IRepository: noterepo,
			}

			handler := &Handler{
				service: service,
			}

			mux := mux.NewRouter()
			mux.HandleFunc("/notes", handler.getNotesHandler).Methods(http.MethodGet)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/notes", nil)

			mux.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestGetNotesExtendedHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockINotesRepository, limit, offset int, status string, date time.Time)

	dateFormatted, _ := time.Parse(dateFormat, "2020-01-20")

	tests := []struct {
		name                 string
		inputBody            string
		inputLimit           int
		inputStatus          string
		inputDate            time.Time
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			inputLimit:  5,
			inputStatus: "not_done",
			inputDate:   time.Time{},
			inputBody:   `{"status": "not_done"}`,
			mockBehavior: func(s *mock_service.MockINotesRepository, limit, offset int, status string, date time.Time) {
				s.EXPECT().GetNotesExtended(gomock.Any(), limit, offset, status, date).Return([]entity.Note{
					{
						ID:          1,
						Title:       "Test one",
						Description: "Test description",
						Status:      "not_done",
						Date:        dateFormatted,
					},
					{
						ID:          2,
						Title:       "Test two",
						Description: "Test description",
						Status:      "not_done",
						Date:        dateFormatted,
					},
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(getNotesResponse{Notes: []entity.Note{
					{
						ID:          1,
						Title:       "Test one",
						Description: "Test description",
						Status:      "not_done",
						Date:        dateFormatted,
					},
					{
						ID:          2,
						Title:       "Test two",
						Description: "Test description",
						Status:      "not_done",
						Date:        dateFormatted,
					},
				}}, "", "    ")
				return string(resp)
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			noterepo := mock_service.NewMockINotesRepository(c)
			tt.mockBehavior(noterepo, tt.inputLimit, 0, tt.inputStatus, tt.inputDate)

			service := &service.Service{
				IRepository: noterepo,
			}

			handler := &Handler{
				service: service,
			}

			mux := mux.NewRouter()
			mux.HandleFunc("/notes/{page}", handler.getNotesExtendedHandler).Methods(http.MethodPost)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/notes/1", bytes.NewBufferString(tt.inputBody))

			mux.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestDeleteById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockINotesRepository, id int)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(s *mock_service.MockINotesRepository, id int) {
				s.EXPECT().GetById(gomock.Any(), id).Return(entity.Note{}, nil)
				s.EXPECT().DeleteById(gomock.Any(), id).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(successCUDResponse{Message: "note deleted succesfully"}, "", "    ")
				return string(resp)
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			noterepo := mock_service.NewMockINotesRepository(c)
			tt.mockBehavior(noterepo, 1)

			service := &service.Service{
				IRepository: noterepo,
			}

			handler := &Handler{
				service: service,
			}

			mux := mux.NewRouter()
			mux.HandleFunc("/note/{id}", handler.deleteNoteHandler).Methods(http.MethodDelete)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/note/1", nil)

			mux.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestDeleteNotes(t *testing.T) {
	type mockBehavior func(s *mock_service.MockINotesRepository)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(s *mock_service.MockINotesRepository) {
				s.EXPECT().DeleteNotes(gomock.Any()).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(successCUDResponse{Message: "notes deleted succesfully"}, "", "    ")
				return string(resp)
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			noterepo := mock_service.NewMockINotesRepository(c)
			tt.mockBehavior(noterepo)

			service := &service.Service{
				IRepository: noterepo,
			}

			handler := &Handler{
				service: service,
			}

			mux := mux.NewRouter()
			mux.HandleFunc("/notes", handler.deleteNotesHandler).Methods(http.MethodDelete)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/notes", nil)

			mux.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestUpdateNoteHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockINotesRepository, id int, title, description, status string)

	tests := []struct {
		name                 string
		inputId              int
		inputTitle           string
		inputStatus          string
		inputDescription     string
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:             "Ok",
			inputId:          1,
			inputTitle:       "Test one",
			inputStatus:      "not_done",
			inputDescription: "Test description",
			inputBody: `{
					"description": "Test description",
					"status": "not_done",
					"title": "Test one"
				}`,
			mockBehavior: func(s *mock_service.MockINotesRepository, id int, title, description, status string) {
				s.EXPECT().GetById(gomock.Any(), id).Return(entity.Note{}, nil)
				s.EXPECT().GetByTitle(gomock.Any(), title).Return(entity.Note{}, entity.ErrNoteNotExists)
				s.EXPECT().UpdateNote(gomock.Any(), id, title, description, status).Return(nil)
			},
			expectedStatusCode: http.StatusAccepted,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(successCUDResponse{Message: "note updated successfully"}, "", "    ")
				return string(resp)
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			noterepo := mock_service.NewMockINotesRepository(c)
			tt.mockBehavior(noterepo, tt.inputId, tt.inputTitle, tt.inputDescription, tt.inputStatus)

			service := &service.Service{
				IRepository: noterepo,
			}

			handler := NewHandler(service)

			// Init endpoint
			mux := mux.NewRouter()
			mux.HandleFunc("/note/{id}", handler.updateNoteHandler).Methods(http.MethodPatch)

			// Create request
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PATCH", "/note/1", bytes.NewBufferString(tt.inputBody))

			mux.ServeHTTP(w, r)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}
