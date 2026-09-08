package postgres

import (
	"context"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"

	"gorm.io/gorm"
)

type DealOfGroupPostgres struct {
	db *gorm.DB
}

func NewDealOfGroupPostgres(db *gorm.DB) repository.DealOfGroupRepo {
	return &DealOfGroupPostgres{db: db}
}

func (r *DealOfGroupPostgres) Create(ctx context.Context, dealOfGroup *entity.DealOfGroup) error {
	return r.db.WithContext(ctx).Create(dealOfGroup).Error
}

func (r *DealOfGroupPostgres) GetByGroupID(ctx context.Context, groupID uint64, page, size int) ([]entity.DealOfGroup, int64, error) {
	var deals []entity.DealOfGroup
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.DealOfGroup{}).Where("group_id = ?", groupID)

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

func (r *DealOfGroupPostgres) GetByDealID(ctx context.Context, dealID uint64) (*entity.DealOfGroup, error) {
	var deal entity.DealOfGroup
	if err := r.db.WithContext(ctx).Where("deal_id = ?", dealID).First(&deal).Error; err != nil {
		return nil, err
	}
	return &deal, nil
}

func (r *DealOfGroupPostgres) DeleteByDealID(ctx context.Context, dealID uint64) error {
	return r.db.WithContext(ctx).Where("deal_id = ?", dealID).Delete(&entity.DealOfGroup{}).Error
}
