package transport

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pintoter/persons/internal/entity"
)

// @Summary Create note
// @Description create note
// @Tags notes
// @Accept json
// @Produce json
// @Param input body createNoteInput true "note info"
// @Success 201 {object} successCUDResponse
// @Failure 400 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/note [post]
func (h *Handler) createNote(w http.ResponseWriter, r *http.Request) {
	var input createNoteInput
	if err := input.Set(r); err != nil {
		renderJSON(w, r, http.StatusBadRequest, errorResponse{entity.ErrInvalidInput.Error()})
		return
	}

	id, err := h.service.Create(r.Context(), entity.Person{
		Name:  input,
		Title: input.Title,
	})

	if err != nil {
		if errors.Is(err, entity.ErrNoteExists) {
			renderJSON(w, r, http.StatusConflict, errorResponse{err.Error()})
		} else {
			renderJSON(w, r, http.StatusInternalServerError, errorResponse{err.Error()})
		}
		return
	}

	renderJSON(w, r, http.StatusCreated, successCUDResponse{"note created successfully"})
}

// @Summary Get note by id
// @Description Get note by id
// @Tags notes
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} getNoteResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/note/{id} [get]
func (h *Handler) getNote(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if id == 0 {
		renderJSON(w, r, http.StatusBadRequest, errorResponse{entity.ErrInvalidId.Error()})
		return
	}

	note, err := h.service.GetPerson(r.Context(), id)
	if err != nil {
		if errors.Is(err, entity.ErrPersonNotExists) {
			renderJSON(w, r, http.StatusNotFound, errorResponse{entity.ErrPersonNotExists.Error()})
		} else {
			renderJSON(w, r, http.StatusInternalServerError, errorResponse{err.Error()})
		}
		return
	}

	renderJSON(w, r, http.StatusOK, getNoteResponse{Note: note})
}

// @Summary Get all notes
// @Description Get all notes
// @Tags notes
// @Produce json
// @Success 200 {object} getNotesResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/notes [get]
func (h *Handler) getNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.service.GetPersons(r.Context(), r.Context().Value("user_id").(int))
	if err != nil {
		renderJSON(w, r, http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}

	renderJSON(w, r, http.StatusOK, getNotesResponse{Notes: notes})
}

// @Summary Update note
// @Description update note by id
// @Tags notes
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param input body updateNoteInput true "updating params"
// @Success 202 {object} successCUDResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/note/{id} [patch]
func (h *Handler) updateNote(w http.ResponseWriter, r *http.Request) {
	var input updateNoteInput
	var err error
	if err = input.Set(r); err != nil {
		renderJSON(w, r, http.StatusBadRequest, errorResponse{entity.ErrInvalidInput.Error()})
		return
	}

	userId, ok := r.Context().Value("user_id").(int)
	if !ok {
		renderJSON(w, r, http.StatusInternalServerError, errorResponse{"bad userID"})
		return
	}

	err = h.service.UpdateNote(r.Context(), input.ID, input.Title, input.Description, input.Status, userId)
	if err != nil {
		if errors.Is(err, entity.ErrNoteNotExists) {
			renderJSON(w, r, http.StatusBadRequest, errorResponse{entity.ErrNoteNotExists.Error()})
		} else if errors.Is(err, entity.ErrNoteExists) {
			renderJSON(w, r, http.StatusBadRequest, errorResponse{entity.ErrNoteExists.Error() + " with title: " + input.Title})
		} else {
			renderJSON(w, r, http.StatusInternalServerError, errorResponse{err.Error()})
		}
		return
	}

	renderJSON(w, r, http.StatusAccepted, successCUDResponse{Message: "note updated successfully"})
}

// @Summary Delete note
// @Description Delete note by id
// @Tags notes
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} successCUDResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/note/{id} [delete]
func (h *Handler) deleteNote(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if id == 0 {
		renderJSON(w, r, http.StatusBadRequest, errorResponse{entity.ErrInvalidId.Error()})
		return
	}

	userId, ok := r.Context().Value("user_id").(int)
	if !ok {
		renderJSON(w, r, http.StatusInternalServerError, errorResponse{"bad userID"})
		return
	}

	if err := h.service.DeleteNoteById(r.Context(), id, userId); err != nil {
		if errors.Is(err, entity.ErrNoteExists) {
			renderJSON(w, r, http.StatusBadRequest, errorResponse{entity.ErrNoteExists.Error()})
			return
		} else {
			renderJSON(w, r, http.StatusInternalServerError, errorResponse{err.Error()})
			return
		}
	}

	renderJSON(w, r, http.StatusOK, successCUDResponse{Message: "note deleted succesfully"})
}
