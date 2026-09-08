package dto

import (
	"crm/internal/domain"
	"time"
)

// CampaignCreateDTO represents the data for creating a new campaign
type CampaignCreateDTO struct {
	Name           string              `json:"name" binding:"required,max=100"`
	Description    string              `json:"description"`
	Type           domain.CampaignType `json:"type" binding:"required"`
	ProductID      *uint64             `json:"product_id"`
	ProductName    string              `json:"product_name"`
	ProductType    string              `json:"product_type"`
	PackageID      *uint64             `json:"package_id" binding:"required"`
	PackageName    string              `json:"package_name"`
	Budget         float64             `json:"budget" binding:"min=0"`
	Duration       int                 `json:"duration" binding:"required,min=1"`
	StartDate      *time.Time          `json:"start_date"`
	EndDate        *time.Time          `json:"end_date"`
	TargetAudience string              `json:"target_audience"`
	TargetLocation string              `json:"target_location"`
	Priority       int                 `json:"priority"`
	InternalNote   string              `json:"internal_note"`
	TermsAccepted  bool                `json:"terms_accepted" binding:"required"`
}

// CampaignUpdateDTO represents the data for updating a campaign
type CampaignUpdateDTO struct {
	Name           string                `json:"name" binding:"max=100"`
	Description    string                `json:"description"`
	Type           domain.CampaignType   `json:"type"`
	ProductID      *uint64               `json:"product_id"`
	ProductName    string                `json:"product_name"`
	ProductType    string                `json:"product_type"`
	PackageID      *uint64               `json:"package_id"`
	PackageName    string                `json:"package_name"`
	Budget         float64               `json:"budget"`
	Duration       int                   `json:"duration"`
	StartDate      *time.Time            `json:"start_date"`
	EndDate        *time.Time            `json:"end_date"`
	TargetAudience string                `json:"target_audience"`
	TargetLocation string                `json:"target_location"`
	Status         domain.CampaignStatus `json:"status"`
	IsActive       *bool                 `json:"is_active"`
	Priority       int                   `json:"priority"`
	InternalNote   string                `json:"internal_note"`
}

// CampaignSearchDTO represents the search criteria for campaigns
type CampaignSearchDTO struct {
	Name        string                `form:"name"`
	Type        domain.CampaignType   `form:"type"`
	Status      domain.CampaignStatus `form:"status"`
	ProductID   *uint64               `form:"product_id"`
	ProductType string                `form:"product_type"`
	PackageID   *uint64               `form:"package_id"`
	IsActive    *bool                 `form:"is_active"`
	FromDate    *time.Time            `form:"from_date"`
	ToDate      *time.Time            `form:"to_date"`
	Page        int                   `form:"page" binding:"min=1"`
	Limit       int                   `form:"limit" binding:"min=1,max=100"`
}

// CampaignResponseDTO represents the response data for a campaign
type CampaignResponseDTO struct {
	ID             uint64                `json:"id"`
	Name           string                `json:"name"`
	Description    string                `json:"description"`
	Type           domain.CampaignType   `json:"type"`
	Status         domain.CampaignStatus `json:"status"`
	ProductID      *uint64               `json:"product_id"`
	ProductName    string                `json:"product_name"`
	ProductType    string                `json:"product_type"`
	PackageID      *uint64               `json:"package_id"`
	PackageName    string                `json:"package_name"`
	Budget         float64               `json:"budget"`
	Duration       int                   `json:"duration"`
	StartDate      *time.Time            `json:"start_date"`
	EndDate        *time.Time            `json:"end_date"`
	TargetAudience string                `json:"target_audience"`
	TargetLocation string                `json:"target_location"`
	IsActive       bool                  `json:"is_active"`
	Priority       int                   `json:"priority"`
	InternalNote   string                `json:"internal_note"`
	CreatedBy      uint64                `json:"created_by"`
	OrganizationID uint64                `json:"organization_id"`
	WalletID       *string               `json:"wallet_id"`
	WalletBalance  float64               `json:"wallet_balance"`
	Currency       string                `json:"currency"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
}

// GetOffset returns the offset for pagination
func (d *CampaignSearchDTO) GetOffset() int {
	return (d.Page - 1) * d.GetLimit()
}

// GetLimit returns the limit for pagination
func (d *CampaignSearchDTO) GetLimit() int {
	if d.Limit <= 0 {
		return 10
	}
	return d.Limit
}

// CampaignBudgetTopupDTO for adding budget to campaign
type CampaignBudgetTopupDTO struct {
	CampaignID        uint64  `json:"campaign_id" binding:"required"`
	Amount            float64 `json:"amount" binding:"required,min=0"`
	PaymentMethod     string  `json:"payment_method" binding:"required"`
	ExternalPaymentID string  `json:"external_payment_id"`
	Description       string  `json:"description"`
}

// CampaignBudgetTopupResponseDTO for budget topup response
type CampaignBudgetTopupResponseDTO struct {
	CampaignID        uint64  `json:"campaign_id"`
	TransactionID     uint64  `json:"transaction_id"`
	Amount            float64 `json:"amount"`
	PreviousBalance   float64 `json:"previous_balance"`
	NewBalance        float64 `json:"new_balance"`
	Currency          string  `json:"currency"`
	Status            string  `json:"status"`
	TransactionCode   string  `json:"transaction_code"`
	PaymentMethod     string  `json:"payment_method"`
	ExternalPaymentID string  `json:"external_payment_id"`
	CreatedAt         string  `json:"created_at"`
	Message           string  `json:"message"`
}

// CampaignWalletInfoDTO for campaign wallet information
type CampaignWalletInfoDTO struct {
	CampaignID      uint64  `json:"campaign_id"`
	WalletID        string  `json:"wallet_id"`
	Balance         float64 `json:"balance"`
	Currency        string  `json:"currency"`
	TotalDeposited  float64 `json:"total_deposited"`
	TotalSpent      float64 `json:"total_spent"`
	AvailableBudget float64 `json:"available_budget"`
	CampaignBudget  float64 `json:"campaign_budget"`
	RemainingBudget float64 `json:"remaining_budget"`
	LastUpdated     string  `json:"last_updated"`
}