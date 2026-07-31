package postgres

import (
	"context"
	"time"

	_db "common/db"
	"user/internal/domain"
)

// @bind: auth/internal/interface/repo.StatusRepository
type StatusPostgres struct {
	// TODO: Add database connection
}

// NewStatusRepository tạo mới StatusRepository
func NewStatusRepository() *StatusPostgres {
	return &StatusPostgres{}
}

// GetByID lấy status theo authID
func (r *StatusPostgres) GetByID(ctx context.Context, authID uint64) (*domain.UserStatusEntity, error) {
	var status domain.UserStatusEntity
	err := _db.DB.Where("auth_id = ?", authID).First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// CreateStatus tạo mới status
func (r *StatusPostgres) CreateStatus(ctx context.Context, status *domain.UserStatusEntity) error {
	return _db.DB.Create(status).Error
}

// UpdateStatus cập nhật status
func (r *StatusPostgres) UpdateStatus(ctx context.Context, status *domain.UserStatusEntity) error {
	return _db.DB.Save(status).Error
}

// DeleteStatus xóa status
func (r *StatusPostgres) DeleteStatus(ctx context.Context, authID uint64) error {
	return _db.DB.Where("auth_id = ?", authID).Delete(&domain.UserStatusEntity{}).Error
}

// GetStatusByUserID lấy status theo userID
func (r *StatusPostgres) GetStatusByUserID(ctx context.Context, userID uint64) (*domain.UserStatusEntity, error) {
	var status domain.UserStatusEntity
	err := _db.DB.Joins("JOIN auth_method_domains ON user_status_entities.auth_id = auth_method_domains.id").
		Where("auth_method_domains.user_id = ?", userID).
		First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// GetLockedUsers lấy danh sách user bị khóa
func (r *StatusPostgres) GetLockedUsers(ctx context.Context) ([]*domain.UserStatusEntity, error) {
	var statuses []*domain.UserStatusEntity
	err := _db.DB.Where("locked_until > ?", time.Now()).Find(&statuses).Error
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// GetInactiveUsers lấy danh sách user không hoạt động
func (r *StatusPostgres) GetInactiveUsers(ctx context.Context) ([]*domain.UserStatusEntity, error) {
	var statuses []*domain.UserStatusEntity
	err := _db.DB.Where("active = ?", false).Find(&statuses).Error
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// UnlockUser mở khóa user
func (r *StatusPostgres) UnlockUser(ctx context.Context, authID uint64) error {
	return _db.DB.Model(&domain.UserStatusEntity{}).
		Where("auth_id = ?", authID).
		Update("locked_until", time.Time{}).Error
}

// ActivateUser kích hoạt user
func (r *StatusPostgres) ActivateUser(ctx context.Context, authID uint64) error {
	return _db.DB.Model(&domain.UserStatusEntity{}).
		Where("auth_id = ?", authID).
		Updates(map[string]interface{}{
			"active":   true,
			"verified": true,
		}).Error
}

// DeactivateUser vô hiệu hóa user
func (r *StatusPostgres) DeactivateUser(ctx context.Context, authID uint64) error {
	return _db.DB.Model(&domain.UserStatusEntity{}).
		Where("auth_id = ?", authID).
		Update("active", false).Error
}

// SetLockUntil thiết lập thời gian khóa user
func (r *StatusPostgres) SetLockUntil(ctx context.Context, authID uint64, lockUntil time.Time) error {
	return _db.DB.Model(&domain.UserStatusEntity{}).
		Where("auth_id = ?", authID).
		Update("locked_until", lockUntil).Error
}
