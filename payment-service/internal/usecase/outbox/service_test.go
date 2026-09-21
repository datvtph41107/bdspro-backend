package outbox

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

type testStore struct {
	message       Message
	found         bool
	claimErr      error
	retryAt       time.Time
	retryReason   string
	retryCalls    int
	publishedAt   time.Time
	publishedCall int
}

func (s *testStore) ClaimOutbox(context.Context, string, time.Time, time.Duration) (Message, bool, error) {
	return s.message, s.found, s.claimErr
}

func (s *testStore) MarkOutboxPublished(_ context.Context, _ uint64, _ string, _ uint64, now time.Time) error {
	s.publishedAt = now
	s.publishedCall++
	return nil
}

func (s *testStore) RetryOutbox(
	_ context.Context,
	_ uint64,
	_ string,
	_ uint64,
	availableAt time.Time,
	reason string,
) error {
	s.retryAt = availableAt
	s.retryReason = reason
	s.retryCalls++
	return nil
}

type testPublisher struct{ err error }

func (p testPublisher) Publish(context.Context, Message) error { return p.err }

func TestRetryPolicyIsDeterministicBoundedAndAttemptAware(t *testing.T) {
	policy := RetryPolicy{Base: 8 * time.Second, Max: 30 * time.Second}

	first := policy.Delay("evt-1", 1)
	if first < 4*time.Second || first > 8*time.Second {
		t.Fatalf("first delay = %s, want within [4s,8s]", first)
	}
	if got := policy.Delay("evt-1", 1); got != first {
		t.Fatalf("deterministic delay = %s, want %s", got, first)
	}

	third := policy.Delay("evt-1", 3)
	if third < 15*time.Second || third > 30*time.Second {
		t.Fatalf("third delay = %s, want capped equal-jitter within [15s,30s]", third)
	}
	if third <= first {
		t.Fatalf("third delay = %s, want greater than first delay %s", third, first)
	}
}

func TestRunOneSchedulesAttemptAwareRetryForNonTransportFailure(t *testing.T) {
	now := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)
	policy := RetryPolicy{Base: 5 * time.Second, Max: time.Minute}
	store := &testStore{
		found: true,
		message: Message{
			ID: 7, EventID: "payment-completed:42", EventType: "payment.completed.v1",
			RoutingKey: "payment.completed.v1", Payload: []byte("{}"),
			AttemptCount: 3, ClaimVersion: 9,
		},
	}
	publishErr := errors.New("negative publisher confirm")
	service := NewService(store, testPublisher{err: publishErr}, func() time.Time { return now }, policy)

	processed, err := service.RunOne(context.Background(), "worker-a", 30*time.Second)
	if !processed {
		t.Fatal("RunOne processed = false, want true")
	}
	var scheduled *RetryScheduledError
	if !errors.As(err, &scheduled) {
		t.Fatalf("RunOne error = %T %v, want RetryScheduledError", err, err)
	}
	wantDelay := policy.Delay(store.message.EventID, store.message.AttemptCount)
	if scheduled.Delay != wantDelay || scheduled.Attempt != 3 || scheduled.EventID != store.message.EventID {
		t.Fatalf("scheduled = %#v, want delay=%s attempt=3 event=%s", scheduled, wantDelay, store.message.EventID)
	}
	if !errors.Is(err, publishErr) {
		t.Fatalf("RunOne error = %v, want wrapped publish error", err)
	}
	if got, want := store.retryAt, now.Add(wantDelay); !got.Equal(want) {
		t.Fatalf("retryAt = %s, want %s", got, want)
	}
	if store.retryCalls != 1 {
		t.Fatalf("retry calls = %d, want 1", store.retryCalls)
	}
}

func TestRunOneReleasesClaimImmediatelyForPublisherUnavailable(t *testing.T) {
	now := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)
	store := &testStore{
		found: true,
		message: Message{
			ID: 8, EventID: "payment-completed:43", EventType: "payment.completed.v1",
			RoutingKey: "payment.completed.v1", Payload: []byte("{}"),
			AttemptCount: 4, ClaimVersion: 10,
		},
	}
	unavailable := fmt.Errorf("%w: connection closed", ErrPublisherUnavailable)
	service := NewService(
		store,
		testPublisher{err: unavailable},
		func() time.Time { return now },
		RetryPolicy{Base: 5 * time.Second, Max: time.Minute},
	)

	processed, err := service.RunOne(context.Background(), "worker-a", 30*time.Second)
	if !processed {
		t.Fatal("RunOne processed = false, want true")
	}
	if !errors.Is(err, ErrPublisherUnavailable) {
		t.Fatalf("RunOne error = %v, want ErrPublisherUnavailable", err)
	}
	if !store.retryAt.Equal(now) {
		t.Fatalf("transport retryAt = %s, want immediate release at %s", store.retryAt, now)
	}
}

func TestRunOneMarksPublishedAfterConfirmedPublish(t *testing.T) {
	now := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)
	store := &testStore{
		found: true,
		message: Message{
			ID: 9, EventID: "payment-completed:44", EventType: "payment.completed.v1",
			RoutingKey: "payment.completed.v1", Payload: []byte("{}"),
			AttemptCount: 1, ClaimVersion: 11,
		},
	}
	service := NewService(
		store,
		testPublisher{},
		func() time.Time { return now },
		RetryPolicy{Base: 5 * time.Second, Max: time.Minute},
	)

	processed, err := service.RunOne(context.Background(), "worker-a", 30*time.Second)
	if err != nil || !processed {
		t.Fatalf("RunOne = processed=%v err=%v, want true,nil", processed, err)
	}
	if store.publishedCall != 1 || !store.publishedAt.Equal(now) {
		t.Fatalf("published calls=%d at=%s, want 1 at %s", store.publishedCall, store.publishedAt, now)
	}
}
