package repo

import (
	"context"
	"notification/internal/domain"
	"notification/internal/dto"
)

type PropertyHistoryRepo interface {
	Create(ctx context.Context, entity *domain.PropertyHistory) (*domain.PropertyHistory, error)
	GetByID(ctx context.Context, id uint64) (*domain.PropertyHistory, error)
	Delete(ctx context.Context, id uint64) error
	Search(ctx context.Context, query *dto.PropertyHistoryQueryDTO) ([]*domain.PropertyHistory, int64, error)
}
