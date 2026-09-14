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
// own a database or broker connection; the process composition root supplies
// those resources through the supervisor factory.
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
			// Broker connection/channel loss is a process-local infrastructure
			// failure. Return it to OutboxSupervisor so the composition-root
			// factory can recreate those resources without terminating Payment.
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

type OutboxServiceFactory func() (*outbox.Service, func(), error)

// OutboxSupervisor preserves the Payment process across a broker restart while
// keeping AMQP construction/close ownership in the composition root. The
// durable outbox service remains authoritative for message retry timing and
// idempotency; this actor only reconstructs failed transport resources.
type OutboxSupervisor struct {
	factory        OutboxServiceFactory
	workerID       string
	lease          time.Duration
	poll           time.Duration
	reconnectDelay time.Duration
}

func NewOutboxSupervisor(
	factory OutboxServiceFactory,
	workerID string,
	lease, poll, reconnectDelay time.Duration,
) *OutboxSupervisor {
	return &OutboxSupervisor{
		factory:        factory,
		workerID:       workerID,
		lease:          lease,
		poll:           poll,
		reconnectDelay: reconnectDelay,
	}
}

func (s *OutboxSupervisor) Run(ctx context.Context) error {
	if s == nil || s.factory == nil || s.workerID == "" || s.lease <= 0 || s.poll <= 0 || s.reconnectDelay <= 0 {
		return errors.New("payment outbox supervisor is not configured")
	}

	established := false
	for {
		if ctx.Err() != nil {
			return nil
		}

		service, closeResources, err := s.factory()
		if err != nil {
			// Preserve startup fail-fast semantics. Reconnection retries become
			// valid only after this process has established its publisher once.
			if !established {
				return err
			}
			log.Printf("service=payment-service component=outbox status=reconnecting retry_in=%s error=%q", s.reconnectDelay, err)
			if !waitForOutboxReconnect(ctx, s.reconnectDelay) {
				return nil
			}
			continue
		}
		if service == nil {
			if closeResources != nil {
				closeResources()
			}
			return errors.New("payment outbox supervisor factory returned nil service")
		}
		established = true

		err = NewOutboxPublisher(service, s.workerID, s.lease, s.poll).Run(ctx)
		if closeResources != nil {
			closeResources()
		}
		if ctx.Err() != nil {
			return nil
		}
		if !errors.Is(err, rabbit.ErrPublisherUnavailable) {
			if err == nil {
				return errors.New("payment outbox publisher exited unexpectedly")
			}
			return err
		}

		log.Printf("service=payment-service component=outbox status=reconnecting retry_in=%s error=%q", s.reconnectDelay, err)
		if !waitForOutboxReconnect(ctx, s.reconnectDelay) {
			return nil
		}
	}
}

func waitForOutboxReconnect(ctx context.Context, delay time.Duration) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func ProcessID(component string) string {
	host, _ := os.Hostname()
	if host == "" {
		host = "unknown"
	}
	return component + "-" + host + "-" + strconv.Itoa(os.Getpid())
}
