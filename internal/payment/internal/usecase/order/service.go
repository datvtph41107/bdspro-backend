package order

import (
	"context"
	"errors"
	"strings"
	"time"

	payment "payment/internal/domain/payment"
)

type Clock func() time.Time

type ReferenceGenerator interface {
	NewOrderReference() (string, error)
}

type Store interface {
	FindOrderByCommand(ctx context.Context, subject payment.Subject, productCode, commandKey string) (payment.Order, bool, error)
	CreateOrder(ctx context.Context, order payment.Order) (payment.Order, bool, error)
}

type CreateCommand struct {
	Subject    payment.Subject
	Terms      payment.CommercialTerms
	CommandKey string
}

type Service struct {
	store     Store
	reference ReferenceGenerator
	now       Clock
	ttl       time.Duration
}

func NewService(store Store, reference ReferenceGenerator, now Clock, ttl time.Duration) *Service {
	return &Service{store: store, reference: reference, now: now, ttl: ttl}
}

func (s *Service) Create(ctx context.Context, command CreateCommand) (payment.Order, bool, error) {
	if s == nil || s.store == nil || s.reference == nil || s.now == nil || s.ttl <= 0 ||
		!command.Subject.IsValid() || strings.TrimSpace(command.CommandKey) == "" {
		return payment.Order{}, false, payment.ErrInvalidCommand
	}
	if !command.Terms.IsValid() {
		return payment.Order{}, false, payment.ErrInvalidCommercialTerms
	}
	if err := ctx.Err(); err != nil {
		return payment.Order{}, false, err
	}

	fingerprint := payment.OrderFingerprint(command.Subject, command.Terms)
	existing, found, err := s.store.FindOrderByCommand(ctx, command.Subject, command.Terms.ProductCode, command.CommandKey)
	if err != nil {
		return payment.Order{}, false, err
	}
	if found {
		if existing.CommandFingerprint != fingerprint {
			return payment.Order{}, false, payment.ErrOrderCommandConflict
		}
		return existing, false, nil
	}

	reference, err := s.reference.NewOrderReference()
	if err != nil {
		return payment.Order{}, false, err
	}
	now := s.now().UTC()
	expiresAt := now.Add(s.ttl)

	value := payment.Order{
		Subject:            command.Subject,
		Terms:              command.Terms,
		CommandKey:         strings.TrimSpace(command.CommandKey),
		CommandFingerprint: fingerprint,
		Reference:          reference,
		Status:             payment.OrderPendingFunds,
		ExpiresAt:          expiresAt,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	saved, created, err := s.store.CreateOrder(ctx, value)
	if errors.Is(err, payment.ErrOrderCommandConflict) {
		return payment.Order{}, false, err
	}
	if err != nil {
		return payment.Order{}, false, err
	}
	return saved, created, nil
}
