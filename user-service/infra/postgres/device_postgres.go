package postgres

import (
	"context"
	"gorm.io/gorm"
	"user/internal/domain/auth"
)

// DevicePostgresRepo triển khai DeviceRepository với PostgreSQL
// @bind: user/internal/interface/repo.DeviceRepository
type DevicePostgresRepo struct {
	db *gorm.DB
}

// NewDevicePostgresRepo tạo mới repository
func NewDevicePostgresRepo(db *gorm.DB) *DevicePostgresRepo {
	return &DevicePostgresRepo{db: db}
}

// FindByDeviceID tìm thiết bị theo device_id
func (r *DevicePostgresRepo) FindByDeviceID(ctx context.Context, deviceID uint64) (*auth.DeviceEntity, error) {
	var device auth.DeviceEntity
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", deviceID).First(&device).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &device, nil
}

func (r *DevicePostgresRepo) FindByDeviceIDAndModelAndManufacturer(ctx context.Context, deviceID string, model string, manufacturer string) (*auth.DeviceEntity, error) {
	var device auth.DeviceEntity
	err := r.db.WithContext(ctx).Where("device_id = ? AND model = ? AND manufacturer = ? AND deleted_at IS NULL", deviceID, model, manufacturer).First(&device).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
	}
	return &device, nil
}

// Create tạo mới thông tin thiết bị
func (r *DevicePostgresRepo) Create(ctx context.Context, device *auth.DeviceEntity) (*auth.DeviceEntity, error) {
	if err := r.db.WithContext(ctx).Create(device).Error; err != nil {
		return nil, err
	}
	return device, nil
}

// Update cập nhật thông tin thiết bị
func (r *DevicePostgresRepo) Update(ctx context.Context, device *auth.DeviceEntity) (*auth.DeviceEntity, error) {
	if err := r.db.WithContext(ctx).Save(device).Error; err != nil {
		return nil, err
	}
	return device, nil
}

// ListByAuthID lấy danh sách thiết bị theo auth_id
func (r *DevicePostgresRepo) ListByAuthID(ctx context.Context, authID uint64) ([]*auth.DeviceEntity, error) {
	var devices []*auth.DeviceEntity
	if err := r.db.WithContext(ctx).
		Where("auth_id = ? AND deleted_at IS NULL", authID).
		Order("updated_at DESC").
		Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

// ListByProfileID lấy danh sách thiết bị theo profile_id
func (r *DevicePostgresRepo) ListByProfileID(ctx context.Context, profileID uint64) ([]*auth.DeviceEntity, error) {
	var devices []*auth.DeviceEntity
	if err := r.db.WithContext(ctx).
		Where("profile_id = ? AND deleted_at IS NULL AND push_token IS NOT NULL AND push_token != ''", profileID).
		Order("updated_at DESC").
		Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}
