package service

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/pintoter/persons/internal/entity"
	"github.com/pintoter/persons/pkg/logger"
)

func (s *Service) Create(ctx context.Context, person entity.Person) (int, error) {
	var wg sync.WaitGroup

	errChan := make(chan error, 1)

	wg.Add(3)
	go func() {
		defer wg.Done()
		var err error
		person.Age, err = s.GenerateAge(ctx, person.Name)
		logger.DebugKV(ctx, "generate age", "age", person.Age, "err", err)
		if err != nil {
			errChan <- entity.ErrInvalidInput
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		person.Gender, err = s.GenerateGender(ctx, person.Name)
		logger.DebugKV(ctx, "generate gender", "gender", person.Gender, "err", err)
		if err != nil {
			errChan <- entity.ErrInvalidInput
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		person.Nationalize, err = s.GenerateNationalize(ctx, person.Name)
		logger.DebugKV(ctx, "generate gender", "nationalize", person.Nationalize, "err", err)
		if err != nil {
			errChan <- entity.ErrInvalidInput
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return 0, err
		}
	}

	id, err := s.Create(ctx, person)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) GetPerson(ctx context.Context, id int) (entity.Person, error) {
	person, err := s.GetPerson(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
	persons, err := s.GetPersons(ctx, filters)
	if err != nil {
		return nil, entity.ErrInternalService
	}

	return persons, nil
}

type UpdateParams struct {
	Name        *string
	Surname     *string
	Patronymic  *string
	Age         *int
	Gender      *string
	Nationalize *string
}

func (s *Service) Update(ctx context.Context, id int, params *UpdateParams) error {
	if !s.isPersonExists(ctx, id) {
		return entity.ErrPersonNotExists
	}

	return s.Update(ctx, id, params)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if !s.isPersonExists(ctx, id) {
		return entity.ErrPersonNotExists
	}

	return s.Delete(ctx, id)
}

func (s *Service) isPersonExists(ctx context.Context, id int) bool {
	_, err := s.GetPerson(ctx, id)

	return err == nil
}
