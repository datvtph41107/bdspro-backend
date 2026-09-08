package dto

// CampaignDetailDTO for campaign detail display
type CampaignDetailDTO struct {
	ID              uint64  `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Type            string  `json:"type"`
	Status          string  `json:"status"`
	
	// Product information
	ProductID       uint64  `json:"product_id"`
	ProductName     string  `json:"product_name"`
	ProductImage    string  `json:"product_image"`
	ProductType     string  `json:"product_type"`
	
	// Package information
	PackageID       uint64  `json:"package_id"`
	PackageName     string  `json:"package_name"`
	PackageType     string  `json:"package_type"`
	
	// Campaign settings
	Budget          float64 `json:"budget"`
	Duration        int     `json:"duration"`
	StartDate       string  `json:"start_date"`
	EndDate         string  `json:"end_date"`
	Currency        string  `json:"currency"`
	
	// Performance metrics
	Impressions     int     `json:"impressions"`
	Clicks          int     `json:"clicks"`
	CTR             float64 `json:"ctr"` // Click-through rate
	CPC             float64 `json:"cpc"` // Cost per click
	SpentAmount     float64 `json:"spent_amount"`
	RemainingBudget float64 `json:"remaining_budget"`
	
	// Wallet information
	WalletID        string  `json:"wallet_id"`
	WalletBalance   float64 `json:"wallet_balance"`
	TransactionID   uint64  `json:"transaction_id"`
	
	// Permissions
	CanEdit         bool    `json:"can_edit"`
	CanPause        bool    `json:"can_pause"`
	CanResume       bool    `json:"can_resume"`
	CanExtend       bool    `json:"can_extend"`
	CanTopup        bool    `json:"can_topup"`
	
	// Metadata
	CreatedBy       uint64  `json:"created_by"`
	OrganizationID  uint64  `json:"organization_id"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// CampaignPerformanceDTO for performance metrics
type CampaignPerformanceDTO struct {
	CampaignID  uint64                    `json:"campaign_id"`
	Date        string                    `json:"date"`
	Impressions int                       `json:"impressions"`
	Clicks      int                       `json:"clicks"`
	CTR         float64                   `json:"ctr"`
	CPC         float64                   `json:"cpc"`
	Spent       float64                   `json:"spent"`
	DailyData   []CampaignDailyDataDTO    `json:"daily_data"`
	Summary     CampaignPerformanceSummaryDTO `json:"summary"`
}

// CampaignDailyDataDTO for daily performance data
type CampaignDailyDataDTO struct {
	Date        string  `json:"date"`
	Impressions int     `json:"impressions"`
	Clicks      int     `json:"clicks"`
	CTR         float64 `json:"ctr"`
	CPC         float64 `json:"cpc"`
	Spent       float64 `json:"spent"`
}

// CampaignPerformanceSummaryDTO for performance summary
type CampaignPerformanceSummaryDTO struct {
	TotalImpressions int     `json:"total_impressions"`
	TotalClicks      int     `json:"total_clicks"`
	AverageCTR       float64 `json:"average_ctr"`
	AverageCPC       float64 `json:"average_cpc"`
	TotalSpent       float64 `json:"total_spent"`
	RemainingBudget  float64 `json:"remaining_budget"`
	ROI              float64 `json:"roi"` // Return on Investment
}

// CampaignActionDTO for campaign actions
type CampaignActionDTO struct {
	ID          uint64 `json:"id"`
	CampaignID  uint64 `json:"campaign_id"`
	Action      string `json:"action"` // "PAUSE", "RESUME", "EXTEND", "TOPUP", "EDIT"
	Description string `json:"description"`
	PerformedBy uint64 `json:"performed_by"`
	PerformedAt string `json:"performed_at"`
	OldValue    string `json:"old_value,omitempty"`
	NewValue    string `json:"new_value,omitempty"`
}

// CampaignHistoryDTO for campaign history
type CampaignHistoryDTO struct {
	ID              uint64 `json:"id"`
	CampaignID      uint64 `json:"campaign_id"`
	Action          string `json:"action"`
	Description     string `json:"description"`
	PerformedBy     uint64 `json:"performed_by"`
	PerformedByUser string `json:"performed_by_user"`
	PerformedAt     string `json:"performed_at"`
	IPAddress       string `json:"ip_address,omitempty"`
	UserAgent       string `json:"user_agent,omitempty"`
}

// CampaignExportDTO for campaign data export
type CampaignExportDTO struct {
	CampaignID      uint64                    `json:"campaign_id"`
	CampaignName    string                    `json:"campaign_name"`
	ExportType      string                    `json:"export_type"` // "PDF", "EXCEL"
	DateRange       string                    `json:"date_range"`
	PerformanceData []CampaignDailyDataDTO    `json:"performance_data"`
	Summary         CampaignPerformanceSummaryDTO `json:"summary"`
	GeneratedAt     string                    `json:"generated_at"`
	GeneratedBy     uint64                    `json:"generated_by"`
}

// CampaignOptimizationDTO for AI optimization suggestions
type CampaignOptimizationDTO struct {
	CampaignID      uint64   `json:"campaign_id"`
	Suggestions     []string `json:"suggestions"`
	Priority        string   `json:"priority"` // "HIGH", "MEDIUM", "LOW"
	EstimatedImpact string   `json:"estimated_impact"`
	GeneratedAt     string   `json:"generated_at"`
}

// CampaignStatusUpdateDTO for status updates
type CampaignStatusUpdateDTO struct {
	CampaignID  uint64 `json:"campaign_id" binding:"required"`
	NewStatus   string `json:"new_status" binding:"required"`
	Reason      string `json:"reason,omitempty"`
	Description string `json:"description,omitempty"`
}

// CampaignExtendDTO for campaign extension
type CampaignExtendDTO struct {
	CampaignID      uint64  `json:"campaign_id" binding:"required"`
	ExtensionDays   int     `json:"extension_days" binding:"required,min=1"`
	AdditionalBudget float64 `json:"additional_budget,omitempty"`
	PaymentMethod   string  `json:"payment_method" binding:"required"`
	TermsAccepted   bool    `json:"terms_accepted" binding:"required"`
} 