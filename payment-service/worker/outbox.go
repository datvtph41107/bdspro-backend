package worker

import (
	"common/logging"
	"context"
	"errors"
	"hash/fnv"
	"log/slog"
	"os"
	"strconv"
	"time"

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
			if errors.Is(err, outbox.ErrPublisherUnavailable) {
				return err
			}

			attrs := []any{slog.Any("error", err)}
			var scheduled *outbox.RetryScheduledError
			if errors.As(err, &scheduled) {
				attrs = append(attrs,
					slog.Duration("retry_in", scheduled.Delay),
					slog.Int("attempt_count", scheduled.Attempt),
					slog.String("event_id", scheduled.EventID),
				)
			}
			logging.WithComponent(ctx, "outbox").Warn("outbox publish retry", attrs...)
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

// ReconnectPolicy owns only broker transport recovery timing. It deliberately
// does not own durable event retry timing.
type ReconnectPolicy struct {
	Base time.Duration
	Max  time.Duration
}

func (p ReconnectPolicy) Valid() bool {
	return p.Base > 0 && p.Max >= p.Base
}

func (p ReconnectPolicy) Delay(workerID string, failureCount int) (time.Duration, bool) {
	if !p.Valid() {
		return 0, false
	}
	nominal, capped := boundedReconnectDelay(p.Base, p.Max, failureCount)
	return reconnectJitter(nominal, workerID, failureCount), capped
}

// OutboxSupervisor preserves the Payment process across broker unavailability,
// including broker absence at process startup. Durable claim/retry semantics
// remain in the outbox service; this actor only reconstructs transport.
type OutboxSupervisor struct {
	factory   OutboxServiceFactory
	workerID  string
	lease     time.Duration
	poll      time.Duration
	reconnect ReconnectPolicy
}

func NewOutboxSupervisor(
	factory OutboxServiceFactory,
	workerID string,
	lease, poll time.Duration,
	reconnect ReconnectPolicy,
) *OutboxSupervisor {
	return &OutboxSupervisor{
		factory: factory, workerID: workerID, lease: lease, poll: poll, reconnect: reconnect,
	}
}

func (s *OutboxSupervisor) Run(ctx context.Context) error {
	if s == nil || s.factory == nil || s.workerID == "" || s.lease <= 0 || s.poll <= 0 || !s.reconnect.Valid() {
		return errors.New("payment outbox supervisor is not configured")
	}

	logger := logging.WithComponent(ctx, "outbox")
	failureCount := 0
	var degradedSince time.Time
	escalated := false

	for {
		if ctx.Err() != nil {
			return nil
		}

		service, closeResources, err := s.factory()
		if err != nil {
			failureCount++
			if degradedSince.IsZero() {
				degradedSince = time.Now().UTC()
			}
			delay, capped := s.reconnect.Delay(s.workerID, failureCount)
			logOutboxReconnectFailure(logger, err, delay, failureCount, time.Since(degradedSince), capped, &escalated)
			if !waitForOutboxReconnect(ctx, delay) {
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

		if failureCount > 0 {
			logger.Info(
				"outbox broker recovered",
				slog.Int("failure_count", failureCount),
				slog.Duration("degraded_for", time.Since(degradedSince)),
			)
			failureCount = 0
			degradedSince = time.Time{}
			escalated = false
		}

		err = NewOutboxPublisher(service, s.workerID, s.lease, s.poll).Run(ctx)
		if closeResources != nil {
			closeResources()
		}
		if ctx.Err() != nil {
			return nil
		}
		if !errors.Is(err, outbox.ErrPublisherUnavailable) {
			if err == nil {
				return errors.New("payment outbox publisher exited unexpectedly")
			}
			return err
		}

		failureCount++
		if degradedSince.IsZero() {
			degradedSince = time.Now().UTC()
		}
		delay, capped := s.reconnect.Delay(s.workerID, failureCount)
		logOutboxReconnectFailure(logger, err, delay, failureCount, time.Since(degradedSince), capped, &escalated)
		if !waitForOutboxReconnect(ctx, delay) {
			return nil
		}
	}
}

func logOutboxReconnectFailure(
	logger *slog.Logger,
	err error,
	delay time.Duration,
	failureCount int,
	degradedFor time.Duration,
	capped bool,
	escalated *bool,
) {
	attrs := []any{
		slog.Duration("retry_in", delay),
		slog.Int("failure_count", failureCount),
		slog.Duration("degraded_for", degradedFor),
		slog.Any("error", err),
	}
	if capped && escalated != nil && !*escalated {
		logger.Error("outbox broker degraded", attrs...)
		*escalated = true
		return
	}
	logger.Warn("outbox broker reconnecting", attrs...)
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

func boundedReconnectDelay(base, max time.Duration, failureCount int) (time.Duration, bool) {
	if failureCount < 1 {
		failureCount = 1
	}
	delay := base
	for i := 1; i < failureCount; i++ {
		if delay >= max || delay > max/2 {
			return max, true
		}
		delay *= 2
	}
	if delay >= max {
		return max, true
	}
	return delay, false
}

func reconnectJitter(nominal time.Duration, workerID string, failureCount int) time.Duration {
	if nominal <= 1 {
		return nominal
	}
	half := nominal / 2
	span := nominal - half
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(workerID))
	_, _ = hash.Write([]byte(":"))
	_, _ = hash.Write([]byte(strconv.Itoa(failureCount)))
	return half + time.Duration(hash.Sum64()%uint64(span+1))
}

func ProcessID(component string) string {
	host, _ := os.Hostname()
	if host == "" {
		host = "unknown"
	}
	return component + "-" + host + "-" + strconv.Itoa(os.Getpid())
}
