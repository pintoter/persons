package transport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/pintoter/persons/internal/entity"
	"github.com/pintoter/persons/pkg/logger"
)

type createPersonInput struct {
	Name       string `json:"name" binding:"required,min=2,max=64"`
	Surname    string `json:"surname" binding:"required,min=2,max=64"`
	Patronymic string `json:"patronymic,omitempty"`
}

func (p *createPersonInput) Set(r *http.Request) error {
	var err error
	if err = json.NewDecoder(r.Body).Decode(p); err != nil {
		return entity.ErrInvalidInput
	}

	return nil
}

type getPersonsRequest struct {
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	Gender      string
	Nationalize string
	Limit       int
	Page        int
}

func (p *getPersonsRequest) Set(r *http.Request) error {
	var err error

	p.Name = r.URL.Query().Get("name")
	p.Surname = r.URL.Query().Get("surname")

	if r.URL.Query().Has("patronymic") {
		p.Patronymic = r.URL.Query().Get("patronymic")
	}

	if r.URL.Query().Has("age") {
		p.Age, err = strconv.Atoi(r.URL.Query().Get("age"))
		if err != nil {
			logger.DebugKV(r.Context(), "get persons request", "err", err)
			return entity.ErrInvalidInput
		}
	}

	p.Gender = r.URL.Query().Get("gender")
	if p.Gender != entity.Male && p.Gender != entity.Female {
		logger.DebugKV(r.Context(), "get persons request", "err", err)
		return entity.ErrInvalidInput
	}

	p.Nationalize = r.URL.Query().Get("nationalize")

	if r.URL.Query().Has("limit") {
		p.Limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			logger.DebugKV(r.Context(), "get persons request", "err", err)
			return entity.ErrInvalidInput
		}

		if p.Limit < 0 {
			return entity.ErrInvalidInput
		}

		if p.Limit == 0 {
			p.Limit = 5
		}
	}

	if r.URL.Query().Has("page") {
		p.Page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			logger.DebugKV(r.Context(), "get persons request", "err", err)
			return entity.ErrInvalidInput
		}

		if p.Page <= 0 {
			logger.DebugKV(r.Context(), "get persons request", "err", err)
			return entity.ErrInvalidInput
		}
	}

	logger.DebugKV(r.Context(), "get persons request", "err", err)
	return nil
}

type updatePersonInput struct {
	ID          int    `json:"-"`
	Name        string `json:"name,omitempty"`
	Surname     string `json:"surname,omitempty"`
	Patronymic  string `json:"patronymic,omitempty"`
	Age         int    `json:"age,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Nationalize string `json:"nationalize,omitempty"`
}

func (p *updatePersonInput) Set(r *http.Request) error {
	p.ID, _ = strconv.Atoi(mux.Vars(r)["id"])
	if p.ID == 0 {
		return entity.ErrInvalidQueryId
	}

	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		return entity.ErrInvalidInput
	}

	if p.Name == "" && p.Surname == "" && p.Age <= 0 && p.Gender == "" && p.Nationalize == "" {
		return entity.ErrInvalidInput
	}
	return nil
}
