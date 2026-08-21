package application

import (
	"context"
	"strings"
	"time"

	"payment/internal/commerce/domain"
)

type OutboxPublisherService struct {
	store      OutboxStore
	publisher  EventPublisher
	now        Clock
	retryDelay time.Duration
}

func NewOutboxPublisherService(store OutboxStore, publisher EventPublisher, now Clock, retryDelay time.Duration) *OutboxPublisherService {
	return &OutboxPublisherService{store: store, publisher: publisher, now: now, retryDelay: retryDelay}
}

func (s *OutboxPublisherService) RunOne(ctx context.Context, workerID string, lease time.Duration) (bool, error) {
	if s == nil || s.store == nil || s.publisher == nil || s.now == nil || strings.TrimSpace(workerID) == "" || lease <= 0 {
		return false, domain.ErrInvalidCommand
	}
	now := s.now().UTC()
	// Claiming is a durable responsibility transfer, not an in-memory queue pop.
	// Store fencing makes a worker that wakes after lease expiry unable to mark a
	// message published/retry after another process has reclaimed it.
	message, found, err := s.store.ClaimOutbox(ctx, workerID, now, lease)
	if err != nil || !found {
		return found, err
	}
	// Publish failure includes UNKNOWN confirmation outcomes. We deliberately do
	// not infer "not delivered" from a timeout; responsibility stays durable and
	// a later retry may duplicate transport delivery. Consumer Inbox handles it.
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
