package common

import (
	"bdspro/internal/enums"
	"context"

	"gorm.io/gorm"
)

type IOwnerRepo[T any, D DTO] interface {
	QuerySearch(c context.Context, profileId uint64, ownerType enums.EOwnerOf, dto D) (*gorm.DB, error)
	Search(c context.Context, profileId uint64, ownerType enums.EOwnerOf, dto D) ([]T, int64, error)
	GetByID(c context.Context, id uint64) (*T, error)
	Detail(c context.Context, id uint64) (*T, error)
	Create(c context.Context, entity *T) error
	Update(c context.Context, id uint64, entity *T) error
	Delete(c context.Context, id uint64) error
	// GetAll(c context.Context) ([]T, error)
}
