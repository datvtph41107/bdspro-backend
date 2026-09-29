package worker

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	qh_usecase "tqd/internal/usecase/qh"
)

const (
	defaultClassifyInterval = 15 * time.Second
	defaultClassifyBatch    = 5
)

// QHPlanningClassifyWorker định kỳ gọi RunOnce để phân loại tài liệu Pending → Classified.
type ClassifyWorkerConfig struct {
	Enabled   bool
	Interval  time.Duration
	BatchSize int
}

type QHPlanningClassifyWorker struct {
	usecase qh_usecase.IQHPlanningClassifyUsecase
	config  ClassifyWorkerConfig
}

func NewQHPlanningClassifyWorker(usecase qh_usecase.IQHPlanningClassifyUsecase, cfg ClassifyWorkerConfig) *QHPlanningClassifyWorker {
	if cfg.Interval <= 0 {
		cfg.Interval = defaultClassifyInterval
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = defaultClassifyBatch
	}
	return &QHPlanningClassifyWorker{usecase: usecase, config: cfg}
}

// Run chạy vòng lặp classify đến khi ctx bị cancel. Run là blocking:
// composition root sở hữu goroutine và quyết định khi nào phải Wait.
func (w *QHPlanningClassifyWorker) Close() error {
	if w == nil || w.usecase == nil {
		return nil
	}
	return w.usecase.Close()
}

func (w *QHPlanningClassifyWorker) Run(ctx context.Context) {
	if w == nil || w.usecase == nil {
		slog.WarnContext(ctx, strings.TrimSuffix(fmt.Sprintln("[QHPlanningClassifyWorker] skip: usecase nil"), "\n"))
		return
	}

	interval := w.config.Interval
	batchSize := w.config.BatchSize
	if !w.config.Enabled {
		slog.WarnContext(ctx, strings.TrimSuffix(fmt.Sprintln("[QHPlanningClassifyWorker] disabled by config"), "\n"))
		return
	}
	slog.InfoContext(ctx, fmt.Sprintf("[QHPlanningClassifyWorker] started interval=%s batch=%d", interval, batchSize))
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Chạy ngay 1 lần khi start, không chờ tick đầu.
	w.runOnce(ctx, batchSize)

	for {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, strings.TrimSuffix(fmt.Sprintln("[QHPlanningClassifyWorker] stopped"), "\n"))
			return
		case <-ticker.C:
			w.runOnce(ctx, batchSize)
		}
	}
}

func (w *QHPlanningClassifyWorker) runOnce(ctx context.Context, batchSize int) {
	n, err := w.usecase.RunOnce(ctx, batchSize)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[QHPlanningClassifyWorker] RunOnce lỗi: %v", err))
		return
	}
	if n > 0 {
		slog.InfoContext(ctx, fmt.Sprintf("[QHPlanningClassifyWorker] đã xử lý %d tài liệu", n))
	}
}
