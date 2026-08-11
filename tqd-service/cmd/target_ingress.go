package cmd

import (
	"common/identitygrpc"
	_middleware "common/middleware"
	"common/requestgrpc"
	"common/transporttrust"
	"context"

	"google.golang.org/grpc"
)

/**
 * targetIngress dùng metadata flow mới cho RPC đã migrate.
 * RPC còn lại vẫn đi parser cũ cho tới khi được chuyển từng use case.
 */
func targetIngress(trust transporttrust.Config, targetEnabled bool) grpc.UnaryServerInterceptor {
	requestIdentity := requestgrpc.UnaryServerInterceptor
	trustedIdentity := identitygrpc.UnaryServerInterceptor(trust)

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info == nil || !usesTargetQuota(targetEnabled, info.FullMethod) {
			return _middleware.ParseGrpcMetadataContextMiddleware(ctx, req, info, handler)
		}

		return requestIdentity(ctx, req, info, func(requestCtx context.Context, request interface{}) (interface{}, error) {
			return trustedIdentity(requestCtx, request, info, handler)
		})
	}
}

/**
 * skipTargetOperationRuntime giữ operation runtime cũ cho RPC chưa migrate.
 */
func skipTargetOperationRuntime(
	legacy grpc.UnaryServerInterceptor,
	targetEnabled bool,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info != nil && usesTargetQuota(targetEnabled, info.FullMethod) {
			return handler(ctx, req)
		}
		return legacy(ctx, req, info, handler)
	}
}
