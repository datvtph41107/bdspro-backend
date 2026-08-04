package _middleware

import (
	_callercontext "common/callercontext"
	_quotasubject "common/quotasubject"
	_transporttrust "common/transporttrust"

	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestParseGrpcMetadataStoresVerifiedQuotaSubject(t *testing.T) {
	t.Parallel()

	now := time.Unix(1_700_000_000, 0)
	cfg := _transporttrust.Config{
		ServiceID: "gateway-service",
		Secret:    "test-secret",
		Mode:      _transporttrust.ModeEnforce,
		Now:       func() time.Time { return now },
	}
	md := metadata.Pairs(
		_callercontext.CallerKindMetadataKey, string(_callercontext.KindAnonymous),
	)
	want := _quotasubject.Subject{
		Kind: _quotasubject.KindAnonymous, ID: "public_bucket_01", Source: _quotasubject.SourceAnonymousEdge,
	}
	_quotasubject.AppendMetadata(md, want)
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
			got, ok := _quotasubject.FromContext(handlerCtx)
			if !ok || got != want {
				t.Fatalf("subject = %+v, %v, want %+v", got, ok, want)
			}
			return nil, nil
		},
	)
	if err != nil {
		t.Fatalf("middleware error = %v", err)
	}
}

func TestParseGrpcMetadataAuditRejectsUnsignedQuotaSubject(t *testing.T) {
	t.Parallel()

	md := metadata.Pairs(
		_callercontext.CallerKindMetadataKey, string(_callercontext.KindAnonymous),
	)
	_quotasubject.AppendMetadata(md, _quotasubject.Subject{
		Kind: _quotasubject.KindAnonymous, ID: "public_bucket_01", Source: _quotasubject.SourceAnonymousEdge,
	})

	called := false
	_, err := parseGrpcMetadataContextMiddlewareWithTrust(
		_transporttrust.Config{Mode: _transporttrust.ModeAudit},
		metadata.NewIncomingContext(context.Background(), md),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.PublicService/Get"},
		func(context.Context, interface{}) (interface{}, error) {
			called = true
			return nil, nil
		},
	)
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("status code = %v, want Unauthenticated", status.Code(err))
	}
	if called {
		t.Fatal("handler was called with an unsigned quota subject")
	}
}

func TestGrpcClientMetadataPropagatesQuotaSubject(t *testing.T) {
	t.Parallel()

	want := _quotasubject.Subject{
		Kind: _quotasubject.KindOrganization, ID: "12", Source: _quotasubject.SourceOrganizationMembership,
	}
	ctx := _quotasubject.WithSubject(context.Background(), want)
	md := GrpcClientMetadataFromContext(ctx)

	got, ok, err := _quotasubject.FromMetadata(md)
	if err != nil || !ok || got != want {
		t.Fatalf("metadata subject = %+v, %v, %v, want %+v", got, ok, err, want)
	}
}
