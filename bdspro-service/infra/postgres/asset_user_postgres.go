package postgres

import (
	"bdspro/internal/domain"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.AssetUserRepo
type PostgreAssetUser struct {
	DB *gorm.DB
}

func NewPostgreAssetUser(db *gorm.DB) *PostgreAssetUser {
	return &PostgreAssetUser{
		DB: db,
	}
}

func (r *PostgreAssetUser) CreateOwner(ctx context.Context, assetID uint64, profileID uint64) error {
	return r.DB.WithContext(ctx).Create(&domain.AssetUser{
		AssetID:   assetID,
		ProfileID: profileID,
		IsOwner:   true,
	}).Error
}
