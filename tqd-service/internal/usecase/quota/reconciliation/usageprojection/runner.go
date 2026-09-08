package usageprojection

import (
	"context"
	"errors"
	"fmt"
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
			if err := r.runOnce(ctx); err != nil {
				r.logger.Error("reconcile runtime usage failed", "error", err)
			}
		}
	}
}

/**
 * ReconcileAll đồng bộ toàn bộ counter còn trong retention trước khi mở traffic.
 */
func (r *Runner) ReconcileAll(ctx context.Context) error {
	if r == nil || r.service == nil || r.scopes == nil {
		return errors.New("usage projection runner is not configured")
	}
	r.offset = 0
	for {
		if err := r.runOnce(ctx); err != nil {
			return err
		}
		if r.offset == 0 {
			return nil
		}
	}
}

func (r *Runner) runOnce(ctx context.Context) error {
	scopes, err := r.scopes.ListUsageScopes(ctx, r.offset, r.batch)
	if err != nil {
		return fmt.Errorf("load usage scopes: %w", err)
	}
	if len(scopes) == 0 {
		r.offset = 0
		return nil
	}

	for _, scope := range scopes {
		result, err := r.service.CheckAndFix(
			ctx,
			scope.Subject,
			scope.MeterCode,
			scope.PeriodStart,
			scope.PeriodEnd,
		)
		if err != nil {
			return fmt.Errorf(
				"reconcile subject %s:%s meter %s: %w",
				scope.Subject.Type,
				scope.Subject.ID,
				scope.MeterCode,
				err,
			)
		}
		if result.Fixed {
			r.logger.Warn(
				"runtime usage fixed",
				"subject_type", string(scope.Subject.Type),
				"subject_id", scope.Subject.ID,
				"meter_code", string(scope.MeterCode),
				"period_start", scope.PeriodStart.UTC().Format(time.RFC3339),
				"period_end", scope.PeriodEnd.UTC().Format(time.RFC3339),
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
	return nil
}
