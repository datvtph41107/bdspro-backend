package postgre

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/repo"
	"fmt"

	"gorm.io/gorm"
)

type AdvertisingSettingsPostgresRepository struct {
	db *gorm.DB
}

func NewAdvertisingSettingsPostgresRepository(db *gorm.DB) repo.AdvertisingSettingsRepo {
	return &AdvertisingSettingsPostgresRepository{
		db: db,
	}
}

// Advertising settings management

func (r *AdvertisingSettingsPostgresRepository) CreateAdvertisingSettings(ctx context.Context, settings *domain.AdvertisingSettings) (*domain.AdvertisingSettings, error) {
	if err := r.db.WithContext(ctx).Create(settings).Error; err != nil {
		return nil, fmt.Errorf("failed to create advertising settings: %w", err)
	}
	return settings, nil
}

func (r *AdvertisingSettingsPostgresRepository) UpdateAdvertisingSettings(ctx context.Context, id uint64, settings *domain.AdvertisingSettings) (*domain.AdvertisingSettings, error) {
	if err := r.db.WithContext(ctx).Model(&domain.AdvertisingSettings{}).Where("id = ?", id).Updates(settings).Error; err != nil {
		return nil, fmt.Errorf("failed to update advertising settings: %w", err)
	}
	return r.GetAdvertisingSettingsByUser(ctx, settings.UserID)
}

func (r *AdvertisingSettingsPostgresRepository) GetAdvertisingSettingsByUser(ctx context.Context, userID uint64) (*domain.AdvertisingSettings, error) {
	var settings domain.AdvertisingSettings
	if err := r.db.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true).First(&settings).Error; err != nil {
		return nil, fmt.Errorf("advertising settings not found: %w", err)
	}
	return &settings, nil
}

func (r *AdvertisingSettingsPostgresRepository) DeleteAdvertisingSettings(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).Delete(&domain.AdvertisingSettings{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete advertising settings: %w", err)
	}
	return nil
}

// Search and list

func (r *AdvertisingSettingsPostgresRepository) SearchAdvertisingSettings(ctx context.Context, searchDTO dto.AdvertisingSettingsSearchDTO) ([]*domain.AdvertisingSettings, int64, error) {
	var settings []*domain.AdvertisingSettings
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.AdvertisingSettings{})

	// Apply filters
	if searchDTO.OrganizationID != nil {
		query = query.Where("organization_id = ?", *searchDTO.OrganizationID)
	}
	if searchDTO.UserID != nil {
		query = query.Where("user_id = ?", *searchDTO.UserID)
	}
	if searchDTO.DefaultCampaignType != "" {
		query = query.Where("default_campaign_type = ?", searchDTO.DefaultCampaignType)
	}
	if searchDTO.AutoSuggestEnabled != nil {
		query = query.Where("auto_suggest_enabled = ?", *searchDTO.AutoSuggestEnabled)
	}
	if searchDTO.AutoRunEnabled != nil {
		query = query.Where("auto_run_enabled = ?", *searchDTO.AutoRunEnabled)
	}
	if searchDTO.IsActive != nil {
		query = query.Where("is_active = ?", *searchDTO.IsActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count advertising settings: %w", err)
	}

	// Apply pagination
	offset := (searchDTO.Page - 1) * searchDTO.Limit
	if err := query.Offset(offset).Limit(searchDTO.Limit).Find(&settings).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search advertising settings: %w", err)
	}

	return settings, total, nil
}

// Settings validation

func (r *AdvertisingSettingsPostgresRepository) ValidateSettings(ctx context.Context, settings *domain.AdvertisingSettings) (*dto.SettingsValidationDTO, error) {
	validation := &dto.SettingsValidationDTO{
		IsValid: true,
	}

	// Validate budget
	if settings.DefaultBudget < 10000 {
		validation.IsValid = false
		validation.Errors = append(validation.Errors, "Default budget must be at least 10,000 VND")
	}

	// Validate duration
	if settings.DefaultDuration < 1 {
		validation.IsValid = false
		validation.Errors = append(validation.Errors, "Default duration must be at least 1 day")
	}

	// Validate campaign type
	if settings.DefaultCampaignType == "" {
		validation.IsValid = false
		validation.Errors = append(validation.Errors, "Default campaign type is required")
	}

	// Validate auto run budget if auto run is enabled
	if settings.AutoRunEnabled && settings.AutoRunBudget < 10000 {
		validation.IsValid = false
		validation.Errors = append(validation.Errors, "Auto run budget must be at least 10,000 VND when auto run is enabled")
	}

	// Check permissions for organization-wide settings
	if settings.ApplyToAllOrganization {
		validation.RequiresPermission = true
		validation.PermissionLevel = "ADMIN"
		validation.CanApplyToAll = true
	}

	return validation, nil
}