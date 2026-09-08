package request

import (
	"context"
	"testing"
	"time"
)

func TestRequestIDContextRoundTrip(
	t *testing.T,
) {
	t.Parallel()

	ctx := WithRequestID(
		context.Background(),
		"request-context-123",
	)

	actual, ok := RequestIDFromContext(ctx)

	if !ok {
		t.Fatal("request ID missing from context")
	}

	if actual != "request-context-123" {
		t.Fatalf(
			"request ID = %q",
			actual,
		)
	}
}

func TestWithRequestIDPreservesCancellation(
	t *testing.T,
) {
	t.Parallel()

	parent, cancel := context.WithCancel(
		context.Background(),
	)

	ctx := WithRequestID(
		parent,
		"request-context-123",
	)

	cancel()

	select {
	case <-ctx.Done():
		// Pass
	case <-time.After(time.Second):
		t.Fatal(
			"request ID context lost parent cancellation",
		)
	}
}

func TestWithRequestIDPreservesDeadline(
	t *testing.T,
) {
	t.Parallel()

	parentDeadline := time.Now().Add(time.Minute)

	parent, cancel := context.WithDeadline(
		context.Background(),
		parentDeadline,
	)

	defer cancel()

	ctx := WithRequestID(
		parent,
		"request-context-123",
	)

	actualDeadline, ok := ctx.Deadline()

	if !ok {
		t.Fatal("derived context lost deadline")
	}

	if !actualDeadline.Equal(parentDeadline) {
		t.Fatalf(
			"deadline = %v, want %v",
			actualDeadline,
			parentDeadline,
		)
	}
}

func TestWithRequestIDIgnoresInvalidValue(
	t *testing.T,
) {
	t.Parallel()

	ctx := WithRequestID(
		context.Background(),
		"unsafe rquest id",
	)

	if actual, ok := RequestIDFromContext(ctx); ok {
		t.Fatalf(
			"invalid request ID reached context: %q",
			actual,
		)
	}
}
