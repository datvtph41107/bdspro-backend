package postgre_repo

import (
	"bdspro/v2/internal/domain"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/v2/repo.AssetHistoryRepo
type GormAssetHistoryRepo struct {
	DB *gorm.DB
}

func NewGormHistoryRepo(db *gorm.DB) *GormAssetHistoryRepo {
	return &GormAssetHistoryRepo{DB: db}
}

func (r *GormAssetHistoryRepo) Create(ctx context.Context, history *domain.AssetSplitMergeHistory) (*domain.AssetSplitMergeHistory, error) {
	if err := r.DB.WithContext(ctx).Create(history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

func (r *GormAssetHistoryRepo) GetByAssetID(ctx context.Context, assetID uint64) ([]domain.AssetSplitMergeHistory, error) {
	var histories []domain.AssetSplitMergeHistory
	if err := r.DB.WithContext(ctx).Where("asset_id = ?", assetID).Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}
