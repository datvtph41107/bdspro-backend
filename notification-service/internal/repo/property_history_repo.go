package repo

import (
	"context"
	"notification/internal/domain"
	"notification/internal/dto"
)

type PropertyHistoryRepo interface {
	Create(ctx context.Context, entity *domain.PropertyHistoryEntity) (*domain.PropertyHistoryEntity, error)
	GetByID(ctx context.Context, id uint64) (*domain.PropertyHistoryEntity, error)
	Delete(ctx context.Context, id uint64) error
	Search(ctx context.Context, query dto.PropertyHistoryQueryDTO) ([]domain.PropertyHistoryEntity, int64, error)
}
