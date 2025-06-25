package repository

import (
	"context"

	"organization/internal/domain/entity"
)

type GroupRepository interface {
	Create(ctx context.Context, group *entity.Group) (*entity.Group, error)
	Update(ctx context.Context, group *entity.Group) (*entity.Group, error)
	Delete(ctx context.Context, id uint32) error
	GetByID(ctx context.Context, id uint32) (*entity.Group, error)
	List(ctx context.Context, offset, limit int, currentUserId uint32) ([]*entity.Group, uint32, error)
	FindByUserId(ctx context.Context, userId uint32) ([]*entity.Group, error)
}
