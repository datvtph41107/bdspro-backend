package application

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"payment/internal/commerce/domain"
)

func TestFulfillmentResponseLossRetriesSameEffectWithoutDoubleBusinessApply(t *testing.T) {
	now := time.Date(2026, 8, 19, 9, 0, 0, 0, time.UTC)
	store := newMemoryStore()
	orders := NewOrderService(store, staticReference{value: "QHPORDER1"}, func() time.Time { return now })
	order, _, err := orders.CreateOrder(context.Background(), CreateOrderCommand{Subject: subject(), Terms: terms(), CommandKey: "checkout-1", ExpiresAt: now.Add(time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	settlements := NewSettlementService(store, store, func() time.Time { return now.Add(time.Minute) })
	_, err = settlements.RecordProviderEvent(context.Background(), domain.ProviderEvent{Provider: "sepay", TransactionID: "tx-1", Reference: order.Reference, Amount: terms().Price, OccurredAt: now.Add(30 * time.Second), Type: domain.ProviderFundsConfirmed})
	if err != nil {
		t.Fatal(err)
	}

	consumer := newIdempotentSubscription(true)
	service := NewFulfillmentService(store, consumer, func() time.Time { return now.Add(2 * time.Minute) }, time.Second)
	claim1, ok, err := store.ClaimFulfillment(context.Background(), "worker-a", now.Add(2*time.Minute), time.Minute)
	if err != nil || !ok {
		t.Fatalf("claim1 ok=%v err=%v", ok, err)
	}
	if err := service.ProcessClaim(context.Background(), claim1, "worker-a"); err != nil {
		t.Fatalf("process response-loss claim: %v", err)
	}
	if consumer.businessApplyCount != 1 {
		t.Fatalf("business apply count after lost response=%d want 1", consumer.businessApplyCount)
	}

	claim2, ok, err := store.ClaimFulfillment(context.Background(), "worker-b", now.Add(3*time.Minute), time.Minute)
	if err != nil || !ok {
		t.Fatalf("claim2 ok=%v err=%v", ok, err)
	}
	if err := service.ProcessClaim(context.Background(), claim2, "worker-b"); err != nil {
		t.Fatalf("process replay claim: %v", err)
	}
	if consumer.businessApplyCount != 1 {
		t.Fatalf("business apply count after replay=%d want 1", consumer.businessApplyCount)
	}
	fulfillment := store.fulfillments[claim2.Fulfillment.ID]
	if fulfillment.Status != domain.FulfillmentCompleted {
		t.Fatalf("fulfillment status=%s want completed", fulfillment.Status)
	}
	if effect := consumer.effects[fmt.Sprintf("payment.order.%d", order.ID)]; !effect.OccurredAt.Equal(now.Add(30 * time.Second)) {
		t.Fatalf("effect occurred_at=%v", effect.OccurredAt)
	}
}

func TestFulfillmentFencingRejectsStaleClaim(t *testing.T) {
	now := time.Date(2026, 8, 19, 9, 0, 0, 0, time.UTC)
	store := newMemoryStore()
	order := domain.Order{ID: 1, Subject: subject(), Terms: terms(), CommandKey: "c", CommandFingerprint: "fp", Reference: "ref", Status: domain.OrderFundsConfirmed, CreatedAt: now}
	confirmed := now
	order.FundsConfirmedAt = &confirmed
	store.orders[1] = order
	store.fulfillments[1] = domain.Fulfillment{ID: 1, OrderID: 1, Status: domain.FulfillmentPending, AvailableAt: now, CreatedAt: now, UpdatedAt: now}
	store.fulfillmentByOrder[1] = 1

	claimA, ok, err := store.ClaimFulfillment(context.Background(), "worker-a", now, time.Minute)
	if err != nil || !ok || claimA.Fulfillment.ClaimVersion != 1 {
		t.Fatalf("claim A=%+v ok=%v err=%v", claimA, ok, err)
	}
	claimB, ok, err := store.ClaimFulfillment(context.Background(), "worker-b", now.Add(2*time.Minute), time.Minute)
	if err != nil || !ok || claimB.Fulfillment.ClaimVersion != 2 {
		t.Fatalf("claim B=%+v ok=%v err=%v", claimB, ok, err)
	}
	if err := store.CompleteFulfillment(context.Background(), 1, "worker-a", 1, now.Add(2*time.Minute)); !errors.Is(err, domain.ErrClaimLost) {
		t.Fatalf("stale finalize err=%v", err)
	}
	if err := store.CompleteFulfillment(context.Background(), 1, "worker-b", 2, now.Add(2*time.Minute)); err != nil {
		t.Fatalf("current finalize err=%v", err)
	}
}
