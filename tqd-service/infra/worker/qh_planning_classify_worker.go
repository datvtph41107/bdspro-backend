package worker

import (
	"context"
	"log"
	"time"

	qh_usecase "tqd/internal/usecase/qh"

	"github.com/spf13/viper"
)

const (
	defaultClassifyInterval = 15 * time.Second
	defaultClassifyBatch    = 5
)

// QHPlanningClassifyWorker định kỳ gọi RunOnce để phân loại tài liệu Pending → Classified.
type QHPlanningClassifyWorker struct {
	usecase qh_usecase.IQHPlanningClassifyUsecase
}

func NewQHPlanningClassifyWorker(usecase qh_usecase.IQHPlanningClassifyUsecase) *QHPlanningClassifyWorker {
	return &QHPlanningClassifyWorker{usecase: usecase}
}

// Start chạy vòng lặp classify đến khi ctx bị cancel. Non-blocking (spawn goroutine).
func (w *QHPlanningClassifyWorker) Start(ctx context.Context) {
	if w == nil || w.usecase == nil {
		log.Println("[QHPlanningClassifyWorker] skip: usecase nil")
		return
	}

	interval := defaultClassifyInterval
	if sec := viper.GetInt("qh_planning_classify.interval_sec"); sec > 0 {
		interval = time.Duration(sec) * time.Second
	}
	batchSize := defaultClassifyBatch
	if n := viper.GetInt("qh_planning_classify.batch_size"); n > 0 {
		batchSize = n
	}
	enabled := true
	if viper.IsSet("qh_planning_classify.enabled") {
		enabled = viper.GetBool("qh_planning_classify.enabled")
	}
	if !enabled {
		log.Println("[QHPlanningClassifyWorker] disabled by config")
		return
	}

	go func() {
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
	}()
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
