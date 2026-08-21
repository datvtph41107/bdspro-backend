package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"payment/internal/commerce/domain"
)

func TestSettlementPreservesMoneyEvidenceAndCreatesOneFulfillment(t *testing.T) {
	now := time.Date(2026, 8, 19, 9, 0, 0, 0, time.UTC)
	store := newMemoryStore()
	orders := NewOrderService(store, staticReference{value: "QHPORDER1"}, func() time.Time { return now })
	order, _, err := orders.CreateOrder(context.Background(), CreateOrderCommand{
		Subject: subject(), Terms: terms(), CommandKey: "checkout-1", ExpiresAt: now.Add(time.Hour),
	})
	if err != nil {
		t.Fatal(err)
	}
	service := NewSettlementService(store, store, func() time.Time { return now.Add(2 * time.Minute) })
	event := domain.ProviderEvent{
		Provider: "sepay", TransactionID: "tx-1", Reference: order.Reference,
		Amount: terms().Price, OccurredAt: now.Add(time.Minute), Type: domain.ProviderFundsConfirmed,
	}
	result, err := service.RecordProviderEvent(context.Background(), event)
	if err != nil || result.Fulfillment == nil || result.Settlement.Status != domain.SettlementFundsConfirmed {
		t.Fatalf("settlement result=%+v err=%v", result, err)
	}
	if len(store.fulfillments) != 1 {
		t.Fatalf("fulfillments=%d want 1", len(store.fulfillments))
	}
	updated, _, _ := store.FindOrderByReference(context.Background(), order.Reference)
	if updated.FundsConfirmedAt == nil || !updated.FundsConfirmedAt.Equal(event.OccurredAt) {
		t.Fatalf("funds confirmed at=%v want %v", updated.FundsConfirmedAt, event.OccurredAt)
	}

	replay, err := service.RecordProviderEvent(context.Background(), event)
	if err != nil || !replay.Replay || len(store.fulfillments) != 1 {
		t.Fatalf("replay=%+v err=%v fulfillments=%d", replay, err, len(store.fulfillments))
	}

	conflict := event
	conflict.Amount.AmountMinor++
	if _, err := service.RecordProviderEvent(context.Background(), conflict); !errors.Is(err, domain.ErrProviderEvidenceConflict) {
		t.Fatalf("conflicting provider evidence err=%v", err)
	}
	settlement, _, _ := store.FindSettlementByProviderTransaction(context.Background(), event.Provider, event.TransactionID)
	if settlement.Status != domain.SettlementRequiresReview {
		t.Fatalf("settlement status=%s want requires_review", settlement.Status)
	}
}

func TestSettlementDoesNotForceAuthorizationIntoFundsConfirmation(t *testing.T) {
	now := time.Date(2026, 8, 19, 9, 0, 0, 0, time.UTC)
	store := newMemoryStore()
	service := NewSettlementService(store, store, func() time.Time { return now })
	result, err := service.RecordProviderEvent(context.Background(), domain.ProviderEvent{
		Provider: "card", TransactionID: "auth-1", Reference: "order-ref",
		Amount: terms().Price, OccurredAt: now, Type: domain.ProviderAuthorized,
	})
	if err != nil || result.Settlement.ID != 0 || len(store.settlements) != 0 || len(store.fulfillments) != 0 {
		t.Fatalf("authorization result=%+v err=%v settlements=%d fulfillments=%d", result, err, len(store.settlements), len(store.fulfillments))
	}
}

func TestSettlementRejectsUnmatchedMismatchAndFundsAfterExpiry(t *testing.T) {
	now := time.Date(2026, 8, 19, 9, 0, 0, 0, time.UTC)
	store := newMemoryStore()
	service := NewSettlementService(store, store, func() time.Time { return now })
	base := domain.ProviderEvent{Provider: "sepay", TransactionID: "unmatched", Reference: "missing", Amount: terms().Price, OccurredAt: now, Type: domain.ProviderFundsConfirmed}
	result, err := service.RecordProviderEvent(context.Background(), base)
	if err != nil || result.Settlement.Status != domain.SettlementUnmatched || result.Fulfillment != nil {
		t.Fatalf("unmatched result=%+v err=%v", result, err)
	}

	orders := NewOrderService(store, staticReference{value: "QHPORDER2"}, func() time.Time { return now })
	order, _, err := orders.CreateOrder(context.Background(), CreateOrderCommand{Subject: subject(), Terms: terms(), CommandKey: "c2", ExpiresAt: now.Add(time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	mismatch := domain.ProviderEvent{Provider: "sepay", TransactionID: "mismatch", Reference: order.Reference, Amount: terms().Price, OccurredAt: now.Add(time.Minute), Type: domain.ProviderFundsConfirmed}
	mismatch.Amount.AmountMinor--
	result, err = service.RecordProviderEvent(context.Background(), mismatch)
	if err != nil || result.Settlement.Status != domain.SettlementRequiresReview || result.Fulfillment != nil {
		t.Fatalf("mismatch result=%+v err=%v", result, err)
	}

	late := domain.ProviderEvent{Provider: "sepay", TransactionID: "late", Reference: order.Reference, Amount: terms().Price, OccurredAt: now.Add(2 * time.Hour), Type: domain.ProviderFundsConfirmed}
	result, err = service.RecordProviderEvent(context.Background(), late)
	if err != nil || result.Settlement.Status != domain.SettlementLate || result.Fulfillment != nil {
		t.Fatalf("late result=%+v err=%v", result, err)
	}
}
