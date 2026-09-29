package _middleware

import (
	_rpc "common/rpc"
	"context"
	"errors"
	"testing"
	"time"

	_enums "common/domain/enum"
	_request "common/request"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestUnaryClientInterceptorPropagatesCanonicalIDsAndPreservesMetadata(
	t *testing.T,
) {
	t.Parallel()

	ctx := _request.WithRequestID(
		context.Background(),
		"req-outbound-123",
	)
	ctx = mustBindOperationID(t, ctx, "op-outbound-456")
	ctx = mustBindIdempotencyKey(t, ctx, "idem-outbound-789")
	ctx = context.WithValue(ctx, _enums.ProfileIDKey, uint64(42))
	ctx = metadata.NewOutgoingContext(
		ctx,
		metadata.Pairs(
			"x-custom-metadata", "keep-me",
			_request.RequestIDMetadataKey, "stale-request-id",
			_request.OperationIDMetadataKey, "stale-operation-id",
			_request.IdempotencyKeyMetadataKey, "stale-idempotency-key",
		),
	)

	var captured metadata.MD
	invoker := func(
		callCtx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		opts ...grpc.CallOption,
	) error {
		var ok bool
		captured, ok = metadata.FromOutgoingContext(callCtx)
		if !ok {
			t.Fatal("outgoing metadata is missing")
		}
		return nil
	}

	err := UnaryClientInterceptor(
		ctx,
		"/test.Service/Call",
		nil,
		nil,
		nil,
		invoker,
	)
	if err != nil {
		t.Fatalf("UnaryClientInterceptor() error = %v", err)
	}

	if got := captured.Get("x-custom-metadata"); len(got) != 1 || got[0] != "keep-me" {
		t.Fatalf("custom metadata = %v, want [keep-me]", got)
	}

	requestIDs := captured.Get(_request.RequestIDMetadataKey)
	if len(requestIDs) != 1 || requestIDs[0] != "req-outbound-123" {
		t.Fatalf(
			"request metadata = %v, want [req-outbound-123]",
			requestIDs,
		)
	}

	operationIDs := captured.Get(_request.OperationIDMetadataKey)
	if len(operationIDs) != 1 || operationIDs[0] != "op-outbound-456" {
		t.Fatalf(
			"operation metadata = %v, want [op-outbound-456]",
			operationIDs,
		)
	}

	idempotencyKeys := captured.Get(_request.IdempotencyKeyMetadataKey)
	if len(idempotencyKeys) != 1 || idempotencyKeys[0] != "idem-outbound-789" {
		t.Fatalf(
			"idempotency metadata = %v, want [idem-outbound-789]",
			idempotencyKeys,
		)
	}

	profileIDs := captured.Get("profileId")
	if len(profileIDs) != 1 || profileIDs[0] != "42" {
		t.Fatalf("profile metadata = %v, want [42]", profileIDs)
	}
}

func TestUnaryClientInterceptorReturnsInvokerError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("downstream failure")
	invoker := func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		opts ...grpc.CallOption,
	) error {
		return wantErr
	}

	err := UnaryClientInterceptor(
		context.Background(),
		"/test.Service/Call",
		nil,
		nil,
		nil,
		invoker,
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf("UnaryClientInterceptor() error = %v, want %v", err, wantErr)
	}
}

func TestUnaryClientInterceptorPreservesCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	invoker := func(
		callCtx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		opts ...grpc.CallOption,
	) error {
		if !errors.Is(callCtx.Err(), context.Canceled) {
			t.Fatalf("context error = %v, want context.Canceled", callCtx.Err())
		}
		return callCtx.Err()
	}

	err := UnaryClientInterceptor(
		ctx,
		"/test.Service/Call",
		nil,
		nil,
		nil,
		invoker,
	)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("UnaryClientInterceptor() error = %v, want context.Canceled", err)
	}
}

func TestUnaryClientInterceptorPreservesDeadline(t *testing.T) {
	t.Parallel()

	deadline := time.Now().Add(time.Minute)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	invoker := func(
		callCtx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		opts ...grpc.CallOption,
	) error {
		got, ok := callCtx.Deadline()
		if !ok {
			t.Fatal("deadline is missing")
		}
		if !got.Equal(deadline) {
			t.Fatalf("deadline = %v, want %v", got, deadline)
		}
		return nil
	}

	if err := UnaryClientInterceptor(
		ctx,
		"/test.Service/Call",
		nil,
		nil,
		nil,
		invoker,
	); err != nil {
		t.Fatalf("UnaryClientInterceptor() error = %v", err)
	}
}

func TestUnaryClientInterceptorSignsCanonicalMetadata(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := _rpc.TransportConfig{
		ServiceAssertion: _rpc.ServiceAssertionConfig{
			ServiceID: "tqd-service",
			Secret:    "test-secret",

			MaxAge: 30 * time.Second,
			Now:    func() time.Time { return now },
		},
		RequireServiceAssertion: true,
	}
	ctx := context.WithValue(context.Background(), _enums.ProfileIDKey, uint64(42))

	var captured metadata.MD
	invoker := func(
		callCtx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		opts ...grpc.CallOption,
	) error {
		captured, _ = metadata.FromOutgoingContext(callCtx)
		return nil
	}

	if err := unaryClientInterceptorWithTrust(
		cfg,
		ctx,
		"/test.Service/Call",
		nil,
		nil,
		nil,
		invoker,
	); err != nil {
		t.Fatalf("interceptor error = %v", err)
	}

	caller, err := _rpc.VerifyServiceAssertion("/test.Service/Call", captured, cfg.ServiceAssertion)
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if caller.ServiceID != "tqd-service" {
		t.Fatalf("caller = %+v", caller)
	}
}

func TestUnaryClientInterceptorEnforceRequiresSignerConfig(t *testing.T) {
	called := false
	invoker := func(
		callCtx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		opts ...grpc.CallOption,
	) error {
		called = true
		return nil
	}

	err := unaryClientInterceptorWithTrust(
		_rpc.TransportConfig{RequireServiceAssertion: true},
		context.Background(),
		"/test.Service/Call",
		nil,
		nil,
		nil,
		invoker,
	)
	if status.Code(err) != codes.FailedPrecondition {
		t.Fatalf("status code = %v, want FailedPrecondition", status.Code(err))
	}
	if called {
		t.Fatal("network invoker was called without signer config")
	}
}

func TestStreamClientInterceptorSignsCallerMetadata(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := _rpc.TransportConfig{
		ServiceAssertion: _rpc.ServiceAssertionConfig{
			ServiceID: "gateway-service",
			Secret:    "test-secret",

			Now: func() time.Time { return now },
		},
		RequireServiceAssertion: true,
	}
	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("x-qhpro-caller-kind", "anonymous"),
	)

	var captured metadata.MD
	streamer := func(
		callCtx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		captured, _ = metadata.FromOutgoingContext(callCtx)
		return nil, nil
	}

	_, err := streamClientInterceptorWithTrust(
		cfg,
		ctx,
		&grpc.StreamDesc{ServerStreams: true},
		nil,
		"/test.StreamService/Watch",
		streamer,
	)
	if err != nil {
		t.Fatalf("stream interceptor error = %v", err)
	}
	if _, err := _rpc.VerifyServiceAssertion("/test.StreamService/Watch", captured, cfg.ServiceAssertion); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if got := captured.Get("x-qhpro-caller-kind"); len(got) != 1 || got[0] != "anonymous" {
		t.Fatalf("caller metadata = %v", got)
	}
}
