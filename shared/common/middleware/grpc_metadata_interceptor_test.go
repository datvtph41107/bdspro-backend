package _middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	_callercontext "common/callercontext"
	_enums "common/domain/enum"
	_identitycontext "common/identitycontext"
	_jwt "common/jwt"
	_requestcontext "common/requestcontext"
	_transporttrust "common/transporttrust"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func mustBindActor(ctx context.Context, actor _identitycontext.Actor) context.Context {
	bound, err := _identitycontext.BindActor(ctx, actor)
	if err != nil {
		panic(err)
	}
	return bound
}

func mustBindCaller(ctx context.Context, caller _callercontext.Caller) context.Context {
	bound, err := _callercontext.BindCaller(ctx, caller)
	if err != nil {
		panic(err)
	}
	return bound
}

func TestInjectGrpcMetadataIncludesCanonicalIDs(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "http://example.test/v2/tqd/public/enums", nil)
	ctx := _requestcontext.WithRequestID(req.Context(), "req-http-123")
	ctx = mustBindOperationID(t, ctx, "op-http-456")
	req = req.WithContext(ctx)

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
	const malformed = "unsafe operation id"

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(
			_requestcontext.OperationIDMetadataKey,
			malformed,
		),
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
		t.Fatal("handler was called with invalid operation identity")
	}
	if strings.Contains(status.Convert(err).Message(), malformed) {
		t.Fatalf("malformed operation identity leaked into status: %v", err)
	}
}

func TestParseGrpcMetadataAllowsMissingOperationIDWithoutOriginating(t *testing.T) {
	t.Parallel()

	called := false
	_, err := ParseGrpcMetadataContextMiddleware(
		metadata.NewIncomingContext(context.Background(), metadata.MD{}),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
		func(handlerCtx context.Context, _ interface{}) (interface{}, error) {
			called = true
			if operationID, ok := _requestcontext.OperationIDFromContext(handlerCtx); ok || operationID != "" {
				t.Fatalf("missing operation identity was originated: %q", operationID)
			}
			return nil, nil
		},
	)
	if err != nil || !called {
		t.Fatalf("error/called = %v/%v", err, called)
	}
}

func TestParseGrpcMetadataRejectsMultipleOperationIDs(t *testing.T) {
	t.Parallel()

	for _, values := range [][]string{
		{"op-first-123", "op-second-456"},
		{"op-repeat-123", "op-repeat-123"},
	} {
		values := values
		t.Run(strings.Join(values, "+"), func(t *testing.T) {
			t.Parallel()
			md := metadata.MD{}
			for _, value := range values {
				md.Append(_requestcontext.OperationIDMetadataKey, value)
			}
			called := false
			_, err := ParseGrpcMetadataContextMiddleware(
				metadata.NewIncomingContext(context.Background(), md),
				nil,
				&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
				func(context.Context, interface{}) (interface{}, error) {
					called = true
					return nil, nil
				},
			)
			if status.Code(err) != codes.InvalidArgument {
				t.Fatalf("status code = %v, want InvalidArgument", status.Code(err))
			}
			if called {
				t.Fatal("handler was called with multiple operation identities")
			}
		})
	}
}

func TestParseGrpcMetadataAcceptsSamePreboundOperationID(t *testing.T) {
	t.Parallel()

	ctx := mustBindOperationID(t, context.Background(), "op-same-123")
	ctx = metadata.NewIncomingContext(
		ctx,
		metadata.Pairs(_requestcontext.OperationIDMetadataKey, "op-same-123"),
	)
	called := false
	_, err := ParseGrpcMetadataContextMiddleware(
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
		func(handlerCtx context.Context, _ interface{}) (interface{}, error) {
			called = true
			operationID, ok := _requestcontext.OperationIDFromContext(handlerCtx)
			if !ok || operationID != "op-same-123" {
				t.Fatalf("operation identity = %q, %v", operationID, ok)
			}
			return nil, nil
		},
	)
	if err != nil || !called {
		t.Fatalf("error/called = %v/%v", err, called)
	}
}

func TestParseGrpcMetadataMapsOperationIDContextConflictToInternal(t *testing.T) {
	t.Parallel()

	ctx := mustBindOperationID(t, context.Background(), "op-existing-123")
	ctx = metadata.NewIncomingContext(
		ctx,
		metadata.Pairs(_requestcontext.OperationIDMetadataKey, "op-different-456"),
	)
	called := false
	_, err := ParseGrpcMetadataContextMiddleware(
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
		func(context.Context, interface{}) (interface{}, error) {
			called = true
			return nil, nil
		},
	)
	if status.Code(err) != codes.Internal {
		t.Fatalf("status code = %v, want Internal", status.Code(err))
	}
	if called {
		t.Fatal("handler was called after operation identity conflict")
	}
	if strings.Contains(status.Convert(err).Message(), "op-different-456") {
		t.Fatalf("conflicting operation identity leaked into status: %v", err)
	}
}

func TestGrpcClientMetadataFromContextIncludesCanonicalIDs(t *testing.T) {
	t.Parallel()

	ctx := _requestcontext.WithRequestID(context.Background(), "req-client-123")
	ctx = mustBindOperationID(t, ctx, "op-client-456")

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
	ctx := mustBindIdempotencyKey(t, req.Context(), "idem-http-123")
	req = req.WithContext(ctx)

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

func TestInjectGrpcMetadataDoesNotReparseUnboundIdentityHeaders(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/commands", nil)
	req.Header.Set(_requestcontext.RequestIDHeader, "req-raw-123")
	req.Header.Set(_requestcontext.OperationIDHeader, "op-raw-456")
	req.Header.Set(_requestcontext.IdempotencyKeyHeader, "idem-raw-789")

	md := InjectGrpcMetadataContextMiddleware(context.Background(), req)
	for _, key := range []string{
		_requestcontext.RequestIDMetadataKey,
		_requestcontext.OperationIDMetadataKey,
		_requestcontext.IdempotencyKeyMetadataKey,
	} {
		if values := md.Get(key); len(values) != 0 {
			t.Fatalf("unbound %s metadata = %v, want empty", key, values)
		}
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

func TestParseGrpcMetadataRejectsMultipleIdempotencyKeys(t *testing.T) {
	t.Parallel()

	for _, values := range [][]string{
		{"idem-first-123", "idem-second-456"},
		{"idem-repeat-123", "idem-repeat-123"},
	} {
		values := values
		t.Run(strings.Join(values, "+"), func(t *testing.T) {
			t.Parallel()
			md := metadata.MD{}
			for _, value := range values {
				md.Append(_requestcontext.IdempotencyKeyMetadataKey, value)
			}
			called := false
			_, err := ParseGrpcMetadataContextMiddleware(
				metadata.NewIncomingContext(context.Background(), md),
				nil,
				&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
				func(context.Context, interface{}) (interface{}, error) {
					called = true
					return nil, nil
				},
			)
			if status.Code(err) != codes.InvalidArgument {
				t.Fatalf("status code = %v, want InvalidArgument", status.Code(err))
			}
			if called {
				t.Fatal("handler was called with multiple idempotency identities")
			}
		})
	}
}

func TestParseGrpcMetadataAcceptsSamePreboundIdempotencyKey(t *testing.T) {
	t.Parallel()

	ctx := mustBindIdempotencyKey(t, context.Background(), "idem-same-123")
	ctx = metadata.NewIncomingContext(
		ctx,
		metadata.Pairs(_requestcontext.IdempotencyKeyMetadataKey, "idem-same-123"),
	)
	called := false
	_, err := ParseGrpcMetadataContextMiddleware(
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
		func(handlerCtx context.Context, _ interface{}) (interface{}, error) {
			called = true
			key, ok := _requestcontext.IdempotencyKeyFromContext(handlerCtx)
			if !ok || key != "idem-same-123" {
				t.Fatalf("idempotency identity = %q, %v", key, ok)
			}
			return nil, nil
		},
	)
	if err != nil || !called {
		t.Fatalf("error/called = %v/%v", err, called)
	}
}

func TestParseGrpcMetadataMapsIdempotencyContextConflictToInternal(t *testing.T) {
	t.Parallel()

	ctx := mustBindIdempotencyKey(t, context.Background(), "idem-existing-123")
	ctx = metadata.NewIncomingContext(
		ctx,
		metadata.Pairs(_requestcontext.IdempotencyKeyMetadataKey, "idem-different-456"),
	)
	called := false
	_, err := ParseGrpcMetadataContextMiddleware(
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
		func(context.Context, interface{}) (interface{}, error) {
			called = true
			return nil, nil
		},
	)
	if status.Code(err) != codes.Internal {
		t.Fatalf("status code = %v, want Internal", status.Code(err))
	}
	if called {
		t.Fatal("handler was called after idempotency identity conflict")
	}
	if strings.Contains(status.Convert(err).Message(), "idem-different-456") {
		t.Fatalf("conflicting idempotency identity leaked into status: %v", err)
	}
}

func TestGrpcClientMetadataFromContextIncludesIdempotencyKey(t *testing.T) {
	t.Parallel()

	ctx := mustBindIdempotencyKey(t, context.Background(), "idem-client-123")
	md := GrpcClientMetadataFromContext(ctx)
	if got := md.Get(_requestcontext.IdempotencyKeyMetadataKey); len(got) != 1 || got[0] != "idem-client-123" {
		t.Fatalf("idempotency metadata = %v, want [idem-client-123]", got)
	}
}

func mustBindOperationID(t *testing.T, ctx context.Context, operationID string) context.Context {
	t.Helper()
	bound, err := _requestcontext.BindOperationID(ctx, operationID)
	if err != nil {
		t.Fatalf("BindOperationID() error = %v", err)
	}
	return bound
}

func mustBindIdempotencyKey(t *testing.T, ctx context.Context, key string) context.Context {
	t.Helper()
	bound, err := _requestcontext.BindIdempotencyKey(ctx, key)
	if err != nil {
		t.Fatalf("BindIdempotencyKey() error = %v", err)
	}
	return bound
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

func TestInjectGrpcMetadataNeverParsesRawAuthorization(t *testing.T) {
	viper.Set("jwt.key-generate", "first-hop-test-secret")
	viper.Set("jwt.accessExpAfterMinutes", int64(5))
	t.Cleanup(viper.Reset)

	token, err := _jwt.GenerateToken(_jwt.JwtTokenProperties{
		ProfileID: 42,
		SessionID: 9,
		Type:      _jwt.AccessToken,
	})
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "http://example.test/v2/resource", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	md := InjectGrpcMetadataContextMiddleware(context.Background(), req)
	if got := md.Get(_identitycontext.ProfileIDMetadataKey); len(got) != 0 {
		t.Fatalf("raw Authorization became trusted Actor metadata: %v", got)
	}
}

func TestParseGrpcMetadataStoresCanonicalActorAndLegacyProjection(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := _transporttrust.Config{
		ServiceID: "gateway-service",
		Secret:    "test-secret",
		Mode:      _transporttrust.ModeEnforce,
		MaxAge:    30 * time.Second,
		Now:       func() time.Time { return now },
	}
	md := metadata.Pairs(
		_identitycontext.AuthIDMetadataKey, "11",
		_identitycontext.ProfileIDMetadataKey, "12",
		_identitycontext.OriginIDMetadataKey, "13",
		_identitycontext.OrganizationIDMetadataKey, "19",
		_identitycontext.SessionIDMetadataKey, "14",
		_identitycontext.RoleMetadataKey, "user",
		_identitycontext.TokenTypeMetadataKey, "ACCESS",
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

func TestParseGrpcMetadataMapsActorContextConflictToInternal(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := _transporttrust.Config{
		ServiceID: "gateway-service",
		Secret:    "test-secret",
		Mode:      _transporttrust.ModeEnforce,
		MaxAge:    30 * time.Second,
		Now:       func() time.Time { return now },
	}
	md := metadata.Pairs(
		_identitycontext.ProfileIDMetadataKey, "99",
		_callercontext.CallerKindMetadataKey, string(_callercontext.KindUser),
	)
	if err := _transporttrust.Sign("/test.Service/Call", md, cfg); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	parent := mustBindActor(context.Background(), _identitycontext.Actor{ProfileID: 12})
	ctx := metadata.NewIncomingContext(parent, md)
	called := false
	_, err := parseGrpcMetadataContextMiddlewareWithTrust(
		cfg,
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"},
		func(context.Context, interface{}) (interface{}, error) {
			called = true
			return nil, nil
		},
	)
	if status.Code(err) != codes.Internal {
		t.Fatalf("status code = %v, want Internal", status.Code(err))
	}
	if called {
		t.Fatal("handler was called after Actor context conflict")
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
	ctx := mustBindActor(context.Background(), actor)
	ctx = context.WithValue(ctx, _enums.ProfileIDKey, uint64(99))

	md := GrpcClientMetadataFromContext(ctx)
	if got := md.Get(_identitycontext.ProfileIDMetadataKey); len(got) != 1 || got[0] != "12" {
		t.Fatalf("profile metadata = %v, want [12]", got)
	}
}

func TestPrincipalFromContextProjectsCanonicalActor(t *testing.T) {
	t.Parallel()

	organizationID := uint64(19)
	ctx := mustBindActor(context.Background(), _identitycontext.Actor{
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
	req = req.WithContext(mustBindCaller(req.Context(), _callercontext.Caller{
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

func TestStreamMetadataInterceptorAppliesOperationIDAcceptancePolicy(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		values   []string
		wantCode codes.Code
		called   bool
	}{
		{name: "missing allowed", values: nil, wantCode: codes.OK, called: true},
		{name: "valid accepted", values: []string{"op-stream-123"}, wantCode: codes.OK, called: true},
		{name: "invalid rejected", values: []string{"unsafe stream operation"}, wantCode: codes.InvalidArgument},
		{name: "multiple rejected", values: []string{"op-first", "op-second"}, wantCode: codes.InvalidArgument},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			md := metadata.MD{}
			for _, value := range testCase.values {
				md.Append(_requestcontext.OperationIDMetadataKey, value)
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
					operationID, ok := _requestcontext.OperationIDFromContext(handlerStream.Context())
					if len(testCase.values) == 1 && testCase.wantCode == codes.OK {
						if !ok || operationID != "op-stream-123" {
							t.Fatalf("operation identity = %q, %v", operationID, ok)
						}
					} else if ok || operationID != "" {
						t.Fatalf("operation identity = %q, %v; want absent", operationID, ok)
					}
					return nil
				},
			)
			if status.Code(err) != testCase.wantCode {
				t.Fatalf("status code = %v, want %v", status.Code(err), testCase.wantCode)
			}
			if called != testCase.called {
				t.Fatalf("handler called = %v, want %v", called, testCase.called)
			}
		})
	}
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

func TestParseGrpcMetadataDoesNotAcceptAmbiguousRequestID(
	t *testing.T,
) {
	t.Parallel()

	ctx :=
		metadata.NewIncomingContext(
			context.Background(),
			metadata.Pairs(
				_requestcontext.
					RequestIDMetadataKey,
				"req-first-123",

				_requestcontext.
					RequestIDMetadataKey,
				"req-second-456",
			),
		)

	_, err :=
		ParseGrpcMetadataContextMiddleware(
			ctx,
			nil,
			&grpc.UnaryServerInfo{
				FullMethod: "/test.Service/Call",
			},
			func(
				handlerCtx context.Context,
				_ interface{},
			) (
				interface{},
				error,
			) {
				if requestID, ok :=
					_requestcontext.
						FromContext(
							handlerCtx,
						); ok || requestID != "" {

					t.Fatalf(
						"ambiguous request ID reached handler: %q",
						requestID,
					)
				}

				return nil, nil
			},
		)

	if err != nil {
		t.Fatalf(
			"ParseGrpcMetadataContextMiddleware() error = %v",
			err,
		)
	}
}

func TestParseGrpcMetadataDoesNotInventMissingRequestID(
	t *testing.T,
) {
	t.Parallel()

	ctx :=
		metadata.NewIncomingContext(
			context.Background(),
			metadata.MD{},
		)

	_, err :=
		ParseGrpcMetadataContextMiddleware(
			ctx,
			nil,
			&grpc.UnaryServerInfo{
				FullMethod: "/test.Service/Call",
			},
			func(
				handlerCtx context.Context,
				_ interface{},
			) (
				interface{},
				error,
			) {
				if requestID, ok :=
					_requestcontext.
						FromContext(
							handlerCtx,
						); ok || requestID != "" {

					t.Fatalf(
						"missing request ID was invented: %q",
						requestID,
					)
				}

				return nil, nil
			},
		)

	if err != nil {
		t.Fatalf(
			"ParseGrpcMetadataContextMiddleware() error = %v",
			err,
		)
	}
}

func TestParseGrpcMetadataDoesNotBindInvalidRequestID(
	t *testing.T,
) {
	t.Parallel()

	ctx :=
		metadata.NewIncomingContext(
			context.Background(),
			metadata.Pairs(
				_requestcontext.
					RequestIDMetadataKey,
				"unsafe request id",
			),
		)

	_, err :=
		ParseGrpcMetadataContextMiddleware(
			ctx,
			nil,
			&grpc.UnaryServerInfo{
				FullMethod: "/test.Service/Call",
			},
			func(
				handlerCtx context.Context,
				_ interface{},
			) (
				interface{},
				error,
			) {
				if requestID, ok :=
					_requestcontext.
						FromContext(
							handlerCtx,
						); ok || requestID != "" {

					t.Fatalf(
						"invalid request ID reached handler: %q",
						requestID,
					)
				}

				return nil, nil
			},
		)

	if err != nil {
		t.Fatalf(
			"ParseGrpcMetadataContextMiddleware() error = %v",
			err,
		)
	}
}
