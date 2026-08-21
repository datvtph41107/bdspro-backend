package cmd

import (
	_middleware "common/middleware"
	"common/rpc"
	"common/transporttrust"
	"context"

	"google.golang.org/grpc"
)

// targetIngress keeps legacy parsing on RPCs that have not migrated yet and
// uses the canonical Foundation RPC transport on migrated RPCs. The selector
// decides migration scope; transport semantics themselves are owned by
// common/rpc.
func targetIngress(trust transporttrust.Config) grpc.UnaryServerInterceptor {
	canonical := rpc.UnaryServerInterceptor(rpc.TransportConfig{Trust: trust})

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info == nil || !usesReportTargetRuntime(info.FullMethod) {
			return _middleware.ParseGrpcMetadataContextMiddleware(
				ctx,
				req,
				info,
				handler,
			)
		}

		return canonical(ctx, req, info, handler)
	}
}
