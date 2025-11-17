package repo

import (
	"auth/internal/domain"
	"context"
)

// DeviceRepository định nghĩa các hành vi thao tác thiết bị trong database
type DeviceRepository interface {
	FindByDeviceID(ctx context.Context, deviceID string) (*domain.DeviceEntity, error)
	Create(ctx context.Context, device *domain.DeviceEntity) (*domain.DeviceEntity, error)
	Update(ctx context.Context, device *domain.DeviceEntity) (*domain.DeviceEntity, error)
	ListByAuthID(ctx context.Context, authID uint64) ([]*domain.DeviceEntity, error)
}
