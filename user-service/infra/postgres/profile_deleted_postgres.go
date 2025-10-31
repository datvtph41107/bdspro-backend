package postgres

import (
	"context"
	"user/internal/models"

	"gorm.io/gorm"
)

// @bind: user/internal/interface/repo.IProfileDeletedRepo
type ProfileDeletedPostgres struct {
	DB *gorm.DB
}

func NewProfileDeletedPostgres(db *gorm.DB) *ProfileDeletedPostgres {
	return &ProfileDeletedPostgres{
		DB: db,
	}
}

// Create lưu thông tin profile đã xóa vào bảng profile_deleted
func (r *ProfileDeletedPostgres) Create(ctx context.Context, profileDeleted *models.ProfileDeleted) error {
	return r.DB.WithContext(ctx).Create(profileDeleted).Error
}

// GetByProfileID lấy thông tin profile đã xóa theo profileID
func (r *ProfileDeletedPostgres) GetByProfileID(ctx context.Context, profileID uint64) (*models.ProfileDeleted, error) {
	var profileDeleted models.ProfileDeleted
	err := r.DB.WithContext(ctx).
		Where("profile_id = ?", profileID).
		First(&profileDeleted).Error
	if err != nil {
		return nil, err
	}
	return &profileDeleted, nil
}

// GetList lấy danh sách profile đã xóa với phân trang
func (r *ProfileDeletedPostgres) GetList(ctx context.Context, offset, limit int) ([]*models.ProfileDeleted, int64, error) {
	var profiles []*models.ProfileDeleted
	var total int64

	// Get total count
	err := r.DB.WithContext(ctx).
		Model(&models.ProfileDeleted{}).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err = r.DB.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("deleted_at DESC").
		Find(&profiles).Error
	if err != nil {
		return nil, 0, err
	}

	return profiles, total, nil
}
