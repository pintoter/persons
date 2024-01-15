package transport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pintoter/todo-list/internal/entity"
)

type createNoteInput struct {
	Title         string    `json:"title" binding:"required,min=1,max=80"`
	Description   string    `json:"description,omitempty"`
	Date          string    `json:"date,omitempty" binding:"min=9,max=10"`
	DateFormatted time.Time `json:"-"`
	Status        string    `json:"status,omitempty"`
}

func (n *createNoteInput) Set(r *http.Request) error {
	var err error
	if err = json.NewDecoder(r.Body).Decode(n); err != nil {
		return entity.ErrInvalidInput
	}

	if n.Date != "" {
		n.DateFormatted, err = time.Parse(dateFormat, n.Date)
		if err != nil {
			return err
		}
	} else {
		n.DateFormatted = time.Now()
		if err != nil {
			return err
		}
	}

	if n.Title == "" {
		return entity.ErrInvalidInput
	}

	if n.Status == "" {
		n.Status = entity.StatusNotDone
	}

	return nil
}

type updateNoteInput struct {
	ID          int    `json:"-"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

func (n *updateNoteInput) Set(r *http.Request) error {
	n.ID, _ = strconv.Atoi(mux.Vars(r)["id"])
	if n.ID == 0 {
		return entity.ErrInvalidId
	}

	err := json.NewDecoder(r.Body).Decode(n)
	if err != nil {
		return entity.ErrInvalidInput
	}

	if n.Title == "" && n.Description == "" && n.Status == "" {
		return entity.ErrInvalidInput
	}
	return nil
}

type getNotesRequest struct {
	Page          int       `json:"-"`
	Status        string    `json:"status,omitempty"`
	Date          string    `json:"date,omitempty"`
	DateFormatted time.Time `json:"-"`
	Limit         int       `json:"limit,omitempty"`
}

func (n *getNotesRequest) Set(r *http.Request) error {
	var err error
	n.Page, _ = strconv.Atoi(mux.Vars(r)["page"])
	if n.Page == 0 {
		return entity.ErrInvalidPage
	}

	err = json.NewDecoder(r.Body).Decode(n)
	if err != nil {
		return entity.ErrInvalidInput
	}

	if n.Limit < 0 {
		return entity.ErrInvalidInput
	}

	if n.Limit == 0 {
		n.Limit = 5
	}

	if n.Date != "" {
		n.DateFormatted, err = time.Parse(dateFormat, n.Date)
		if err != nil {
			return entity.ErrInvalidDate
		}
	}

	if n.Status != "" && n.Status != entity.StatusDone && n.Status != entity.StatusNotDone {
		return entity.ErrInvalidStatus
	}

	return nil
}
