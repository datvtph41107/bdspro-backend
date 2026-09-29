package repository

import (
	"context"
	"organization/internal/domain/entity"
)

type DealOfGroupRepo interface {
	Create(ctx context.Context, dealOfGroup *entity.DealOfGroup) error
	GetByGroupID(ctx context.Context, groupID uint64, page, size int) ([]entity.DealOfGroup, int64, error)
	GetByDealID(ctx context.Context, dealID uint64) (*entity.DealOfGroup, error)
	DeleteByDealID(ctx context.Context, dealID uint64) error
}
