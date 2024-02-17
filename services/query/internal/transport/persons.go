package transport

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pintoter/persons/pkg/logger"
	"github.com/pintoter/persons/services/query/internal/entity"
	"github.com/pintoter/persons/services/query/internal/service"
)

// @Summary Get person by id
// @Description Get person by id
// @Tags persons
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} getPersonResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/persons/{id} [get]
func (h *Handler) getPerson(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if id == 0 {
		renderJSON(w, r, http.StatusBadRequest, errorResponse{entity.ErrInvalidInput.Error()})
		return
	}

	person, err := h.service.GetPerson(r.Context(), id)
	if err != nil {
		if errors.Is(err, entity.ErrPersonNotExists) {
			renderJSON(w, r, http.StatusNotFound, errorResponse{entity.ErrPersonNotExists.Error()})
		} else {
			renderJSON(w, r, http.StatusInternalServerError, errorResponse{err.Error()})
		}
		return
	}

	renderJSON(w, r, http.StatusOK, getPersonResponse{Person: person})
}

// @Summary Get all persons
// @Description Get all persons
// @Tags persons
// @Produce json
// @Param name query string false "name"
// @Param surname query string false "surname"
// @Param patronymic query string false "patronymic"
// @Param age query int false "age"
// @Param gender query string false "gender"
// @Param nationalize query string false "nationalize"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Success 200 {object} getPersonsResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/persons [get]
func (h *Handler) getPersons(w http.ResponseWriter, r *http.Request) {
	var input getPersonsRequest

	if err := input.Set(r); err != nil {
		renderJSON(w, r, http.StatusBadRequest, errorResponse{err.Error()})
		return
	}
	logger.DebugKV(r.Context(), "get persons request", "input", input)

	data := &service.GetFilters{}
	convertInputToGetFilters(data, &input)

	logger.DebugKV(r.Context(), "get persons request", "input filters", data)

	persons, err := h.service.GetPersons(r.Context(), data)
	if err != nil {
		renderJSON(w, r, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	renderJSON(w, r, http.StatusOK, getPersonsResponse{Persons: persons})
}

func convertInputToGetFilters(data *service.GetFilters, input *getPersonsRequest) {
	if input.Name != "" {
		data.Name = &input.Name
	}

	if input.Surname != "" {
		data.Surname = &input.Surname
	}

	if input.Patronymic != "" {
		data.Patronymic = &input.Patronymic
	}

	if input.Age != 0 {
		data.Age = &input.Age
	}

	if input.Gender != "" {
		data.Gender = &input.Gender
	}

	if input.Nationalize != "" {
		data.Nationalize = &input.Nationalize
	}

	data.Limit = int64(input.Limit)
	data.Offset = (int64(input.Page) - 1) * data.Limit
}
