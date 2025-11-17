package postgre

import (
	"auth/internal/domain"
	_db "common/db"
	"context"

	"gorm.io/gorm"
)

// DevicePostgresRepo triển khai DeviceRepository với PostgreSQL
// @bind: auth/internal/interface/repo.DeviceRepository
type DevicePostgresRepo struct {
}

// NewDevicePostgresRepo tạo mới repository
func NewDevicePostgresRepo() *DevicePostgresRepo {
	return &DevicePostgresRepo{}
}

// FindByDeviceID tìm thiết bị theo device_id
func (r *DevicePostgresRepo) FindByDeviceID(ctx context.Context, deviceID string) (*domain.DeviceEntity, error) {
	var device domain.DeviceEntity
	err := _db.DB.WithContext(ctx).Where("device_id = ? AND deleted_at IS NULL", deviceID).First(&device).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &device, nil
}

// Create tạo mới thông tin thiết bị
func (r *DevicePostgresRepo) Create(ctx context.Context, device *domain.DeviceEntity) (*domain.DeviceEntity, error) {
	if err := _db.DB.WithContext(ctx).Create(device).Error; err != nil {
		return nil, err
	}
	return device, nil
}

// Update cập nhật thông tin thiết bị
func (r *DevicePostgresRepo) Update(ctx context.Context, device *domain.DeviceEntity) (*domain.DeviceEntity, error) {
	if err := _db.DB.WithContext(ctx).Save(device).Error; err != nil {
		return nil, err
	}
	return device, nil
}

// ListByAuthID lấy danh sách thiết bị theo auth_id
func (r *DevicePostgresRepo) ListByAuthID(ctx context.Context, authID uint64) ([]*domain.DeviceEntity, error) {
	var devices []*domain.DeviceEntity
	if err := _db.DB.WithContext(ctx).
		Where("auth_id = ? AND deleted_at IS NULL", authID).
		Order("updated_at DESC").
		Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}
