package worker

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	rabbit "payment/infra/broker/rabbitmq"
	"payment/internal/usecase/outbox"
)

// OutboxPublisher transfers durable Payment events to RabbitMQ. It does not
// own a database or broker connection; main.go creates and closes those
// resources so the complete process remains visible from one composition root.
type OutboxPublisher struct {
	service  *outbox.Service
	workerID string
	lease    time.Duration
	poll     time.Duration
}

func NewOutboxPublisher(service *outbox.Service, workerID string, lease, poll time.Duration) *OutboxPublisher {
	return &OutboxPublisher{service: service, workerID: workerID, lease: lease, poll: poll}
}

func (w *OutboxPublisher) Run(ctx context.Context) error {
	if w == nil || w.service == nil {
		return errors.New("payment outbox worker is not configured")
	}
	ticker := time.NewTicker(w.poll)
	defer ticker.Stop()

	for {
		processed, err := w.service.RunOne(ctx, w.workerID, w.lease)
		if err != nil && !errors.Is(err, context.Canceled) {
			// A closed Rabbit channel cannot recover inside the current process.
			// Returning restarts the whole Payment service; the durable outbox row
			// remains available for the next process instance.
			if errors.Is(err, rabbit.ErrPublisherUnavailable) {
				return err
			}
			log.Printf("service=payment-service component=outbox status=retry error=%q", err)
		}
		if ctx.Err() != nil {
			return nil
		}
		if processed {
			continue
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func ProcessID(component string) string {
	host, _ := os.Hostname()
	if host == "" {
		host = "unknown"
	}
	return component + "-" + host + "-" + strconv.Itoa(os.Getpid())
}
