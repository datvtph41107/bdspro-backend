package _middleware

import (
	_callercontext "common/callercontext"
	_operationcontext "common/operationcontext"
	_transporttrust "common/transporttrust"

	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestParseGrpcMetadataStoresVerifiedOperation(t *testing.T) {
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
	want, err := _operationcontext.New(
		"parcel.lookup",
		"parcel.lookup_transport",
		_operationcontext.SourceDirectBinding,
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	_operationcontext.AppendMetadata(md, want)
	if err := _transporttrust.Sign("/test.ParcelService/Get", md, cfg); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	ctx := metadata.NewIncomingContext(context.Background(), md)
	_, err = parseGrpcMetadataContextMiddlewareWithTrust(
		cfg,
		ctx,
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.ParcelService/Get"},
		func(handlerCtx context.Context, req interface{}) (interface{}, error) {
			got, ok := _operationcontext.FromContext(handlerCtx)
			if !ok || got != want {
				t.Fatalf("operation = %+v, %v, want %+v", got, ok, want)
			}
			return nil, nil
		},
	)
	if err != nil {
		t.Fatalf("middleware error = %v", err)
	}
}

func TestParseGrpcMetadataRejectsUnsignedOperation(t *testing.T) {
	t.Parallel()

	md := metadata.Pairs(
		_callercontext.CallerKindMetadataKey, string(_callercontext.KindAnonymous),
	)
	selection, _ := _operationcontext.New(
		"parcel.lookup",
		"parcel.lookup_transport",
		_operationcontext.SourceDirectBinding,
	)
	_operationcontext.AppendMetadata(md, selection)

	called := false
	_, err := parseGrpcMetadataContextMiddlewareWithTrust(
		_transporttrust.Config{Mode: _transporttrust.ModeAudit},
		metadata.NewIncomingContext(context.Background(), md),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.ParcelService/Get"},
		func(context.Context, interface{}) (interface{}, error) {
			called = true
			return nil, nil
		},
	)
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("status code = %v, want Unauthenticated", status.Code(err))
	}
	if called {
		t.Fatal("handler was called with unsigned operation metadata")
	}
}

func TestGrpcClientMetadataPropagatesOperation(t *testing.T) {
	t.Parallel()

	want, _ := _operationcontext.New(
		"parcel.lookup",
		"parcel.lookup_transport",
		_operationcontext.SourceDirectBinding,
	)
	ctx, err := _operationcontext.Bind(context.Background(), want)
	if err != nil {
		t.Fatalf("Bind() error = %v", err)
	}
	md := GrpcClientMetadataFromContext(ctx)

	got, ok, err := _operationcontext.FromMetadata(md)
	if err != nil || !ok || got != want {
		t.Fatalf("metadata operation = %+v, %v, %v, want %+v", got, ok, err, want)
	}
}
