package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pintoter/persons/pkg/logger"
	"github.com/pintoter/persons/services/query/internal/entity"
)

func (s *Service) GetPerson(ctx context.Context, id int) (entity.Person, error) {
	layer := "service.GetPerson"

	person, err := s.repo.GetPerson(ctx, id)
	logger.DebugKV(ctx, "result get person", "layer", layer, "person", person, "err", err)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, entity.ErrPersonNotExists) {
			return entity.Person{}, entity.ErrPersonNotExists
		} else {
			return entity.Person{}, entity.ErrInternalService
		}
	}

	return person, nil
}

type GetFilters struct {
	Name        *string
	Surname     *string
	Patronymic  *string
	Age         *int
	Gender      *string
	Nationalize *string
	Limit       int64
	Offset      int64
}

func (s *Service) GetPersons(ctx context.Context, filters *GetFilters) ([]entity.Person, error) {
	layer := "service.GetPersons"

	persons, err := s.repo.GetPersons(ctx, filters)
	logger.DebugKV(ctx, "get persons request", "layer", layer, "persons", persons)
	if err != nil {
		return nil, entity.ErrInternalService
	}

	return persons, nil
}
