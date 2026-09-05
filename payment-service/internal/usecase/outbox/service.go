package outbox

import (
	"context"
	"strings"
	"time"

	payment "payment/internal/domain/payment"
)

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

type Service struct {
	store      Store
	publisher  Publisher
	now        Clock
	retryDelay time.Duration
}

func NewService(store Store, publisher Publisher, now Clock, retryDelay time.Duration) *Service {
	return &Service{store: store, publisher: publisher, now: now, retryDelay: retryDelay}
}

func (s *Service) RunOne(ctx context.Context, workerID string, lease time.Duration) (bool, error) {
	if s == nil || s.store == nil || s.publisher == nil || s.now == nil || strings.TrimSpace(workerID) == "" || lease <= 0 {
		return false, payment.ErrInvalidCommand
	}
	now := s.now().UTC()
	message, found, err := s.store.ClaimOutbox(ctx, workerID, now, lease)
	if err != nil || !found {
		return found, err
	}
	if err := s.publisher.Publish(ctx, message); err != nil {
		if retryErr := s.store.RetryOutbox(ctx, message.ID, workerID, message.ClaimVersion, now.Add(s.retryDelay), err.Error()); retryErr != nil {
			return true, retryErr
		}
		return true, err
	}
	if err := s.store.MarkOutboxPublished(ctx, message.ID, workerID, message.ClaimVersion, now); err != nil {
		return true, err
	}
	return true, nil
}
