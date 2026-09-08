package fulfillmentworker

import (
	"context"
	"log"
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
			log.Printf("commerce fulfillment claim: %v", err)
		}
		if ok {
			processed = true
			if err := w.service.ProcessClaim(ctx, claim, w.workerID); err != nil && ctx.Err() == nil {
				log.Printf("commerce fulfillment process: %v", err)
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
