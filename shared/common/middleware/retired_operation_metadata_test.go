package _middleware

import (
	_transporttrust "common/transporttrust"

	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestParseGrpcMetadataIgnoresRetiredOperationHeaders(
	t *testing.T,
) {
	t.Parallel()

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(
			"x-qhpro-operation-code",
			"parcel.lookup",
			"x-qhpro-operation-binding",
			"parcel.lookup_transport",
			"x-qhpro-operation-source",
			"direct_binding",
		),
	)

	called := false

	_, err := parseGrpcMetadataContextMiddlewareWithTrust(
		_transporttrust.Config{
			Mode: _transporttrust.ModeAudit,
		},
		ctx,
		nil,
		&grpc.UnaryServerInfo{
			FullMethod: "/test.ParcelService/Get",
		},
		func(
			context.Context,
			interface{},
		) (interface{}, error) {
			called = true
			return nil, nil
		},
	)
	if err != nil {
		t.Fatalf(
			"retired operation headers caused middleware error: %v",
			err,
		)
	}
	if !called {
		t.Fatal(
			"handler was not called for retired operation headers",
		)
	}
}
