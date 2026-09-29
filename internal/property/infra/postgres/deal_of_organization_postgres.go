package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"

	"gorm.io/gorm"
)

type DealOfOrganizationPostgres struct {
	db *gorm.DB
}

func NewDealOfOrganizationPostgres(db *gorm.DB) repo.DealOfOrganizationRepo {
	return &DealOfOrganizationPostgres{db: db}
}

func (r *DealOfOrganizationPostgres) Create(ctx context.Context, dealOfOrganization *domain.DealOfOrganization) error {
	return r.db.WithContext(ctx).Create(dealOfOrganization).Error
}

func (r *DealOfOrganizationPostgres) GetByOrganizationID(ctx context.Context, organizationID uint64, page, size int) ([]domain.DealOfOrganization, int64, error) {
	var deals []domain.DealOfOrganization
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.DealOfOrganization{}).Where("organization_id = ?", organizationID)

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

func (r *DealOfOrganizationPostgres) GetByDealID(ctx context.Context, dealID uint64) (*domain.DealOfOrganization, error) {
	var deal domain.DealOfOrganization
	if err := r.db.WithContext(ctx).Where("deal_id = ?", dealID).First(&deal).Error; err != nil {
		return nil, err
	}
	return &deal, nil
}

func (r *DealOfOrganizationPostgres) DeleteByDealID(ctx context.Context, dealID uint64) error {
	return r.db.WithContext(ctx).Where("deal_id = ?", dealID).Delete(&domain.DealOfOrganization{}).Error
}
