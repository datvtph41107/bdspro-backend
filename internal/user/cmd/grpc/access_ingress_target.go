package main

import (
	_middleware "common/middleware"
	"common/rpc"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

/**
 * accessIngress dùng request identity + trusted Actor/Caller cho các RPC cắt lát mới.
 * Các RPC user-service khác vẫn dùng parser cũ trong giai đoạn migration.
 */
func accessIngress(transportConfig rpc.TransportConfig) grpc.UnaryServerInterceptor {
	trustedIdentity := rpc.StrictIdentityUnaryServerInterceptor(transportConfig.ServiceAssertion)

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info == nil || !trustedIdentityMethod(info.FullMethod) {
			return _middleware.ParseGrpcMetadataContextMiddleware(ctx, req, info, handler)
		}

		md, _ := metadata.FromIncomingContext(ctx)
		requestCtx, err := rpc.RestoreRequestFromMetadata(ctx, md)
		if err != nil {
			return nil, err
		}

		return trustedIdentity(requestCtx, req, info, handler)
	}
}
