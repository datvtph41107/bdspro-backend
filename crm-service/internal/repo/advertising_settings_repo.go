package repo

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
)

type AdvertisingSettingsRepo interface {
	// Advertising settings management
	CreateAdvertisingSettings(ctx context.Context, settings *domain.AdvertisingSettings) (*domain.AdvertisingSettings, error)
	UpdateAdvertisingSettings(ctx context.Context, id uint64, settings *domain.AdvertisingSettings) (*domain.AdvertisingSettings, error)
	GetAdvertisingSettingsByUser(ctx context.Context, userID uint64) (*domain.AdvertisingSettings, error)
	DeleteAdvertisingSettings(ctx context.Context, id uint64) error

	// Search and list
	SearchAdvertisingSettings(ctx context.Context, searchDTO dto.AdvertisingSettingsSearchDTO) ([]*domain.AdvertisingSettings, int64, error)

	// Settings validation
	ValidateSettings(ctx context.Context, settings *domain.AdvertisingSettings) (*dto.SettingsValidationDTO, error)
}