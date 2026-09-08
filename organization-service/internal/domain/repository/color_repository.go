package repository

import (
	"context"
	"organization/internal/domain/entity"
)

type ColorRepository interface {
	Create(ctx context.Context, color *entity.Color) error
	Update(ctx context.Context, color *entity.Color) error
	Delete(ctx context.Context, id uint32) error
	GetByID(ctx context.Context, id uint32) (*entity.Color, error)
	GetAll(ctx context.Context, page, size int, isActive *bool) ([]*entity.Color, int64, error)
} 