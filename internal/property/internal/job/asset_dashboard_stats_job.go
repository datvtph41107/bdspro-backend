package job

import (
	dashboard "base/dashboard"
	"bdspro/internal/repo"
	"context"
	"time"

	"gorm.io/gorm"
)

type AssetDashboardStatsJob struct {
	assetRepo repo.AssetRepo
	job       dashboard.DashboardStatsJob
}

func NewAssetDashboardStatsJob(
	assetRepo repo.AssetRepo,
	db *gorm.DB,
) *AssetDashboardStatsJob {
	// Tạo dashboard stats job với table name và function lấy count
	dashboardJob := dashboard.NewDashboardStatsJob(
		db,
		"asset_stats", // table name
		func(ctx context.Context) int64 {
			// Lấy count từ asset repo
			// count, err := assetRepo.CountCurrent(ctx)
			// if err != nil {
			// 	return 0
			// }
			return 0
		},
	)

	return &AssetDashboardStatsJob{
		assetRepo: assetRepo,
		job:       dashboardJob,
	}
}

func (j *AssetDashboardStatsJob) Run(ctx context.Context) error {
	return j.job.Run(ctx)
}

func (j *AssetDashboardStatsJob) RunDailyAtMidnight(ctx context.Context) {
	j.job.RunDailyAtMidnight(ctx)
}

func (j *AssetDashboardStatsJob) GetStats(ctx context.Context, firstDate time.Time, secondDate time.Time) (*dashboard.DashboardStats, *dashboard.DashboardStats, error) {
	return j.job.GetStats(ctx, firstDate, secondDate)
}
