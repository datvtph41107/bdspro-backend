package job

import (
	"context"
	"payment/internal/domain/wallet"
	walletuc "payment/internal/usecase/wallet"
	"time"

	"github.com/robfig/cron/v3"
)

type TransactionDashboardJob struct {
	paymentUsecase walletuc.PaymentUsecase
	cronJob        *cron.Cron
}

func NewTransactionDashboardJob(
	paymentUsecase walletuc.PaymentUsecase,
) *TransactionDashboardJob {
	return &TransactionDashboardJob{
		paymentUsecase: paymentUsecase,
		cronJob:        cron.New(),
	}
}
func (j *TransactionDashboardJob) Run(ctx context.Context) error {
	_, err := j.cronJob.AddFunc("0 0 * * *", func() {
		processingCount, err := j.paymentUsecase.CountProcessingTransactions(ctx)
		if err != nil {
			return
		}
		completedCount, err := j.paymentUsecase.CountCompletedTransactions(ctx)
		if err != nil {
			return
		}
		err = j.paymentUsecase.CreateOrUpdateDashboardMetric(ctx, &wallet.DashboardMetric{
			Time:            time.Now(),
			ProcessingCount: uint64(processingCount),
			CompletedCount:  uint64(completedCount),
		})
		if err != nil {
			return
		}
	})
	if err != nil {
		return err
	}
	j.cronJob.Start()
	<-ctx.Done()
	stopCtx := j.cronJob.Stop()
	<-stopCtx.Done()
	return nil
}
