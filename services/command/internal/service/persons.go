package service

import (
	"context"
	"sync"

	"github.com/pintoter/persons/pkg/logger"
	"github.com/pintoter/persons/services/command/internal/entity"
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

type UpdateParams struct {
	Name       *string
	Surname    *string
	Patronymic *string
}

func (s *Service) Update(ctx context.Context, id int, params *UpdateParams) error {
	return s.repo.Update(ctx, id, params)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
