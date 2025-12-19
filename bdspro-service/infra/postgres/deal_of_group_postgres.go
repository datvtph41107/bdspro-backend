package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"

	"gorm.io/gorm"
)

type DealOfGroupPostgres struct {
	db *gorm.DB
}

func NewDealOfGroupPostgres(db *gorm.DB) repo.DealOfGroupRepo {
	return &DealOfGroupPostgres{db: db}
}

func (r *DealOfGroupPostgres) Create(ctx context.Context, dealOfGroup *domain.DealOfGroup) error {
	return r.db.WithContext(ctx).Create(dealOfGroup).Error
}

func (r *DealOfGroupPostgres) GetByGroupID(ctx context.Context, groupID uint64, page, size int) ([]domain.DealOfGroup, int64, error) {
	var deals []domain.DealOfGroup
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.DealOfGroup{}).Where("group_id = ?", groupID)

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

func (r *DealOfGroupPostgres) GetByDealID(ctx context.Context, dealID uint64) (*domain.DealOfGroup, error) {
	var deal domain.DealOfGroup
	if err := r.db.WithContext(ctx).Where("deal_id = ?", dealID).First(&deal).Error; err != nil {
		return nil, err
	}
	return &deal, nil
}

func (r *DealOfGroupPostgres) DeleteByDealID(ctx context.Context, dealID uint64) error {
	return r.db.WithContext(ctx).Where("deal_id = ?", dealID).Delete(&domain.DealOfGroup{}).Error
}
