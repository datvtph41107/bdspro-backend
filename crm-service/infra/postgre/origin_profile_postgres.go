package postgre

import (
	"common/case/crud3"
	"context"
	"crm/internal/domain"
	"crm/internal/repo"

	"gorm.io/gorm"
)

type PostgreOriginProfileRepo struct {
	crud3.BaseRepo[domain.OriginProfile]
}

func NewPostgreOriginProfileRepo(db *gorm.DB) repo.OriginProfileRepo {
	return &PostgreOriginProfileRepo{
		BaseRepo: crud3.BaseRepo[domain.OriginProfile]{DB: db},
	}
}

func (r *PostgreOriginProfileRepo) GetByOriginID(ctx context.Context, originID uint64) (*domain.OriginProfile, error) {
	var profile domain.OriginProfile
	err := r.DB.WithContext(ctx).
		Where("origin_id = ? AND deleted_at IS NULL", originID).
		First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *PostgreOriginProfileRepo) GetByOriginIDs(ctx context.Context, originIDs []uint64) ([]*domain.OriginProfile, error) {
	var profiles []*domain.OriginProfile
	err := r.DB.WithContext(ctx).
		Where("origin_id IN (?) AND deleted_at IS NULL", originIDs).
		Find(&profiles).Error
	if err != nil {
		return nil, err
	}
	return profiles, nil
}
