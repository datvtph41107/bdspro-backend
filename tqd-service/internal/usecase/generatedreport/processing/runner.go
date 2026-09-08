package processing

import (
	"context"
	"log/slog"
	"time"
)

/**
 * Runner poll report job đều đặn và gọi Worker.RunOnce.
 *
 * Không cần cron framework cho job chạy liên tục này.
 */
type Runner struct {
	worker           *Worker
	interval         time.Duration
	recoveryInterval time.Duration
	staleAfter       time.Duration
	recoveryBatch    int
	logger           *slog.Logger
}

func NewRunner(worker *Worker, interval time.Duration, logger *slog.Logger) *Runner {
	if interval <= 0 {
		interval = time.Second
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Runner{
		worker:           worker,
		interval:         interval,
		recoveryInterval: time.Minute,
		staleAfter:       5 * time.Minute,
		recoveryBatch:    100,
		logger:           logger,
	}
}

/**
 * Run chạy cho tới khi context bị cancel.
 */
func (r *Runner) Run(ctx context.Context) {
	if r == nil || r.worker == nil {
		return
	}

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	recoveryTicker := time.NewTicker(r.recoveryInterval)
	defer recoveryTicker.Stop()

	r.releaseStale(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-recoveryTicker.C:
			r.releaseStale(ctx)
		case <-ticker.C:
			processed, err := r.worker.RunOnce(ctx)
			if err != nil {
				r.logger.Error("report job failed", "error", err)
				continue
			}
			if processed {
				r.logger.Debug("report job processed")
			}
		}
	}
}

func (r *Runner) releaseStale(ctx context.Context) {
	released, err := r.worker.ReleaseStale(ctx, r.staleAfter, r.recoveryBatch)
	if err != nil {
		r.logger.Error("release stale report jobs failed", "error", err)
		return
	}
	if released > 0 {
		r.logger.Warn("stale report jobs released", "count", released)
	}
}
