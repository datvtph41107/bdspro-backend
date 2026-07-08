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

type CampaignDetailUsecase interface {
	// D.10.2.4 specific methods
	GetCampaignDetail(ctx context.Context, campaignID uint64) (*dto.CampaignDetailDTO, error)
	GetCampaignPerformance(ctx context.Context, campaignID uint64, startDate, endDate string) (*dto.CampaignPerformanceDTO, error)
	GetCampaignHistory(ctx context.Context, campaignID uint64, page, limit int) ([]*dto.CampaignHistoryDTO, int64, error)

	// Campaign actions
	UpdateCampaignStatus(ctx context.Context, statusDTO *dto.CampaignStatusUpdateDTO) error
	ExtendCampaign(ctx context.Context, extendDTO *dto.CampaignExtendDTO) (*dto.PaymentResponseDTO, error)

	// Export and optimization
	ExportCampaignData(ctx context.Context, campaignID uint64, exportType, dateRange string) (*dto.CampaignExportDTO, error)
	GetCampaignOptimization(ctx context.Context, campaignID uint64) (*dto.CampaignOptimizationDTO, error)

	// Performance tracking
	TrackCampaignPerformance(ctx context.Context, campaignID uint64) error
	UpdatePerformanceMetrics(ctx context.Context, campaignID uint64, impressions, clicks int, spent float64) error
}

type campaignDetailUsecase struct {
	campaignRepo  repo.CampaignRepo
	packageRepo   repo.PackageRepo
	paymentClient *client.PaymentClient
}

func NewCampaignDetailUsecase(campaignRepo repo.CampaignRepo, packageRepo repo.PackageRepo, paymentClient *client.PaymentClient) CampaignDetailUsecase {
	return &campaignDetailUsecase{
		campaignRepo:  campaignRepo,
		packageRepo:   packageRepo,
		paymentClient: paymentClient,
	}
}

// GetCampaignDetail gets detailed campaign information
func (u *campaignDetailUsecase) GetCampaignDetail(ctx context.Context, campaignID uint64) (*dto.CampaignDetailDTO, error) {
	// Get campaign
	campaign, err := u.campaignRepo.GetByID(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	// Validate ownership
	organizationID := _utils.GetOrganizationIdFromContext(ctx)
	userID := uint64(123) // TODO: Get from context when available
	if campaign.OrganizationID != organizationID {
		return nil, fmt.Errorf("unauthorized: campaign does not belong to your organization")
	}

	// Get package details if available
	var packageName, packageType string
	if campaign.PackageID != nil {
		pkg, err := u.packageRepo.GetByID(ctx, *campaign.PackageID)
		if err == nil {
			packageName = pkg.Name
			packageType = string(pkg.Type)
		}
	}

	// Calculate performance metrics
	performance, err := u.GetCampaignPerformance(ctx, campaignID, "", "")
	if err != nil {
		// Use default values if performance data not available
		performance = &dto.CampaignPerformanceDTO{
			Summary: dto.CampaignPerformanceSummaryDTO{
				TotalImpressions: 0,
				TotalClicks:      0,
				AverageCTR:       0,
				AverageCPC:       0,
				TotalSpent:       0,
				RemainingBudget:  campaign.Budget,
			},
		}
	}

	// Determine permissions
	canEdit := campaign.CreatedBy == userID
	canPause := campaign.Status == domain.CampaignStatusActive
	canResume := campaign.Status == domain.CampaignStatusPaused
	canExtend := campaign.Status == domain.CampaignStatusActive || campaign.Status == domain.CampaignStatusPaused
	canTopup := campaign.Status == domain.CampaignStatusActive

	// Check if campaign can be extended (within 7 days of end date)
	if canExtend && campaign.EndDate != nil {
		daysSinceEnd := time.Since(*campaign.EndDate).Hours() / 24
		if daysSinceEnd > 7 {
			canExtend = false
		}
	}

	// TODO: Get product details from product service
	productName := campaign.ProductName
	productImage := "product_image.jpg" // Placeholder

	return &dto.CampaignDetailDTO{
		ID:          campaign.ID,
		Name:        campaign.Name,
		Description: campaign.Description,
		Type:        string(campaign.Type),
		Status:      string(campaign.Status),

		// Product information
		ProductID:    getUint64Value(campaign.ProductID),
		ProductName:  productName,
		ProductImage: productImage,
		ProductType:  campaign.ProductType,

		// Package information
		PackageID:   getUint64Value(campaign.PackageID),
		PackageName: packageName,
		PackageType: packageType,

		// Campaign settings
		Budget:    campaign.Budget,
		Duration:  campaign.Duration,
		StartDate: formatTime(campaign.StartDate),
		EndDate:   formatTime(campaign.EndDate),
		Currency:  campaign.Currency,

		// Performance metrics
		Impressions:     performance.Summary.TotalImpressions,
		Clicks:          performance.Summary.TotalClicks,
		CTR:             performance.Summary.AverageCTR,
		CPC:             performance.Summary.AverageCPC,
		SpentAmount:     performance.Summary.TotalSpent,
		RemainingBudget: performance.Summary.RemainingBudget,

		// Wallet information
		WalletBalance: campaign.WalletBalance,
		TransactionID: getUint64Value(campaign.TransactionID),

		// Permissions
		CanEdit:   canEdit,
		CanPause:  canPause,
		CanResume: canResume,
		CanExtend: canExtend,
		CanTopup:  canTopup,

		// Metadata
		CreatedBy:      campaign.CreatedBy,
		OrganizationID: campaign.OrganizationID,
		CreatedAt:      campaign.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      campaign.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetCampaignPerformance gets campaign performance metrics
func (u *campaignDetailUsecase) GetCampaignPerformance(ctx context.Context, campaignID uint64, startDate, endDate string) (*dto.CampaignPerformanceDTO, error) {
	// TODO: Implement performance data retrieval from analytics service
	// For now, return mock data

	// Parse date range
	var start, end time.Time
	var err error

	if startDate == "" {
		start = time.Now().AddDate(0, 0, -30) // Last 30 days
	} else {
		start, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date: %w", err)
		}
	}

	if endDate == "" {
		end = time.Now()
	} else {
		end, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date: %w", err)
		}
	}

	// Generate mock daily data
	var dailyData []dto.CampaignDailyDataDTO
	totalImpressions := 0
	totalClicks := 0
	totalSpent := 0.0

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		impressions := 1000 + (int(d.Unix()) % 500) // Mock data
		clicks := impressions/20 + (int(d.Unix()) % 50)
		spent := float64(impressions) * 0.001 // Mock CPC

		totalImpressions += impressions
		totalClicks += clicks
		totalSpent += spent

		dailyData = append(dailyData, dto.CampaignDailyDataDTO{
			Date:        d.Format("2006-01-02"),
			Impressions: impressions,
			Clicks:      clicks,
			CTR:         float64(clicks) / float64(impressions) * 100,
			CPC:         spent / float64(clicks),
			Spent:       spent,
		})
	}

	// Calculate averages
	var averageCTR, averageCPC float64
	if totalImpressions > 0 {
		averageCTR = float64(totalClicks) / float64(totalImpressions) * 100
	}
	if totalClicks > 0 {
		averageCPC = totalSpent / float64(totalClicks)
	}

	return &dto.CampaignPerformanceDTO{
		CampaignID:  campaignID,
		Date:        time.Now().Format("2006-01-02"),
		Impressions: totalImpressions,
		Clicks:      totalClicks,
		CTR:         averageCTR,
		CPC:         averageCPC,
		Spent:       totalSpent,
		DailyData:   dailyData,
		Summary: dto.CampaignPerformanceSummaryDTO{
			TotalImpressions: totalImpressions,
			TotalClicks:      totalClicks,
			AverageCTR:       averageCTR,
			AverageCPC:       averageCPC,
			TotalSpent:       totalSpent,
			RemainingBudget:  1000000 - totalSpent,                   // Mock remaining budget
			ROI:              float64(totalClicks*1000) / totalSpent, // Mock ROI
		},
	}, nil
}

// GetCampaignHistory gets campaign action history
func (u *campaignDetailUsecase) GetCampaignHistory(ctx context.Context, campaignID uint64, page, limit int) ([]*dto.CampaignHistoryDTO, int64, error) {
	// TODO: Implement history retrieval from audit log service
	// For now, return mock data

	mockHistory := []*dto.CampaignHistoryDTO{
		{
			ID:              1,
			CampaignID:      campaignID,
			Action:          "CREATED",
			Description:     "Campaign created",
			PerformedBy:     123,
			PerformedByUser: "John Doe",
			PerformedAt:     time.Now().AddDate(0, 0, -5).Format("2006-01-02T15:04:05Z07:00"),
		},
		{
			ID:              2,
			CampaignID:      campaignID,
			Action:          "PAYMENT",
			Description:     "Payment processed successfully",
			PerformedBy:     123,
			PerformedByUser: "John Doe",
			PerformedAt:     time.Now().AddDate(0, 0, -5).Format("2006-01-02T15:04:05Z07:00"),
		},
		{
			ID:              3,
			CampaignID:      campaignID,
			Action:          "ACTIVATED",
			Description:     "Campaign activated",
			PerformedBy:     123,
			PerformedByUser: "John Doe",
			PerformedAt:     time.Now().AddDate(0, 0, -4).Format("2006-01-02T15:04:05Z07:00"),
		},
	}

	return mockHistory, int64(len(mockHistory)), nil
}

// UpdateCampaignStatus updates campaign status
func (u *campaignDetailUsecase) UpdateCampaignStatus(ctx context.Context, statusDTO *dto.CampaignStatusUpdateDTO) error {
	campaign, err := u.campaignRepo.GetByID(ctx, statusDTO.CampaignID)
	if err != nil {
		return fmt.Errorf("campaign not found: %w", err)
	}

	// Validate ownership
	organizationID := _utils.GetOrganizationIdFromContext(ctx)
	if campaign.OrganizationID != organizationID {
		return fmt.Errorf("unauthorized: campaign does not belong to your organization")
	}

	// Validate status transition
	oldStatus := campaign.Status
	newStatus := domain.CampaignStatus(statusDTO.NewStatus)

	if !u.isValidStatusTransition(oldStatus, newStatus) {
		return fmt.Errorf("invalid status transition from %s to %s", oldStatus, newStatus)
	}

	// Update status
	campaign.Status = newStatus
	_, err = u.campaignRepo.Update(ctx, campaign.ID, campaign)
	if err != nil {
		return fmt.Errorf("failed to update campaign status: %w", err)
	}

	// TODO: Log action to audit system
	// TODO: Send notifications if needed

	return nil
}

// ExtendCampaign extends campaign duration
func (u *campaignDetailUsecase) ExtendCampaign(ctx context.Context, extendDTO *dto.CampaignExtendDTO) (*dto.PaymentResponseDTO, error) {
	campaign, err := u.campaignRepo.GetByID(ctx, extendDTO.CampaignID)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	// Validate ownership
	organizationID := _utils.GetOrganizationIdFromContext(ctx)
	if campaign.OrganizationID != organizationID {
		return nil, fmt.Errorf("unauthorized: campaign does not belong to your organization")
	}

	// Check if campaign can be extended
	if campaign.EndDate != nil {
		daysSinceEnd := time.Since(*campaign.EndDate).Hours() / 24
		if daysSinceEnd > 7 {
			return nil, fmt.Errorf("campaign cannot be extended after 7 days from end date")
		}
	}

	// Calculate extension cost
	extensionCost := extendDTO.AdditionalBudget
	if extensionCost == 0 {
		// Calculate based on current package price and extension days
		if campaign.PackageID != nil {
			pkg, err := u.packageRepo.GetByID(ctx, *campaign.PackageID)
			if err == nil {
				dailyCost := pkg.Price / float64(pkg.Duration)
				extensionCost = dailyCost * float64(extendDTO.ExtensionDays)
			}
		}
	}

	// Process payment if needed
	if extensionCost > 0 {
		// TODO: Implement payment processing for extension
		// paymentDTO := &dto.PaymentRequestDTO{
		// 	CampaignID:    extendDTO.CampaignID,
		// 	Amount:        extensionCost,
		// 	PaymentMethod: extendDTO.PaymentMethod,
		// 	TermsAccepted: extendDTO.TermsAccepted,
		// }

		// For now, return mock response
		return &dto.PaymentResponseDTO{
			Success:       true,
			TransactionID: 999,
			CampaignID:    extendDTO.CampaignID,
			Amount:        extensionCost,
			Currency:      campaign.Currency,
			Status:        "SUCCESS",
			Message:       "Campaign extended successfully",
			NextStep:      "campaign_details",
		}, nil
	}

	// Update campaign duration
	if campaign.EndDate != nil {
		newEndDate := campaign.EndDate.AddDate(0, 0, extendDTO.ExtensionDays)
		campaign.EndDate = &newEndDate
	}
	campaign.Duration += extendDTO.ExtensionDays
	campaign.Budget += extensionCost

	_, err = u.campaignRepo.Update(ctx, campaign.ID, campaign)
	if err != nil {
		return nil, fmt.Errorf("failed to update campaign: %w", err)
	}

	return &dto.PaymentResponseDTO{
		Success:       true,
		TransactionID: 0,
		CampaignID:    extendDTO.CampaignID,
		Amount:        extensionCost,
		Currency:      campaign.Currency,
		Status:        "SUCCESS",
		Message:       "Campaign extended successfully",
		NextStep:      "campaign_details",
	}, nil
}

// ExportCampaignData exports campaign data
func (u *campaignDetailUsecase) ExportCampaignData(ctx context.Context, campaignID uint64, exportType, dateRange string) (*dto.CampaignExportDTO, error) {
	// TODO: Implement export functionality
	// For now, return mock data

	performance, err := u.GetCampaignPerformance(ctx, campaignID, "", "")
	if err != nil {
		return nil, err
	}

	return &dto.CampaignExportDTO{
		CampaignID:      campaignID,
		CampaignName:    "Sample Campaign",
		ExportType:      exportType,
		DateRange:       dateRange,
		PerformanceData: performance.DailyData,
		Summary:         performance.Summary,
		GeneratedAt:     time.Now().Format("2006-01-02T15:04:05Z07:00"),
		GeneratedBy:     uint64(123), // TODO: Get from context when available
	}, nil
}

// GetCampaignOptimization gets AI optimization suggestions
func (u *campaignDetailUsecase) GetCampaignOptimization(ctx context.Context, campaignID uint64) (*dto.CampaignOptimizationDTO, error) {
	// TODO: Implement AI optimization logic
	// For now, return mock suggestions

	return &dto.CampaignOptimizationDTO{
		CampaignID: campaignID,
		Suggestions: []string{
			"Tối ưu từ khóa quảng cáo để tăng CTR",
			"Điều chỉnh ngân sách theo thời gian cao điểm",
			"Cải thiện hình ảnh sản phẩm",
		},
		Priority:        "MEDIUM",
		EstimatedImpact: "Tăng 15-20% hiệu quả quảng cáo",
		GeneratedAt:     time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// TrackCampaignPerformance tracks campaign performance
func (u *campaignDetailUsecase) TrackCampaignPerformance(ctx context.Context, campaignID uint64) error {
	// TODO: Implement performance tracking logic
	// This would typically involve:
	// 1. Fetching data from external ad platforms (Facebook, Google)
	// 2. Updating local performance metrics
	// 3. Triggering alerts if performance is poor

	return nil
}

// UpdatePerformanceMetrics updates performance metrics
func (u *campaignDetailUsecase) UpdatePerformanceMetrics(ctx context.Context, campaignID uint64, impressions, clicks int, spent float64) error {
	// TODO: Implement metrics update logic
	// This would typically involve:
	// 1. Updating campaign performance table
	// 2. Calculating new CTR, CPC, ROI
	// 3. Storing historical data

	return nil
}

// Helper methods

func (u *campaignDetailUsecase) isValidStatusTransition(oldStatus, newStatus domain.CampaignStatus) bool {
	validTransitions := map[domain.CampaignStatus][]domain.CampaignStatus{
		domain.CampaignStatusDraft: {
			domain.CampaignStatusActive,
			domain.CampaignStatusCancelled,
		},
		domain.CampaignStatusActive: {
			domain.CampaignStatusPaused,
			domain.CampaignStatusCompleted,
			domain.CampaignStatusCancelled,
		},
		domain.CampaignStatusPaused: {
			domain.CampaignStatusActive,
			domain.CampaignStatusCompleted,
			domain.CampaignStatusCancelled,
		},
		domain.CampaignStatusCompleted: {
			// No valid transitions from completed
		},
		domain.CampaignStatusCancelled: {
			// No valid transitions from cancelled
		},
	}

	allowedTransitions, exists := validTransitions[oldStatus]
	if !exists {
		return false
	}

	for _, allowed := range allowedTransitions {
		if allowed == newStatus {
			return true
		}
	}
	return false
}

func getUint64Value(ptr *uint64) uint64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02T15:04:05Z07:00")
}