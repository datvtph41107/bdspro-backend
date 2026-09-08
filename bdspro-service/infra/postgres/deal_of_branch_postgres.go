package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"

	"gorm.io/gorm"
)

type DealOfBranchPostgres struct {
	db *gorm.DB
}

func NewDealOfBranchPostgres(db *gorm.DB) repo.DealOfBranchRepo {
	return &DealOfBranchPostgres{db: db}
}

func (r *DealOfBranchPostgres) Create(ctx context.Context, dealOfBranch *domain.DealOfBranch) error {
	return r.db.WithContext(ctx).Create(dealOfBranch).Error
}

func (r *DealOfBranchPostgres) GetByBranchID(ctx context.Context, branchID uint64, page, size int) ([]domain.DealOfBranch, int64, error) {
	var deals []domain.DealOfBranch
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.DealOfBranch{}).Where("branch_id = ?", branchID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&deals).Error; err != nil {
		return nil, 0, err
	}

	return deals, total, nil
}

func (r *DealOfBranchPostgres) GetByDealID(ctx context.Context, dealID uint64) (*domain.DealOfBranch, error) {
	var deal domain.DealOfBranch
	if err := r.db.WithContext(ctx).Where("deal_id = ?", dealID).First(&deal).Error; err != nil {
		return nil, err
	}
	return &deal, nil
}

func (r *DealOfBranchPostgres) DeleteByDealID(ctx context.Context, dealID uint64) error {
	return r.db.WithContext(ctx).Where("deal_id = ?", dealID).Delete(&domain.DealOfBranch{}).Error
}
