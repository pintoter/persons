package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sync"

	"github.com/pintoter/persons/internal/entity"
)

func (s *Service) Create(ctx context.Context, person entity.Person) (int, error) {
	var wg sync.WaitGroup

	errChan := make(chan error, 1)

	wg.Add(3)
	go func() {
		defer wg.Done()
		var err error
		person.Age, err = s.client.GetAge(person.Name)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		person.Gender, err = s.client.GetGender(person.Name)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		person.Nationalize, err = s.client.GetNationalize(person.Name)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()

	go func() { close(errChan) }()

	for err := range errChan {
		log.Println(err)
		return 0, err
	}

	id, err := s.repo.Create(ctx, person)

	return id, err
}

func (s *Service) GetPerson(ctx context.Context, id int) (entity.Person, error) {
	person, err := s.repo.GetPerson(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Person{}, entity.ErrPersonNotExists
		} else {
			return entity.Person{}, entity.ErrInternalService
		}
	}

	return person, nil
}

func (s *Service) GetPersons(ctx context.Context, limit, offset int) ([]entity.Person, error) {
	Persons, err := s.repo.GetPersons(ctx, limit, offset)
	if err != nil {
		return nil, entity.ErrInternalService
	}

	return Persons, nil
}

func (s *Service) Update(ctx context.Context, id int) error {
	if !s.isPersonExists(ctx, id) {
		return entity.ErrPersonNotExists
	}

	return s.repo.Update(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if s.isPersonExists(ctx, id) {
		return s.repo.Delete(ctx, id)
	}

	return entity.ErrPersonNotExists
}

func (s *Service) isPersonExists(ctx context.Context, id int) bool {
	_, err := s.repo.GetPerson(ctx, id)

	return err == nil
}
