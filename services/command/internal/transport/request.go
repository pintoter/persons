package transport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/pintoter/persons/services/command/internal/entity"
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

type updatePersonInput struct {
	ID         int    `json:"-"`
	Name       string `json:"name,omitempty"`
	Surname    string `json:"surname,omitempty"`
	Patronymic string `json:"patronymic,omitempty"`
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

	if p.Name == "" && p.Surname == "" && p.Patronymic == "" {
		return entity.ErrInvalidInput
	}
	return nil
}
