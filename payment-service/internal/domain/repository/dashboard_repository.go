package repository

import (
	"context"
	"payment/internal/domain/entity"
	"time"
)

type DashboardMetricRepository interface {
	CreateOrUpdate(ctx context.Context, metric *entity.DashboardMetric) (*entity.DashboardMetric, error)
	GetDashboardMetrics(ctx context.Context, fromDate time.Time, toDate time.Time) ([]*entity.DashboardMetric, error)
}
