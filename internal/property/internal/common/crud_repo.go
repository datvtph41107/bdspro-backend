package common

import (
	"bdspro/internal/enums"
	"context"

	"gorm.io/gorm"
)

type DTO interface {
	GetOffset() int
	GetLimit() int
}

type IBaseRepo[T any, D DTO] interface {
	Search(c context.Context, profileId uint64, ownerType enums.EOwnerOf, dto D) ([]T, int64, error)
	GetByID(c context.Context, id uint64) (*T, error)
	Detail(c context.Context, id uint64) (*T, error)
	Create(c context.Context, entity *T) error
	Update(c context.Context, id uint64, entity *T) error
	Delete(c context.Context, id uint64) error
	GetAll(c context.Context, dto D) ([]T, error)
	QuerySearch(c context.Context, profileId uint64, ownerType enums.EOwnerOf, dto D) (*gorm.DB, error)
}
