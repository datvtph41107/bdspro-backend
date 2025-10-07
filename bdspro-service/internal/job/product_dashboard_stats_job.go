package job

import (
	"context"
	"time"

	dashboard "base/dashboard"
	"bdspro/internal/repo"

	"gorm.io/gorm"
)

type ProductDashboardStatsJob struct {
	productRepo repo.ProductRepo
	job         dashboard.DashboardStatsJob
}

func NewProductDashboardStatsJob(
	productRepo repo.ProductRepo,
	db *gorm.DB,
) *ProductDashboardStatsJob {
	// Tạo dashboard stats job với table name và function lấy count
	dashboardJob := dashboard.NewDashboardStatsJob(
		db,
		"product_stats", // table name
		func() int64 {
			// Lấy count từ product repo
			count, err := productRepo.CountCurrent(context.Background())
			if err != nil {
				return 0
			}
			return count
		},
	)

	return &ProductDashboardStatsJob{
		productRepo: productRepo,
		job:         dashboardJob,
	}
}

func (j *ProductDashboardStatsJob) Run(ctx context.Context) error {
	return j.job.Run(ctx)
}

func (j *ProductDashboardStatsJob) RunDailyAtMidnight(ctx context.Context) {
	j.job.RunDailyAtMidnight(ctx)
}
func (j *ProductDashboardStatsJob) GetStats(ctx context.Context, firstDate time.Time, secondDate time.Time) (*dashboard.DashboardStats, *dashboard.DashboardStats, error) {
	return j.job.GetStats(ctx, firstDate, secondDate)
}
