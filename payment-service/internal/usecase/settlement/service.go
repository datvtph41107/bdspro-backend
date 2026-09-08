package settlement

import (
	"context"
	"errors"
	"fmt"
	"time"

	payment "payment/internal/domain/payment"
)

type Clock func() time.Time

type OrderStore interface {
	FindOrderByReference(ctx context.Context, reference string) (payment.Order, bool, error)
}

type Decision struct {
	Settlement        payment.Settlement
	ConfirmOrderID    uint64
	FundsConfirmedAt  *time.Time
	CreateFulfillment bool
	CommandEffect     *payment.CommandEffect
}

type Result struct {
	Settlement  payment.Settlement
	Fulfillment *payment.Fulfillment
	Replay      bool
}

type Store interface {
	FindSettlementByProviderTransaction(ctx context.Context, provider, transactionID string) (payment.Settlement, bool, error)
	AcceptSettlement(ctx context.Context, decision Decision) (Result, error)
	MarkSettlementRequiresReview(ctx context.Context, settlementID uint64, reason string) error
}

type Service struct {
	orders      OrderStore
	settlements Store
	now         Clock
}

func NewService(orders OrderStore, settlements Store, now Clock) *Service {
	return &Service{orders: orders, settlements: settlements, now: now}
}

func (s *Service) RecordProviderEvent(ctx context.Context, event payment.ProviderEvent) (Result, error) {
	return s.recordProviderEvent(ctx, event, nil)
}

// RecordOperatorFundsConfirmation enters manually reconciled bank evidence
// through the same settlement theorem used by provider webhooks. The audit
// effect is committed atomically with settlement/order/fulfillment state.
func (s *Service) RecordOperatorFundsConfirmation(ctx context.Context, event payment.ProviderEvent, effect payment.CommandEffect) (Result, error) {
	if !effect.IsValid() || effect.EffectType != payment.CommandEffectManualFundsConfirmation {
		return Result{}, payment.ErrInvalidCommand
	}
	return s.recordProviderEvent(ctx, event, &effect)
}

func (s *Service) recordProviderEvent(ctx context.Context, event payment.ProviderEvent, effect *payment.CommandEffect) (Result, error) {
	if s == nil || s.orders == nil || s.settlements == nil || s.now == nil || !event.IsValid() {
		return Result{}, payment.ErrInvalidCommand
	}
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	switch event.Type {
	case payment.ProviderAuthorized, payment.ProviderCaptured:
		return Result{}, nil
	case payment.ProviderFundsConfirmed:
	default:
		return Result{}, payment.ErrUnsupportedProviderEvent
	}

	evidenceHash := payment.ProviderEvidenceHash(event)
	existing, found, err := s.settlements.FindSettlementByProviderTransaction(ctx, event.Provider, event.TransactionID)
	if err != nil {
		return Result{}, err
	}
	if found {
		if existing.EvidenceHash != evidenceHash {
			if err := s.settlements.MarkSettlementRequiresReview(ctx, existing.ID, "provider transaction evidence changed"); err != nil {
				return Result{}, fmt.Errorf("mark provider evidence conflict for review: %w", err)
			}
			return Result{}, payment.ErrProviderEvidenceConflict
		}
		return Result{Settlement: existing, Replay: true}, nil
	}

	order, orderFound, err := s.orders.FindOrderByReference(ctx, event.Reference)
	if err != nil {
		return Result{}, err
	}
	value := payment.Settlement{
		Provider:              event.Provider,
		ProviderTransactionID: event.TransactionID,
		Reference:             event.Reference,
		Amount:                event.Amount,
		OccurredAt:            event.OccurredAt.UTC(),
		EvidenceHash:          evidenceHash,
		CreatedAt:             s.now().UTC(),
	}
	decision := Decision{Settlement: value, CommandEffect: effect}

	if !orderFound {
		decision.Settlement.Status = payment.SettlementUnmatched
		decision.Settlement.ReviewReason = "order reference not found"
		return s.settlements.AcceptSettlement(ctx, decision)
	}
	decision.Settlement.OrderID = order.ID

	if order.Terms.Price != event.Amount {
		decision.Settlement.Status = payment.SettlementRequiresReview
		decision.Settlement.ReviewReason = "provider amount or currency differs from frozen order terms"
		return s.settlements.AcceptSettlement(ctx, decision)
	}
	if event.OccurredAt.After(order.ExpiresAt) {
		decision.Settlement.Status = payment.SettlementLate
		decision.Settlement.ReviewReason = "funds occurred after order expiry"
		return s.settlements.AcceptSettlement(ctx, decision)
	}

	occurredAt := event.OccurredAt.UTC()
	decision.Settlement.Status = payment.SettlementFundsConfirmed
	if decision.CommandEffect != nil {
		decision.CommandEffect.Outcome = string(decision.Settlement.Status)
	}
	decision.ConfirmOrderID = order.ID
	decision.FundsConfirmedAt = &occurredAt
	decision.CreateFulfillment = true
	result, err := s.settlements.AcceptSettlement(ctx, decision)
	if errors.Is(err, payment.ErrProviderEvidenceConflict) {
		return Result{}, err
	}
	if err != nil {
		return Result{}, fmt.Errorf("accept provider funds: %w", err)
	}
	return result, nil
}
