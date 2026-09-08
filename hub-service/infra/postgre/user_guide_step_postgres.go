package postgres

import (
	"context"

	_db "common/db"
	"hub/internal/domain"
	"hub/internal/repo"
)

type UserGuideStepRepo struct {
	db *_db.TransactionRepo
}

func NewUserGuideStepRepo(db *_db.TransactionRepo) repo.IUserGuideStepRepo {
	return &UserGuideStepRepo{db: db}
}

// GetByUserGuideID lấy tất cả steps của một user guide
func (r *UserGuideStepRepo) GetByUserGuideID(ctx context.Context, userGuideID uint64) ([]*domain.UserGuideStepEntity, error) {
	var steps []*domain.UserGuideStepEntity

	err := r.db.GetDB(ctx).
		Where("user_guide_id = ?", userGuideID).
		Where("deleted_at IS NULL").
		Order("step_order ASC").
		Find(&steps).Error

	if err != nil {
		return nil, err
	}

	return steps, nil
}

// CreateBatch tạo nhiều steps cùng lúc
func (r *UserGuideStepRepo) CreateBatch(ctx context.Context, steps []*domain.UserGuideStepEntity) error {
	if len(steps) == 0 {
		return nil
	}

	return r.db.GetDB(ctx).Create(&steps).Error
}

// DeleteByUserGuideID xóa tất cả steps của một user guide
func (r *UserGuideStepRepo) DeleteByUserGuideID(ctx context.Context, userGuideID uint64) error {
	return r.db.GetDB(ctx).
		Where("user_guide_id = ?", userGuideID).
		Delete(&domain.UserGuideStepEntity{}).Error
}
