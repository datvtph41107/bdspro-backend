package requestgrpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

/**
 * UnaryServerInterceptor restore request identity trước khi handler chạy.
 */
func UnaryServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	bound, err := RestoreFromMetadata(ctx, md)
	if err != nil {
		return nil, err
	}
	return handler(bound, req)
}
