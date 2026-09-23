package postgres

import (
	"context"
	"user/internal/domain/user"

	"gorm.io/gorm"
)

// @bind: user/internal/interface/repo.IDemoProductLogRepo
type DemoProductLogPostgres struct {
	DB *gorm.DB
}

func NewDemoProductLogPostgres(db *gorm.DB) *DemoProductLogPostgres {
	return &DemoProductLogPostgres{DB: db}
}

func (r *DemoProductLogPostgres) CreateBatch(ctx context.Context, logs []*user.DemoProductLog) error {
	if len(logs) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).Create(&logs).Error
}

func (r *DemoProductLogPostgres) ListProductIDsByProfileID(ctx context.Context, profileID uint64) ([]uint64, error) {
	var ids []uint64
	err := r.DB.WithContext(ctx).
		Model(&user.DemoProductLog{}).
		Where("profile_id = ?", profileID).
		Order("id ASC").
		Pluck("product_id", &ids).Error
	return ids, err
}

func (r *DemoProductLogPostgres) DeleteByProfileID(ctx context.Context, profileID uint64) error {
	return r.DB.WithContext(ctx).
		Where("profile_id = ?", profileID).
		Delete(&user.DemoProductLog{}).Error
}

func (r *DemoProductLogPostgres) CreateContactBatch(ctx context.Context, logs []*user.DemoContactLog) error {
	if len(logs) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).Create(&logs).Error
}

func (r *DemoProductLogPostgres) ListContactIDsByProfileID(ctx context.Context, profileID uint64) ([]uint64, error) {
	var ids []uint64
	err := r.DB.WithContext(ctx).
		Model(&user.DemoContactLog{}).
		Where("profile_id = ?", profileID).
		Order("id ASC").
		Pluck("contact_id", &ids).Error
	return ids, err
}

func (r *DemoProductLogPostgres) DeleteContactsByProfileID(ctx context.Context, profileID uint64) error {
	return r.DB.WithContext(ctx).
		Where("profile_id = ?", profileID).
		Delete(&user.DemoContactLog{}).Error
}
