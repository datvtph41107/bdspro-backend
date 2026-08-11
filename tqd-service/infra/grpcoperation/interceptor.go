package grpcoperation

import (
	"common/businessoperation"
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/**
 * UnaryServerInterceptor chọn business operation từ gRPC method.
 *
 * Method chưa được khai báo vẫn được chạy bình thường.
 * Việc một method bắt buộc phải có operation sẽ được kiểm tra ở use case cần nó,
 * thay vì biến mọi RPC trong service thành quota operation.
 */
func UnaryServerInterceptor(registry *Registry) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info == nil {
			return nil, status.Error(codes.Internal, "gRPC method info is missing")
		}

		operation, ok := registry.Find(info.FullMethod)
		if !ok {
			return handler(ctx, req)
		}

		bound, err := businessoperation.Bind(ctx, operation)
		if err != nil {
			if errors.Is(err, businessoperation.ErrContextConflict) {
				return nil, status.Error(codes.Internal, "business operation context conflict")
			}
			return nil, status.Error(codes.Internal, "cannot bind business operation")
		}

		return handler(bound, req)
	}
}
