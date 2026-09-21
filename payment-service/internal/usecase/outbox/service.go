package outbox

import (
	"context"
	"errors"
	"hash/fnv"
	"strconv"
	"strings"
	"time"

	payment "payment/internal/domain/payment"
)

var ErrPublisherUnavailable = errors.New("outbox publisher unavailable")

type Clock func() time.Time

type Message struct {
	ID            uint64
	EventID       string
	EventType     string
	SchemaVersion int
	RoutingKey    string
	Payload       []byte
	AttemptCount  int
	ClaimVersion  uint64
	LockedBy      string
}

type Store interface {
	ClaimOutbox(ctx context.Context, workerID string, now time.Time, lease time.Duration) (Message, bool, error)
	MarkOutboxPublished(ctx context.Context, id uint64, workerID string, claimVersion uint64, now time.Time) error
	RetryOutbox(ctx context.Context, id uint64, workerID string, claimVersion uint64, availableAt time.Time, reason string) error
}

type Publisher interface {
	Publish(ctx context.Context, message Message) error
}

// RetryPolicy owns durable message retry timing. Base is the first-attempt
// nominal ceiling; Max caps exponential growth. Deterministic equal jitter
// spreads different event identities without introducing a process-global RNG.
type RetryPolicy struct {
	Base time.Duration
	Max  time.Duration
}

func (p RetryPolicy) Valid() bool {
	return p.Base > 0 && p.Max >= p.Base
}

func (p RetryPolicy) Delay(eventID string, attempt int) time.Duration {
	if !p.Valid() {
		return 0
	}
	nominal, _ := boundedExponential(p.Base, p.Max, attempt)
	return deterministicEqualJitter(nominal, eventID, attempt)
}

type RetryScheduledError struct {
	Cause   error
	Delay   time.Duration
	Attempt int
	EventID string
}

func (e *RetryScheduledError) Error() string {
	if e == nil || e.Cause == nil {
		return "outbox retry scheduled"
	}
	return e.Cause.Error()
}

func (e *RetryScheduledError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

type Service struct {
	store     Store
	publisher Publisher
	now       Clock
	retry     RetryPolicy
}

func NewService(store Store, publisher Publisher, now Clock, retry RetryPolicy) *Service {
	return &Service{store: store, publisher: publisher, now: now, retry: retry}
}

func (s *Service) RunOne(ctx context.Context, workerID string, lease time.Duration) (bool, error) {
	if s == nil || s.store == nil || s.publisher == nil || s.now == nil || !s.retry.Valid() ||
		strings.TrimSpace(workerID) == "" || lease <= 0 {
		return false, payment.ErrInvalidCommand
	}
	now := s.now().UTC()
	message, found, err := s.store.ClaimOutbox(ctx, workerID, now, lease)
	if err != nil || !found {
		return found, err
	}
	if err := s.publisher.Publish(ctx, message); err != nil {
		availableAt := now
		if !errors.Is(err, ErrPublisherUnavailable) {
			delay := s.retry.Delay(message.EventID, message.AttemptCount)
			availableAt = now.Add(delay)
			if retryErr := s.store.RetryOutbox(
				ctx,
				message.ID,
				workerID,
				message.ClaimVersion,
				availableAt,
				err.Error(),
			); retryErr != nil {
				return true, retryErr
			}
			return true, &RetryScheduledError{
				Cause: err, Delay: delay, Attempt: message.AttemptCount, EventID: message.EventID,
			}
		}

		// Transport recovery has a separate owner in OutboxSupervisor. Release
		// this durable claim immediately so reconnect backoff is not stacked
		// with a second message-level delay.
		if retryErr := s.store.RetryOutbox(
			ctx,
			message.ID,
			workerID,
			message.ClaimVersion,
			availableAt,
			err.Error(),
		); retryErr != nil {
			return true, retryErr
		}
		return true, err
	}
	if err := s.store.MarkOutboxPublished(ctx, message.ID, workerID, message.ClaimVersion, now); err != nil {
		return true, err
	}
	return true, nil
}

func boundedExponential(base, max time.Duration, attempt int) (time.Duration, bool) {
	if attempt < 1 {
		attempt = 1
	}
	delay := base
	for i := 1; i < attempt; i++ {
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

func deterministicEqualJitter(nominal time.Duration, identity string, attempt int) time.Duration {
	if nominal <= 1 {
		return nominal
	}
	half := nominal / 2
	span := nominal - half
	if span <= 0 {
		return nominal
	}
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(identity))
	_, _ = hash.Write([]byte(":"))
	_, _ = hash.Write([]byte(strconv.Itoa(attempt)))
	jitter := time.Duration(hash.Sum64() % uint64(span+1))
	return half + jitter
}
