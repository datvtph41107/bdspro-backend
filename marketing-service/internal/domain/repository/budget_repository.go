package repository

import (
	"context"
	"marketing/internal/domain/model"
	"marketing/internal/dto"
)

type BudgetRepository interface {
	// Budget config management
	CreateBudgetConfig(ctx context.Context, config *model.BudgetConfig) (*model.BudgetConfig, error)
	UpdateBudgetConfig(ctx context.Context, id uint64, config *model.BudgetConfig) (*model.BudgetConfig, error)
	GetBudgetConfigByCampaignID(ctx context.Context, campaignID uint64) (*model.BudgetConfig, error)
	DeleteBudgetConfig(ctx context.Context, id uint64) error
	
	// Budget list and search
	SearchBudgetConfigs(ctx context.Context, organizationID uint64, searchDTO dto.BudgetSearchDTO) ([]*model.BudgetConfig, int64, error)
	GetBudgetList(ctx context.Context, organizationID uint64, page, limit int) ([]*dto.BudgetListDTO, int64, error)
	
	// Budget usage tracking
	CreateBudgetUsage(ctx context.Context, usage *model.BudgetUsage) (*model.BudgetUsage, error)
	UpdateBudgetUsage(ctx context.Context, id uint64, usage *model.BudgetUsage) (*model.BudgetUsage, error)
	GetBudgetUsageByDate(ctx context.Context, campaignID uint64, date string) (*model.BudgetUsage, error)
	GetBudgetUsageSummary(ctx context.Context, campaignID uint64, startDate, endDate string) (*dto.BudgetUsageSummaryDTO, error)
	
	// Budget alerts
	CreateBudgetAlert(ctx context.Context, alert *model.BudgetAlert) (*model.BudgetAlert, error)
	GetBudgetAlerts(ctx context.Context, campaignID uint64, page, limit int) ([]*model.BudgetAlert, int64, error)
	MarkAlertAsRead(ctx context.Context, alertID uint64) error
	GetUnreadAlerts(ctx context.Context, organizationID uint64) ([]*model.BudgetAlert, error)
	
	// Budget status and calculations
	GetBudgetStatus(ctx context.Context, campaignID uint64) (*dto.BudgetStatusDTO, error)
	CalculateMonthlyUsage(ctx context.Context, campaignID uint64, year, month int) (float64, error)
	CalculateDailyUsage(ctx context.Context, campaignID uint64, date string) (float64, error)
	CheckAlertThreshold(ctx context.Context, campaignID uint64) (bool, error)
} 