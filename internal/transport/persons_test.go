package transport

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pintoter/persons/internal/entity"
	"github.com/pintoter/persons/internal/service"
	mock_service "github.com/pintoter/persons/internal/service/mocks"

	"github.com/stretchr/testify/assert"
)

func Test_CreatePersonHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockRepository, b *mock_service.MockGenerator, person entity.Person)

	tests := []struct {
		name                 string
		inputBody            string
		inputPerson          entity.Person
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Success",
			inputBody: `{
					"name": "Ivan",
					"surname": "Ivanov",
					"patronymic": "Ivanovich"
				}`,
			inputPerson: entity.Person{
				Name:       "Ivan",
				Surname:    "Ivanov",
				Patronymic: "Ivanovich",
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
			},
			mockBehavior: func(r *mock_service.MockRepository, g *mock_service.MockGenerator, person entity.Person) {
				g.EXPECT().GenerateAge(gomock.Any(), person.Name).Times(1).Return(18, nil)
				g.EXPECT().GenerateGender(gomock.Any(), person.Name).Times(1).Return("male", nil)
				g.EXPECT().
					GenerateNationalize(gomock.Any(), person.Name).Times(1).
					Return([]entity.Nationality{{Country: "RU", Probability: 0.1}, {Country: "KZ", Probability: 0.05}}, nil)
				r.EXPECT().Create(gomock.Any(), person).Return(1, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(successResponse{Message: "created new person ID: 1"}, "", "    ")
				return string(resp)
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockRepository(c)
			gen := mock_service.NewMockGenerator(c)
			tt.mockBehavior(repo, gen, tt.inputPerson)

			service := service.New(repo, gen)

			handler := NewHandler(service)

			// Create request
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/persons", bytes.NewBufferString(tt.inputBody))

			handler.ServeHTTP(w, r)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}

func Test_GetPersonHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockRepository, g *mock_service.MockGenerator, id int)

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
			mockBehavior: func(s *mock_service.MockRepository, g *mock_service.MockGenerator, id int) {
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
			mockBehavior: func(s *mock_service.MockRepository, g *mock_service.MockGenerator, id int) {
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
			mockBehavior: func(s *mock_service.MockRepository, g *mock_service.MockGenerator, id int) {
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
			mockBehavior: func(s *mock_service.MockRepository, g *mock_service.MockGenerator, id int) {
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
			gen := mock_service.NewMockGenerator(c)
			tt.mockBehavior(repo, gen, tt.inputId)

			service := service.New(repo, gen)

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
	type mockBehavior func(s *mock_service.MockRepository, g *mock_service.MockGenerator)

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
			mockBehavior: func(s *mock_service.MockRepository, g *mock_service.MockGenerator) {
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
			mockBehavior: func(s *mock_service.MockRepository, g *mock_service.MockGenerator) {
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
			mockBehavior: func(s *mock_service.MockRepository, g *mock_service.MockGenerator) {
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
			mockBehavior: func(s *mock_service.MockRepository, g *mock_service.MockGenerator) {
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
			gen := mock_service.NewMockGenerator(c)
			tt.mockBehavior(repo, gen)

			service := service.New(repo, gen)

			handler := NewHandler(service)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/persons"+tt.path, nil)

			handler.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func Test_Delete(t *testing.T) {
	type mockBehavior func(s *mock_service.MockRepository, g *mock_service.MockGenerator, id int)

	tests := []struct {
		name                 string
		id                   int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			id:   1,
			mockBehavior: func(s *mock_service.MockRepository, g *mock_service.MockGenerator, id int) {
				s.EXPECT().GetPerson(gomock.Any(), id).Return(entity.Person{}, nil)
				s.EXPECT().Delete(gomock.Any(), id).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(successResponse{Message: "person deleted succesfully"}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name:               "FailedWithInput",
			id:                 0,
			mockBehavior:       func(s *mock_service.MockRepository, g *mock_service.MockGenerator, id int) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(errorResponse{Err: entity.ErrInvalidInput.Error()}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name: "FailedWithNoPerson",
			id:   5,
			mockBehavior: func(s *mock_service.MockRepository, g *mock_service.MockGenerator, id int) {
				s.EXPECT().GetPerson(gomock.Any(), id).Return(entity.Person{}, sql.ErrNoRows)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(errorResponse{Err: entity.ErrPersonNotExists.Error()}, "", "    ")
				return string(resp)
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockRepository(c)
			gen := mock_service.NewMockGenerator(c)
			tt.mockBehavior(repo, gen, tt.id)

			service := service.New(repo, gen)

			handler := NewHandler(service)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/persons/"+fmt.Sprintf("%d", tt.id), nil)

			handler.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func Test_Update(t *testing.T) {
	type mockBehavior func(s *mock_service.MockRepository, g *mock_service.MockGenerator, id int)

	tests := []struct {
		name                 string
		id                   int
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			id:        1,
			inputBody: `{"name": "Ivan"}`,
			mockBehavior: func(s *mock_service.MockRepository, g *mock_service.MockGenerator, id int) {
				s.EXPECT().GetPerson(gomock.Any(), id).Return(entity.Person{}, nil)
				s.EXPECT().Update(gomock.Any(), id, gomock.Any()).Return(nil)
			},
			expectedStatusCode: http.StatusAccepted,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(successResponse{Message: "person updated successfully"}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name:      "Failed",
			id:        0,
			inputBody: `{"name": "Ivan"}`,
			mockBehavior: func(s *mock_service.MockRepository, g *mock_service.MockGenerator, id int) {
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(errorResponse{Err: entity.ErrInvalidQueryId.Error()}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name:      "FailedWithId",
			id:        5,
			inputBody: `{"name": "Ivan"}`,
			mockBehavior: func(s *mock_service.MockRepository, g *mock_service.MockGenerator, id int) {
				s.EXPECT().GetPerson(gomock.Any(), id).Return(entity.Person{}, sql.ErrNoRows)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(errorResponse{Err: entity.ErrPersonNotExists.Error()}, "", "    ")
				return string(resp)
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockRepository(c)
			gen := mock_service.NewMockGenerator(c)
			tt.mockBehavior(repo, gen, tt.id)

			service := service.New(repo, gen)

			handler := NewHandler(service)

			// Create request
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PATCH", "/api/v1/persons/"+fmt.Sprintf("%d", tt.id), bytes.NewBufferString(tt.inputBody))

			handler.ServeHTTP(w, r)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}
