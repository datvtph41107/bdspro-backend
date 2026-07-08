package domain

import (
	_models "common/models"
	"database/sql/driver"
	"encoding/json"
	"time"
)

// AdvertisingSettings represents advertising settings for users/organizations
type AdvertisingSettings struct {
	_models.BaseEntity
	OrganizationID uint64 `json:"organization_id" gorm:"not null;index"`
	UserID         uint64 `json:"user_id" gorm:"not null;index"`

	// Default campaign settings
	DefaultCampaignType string  `json:"default_campaign_type" gorm:"not null;default:'NONE'"` // "VIP", "META_ADS", "GOOGLE_ADS", "NONE"
	DefaultPackageID    *uint64 `json:"default_package_id"`
	DefaultPackageName  string  `json:"default_package_name"`
	DefaultBudget       float64 `json:"default_budget" gorm:"not null;default:0"`
	DefaultDuration     int     `json:"default_duration" gorm:"not null;default:7"`

	// Auto advertising settings
	AutoSuggestEnabled bool    `json:"auto_suggest_enabled" gorm:"default:false"` // Tự động đề xuất quảng cáo
	AutoRunEnabled     bool    `json:"auto_run_enabled" gorm:"default:false"`     // Tự động chạy quảng cáo
	AutoRunBudget      float64 `json:"auto_run_budget" gorm:"default:0"`
	AutoRunPackageID   *uint64 `json:"auto_run_package_id"`
	AutoRunPackageName string  `json:"auto_run_package_name"`

	// Budget limits and alerts
	MonthlyBudgetLimit  float64 `json:"monthly_budget_limit" gorm:"not null;default:0"`
	AlertOnBudgetExceed bool    `json:"alert_on_budget_exceed" gorm:"default:false"`
	AlertEmail          string  `json:"alert_email"`
	AlertPhone          string  `json:"alert_phone"`

	// Product type specific settings (JSON)
	ProductTypeSettings ProductTypeSettings `json:"product_type_settings" gorm:"type:json"`

	// System settings
	ApplyToAllOrganization bool `json:"apply_to_all_organization" gorm:"default:false"` // Lưu & áp dụng cho toàn hệ thống
	IsActive               bool `json:"is_active" gorm:"default:true"`

	// Metadata
	CreatedBy uint64 `json:"created_by"`
	UpdatedBy uint64 `json:"updated_by"`
}

func (AdvertisingSettings) TableName() string {
	return "advertising_settings"
}

// ProductTypeSettings represents product type specific settings
type ProductTypeSettings map[string]ProductTypeSetting

// ProductTypeSetting represents settings for a specific product type
type ProductTypeSetting struct {
	ProductType      string  `json:"product_type"` // "LAND", "HOUSE", "APARTMENT", "COMMERCIAL"
	DefaultBudget    float64 `json:"default_budget"`
	DefaultDuration  int     `json:"default_duration"`
	DefaultPackageID *uint64 `json:"default_package_id,omitempty"`
	AutoRunEnabled   bool    `json:"auto_run_enabled"`
	AutoRunBudget    float64 `json:"auto_run_budget"`
	Priority         int     `json:"priority"` // 1-10, higher = more priority
}

// Value implements driver.Valuer interface for JSON serialization
func (pts ProductTypeSettings) Value() (driver.Value, error) {
	if pts == nil {
		return nil, nil
	}
	return json.Marshal(pts)
}

// Scan implements sql.Scanner interface for JSON deserialization
func (pts *ProductTypeSettings) Scan(value interface{}) error {
	if value == nil {
		*pts = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, pts)
}

// AutoCampaignSuggestion represents auto campaign suggestions
type AutoCampaignSuggestion struct {
	ID             uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	OrganizationID uint64 `json:"organization_id" gorm:"not null;index"`
	UserID         uint64 `json:"user_id" gorm:"not null;index"`
	ProductID      uint64 `json:"product_id" gorm:"not null;index"`

	// Suggestion details
	SuggestedCampaignType string  `json:"suggested_campaign_type" gorm:"not null"`
	SuggestedPackageID    uint64  `json:"suggested_package_id"`
	SuggestedPackageName  string  `json:"suggested_package_name"`
	SuggestedBudget       float64 `json:"suggested_budget" gorm:"not null"`
	SuggestedDuration     int     `json:"suggested_duration" gorm:"not null"`
	Reason                string  `json:"reason" gorm:"not null"`
	Confidence            float64 `json:"confidence" gorm:"not null;default:0"`
	EstimatedReach        *int    `json:"estimated_reach"`
	EstimatedClicks       *int    `json:"estimated_clicks"`

	// Status
	IsAccepted        bool    `json:"is_accepted" gorm:"default:false"`
	IsRejected        bool    `json:"is_rejected" gorm:"default:false"`
	IsApplied         bool    `json:"is_applied" gorm:"default:false"`
	AppliedCampaignID *uint64 `json:"applied_campaign_id"`

	// Metadata
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (AutoCampaignSuggestion) TableName() string {
	return "auto_campaign_suggestions"
}