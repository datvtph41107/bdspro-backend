package usagecheck

import (
	"context"
	"log/slog"
	"time"
)

/**
 * Runner định kỳ kiểm saved usage với Redis runtime.
 */
type Runner struct {
	service  *Service
	scopes   ScopeStore
	interval time.Duration
	batch    int
	offset   int
	logger   *slog.Logger
}

func NewRunner(
	service *Service,
	scopes ScopeStore,
	interval time.Duration,
	batch int,
	logger *slog.Logger,
) *Runner {
	if interval <= 0 {
		interval = 5 * time.Minute
	}
	if batch <= 0 {
		batch = 200
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Runner{service: service, scopes: scopes, interval: interval, batch: batch, logger: logger}
}

func (r *Runner) Run(ctx context.Context) {
	if r == nil || r.service == nil || r.scopes == nil {
		return
	}

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.runOnce(ctx)
		}
	}
}

func (r *Runner) runOnce(ctx context.Context) {
	scopes, err := r.scopes.ListUsageScopes(ctx, r.offset, r.batch)
	if err != nil {
		r.logger.Error("load usage scopes failed", "error", err)
		return
	}
	if len(scopes) == 0 {
		r.offset = 0
		return
	}

	for _, scope := range scopes {
		result, err := r.service.CheckAndFix(
			ctx,
			scope.Subject,
			scope.Operation,
			scope.PeriodStart,
			scope.PeriodEnd,
		)
		if err != nil {
			r.logger.Error(
				"usage check failed",
				"subject_type", string(scope.Subject.Type),
				"subject_id", scope.Subject.ID,
				"business_operation", string(scope.Operation),
				"error", err,
			)
			continue
		}
		if result.Fixed {
			r.logger.Warn(
				"runtime usage fixed",
				"subject_type", string(scope.Subject.Type),
				"subject_id", scope.Subject.ID,
				"business_operation", string(scope.Operation),
				"saved_used", result.SavedUsed,
				"runtime_used", result.RuntimeUsed,
				"reserved", result.Reserved,
			)
		}
	}

	if len(scopes) < r.batch {
		r.offset = 0
	} else {
		r.offset += len(scopes)
	}
}
