package crud3

import "context"

type CrudRepo[T any] interface {
	Create(c context.Context, entity *T) error
	Update(c context.Context, id uint64, entity *T) error
	Delete(c context.Context, id uint64) error
	GetAll() ([]T, error)
	GetByID(id uint64) (*T, error)
}
