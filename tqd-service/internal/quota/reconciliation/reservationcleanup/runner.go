package reservationcleanup

import (
	"context"
	"log/slog"
	"time"
)

/**
 * Runner xử lý reservation quá hạn theo batch nhỏ.
 */
type Runner struct {
	service  *Service
	interval time.Duration
	batch    int
	logger   *slog.Logger
}

func NewRunner(service *Service, interval time.Duration, batch int, logger *slog.Logger) *Runner {
	if interval <= 0 {
		interval = time.Minute
	}
	if batch <= 0 {
		batch = 100
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Runner{service: service, interval: interval, batch: batch, logger: logger}
}

func (r *Runner) Run(ctx context.Context) {
	if r == nil || r.service == nil {
		return
	}

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			result, err := r.service.RunOnce(ctx, now.UTC(), r.batch)
			if err != nil {
				r.logger.Error("quota cleanup failed", "error", err)
				continue
			}
			for _, action := range result.Actions {
				r.logger.Info(
					"quota reservation repaired",
					"reservation_id", action.ReservationID,
					"action", action.Outcome,
					"durable_usage_found", action.DurableUsageFound,
				)
			}

			if result.Committed+result.Canceled+result.Skipped > 0 {
				r.logger.Info(
					"quota cleanup completed",
					"committed", result.Committed,
					"canceled", result.Canceled,
					"skipped", result.Skipped,
				)
			}
		}
	}
}
