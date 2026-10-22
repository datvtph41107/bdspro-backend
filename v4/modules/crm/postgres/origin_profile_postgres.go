package postgres

import (
	"common/case/crud3"
	"context"
	"mdn4/modules/crm/enums"
	repo "mdn4/modules/crm/repository"
	entity "mdn4/modules/user/domain"

	"gorm.io/gorm"
)

type PostgreOriginProfileRepo struct {
	crud3.BaseRepo[entity.OriginProfile]
}

func NewPostgreOriginProfileRepo(db *gorm.DB) repo.OriginProfileRepo {
	return &PostgreOriginProfileRepo{
		BaseRepo: crud3.BaseRepo[entity.OriginProfile]{DB: db},
	}
}

func (r *PostgreOriginProfileRepo) GetByOriginID(ctx context.Context, originID uint64) (*entity.OriginProfile, error) {
	var profile entity.OriginProfile
	err := r.DB.WithContext(ctx).
		Where("origin_id = ? AND deleted_at IS NULL", originID).
		First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *PostgreOriginProfileRepo) GetByIDs(ctx context.Context, originIDs []uint64) ([]*entity.OriginProfile, error) {
	var profiles []*entity.OriginProfile
	err := r.DB.WithContext(ctx).
		Where("origin_id IN (?) AND deleted_at IS NULL", originIDs).
		Find(&profiles).Error
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

func (r *PostgreOriginProfileRepo) GetByProfileID(ctx context.Context, profileID uint64) (*entity.OriginProfile, error) {
	var profile entity.OriginProfile
	err := r.DB.WithContext(ctx).
		Where("owner_id = ? AND owner_of = ? AND deleted_at IS NULL", profileID, uint32(enums.EOwnerOfMember)).
		First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}
