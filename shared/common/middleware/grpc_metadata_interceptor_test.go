package _middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_callercontext "common/callercontext"
	_enums "common/domain/enum"
	_identitycontext "common/identitycontext"
	_jwt "common/jwt"
	_requestcontext "common/requestcontext"
	_transporttrust "common/transporttrust"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestInjectGrpcMetadataIncludesCanonicalIDs(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "http://example.test/v2/tqd/public/enums", nil)
	req.Header.Set(_requestcontext.RequestIDHeader, "req-http-123")
	req.Header.Set(_requestcontext.OperationIDHeader, "op-http-456")

	md := InjectGrpcMetadataContextMiddleware(context.Background(), req)

	if got := md.Get(_requestcontext.RequestIDMetadataKey); len(got) != 1 || got[0] != "req-http-123" {
		t.Fatalf("request metadata = %v, want [req-http-123]", got)
	}
	if got := md.Get(_requestcontext.OperationIDMetadataKey); len(got) != 1 || got[0] != "op-http-456" {
		t.Fatalf("operation metadata = %v, want [op-http-456]", got)
	}
}

func TestParseGrpcMetadataStoresCanonicalIDs(t *testing.T) {
	t.Parallel()

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(
			_requestcontext.RequestIDMetadataKey, "req-service-123",
			_requestcontext.OperationIDMetadataKey, "op-service-456",
		),
	)

	_, err := ParseGrpcMetadataContextMiddleware(
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
		func(handlerCtx context.Context, req interface{}) (interface{}, error) {
			requestID, requestOK := _requestcontext.FromContext(handlerCtx)
			if !requestOK || requestID != "req-service-123" {
				t.Fatalf("request context = %q, %v", requestID, requestOK)
			}

			operationID, operationOK := _requestcontext.OperationIDFromContext(handlerCtx)
			if !operationOK || operationID != "op-service-456" {
				t.Fatalf("operation context = %q, %v", operationID, operationOK)
			}

			return nil, nil
		},
	)
	if err != nil {
		t.Fatalf("ParseGrpcMetadataContextMiddleware() error = %v", err)
	}
}

func TestParseGrpcMetadataRejectsInvalidOperationID(t *testing.T) {
	t.Parallel()

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(
			_requestcontext.OperationIDMetadataKey,
			"unsafe operation id",
		),
	)

	_, err := ParseGrpcMetadataContextMiddleware(
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
		func(handlerCtx context.Context, req interface{}) (interface{}, error) {
			if operationID, ok := _requestcontext.OperationIDFromContext(handlerCtx); ok || operationID != "" {
				t.Fatalf("invalid operation ID reached context: %q", operationID)
			}
			return nil, nil
		},
	)
	if err != nil {
		t.Fatalf("ParseGrpcMetadataContextMiddleware() error = %v", err)
	}
}

func TestGrpcClientMetadataFromContextIncludesCanonicalIDs(t *testing.T) {
	t.Parallel()

	ctx := _requestcontext.WithRequestID(context.Background(), "req-client-123")
	ctx = _requestcontext.WithOperationID(ctx, "op-client-456")

	md := GrpcClientMetadataFromContext(ctx)
	if got := md.Get(_requestcontext.RequestIDMetadataKey); len(got) != 1 || got[0] != "req-client-123" {
		t.Fatalf("request metadata = %v, want [req-client-123]", got)
	}
	if got := md.Get(_requestcontext.OperationIDMetadataKey); len(got) != 1 || got[0] != "op-client-456" {
		t.Fatalf("operation metadata = %v, want [op-client-456]", got)
	}
}

func TestInjectGrpcMetadataIncludesOptionalIdempotencyKey(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/reports", nil)
	req.Header.Set(_requestcontext.IdempotencyKeyHeader, "idem-http-123")

	md := InjectGrpcMetadataContextMiddleware(context.Background(), req)
	if got := md.Get(_requestcontext.IdempotencyKeyMetadataKey); len(got) != 1 || got[0] != "idem-http-123" {
		t.Fatalf("idempotency metadata = %v, want [idem-http-123]", got)
	}
}

func TestInjectGrpcMetadataOmitsMissingIdempotencyKey(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/resource", nil)
	md := InjectGrpcMetadataContextMiddleware(context.Background(), req)
	if got := md.Get(_requestcontext.IdempotencyKeyMetadataKey); len(got) != 0 {
		t.Fatalf("idempotency metadata = %v, want empty", got)
	}
}

func TestParseGrpcMetadataStoresIdempotencyKey(t *testing.T) {
	t.Parallel()

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(_requestcontext.IdempotencyKeyMetadataKey, "idem-service-123"),
	)

	_, err := ParseGrpcMetadataContextMiddleware(
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
		func(handlerCtx context.Context, req interface{}) (interface{}, error) {
			key, ok := _requestcontext.IdempotencyKeyFromContext(handlerCtx)
			if !ok || key != "idem-service-123" {
				t.Fatalf("idempotency context = %q, %v", key, ok)
			}
			return nil, nil
		},
	)
	if err != nil {
		t.Fatalf("ParseGrpcMetadataContextMiddleware() error = %v", err)
	}
}

func TestParseGrpcMetadataRejectsInvalidIdempotencyKey(t *testing.T) {
	t.Parallel()

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(_requestcontext.IdempotencyKeyMetadataKey, "unsafe key"),
	)

	called := false
	_, err := ParseGrpcMetadataContextMiddleware(
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
		func(handlerCtx context.Context, req interface{}) (interface{}, error) {
			called = true
			return nil, nil
		},
	)
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("status code = %v, want InvalidArgument", status.Code(err))
	}
	if called {
		t.Fatal("handler was called with invalid idempotency key")
	}
}

func TestGrpcClientMetadataFromContextIncludesIdempotencyKey(t *testing.T) {
	t.Parallel()

	ctx := _requestcontext.WithIdempotencyKey(context.Background(), "idem-client-123")
	md := GrpcClientMetadataFromContext(ctx)
	if got := md.Get(_requestcontext.IdempotencyKeyMetadataKey); len(got) != 1 || got[0] != "idem-client-123" {
		t.Fatalf("idempotency metadata = %v, want [idem-client-123]", got)
	}
}

func TestParseGrpcMetadataPreservesCancellation(t *testing.T) {
	t.Parallel()

	parent, cancel := context.WithCancel(context.Background())
	ctx := metadata.NewIncomingContext(
		parent,
		metadata.Pairs(_requestcontext.RequestIDMetadataKey, "req-canceled-123"),
	)
	cancel()

	_, err := ParseGrpcMetadataContextMiddleware(
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
		func(handlerCtx context.Context, req interface{}) (interface{}, error) {
			if !errors.Is(handlerCtx.Err(), context.Canceled) {
				t.Fatalf("context error = %v, want context.Canceled", handlerCtx.Err())
			}
			return nil, handlerCtx.Err()
		},
	)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("ParseGrpcMetadataContextMiddleware() error = %v, want context.Canceled", err)
	}
}

func TestParseGrpcMetadataPreservesDeadline(t *testing.T) {
	t.Parallel()

	deadline := time.Now().Add(time.Minute)
	parent, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	ctx := metadata.NewIncomingContext(
		parent,
		metadata.Pairs(_requestcontext.RequestIDMetadataKey, "req-deadline-123"),
	)

	_, err := ParseGrpcMetadataContextMiddleware(
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
		func(handlerCtx context.Context, req interface{}) (interface{}, error) {
			got, ok := handlerCtx.Deadline()
			if !ok {
				t.Fatal("deadline is missing")
			}
			if !got.Equal(deadline) {
				t.Fatalf("deadline = %v, want %v", got, deadline)
			}
			return nil, nil
		},
	)
	if err != nil {
		t.Fatalf("ParseGrpcMetadataContextMiddleware() error = %v", err)
	}
}

func TestInjectGrpcMetadataUsesVerifiedRequestPrincipal(t *testing.T) {
	t.Parallel()

	organizationID := uint64(19)
	principal := &_jwt.Principal{
		AuthID:         11,
		ProfileId:      12,
		OriginId:       13,
		OrganizationId: &organizationID,
		Session:        14,
		Role:           "user",
		Type:           "ACCESS",
	}

	req := httptest.NewRequest(http.MethodGet, "http://example.test/v2/resource", nil)
	req = req.WithContext(_jwt.WithPrincipal(req.Context(), principal))

	md := InjectGrpcMetadataContextMiddleware(context.Background(), req)
	if got := md.Get(_identitycontext.AuthIDMetadataKey); len(got) != 1 || got[0] != "11" {
		t.Fatalf("auth metadata = %v, want [11]", got)
	}
	if got := md.Get(_identitycontext.ProfileIDMetadataKey); len(got) != 1 || got[0] != "12" {
		t.Fatalf("profile metadata = %v, want [12]", got)
	}
	if got := md.Get(_identitycontext.OrganizationIDMetadataKey); len(got) != 1 || got[0] != "19" {
		t.Fatalf("organization metadata = %v, want [19]", got)
	}
}

func TestParseGrpcMetadataStoresCanonicalActorAndLegacyProjection(t *testing.T) {
	t.Parallel()

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(
			_identitycontext.AuthIDMetadataKey, "11",
			_identitycontext.ProfileIDMetadataKey, "12",
			_identitycontext.OriginIDMetadataKey, "13",
			_identitycontext.OrganizationIDMetadataKey, "19",
			_identitycontext.SessionIDMetadataKey, "14",
			_identitycontext.RoleMetadataKey, "user",
			_identitycontext.TokenTypeMetadataKey, "ACCESS",
		),
	)

	_, err := ParseGrpcMetadataContextMiddleware(
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
		func(handlerCtx context.Context, req interface{}) (interface{}, error) {
			actor, ok := _identitycontext.FromContext(handlerCtx)
			if !ok {
				t.Fatal("canonical actor is missing")
			}
			if actor.AuthID != 11 || actor.ProfileID != 12 || actor.OriginID != 13 {
				t.Fatalf("actor IDs = %+v", actor)
			}
			if actor.OrganizationID == nil || *actor.OrganizationID != 19 {
				t.Fatalf("organization = %v, want 19", actor.OrganizationID)
			}
			if got := handlerCtx.Value(_enums.ProfileIDKey); got != uint64(12) {
				t.Fatalf("legacy profile = %v, want 12", got)
			}
			return nil, nil
		},
	)
	if err != nil {
		t.Fatalf("ParseGrpcMetadataContextMiddleware() error = %v", err)
	}
}

func TestParseGrpcMetadataRejectsDuplicateActorField(t *testing.T) {
	t.Parallel()

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.MD{
			_identitycontext.ProfileIDMetadataKey: {"12", "13"},
		},
	)

	called := false
	_, err := ParseGrpcMetadataContextMiddleware(
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
		func(handlerCtx context.Context, req interface{}) (interface{}, error) {
			called = true
			return nil, nil
		},
	)
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("status code = %v, want Unauthenticated", status.Code(err))
	}
	if called {
		t.Fatal("handler was called with duplicate identity metadata")
	}
}

func TestGrpcClientMetadataPrefersCanonicalActor(t *testing.T) {
	t.Parallel()

	actor := _identitycontext.Actor{AuthID: 11, ProfileID: 12, Role: "user"}
	ctx := _identitycontext.WithActor(context.Background(), actor)
	ctx = context.WithValue(ctx, _enums.ProfileIDKey, uint64(99))

	md := GrpcClientMetadataFromContext(ctx)
	if got := md.Get(_identitycontext.ProfileIDMetadataKey); len(got) != 1 || got[0] != "12" {
		t.Fatalf("profile metadata = %v, want [12]", got)
	}
}

func TestPrincipalFromContextProjectsCanonicalActor(t *testing.T) {
	t.Parallel()

	organizationID := uint64(19)
	ctx := _identitycontext.WithActor(context.Background(), _identitycontext.Actor{
		AuthID:         11,
		ProfileID:      12,
		OriginID:       13,
		OrganizationID: &organizationID,
		SessionID:      14,
		Role:           "user",
		TokenType:      "ACCESS",
	})
	ctx = context.WithValue(ctx, _enums.PlanIDKey, uint64(21))

	principal, err := PrincipalFromContext(ctx)
	if err != nil {
		t.Fatalf("PrincipalFromContext() error = %v", err)
	}
	if principal.AuthID != 11 || principal.ProfileId != 12 || principal.OriginId != 13 {
		t.Fatalf("principal = %+v", principal)
	}
	if principal.PlanId == nil || *principal.PlanId != 21 {
		t.Fatalf("plan ID = %v, want 21", principal.PlanId)
	}
}

func TestParseGrpcMetadataEnforceRejectsUnsignedActor(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := _transporttrust.Config{
		Secret: "test-secret",
		Mode:   _transporttrust.ModeEnforce,
		Now:    func() time.Time { return now },
	}
	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(_identitycontext.ProfileIDMetadataKey, "12"),
	)

	called := false
	_, err := parseGrpcMetadataContextMiddlewareWithTrust(
		cfg,
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
		func(handlerCtx context.Context, req interface{}) (interface{}, error) {
			called = true
			return nil, nil
		},
	)
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("status code = %v, want Unauthenticated", status.Code(err))
	}
	if called {
		t.Fatal("handler was called with unsigned actor metadata")
	}
}

func TestParseGrpcMetadataEnforceAcceptsVerifiedActor(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := _transporttrust.Config{
		ServiceID: "gateway-service",
		Secret:    "test-secret",
		Mode:      _transporttrust.ModeEnforce,
		MaxAge:    30 * time.Second,
		Now:       func() time.Time { return now },
	}
	md := metadata.Pairs(
		_identitycontext.ProfileIDMetadataKey, "12",
		_callercontext.CallerKindMetadataKey, string(_callercontext.KindUser),
	)
	if err := _transporttrust.Sign("/test.Service/Call", md, cfg); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := parseGrpcMetadataContextMiddlewareWithTrust(
		cfg,
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
		func(handlerCtx context.Context, req interface{}) (interface{}, error) {
			if !_transporttrust.IsVerified(handlerCtx) {
				t.Fatal("transport caller is not verified")
			}
			actor, ok := _identitycontext.FromContext(handlerCtx)
			if !ok || actor.ProfileID != 12 {
				t.Fatalf("actor = %+v, %v", actor, ok)
			}
			return nil, nil
		},
	)
	if err != nil {
		t.Fatalf("middleware error = %v", err)
	}
}

func TestParseGrpcMetadataAuditMarksLegacyActorUnverified(t *testing.T) {
	cfg := _transporttrust.Config{Mode: _transporttrust.ModeAudit}
	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(_identitycontext.ProfileIDMetadataKey, "12"),
	)

	_, err := parseGrpcMetadataContextMiddlewareWithTrust(
		cfg,
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
		func(handlerCtx context.Context, req interface{}) (interface{}, error) {
			caller, ok := _transporttrust.FromContext(handlerCtx)
			if !ok || caller.Level != _transporttrust.LevelLegacyUnverified {
				t.Fatalf("caller = %+v, %v", caller, ok)
			}
			if _transporttrust.IsVerified(handlerCtx) {
				t.Fatal("legacy metadata was marked verified")
			}
			return nil, nil
		},
	)
	if err != nil {
		t.Fatalf("middleware error = %v", err)
	}
}

func TestParseGrpcMetadataEnforceRejectsUnsignedDirectAnonymous(t *testing.T) {
	cfg := _transporttrust.Config{
		Secret: "test-secret",
		Mode:   _transporttrust.ModeEnforce,
	}
	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(_requestcontext.RequestIDMetadataKey, "req-public-123"),
	)

	called := false
	_, err := parseGrpcMetadataContextMiddlewareWithTrust(
		cfg,
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.PublicService/Get"},
		func(handlerCtx context.Context, req interface{}) (interface{}, error) {
			called = true
			return nil, nil
		},
	)
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("status code = %v, want Unauthenticated", status.Code(err))
	}
	if called {
		t.Fatal("unsigned direct anonymous handler was called")
	}
}

func TestInjectGrpcMetadataIncludesCallerClassification(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "http://example.test/v2/resource", nil)
	req.Header.Set(_jwt.APIKeyHeader, "0123456789abcdef")
	req = req.WithContext(_callercontext.WithCaller(req.Context(), _callercontext.Caller{
		Kind:      _callercontext.KindAPIKey,
		APIKeyID:  7,
		APIKeyApp: "release-uploader",
	}))

	md := InjectGrpcMetadataContextMiddleware(context.Background(), req)
	caller, ok, err := _callercontext.FromMetadata(md)
	if err != nil || !ok {
		t.Fatalf("caller metadata = %#v, %v, %v", caller, ok, err)
	}
	if caller.Kind != _callercontext.KindAPIKey || caller.APIKeyID != 7 || caller.APIKeyApp != "release-uploader" {
		t.Fatalf("caller = %#v", caller)
	}
	if got := md.Get("x-api-key"); len(got) != 1 || got[0] != "0123456789abcdef" {
		t.Fatalf("x-api-key metadata = %v", got)
	}
}

func TestParseGrpcMetadataStoresVerifiedAnonymousCaller(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := _transporttrust.Config{
		ServiceID: "gateway-service",
		Secret:    "test-secret",
		Mode:      _transporttrust.ModeEnforce,
		MaxAge:    30 * time.Second,
		Now:       func() time.Time { return now },
	}
	md := metadata.Pairs(
		_callercontext.CallerKindMetadataKey, string(_callercontext.KindAnonymous),
	)
	if err := _transporttrust.Sign("/test.PublicService/Get", md, cfg); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	ctx := metadata.NewIncomingContext(context.Background(), md)
	_, err := parseGrpcMetadataContextMiddlewareWithTrust(
		cfg,
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.PublicService/Get"},
		func(handlerCtx context.Context, req interface{}) (interface{}, error) {
			caller, ok := _callercontext.FromContext(handlerCtx)
			if !ok || caller.Kind != _callercontext.KindAnonymous {
				t.Fatalf("caller = %#v, %v", caller, ok)
			}
			if !_transporttrust.IsVerified(handlerCtx) {
				t.Fatal("anonymous gateway caller is not transport verified")
			}
			return nil, nil
		},
	)
	if err != nil {
		t.Fatalf("middleware error = %v", err)
	}
}

func TestParseGrpcMetadataRejectsUserCallerWithoutActor(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := _transporttrust.Config{
		ServiceID: "gateway-service",
		Secret:    "test-secret",
		Mode:      _transporttrust.ModeEnforce,
		Now:       func() time.Time { return now },
	}
	md := metadata.Pairs(
		_callercontext.CallerKindMetadataKey, string(_callercontext.KindUser),
	)
	if err := _transporttrust.Sign("/test.Service/Call", md, cfg); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	ctx := metadata.NewIncomingContext(context.Background(), md)
	_, err := parseGrpcMetadataContextMiddlewareWithTrust(
		cfg,
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
		func(context.Context, interface{}) (interface{}, error) { return nil, nil },
	)
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("status code = %v, want Unauthenticated", status.Code(err))
	}
}

func TestGrpcClientMetadataPreservesExistingGatewayCaller(t *testing.T) {
	t.Parallel()

	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs(
			_callercontext.CallerKindMetadataKey,
			string(_callercontext.KindAnonymous),
		),
	)
	md := GrpcClientMetadataFromContext(ctx)
	if got := md.Get(_callercontext.CallerKindMetadataKey); len(got) != 0 {
		t.Fatalf("canonical metadata overwrote existing gateway caller: %v", got)
	}
}

type metadataTestServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *metadataTestServerStream) Context() context.Context {
	return s.ctx
}

func TestStreamMetadataInterceptorStoresCallerContext(t *testing.T) {
	t.Setenv("QHPRO_SERVICE_ID", "gateway-service")
	t.Setenv("QHPRO_INTERNAL_METADATA_SECRET", "test-secret")
	t.Setenv("QHPRO_TRUSTED_METADATA_MODE", "enforce")

	cfg := _transporttrust.ConfigFromEnv()
	md := metadata.Pairs(
		_callercontext.CallerKindMetadataKey, string(_callercontext.KindAnonymous),
	)
	if err := _transporttrust.Sign("/test.StreamService/Watch", md, cfg); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	stream := &metadataTestServerStream{
		ctx: metadata.NewIncomingContext(context.Background(), md),
	}
	called := false
	err := ParseGrpcMetadataContextStreamMiddleware(
		nil,
		stream,
		&grpc.StreamServerInfo{FullMethod: "/test.StreamService/Watch"},
		func(_ interface{}, handlerStream grpc.ServerStream) error {
			called = true
			caller, ok := _callercontext.FromContext(handlerStream.Context())
			if !ok || caller.Kind != _callercontext.KindAnonymous {
				t.Fatalf("caller = %#v, %v", caller, ok)
			}
			return nil
		},
	)
	if err != nil {
		t.Fatalf("stream interceptor error = %v", err)
	}
	if !called {
		t.Fatal("stream handler was not called")
	}
}
