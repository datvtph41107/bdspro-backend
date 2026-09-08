package operation

import (
	"context"
	"errors"
	"testing"
)

func TestBindStoresOperation(t *testing.T) {
	t.Parallel()

	code := Code("workspace.report.generate")

	ctx, err := Bind(context.Background(), code)
	if err != nil {
		t.Fatalf("Bind() error = %v", err)
	}

	got, ok := FromContext(ctx)
	if !ok || got != code {
		t.Fatalf("FromContext() = %q, %v", got, ok)
	}
}

func TestBindSameOperationIsIdempotent(t *testing.T) {
	t.Parallel()

	code := Code("workspace.report.generate")

	ctx, err := Bind(context.Background(), code)
	if err != nil {
		t.Fatal(err)
	}

	same, err := Bind(ctx, code)
	if err != nil {
		t.Fatalf("Bind() same operation error = %v", err)
	}
	if same != ctx {
		t.Fatal("same operation should keep the existing context")
	}
}

func TestBindRejectsDifferentOperation(t *testing.T) {
	t.Parallel()

	ctx, err := Bind(
		context.Background(),
		Code("workspace.report.generate"),
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Bind(
		ctx,
		Code("workspace.report.delete"),
	)
	if !errors.Is(err, ErrContextConflict) {
		t.Fatalf(
			"Bind() error = %v, want ErrContextConflict",
			err,
		)
	}
}

func TestBindRejectsInvalidOperation(t *testing.T) {
	t.Parallel()

	_, err := Bind(
		context.Background(),
		Code("report"),
	)
	if !errors.Is(err, ErrInvalidCode) {
		t.Fatalf(
			"Bind() error = %v, want ErrInvalidCode",
			err,
		)
	}
}
