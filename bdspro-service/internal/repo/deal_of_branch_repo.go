package repo

import (
	"bdspro/internal/domain"
	"context"
)

type DealOfBranchRepo interface {
	Create(ctx context.Context, dealOfBranch *domain.DealOfBranch) error
	GetByBranchID(ctx context.Context, branchID uint64, page, size int) ([]domain.DealOfBranch, int64, error)
	GetByDealID(ctx context.Context, dealID uint64) (*domain.DealOfBranch, error)
	DeleteByDealID(ctx context.Context, dealID uint64) error
}
