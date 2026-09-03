package walletpostgres

import (
	"context"
	"errors"
	"payment/internal/domain/wallet"
	walletuc "payment/internal/usecase/wallet"
	"time"

	"gorm.io/gorm"
)

type DashboardMetricPostgresRepository struct {
	db *gorm.DB
}

func NewDashboardMetricPostgresRepository(db *gorm.DB) walletuc.DashboardMetricRepository {
	return &DashboardMetricPostgresRepository{db: db}
}

func (r *DashboardMetricPostgresRepository) CreateOrUpdate(ctx context.Context, metric *wallet.DashboardMetric) (*wallet.DashboardMetric, error) {
	if err := dbFromContext(ctx, r.db).Where("time = ?", metric.Time).First(&wallet.DashboardMetric{}).Error; err != nil {
		if err := dbFromContext(ctx, r.db).Create(metric).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			return nil, err
		}
		return metric, nil
	}
	if err := dbFromContext(ctx, r.db).Model(&wallet.DashboardMetric{}).Where("time = ?", metric.Time).Updates(metric).Error; err != nil {
		return nil, err
	}
	return metric, nil
}

func (r *DashboardMetricPostgresRepository) GetDashboardMetrics(ctx context.Context, fromDate time.Time, toDate time.Time) ([]*wallet.DashboardMetric, error) {
	var metrics []*wallet.DashboardMetric
	if err := dbFromContext(ctx, r.db).Where("time BETWEEN ? AND ?", fromDate, toDate).Find(&metrics).Error; err != nil {
		return nil, err
	}
	return metrics, nil
}
