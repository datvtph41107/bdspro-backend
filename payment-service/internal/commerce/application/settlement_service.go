package application

import (
	"context"
	"errors"
	"fmt"

	"payment/internal/commerce/domain"
)

type SettlementService struct {
	orders      OrderStore
	settlements SettlementStore
	now         Clock
}

func NewSettlementService(orders OrderStore, settlements SettlementStore, now Clock) *SettlementService {
	return &SettlementService{orders: orders, settlements: settlements, now: now}
}

func (s *SettlementService) RecordProviderEvent(ctx context.Context, event domain.ProviderEvent) (SettlementResult, error) {
	if s == nil || s.orders == nil || s.settlements == nil || s.now == nil || !event.IsValid() {
		return SettlementResult{}, domain.ErrInvalidCommand
	}
	if err := ctx.Err(); err != nil {
		return SettlementResult{}, err
	}
	switch event.Type {
	case domain.ProviderAuthorized, domain.ProviderCaptured:
		// These provider states are valid protocol evidence but are not enough
		// to create the Subscription effect owned by funds confirmation.
		return SettlementResult{}, nil
	case domain.ProviderFundsConfirmed:
		// Continue into durable Payment settlement.
	default:
		return SettlementResult{}, domain.ErrUnsupportedProviderEvent
	}

	evidenceHash := domain.ProviderEvidenceHash(event)
	existing, found, err := s.settlements.FindSettlementByProviderTransaction(ctx, event.Provider, event.TransactionID)
	if err != nil {
		return SettlementResult{}, err
	}
	if found {
		if existing.EvidenceHash != evidenceHash {
			_ = s.settlements.MarkSettlementRequiresReview(ctx, existing.ID, "provider transaction evidence changed")
			return SettlementResult{}, domain.ErrProviderEvidenceConflict
		}
		return SettlementResult{Settlement: existing, Replay: true}, nil
	}

	order, orderFound, err := s.orders.FindOrderByReference(ctx, event.Reference)
	if err != nil {
		return SettlementResult{}, err
	}

	settlement := domain.Settlement{
		Provider:              event.Provider,
		ProviderTransactionID: event.TransactionID,
		Reference:             event.Reference,
		Amount:                event.Amount,
		OccurredAt:            event.OccurredAt.UTC(),
		EvidenceHash:          evidenceHash,
		CreatedAt:             s.now().UTC(),
	}
	decision := SettlementDecision{Settlement: settlement}

	if !orderFound {
		decision.Settlement.Status = domain.SettlementUnmatched
		decision.Settlement.ReviewReason = "order reference not found"
		return s.settlements.AcceptSettlement(ctx, decision)
	}
	decision.Settlement.OrderID = order.ID

	if order.Terms.Price != event.Amount {
		decision.Settlement.Status = domain.SettlementRequiresReview
		decision.Settlement.ReviewReason = "provider amount or currency differs from frozen order terms"
		return s.settlements.AcceptSettlement(ctx, decision)
	}
	if event.OccurredAt.After(order.ExpiresAt) {
		decision.Settlement.Status = domain.SettlementLate
		decision.Settlement.ReviewReason = "funds occurred after order expiry"
		return s.settlements.AcceptSettlement(ctx, decision)
	}

	occurredAt := event.OccurredAt.UTC()
	decision.Settlement.Status = domain.SettlementFundsConfirmed
	decision.ConfirmOrderID = order.ID
	decision.FundsConfirmedAt = &occurredAt
	decision.CreateFulfillment = true
	result, err := s.settlements.AcceptSettlement(ctx, decision)
	if errors.Is(err, domain.ErrProviderEvidenceConflict) {
		return SettlementResult{}, err
	}
	if err != nil {
		return SettlementResult{}, fmt.Errorf("accept provider funds: %w", err)
	}
	return result, nil
}
