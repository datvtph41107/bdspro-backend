package processing

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

/**
 * Generator tạo file thật cho một generated report.
 *
 * Phần tạo PDF/image có thể thay implementation mà worker flow không đổi.
 */
type Generator interface {
	GenerateReport(ctx context.Context, job Job) (Output, error)
}

/**
 * Worker lấy một pending job, chạy generator rồi cập nhật trạng thái.
 */
type Worker struct {
	store     Store
	generator Generator
	logger    *slog.Logger
	workerID  string
	maxTry    int
	retryWait time.Duration
	now       func() time.Time
}

func NewWorker(
	store Store,
	generator Generator,
	logger *slog.Logger,
	workerID string,
) *Worker {
	if logger == nil {
		logger = slog.Default()
	}
	return &Worker{
		store:     store,
		generator: generator,
		logger:    logger,
		workerID:  workerID,
		maxTry:    3,
		retryWait: 30 * time.Second,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}

/**
 * RunOnce xử lý tối đa một report job.
 *
 * false nghĩa là hiện tại không có pending job, không phải lỗi.
 */
func (w *Worker) RunOnce(ctx context.Context) (bool, error) {
	if w == nil || w.store == nil || w.generator == nil || w.workerID == "" {
		return false, errors.New("report worker is not configured")
	}

	now := w.now()
	job, found, err := w.store.ClaimNextJob(ctx, w.workerID, now)
	if err != nil || !found {
		return found, err
	}

	logger := w.logger.With(
		"report_job_id", job.ID,
		"report_id", job.ReportID,
		"operation_id", job.OperationID,
		"business_operation", string(job.Operation),
		"attempt", job.Attempts,
	)

	output, generateErr := w.generator.GenerateReport(ctx, job)
	if generateErr == nil {
		if err := w.store.CompleteJob(ctx, job, output, w.now()); err != nil {
			logger.Error("complete report job failed", "error", err)
			return true, err
		}
		logger.Info("report job completed")
		return true, nil
	}

	if job.Attempts >= w.maxTry {
		if err := w.store.FailJob(ctx, job, generateErr.Error(), w.now()); err != nil {
			logger.Error("mark report job failed", "error", err, "generate_error", generateErr)
			return true, err
		}
		logger.Error("report job failed", "error", generateErr)
		return true, nil
	}

	availableAt := w.now().Add(w.retryWait * time.Duration(job.Attempts))
	if err := w.store.RetryJob(ctx, job, generateErr.Error(), availableAt); err != nil {
		logger.Error("schedule report retry failed", "error", err, "generate_error", generateErr)
		return true, err
	}

	logger.Warn("report job scheduled for retry", "error", generateErr, "available_at", availableAt)
	return true, nil
}

// ReleaseStale đưa các job bị worker chết giữa chừng về pending để có thể claim lại.
func (w *Worker) ReleaseStale(
	ctx context.Context,
	staleAfter time.Duration,
	limit int,
) (int, error) {
	if w == nil || w.store == nil {
		return 0, errors.New("report worker is not configured")
	}
	if staleAfter <= 0 {
		staleAfter = 5 * time.Minute
	}
	if limit <= 0 {
		limit = 100
	}
	return w.store.ReleaseStaleJobs(ctx, w.now().Add(-staleAfter), limit)
}
