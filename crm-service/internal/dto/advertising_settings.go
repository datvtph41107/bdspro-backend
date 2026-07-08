package dto

// AdvertisingSettingsDTO for advertising settings
type AdvertisingSettingsDTO struct {
	ID             uint64 `json:"id"`
	OrganizationID uint64 `json:"organization_id"`
	UserID         uint64 `json:"user_id"`

	// Default campaign settings
	DefaultCampaignType string  `json:"default_campaign_type"` // "VIP", "META_ADS", "GOOGLE_ADS", "NONE"
	DefaultPackageID    *uint64 `json:"default_package_id,omitempty"`
	DefaultPackageName  string  `json:"default_package_name,omitempty"`
	DefaultBudget       float64 `json:"default_budget"`
	DefaultDuration     int     `json:"default_duration"`

	// Auto advertising settings
	AutoSuggestEnabled bool    `json:"auto_suggest_enabled"` // Tự động đề xuất quảng cáo
	AutoRunEnabled     bool    `json:"auto_run_enabled"`     // Tự động chạy quảng cáo
	AutoRunBudget      float64 `json:"auto_run_budget"`
	AutoRunPackageID   *uint64 `json:"auto_run_package_id,omitempty"`
	AutoRunPackageName string  `json:"auto_run_package_name,omitempty"`

	// Budget limits and alerts
	MonthlyBudgetLimit  float64 `json:"monthly_budget_limit"`
	AlertOnBudgetExceed bool    `json:"alert_on_budget_exceed"`
	AlertEmail          string  `json:"alert_email,omitempty"`
	AlertPhone          string  `json:"alert_phone,omitempty"`

	// Product type specific settings
	ProductTypeSettings map[string]ProductTypeSettingDTO `json:"product_type_settings,omitempty"`

	// System settings
	ApplyToAllOrganization bool `json:"apply_to_all_organization"` // Lưu & áp dụng cho toàn hệ thống
	IsActive               bool `json:"is_active"`

	// Metadata
	CreatedBy uint64 `json:"created_by"`
	UpdatedBy uint64 `json:"updated_by"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ProductTypeSettingDTO for product type specific settings
type ProductTypeSettingDTO struct {
	ProductType      string  `json:"product_type"` // "LAND", "HOUSE", "APARTMENT", "COMMERCIAL"
	DefaultBudget    float64 `json:"default_budget"`
	DefaultDuration  int     `json:"default_duration"`
	DefaultPackageID *uint64 `json:"default_package_id,omitempty"`
	AutoRunEnabled   bool    `json:"auto_run_enabled"`
	AutoRunBudget    float64 `json:"auto_run_budget"`
	Priority         int     `json:"priority"` // 1-10, higher = more priority
}

// AdvertisingSettingsCreateDTO for creating advertising settings
type AdvertisingSettingsCreateDTO struct {
	OrganizationID uint64 `json:"organization_id"`
	UserID         uint64 `json:"user_id"`

	// Default campaign settings
	DefaultCampaignType string  `json:"default_campaign_type" binding:"required"`
	DefaultPackageID    *uint64 `json:"default_package_id,omitempty"`
	DefaultBudget       float64 `json:"default_budget" binding:"required,min=10000"`
	DefaultDuration     int     `json:"default_duration" binding:"required,min=1"`

	// Auto advertising settings
	AutoSuggestEnabled bool    `json:"auto_suggest_enabled"`
	AutoRunEnabled     bool    `json:"auto_run_enabled"`
	AutoRunBudget      float64 `json:"auto_run_budget,omitempty"`
	AutoRunPackageID   *uint64 `json:"auto_run_package_id,omitempty"`

	// Budget limits and alerts
	MonthlyBudgetLimit  float64 `json:"monthly_budget_limit" binding:"required,min=0"`
	AlertOnBudgetExceed bool    `json:"alert_on_budget_exceed"`
	AlertEmail          string  `json:"alert_email,omitempty"`
	AlertPhone          string  `json:"alert_phone,omitempty"`

	// Product type specific settings
	ProductTypeSettings map[string]ProductTypeSettingDTO `json:"product_type_settings,omitempty"`

	// System settings
	ApplyToAllOrganization bool `json:"apply_to_all_organization"`
}

// AdvertisingSettingsUpdateDTO for updating advertising settings
type AdvertisingSettingsUpdateDTO struct {
	// Default campaign settings
	DefaultCampaignType *string  `json:"default_campaign_type,omitempty"`
	DefaultPackageID    *uint64  `json:"default_package_id,omitempty"`
	DefaultBudget       *float64 `json:"default_budget,omitempty"`
	DefaultDuration     *int     `json:"default_duration,omitempty"`

	// Auto advertising settings
	AutoSuggestEnabled *bool    `json:"auto_suggest_enabled,omitempty"`
	AutoRunEnabled     *bool    `json:"auto_run_enabled,omitempty"`
	AutoRunBudget      *float64 `json:"auto_run_budget,omitempty"`
	AutoRunPackageID   *uint64  `json:"auto_run_package_id,omitempty"`

	// Budget limits and alerts
	MonthlyBudgetLimit  *float64 `json:"monthly_budget_limit,omitempty"`
	AlertOnBudgetExceed *bool    `json:"alert_on_budget_exceed,omitempty"`
	AlertEmail          *string  `json:"alert_email,omitempty"`
	AlertPhone          *string  `json:"alert_phone,omitempty"`

	// Product type specific settings
	ProductTypeSettings *map[string]ProductTypeSettingDTO `json:"product_type_settings,omitempty"`

	// System settings
	ApplyToAllOrganization *bool `json:"apply_to_all_organization,omitempty"`
}

// AdvertisingSettingsSearchDTO for searching advertising settings
type AdvertisingSettingsSearchDTO struct {
	OrganizationID      *uint64 `json:"organization_id,omitempty"`
	UserID              *uint64 `json:"user_id,omitempty"`
	DefaultCampaignType string  `json:"default_campaign_type,omitempty"`
	AutoSuggestEnabled  *bool   `json:"auto_suggest_enabled,omitempty"`
	AutoRunEnabled      *bool   `json:"auto_run_enabled,omitempty"`
	IsActive            *bool   `json:"is_active,omitempty"`
	Page                int     `json:"page"`
	Limit               int     `json:"limit"`
}

// AutoCampaignSuggestionDTO for auto campaign suggestions
type AutoCampaignSuggestionDTO struct {
	ProductID             uint64  `json:"product_id"`
	ProductName           string  `json:"product_name"`
	ProductType           string  `json:"product_type"`
	SuggestedCampaignType string  `json:"suggested_campaign_type"`
	SuggestedPackageID    uint64  `json:"suggested_package_id"`
	SuggestedPackageName  string  `json:"suggested_package_name"`
	SuggestedBudget       float64 `json:"suggested_budget"`
	SuggestedDuration     int     `json:"suggested_duration"`
	Reason                string  `json:"reason"`     // Why this suggestion
	Confidence            float64 `json:"confidence"` // 0-1, how confident the suggestion is
	EstimatedReach        *int    `json:"estimated_reach,omitempty"`
	EstimatedClicks       *int    `json:"estimated_clicks,omitempty"`
}

// AutoCampaignSuggestionRequestDTO for requesting auto campaign suggestions
type AutoCampaignSuggestionRequestDTO struct {
	ProductID         uint64   `json:"product_id" binding:"required"`
	ProductName       string   `json:"product_name"`
	ProductType       string   `json:"product_type"`
	ProductPrice      *float64 `json:"product_price,omitempty"`
	ProductLocation   string   `json:"product_location,omitempty"`
	UserBudget        *float64 `json:"user_budget,omitempty"`
	PreferredDuration *int     `json:"preferred_duration,omitempty"`
}

// AutoCampaignSuggestionResponseDTO for auto campaign suggestion response
type AutoCampaignSuggestionResponseDTO struct {
	ProductID        uint64                      `json:"product_id"`
	Suggestions      []AutoCampaignSuggestionDTO `json:"suggestions"`
	TotalSuggestions int                         `json:"total_suggestions"`
	BestSuggestion   *AutoCampaignSuggestionDTO  `json:"best_suggestion,omitempty"`
	SettingsUsed     AdvertisingSettingsDTO      `json:"settings_used"`
}

// CampaignDefaultsDTO for campaign defaults based on settings
type CampaignDefaultsDTO struct {
	CampaignType   string  `json:"campaign_type"`
	PackageID      *uint64 `json:"package_id,omitempty"`
	PackageName    string  `json:"package_name,omitempty"`
	Budget         float64 `json:"budget"`
	Duration       int     `json:"duration"`
	TargetAudience string  `json:"target_audience,omitempty"`
	TargetLocation string  `json:"target_location,omitempty"`
	Priority       int     `json:"priority"`
	AutoRunEnabled bool    `json:"auto_run_enabled"`
	SettingsSource string  `json:"settings_source"` // "USER", "ORGANIZATION", "SYSTEM"
}

// SettingsValidationDTO for settings validation
type SettingsValidationDTO struct {
	IsValid            bool     `json:"is_valid"`
	Errors             []string `json:"errors,omitempty"`
	Warnings           []string `json:"warnings,omitempty"`
	Suggestions        []string `json:"suggestions,omitempty"`
	CanApplyToAll      bool     `json:"can_apply_to_all"`
	RequiresPermission bool     `json:"requires_permission"`
	PermissionLevel    string   `json:"permission_level,omitempty"`
}