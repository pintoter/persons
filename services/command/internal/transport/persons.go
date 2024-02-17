package transport

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pintoter/persons/services/command/internal/entity"
	"github.com/pintoter/persons/services/command/internal/service"
)

// @Summary Create person
// @Description Create person
// @Tags persons
// @Accept json
// @Produce json
// @Param input body createPersonInput true "Person's information"
// @Success 201 {object} successResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/persons [post]
func (h *Handler) createPerson(w http.ResponseWriter, r *http.Request) {
	var input createPersonInput
	if err := input.Set(r); err != nil {
		renderJSON(w, r, http.StatusBadRequest, errorResponse{entity.ErrInvalidInput.Error()})
		return
	}

	id, err := h.service.CreatePerson(r.Context(), entity.Person{
		Name:       input.Name,
		Surname:    input.Surname,
		Patronymic: input.Patronymic,
	})

	if err != nil {
		renderJSON(w, r, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	renderJSON(w, r, http.StatusCreated, successResponse{fmt.Sprintf("created new person ID: %d", id)})
}

// @Summary Update persons
// @Description update person by id
// @Tags persons
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param input body updatePersonInput true "updating params"
// @Success 202 {object} successResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/persons/{id} [patch]
func (h *Handler) updatePerson(w http.ResponseWriter, r *http.Request) {
	var input updatePersonInput
	var err error
	if err = input.Set(r); err != nil {
		renderJSON(w, r, http.StatusBadRequest, errorResponse{err.Error()})
		return
	}

	var data service.UpdateParams
	if input.Name != "" {
		data.Name = &input.Name
	}

	if input.Surname != "" {
		data.Surname = &input.Surname
	}

	if input.Patronymic != "" {
		data.Patronymic = &input.Patronymic
	}

	err = h.service.Update(r.Context(), input.ID, &data)
	if err != nil {
		if errors.Is(err, entity.ErrPersonNotExists) {
			renderJSON(w, r, http.StatusBadRequest, errorResponse{entity.ErrPersonNotExists.Error()})
		} else {
			renderJSON(w, r, http.StatusInternalServerError, errorResponse{err.Error()})
		}
		return
	}

	renderJSON(w, r, http.StatusAccepted, successResponse{Message: "person updated successfully"})
}

// @Summary Delete person
// @Description Delete person by id
// @Tags persons
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} successResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/person/{id} [delete]
func (h *Handler) deletePerson(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if id == 0 {
		renderJSON(w, r, http.StatusBadRequest, errorResponse{entity.ErrInvalidInput.Error()})
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, entity.ErrPersonNotExists) {
			renderJSON(w, r, http.StatusBadRequest, errorResponse{entity.ErrPersonNotExists.Error()})
		} else {
			renderJSON(w, r, http.StatusInternalServerError, errorResponse{err.Error()})
		}
		return
	}

	renderJSON(w, r, http.StatusOK, successResponse{Message: "person deleted succesfully"})
}
