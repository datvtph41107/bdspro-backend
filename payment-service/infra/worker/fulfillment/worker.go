package fulfillmentworker

import (
	"common/logging"
	"context"
	"log/slog"
	"payment/internal/usecase/fulfillment"
	"time"
)

type Worker struct {
	service     *fulfillment.Service
	store       fulfillment.Store
	workerID    string
	lease, poll time.Duration
	now         fulfillment.Clock
}

func New(s *fulfillment.Service, store fulfillment.Store, id string, lease, poll time.Duration, now fulfillment.Clock) *Worker {
	return &Worker{service: s, store: store, workerID: id, lease: lease, poll: poll, now: now}
}
func (w *Worker) Run(ctx context.Context) {
	t := time.NewTicker(w.poll)
	defer t.Stop()
	for {
		processed := false
		claim, ok, err := w.store.ClaimFulfillment(ctx, w.workerID, w.now().UTC(), w.lease)
		if err != nil && ctx.Err() == nil {
			logging.WithComponent(ctx, "fulfillment").Error("claim fulfillment failed", slog.Any("error", err))
		}
		if ok {
			processed = true
			if err := w.service.ProcessClaim(ctx, claim, w.workerID); err != nil && ctx.Err() == nil {
				logging.WithComponent(ctx, "fulfillment").Error("process fulfillment claim failed", slog.Any("error", err))
			}
		}
		if processed {
			continue
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}
