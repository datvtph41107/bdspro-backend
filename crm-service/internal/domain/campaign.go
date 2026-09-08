package domain

import (
	_models "common/models"
	"time"
)

// CampaignType represents the type of advertising campaign
type CampaignType string

const (
	CampaignTypeVIP       CampaignType = "VIP"
	CampaignTypeMetaAds   CampaignType = "META_ADS"
	CampaignTypeGoogleAds CampaignType = "GOOGLE_ADS"
	CampaignTypePR        CampaignType = "PR"
)

// CampaignStatus represents the status of a campaign
type CampaignStatus string

const (
	CampaignStatusDraft     CampaignStatus = "DRAFT"
	CampaignStatusActive    CampaignStatus = "ACTIVE"
	CampaignStatusPaused    CampaignStatus = "PAUSED"
	CampaignStatusCompleted CampaignStatus = "COMPLETED"
	CampaignStatusCancelled CampaignStatus = "CANCELLED"
)

// Campaign represents an advertising campaign
type Campaign struct {
	_models.BaseEntity
	Name        string         `json:"name" gorm:"not null;size:100"`
	Description string         `json:"description"`
	Type        CampaignType   `json:"type" gorm:"not null"`
	Status      CampaignStatus `json:"status" gorm:"default:'DRAFT'"`

	// Product/Listing information
	ProductID   *uint64 `json:"product_id"`
	ProductName string  `json:"product_name"`
	ProductType string  `json:"product_type"` // "PRODUCT", "LISTING"

	// Campaign settings
	PackageID   *uint64    `json:"package_id"`
	PackageName string     `json:"package_name"`
	Budget      float64    `json:"budget"`
	Duration    int        `json:"duration"` // in days
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`

	// Target audience
	TargetAudience string `json:"target_audience"`
	TargetLocation string `json:"target_location"`

	// Campaign settings
	IsActive     bool   `json:"is_active" gorm:"default:true"`
	Priority     int    `json:"priority" gorm:"default:0"`
	InternalNote string `json:"internal_note"`

	// Payment/Wallet information
	WalletID      *uint64 `json:"wallet_id"`
	WalletBalance float64 `json:"wallet_balance" gorm:"default:0"`
	Currency      string  `json:"currency" gorm:"default:'VND'"`
	TransactionID *uint64 `json:"transaction_id"` // Payment transaction ID

	// Metadata
	CreatedBy      uint64 `json:"created_by"`
	UpdatedBy      uint64 `json:"updated_by"`
	OrganizationID uint64 `json:"organization_id"`
}

func (Campaign) TableName() string {
	return "campaigns"
}