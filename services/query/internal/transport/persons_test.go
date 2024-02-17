package transport

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pintoter/persons/services/query/internal/entity"
	"github.com/pintoter/persons/services/query/internal/service"
	mock_service "github.com/pintoter/persons/services/query/internal/service/mocks"

	"github.com/stretchr/testify/assert"
)

func Test_GetPersonHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockRepository, id int)

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
			mockBehavior: func(s *mock_service.MockRepository, id int) {
				s.EXPECT().GetPerson(gomock.Any(), id).Return(entity.Person{
					ID:         1,
					Name:       "name",
					Surname:    "surname",
					Patronymic: "patronymic",
					Age:        18,
					Gender:     "male",
					Nationalize: []entity.Nationality{
						{
							Country:     "RU",
							Probability: 0.1,
						},
						{
							Country:     "KZ",
							Probability: 0.05,
						},
					},
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(getPersonResponse{Person: entity.Person{
					ID:         1,
					Name:       "name",
					Surname:    "surname",
					Patronymic: "patronymic",
					Age:        18,
					Gender:     "male",
					Nationalize: []entity.Nationality{
						{
							Country:     "RU",
							Probability: 0.1,
						},
						{
							Country:     "KZ",
							Probability: 0.05,
						},
					},
				}}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name:    "FailedNotExist",
			inputId: 2,
			mockBehavior: func(s *mock_service.MockRepository, id int) {
				s.EXPECT().GetPerson(gomock.Any(), id).Return(entity.Person{}, sql.ErrNoRows)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(errorResponse{
					Err: entity.ErrPersonNotExists.Error(),
				}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name:    "FailedWithId",
			inputId: 0,
			mockBehavior: func(s *mock_service.MockRepository, id int) {
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(errorResponse{
					Err: entity.ErrInvalidInput.Error(),
				}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name:    "FailedWithErr",
			inputId: 2,
			mockBehavior: func(s *mock_service.MockRepository, id int) {
				s.EXPECT().GetPerson(gomock.Any(), id).Return(entity.Person{}, sql.ErrNoRows)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(errorResponse{
					Err: entity.ErrPersonNotExists.Error(),
				}, "", "    ")
				return string(resp)
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockRepository(c)
			tt.mockBehavior(repo, tt.inputId)

			service := service.New(repo)

			handler := NewHandler(service)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/persons/"+fmt.Sprintf("%d", tt.inputId), nil)

			handler.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func Test_GetPersonsHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockRepository)

	tests := []struct {
		name                 string
		path                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			path: "?nationalize=RU",
			mockBehavior: func(s *mock_service.MockRepository) {
				s.EXPECT().GetPersons(gomock.Any(), gomock.Any()).Return(
					[]entity.Person{
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
									Probability: 0.1,
								},
								{
									Country:     "KZ",
									Probability: 0.05,
								},
							},
						}, {
							ID:         2,
							Name:       "name1",
							Surname:    "surname1",
							Patronymic: "patronymic1",
							Age:        19,
							Gender:     "male",
							Nationalize: []entity.Nationality{
								{
									Country:     "RU",
									Probability: 0.1,
								},
								{
									Country:     "KZ",
									Probability: 0.05,
								},
							},
						},
					},
					nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(getPersonsResponse{Persons: []entity.Person{
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
								Probability: 0.1,
							},
							{
								Country:     "KZ",
								Probability: 0.05,
							},
						},
					}, {
						ID:         2,
						Name:       "name1",
						Surname:    "surname1",
						Patronymic: "patronymic1",
						Age:        19,
						Gender:     "male",
						Nationalize: []entity.Nationality{
							{
								Country:     "RU",
								Probability: 0.1,
							},
							{
								Country:     "KZ",
								Probability: 0.05,
							},
						},
					},
				}}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name: "Ok2",
			path: "?name=name",
			mockBehavior: func(s *mock_service.MockRepository) {
				s.EXPECT().GetPersons(gomock.Any(), gomock.Any()).Return(
					[]entity.Person{
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
									Probability: 0.1,
								},
								{
									Country:     "KZ",
									Probability: 0.05,
								},
							},
						}, {
							ID:         2,
							Name:       "name",
							Surname:    "surname1",
							Patronymic: "patronymic1",
							Age:        19,
							Gender:     "male",
							Nationalize: []entity.Nationality{
								{
									Country:     "RU",
									Probability: 0.1,
								},
								{
									Country:     "KZ",
									Probability: 0.05,
								},
							},
						},
					},
					nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(getPersonsResponse{Persons: []entity.Person{
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
								Probability: 0.1,
							},
							{
								Country:     "KZ",
								Probability: 0.05,
							},
						},
					}, {
						ID:         2,
						Name:       "name",
						Surname:    "surname1",
						Patronymic: "patronymic1",
						Age:        19,
						Gender:     "male",
						Nationalize: []entity.Nationality{
							{
								Country:     "RU",
								Probability: 0.1,
							},
							{
								Country:     "KZ",
								Probability: 0.05,
							},
						},
					},
				}}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name: "FailedWithErr",
			mockBehavior: func(s *mock_service.MockRepository) {
				s.EXPECT().GetPersons(gomock.Any(), gomock.Any()).Return(nil, errors.New("any error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(errorResponse{
					Err: entity.ErrInternalService.Error(),
				}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name: "FailedWithErr2",
			path: "?gender=mala",
			mockBehavior: func(s *mock_service.MockRepository) {
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(errorResponse{
					Err: entity.ErrInvalidInput.Error(),
				}, "", "    ")
				return string(resp)
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockRepository(c)
			tt.mockBehavior(repo)

			service := service.New(repo)

			handler := NewHandler(service)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/persons"+tt.path, nil)

			handler.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
