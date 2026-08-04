package operationruntime

import (
	_operationbinding "common/operationbinding"

	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/** UnaryServerInterceptor binds approved operations before handlers. */
func (r *Runtime) UnaryServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if info == nil {
		return nil, status.Error(codes.Internal, "gRPC method information is required")
	}

	bound, _, err := r.Bind(ctx, req, _operationbinding.GRPCUnary(info.FullMethod))
	if err != nil {
		if current, ok := status.FromError(err); ok && current.Code() != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	return handler(bound, req)
}
