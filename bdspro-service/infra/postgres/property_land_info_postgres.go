package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/repo"
	"context"

	"gorm.io/gorm"
)

type PostgrePropertyLandInfoRepo struct {
	DB *gorm.DB
}

func NewPostgrePropertyLandInfoRepo(db *gorm.DB) repo.PropertyLandInfoRepo {
	return &PostgrePropertyLandInfoRepo{
		DB: db,
	}
}

func (r *PostgrePropertyLandInfoRepo) Create(ctx context.Context, landInfo *domain.PropertyLandInfo) error {
	return GetDB(ctx, r.DB).Create(landInfo).Error
}

func (r *PostgrePropertyLandInfoRepo) Update(ctx context.Context, landInfo *domain.PropertyLandInfo) error {
	return GetDB(ctx, r.DB).
		Model(&domain.PropertyLandInfo{}).
		Where("id = ? AND deleted_at IS NULL", landInfo.ID).
		Updates(landInfo).Error
}

func (r *PostgrePropertyLandInfoRepo) GetByPropertyID(ctx context.Context, propertyID uint64) (*domain.PropertyLandInfo, error) {
	var landInfo domain.PropertyLandInfo
	err := GetDB(ctx, r.DB).
		Where("property_id = ? AND deleted_at IS NULL", propertyID).
		First(&landInfo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &landInfo, nil
}
