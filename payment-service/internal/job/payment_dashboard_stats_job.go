package job

import (
	"context"
	"time"

	dashboard "base/dashboard"

	walletuc "payment/internal/usecase/wallet"

	"gorm.io/gorm"
)

type PaymentDashboardStatsJob struct {
	paymentUsecase walletuc.PaymentUsecase
	job            dashboard.DashboardStatsJob
}

func NewPaymentDashboardStatsJob(
	paymentUsecase walletuc.PaymentUsecase,
	db *gorm.DB,
) *PaymentDashboardStatsJob {

	// Preserve technical CountProcessingTransactions failures instead of materializing them as count=0.
	dashboardJob := dashboard.NewDashboardStatsJobWithError(
		db,
		"payment_stats",
		func(ctx context.Context) (int64, error) {
			return paymentUsecase.CountProcessingTransactions(ctx)
		},
	)

	return &PaymentDashboardStatsJob{
		paymentUsecase: paymentUsecase,
		job:            dashboardJob,
	}
}

func (j *PaymentDashboardStatsJob) Run(ctx context.Context) error {
	return j.job.Run(ctx)
}

func (j *PaymentDashboardStatsJob) RunDailyAtMidnight(ctx context.Context) {
	j.job.RunDailyAtMidnight(ctx)
}

func (j *PaymentDashboardStatsJob) GetStats(ctx context.Context, firstDate time.Time, secondDate time.Time) (*dashboard.DashboardStats, *dashboard.DashboardStats, error) {
	return j.job.GetStats(ctx, firstDate, secondDate)
}
