package repository

import (
	"context"

	"github.com/pintoter/persons/internal/entity"
)

type PersonRepository interface {
	Create(ctx context.Context, person entity.Person) (int, error)
	GetPerson(ctx context.Context, id int) (entity.Person, error)
	GetPersons(ctx context.Context, limit, offset int) ([]entity.Person, error)
	Update(ctx context.Context, id int) error
	Delete(ctx context.Context, id int) error
}

type Repository interface {
	PersonRepository
}
