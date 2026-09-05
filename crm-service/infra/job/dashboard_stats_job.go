package job

import (
	"context"
	"time"

	dashboard "base/dashboard"

	"crm/internal/usecase"

	"gorm.io/gorm"
)

type DashboardStatsJob struct {
	appointmentUsecase usecase.AppointmentUsecase
	job                dashboard.DashboardStatsJob
}

func NewDashboardStatsJob(
	appointmentUsecase usecase.AppointmentUsecase,
	db *gorm.DB,
) *DashboardStatsJob {

	// Tạo dashboard stats job với table name và function lấy count
	dashboardJob := dashboard.NewDashboardStatsJob(
		db,
		"appointment_stats", // table name
		func(ctx context.Context) int64 {
			// Lấy count từ appointment usecase
			count, err := appointmentUsecase.CountCurrent(ctx)
			if err != nil {
				return 0
			}
			return count
		},
	)

	return &DashboardStatsJob{
		appointmentUsecase: appointmentUsecase,
		job:                dashboardJob,
	}
}

func (j *DashboardStatsJob) Run(ctx context.Context) error {
	return j.job.Run(ctx)
}

func (j *DashboardStatsJob) RunDailyAtMidnight(ctx context.Context) {
	j.job.RunDailyAtMidnight(ctx)
}

func (j *DashboardStatsJob) GetStats(ctx context.Context, firstDate time.Time, secondDate time.Time) (*dashboard.DashboardStats, *dashboard.DashboardStats, error) {
	return j.job.GetStats(ctx, firstDate, secondDate)
}
