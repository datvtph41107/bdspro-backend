package u_asset

import (
	"bdspro/enums"
	"bdspro/models"
	"common/case/crud"
	"context"

	"gorm.io/gorm"
)

type AssetRepo struct {
	crud.BaseRepo[models.Asset]
}

func NewAssetRepo(db *gorm.DB) *AssetRepo {
	return &AssetRepo{
		BaseRepo: crud.BaseRepo[models.Asset]{DB: db},
	}
}

func (r *AssetRepo) GetAssetByID(id uint64) (*models.Asset, error) {
	var asset models.Asset
	if err := r.DB.First(&asset, id).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *AssetRepo) UpdateStatus(ctx context.Context, assetID uint64, status enums.AssetStatus) error {
	return r.DB.WithContext(ctx).Model(&models.Asset{}).
		Where("id = ? and deleted_at is null", assetID).
		Update("status", status).Error
}
