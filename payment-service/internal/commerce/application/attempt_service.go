package application

import (
	"context"
	"strings"

	"payment/internal/commerce/domain"
)

type CreateAttemptCommand struct {
	OrderID    uint64
	Method     domain.PaymentMethod
	CommandKey string
}

type AttemptService struct {
	orders    OrderStore
	attempts  PaymentAttemptStore
	providers []AttemptProvider
	now       Clock
}

func NewAttemptService(orders OrderStore, attempts PaymentAttemptStore, providers []AttemptProvider, now Clock) *AttemptService {
	ordered := make([]AttemptProvider, 0, len(providers))
	for _, provider := range providers {
		if provider == nil || strings.TrimSpace(provider.Code()) == "" {
			continue
		}
		ordered = append(ordered, provider)
	}
	return &AttemptService{orders: orders, attempts: attempts, providers: ordered, now: now}
}

func (s *AttemptService) CreateAttempt(ctx context.Context, command CreateAttemptCommand) (domain.PaymentAttempt, bool, error) {
	if s == nil || s.orders == nil || s.attempts == nil || s.now == nil ||
		command.OrderID == 0 || !command.Method.IsValid() || strings.TrimSpace(command.CommandKey) == "" {
		return domain.PaymentAttempt{}, false, domain.ErrInvalidCommand
	}
	if err := ctx.Err(); err != nil {
		return domain.PaymentAttempt{}, false, err
	}
	// The caller chooses a method, never a provider. Provider routing is Payment
	// policy/capability ownership; exposing it to RN/User would make external
	// callers author infrastructure routing and would make fallback unsafe.
	provider := s.selectProvider(command.Method)
	if provider == nil {
		return domain.PaymentAttempt{}, false, domain.ErrUnsupportedPaymentMethod
	}

	// CommandKey is checked before contacting the provider. A transport retry
	// must replay an existing Attempt instead of creating a second provider-side
	// payment session. The same identity is also passed to the provider so a DB
	// failure after provider creation can be retried idempotently there.
	existing, found, err := s.attempts.FindAttemptByCommand(ctx, command.OrderID, command.CommandKey)
	if err != nil {
		return domain.PaymentAttempt{}, false, err
	}
	fingerprint := domain.PaymentAttemptFingerprint(command.OrderID, command.Method, provider.Code())
	if found {
		if existing.CommandFingerprint != fingerprint {
			return domain.PaymentAttempt{}, false, domain.ErrAttemptCommandConflict
		}
		return existing, false, nil
	}

	order, found, err := s.orders.FindOrderByID(ctx, command.OrderID)
	if err != nil {
		return domain.PaymentAttempt{}, false, err
	}
	if !found {
		return domain.PaymentAttempt{}, false, domain.ErrOrderNotFound
	}
	now := s.now().UTC()
	if order.Status != domain.OrderPendingFunds || !order.ExpiresAt.After(now) {
		return domain.PaymentAttempt{}, false, domain.ErrAttemptNotAllowed
	}

	providerAttempt, err := provider.CreateAttempt(ctx, order, command.Method, command.CommandKey)
	if err != nil {
		return domain.PaymentAttempt{}, false, err
	}
	if !providerAttempt.IsValid() {
		return domain.PaymentAttempt{}, false, domain.ErrProviderUnavailable
	}
	attempt := domain.PaymentAttempt{
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
	return s.attempts.CreateAttempt(ctx, attempt)
}

func (s *AttemptService) selectProvider(method domain.PaymentMethod) AttemptProvider {
	for _, provider := range s.providers {
		if provider.Supports(method) {
			return provider
		}
	}
	return nil
}
