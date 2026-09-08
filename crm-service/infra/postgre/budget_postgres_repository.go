package postgre

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/repo"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BudgetPostgresRepository struct {
	db *gorm.DB
}

func NewBudgetPostgresRepository(db *gorm.DB) repo.BudgetRepo {
	return &BudgetPostgresRepository{
		db: db,
	}
}

// Budget config management

func (r *BudgetPostgresRepository) CreateBudgetConfig(ctx context.Context, config *domain.BudgetConfig) (*domain.BudgetConfig, error) {
	if err := r.db.WithContext(ctx).Create(config).Error; err != nil {
		return nil, fmt.Errorf("failed to create budget config: %w", err)
	}
	return config, nil
}

func (r *BudgetPostgresRepository) UpdateBudgetConfig(ctx context.Context, id uint64, config *domain.BudgetConfig) (*domain.BudgetConfig, error) {
	if err := r.db.WithContext(ctx).Model(&domain.BudgetConfig{}).Where("id = ?", id).Updates(config).Error; err != nil {
		return nil, fmt.Errorf("failed to update budget config: %w", err)
	}
	return r.GetBudgetConfigByCampaignID(ctx, config.CampaignID)
}

func (r *BudgetPostgresRepository) GetBudgetConfigByCampaignID(ctx context.Context, campaignID uint64) (*domain.BudgetConfig, error) {
	var config domain.BudgetConfig
	if err := r.db.WithContext(ctx).Where("campaign_id = ?", campaignID).First(&config).Error; err != nil {
		return nil, fmt.Errorf("budget config not found: %w", err)
	}
	return &config, nil
}

func (r *BudgetPostgresRepository) DeleteBudgetConfig(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).Delete(&domain.BudgetConfig{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete budget config: %w", err)
	}
	return nil
}

// Budget list and search

func (r *BudgetPostgresRepository) SearchBudgetConfigs(ctx context.Context, organizationID uint64, searchDTO dto.BudgetSearchDTO) ([]*domain.BudgetConfig, int64, error) {
	var configs []*domain.BudgetConfig
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.BudgetConfig{}).Where("organization_id = ?", organizationID)

	// Apply filters
	if searchDTO.CampaignID != nil {
		query = query.Where("campaign_id = ?", *searchDTO.CampaignID)
	}
	if searchDTO.AlertEnabled != nil {
		query = query.Where("alert_enabled = ?", *searchDTO.AlertEnabled)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count budget configs: %w", err)
	}

	// Apply pagination
	offset := searchDTO.Page * searchDTO.Size
	if err := query.Offset(offset).Limit(searchDTO.Size).Find(&configs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search budget configs: %w", err)
	}

	return configs, total, nil
}

func (r *BudgetPostgresRepository) GetBudgetList(ctx context.Context, organizationID uint64, page, size int) ([]*dto.BudgetListDTO, int64, error) {
	var total int64
	var budgetList []*dto.BudgetListDTO

	// Build query with campaign join
	query := r.db.WithContext(ctx).
		Table("budget_configs bc").
		Select(`
			bc.id,
			bc.campaign_id,
			c.name as campaign_name,
			c.status as campaign_status,
			c.type as campaign_type,
			c.budget as total_budget,
			COALESCE(SUM(bu.spent_amount), 0) as spent_amount,
			c.budget - COALESCE(SUM(bu.spent_amount), 0) as remaining_budget,
			CASE 
				WHEN c.budget > 0 THEN (COALESCE(SUM(bu.spent_amount), 0) / c.budget) * 100 
				ELSE 0 
			END as usage_percentage,
			bc.monthly_limit,
			COALESCE(SUM(CASE 
				WHEN DATE_TRUNC('month', bu.date) = DATE_TRUNC('month', CURRENT_DATE) 
				THEN bu.spent_amount 
				ELSE 0 
			END), 0) as monthly_spent,
			CASE 
				WHEN bc.monthly_limit > 0 THEN (COALESCE(SUM(CASE 
					WHEN DATE_TRUNC('month', bu.date) = DATE_TRUNC('month', CURRENT_DATE) 
					THEN bu.spent_amount 
					ELSE 0 
				END), 0) / bc.monthly_limit) * 100 
				ELSE 0 
			END as monthly_usage,
			bc.alert_enabled,
			CASE 
				WHEN bc.monthly_limit > 0 AND (COALESCE(SUM(CASE 
					WHEN DATE_TRUNC('month', bu.date) = DATE_TRUNC('month', CURRENT_DATE) 
					THEN bu.spent_amount 
					ELSE 0 
				END), 0) / bc.monthly_limit) * 100 >= bc.alert_threshold 
				THEN true 
				ELSE false 
			END as is_alert_triggered,
			c.currency,
			bc.created_at,
			bc.updated_at
		`).
		Joins("LEFT JOIN campaigns c ON bc.campaign_id = c.id").
		Joins("LEFT JOIN budget_usages bu ON bc.campaign_id = bu.campaign_id").
		Where("bc.organization_id = ?", organizationID).
		Group("bc.id, c.id")

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count budget list: %w", err)
	}

	// Apply pagination
	offset := page * size
	if err := query.Offset(offset).Limit(size).Find(&budgetList).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get budget list: %w", err)
	}

	return budgetList, total, nil
}

// Budget usage tracking

func (r *BudgetPostgresRepository) CreateBudgetUsage(ctx context.Context, usage *domain.BudgetUsage) (*domain.BudgetUsage, error) {
	if err := r.db.WithContext(ctx).Create(usage).Error; err != nil {
		return nil, fmt.Errorf("failed to create budget usage: %w", err)
	}
	return usage, nil
}

func (r *BudgetPostgresRepository) UpdateBudgetUsage(ctx context.Context, id uint64, usage *domain.BudgetUsage) (*domain.BudgetUsage, error) {
	if err := r.db.WithContext(ctx).Model(&domain.BudgetUsage{}).Where("id = ?", id).Updates(usage).Error; err != nil {
		return nil, fmt.Errorf("failed to update budget usage: %w", err)
	}
	return usage, nil
}

func (r *BudgetPostgresRepository) GetBudgetUsageByDate(ctx context.Context, campaignID uint64, date string) (*domain.BudgetUsage, error) {
	var usage domain.BudgetUsage
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	if err := r.db.WithContext(ctx).
		Where("campaign_id = ? AND DATE(date) = DATE(?)", campaignID, parsedDate).
		First(&usage).Error; err != nil {
		return nil, fmt.Errorf("budget usage not found: %w", err)
	}
	return &usage, nil
}

func (r *BudgetPostgresRepository) GetBudgetUsageSummary(ctx context.Context, campaignID uint64, startDate, endDate string) (*dto.BudgetUsageSummaryDTO, error) {
	var usageData []dto.BudgetUsageDTO
	var summary dto.BudgetUsageSummaryDTO

	// Parse dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	// Get daily usage data
	if err := r.db.WithContext(ctx).
		Table("budget_usages").
		Select(`
			campaign_id,
			DATE(date) as date,
			SUM(impressions) as impressions,
			SUM(clicks) as clicks,
			SUM(spent_amount) as spent_amount,
			CASE 
				WHEN SUM(impressions) > 0 THEN (SUM(clicks)::float / SUM(impressions)) * 100 
				ELSE 0 
			END as ctr,
			CASE 
				WHEN SUM(clicks) > 0 THEN SUM(spent_amount) / SUM(clicks) 
				ELSE 0 
			END as cpc
		`).
		Where("campaign_id = ? AND date BETWEEN ? AND ?", campaignID, start, end).
		Group("campaign_id, DATE(date)").
		Order("date").
		Find(&usageData).Error; err != nil {
		return nil, fmt.Errorf("failed to get usage data: %w", err)
	}

	// Calculate summary
	var totalImpressions, totalClicks int
	var totalSpent float64
	for _, data := range usageData {
		totalImpressions += data.Impressions
		totalClicks += data.Clicks
		totalSpent += data.SpentAmount
	}

	averageCTR := 0.0
	if totalImpressions > 0 {
		averageCTR = float64(totalClicks) / float64(totalImpressions) * 100
	}

	averageCPC := 0.0
	if totalClicks > 0 {
		averageCPC = totalSpent / float64(totalClicks)
	}

	summary = dto.BudgetUsageSummaryDTO{
		CampaignID:       campaignID,
		Period:           "daily",
		TotalImpressions: totalImpressions,
		TotalClicks:      totalClicks,
		TotalSpent:       totalSpent,
		AverageCTR:       averageCTR,
		AverageCPC:       averageCPC,
		UsageData:        usageData,
	}

	return &summary, nil
}

// Budget alerts

func (r *BudgetPostgresRepository) CreateBudgetAlert(ctx context.Context, alert *domain.BudgetAlert) (*domain.BudgetAlert, error) {
	if err := r.db.WithContext(ctx).Create(alert).Error; err != nil {
		return nil, fmt.Errorf("failed to create budget alert: %w", err)
	}
	return alert, nil
}

func (r *BudgetPostgresRepository) GetBudgetAlerts(ctx context.Context, campaignID uint64, page, limit int) ([]*domain.BudgetAlert, int64, error) {
	var alerts []*domain.BudgetAlert
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.BudgetAlert{}).Where("campaign_id = ?", campaignID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count budget alerts: %w", err)
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&alerts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get budget alerts: %w", err)
	}

	return alerts, total, nil
}

func (r *BudgetPostgresRepository) MarkAlertAsRead(ctx context.Context, alertID uint64) error {
	if err := r.db.WithContext(ctx).Model(&domain.BudgetAlert{}).Where("id = ?", alertID).Update("is_read", true).Error; err != nil {
		return fmt.Errorf("failed to mark alert as read: %w", err)
	}
	return nil
}

func (r *BudgetPostgresRepository) GetUnreadAlerts(ctx context.Context, organizationID uint64) ([]*domain.BudgetAlert, error) {
	var alerts []*domain.BudgetAlert
	if err := r.db.WithContext(ctx).Where("organization_id = ? AND is_read = ?", organizationID, false).Find(&alerts).Error; err != nil {
		return nil, fmt.Errorf("failed to get unread alerts: %w", err)
	}
	return alerts, nil
}

// Budget status and calculations

func (r *BudgetPostgresRepository) GetBudgetStatus(ctx context.Context, campaignID uint64) (*dto.BudgetStatusDTO, error) {
	// This method is implemented in usecase layer
	// Repository should focus on data access only
	return nil, fmt.Errorf("get budget status should be implemented in usecase layer")
}

func (r *BudgetPostgresRepository) CalculateMonthlyUsage(ctx context.Context, campaignID uint64, year, month int) (float64, error) {
	var totalSpent float64

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	if err := r.db.WithContext(ctx).
		Model(&domain.BudgetUsage{}).
		Select("COALESCE(SUM(spent_amount), 0)").
		Where("campaign_id = ? AND date BETWEEN ? AND ?", campaignID, startDate, endDate).
		Scan(&totalSpent).Error; err != nil {
		return 0, fmt.Errorf("failed to calculate monthly usage: %w", err)
	}

	return totalSpent, nil
}

func (r *BudgetPostgresRepository) CalculateDailyUsage(ctx context.Context, campaignID uint64, date string) (float64, error) {
	var totalSpent float64

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0, fmt.Errorf("invalid date format: %w", err)
	}

	if err := r.db.WithContext(ctx).
		Model(&domain.BudgetUsage{}).
		Select("COALESCE(SUM(spent_amount), 0)").
		Where("campaign_id = ? AND DATE(date) = DATE(?)", campaignID, parsedDate).
		Scan(&totalSpent).Error; err != nil {
		return 0, fmt.Errorf("failed to calculate daily usage: %w", err)
	}

	return totalSpent, nil
}

func (r *BudgetPostgresRepository) CheckAlertThreshold(ctx context.Context, campaignID uint64) (bool, error) {
	// Get budget config
	config, err := r.GetBudgetConfigByCampaignID(ctx, campaignID)
	if err != nil {
		return false, err
	}

	// Calculate current month usage
	now := time.Now()
	monthlyUsage, err := r.CalculateMonthlyUsage(ctx, campaignID, now.Year(), int(now.Month()))
	if err != nil {
		return false, err
	}

	// Check if usage exceeds threshold
	if config.MonthlyLimit > 0 {
		usagePercentage := (monthlyUsage / config.MonthlyLimit) * 100
		return usagePercentage >= config.AlertThreshold, nil
	}

	return false, nil
}