package repository

import (
	"context"

	"marketing/internal/domain/model"
	"marketing/internal/dto"
)

type AdvertisingSettingsRepository interface {
	// Advertising settings management
	CreateAdvertisingSettings(ctx context.Context, settings *model.AdvertisingSettings) (*model.AdvertisingSettings, error)
	UpdateAdvertisingSettings(ctx context.Context, id uint64, settings *model.AdvertisingSettings) (*model.AdvertisingSettings, error)
	GetAdvertisingSettingsByUser(ctx context.Context, userID uint64) (*model.AdvertisingSettings, error)
	DeleteAdvertisingSettings(ctx context.Context, id uint64) error

	// Search and list
	SearchAdvertisingSettings(ctx context.Context, searchDTO dto.AdvertisingSettingsSearchDTO) ([]*model.AdvertisingSettings, int64, error)

	// Settings validation
	ValidateSettings(ctx context.Context, settings *model.AdvertisingSettings) (*dto.SettingsValidationDTO, error)
}
