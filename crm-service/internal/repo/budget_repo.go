package repo

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
)

type BudgetRepo interface {
	// Budget config management
	CreateBudgetConfig(ctx context.Context, config *domain.BudgetConfig) (*domain.BudgetConfig, error)
	UpdateBudgetConfig(ctx context.Context, id uint64, config *domain.BudgetConfig) (*domain.BudgetConfig, error)
	GetBudgetConfigByCampaignID(ctx context.Context, campaignID uint64) (*domain.BudgetConfig, error)
	DeleteBudgetConfig(ctx context.Context, id uint64) error

	// Budget list and search
	SearchBudgetConfigs(ctx context.Context, organizationID uint64, searchDTO dto.BudgetSearchDTO) ([]*domain.BudgetConfig, int64, error)
	GetBudgetList(ctx context.Context, organizationID uint64, page, limit int) ([]*dto.BudgetListDTO, int64, error)

	// Budget usage tracking
	CreateBudgetUsage(ctx context.Context, usage *domain.BudgetUsage) (*domain.BudgetUsage, error)
	UpdateBudgetUsage(ctx context.Context, id uint64, usage *domain.BudgetUsage) (*domain.BudgetUsage, error)
	GetBudgetUsageByDate(ctx context.Context, campaignID uint64, date string) (*domain.BudgetUsage, error)
	GetBudgetUsageSummary(ctx context.Context, campaignID uint64, startDate, endDate string) (*dto.BudgetUsageSummaryDTO, error)

	// Budget alerts
	CreateBudgetAlert(ctx context.Context, alert *domain.BudgetAlert) (*domain.BudgetAlert, error)
	GetBudgetAlerts(ctx context.Context, campaignID uint64, page, limit int) ([]*domain.BudgetAlert, int64, error)
	MarkAlertAsRead(ctx context.Context, alertID uint64) error
	GetUnreadAlerts(ctx context.Context, organizationID uint64) ([]*domain.BudgetAlert, error)

	// Budget status and calculations
	GetBudgetStatus(ctx context.Context, campaignID uint64) (*dto.BudgetStatusDTO, error)
	CalculateMonthlyUsage(ctx context.Context, campaignID uint64, year, month int) (float64, error)
	CalculateDailyUsage(ctx context.Context, campaignID uint64, date string) (float64, error)
	CheckAlertThreshold(ctx context.Context, campaignID uint64) (bool, error)
}