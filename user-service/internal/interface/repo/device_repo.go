package repo

import (
	"context"
	"user/internal/domain/auth"
)

// DeviceRepository định nghĩa các hành vi thao tác thiết bị trong database
type DeviceRepository interface {
	FindByDeviceID(ctx context.Context, deviceID uint64) (*auth.DeviceEntity, error)
	FindByDeviceIDAndModelAndManufacturer(ctx context.Context, deviceID string, model string, manufacturer string) (*auth.DeviceEntity, error)
	Create(ctx context.Context, device *auth.DeviceEntity) (*auth.DeviceEntity, error)
	Update(ctx context.Context, device *auth.DeviceEntity) (*auth.DeviceEntity, error)
	ListByAuthID(ctx context.Context, authID uint64) ([]*auth.DeviceEntity, error)
	ListByProfileID(ctx context.Context, profileID uint64) ([]*auth.DeviceEntity, error)
}
