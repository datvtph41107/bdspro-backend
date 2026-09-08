package dto

// BudgetConfigDTO for budget configuration
type BudgetConfigDTO struct {
	ID             uint64 `json:"id"`
	CampaignID     uint64 `json:"campaign_id" binding:"required"`
	OrganizationID uint64 `json:"organization_id"`

	// Budget limits
	MonthlyLimit float64  `json:"monthly_limit" binding:"required,min=0"`
	DailyLimit   *float64 `json:"daily_limit,omitempty"`
	TotalLimit   *float64 `json:"total_limit,omitempty"`

	// Alert thresholds
	AlertThreshold float64 `json:"alert_threshold" binding:"required,min=0,max=100"` // Percentage
	AlertEnabled   bool    `json:"alert_enabled" gorm:"default:true"`
	AlertEmail     string  `json:"alert_email,omitempty"`
	AlertPhone     string  `json:"alert_phone,omitempty"`

	// Auto settings
	AutoTopup          bool    `json:"auto_topup" gorm:"default:false"`
	AutoTopupAmount    float64 `json:"auto_topup_amount,omitempty"`
	AutoTopupThreshold float64 `json:"auto_topup_threshold,omitempty"` // Percentage

	// Metadata
	CreatedBy uint64 `json:"created_by"`
	UpdatedBy uint64 `json:"updated_by"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// BudgetConfigCreateDTO for creating budget config
type BudgetConfigCreateDTO struct {
	CampaignID         uint64   `json:"campaign_id" binding:"required"`
	MonthlyLimit       float64  `json:"monthly_limit" binding:"required,min=0"`
	DailyLimit         *float64 `json:"daily_limit,omitempty"`
	TotalLimit         *float64 `json:"total_limit,omitempty"`
	AlertThreshold     float64  `json:"alert_threshold" binding:"required,min=0,max=100"`
	AlertEnabled       bool     `json:"alert_enabled"`
	AlertEmail         string   `json:"alert_email,omitempty"`
	AlertPhone         string   `json:"alert_phone,omitempty"`
	AutoTopup          bool     `json:"auto_topup"`
	AutoTopupAmount    float64  `json:"auto_topup_amount,omitempty"`
	AutoTopupThreshold float64  `json:"auto_topup_threshold,omitempty"`
}

// BudgetConfigUpdateDTO for updating budget config
type BudgetConfigUpdateDTO struct {
	MonthlyLimit       *float64 `json:"monthly_limit,omitempty"`
	DailyLimit         *float64 `json:"daily_limit,omitempty"`
	TotalLimit         *float64 `json:"total_limit,omitempty"`
	AlertThreshold     *float64 `json:"alert_threshold,omitempty"`
	AlertEnabled       *bool    `json:"alert_enabled,omitempty"`
	AlertEmail         *string  `json:"alert_email,omitempty"`
	AlertPhone         *string  `json:"alert_phone,omitempty"`
	AutoTopup          *bool    `json:"auto_topup,omitempty"`
	AutoTopupAmount    *float64 `json:"auto_topup_amount,omitempty"`
	AutoTopupThreshold *float64 `json:"auto_topup_threshold,omitempty"`
}

// BudgetStatusDTO for budget status and usage
type BudgetStatusDTO struct {
	CampaignID     uint64 `json:"campaign_id"`
	CampaignName   string `json:"campaign_name"`
	CampaignStatus string `json:"campaign_status"`

	// Budget information
	TotalBudget     float64 `json:"total_budget"`
	SpentAmount     float64 `json:"spent_amount"`
	RemainingBudget float64 `json:"remaining_budget"`
	UsagePercentage float64 `json:"usage_percentage"`

	// Monthly tracking
	MonthlyLimit     float64 `json:"monthly_limit"`
	MonthlySpent     float64 `json:"monthly_spent"`
	MonthlyRemaining float64 `json:"monthly_remaining"`
	MonthlyUsage     float64 `json:"monthly_usage"`

	// Daily tracking
	DailyLimit     *float64 `json:"daily_limit,omitempty"`
	DailySpent     float64  `json:"daily_spent"`
	DailyRemaining *float64 `json:"daily_remaining,omitempty"`
	DailyUsage     *float64 `json:"daily_usage,omitempty"`

	// Alert status
	AlertThreshold   float64 `json:"alert_threshold"`
	AlertEnabled     bool    `json:"alert_enabled"`
	IsAlertTriggered bool    `json:"is_alert_triggered"`
	AlertMessage     string  `json:"alert_message,omitempty"`

	// Currency
	Currency    string `json:"currency"`
	LastUpdated string `json:"last_updated"`
}

// BudgetListDTO for budget list response
type BudgetListDTO struct {
	ID             uint64 `json:"id"`
	CampaignID     uint64 `json:"campaign_id"`
	CampaignName   string `json:"campaign_name"`
	CampaignStatus string `json:"campaign_status"`
	CampaignType   string `json:"campaign_type"`

	// Budget summary
	TotalBudget     float64 `json:"total_budget"`
	SpentAmount     float64 `json:"spent_amount"`
	RemainingBudget float64 `json:"remaining_budget"`
	UsagePercentage float64 `json:"usage_percentage"`

	// Monthly info
	MonthlyLimit float64 `json:"monthly_limit"`
	MonthlySpent float64 `json:"monthly_spent"`
	MonthlyUsage float64 `json:"monthly_usage"`

	// Alert status
	AlertEnabled     bool `json:"alert_enabled"`
	IsAlertTriggered bool `json:"is_alert_triggered"`

	// Metadata
	Currency  string `json:"currency"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// BudgetSearchDTO for searching budgets
type BudgetSearchDTO struct {
	CampaignID       *uint64  `json:"campaign_id,omitempty"`
	CampaignStatus   string   `json:"campaign_status,omitempty"`
	CampaignType     string   `json:"campaign_type,omitempty"`
	AlertEnabled     *bool    `json:"alert_enabled,omitempty"`
	IsAlertTriggered *bool    `json:"is_alert_triggered,omitempty"`
	UsageMin         *float64 `json:"usage_min,omitempty"`
	UsageMax         *float64 `json:"usage_max,omitempty"`
	Page             int      `json:"page"`
	Size             int      `json:"size"`
}

// BudgetAlertDTO for budget alerts
type BudgetAlertDTO struct {
	ID           uint64  `json:"id"`
	CampaignID   uint64  `json:"campaign_id"`
	CampaignName string  `json:"campaign_name"`
	AlertType    string  `json:"alert_type"`  // "THRESHOLD", "LIMIT_EXCEEDED", "LOW_BALANCE"
	AlertLevel   string  `json:"alert_level"` // "WARNING", "CRITICAL"
	Message      string  `json:"message"`
	CurrentUsage float64 `json:"current_usage"`
	Threshold    float64 `json:"threshold"`
	IsRead       bool    `json:"is_read"`
	CreatedAt    string  `json:"created_at"`
}

// BudgetTopupDTO for budget topup
type BudgetTopupDTO struct {
	CampaignID    uint64  `json:"campaign_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,min=10000"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
	Description   string  `json:"description,omitempty"`
	TermsAccepted bool    `json:"terms_accepted" binding:"required"`
}

// BudgetTopupResponseDTO for budget topup response
type BudgetTopupResponseDTO struct {
	Success        bool    `json:"success"`
	TransactionID  uint64  `json:"transaction_id"`
	CampaignID     uint64  `json:"campaign_id"`
	Amount         float64 `json:"amount"`
	PreviousBudget float64 `json:"previous_budget"`
	NewBudget      float64 `json:"new_budget"`
	Currency       string  `json:"currency"`
	Status         string  `json:"status"`
	Message        string  `json:"message"`
	CreatedAt      string  `json:"created_at"`
}

// BudgetUsageDTO for budget usage tracking
type BudgetUsageDTO struct {
	CampaignID  uint64   `json:"campaign_id"`
	Date        string   `json:"date"`
	Impressions int      `json:"impressions"`
	Clicks      int      `json:"clicks"`
	SpentAmount float64  `json:"spent_amount"`
	CTR         float64  `json:"ctr"`
	CPC         float64  `json:"cpc"`
	DailyLimit  *float64 `json:"daily_limit,omitempty"`
	DailyUsage  *float64 `json:"daily_usage,omitempty"`
}

// BudgetUsageSummaryDTO for budget usage summary
type BudgetUsageSummaryDTO struct {
	CampaignID       uint64           `json:"campaign_id"`
	Period           string           `json:"period"` // "daily", "weekly", "monthly"
	TotalImpressions int              `json:"total_impressions"`
	TotalClicks      int              `json:"total_clicks"`
	TotalSpent       float64          `json:"total_spent"`
	AverageCTR       float64          `json:"average_ctr"`
	AverageCPC       float64          `json:"average_cpc"`
	UsageData        []BudgetUsageDTO `json:"usage_data"`
}