package postgres

import (
	"bdspro/internal/domain"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.ProductAssetRepo
type ProductAssetPostgres struct {
	db *gorm.DB
}

func NewProductAssetPostgres(db *gorm.DB) *ProductAssetPostgres {
	return &ProductAssetPostgres{db: db}
}

func (r *ProductAssetPostgres) Link(ctx context.Context, productID, assetID uint64) error {
	productAsset := &domain.ProductAsset{
		ProductID: productID,
		AssetID:   assetID,
	}
	return GetDB(ctx, r.db).Create(productAsset).Error
}

func (r *ProductAssetPostgres) Unlink(ctx context.Context, productID, assetID uint64) error {
	return GetDB(ctx, r.db).
		Where("product_id = ? AND asset_id = ? AND deleted_at IS NULL", productID, assetID).
		Delete(&domain.ProductAsset{}).Error
}

func (r *ProductAssetPostgres) GetByProductID(ctx context.Context, productID uint64) ([]*domain.ProductAsset, error) {
	var productAssets []*domain.ProductAsset
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND deleted_at IS NULL", productID).
		Find(&productAssets).Error
	if err != nil {
		return nil, err
	}
	return productAssets, nil
}

func (r *ProductAssetPostgres) GetByAssetID(ctx context.Context, assetID uint64) ([]*domain.ProductAsset, error) {
	var productAssets []*domain.ProductAsset
	err := r.db.WithContext(ctx).
		Where("asset_id = ? AND deleted_at IS NULL", assetID).
		Find(&productAssets).Error
	if err != nil {
		return nil, err
	}
	return productAssets, nil
}

func (r *ProductAssetPostgres) CheckExists(ctx context.Context, productID, assetID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.ProductAsset{}).
		Where("product_id = ? AND asset_id = ? AND deleted_at IS NULL", productID, assetID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ProductAssetPostgres) GetAssetIDsByProductID(ctx context.Context, productID uint64) ([]uint64, error) {
	var assetIDs []uint64
	err := r.db.WithContext(ctx).
		Model(&domain.ProductAsset{}).
		Where("product_id = ? AND deleted_at IS NULL", productID).
		Pluck("asset_id", &assetIDs).Error
	if err != nil {
		return nil, err
	}
	return assetIDs, nil
}

func (r *ProductAssetPostgres) GetProductIDsByAssetID(ctx context.Context, assetID uint64) ([]uint64, error) {
	var productIDs []uint64
	err := r.db.WithContext(ctx).
		Model(&domain.ProductAsset{}).
		Where("asset_id = ? AND deleted_at IS NULL", assetID).
		Pluck("product_id", &productIDs).Error
	if err != nil {
		return nil, err
	}
	return productIDs, nil
}

func (r *ProductAssetPostgres) CountAssetsByProductID(ctx context.Context, productID uint64) (int64, error) {
	var count int64
	err := GetDB(ctx, r.db).
		Model(&domain.ProductAsset{}).
		Where("product_id = ? AND deleted_at IS NULL", productID).
		Count(&count).Error
	return count, err
}
