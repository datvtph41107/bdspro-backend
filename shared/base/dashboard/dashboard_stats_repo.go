package dashboard

import (
	"context"
	"time"
)

type DashboardStatsRepository interface {
	Create(ctx context.Context, tableName string, stats *DashboardStats) (*DashboardStats, error)
	GetByTime(ctx context.Context, tableName string, time time.Time) (*DashboardStats, error)
	GetByTwoTime(ctx context.Context, tableName string, firstTime time.Time, secondTime time.Time) (*DashboardStats, *DashboardStats, error)
	UpdateOrCreate(ctx context.Context, tableName string, stats *DashboardStats) (*DashboardStats, error)
}
