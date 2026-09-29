package job

import (
	"context"
	"time"

	dashboard "base/dashboard"
	"bdspro/internal/repo"

	"gorm.io/gorm"
)

type PostDashboardStatsJob struct {
	postRepo repo.PostRepo
	job      dashboard.DashboardStatsJob
}

func NewPostDashboardStatsJob(
	postRepo repo.PostRepo,
	db *gorm.DB,
) *PostDashboardStatsJob {
	// Tạo dashboard stats job với table name và function lấy count
	dashboardJob := dashboard.NewDashboardStatsJob(
		db,
		"post_stats", // table name
		func(ctx context.Context) int64 {
			// Lấy count từ post repo
			count, err := postRepo.CountCurrent(ctx)
			if err != nil {
				return 0
			}
			return count
		},
	)

	return &PostDashboardStatsJob{
		postRepo: postRepo,
		job:      dashboardJob,
	}
}

func (j *PostDashboardStatsJob) Run(ctx context.Context) error {
	return j.job.Run(ctx)
}

func (j *PostDashboardStatsJob) RunDailyAtMidnight(ctx context.Context) {
	j.job.RunDailyAtMidnight(ctx)
}
func (j *PostDashboardStatsJob) GetStats(ctx context.Context, firstDate time.Time, secondDate time.Time) (*dashboard.DashboardStats, *dashboard.DashboardStats, error) {
	return j.job.GetStats(ctx, firstDate, secondDate)
}
