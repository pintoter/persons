package service

import (
	"context"

	"github.com/pintoter/persons/services/query/internal/entity"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Repository interface {
	GetPerson(ctx context.Context, id int) (entity.Person, error)
	GetPersons(ctx context.Context, filters *GetFilters) ([]entity.Person, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}
