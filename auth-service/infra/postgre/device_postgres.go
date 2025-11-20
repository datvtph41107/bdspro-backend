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
func (r *DevicePostgresRepo) FindByDeviceID(ctx context.Context, deviceID uint64) (*domain.DeviceEntity, error) {
	var device domain.DeviceEntity
	err := _db.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", deviceID).First(&device).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &device, nil
}

func (r *DevicePostgresRepo) FindByDeviceIDAndModelAndManufacturer(ctx context.Context, deviceID string, model string, manufacturer string) (*domain.DeviceEntity, error) {
	var device domain.DeviceEntity
	err := _db.DB.WithContext(ctx).Where("device_id = ? AND model = ? AND manufacturer = ? AND deleted_at IS NULL", deviceID, model, manufacturer).First(&device).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
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

// ListByProfileID lấy danh sách thiết bị theo profile_id
func (r *DevicePostgresRepo) ListByProfileID(ctx context.Context, profileID uint64) ([]*domain.DeviceEntity, error) {
	var devices []*domain.DeviceEntity
	if err := _db.DB.WithContext(ctx).
		Where("profile_id = ? AND deleted_at IS NULL AND push_token IS NOT NULL AND push_token != ''", profileID).
		Order("updated_at DESC").
		Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}
