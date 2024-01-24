package service

import (
	"context"

	"github.com/pintoter/persons/internal/entity"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Repository interface {
	Create(ctx context.Context, person entity.Person) (int, error)
	GetPerson(ctx context.Context, id int) (entity.Person, error)
	GetPersons(ctx context.Context, filters *GetFilters) ([]entity.Person, error)
	Update(ctx context.Context, id int, params *UpdateParams) error
	Delete(ctx context.Context, id int) error
}

type Generator interface {
	GenerateAge(ctx context.Context, name string) (int, error)
	GenerateGender(ctx context.Context, name string) (string, error)
	GenerateNationalize(ctx context.Context, name string) ([]entity.Nationality, error)
}

type Service struct {
	repo Repository
	gen  Generator
}

func New(repo Repository, gen Generator) *Service {
	return &Service{
		repo: repo,
		gen:  gen,
	}
}
