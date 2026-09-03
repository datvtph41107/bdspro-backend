package worker

import (
	"context"
	"log"
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
		log.Println("[QHPlanningClassifyWorker] skip: usecase nil")
		return
	}

	interval := w.config.Interval
	batchSize := w.config.BatchSize
	if !w.config.Enabled {
		log.Println("[QHPlanningClassifyWorker] disabled by config")
		return
	}

	log.Printf("[QHPlanningClassifyWorker] started interval=%s batch=%d", interval, batchSize)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Chạy ngay 1 lần khi start, không chờ tick đầu.
	w.runOnce(ctx, batchSize)

	for {
		select {
		case <-ctx.Done():
			log.Println("[QHPlanningClassifyWorker] stopped")
			return
		case <-ticker.C:
			w.runOnce(ctx, batchSize)
		}
	}
}

func (w *QHPlanningClassifyWorker) runOnce(ctx context.Context, batchSize int) {
	n, err := w.usecase.RunOnce(ctx, batchSize)
	if err != nil {
		log.Printf("[QHPlanningClassifyWorker] RunOnce lỗi: %v", err)
		return
	}
	if n > 0 {
		log.Printf("[QHPlanningClassifyWorker] đã xử lý %d tài liệu", n)
	}
}
