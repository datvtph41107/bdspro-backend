package postgres

import (
	"bdspro/internal/domain"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.AssetOrganizationRepo
type PostgreAssetOrganization struct {
	DB *gorm.DB
}

func NewPostgreAssetOrganization(db *gorm.DB) *PostgreAssetOrganization {
	return &PostgreAssetOrganization{
		DB: db,
	}
}

func (r *PostgreAssetOrganization) CreateOwner(ctx context.Context, assetID uint64, organizationID uint64) error {
	return r.DB.WithContext(ctx).Create(&domain.AssetOrganization{
		AssetID:        assetID,
		OrganizationID: organizationID,
	}).Error
}
