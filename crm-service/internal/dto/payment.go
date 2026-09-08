package dto

// PaymentSummaryDTO for payment summary display
type PaymentSummaryDTO struct {
	PackageID       uint64  `json:"package_id"`
	PackageName     string  `json:"package_name"`
	PackageType     string  `json:"package_type"`
	ProductID       uint64  `json:"product_id"`
	ProductName     string  `json:"product_name"`
	ProductImage    string  `json:"product_image"`
	Duration        int     `json:"duration"`
	TotalCost       float64 `json:"total_cost"`
	Currency        string  `json:"currency"`
	CustomBudget    *float64 `json:"custom_budget,omitempty"`
	CustomDuration  *int    `json:"custom_duration,omitempty"`
	EstimatedReach  *int    `json:"estimated_reach,omitempty"`
	EstimatedClicks *int    `json:"estimated_clicks,omitempty"`
}

// PaymentRequestDTO for payment processing
type PaymentRequestDTO struct {
	CampaignID      uint64  `json:"campaign_id" binding:"required"`
	PackageID       uint64  `json:"package_id" binding:"required"`
	Amount          float64 `json:"amount" binding:"required,min=0"`
	PaymentMethod   string  `json:"payment_method" binding:"required"`
	WalletID        string  `json:"wallet_id"`
	ProductID       uint64  `json:"product_id"`
	ProductName     string  `json:"product_name"`
	ProductType     string  `json:"product_type"`
	CustomBudget    *float64 `json:"custom_budget,omitempty"`
	CustomDuration  *int    `json:"custom_duration,omitempty"`
	RequireInvoice  bool    `json:"require_invoice"`
	InvoiceEmail    string  `json:"invoice_email,omitempty"`
	TermsAccepted   bool    `json:"terms_accepted" binding:"required"`
}

// PaymentResponseDTO for payment result
type PaymentResponseDTO struct {
	Success         bool    `json:"success"`
	TransactionID   uint64  `json:"transaction_id"`
	CampaignID      uint64  `json:"campaign_id"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	PaymentMethod   string  `json:"payment_method"`
	Status          string  `json:"status"`
	TransactionCode string  `json:"transaction_code"`
	WalletBalance   float64 `json:"wallet_balance"`
	PreviousBalance float64 `json:"previous_balance"`
	CreatedAt       string  `json:"created_at"`
	Message         string  `json:"message"`
	NextStep        string  `json:"next_step"` // "campaign_details" or "dashboard"
}

// WalletBalanceDTO for wallet balance check
type WalletBalanceDTO struct {
	WalletID        string  `json:"wallet_id"`
	Balance         float64 `json:"balance"`
	Currency        string  `json:"currency"`
	IsSufficient    bool    `json:"is_sufficient"`
	RequiredAmount  float64 `json:"required_amount"`
	Shortfall       float64 `json:"shortfall"`
	LastUpdated     string  `json:"last_updated"`
}

// PaymentMethodDTO for available payment methods
type PaymentMethodDTO struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	IsDefault   bool   `json:"is_default"`
	Icon        string `json:"icon"`
}

// PaymentValidationDTO for payment validation
type PaymentValidationDTO struct {
	IsValid         bool     `json:"is_valid"`
	Errors          []string `json:"errors,omitempty"`
	Warnings        []string `json:"warnings,omitempty"`
	WalletBalance   *WalletBalanceDTO `json:"wallet_balance,omitempty"`
	PaymentMethods  []PaymentMethodDTO `json:"payment_methods,omitempty"`
	CanProceed      bool     `json:"can_proceed"`
	RequiresTopup   bool     `json:"requires_topup"`
	TopupAmount     float64  `json:"topup_amount,omitempty"`
}

// CampaignPaymentDTO for campaign payment info
type CampaignPaymentDTO struct {
	CampaignID      uint64  `json:"campaign_id"`
	PackageID       uint64  `json:"package_id"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	PaymentStatus   string  `json:"payment_status"`
	TransactionID   *uint64 `json:"transaction_id,omitempty"`
	PaymentMethod   string  `json:"payment_method"`
	PaidAt          *string `json:"paid_at,omitempty"`
	InvoiceNumber   *string `json:"invoice_number,omitempty"`
	RequireInvoice  bool    `json:"require_invoice"`
	InvoiceEmail    *string `json:"invoice_email,omitempty"`
}

// PaymentHistoryDTO for payment history
type PaymentHistoryDTO struct {
	ID              uint64  `json:"id"`
	CampaignID      uint64  `json:"campaign_id"`
	PackageID       uint64  `json:"package_id"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	PaymentMethod   string  `json:"payment_method"`
	Status          string  `json:"status"`
	TransactionCode string  `json:"transaction_code"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	CampaignName    string  `json:"campaign_name"`
	PackageName     string  `json:"package_name"`
} 