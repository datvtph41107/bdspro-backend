package requestgrpc

import (
	"common/rpc"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryServerInterceptor preserves the historical request-only facade. New
// server composition should use rpc.ServerTransport so trust, request, and
// identity ordering remains one canonical boundary.
func UnaryServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	bound, err := rpc.RestoreRequestFromMetadata(ctx, md)
	if err != nil {
		return nil, err
	}
	return handler(bound, req)
}
