package rpc

import (
	"common/identity"
	"common/request"
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func transportTestConfig(now time.Time, requireServiceAssertion bool) TransportConfig {
	return TransportConfig{
		ServiceAssertion: ServiceAssertionConfig{
			ServiceID:                  "gateway-service",
			Secret:                     "test-secret",
			VerificationSecrets:        map[string]string{"gateway-service": "test-secret"},
			VerificationKeysConfigured: true,
			MaxAge:                     30 * time.Second,
			ClockSkew:                  5 * time.Second,
			Now:                        func() time.Time { return now },
		},
		RequireServiceAssertion: requireServiceAssertion,
	}
}

func TestWithOutgoingContextPreservesExistingAndSignsCanonicalMetadata(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
		"x-custom", "keep",
		RequestIDMetadataKey, "stale",
	))
	ctx = request.WithRequestID(ctx, "req_123")
	var err error
	ctx, err = identity.BindCaller(ctx, identity.Caller{Kind: identity.CallerAnonymous})
	if err != nil {
		t.Fatal(err)
	}

	out, err := WithOutgoingContext(ctx, "/test.Service/Call", transportTestConfig(now, true))
	if err != nil {
		t.Fatal(err)
	}
	md, ok := metadata.FromOutgoingContext(out)
	if !ok {
		t.Fatal("outgoing metadata missing")
	}
	if got := md.Get("x-custom"); len(got) != 1 || got[0] != "keep" {
		t.Fatalf("custom = %v", got)
	}
	if got := md.Get(RequestIDMetadataKey); len(got) != 1 || got[0] != "req_123" {
		t.Fatalf("request = %v", got)
	}
	if _, err := VerifyServiceAssertion("/test.Service/Call", md, transportTestConfig(now, true).ServiceAssertion); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
}

func TestPrepareIncomingContextVerifiedRestoresCanonicalState(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := transportTestConfig(now, true)
	md := metadata.Pairs(
		RequestIDMetadataKey, "req_123",
		CallerKindMetadataKey, string(identity.CallerUser),
		ProfileIDMetadataKey, "42",
	)
	if err := SignServiceAssertion("/test.Service/Call", md, cfg.ServiceAssertion); err != nil {
		t.Fatal(err)
	}
	ctx := metadata.NewIncomingContext(context.Background(), md)
	bound, err := PrepareIncomingContext(ctx, "/test.Service/Call", cfg)
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := request.RequestIDFromContext(bound); !ok || got != "req_123" {
		t.Fatalf("request = %q, %v", got, ok)
	}
	actor, ok := identity.ActorFromContext(bound)
	if !ok || actor.ProfileID != 42 {
		t.Fatalf("actor = %#v, %v", actor, ok)
	}
	caller, ok := identity.CallerFromContext(bound)
	if !ok || caller.Kind != identity.CallerUser {
		t.Fatalf("caller = %#v, %v", caller, ok)
	}
	if _, ok := identity.ServiceCallerFromContext(bound); !ok {
		t.Fatal("transport caller is not verified")
	}
}

func TestPrepareIncomingContextAuditDoesNotPromoteUnsignedIdentity(t *testing.T) {
	cfg := transportTestConfig(time.Unix(1_700_000_000, 0), false)
	cfg.ServiceAssertion.Secret = ""
	cfg.ServiceAssertion.VerificationSecrets = nil
	cfg.ServiceAssertion.VerificationKeysConfigured = false
	md := metadata.Pairs(
		RequestIDMetadataKey, "req_123",
		CallerKindMetadataKey, string(identity.CallerUser),
		ProfileIDMetadataKey, "42",
	)
	bound, err := PrepareIncomingContext(metadata.NewIncomingContext(context.Background(), md), "/test.Service/Call", cfg)
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := request.RequestIDFromContext(bound); !ok || got != "req_123" {
		t.Fatalf("request = %q, %v", got, ok)
	}
	if _, ok := identity.ActorFromContext(bound); ok {
		t.Fatal("unsigned actor was promoted")
	}
	if caller, ok := identity.ServiceCallerFromContext(bound); ok {
		t.Fatalf("unverified transport metadata created verified service caller: %#v", caller)
	}
}

func TestPrepareIncomingContextEnforceRejectsUnsignedBusinessRPC(t *testing.T) {
	cfg := transportTestConfig(time.Unix(1_700_000_000, 0), true)
	_, err := PrepareIncomingContext(
		metadata.NewIncomingContext(context.Background(), metadata.MD{}),
		"/test.Service/Call",
		cfg,
	)
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("code = %v, err = %v", status.Code(err), err)
	}
}

func TestPrepareIncomingContextAllowsUnsignedHealthOnlyWithoutPrivilege(t *testing.T) {
	cfg := transportTestConfig(time.Unix(1_700_000_000, 0), true)
	if _, err := PrepareIncomingContext(
		metadata.NewIncomingContext(context.Background(), metadata.MD{}),
		"/grpc.health.v1.Health/Check", cfg,
	); err != nil {
		t.Fatalf("health error = %v", err)
	}
	_, err := PrepareIncomingContext(
		metadata.NewIncomingContext(context.Background(), metadata.Pairs(ProfileIDMetadataKey, "42")),
		"/grpc.health.v1.Health/Check", cfg,
	)
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("smuggling code = %v", status.Code(err))
	}
}

func TestUnaryClientInterceptorUsesCanonicalOutgoingContext(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	ctx, err := identity.BindCaller(context.Background(), identity.Caller{Kind: identity.CallerAnonymous})
	if err != nil {
		t.Fatal(err)
	}
	called := false
	err = UnaryClientInterceptor(transportTestConfig(now, true))(
		ctx, "/test.Service/Call", nil, nil, nil,
		func(callCtx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			called = true
			md, ok := metadata.FromOutgoingContext(callCtx)
			if !ok {
				t.Fatal("metadata missing")
			}
			if _, verifyErr := VerifyServiceAssertion(method, md, transportTestConfig(now, true).ServiceAssertion); verifyErr != nil {
				t.Fatalf("verify error = %v", verifyErr)
			}
			return nil
		},
	)
	if err != nil || !called {
		t.Fatalf("err=%v called=%v", err, called)
	}
}

type testServerStream struct {
	ctx context.Context
}

func (s *testServerStream) SetHeader(metadata.MD) error {
	return nil
}

func (s *testServerStream) SendHeader(metadata.MD) error {
	return nil
}

func (s *testServerStream) SetTrailer(metadata.MD) {}

func (s *testServerStream) Context() context.Context {
	return s.ctx
}

func (s *testServerStream) SendMsg(interface{}) error {
	return nil
}

func (s *testServerStream) RecvMsg(interface{}) error {
	return nil
}

type testClientStream struct {
	ctx context.Context
}

func (s *testClientStream) Header() (metadata.MD, error) {
	return metadata.MD{}, nil
}

func (s *testClientStream) Trailer() metadata.MD {
	return metadata.MD{}
}

func (s *testClientStream) CloseSend() error {
	return nil
}

func (s *testClientStream) Context() context.Context {
	return s.ctx
}

func (s *testClientStream) SendMsg(interface{}) error {
	return nil
}

func (s *testClientStream) RecvMsg(interface{}) error {
	return nil
}

func TestStreamClientInterceptorUsesSameSignedMetadataContract(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	ctx, err := identity.BindCaller(context.Background(), identity.Caller{Kind: identity.CallerAnonymous})
	if err != nil {
		t.Fatal(err)
	}
	called := false
	stream, err := StreamClientInterceptor(transportTestConfig(now, true))(
		ctx, &grpc.StreamDesc{}, nil, "/test.Service/Watch",
		func(callCtx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			called = true
			md, ok := metadata.FromOutgoingContext(callCtx)
			if !ok {
				t.Fatal("metadata missing")
			}
			if _, verifyErr := VerifyServiceAssertion(method, md, transportTestConfig(now, true).ServiceAssertion); verifyErr != nil {
				t.Fatalf("verify error = %v", verifyErr)
			}
			return &testClientStream{
				ctx: callCtx,
			}, nil
		},
	)
	if err != nil || !called || stream == nil {
		t.Fatalf("stream=%v err=%v called=%v", stream, err, called)
	}
}

func TestStreamServerInterceptorRestoresSameCanonicalContext(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := transportTestConfig(now, true)
	md := metadata.Pairs(
		RequestIDMetadataKey, "req_stream",
		CallerKindMetadataKey, string(identity.CallerUser),
		ProfileIDMetadataKey, "42",
	)
	if err := SignServiceAssertion("/test.Service/Watch", md, cfg.ServiceAssertion); err != nil {
		t.Fatal(err)
	}
	stream := &testServerStream{
		ctx: metadata.NewIncomingContext(context.Background(), md),
	}
	called := false
	err := StreamServerInterceptor(cfg)(
		nil, stream, &grpc.StreamServerInfo{FullMethod: "/test.Service/Watch"},
		func(_ interface{}, bound grpc.ServerStream) error {
			called = true
			if got, ok := request.RequestIDFromContext(bound.Context()); !ok || got != "req_stream" {
				t.Fatalf("request = %q, %v", got, ok)
			}
			actor, ok := identity.ActorFromContext(bound.Context())
			if !ok || actor.ProfileID != 42 {
				t.Fatalf("actor = %#v, %v", actor, ok)
			}
			return nil
		},
	)
	if err != nil || !called {
		t.Fatalf("err=%v called=%v", err, called)
	}
}
