package request

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestOperationIDContextRoundTrip(
	t *testing.T,
) {
	t.Parallel()

	ctx, err :=
		BindOperationID(
			context.Background(),
			"op-context-123",
		)

	if err != nil {
		t.Fatalf(
			"BindOperationID() error = %v",
			err,
		)
	}

	got, ok :=
		OperationIDFromContext(
			ctx,
		)

	if !ok {
		t.Fatal(
			"operation ID missing from context",
		)
	}

	if got !=
		"op-context-123" {

		t.Fatalf(
			"operation ID = %q, want %q",
			got,
			"op-context-123",
		)
	}
}

func TestBindOperationIDDoesNotMutateParent(
	t *testing.T,
) {
	t.Parallel()

	parent :=
		context.Background()

	child, err :=
		BindOperationID(
			parent,
			"op-child-123",
		)

	if err != nil {
		t.Fatalf(
			"BindOperationID() error = %v",
			err,
		)
	}

	if got, ok :=
		OperationIDFromContext(
			parent,
		); ok {

		t.Fatalf(
			"parent unexpectedly contains operation ID %q",
			got,
		)
	}

	got, ok :=
		OperationIDFromContext(
			child,
		)

	if !ok ||
		got !=
			"op-child-123" {

		t.Fatalf(
			"child operation ID = %q, %v; want %q",
			got,
			ok,
			"op-child-123",
		)
	}
}

func TestBindOperationIDAllowsSameIdentity(
	t *testing.T,
) {
	t.Parallel()

	ctx, err :=
		BindOperationID(
			context.Background(),
			"op-same-123",
		)

	if err != nil {
		t.Fatalf(
			"first bind error = %v",
			err,
		)
	}

	rebound, err :=
		BindOperationID(
			ctx,
			"op-same-123",
		)

	if err != nil {
		t.Fatalf(
			"same-ID rebind error = %v",
			err,
		)
	}

	got, ok :=
		OperationIDFromContext(
			rebound,
		)

	if !ok ||
		got !=
			"op-same-123" {

		t.Fatalf(
			"operation ID = %q, %v; want %q",
			got,
			ok,
			"op-same-123",
		)
	}
}

func TestBindOperationIDRejectsDifferentIdentity(
	t *testing.T,
) {
	t.Parallel()

	ctx, err :=
		BindOperationID(
			context.Background(),
			"op-original-123",
		)

	if err != nil {
		t.Fatalf(
			"first bind error = %v",
			err,
		)
	}

	rebound, err :=
		BindOperationID(
			ctx,
			"op-different-456",
		)

	if !errors.Is(
		err,
		ErrOperationIDContextConflict,
	) {
		t.Fatalf(
			"rebind error = %v, want %v",
			err,
			ErrOperationIDContextConflict,
		)
	}

	got, ok :=
		OperationIDFromContext(
			rebound,
		)

	if !ok {
		t.Fatal(
			"original operation ID disappeared",
		)
	}

	if got !=
		"op-original-123" {

		t.Fatalf(
			"operation ID = %q, want original %q",
			got,
			"op-original-123",
		)
	}
}

func TestBindOperationIDRejectsInvalidValue(
	t *testing.T,
) {
	t.Parallel()

	parent :=
		context.Background()

	ctx, err :=
		BindOperationID(
			parent,
			" unsafe operation id ",
		)

	if !errors.Is(
		err,
		ErrInvalidOperationID,
	) {
		t.Fatalf(
			"error = %v, want %v",
			err,
			ErrInvalidOperationID,
		)
	}

	if got, ok :=
		OperationIDFromContext(
			ctx,
		); ok {

		t.Fatalf(
			"invalid operation ID reached context: %q",
			got,
		)
	}
}

func TestBindOperationIDDoesNotNormalize(
	t *testing.T,
) {
	t.Parallel()

	_, err :=
		BindOperationID(
			context.Background(),
			"  op-valid-but-not-normalized  ",
		)

	if !errors.Is(
		err,
		ErrInvalidOperationID,
	) {
		t.Fatalf(
			"error = %v, want %v",
			err,
			ErrInvalidOperationID,
		)
	}
}

func TestBindOperationIDPreservesDeadline(
	t *testing.T,
) {
	t.Parallel()

	deadline :=
		time.Now().
			Add(
				time.Minute,
			)

	parent, cancel :=
		context.WithDeadline(
			context.Background(),
			deadline,
		)

	defer cancel()

	ctx, err :=
		BindOperationID(
			parent,
			"op-deadline-123",
		)

	if err != nil {
		t.Fatalf(
			"BindOperationID() error = %v",
			err,
		)
	}

	got, ok :=
		ctx.Deadline()

	if !ok {
		t.Fatal(
			"operation ID context lost deadline",
		)
	}

	if !got.Equal(
		deadline,
	) {
		t.Fatalf(
			"deadline = %v, want %v",
			got,
			deadline,
		)
	}
}

func TestOperationIDContextDoesNotUseStringKey(
	t *testing.T,
) {
	t.Parallel()

	ctx :=
		context.WithValue(
			context.Background(),
			"operation_id",
			"op-spoofed-123",
		)

	if got, ok :=
		OperationIDFromContext(
			ctx,
		); ok {

		t.Fatalf(
			"string-key operation ID leaked into canonical context: %q",
			got,
		)
	}
}

func TestBindOperationIDIgnoresUnrelatedStringKey(
	t *testing.T,
) {
	t.Parallel()

	parent :=
		context.WithValue(
			context.Background(),
			"operation_id",
			"op-string-key-123",
		)

	ctx, err :=
		BindOperationID(
			parent,
			"op-canonical-456",
		)

	if err != nil {
		t.Fatalf(
			"BindOperationID() error = %v",
			err,
		)
	}

	got, ok :=
		OperationIDFromContext(
			ctx,
		)

	if !ok ||
		got !=
			"op-canonical-456" {

		t.Fatalf(
			"canonical operation ID = %q, %v; want %q",
			got,
			ok,
			"op-canonical-456",
		)
	}
}

func TestBindOperationIDPanicsOnNilContext(
	t *testing.T,
) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal(
				"BindOperationID(nil, ...) did not panic",
			)
		}
	}()

	_, _ =
		BindOperationID(
			nil,
			"op-nil-123",
		)
}

func TestOperationIDFromNilContextIsMissing(
	t *testing.T,
) {
	t.Parallel()

	got, ok :=
		OperationIDFromContext(
			nil,
		)

	if ok ||
		got != "" {

		t.Fatalf(
			"OperationIDFromContext(nil) = %q, %v; want empty,false",
			got,
			ok,
		)
	}
}
