package cmd_grpc

import (
	"common/identitygrpc"
	_middleware "common/middleware"
	"common/requestgrpc"
	"common/transporttrust"
	"context"

	"google.golang.org/grpc"
)

/**
 * accessIngress dùng request identity + trusted Actor/Caller cho các RPC cắt lát mới.
 * Các RPC user-service khác vẫn dùng parser cũ trong giai đoạn migration.
 */
func accessIngress(trust transporttrust.Config) grpc.UnaryServerInterceptor {
	requestIdentity := requestgrpc.UnaryServerInterceptor
	trustedIdentity := identitygrpc.UnaryServerInterceptor(trust)

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info == nil || !trustedIdentityMethod(info.FullMethod) {
			return _middleware.ParseGrpcMetadataContextMiddleware(ctx, req, info, handler)
		}
		return requestIdentity(ctx, req, info, func(requestCtx context.Context, request interface{}) (interface{}, error) {
			return trustedIdentity(requestCtx, request, info, handler)
		})
	}
}
