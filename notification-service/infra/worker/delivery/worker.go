package delivery

import (
	"context"
	"errors"
	"log"
	"time"

	usecase "notification/internal/usecase/delivery"
)

type Worker struct {
	id      string
	service *usecase.Service
	poll    time.Duration
}

func New(id string, service *usecase.Service, poll time.Duration) *Worker {
	if poll <= 0 {
		poll = time.Second
	}
	return &Worker{id: id, service: service, poll: poll}
}
func (w *Worker) Run(ctx context.Context) error {
	if w == nil || w.service == nil {
		return errors.New("delivery worker missing service")
	}
	t := time.NewTicker(w.poll)
	defer t.Stop()
	for {
		processed, err := w.service.ProcessOne(ctx, w.id)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			// Durable intent state owns business retry. Infrastructure failures
			// remain visible instead of being silently swallowed.
			log.Printf("notification delivery worker %s: %v", w.id, err)
		}
		if processed {
			continue
		}
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
		}
	}
}
