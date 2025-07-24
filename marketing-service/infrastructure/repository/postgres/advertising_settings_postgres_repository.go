package postgres

import (
	"context"
	"fmt"
	"marketing/internal/domain/model"
	"marketing/internal/domain/repository"
	"marketing/internal/dto"

	"gorm.io/gorm"
)

type AdvertisingSettingsPostgresRepository struct {
	db *gorm.DB
}

func NewAdvertisingSettingsPostgresRepository(db *gorm.DB) repository.AdvertisingSettingsRepository {
	return &AdvertisingSettingsPostgresRepository{
		db: db,
	}
}

// Advertising settings management

func (r *AdvertisingSettingsPostgresRepository) CreateAdvertisingSettings(ctx context.Context, settings *model.AdvertisingSettings) (*model.AdvertisingSettings, error) {
	if err := r.db.WithContext(ctx).Create(settings).Error; err != nil {
		return nil, fmt.Errorf("failed to create advertising settings: %w", err)
	}
	return settings, nil
}

func (r *AdvertisingSettingsPostgresRepository) UpdateAdvertisingSettings(ctx context.Context, id uint64, settings *model.AdvertisingSettings) (*model.AdvertisingSettings, error) {
	if err := r.db.WithContext(ctx).Model(&model.AdvertisingSettings{}).Where("id = ?", id).Updates(settings).Error; err != nil {
		return nil, fmt.Errorf("failed to update advertising settings: %w", err)
	}
	return r.GetAdvertisingSettingsByUser(ctx, settings.UserID)
}

func (r *AdvertisingSettingsPostgresRepository) GetAdvertisingSettingsByUser(ctx context.Context, userID uint64) (*model.AdvertisingSettings, error) {
	var settings model.AdvertisingSettings
	if err := r.db.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true).First(&settings).Error; err != nil {
		return nil, fmt.Errorf("advertising settings not found: %w", err)
	}
	return &settings, nil
}

func (r *AdvertisingSettingsPostgresRepository) GetAdvertisingSettingsByOrganization(ctx context.Context, organizationID uint64) (*model.AdvertisingSettings, error) {
	var settings model.AdvertisingSettings
	if err := r.db.WithContext(ctx).Where("organization_id = ? AND apply_to_all_organization = ? AND is_active = ?", organizationID, true, true).First(&settings).Error; err != nil {
		return nil, fmt.Errorf("organization advertising settings not found: %w", err)
	}
	return &settings, nil
}

func (r *AdvertisingSettingsPostgresRepository) DeleteAdvertisingSettings(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).Delete(&model.AdvertisingSettings{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete advertising settings: %w", err)
	}
	return nil
}

// Search and list

func (r *AdvertisingSettingsPostgresRepository) SearchAdvertisingSettings(ctx context.Context, searchDTO dto.AdvertisingSettingsSearchDTO) ([]*model.AdvertisingSettings, int64, error) {
	var settings []*model.AdvertisingSettings
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AdvertisingSettings{})

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

// Campaign defaults

func (r *AdvertisingSettingsPostgresRepository) GetCampaignDefaults(ctx context.Context, userID uint64, productType string) (*dto.CampaignDefaultsDTO, error) {
	// Get user settings
	settings, err := r.GetAdvertisingSettingsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if there are product type specific settings
	if settings.ProductTypeSettings != nil {
		if productSetting, exists := settings.ProductTypeSettings[productType]; exists {
			return &dto.CampaignDefaultsDTO{
				CampaignType:   settings.DefaultCampaignType,
				PackageID:      productSetting.DefaultPackageID,
				Budget:         productSetting.DefaultBudget,
				Duration:       productSetting.DefaultDuration,
				Priority:       productSetting.Priority,
				AutoRunEnabled: productSetting.AutoRunEnabled,
				SettingsSource: "USER",
			}, nil
		}
	}

	// Return general defaults
	return &dto.CampaignDefaultsDTO{
		CampaignType:   settings.DefaultCampaignType,
		PackageID:      settings.DefaultPackageID,
		Budget:         settings.DefaultBudget,
		Duration:       settings.DefaultDuration,
		Priority:       5,
		AutoRunEnabled: settings.AutoRunEnabled,
		SettingsSource: "USER",
	}, nil
}

func (r *AdvertisingSettingsPostgresRepository) GetOrganizationDefaults(ctx context.Context, organizationID uint64, productType string) (*dto.CampaignDefaultsDTO, error) {
	// Get organization settings
	settings, err := r.GetAdvertisingSettingsByOrganization(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	// Check if there are product type specific settings
	if settings.ProductTypeSettings != nil {
		if productSetting, exists := settings.ProductTypeSettings[productType]; exists {
			return &dto.CampaignDefaultsDTO{
				CampaignType:   settings.DefaultCampaignType,
				PackageID:      productSetting.DefaultPackageID,
				Budget:         productSetting.DefaultBudget,
				Duration:       productSetting.DefaultDuration,
				Priority:       productSetting.Priority,
				AutoRunEnabled: productSetting.AutoRunEnabled,
				SettingsSource: "ORGANIZATION",
			}, nil
		}
	}

	// Return general defaults
	return &dto.CampaignDefaultsDTO{
		CampaignType:   settings.DefaultCampaignType,
		PackageID:      settings.DefaultPackageID,
		Budget:         settings.DefaultBudget,
		Duration:       settings.DefaultDuration,
		Priority:       5,
		AutoRunEnabled: settings.AutoRunEnabled,
		SettingsSource: "ORGANIZATION",
	}, nil
}

// Auto campaign suggestions

func (r *AdvertisingSettingsPostgresRepository) CreateAutoCampaignSuggestion(ctx context.Context, suggestion *model.AutoCampaignSuggestion) (*model.AutoCampaignSuggestion, error) {
	if err := r.db.WithContext(ctx).Create(suggestion).Error; err != nil {
		return nil, fmt.Errorf("failed to create auto campaign suggestion: %w", err)
	}
	return suggestion, nil
}

func (r *AdvertisingSettingsPostgresRepository) UpdateAutoCampaignSuggestion(ctx context.Context, id uint64, suggestion *model.AutoCampaignSuggestion) (*model.AutoCampaignSuggestion, error) {
	if err := r.db.WithContext(ctx).Model(&model.AutoCampaignSuggestion{}).Where("id = ?", id).Updates(suggestion).Error; err != nil {
		return nil, fmt.Errorf("failed to update auto campaign suggestion: %w", err)
	}
	return suggestion, nil
}

func (r *AdvertisingSettingsPostgresRepository) GetAutoCampaignSuggestions(ctx context.Context, userID uint64, page, limit int) ([]*model.AutoCampaignSuggestion, int64, error) {
	var suggestions []*model.AutoCampaignSuggestion
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AutoCampaignSuggestion{}).Where("user_id = ?", userID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count auto campaign suggestions: %w", err)
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&suggestions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get auto campaign suggestions: %w", err)
	}

	return suggestions, total, nil
}

func (r *AdvertisingSettingsPostgresRepository) GetPendingAutoCampaignSuggestions(ctx context.Context, userID uint64) ([]*model.AutoCampaignSuggestion, error) {
	var suggestions []*model.AutoCampaignSuggestion
	if err := r.db.WithContext(ctx).Where("user_id = ? AND is_accepted = ? AND is_rejected = ?", userID, false, false).Find(&suggestions).Error; err != nil {
		return nil, fmt.Errorf("failed to get pending auto campaign suggestions: %w", err)
	}
	return suggestions, nil
}

func (r *AdvertisingSettingsPostgresRepository) AcceptAutoCampaignSuggestion(ctx context.Context, suggestionID uint64, campaignID uint64) error {
	if err := r.db.WithContext(ctx).Model(&model.AutoCampaignSuggestion{}).Where("id = ?", suggestionID).Updates(map[string]interface{}{
		"is_accepted":        true,
		"applied_campaign_id": campaignID,
	}).Error; err != nil {
		return fmt.Errorf("failed to accept auto campaign suggestion: %w", err)
	}
	return nil
}

func (r *AdvertisingSettingsPostgresRepository) RejectAutoCampaignSuggestion(ctx context.Context, suggestionID uint64) error {
	if err := r.db.WithContext(ctx).Model(&model.AutoCampaignSuggestion{}).Where("id = ?", suggestionID).Update("is_rejected", true).Error; err != nil {
		return fmt.Errorf("failed to reject auto campaign suggestion: %w", err)
	}
	return nil
}

// Settings validation

func (r *AdvertisingSettingsPostgresRepository) ValidateSettings(ctx context.Context, settings *model.AdvertisingSettings) (*dto.SettingsValidationDTO, error) {
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

func (r *AdvertisingSettingsPostgresRepository) CheckUserPermissions(ctx context.Context, userID uint64, organizationID uint64) (bool, string, error) {
	// TODO: Implement actual permission check
	// For now, return mock data
	return true, "ADMIN", nil
}

// Auto suggestion generation

func (r *AdvertisingSettingsPostgresRepository) GenerateAutoCampaignSuggestions(ctx context.Context, requestDTO *dto.AutoCampaignSuggestionRequestDTO) (*dto.AutoCampaignSuggestionResponseDTO, error) {
	// TODO: Implement actual suggestion generation logic
	// For now, return mock suggestions
	
	userID := uint64(123) // TODO: Get from context
	organizationID := uint64(1) // TODO: Get from context

	// Create mock suggestion
	suggestion := &model.AutoCampaignSuggestion{
		OrganizationID:        organizationID,
		UserID:                userID,
		ProductID:             requestDTO.ProductID,
		SuggestedCampaignType: "VIP",
		SuggestedPackageID:    1,
		SuggestedPackageName:  "VIP 7 ngày",
		SuggestedBudget:       100000,
		SuggestedDuration:     7,
		Reason:                "Sản phẩm mới, đề xuất gói VIP để tăng hiển thị",
		Confidence:            0.8,
		EstimatedReach:        &[]int{1000}[0],
		EstimatedClicks:       &[]int{50}[0],
	}

	// Save suggestion
	_, err := r.CreateAutoCampaignSuggestion(ctx, suggestion)
	if err != nil {
		return nil, err
	}

	// Get user settings for response
	settings, err := r.GetAdvertisingSettingsByUser(ctx, userID)
	if err != nil {
		// Use default settings if not found
		settings = &model.AdvertisingSettings{
			DefaultCampaignType: "VIP",
			DefaultBudget:       100000,
			DefaultDuration:     7,
		}
	}

	// Convert to DTO
	suggestionDTO := &dto.AutoCampaignSuggestionDTO{
		ProductID:             suggestion.ProductID,
		ProductName:           requestDTO.ProductName,
		ProductType:           requestDTO.ProductType,
		SuggestedCampaignType: suggestion.SuggestedCampaignType,
		SuggestedPackageID:    suggestion.SuggestedPackageID,
		SuggestedPackageName:  suggestion.SuggestedPackageName,
		SuggestedBudget:       suggestion.SuggestedBudget,
		SuggestedDuration:     suggestion.SuggestedDuration,
		Reason:                suggestion.Reason,
		Confidence:            suggestion.Confidence,
		EstimatedReach:        suggestion.EstimatedReach,
		EstimatedClicks:       suggestion.EstimatedClicks,
	}

	// Convert settings to DTO
	settingsDTO := dto.AdvertisingSettingsDTO{
		DefaultCampaignType: settings.DefaultCampaignType,
		DefaultBudget:       settings.DefaultBudget,
		DefaultDuration:     settings.DefaultDuration,
	}

	return &dto.AutoCampaignSuggestionResponseDTO{
		ProductID:        requestDTO.ProductID,
		Suggestions:      []dto.AutoCampaignSuggestionDTO{*suggestionDTO},
		TotalSuggestions: 1,
		BestSuggestion:   suggestionDTO,
		SettingsUsed:     settingsDTO,
	}, nil
}

// Settings inheritance

func (r *AdvertisingSettingsPostgresRepository) GetInheritedSettings(ctx context.Context, userID uint64, organizationID uint64) (*model.AdvertisingSettings, error) {
	// Try to get user settings first
	settings, err := r.GetAdvertisingSettingsByUser(ctx, userID)
	if err == nil {
		return settings, nil
	}

	// Fallback to organization settings
	return r.GetAdvertisingSettingsByOrganization(ctx, organizationID)
}

func (r *AdvertisingSettingsPostgresRepository) ApplySettingsToOrganization(ctx context.Context, settingsID uint64, organizationID uint64) error {
	// TODO: Implement organization-wide settings application
	// This would involve copying settings to all users in the organization
	return nil
} 