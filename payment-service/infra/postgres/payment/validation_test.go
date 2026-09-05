package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "payment/internal/domain/payment"
)

func TestClaimFulfillmentRejectsInvalidRuntimeIdentityBeforeDatabaseUse(t *testing.T) {
	store := &Store{}
	now := time.Date(2026, 8, 19, 10, 0, 0, 0, time.UTC)

	cases := []struct {
		name     string
		workerID string
		now      time.Time
		lease    time.Duration
	}{
		{name: "empty worker", workerID: "", now: now, lease: time.Minute},
		{name: "zero time", workerID: "worker-a", now: time.Time{}, lease: time.Minute},
		{name: "zero lease", workerID: "worker-a", now: now, lease: 0},
		{name: "negative lease", workerID: "worker-a", now: now, lease: -time.Second},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := store.ClaimFulfillment(context.Background(), tc.workerID, tc.now, tc.lease)
			if !errors.Is(err, domain.ErrInvalidCommand) {
				t.Fatalf("ClaimFulfillment error = %v, want ErrInvalidCommand", err)
			}
		})
	}
}

func TestFinalizeAndRedriveRejectInvalidDurableIdentityBeforeDatabaseUse(t *testing.T) {
	store := &Store{}
	now := time.Date(2026, 8, 19, 10, 0, 0, 0, time.UTC)

	if err := store.CompleteFulfillment(context.Background(), 0, "worker-a", 1, now); !errors.Is(err, domain.ErrInvalidCommand) {
		t.Fatalf("CompleteFulfillment invalid id error = %v", err)
	}
	if err := store.CompleteFulfillment(context.Background(), 1, "", 1, now); !errors.Is(err, domain.ErrInvalidCommand) {
		t.Fatalf("CompleteFulfillment invalid worker error = %v", err)
	}
	if err := store.CompleteFulfillment(context.Background(), 1, "worker-a", 0, now); !errors.Is(err, domain.ErrInvalidCommand) {
		t.Fatalf("CompleteFulfillment invalid claim version error = %v", err)
	}

	invalid := []domain.RedriveCommand{
		{},
		{FulfillmentID: 1, ActorID: "admin-1"},
		{FulfillmentID: 1, CommandKey: "redrive-1"},
	}
	for _, command := range invalid {
		_, _, err := store.RedriveFulfillment(context.Background(), command, now)
		if !errors.Is(err, domain.ErrInvalidCommand) {
			t.Fatalf("RedriveFulfillment(%+v) error = %v, want ErrInvalidCommand", command, err)
		}
	}
	_, _, err := store.RedriveFulfillment(context.Background(), domain.RedriveCommand{
		FulfillmentID: 1,
		CommandKey:    "redrive-1",
		ActorID:       "admin-1",
	}, time.Time{})
	if !errors.Is(err, domain.ErrInvalidCommand) {
		t.Fatalf("RedriveFulfillment zero time error = %v, want ErrInvalidCommand", err)
	}
}
