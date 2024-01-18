package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
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
				Name:        "Ivan",
				Surname:     "Ivanov",
				Patronymic:  "Ivanovich",
				Age:         18,
				Gender:      "male",
				Nationalize: "RU",
			},
			mockBehavior: func(r *mock_service.MockRepository, g *mock_service.MockGenerator, person entity.Person) {
				g.EXPECT().GenerateAge(gomock.Any(), person.Name).Times(1).Return(18, nil)
				g.EXPECT().GenerateGender(gomock.Any(), person.Name).Times(1).Return("male", nil)
				g.EXPECT().GenerateNationalize(gomock.Any(), person.Name).Times(1).Return("RU", nil)
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

			service := &service.Service{
				Repository: repo,
				Generator:  gen,
			}

			handler := Handler{}
			handler.service = service

			// Init endpoint
			mux := mux.NewRouter()
			mux.HandleFunc("/persons", handler.createPerson).Methods(http.MethodPost)

			// Create request
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/persons", bytes.NewBufferString(tt.inputBody))

			mux.ServeHTTP(w, r)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}

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
				s.EXPECT().GetPerson(context.Background(), id).Return(entity.Person{
					ID:          1,
					Name:        "name",
					Surname:     "surname",
					Patronymic:  "patronymic",
					Age:         18,
					Gender:      "male",
					Nationalize: "RU",
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(getPersonResponse{Person: entity.Person{
					ID:          1,
					Name:        "name",
					Surname:     "surname",
					Patronymic:  "patronymic",
					Age:         18,
					Gender:      "male",
					Nationalize: "RU",
				}}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name:    "FailedNotExist",
			inputId: 2,
			mockBehavior: func(s *mock_service.MockRepository, id int) {
				s.EXPECT().GetPerson(context.Background(), id).Return(entity.Person{}, entity.ErrPersonNotExists)
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
			name:               "FailedWithId",
			inputId:            0,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(errorResponse{
					Err: entity.ErrPersonNotExists.Error(),
				}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name:    "FailedWithErr",
			inputId: 2,
			mockBehavior: func(s *mock_service.MockRepository, id int) {
				s.EXPECT().GetPerson(context.Background(), id).Return(entity.Person{}, entity.ErrPersonNotExists)
			},
			expectedStatusCode: http.StatusInternalServerError,
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

			service := &service.Service{
				Repository: repo,
			}

			handler := &Handler{
				service: service,
			}

			mux := mux.NewRouter()
			mux.HandleFunc("/persons/{id}", handler.getPerson).Methods(http.MethodGet)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/persons/1", nil)

			mux.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func Test_GetPersonsHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockRepository, data *service.GetFilters)

	tests := []struct {
		name                 string
		path                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			path: "/?nationalize=RU",
			mockBehavior: func(s *mock_service.MockRepository, data *service.GetFilters) {
				s.EXPECT().GetPersons(context.Background(), data).Return(
					[]entity.Person{
						{
							ID:          1,
							Name:        "name",
							Surname:     "surname",
							Patronymic:  "patronymic",
							Age:         18,
							Gender:      "male",
							Nationalize: "RU",
						}, {
							ID:          2,
							Name:        "name1",
							Surname:     "surname1",
							Patronymic:  "patronymic1",
							Age:         19,
							Gender:      "male",
							Nationalize: "RU",
						},
					},
					nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(getPersonsResponse{Persons: []entity.Person{
					{
						ID:          1,
						Name:        "name",
						Surname:     "surname",
						Patronymic:  "patronymic",
						Age:         18,
						Gender:      "male",
						Nationalize: "RU",
					}, {
						ID:          2,
						Name:        "name1",
						Surname:     "surname1",
						Patronymic:  "patronymic1",
						Age:         19,
						Gender:      "male",
						Nationalize: "RU",
					},
				}}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name: "Ok2",
			path: "/?name=name",
			mockBehavior: func(s *mock_service.MockRepository, data *service.GetFilters) {
				s.EXPECT().GetPersons(context.Background(), data).Return(
					[]entity.Person{
						{
							ID:          1,
							Name:        "name",
							Surname:     "surname",
							Patronymic:  "patronymic",
							Age:         18,
							Gender:      "male",
							Nationalize: "RU",
						}, {
							ID:          2,
							Name:        "name",
							Surname:     "surname1",
							Patronymic:  "patronymic1",
							Age:         19,
							Gender:      "male",
							Nationalize: "RU",
						},
					},
					nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(getPersonsResponse{Persons: []entity.Person{
					{
						ID:          1,
						Name:        "name",
						Surname:     "surname",
						Patronymic:  "patronymic",
						Age:         18,
						Gender:      "male",
						Nationalize: "RU",
					}, {
						ID:          2,
						Name:        "name1",
						Surname:     "surname1",
						Patronymic:  "patronymic1",
						Age:         19,
						Gender:      "male",
						Nationalize: "RU",
					},
				}}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name: "FailedWithErr",
			mockBehavior: func(s *mock_service.MockRepository, data *service.GetFilters) {
				s.EXPECT().GetPersons(context.Background(), data).Return(nil, errors.New("any error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(errorResponse{
					Err: "any error",
				}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name:               "FailedWithErr2",
			path:               "/?gender=mala",
			expectedStatusCode: http.StatusInternalServerError,
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

			filters := &service.GetFilters{}
			repo := mock_service.NewMockRepository(c)
			tt.mockBehavior(repo, filters)

			service := &service.Service{
				Repository: repo,
			}

			handler := &Handler{
				service: service,
			}

			mux := mux.NewRouter()
			mux.HandleFunc("/persons", handler.getPersons).Methods(http.MethodGet)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/notes"+tt.path, nil)

			mux.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func Test_Delete(t *testing.T) {
	type mockBehavior func(s *mock_service.MockRepository, id int)

	tests := []struct {
		name                 string
		path                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			path: "1",
			mockBehavior: func(s *mock_service.MockRepository, id int) {
				s.EXPECT().Delete(context.Background(), id).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(successResponse{Message: "note deleted succesfully"}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name:               "FailedWithInput",
			path:               "0",
			mockBehavior:       func(s *mock_service.MockRepository, id int) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(errorResponse{Err: entity.ErrInvalidInput.Error()}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name: "FailedWithNoPerson",
			path: "5",
			mockBehavior: func(s *mock_service.MockRepository, id int) {
				s.EXPECT().Delete(context.Background(), id).Return(entity.ErrPersonNotExists)
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
			tt.mockBehavior(repo, 1)

			service := &service.Service{
				Repository: repo,
			}

			handler := &Handler{
				service: service,
			}

			mux := mux.NewRouter()
			mux.HandleFunc("/note/{id}", handler.deletePerson).Methods(http.MethodDelete)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/persons/"+tt.path, nil)

			mux.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestUpdateNoteHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockRepository, id int, params *service.UpdateParams)

	tests := []struct {
		name                 string
		id                   int
		params               *service.UpdateParams
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			id:   1,
			inputBody: `{
					"name": "Ivan",
				}`,
			mockBehavior: func(s *mock_service.MockRepository, id int, params *service.UpdateParams) {
				s.EXPECT().Update(context.Background(), id, params).Return(nil)
			},
			expectedStatusCode: http.StatusAccepted,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(successResponse{Message: "note updated successfully"}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name: "Failed",
			id:   0,
			inputBody: `{
					"name": "Ivan",
				}`,
			mockBehavior:       func(s *mock_service.MockRepository, id int, params *service.UpdateParams) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: func() string {
				resp, _ := json.MarshalIndent(errorResponse{Err: entity.ErrInvalidInput.Error()}, "", "    ")
				return string(resp)
			}(),
		},
		{
			name: "FailedWithId",
			id:   5,
			inputBody: `{
					"name": "Ivan",
				}`,
			mockBehavior:       func(s *mock_service.MockRepository, id int, params *service.UpdateParams) {},
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
			tt.mockBehavior(repo, tt.id, tt.params)

			service := &service.Service{
				Repository: repo,
			}

			handler := &Handler{
				service: service,
			}

			// Init endpoint
			mux := mux.NewRouter()
			mux.HandleFunc("/note/{id}", handler.updatePerson).Methods(http.MethodPatch)

			// Create request
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PATCH", "/note/1", bytes.NewBufferString(tt.inputBody))

			mux.ServeHTTP(w, r)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}
