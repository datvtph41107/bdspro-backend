package attempt

import (
	"context"
	"strings"
	"time"

	payment "payment/internal/domain/payment"
)

type Clock func() time.Time

type OrderStore interface {
	FindOrderByID(ctx context.Context, orderID uint64) (payment.Order, bool, error)
}

type Store interface {
	FindAttemptByCommand(ctx context.Context, orderID uint64, commandKey string) (payment.PaymentAttempt, bool, error)
	CreateAttempt(ctx context.Context, attempt payment.PaymentAttempt) (payment.PaymentAttempt, bool, error)
	UpdateAttemptFromProvider(ctx context.Context, attemptID uint64, provider payment.ProviderAttempt, now time.Time) (payment.PaymentAttempt, error)
}

type Provider interface {
	Code() string
	Supports(method payment.PaymentMethod) bool
	CreateAttempt(ctx context.Context, order payment.Order, method payment.PaymentMethod, commandKey string) (payment.ProviderAttempt, error)
}

type CreateCommand struct {
	Subject    payment.Subject
	OrderID    uint64
	Method     payment.PaymentMethod
	CommandKey string
}

type Service struct {
	orders    OrderStore
	attempts  Store
	providers []Provider
	now       Clock
}

func NewService(orders OrderStore, attempts Store, providers []Provider, now Clock) *Service {
	ordered := make([]Provider, 0, len(providers))
	for _, provider := range providers {
		if provider == nil || strings.TrimSpace(provider.Code()) == "" {
			continue
		}
		ordered = append(ordered, provider)
	}
	return &Service{orders: orders, attempts: attempts, providers: ordered, now: now}
}

func (s *Service) Create(ctx context.Context, command CreateCommand) (payment.PaymentAttempt, bool, error) {
	if s == nil || s.orders == nil || s.attempts == nil || s.now == nil ||
		!command.Subject.IsValid() || command.OrderID == 0 || !command.Method.IsValid() || strings.TrimSpace(command.CommandKey) == "" {
		return payment.PaymentAttempt{}, false, payment.ErrInvalidCommand
	}
	if err := ctx.Err(); err != nil {
		return payment.PaymentAttempt{}, false, err
	}
	provider := s.selectProvider(command.Method)
	if provider == nil {
		return payment.PaymentAttempt{}, false, payment.ErrUnsupportedPaymentMethod
	}

	existing, found, err := s.attempts.FindAttemptByCommand(ctx, command.OrderID, command.CommandKey)
	if err != nil {
		return payment.PaymentAttempt{}, false, err
	}
	fingerprint := payment.PaymentAttemptFingerprint(command.OrderID, command.Method, provider.Code())
	if found {
		if existing.CommandFingerprint != fingerprint {
			return payment.PaymentAttempt{}, false, payment.ErrAttemptCommandConflict
		}
		return existing, false, nil
	}

	order, found, err := s.orders.FindOrderByID(ctx, command.OrderID)
	if err != nil {
		return payment.PaymentAttempt{}, false, err
	}
	if !found {
		return payment.PaymentAttempt{}, false, payment.ErrOrderNotFound
	}
	if order.Subject != command.Subject {
		return payment.PaymentAttempt{}, false, payment.ErrAttemptNotAllowed
	}
	now := s.now().UTC()
	if order.Status != payment.OrderPendingFunds || !order.ExpiresAt.After(now) {
		return payment.PaymentAttempt{}, false, payment.ErrAttemptNotAllowed
	}

	providerAttempt, err := provider.CreateAttempt(ctx, order, command.Method, command.CommandKey)
	if err != nil {
		return payment.PaymentAttempt{}, false, err
	}
	if !providerAttempt.IsValid() {
		return payment.PaymentAttempt{}, false, payment.ErrProviderUnavailable
	}
	value := payment.PaymentAttempt{
		OrderID:            order.ID,
		CommandKey:         strings.TrimSpace(command.CommandKey),
		CommandFingerprint: fingerprint,
		Method:             command.Method,
		Provider:           provider.Code(),
		ProviderReference:  providerAttempt.ProviderReference,
		Status:             providerAttempt.Status,
		NextAction:         providerAttempt.NextAction,
		ExpiresAt:          providerAttempt.ExpiresAt,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	return s.attempts.CreateAttempt(ctx, value)
}

func (s *Service) selectProvider(method payment.PaymentMethod) Provider {
	for _, provider := range s.providers {
		if provider.Supports(method) {
			return provider
		}
	}
	return nil
}
