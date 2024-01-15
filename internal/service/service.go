package service

import (
	"github.com/pintoter/persons/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Client interface {
	GetAge(name string) (int, error)
	GetGender(name string) (string, error)
	GetNationalize(name string) (string, error)
}

type Service struct {
	repo   repository.Repository
	client Client
}

func New(repo repository.Repository, client Client) *Service {
	return &Service{
		repo:   repo,
		client: client,
	}
}
