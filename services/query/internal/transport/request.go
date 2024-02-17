package transport

import (
	"net/http"
	"strconv"

	"github.com/pintoter/persons/pkg/logger"
	"github.com/pintoter/persons/services/query/internal/entity"
)

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
	logger.DebugKV(r.Context(), "get persons request", "query", r.URL.Query())
	if r.URL.Query().Has("name") {
		p.Name = r.URL.Query().Get("name")
	}

	if r.URL.Query().Has("surname") {
		p.Surname = r.URL.Query().Get("surname")
	}

	if r.URL.Query().Has("patronymic") {
		p.Patronymic = r.URL.Query().Get("patronymic")
	}

	if r.URL.Query().Has("age") {
		age, err := strconv.Atoi(r.URL.Query().Get("age"))
		logger.DebugKV(r.Context(), "get persons request", "has age", age)
		if err != nil {
			logger.DebugKV(r.Context(), "get persons request", "err", err)
			return entity.ErrInvalidInput
		}
		p.Age = age
	}
	logger.DebugKV(r.Context(), "get persons request", "age", p.Age)

	if r.URL.Query().Has("gender") {
		gender := r.URL.Query().Get("gender")
		if gender != entity.Male && gender != entity.Female {
			logger.DebugKV(r.Context(), "get persons request", "err", "invalid gender")
			return entity.ErrInvalidInput
		}
		p.Gender = gender
	}

	if r.URL.Query().Has("nationalize") {
		p.Nationalize = r.URL.Query().Get("nationalize")
	}

	if r.URL.Query().Has("limit") {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || limit < 0 {
			logger.DebugKV(r.Context(), "get persons request", "err", err)
			return entity.ErrInvalidInput
		}
		p.Limit = limit
	} else {
		p.Limit = 5
	}

	if r.URL.Query().Has("page") {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page <= 0 {
			logger.DebugKV(r.Context(), "get persons request", "err", err)
			return entity.ErrInvalidInput
		}
		p.Page = page
	} else {
		p.Page = 1
	}

	return nil
}
