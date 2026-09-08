package requestlog

import (
	"bytes"
	"common/identity"
	"common/operation"
	"common/request"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestFromContextAddsTraceFieldsWithoutLoggingIdempotencyKey(t *testing.T) {
	ctx := context.Background()
	ctx = request.WithRequestID(ctx, "req_123")

	var err error
	ctx, err = request.BindOperationID(ctx, "op_123")
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = request.BindIdempotencyKey(ctx, "idem_secret_123")
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = identity.BindActor(ctx, identity.Actor{ProfileID: 42})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = operation.Bind(ctx, operation.Code("workspace.report.generate"))
	if err != nil {
		t.Fatal(err)
	}

	var buffer bytes.Buffer
	base := slog.New(slog.NewTextHandler(&buffer, nil))
	FromContext(ctx, base).Info("request trace")
	output := buffer.String()

	for _, expected := range []string{
		"request_id=req_123",
		"operation_id=op_123",
		"idempotency_key_set=true",
		"profile_id=42",
		"business_operation=workspace.report.generate",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("log output %q missing %q", output, expected)
		}
	}
	if strings.Contains(output, "idem_secret_123") {
		t.Fatalf("log output leaked idempotency key: %q", output)
	}
}

func TestFromContextOmitsOperationFieldWhenCanonicalOperationIsAbsent(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer
	base := slog.New(slog.NewTextHandler(&buffer, nil))
	FromContext(context.Background(), base).Info("request without operation")

	if strings.Contains(buffer.String(), "business_operation=") {
		t.Fatalf(
			"log output %q contains an operation without canonical context",
			buffer.String(),
		)
	}
}
