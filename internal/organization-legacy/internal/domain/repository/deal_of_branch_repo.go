package repository

import (
	"context"
	"organization/internal/domain/entity"
)

type DealOfBranchRepo interface {
	Create(ctx context.Context, dealOfBranch *entity.DealOfBranch) error
	GetByBranchID(ctx context.Context, branchID uint64, page, size int) ([]entity.DealOfBranch, int64, error)
	GetByDealID(ctx context.Context, dealID uint64) (*entity.DealOfBranch, error)
	DeleteByDealID(ctx context.Context, dealID uint64) error
}
