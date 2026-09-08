package operationgrpc

import (
	"context"
	"errors"

	"common/operation"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor lấy operation code từ protobuf contract của RPC
// và gắn operation đó vào context trước khi handler được gọi.
//
// RPC không có operation annotation thì giữ nguyên context và gọi handler.
// Nếu có annotation nhưng contract không hợp lệ hoặc không thể bind vào context,
// interceptor trả về lỗi và không gọi handler.
func UnaryServerInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	if info == nil {
		return nil, status.Error(
			codes.Internal,
			"gRPC method info is missing",
		)
	}

	code, found, err := codeForMethod(info.FullMethod)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			"operation contract lookup failed",
		)
	}

	if !found {
		return handler(ctx, req)
	}

	bound, err := operation.Bind(ctx, code)
	if err != nil {
		if errors.Is(err, operation.ErrContextConflict) {
			return nil, status.Error(
				codes.Internal,
				"operation context conflict",
			)
		}

		return nil, status.Error(
			codes.Internal,
			"operation context binding failed",
		)
	}

	return handler(bound, req)
}
