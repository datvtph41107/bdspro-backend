package postgre_repo

import (
	"bdspro/v2/internal/domain"
	"common/case/crud3"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// @bind: bdspro/v2/repo.AssetRepo
type GormAssetRepo struct {
	crud3.BaseRepo[domain.Asset]
}

func NewGormAssetRepo(db *gorm.DB) *GormAssetRepo {
	return &GormAssetRepo{
		BaseRepo: crud3.BaseRepo[domain.Asset]{DB: db},
	}
}

// func (r *GormAssetRepo) Create(ctx context.Context, asset *domain.Asset) (*domain.Asset, error) {
// 	if err := r.DB.WithContext(ctx).Create(asset).Error; err != nil {
// 		return nil, err
// 	}
// 	return asset, nil
// }

func (r *GormAssetRepo) GetAssets(ctx context.Context, profileId uint64, limit, offset int) ([]domain.Asset, error) {
	var assets []domain.Asset

	err := r.DB.WithContext(ctx).
		Where("created_by = ? and deleted_at is null", profileId). // Lọc theo profileId
		Limit(limit).                                              // Phân trang: giới hạn số lượng bản ghi
		Offset(offset).                                            // Phân trang: bỏ qua số lượng bản ghi
		Find(&assets).Error

	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (r *GormAssetRepo) FilterAssets(ctx context.Context, filter map[string]interface{}) ([]domain.Asset, error) {
	var assets []domain.Asset
	query := r.DB.WithContext(ctx)

	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}

func (r *GormAssetRepo) BulkAction(ctx context.Context, action string, ids []uint64) error {
	switch action {
	case "delete":
		if err := r.DB.WithContext(ctx).Where("id IN ? and deleted_at is null", ids).Update("deleted_at", time.Now()).Error; err != nil {
			return err
		}
	case "archive":
		if err := r.DB.WithContext(ctx).Model(&domain.Asset{}).Where("id IN ? and deleted_at is null", ids).Update("doc_status", 3).Error; err != nil {
			return err
		}
	default:
		return errors.New("unsupported action: " + action)
	}
	return nil
}

func (r *GormAssetRepo) Update(ctx context.Context, id uint64, asset *domain.Asset) error {
	if err := r.DB.WithContext(ctx).Model(&domain.Asset{}).Where("id = ? and deleted_at is null", id).Updates(asset).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormAssetRepo) UpdateArchived(ctx context.Context, id uint64, archived bool) error {
	if err := r.DB.WithContext(ctx).Model(&domain.Asset{}).Where("id = ? and deleted_at is null", id).Update("archived", archived).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormAssetRepo) Delete(ctx context.Context, id uint64) error {
	if err := r.DB.WithContext(ctx).Where("id = ? and deleted_at is null", id).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormAssetRepo) GetAssetDetails(ctx context.Context, id uint64) (*domain.Asset, error) {
	var asset domain.Asset
	if err := r.DB.WithContext(ctx).Where("id = ? and deleted_at is null", id).First(&asset).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *GormAssetRepo) ShareAsset(ctx context.Context, id uint64, userID uint64) error {
	// Example: Add userID to a shared_users column (assumes JSONB array in PostgreSQL)
	if err := r.DB.WithContext(ctx).Exec(
		"UPDATE assets SET shared_users = array_append(shared_users, ?) WHERE id = ? and deleted_at is null", userID, id,
	).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormAssetRepo) SplitAsset(ctx context.Context, id uint64, splitData map[string]interface{}) ([]domain.Asset, error) {
	var asset domain.Asset
	if err := r.DB.WithContext(ctx).Where("id = ? and deleted_at is null", id).First(&asset).Error; err != nil {
		return nil, err
	}

	// Example logic for splitting
	part1 := domain.Asset{
		Name:        asset.Name + " Part 1",
		Description: asset.Description,
		Area:        asset.Area / 2,
		WardId:      asset.WardId,
		DistrictId:  asset.DistrictId,
		ProvinceId:  asset.ProvinceId,
		Address:     asset.Address,
		Price:       asset.Price / 2,
		DocStatus:   asset.DocStatus,
	}
	part2 := domain.Asset{
		Name:        asset.Name + " Part 2",
		Description: asset.Description,
		Area:        asset.Area / 2,
		WardId:      asset.WardId,
		DistrictId:  asset.DistrictId,
		ProvinceId:  asset.ProvinceId,
		Address:     asset.Address,
		Price:       asset.Price / 2,
		DocStatus:   asset.DocStatus,
	}

	if err := r.DB.WithContext(ctx).Create(&part1).Error; err != nil {
		return nil, err
	}
	if err := r.DB.WithContext(ctx).Create(&part2).Error; err != nil {
		return nil, err
	}

	// Delete the original asset
	if err := r.DB.WithContext(ctx).Where("id = ? and deleted_at is null", id).Delete(&domain.Asset{}).Error; err != nil {
		return nil, err
	}

	return []domain.Asset{part1, part2}, nil
}

func (r *GormAssetRepo) MergeAssets(ctx context.Context, ids []uint64) (*domain.Asset, error) {
	var assets []domain.Asset
	if err := r.DB.WithContext(ctx).Where("id IN ? and deleted_at is null", ids).Find(&assets).Error; err != nil {
		return nil, err
	}

	if len(assets) < 2 {
		return nil, errors.New("at least two assets are required to merge")
	}

	merged := domain.Asset{
		Name:        "Merged Asset",
		Description: "Merged from multiple assets",
		Area:        0,
		Price:       0,
		WardId:      assets[0].WardId,
		DistrictId:  assets[0].DistrictId,
		ProvinceId:  assets[0].ProvinceId,
		Address:     assets[0].Address,
		DocStatus:   assets[0].DocStatus,
	}

	for _, asset := range assets {
		merged.Area += asset.Area
		merged.Price += asset.Price
	}

	if err := r.DB.WithContext(ctx).Create(&merged).Error; err != nil {
		return nil, err
	}

	// Delete the original assets
	if err := r.DB.WithContext(ctx).Where("id IN ? and deleted_at is null", ids).Delete(&domain.Asset{}).Error; err != nil {
		return nil, err
	}

	return &merged, nil
}

func (r *GormAssetRepo) GetProductByID(ctx context.Context, productID uint64) (*domain.Product, error) {
	var product domain.Product
	if err := r.DB.WithContext(ctx).Where("id = ? and deleted_at is null", productID).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *GormAssetRepo) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	if err := r.DB.WithContext(ctx).Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}
