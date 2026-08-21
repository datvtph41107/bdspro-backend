package application

import (
	"context"
	"testing"
	"time"

	"payment/internal/commerce/domain"
)

func TestRedriveConsumesIdentityForChangeAndNoop(t *testing.T) {
	now := time.Date(2026, 8, 19, 9, 0, 0, 0, time.UTC)
	store := newMemoryStore()
	store.fulfillments[1] = domain.Fulfillment{ID: 1, OrderID: 1, Status: domain.FulfillmentRequiresReview, CreatedAt: now, UpdatedAt: now}
	store.fulfillments[2] = domain.Fulfillment{ID: 2, OrderID: 2, Status: domain.FulfillmentCompleted, CreatedAt: now, UpdatedAt: now}
	service := NewRecoveryService(store, func() time.Time { return now })

	changed, replay, err := service.RedriveFulfillment(context.Background(), domain.RedriveCommand{FulfillmentID: 1, CommandKey: "abc", ActorID: "admin-1"})
	if err != nil || !changed || replay || store.fulfillments[1].Status != domain.FulfillmentPending {
		t.Fatalf("redrive change changed=%v replay=%v err=%v status=%s", changed, replay, err, store.fulfillments[1].Status)
	}
	changed, replay, err = service.RedriveFulfillment(context.Background(), domain.RedriveCommand{FulfillmentID: 1, CommandKey: "abc", ActorID: "admin-1"})
	if err != nil || changed || !replay {
		t.Fatalf("redrive replay changed=%v replay=%v err=%v", changed, replay, err)
	}

	changed, replay, err = service.RedriveFulfillment(context.Background(), domain.RedriveCommand{FulfillmentID: 2, CommandKey: "noop", ActorID: "admin-1"})
	if err != nil || changed || replay {
		t.Fatalf("noop changed=%v replay=%v err=%v", changed, replay, err)
	}
	// Simulate a later state change. The old no-op command identity must remain consumed.
	f := store.fulfillments[2]
	f.Status = domain.FulfillmentRequiresReview
	store.fulfillments[2] = f
	changed, replay, err = service.RedriveFulfillment(context.Background(), domain.RedriveCommand{FulfillmentID: 2, CommandKey: "noop", ActorID: "admin-1"})
	if err != nil || changed || !replay || store.fulfillments[2].Status != domain.FulfillmentRequiresReview {
		t.Fatalf("noop replay changed=%v replay=%v err=%v status=%s", changed, replay, err, store.fulfillments[2].Status)
	}
}
