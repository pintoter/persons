package service

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/pintoter/persons/internal/entity"
	"github.com/pintoter/persons/pkg/logger"
)

func (s *Service) CreatePerson(ctx context.Context, person entity.Person) (int, error) {
	layer := "service.Create"

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errChan := make(chan error)
	go func() {
		wg.Add(3)
		go func() {
			defer wg.Done()
			var err error
			person.Age, err = s.gen.GenerateAge(ctx, person.Name)
			logger.DebugKV(ctx, "generate age", "layer", layer, "age", person.Age, "err", err)
			if err != nil {
				errChan <- entity.ErrInvalidInput
			}
		}()

		go func() {
			defer wg.Done()
			var err error
			person.Gender, err = s.gen.GenerateGender(ctx, person.Name)
			logger.DebugKV(ctx, "generate gender", "layer", layer, "gender", person.Gender, "err", err)
			if err != nil {
				errChan <- entity.ErrInvalidInput
			}
		}()

		go func() {
			defer wg.Done()
			var err error
			person.Nationalize, err = s.gen.GenerateNationalize(ctx, person.Name)
			logger.DebugKV(ctx, "generate nationalize", "layer", layer, "nationalize", person.Nationalize, "err", err)
			if err != nil || person.Nationalize == nil {
				errChan <- entity.ErrInvalidInput
			}
		}()

		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return 0, err
		}
	}

	id, err := s.repo.Create(ctx, person)
	if err != nil {
		return 0, err
	}

	return id, nil
}

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

type UpdateParams struct {
	Name       *string
	Surname    *string
	Patronymic *string
}

func (s *Service) Update(ctx context.Context, id int, params *UpdateParams) error {
	if !s.isPersonExists(ctx, id) {
		return entity.ErrPersonNotExists
	}

	return s.repo.Update(ctx, id, params)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if !s.isPersonExists(ctx, id) {
		return entity.ErrPersonNotExists
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) isPersonExists(ctx context.Context, id int) bool {
	_, err := s.repo.GetPerson(ctx, id)

	return err == nil
}
