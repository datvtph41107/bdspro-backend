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
	GetAdvertisingSettingsByOrganization(ctx context.Context, organizationID uint64) (*model.AdvertisingSettings, error)
	DeleteAdvertisingSettings(ctx context.Context, id uint64) error
	
	// Search and list
	SearchAdvertisingSettings(ctx context.Context, searchDTO dto.AdvertisingSettingsSearchDTO) ([]*model.AdvertisingSettings, int64, error)
	
	// Campaign defaults
	GetCampaignDefaults(ctx context.Context, userID uint64, productType string) (*dto.CampaignDefaultsDTO, error)
	GetOrganizationDefaults(ctx context.Context, organizationID uint64, productType string) (*dto.CampaignDefaultsDTO, error)
	
	// Auto campaign suggestions
	CreateAutoCampaignSuggestion(ctx context.Context, suggestion *model.AutoCampaignSuggestion) (*model.AutoCampaignSuggestion, error)
	UpdateAutoCampaignSuggestion(ctx context.Context, id uint64, suggestion *model.AutoCampaignSuggestion) (*model.AutoCampaignSuggestion, error)
	GetAutoCampaignSuggestions(ctx context.Context, userID uint64, page, limit int) ([]*model.AutoCampaignSuggestion, int64, error)
	GetPendingAutoCampaignSuggestions(ctx context.Context, userID uint64) ([]*model.AutoCampaignSuggestion, error)
	AcceptAutoCampaignSuggestion(ctx context.Context, suggestionID uint64, campaignID uint64) error
	RejectAutoCampaignSuggestion(ctx context.Context, suggestionID uint64) error
	
	// Settings validation
	ValidateSettings(ctx context.Context, settings *model.AdvertisingSettings) (*dto.SettingsValidationDTO, error)
	CheckUserPermissions(ctx context.Context, userID uint64, organizationID uint64) (bool, string, error)
	
	// Auto suggestion generation
	GenerateAutoCampaignSuggestions(ctx context.Context, requestDTO *dto.AutoCampaignSuggestionRequestDTO) (*dto.AutoCampaignSuggestionResponseDTO, error)
	
	// Settings inheritance
	GetInheritedSettings(ctx context.Context, userID uint64, organizationID uint64) (*model.AdvertisingSettings, error)
	ApplySettingsToOrganization(ctx context.Context, settingsID uint64, organizationID uint64) error
} 