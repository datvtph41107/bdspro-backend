package job

import (
	"context"
	"time"

	dashboard "base/dashboard"

	"user/internal/usecases"

	"gorm.io/gorm"
)

type UserDashboardStatsJob struct {
	userUsecase usecases.UserUsecase
	job         dashboard.DashboardStatsJob
}

func NewUserDashboardStatsJob(
	userUsecase usecases.UserUsecase,
	db *gorm.DB,
) *UserDashboardStatsJob {

	// Tạo dashboard stats job với table name và function lấy count
	dashboardJob := dashboard.NewDashboardStatsJob(
		db,
		"user_stats", // table name
		func(ctx context.Context) int64 {
			// Lấy count từ user usecase
			count, err := userUsecase.CountCurrent(ctx)
			if err != nil {
				return 0
			}
			return count
		},
	)

	return &UserDashboardStatsJob{
		userUsecase: userUsecase,
		job:         dashboardJob,
	}
}

func (j *UserDashboardStatsJob) Run(ctx context.Context) error {
	return j.job.Run(ctx)
}

func (j *UserDashboardStatsJob) RunDailyAtMidnight(ctx context.Context) {
	j.job.RunDailyAtMidnight(ctx)
}

func (j *UserDashboardStatsJob) GetStats(ctx context.Context, firstDate time.Time, secondDate time.Time) (*dashboard.DashboardStats, *dashboard.DashboardStats, error) {
	return j.job.GetStats(ctx, firstDate, secondDate)
}
