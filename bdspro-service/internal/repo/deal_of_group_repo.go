package repo

import (
	"bdspro/internal/domain"
	"context"
)

type DealOfGroupRepo interface {
	Create(ctx context.Context, dealOfGroup *domain.DealOfGroup) error
	GetByGroupID(ctx context.Context, groupID uint64, page, size int) ([]domain.DealOfGroup, int64, error)
	GetByDealID(ctx context.Context, dealID uint64) (*domain.DealOfGroup, error)
	DeleteByDealID(ctx context.Context, dealID uint64) error
}
