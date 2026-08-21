package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"payment/internal/commerce/domain"
)

type CreateOrderCommand struct {
	Subject    domain.Subject
	Terms      domain.CommercialTerms
	CommandKey string
	ExpiresAt  time.Time
}

type OrderService struct {
	store     OrderStore
	reference ReferenceGenerator
	now       Clock
}

func NewOrderService(store OrderStore, reference ReferenceGenerator, now Clock) *OrderService {
	return &OrderService{store: store, reference: reference, now: now}
}

func (s *OrderService) CreateOrder(ctx context.Context, command CreateOrderCommand) (domain.Order, bool, error) {
	if s == nil || s.store == nil || s.reference == nil || s.now == nil ||
		!command.Subject.IsValid() || strings.TrimSpace(command.CommandKey) == "" {
		return domain.Order{}, false, domain.ErrInvalidCommand
	}
	if !command.Terms.IsValid() {
		return domain.Order{}, false, domain.ErrInvalidCommercialTerms
	}
	if err := ctx.Err(); err != nil {
		return domain.Order{}, false, err
	}

	fingerprint := domain.OrderFingerprint(command.Subject, command.Terms)
	existing, found, err := s.store.FindOrderByCommand(ctx, command.Subject, command.Terms.ProductCode, command.CommandKey)
	if err != nil {
		return domain.Order{}, false, err
	}
	if found {
		if existing.CommandFingerprint != fingerprint {
			return domain.Order{}, false, domain.ErrOrderCommandConflict
		}
		return existing, false, nil
	}

	reference, err := s.reference.NewOrderReference()
	if err != nil {
		return domain.Order{}, false, err
	}
	now := s.now().UTC()
	if command.ExpiresAt.IsZero() || !command.ExpiresAt.After(now) {
		return domain.Order{}, false, domain.ErrInvalidCommand
	}

	order := domain.Order{
		Subject:            command.Subject,
		Terms:              command.Terms,
		CommandKey:         strings.TrimSpace(command.CommandKey),
		CommandFingerprint: fingerprint,
		Reference:          reference,
		Status:             domain.OrderPendingFunds,
		ExpiresAt:          command.ExpiresAt.UTC(),
		CreatedAt:          now,
	}
	saved, created, err := s.store.CreateOrder(ctx, order)
	if errors.Is(err, domain.ErrOrderCommandConflict) {
		return domain.Order{}, false, err
	}
	if err != nil {
		return domain.Order{}, false, err
	}
	return saved, created, nil
}
