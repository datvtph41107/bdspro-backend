package postgre

import (
	"context"
	"time"

	"auth/internal/domain"
	_db "common/db"

	"gorm.io/gorm"
)

// PINPostgres implementation của PIN repository
// @bind: auth/internal/interface/repo.PINRepository
type PINPostgres struct {
	db *gorm.DB
}

// NewPINRepository tạo mới PINRepository
func NewPINRepository(db *gorm.DB) *PINPostgres {
	return &PINPostgres{db: db}
}

// GetByAuthID lấy PIN theo authID
func (r *PINPostgres) GetByAuthID(ctx context.Context, authID uint64) (*domain.UserPINEntity, error) {
	var pin domain.UserPINEntity
	err := _db.DB.Where("auth_id = ?", authID).First(&pin).Error
	if err != nil {
		return nil, err
	}
	return &pin, nil
}

// Create tạo mới PIN
func (r *PINPostgres) Create(ctx context.Context, pin *domain.UserPINEntity) error {
	return _db.DB.Create(pin).Error
}

// Update cập nhật PIN, nếu chưa có thì tạo mới
func (r *PINPostgres) Update(ctx context.Context, pin *domain.UserPINEntity) error {
	// Kiểm tra xem PIN đã tồn tại chưa
	var existingPIN domain.UserPINEntity
	err := _db.DB.Where("auth_id = ?", pin.AuthID).First(&existingPIN).Error

	if err != nil {
		// Nếu không tìm thấy, tạo mới
		if err == gorm.ErrRecordNotFound {
			return _db.DB.Create(pin).Error
		}
		return err
	}

	// Nếu đã tồn tại, cập nhật
	return r.db.Model(&domain.UserPINEntity{}).
		Where("auth_id = ?", pin.AuthID).
		Updates(map[string]interface{}{
			"pin":            pin.PIN,
			"pin_check_time": pin.PINCheckTime,
			"pin_date":       pin.PINDate,
			"is_active":      pin.IsActive,
			"locked_until":   pin.LockedUntil,
		}).
		Error
}

// Delete xóa PIN
func (r *PINPostgres) Delete(ctx context.Context, authID uint64) error {
	return _db.DB.Where("auth_id = ?", authID).Delete(&domain.UserPINEntity{}).Error
}

// ResetCheckCounter reset số lần nhập sai PIN
func (r *PINPostgres) ResetCheckCounter(ctx context.Context, authID uint64) error {
	return _db.DB.Model(&domain.UserPINEntity{}).
		Where("auth_id = ?", authID).
		Updates(map[string]interface{}{
			"pin_check_time": 0,
			"locked_until":   nil,
		}).Error
}

// LockPIN khóa PIN do nhập sai nhiều lần
func (r *PINPostgres) LockPIN(ctx context.Context, authID uint64, lockedUntil time.Time) error {
	return _db.DB.Model(&domain.UserPINEntity{}).
		Where("auth_id = ?", authID).
		Updates(map[string]interface{}{
			"is_active":    false,
			"locked_until": lockedUntil,
		}).Error
}

// UnlockPIN mở khóa PIN
func (r *PINPostgres) UnlockPIN(ctx context.Context, authID uint64) error {
	return _db.DB.Model(&domain.UserPINEntity{}).
		Where("auth_id = ?", authID).
		Updates(map[string]interface{}{
			"is_active":      true,
			"locked_until":   nil,
			"pin_check_time": 0,
		}).Error
}

// CheckPINExists kiểm tra PIN có tồn tại không
func (r *PINPostgres) CheckPINExists(ctx context.Context, authID uint64) (bool, error) {
	var count int64
	err := _db.DB.Model(&domain.UserPINEntity{}).
		Where("auth_id = ? AND is_active = ?", authID, true).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// IncrementCheckTime tăng số lần nhập sai PIN
func (r *PINPostgres) IncrementCheckTime(ctx context.Context, authID uint64) error {
	return _db.DB.Model(&domain.UserPINEntity{}).
		Where("auth_id = ?", authID).
		Update("pin_check_time", gorm.Expr("pin_check_time + ?", 1)).
		Error
}
