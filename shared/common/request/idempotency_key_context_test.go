package request

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestBindIdempotencyKey(t *testing.T) {
	t.Parallel()

	ctx, err := BindIdempotencyKey(context.Background(), "idem-command-123")
	if err != nil {
		t.Fatalf("BindIdempotencyKey() error = %v", err)
	}

	key, ok := IdempotencyKeyFromContext(ctx)
	if !ok || key != "idem-command-123" {
		t.Fatalf("IdempotencyKeyFromContext() = %q, %v", key, ok)
	}

	same, err := BindIdempotencyKey(ctx, "idem-command-123")
	if err != nil {
		t.Fatalf("same-value BindIdempotencyKey() error = %v", err)
	}
	if same != ctx {
		t.Fatal("same-value binding replaced the context")
	}

	conflicted, err := BindIdempotencyKey(ctx, "idem-other-456")
	if !errors.Is(err, ErrIdempotencyKeyContextConflict) {
		t.Fatalf("conflicting BindIdempotencyKey() error = %v", err)
	}
	if conflicted != ctx {
		t.Fatal("conflicting binding replaced the context")
	}
}

func TestBindIdempotencyKeyRejectsInvalidWithoutMutation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	got, err := BindIdempotencyKey(ctx, "unsafe key")
	if !errors.Is(err, ErrInvalidIdempotencyKey) {
		t.Fatalf("BindIdempotencyKey() error = %v", err)
	}
	if got != ctx {
		t.Fatal("invalid binding replaced the context")
	}
}

func TestIdempotencyKeyContextPreservesLifecycle(t *testing.T) {
	t.Parallel()

	deadline := time.Now().Add(time.Minute)
	parent, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	ctx, err := BindIdempotencyKey(parent, "idem-lifecycle-123")
	if err != nil {
		t.Fatalf("BindIdempotencyKey() error = %v", err)
	}

	gotDeadline, ok := ctx.Deadline()
	if !ok || !gotDeadline.Equal(deadline) {
		t.Fatalf("deadline = %v, %v; want %v, true", gotDeadline, ok, deadline)
	}

	cancel()
	if !errors.Is(ctx.Err(), context.Canceled) {
		t.Fatalf("context error = %v, want context.Canceled", ctx.Err())
	}
}

func TestIdempotencyKeyContextIsolationAndNilHandling(t *testing.T) {
	t.Parallel()

	type foreignKey string
	ctx := context.WithValue(context.Background(), foreignKey("idempotency-key"), "foreign-value")
	if key, ok := IdempotencyKeyFromContext(ctx); ok || key != "" {
		t.Fatalf("foreign context value leaked: %q, %v", key, ok)
	}

	if key, ok := IdempotencyKeyFromContext(nil); ok || key != "" {
		t.Fatalf("nil context read = %q, %v", key, ok)
	}

	deferredPanic := false
	func() {
		defer func() {
			deferredPanic = recover() != nil
		}()
		_, _ = BindIdempotencyKey(nil, "idem-command-123")
	}()
	if !deferredPanic {
		t.Fatal("BindIdempotencyKey(nil, ...) did not panic")
	}
}
