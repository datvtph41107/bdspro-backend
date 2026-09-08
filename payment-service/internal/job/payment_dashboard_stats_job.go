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

	// Tạo dashboard stats job với table name và function lấy count
	dashboardJob := dashboard.NewDashboardStatsJob(
		db,
		"payment_stats", // table name
		func(ctx context.Context) int64 {
			// Lấy count từ payment usecase (giao dịch đang xử lý)
			count, err := paymentUsecase.CountProcessingTransactions(ctx)
			if err != nil {
				return 0
			}
			return count
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
