package repository

import (
	"context"
	"errors"

	"mdn4/modules/hub/entity"

	"gorm.io/gorm"
)

// ErrApplinkNotFound is returned when a code lookup finds no row.
var ErrApplinkNotFound = errors.New("applink: not found")

// ApplinkRepository manages applink (deep-link code) persistence.
type ApplinkRepository interface {
	ExistsByCode(ctx context.Context, code uint64) (bool, error)
	Create(ctx context.Context, a *entity.Applink) error
	GetByRefIDAndAction(ctx context.Context, refID uint64, action int32) (*entity.Applink, error)
	GetByCode(ctx context.Context, code uint64) (*entity.Applink, error)
}

type applinkPostgres struct{ db *gorm.DB }

func NewApplinkPostgres(db *gorm.DB) ApplinkRepository {
	return &applinkPostgres{db: db}
}

func (r *applinkPostgres) ExistsByCode(ctx context.Context, code uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Applink{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}

func (r *applinkPostgres) Create(ctx context.Context, a *entity.Applink) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *applinkPostgres) GetByRefIDAndAction(ctx context.Context, refID uint64, action int32) (*entity.Applink, error) {
	var a entity.Applink
	err := r.db.WithContext(ctx).Where("ref_id = ? AND action = ?", refID, action).First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &a, err
}

func (r *applinkPostgres) GetByCode(ctx context.Context, code uint64) (*entity.Applink, error) {
	var a entity.Applink
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrApplinkNotFound
	}
	return &a, err
}
