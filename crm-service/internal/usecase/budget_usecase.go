package usecase

import (
	_utils "common/utils"
	"context"
	"crm/infra/client"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/repo"
	"fmt"
	"time"
)

type BudgetUsecase interface {
	// Budget config management - D.10.2.5
	CreateBudgetConfig(ctx context.Context, createDTO *dto.BudgetConfigCreateDTO) (*dto.BudgetConfigDTO, error)
	UpdateBudgetConfig(ctx context.Context, id uint64, updateDTO *dto.BudgetConfigUpdateDTO) (*dto.BudgetConfigDTO, error)
	GetBudgetConfig(ctx context.Context, campaignID uint64) (*dto.BudgetConfigDTO, error)
	DeleteBudgetConfig(ctx context.Context, id uint64) error

	// Budget list and search
	GetBudgetList(ctx context.Context, searchDTO dto.BudgetSearchDTO) ([]*dto.BudgetListDTO, int64, error)
	GetBudgetStatus(ctx context.Context, campaignID uint64) (*dto.BudgetStatusDTO, error)

	// Budget topup
	TopupBudget(ctx context.Context, topupDTO *dto.BudgetTopupDTO) (*dto.BudgetTopupResponseDTO, error)

	// Budget usage tracking
	UpdateBudgetUsage(ctx context.Context, campaignID uint64, impressions, clicks int, spent float64) error
	GetBudgetUsageSummary(ctx context.Context, campaignID uint64, startDate, endDate string) (*dto.BudgetUsageSummaryDTO, error)

	// Budget alerts
	GetBudgetAlerts(ctx context.Context, campaignID uint64, page, limit int) ([]*dto.BudgetAlertDTO, int64, error)
	MarkAlertAsRead(ctx context.Context, alertID uint64) error
	CheckAndCreateAlerts(ctx context.Context, campaignID uint64) error
}

type budgetUsecase struct {
	budgetRepo    repo.BudgetRepo
	campaignRepo  repo.CampaignRepo
	paymentClient *client.PaymentClient
}

func NewBudgetUsecase(budgetRepo repo.BudgetRepo, campaignRepo repo.CampaignRepo, paymentClient *client.PaymentClient) BudgetUsecase {
	return &budgetUsecase{
		budgetRepo:    budgetRepo,
		campaignRepo:  campaignRepo,
		paymentClient: paymentClient,
	}
}

// CreateBudgetConfig creates budget configuration
func (u *budgetUsecase) CreateBudgetConfig(ctx context.Context, createDTO *dto.BudgetConfigCreateDTO) (*dto.BudgetConfigDTO, error) {
	// Validate campaign exists and belongs to organization
	campaign, err := u.campaignRepo.GetByID(ctx, createDTO.CampaignID)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	organizationID := _utils.GetOrganizationIdFromContext(ctx)
	if campaign.OrganizationID != organizationID {
		return nil, fmt.Errorf("unauthorized: campaign does not belong to your organization")
	}

	// Check if budget config already exists
	existingConfig, err := u.budgetRepo.GetBudgetConfigByCampaignID(ctx, createDTO.CampaignID)
	if err == nil && existingConfig != nil {
		return nil, fmt.Errorf("budget configuration already exists for this campaign")
	}

	// Create budget config
	config := &domain.BudgetConfig{
		CampaignID:         createDTO.CampaignID,
		OrganizationID:     organizationID,
		MonthlyLimit:       createDTO.MonthlyLimit,
		DailyLimit:         createDTO.DailyLimit,
		TotalLimit:         createDTO.TotalLimit,
		AlertThreshold:     createDTO.AlertThreshold,
		AlertEnabled:       createDTO.AlertEnabled,
		AlertEmail:         createDTO.AlertEmail,
		AlertPhone:         createDTO.AlertPhone,
		AutoTopup:          createDTO.AutoTopup,
		AutoTopupAmount:    createDTO.AutoTopupAmount,
		AutoTopupThreshold: createDTO.AutoTopupThreshold,
		CreatedBy:          uint64(123), // TODO: Get from context
		UpdatedBy:          uint64(123), // TODO: Get from context
	}

	createdConfig, err := u.budgetRepo.CreateBudgetConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create budget config: %w", err)
	}

	return u.toBudgetConfigDTO(createdConfig), nil
}

// UpdateBudgetConfig updates budget configuration
func (u *budgetUsecase) UpdateBudgetConfig(ctx context.Context, id uint64, updateDTO *dto.BudgetConfigUpdateDTO) (*dto.BudgetConfigDTO, error) {
	// Get existing config
	existingConfig, err := u.budgetRepo.GetBudgetConfigByCampaignID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("budget config not found: %w", err)
	}

	// Validate ownership
	organizationID := _utils.GetOrganizationIdFromContext(ctx)
	if existingConfig.OrganizationID != organizationID {
		return nil, fmt.Errorf("unauthorized: budget config does not belong to your organization")
	}

	// Update fields
	if updateDTO.MonthlyLimit != nil {
		existingConfig.MonthlyLimit = *updateDTO.MonthlyLimit
	}
	if updateDTO.DailyLimit != nil {
		existingConfig.DailyLimit = updateDTO.DailyLimit
	}
	if updateDTO.TotalLimit != nil {
		existingConfig.TotalLimit = updateDTO.TotalLimit
	}
	if updateDTO.AlertThreshold != nil {
		existingConfig.AlertThreshold = *updateDTO.AlertThreshold
	}
	if updateDTO.AlertEnabled != nil {
		existingConfig.AlertEnabled = *updateDTO.AlertEnabled
	}
	if updateDTO.AlertEmail != nil {
		existingConfig.AlertEmail = *updateDTO.AlertEmail
	}
	if updateDTO.AlertPhone != nil {
		existingConfig.AlertPhone = *updateDTO.AlertPhone
	}
	if updateDTO.AutoTopup != nil {
		existingConfig.AutoTopup = *updateDTO.AutoTopup
	}
	if updateDTO.AutoTopupAmount != nil {
		existingConfig.AutoTopupAmount = *updateDTO.AutoTopupAmount
	}
	if updateDTO.AutoTopupThreshold != nil {
		existingConfig.AutoTopupThreshold = *updateDTO.AutoTopupThreshold
	}

	updatedConfig, err := u.budgetRepo.UpdateBudgetConfig(ctx, existingConfig.ID, existingConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget config: %w", err)
	}

	return u.toBudgetConfigDTO(updatedConfig), nil
}

// GetBudgetConfig gets budget configuration
func (u *budgetUsecase) GetBudgetConfig(ctx context.Context, campaignID uint64) (*dto.BudgetConfigDTO, error) {
	config, err := u.budgetRepo.GetBudgetConfigByCampaignID(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("budget config not found: %w", err)
	}

	// Validate ownership
	organizationID := _utils.GetOrganizationIdFromContext(ctx)
	if config.OrganizationID != organizationID {
		return nil, fmt.Errorf("unauthorized: budget config does not belong to your organization")
	}

	return u.toBudgetConfigDTO(config), nil
}

// DeleteBudgetConfig deletes budget configuration
func (u *budgetUsecase) DeleteBudgetConfig(ctx context.Context, id uint64) error {
	config, err := u.budgetRepo.GetBudgetConfigByCampaignID(ctx, id)
	if err != nil {
		return fmt.Errorf("budget config not found: %w", err)
	}

	// Validate ownership
	organizationID := _utils.GetOrganizationIdFromContext(ctx)
	if config.OrganizationID != organizationID {
		return fmt.Errorf("unauthorized: budget config does not belong to your organization")
	}

	return u.budgetRepo.DeleteBudgetConfig(ctx, config.ID)
}

// GetBudgetList gets budget list
func (u *budgetUsecase) GetBudgetList(ctx context.Context, searchDTO dto.BudgetSearchDTO) ([]*dto.BudgetListDTO, int64, error) {
	organizationID := _utils.GetOrganizationIdFromContext(ctx)

	budgets, total, err := u.budgetRepo.GetBudgetList(ctx, organizationID, searchDTO.Page, searchDTO.Size)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get budget list: %w", err)
	}

	return budgets, total, nil
}

// GetBudgetStatus gets budget status
func (u *budgetUsecase) GetBudgetStatus(ctx context.Context, campaignID uint64) (*dto.BudgetStatusDTO, error) {
	// Get campaign info
	campaign, err := u.campaignRepo.GetByID(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	// Validate ownership
	organizationID := _utils.GetOrganizationIdFromContext(ctx)
	if campaign.OrganizationID != organizationID {
		return nil, fmt.Errorf("unauthorized: campaign does not belong to your organization")
	}

	// Get budget config
	config, err := u.budgetRepo.GetBudgetConfigByCampaignID(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("budget config not found: %w", err)
	}

	// Calculate usage
	now := time.Now()
	monthlySpent, err := u.budgetRepo.CalculateMonthlyUsage(ctx, campaignID, now.Year(), int(now.Month()))
	if err != nil {
		monthlySpent = 0 // Default to 0 if calculation fails
	}

	dailySpent, err := u.budgetRepo.CalculateDailyUsage(ctx, campaignID, now.Format("2006-01-02"))
	if err != nil {
		dailySpent = 0 // Default to 0 if calculation fails
	}

	// Calculate percentages
	monthlyUsage := 0.0
	if config.MonthlyLimit > 0 {
		monthlyUsage = (monthlySpent / config.MonthlyLimit) * 100
	}

	// Check if alert is triggered
	isAlertTriggered := monthlyUsage >= config.AlertThreshold
	alertMessage := ""
	if isAlertTriggered {
		alertMessage = fmt.Sprintf("Budget usage has reached %.1f%% of monthly limit", monthlyUsage)
	}

	return &dto.BudgetStatusDTO{
		CampaignID:     campaignID,
		CampaignName:   campaign.Name,
		CampaignStatus: string(campaign.Status),

		// Budget information
		TotalBudget:     campaign.Budget,
		SpentAmount:     monthlySpent,
		RemainingBudget: campaign.Budget - monthlySpent,
		UsagePercentage: monthlyUsage,

		// Monthly tracking
		MonthlyLimit:     config.MonthlyLimit,
		MonthlySpent:     monthlySpent,
		MonthlyRemaining: config.MonthlyLimit - monthlySpent,
		MonthlyUsage:     monthlyUsage,

		// Daily tracking
		DailyLimit:     config.DailyLimit,
		DailySpent:     dailySpent,
		DailyRemaining: getDailyRemaining(config.DailyLimit, dailySpent),
		DailyUsage:     getDailyUsage(config.DailyLimit, dailySpent),

		// Alert status
		AlertThreshold:   config.AlertThreshold,
		AlertEnabled:     config.AlertEnabled,
		IsAlertTriggered: isAlertTriggered,
		AlertMessage:     alertMessage,

		// Currency
		Currency:    campaign.Currency,
		LastUpdated: time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// TopupBudget tops up campaign budget
func (u *budgetUsecase) TopupBudget(ctx context.Context, topupDTO *dto.BudgetTopupDTO) (*dto.BudgetTopupResponseDTO, error) {
	// Validate campaign exists and is active
	campaign, err := u.campaignRepo.GetByID(ctx, topupDTO.CampaignID)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	// Check if campaign is expired
	if campaign.EndDate != nil && time.Now().After(*campaign.EndDate) {
		return nil, fmt.Errorf("cannot topup budget for expired campaign")
	}

	// Validate ownership
	organizationID := _utils.GetOrganizationIdFromContext(ctx)
	if campaign.OrganizationID != organizationID {
		return nil, fmt.Errorf("unauthorized: campaign does not belong to your organization")
	}

	// Validate minimum amount
	if topupDTO.Amount < 10000 {
		return nil, fmt.Errorf("minimum topup amount is 10,000 VND")
	}

	// Process payment via payment service
	// TODO: Implement actual payment processing
	// For now, return mock response

	previousBudget := campaign.Budget
	newBudget := previousBudget + topupDTO.Amount

	// Update campaign budget
	campaign.Budget = newBudget
	_, err = u.campaignRepo.Update(ctx, campaign.ID, campaign)
	if err != nil {
		return nil, fmt.Errorf("failed to update campaign budget: %w", err)
	}

	return &dto.BudgetTopupResponseDTO{
		Success:        true,
		TransactionID:  999, // Mock transaction ID
		CampaignID:     topupDTO.CampaignID,
		Amount:         topupDTO.Amount,
		PreviousBudget: previousBudget,
		NewBudget:      newBudget,
		Currency:       campaign.Currency,
		Status:         "SUCCESS",
		Message:        "Budget topup successful",
		CreatedAt:      time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// UpdateBudgetUsage updates budget usage
func (u *budgetUsecase) UpdateBudgetUsage(ctx context.Context, campaignID uint64, impressions, clicks int, spent float64) error {
	// Get or create budget usage for today
	today := time.Now().Format("2006-01-02")
	usage, err := u.budgetRepo.GetBudgetUsageByDate(ctx, campaignID, today)

	if err != nil {
		// Create new usage record
		usage = &domain.BudgetUsage{
			CampaignID:  campaignID,
			Date:        time.Now(),
			Impressions: impressions,
			Clicks:      clicks,
			SpentAmount: spent,
		}

		// Calculate CTR and CPC
		if impressions > 0 {
			usage.CTR = float64(clicks) / float64(impressions) * 100
		}
		if clicks > 0 {
			usage.CPC = spent / float64(clicks)
		}

		_, err = u.budgetRepo.CreateBudgetUsage(ctx, usage)
		if err != nil {
			return fmt.Errorf("failed to create budget usage: %w", err)
		}
	} else {
		// Update existing usage record
		usage.Impressions += impressions
		usage.Clicks += clicks
		usage.SpentAmount += spent

		// Recalculate CTR and CPC
		if usage.Impressions > 0 {
			usage.CTR = float64(usage.Clicks) / float64(usage.Impressions) * 100
		}
		if usage.Clicks > 0 {
			usage.CPC = usage.SpentAmount / float64(usage.Clicks)
		}

		_, err = u.budgetRepo.UpdateBudgetUsage(ctx, usage.ID, usage)
		if err != nil {
			return fmt.Errorf("failed to update budget usage: %w", err)
		}
	}

	// Check for alerts
	return u.CheckAndCreateAlerts(ctx, campaignID)
}

// GetBudgetUsageSummary gets budget usage summary
func (u *budgetUsecase) GetBudgetUsageSummary(ctx context.Context, campaignID uint64, startDate, endDate string) (*dto.BudgetUsageSummaryDTO, error) {
	// Validate campaign ownership
	campaign, err := u.campaignRepo.GetByID(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	organizationID := _utils.GetOrganizationIdFromContext(ctx)
	if campaign.OrganizationID != organizationID {
		return nil, fmt.Errorf("unauthorized: campaign does not belong to your organization")
	}

	return u.budgetRepo.GetBudgetUsageSummary(ctx, campaignID, startDate, endDate)
}

// GetBudgetAlerts gets budget alerts
func (u *budgetUsecase) GetBudgetAlerts(ctx context.Context, campaignID uint64, page, limit int) ([]*dto.BudgetAlertDTO, int64, error) {
	alerts, total, err := u.budgetRepo.GetBudgetAlerts(ctx, campaignID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get budget alerts: %w", err)
	}

	alertDTOs := make([]*dto.BudgetAlertDTO, len(alerts))
	for i, alert := range alerts {
		alertDTOs[i] = &dto.BudgetAlertDTO{
			ID:           alert.ID,
			CampaignID:   alert.CampaignID,
			CampaignName: "Campaign Name", // TODO: Get from campaign
			AlertType:    alert.AlertType,
			AlertLevel:   alert.AlertLevel,
			Message:      alert.Message,
			CurrentUsage: alert.CurrentUsage,
			Threshold:    alert.Threshold,
			IsRead:       alert.IsRead,
			CreatedAt:    alert.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return alertDTOs, total, nil
}

// MarkAlertAsRead marks alert as read
func (u *budgetUsecase) MarkAlertAsRead(ctx context.Context, alertID uint64) error {
	return u.budgetRepo.MarkAlertAsRead(ctx, alertID)
}

// CheckAndCreateAlerts checks and creates budget alerts
func (u *budgetUsecase) CheckAndCreateAlerts(ctx context.Context, campaignID uint64) error {
	// Get budget status
	status, err := u.GetBudgetStatus(ctx, campaignID)
	if err != nil {
		return err
	}

	// Get budget config
	config, err := u.budgetRepo.GetBudgetConfigByCampaignID(ctx, campaignID)
	if err != nil {
		return err
	}

	// Check if alert should be created
	if !config.AlertEnabled {
		return nil
	}

	if status.MonthlyUsage >= config.AlertThreshold {
		// Get organization ID from context
		organizationID := _utils.GetOrganizationIdFromContext(ctx)

		// Create alert
		alert := &domain.BudgetAlert{
			CampaignID:     campaignID,
			OrganizationID: organizationID,
			AlertType:      "THRESHOLD",
			AlertLevel:     "WARNING",
			Message:        fmt.Sprintf("Budget usage has reached %.1f%% of monthly limit", status.MonthlyUsage),
			CurrentUsage:   status.MonthlyUsage,
			Threshold:      config.AlertThreshold,
		}

		_, err = u.budgetRepo.CreateBudgetAlert(ctx, alert)
		if err != nil {
			return fmt.Errorf("failed to create budget alert: %w", err)
		}
	}

	return nil
}

// Helper methods

func (u *budgetUsecase) toBudgetConfigDTO(config *domain.BudgetConfig) *dto.BudgetConfigDTO {
	return &dto.BudgetConfigDTO{
		ID:                 config.ID,
		CampaignID:         config.CampaignID,
		OrganizationID:     config.OrganizationID,
		MonthlyLimit:       config.MonthlyLimit,
		DailyLimit:         config.DailyLimit,
		TotalLimit:         config.TotalLimit,
		AlertThreshold:     config.AlertThreshold,
		AlertEnabled:       config.AlertEnabled,
		AlertEmail:         config.AlertEmail,
		AlertPhone:         config.AlertPhone,
		AutoTopup:          config.AutoTopup,
		AutoTopupAmount:    config.AutoTopupAmount,
		AutoTopupThreshold: config.AutoTopupThreshold,
		CreatedBy:          config.CreatedBy,
		UpdatedBy:          config.UpdatedBy,
		CreatedAt:          config.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:          config.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func getDailyRemaining(dailyLimit *float64, dailySpent float64) *float64 {
	if dailyLimit == nil {
		return nil
	}
	remaining := *dailyLimit - dailySpent
	return &remaining
}

func getDailyUsage(dailyLimit *float64, dailySpent float64) *float64 {
	if dailyLimit == nil || *dailyLimit == 0 {
		return nil
	}
	usage := (dailySpent / *dailyLimit) * 100
	return &usage
}