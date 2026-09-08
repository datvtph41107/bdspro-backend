package postgres

import (
	"context"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"

	"gorm.io/gorm"
)

type DealOfOrganizationPostgres struct {
	db *gorm.DB
}

func NewDealOfOrganizationPostgres(db *gorm.DB) repository.DealOfOrganizationRepo {
	return &DealOfOrganizationPostgres{db: db}
}

func (r *DealOfOrganizationPostgres) Create(ctx context.Context, dealOfOrganization *entity.DealOfOrganization) error {
	return r.db.WithContext(ctx).Create(dealOfOrganization).Error
}

func (r *DealOfOrganizationPostgres) GetByOrganizationID(ctx context.Context, organizationID uint64, page, size int) ([]entity.DealOfOrganization, int64, error) {
	var deals []entity.DealOfOrganization
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.DealOfOrganization{}).Where("organization_id = ?", organizationID)

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

func (r *DealOfOrganizationPostgres) GetByDealID(ctx context.Context, dealID uint64) (*entity.DealOfOrganization, error) {
	var deal entity.DealOfOrganization
	if err := r.db.WithContext(ctx).Where("deal_id = ?", dealID).First(&deal).Error; err != nil {
		return nil, err
	}
	return &deal, nil
}

func (r *DealOfOrganizationPostgres) DeleteByDealID(ctx context.Context, dealID uint64) error {
	return r.db.WithContext(ctx).Where("deal_id = ?", dealID).Delete(&entity.DealOfOrganization{}).Error
}
