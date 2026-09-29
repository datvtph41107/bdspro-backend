package dashboard

import (
	"context"
	"time"
)

type DashboardStatsUsecase interface {
	Create(ctx context.Context, tableName string, stats *DashboardStats) (*DashboardStats, error)
	GetByTime(ctx context.Context, tableName string, time time.Time) (*DashboardStats, error)
	GetByTwoTime(ctx context.Context, tableName string, firstTime time.Time, secondTime time.Time) (*DashboardStats, *DashboardStats, error)
	UpdateOrCreate(ctx context.Context, tableName string, stats *DashboardStats) (*DashboardStats, error)
}

type dashboardStatsUsecase struct {
	repo DashboardStatsRepository
}

func NewDashboardStatsUsecase(repo DashboardStatsRepository) DashboardStatsUsecase {
	return &dashboardStatsUsecase{repo: repo}
}

func (u *dashboardStatsUsecase) Create(ctx context.Context, tableName string, stats *DashboardStats) (*DashboardStats, error) {
	return u.repo.Create(ctx, tableName, stats)
}

func (u *dashboardStatsUsecase) GetByTime(ctx context.Context, tableName string, time time.Time) (*DashboardStats, error) {
	return u.repo.GetByTime(ctx, tableName, time)
}

func (u *dashboardStatsUsecase) GetByTwoTime(ctx context.Context, tableName string, firstTime time.Time, secondTime time.Time) (*DashboardStats, *DashboardStats, error) {
	return u.repo.GetByTwoTime(ctx, tableName, firstTime, secondTime)
}

func (u *dashboardStatsUsecase) UpdateOrCreate(ctx context.Context, tableName string, stats *DashboardStats) (*DashboardStats, error) {
	return u.repo.UpdateOrCreate(ctx, tableName, stats)
}
