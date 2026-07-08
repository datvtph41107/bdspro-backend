package usecase

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/repo"
	"fmt"
)

type AdvertisingSettingsUsecase interface {
	// Advertising settings management - D.10.2.8
	CreateAdvertisingSettings(ctx context.Context, createDTO *dto.AdvertisingSettingsCreateDTO) (*dto.AdvertisingSettingsDTO, error)
	UpdateAdvertisingSettings(ctx context.Context, id uint64, updateDTO *dto.AdvertisingSettingsUpdateDTO) (*dto.AdvertisingSettingsDTO, error)
	GetAdvertisingSettings(ctx context.Context, userID uint64) (*dto.AdvertisingSettingsDTO, error)
	DeleteAdvertisingSettings(ctx context.Context, id uint64) error

	// Search and list
	SearchAdvertisingSettings(ctx context.Context, searchDTO dto.AdvertisingSettingsSearchDTO) ([]*dto.AdvertisingSettingsDTO, int64, error)
}

type advertisingSettingsUsecase struct {
	settingsRepo repo.AdvertisingSettingsRepo
	packageRepo  repo.PackageRepo
	campaignRepo repo.CampaignRepo
}

func NewAdvertisingSettingsUsecase(settingsRepo repo.AdvertisingSettingsRepo, packageRepo repo.PackageRepo, campaignRepo repo.CampaignRepo) AdvertisingSettingsUsecase {
	return &advertisingSettingsUsecase{
		settingsRepo: settingsRepo,
		packageRepo:  packageRepo,
		campaignRepo: campaignRepo,
	}
}

// CreateAdvertisingSettings creates advertising settings
func (u *advertisingSettingsUsecase) CreateAdvertisingSettings(ctx context.Context, createDTO *dto.AdvertisingSettingsCreateDTO) (*dto.AdvertisingSettingsDTO, error) {
	// Validate settings
	validation, err := u.ValidateSettings(ctx, createDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to validate settings: %w", err)
	}
	if !validation.IsValid {
		return nil, fmt.Errorf("invalid settings: %v", validation.Errors)
	}

	// Check if settings already exist for user
	existingSettings, err := u.settingsRepo.GetAdvertisingSettingsByUser(ctx, createDTO.UserID)
	if err == nil && existingSettings != nil {
		return nil, fmt.Errorf("advertising settings already exist for this user")
	}

	// Convert DTO to model
	settings := &domain.AdvertisingSettings{
		OrganizationID:         createDTO.OrganizationID,
		UserID:                 createDTO.UserID,
		DefaultCampaignType:    createDTO.DefaultCampaignType,
		DefaultPackageID:       createDTO.DefaultPackageID,
		DefaultBudget:          createDTO.DefaultBudget,
		DefaultDuration:        createDTO.DefaultDuration,
		AutoSuggestEnabled:     createDTO.AutoSuggestEnabled,
		AutoRunEnabled:         createDTO.AutoRunEnabled,
		AutoRunBudget:          createDTO.AutoRunBudget,
		AutoRunPackageID:       createDTO.AutoRunPackageID,
		MonthlyBudgetLimit:     createDTO.MonthlyBudgetLimit,
		AlertOnBudgetExceed:    createDTO.AlertOnBudgetExceed,
		AlertEmail:             createDTO.AlertEmail,
		AlertPhone:             createDTO.AlertPhone,
		ApplyToAllOrganization: createDTO.ApplyToAllOrganization,
		CreatedBy:              uint64(123), // TODO: Get from context
		UpdatedBy:              uint64(123), // TODO: Get from context
	}

	// Convert product type settings
	if createDTO.ProductTypeSettings != nil {
		settings.ProductTypeSettings = make(domain.ProductTypeSettings)
		for key, value := range createDTO.ProductTypeSettings {
			settings.ProductTypeSettings[key] = domain.ProductTypeSetting{
				ProductType:      value.ProductType,
				DefaultBudget:    value.DefaultBudget,
				DefaultDuration:  value.DefaultDuration,
				DefaultPackageID: value.DefaultPackageID,
				AutoRunEnabled:   value.AutoRunEnabled,
				AutoRunBudget:    value.AutoRunBudget,
				Priority:         value.Priority,
			}
		}
	}

	// Create settings
	createdSettings, err := u.settingsRepo.CreateAdvertisingSettings(ctx, settings)
	if err != nil {
		return nil, fmt.Errorf("failed to create advertising settings: %w", err)
	}

	return u.toAdvertisingSettingsDTO(createdSettings), nil
}

// UpdateAdvertisingSettings updates advertising settings
func (u *advertisingSettingsUsecase) UpdateAdvertisingSettings(ctx context.Context, id uint64, updateDTO *dto.AdvertisingSettingsUpdateDTO) (*dto.AdvertisingSettingsDTO, error) {
	// Get existing settings
	existingSettings, err := u.settingsRepo.GetAdvertisingSettingsByUser(ctx, uint64(123)) // TODO: Get userID from context
	if err != nil {
		return nil, fmt.Errorf("advertising settings not found: %w", err)
	}

	// Update fields
	if updateDTO.DefaultCampaignType != nil {
		existingSettings.DefaultCampaignType = *updateDTO.DefaultCampaignType
	}
	if updateDTO.DefaultPackageID != nil {
		existingSettings.DefaultPackageID = updateDTO.DefaultPackageID
	}
	if updateDTO.DefaultBudget != nil {
		existingSettings.DefaultBudget = *updateDTO.DefaultBudget
	}
	if updateDTO.DefaultDuration != nil {
		existingSettings.DefaultDuration = *updateDTO.DefaultDuration
	}
	if updateDTO.AutoSuggestEnabled != nil {
		existingSettings.AutoSuggestEnabled = *updateDTO.AutoSuggestEnabled
	}
	if updateDTO.AutoRunEnabled != nil {
		existingSettings.AutoRunEnabled = *updateDTO.AutoRunEnabled
	}
	if updateDTO.AutoRunBudget != nil {
		existingSettings.AutoRunBudget = *updateDTO.AutoRunBudget
	}
	if updateDTO.AutoRunPackageID != nil {
		existingSettings.AutoRunPackageID = updateDTO.AutoRunPackageID
	}
	if updateDTO.MonthlyBudgetLimit != nil {
		existingSettings.MonthlyBudgetLimit = *updateDTO.MonthlyBudgetLimit
	}
	if updateDTO.AlertOnBudgetExceed != nil {
		existingSettings.AlertOnBudgetExceed = *updateDTO.AlertOnBudgetExceed
	}
	if updateDTO.AlertEmail != nil {
		existingSettings.AlertEmail = *updateDTO.AlertEmail
	}
	if updateDTO.AlertPhone != nil {
		existingSettings.AlertPhone = *updateDTO.AlertPhone
	}
	if updateDTO.ApplyToAllOrganization != nil {
		existingSettings.ApplyToAllOrganization = *updateDTO.ApplyToAllOrganization
	}

	// Update product type settings
	if updateDTO.ProductTypeSettings != nil {
		existingSettings.ProductTypeSettings = make(domain.ProductTypeSettings)
		for key, value := range *updateDTO.ProductTypeSettings {
			existingSettings.ProductTypeSettings[key] = domain.ProductTypeSetting{
				ProductType:      value.ProductType,
				DefaultBudget:    value.DefaultBudget,
				DefaultDuration:  value.DefaultDuration,
				DefaultPackageID: value.DefaultPackageID,
				AutoRunEnabled:   value.AutoRunEnabled,
				AutoRunBudget:    value.AutoRunBudget,
				Priority:         value.Priority,
			}
		}
	}

	existingSettings.UpdatedBy = uint64(123) // TODO: Get from context

	// Update settings
	updatedSettings, err := u.settingsRepo.UpdateAdvertisingSettings(ctx, existingSettings.ID, existingSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to update advertising settings: %w", err)
	}

	return u.toAdvertisingSettingsDTO(updatedSettings), nil
}

// GetAdvertisingSettings gets advertising settings
func (u *advertisingSettingsUsecase) GetAdvertisingSettings(ctx context.Context, userID uint64) (*dto.AdvertisingSettingsDTO, error) {
	settings, err := u.settingsRepo.GetAdvertisingSettingsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("advertising settings not found: %w", err)
	}

	return u.toAdvertisingSettingsDTO(settings), nil
}

// DeleteAdvertisingSettings deletes advertising settings
func (u *advertisingSettingsUsecase) DeleteAdvertisingSettings(ctx context.Context, id uint64) error {
	return u.settingsRepo.DeleteAdvertisingSettings(ctx, id)
}

// SearchAdvertisingSettings searches advertising settings
func (u *advertisingSettingsUsecase) SearchAdvertisingSettings(ctx context.Context, searchDTO dto.AdvertisingSettingsSearchDTO) ([]*dto.AdvertisingSettingsDTO, int64, error) {
	settings, total, err := u.settingsRepo.SearchAdvertisingSettings(ctx, searchDTO)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search advertising settings: %w", err)
	}

	settingsDTOs := make([]*dto.AdvertisingSettingsDTO, len(settings))
	for i, setting := range settings {
		settingsDTOs[i] = u.toAdvertisingSettingsDTO(setting)
	}

	return settingsDTOs, total, nil
}

// ValidateSettings validates advertising settings
func (u *advertisingSettingsUsecase) ValidateSettings(ctx context.Context, settingsDTO *dto.AdvertisingSettingsCreateDTO) (*dto.SettingsValidationDTO, error) {
	// Convert DTO to model for validation
	settings := &domain.AdvertisingSettings{
		DefaultCampaignType: settingsDTO.DefaultCampaignType,
		DefaultBudget:       settingsDTO.DefaultBudget,
		DefaultDuration:     settingsDTO.DefaultDuration,
		AutoRunBudget:       settingsDTO.AutoRunBudget,
		MonthlyBudgetLimit:  settingsDTO.MonthlyBudgetLimit,
	}

	return u.settingsRepo.ValidateSettings(ctx, settings)
}

// Helper methods

func (u *advertisingSettingsUsecase) toAdvertisingSettingsDTO(settings *domain.AdvertisingSettings) *dto.AdvertisingSettingsDTO {
	// Convert product type settings
	productTypeSettings := make(map[string]dto.ProductTypeSettingDTO)
	for key, value := range settings.ProductTypeSettings {
		productTypeSettings[key] = dto.ProductTypeSettingDTO{
			ProductType:      value.ProductType,
			DefaultBudget:    value.DefaultBudget,
			DefaultDuration:  value.DefaultDuration,
			DefaultPackageID: value.DefaultPackageID,
			AutoRunEnabled:   value.AutoRunEnabled,
			AutoRunBudget:    value.AutoRunBudget,
			Priority:         value.Priority,
		}
	}

	return &dto.AdvertisingSettingsDTO{
		ID:                     settings.ID,
		OrganizationID:         settings.OrganizationID,
		UserID:                 settings.UserID,
		DefaultCampaignType:    settings.DefaultCampaignType,
		DefaultPackageID:       settings.DefaultPackageID,
		DefaultPackageName:     settings.DefaultPackageName,
		DefaultBudget:          settings.DefaultBudget,
		DefaultDuration:        settings.DefaultDuration,
		AutoSuggestEnabled:     settings.AutoSuggestEnabled,
		AutoRunEnabled:         settings.AutoRunEnabled,
		AutoRunBudget:          settings.AutoRunBudget,
		AutoRunPackageID:       settings.AutoRunPackageID,
		AutoRunPackageName:     settings.AutoRunPackageName,
		MonthlyBudgetLimit:     settings.MonthlyBudgetLimit,
		AlertOnBudgetExceed:    settings.AlertOnBudgetExceed,
		AlertEmail:             settings.AlertEmail,
		AlertPhone:             settings.AlertPhone,
		ProductTypeSettings:    productTypeSettings,
		ApplyToAllOrganization: settings.ApplyToAllOrganization,
		IsActive:               settings.IsActive,
		CreatedBy:              settings.CreatedBy,
		UpdatedBy:              settings.UpdatedBy,
		CreatedAt:              settings.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:              settings.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}