package domain

import (
	_models "common/models"
	"time"
)

// BudgetConfig represents budget configuration for campaigns
type BudgetConfig struct {
	_models.BaseEntity
	CampaignID     uint64 `json:"campaign_id" gorm:"not null;uniqueIndex"`
	OrganizationID uint64 `json:"organization_id" gorm:"not null"`

	// Budget limits
	MonthlyLimit float64  `json:"monthly_limit" gorm:"not null"`
	DailyLimit   *float64 `json:"daily_limit"`
	TotalLimit   *float64 `json:"total_limit"`

	// Alert thresholds
	AlertThreshold float64 `json:"alert_threshold" gorm:"not null;default:80"` // Percentage
	AlertEnabled   bool    `json:"alert_enabled" gorm:"default:true"`
	AlertEmail     string  `json:"alert_email"`
	AlertPhone     string  `json:"alert_phone"`

	// Auto settings
	AutoTopup          bool    `json:"auto_topup" gorm:"default:false"`
	AutoTopupAmount    float64 `json:"auto_topup_amount"`
	AutoTopupThreshold float64 `json:"auto_topup_threshold"` // Percentage

	// Metadata
	CreatedBy uint64 `json:"created_by"`
	UpdatedBy uint64 `json:"updated_by"`
}

func (BudgetConfig) TableName() string {
	return "budget_configs"
}

// BudgetUsage represents daily budget usage tracking
type BudgetUsage struct {
	ID         uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	CampaignID uint64    `json:"campaign_id" gorm:"not null;index"`
	Date       time.Time `json:"date" gorm:"not null;index"`

	// Usage metrics
	Impressions int     `json:"impressions" gorm:"default:0"`
	Clicks      int     `json:"clicks" gorm:"default:0"`
	SpentAmount float64 `json:"spent_amount" gorm:"default:0"`
	CTR         float64 `json:"ctr" gorm:"default:0"`
	CPC         float64 `json:"cpc" gorm:"default:0"`

	// Limits and usage
	DailyLimit *float64 `json:"daily_limit"`
	DailyUsage *float64 `json:"daily_usage"`

	// Metadata
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (BudgetUsage) TableName() string {
	return "budget_usages"
}

// BudgetAlert represents budget alerts
type BudgetAlert struct {
	ID             uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	CampaignID     uint64 `json:"campaign_id" gorm:"not null;index"`
	OrganizationID uint64 `json:"organization_id" gorm:"not null"`

	// Alert details
	AlertType    string  `json:"alert_type" gorm:"not null"`  // "THRESHOLD", "LIMIT_EXCEEDED", "LOW_BALANCE"
	AlertLevel   string  `json:"alert_level" gorm:"not null"` // "WARNING", "CRITICAL"
	Message      string  `json:"message" gorm:"not null"`
	CurrentUsage float64 `json:"current_usage" gorm:"not null"`
	Threshold    float64 `json:"threshold" gorm:"not null"`

	// Status
	IsRead bool       `json:"is_read" gorm:"default:false"`
	IsSent bool       `json:"is_sent" gorm:"default:false"`
	SentAt *time.Time `json:"sent_at"`

	// Metadata
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (BudgetAlert) TableName() string {
	return "budget_alerts"
}