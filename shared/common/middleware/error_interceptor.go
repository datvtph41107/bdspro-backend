package _middleware

import (
	"context"
	"log/slog"

	_errors "common/errors"
	"common/logging"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryErrorInterceptor is the canonical outer application-error normalization
// boundary for unary RPC handlers. Expected application failures are projected
// without being promoted to operational ERROR. Unknown Go errors are logged
// once with request context and sanitized before they cross the RPC boundary.
func UnaryErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}
		return resp, normalizeRPCError(ctx, unaryMethod(info), err)
	}
}

func StreamErrorInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := handler(srv, stream)
		if err == nil {
			return nil
		}
		ctx := context.Background()
		if stream != nil && stream.Context() != nil {
			ctx = stream.Context()
		}
		method := ""
		if info != nil {
			method = info.FullMethod
		}
		return normalizeRPCError(ctx, method, err)
	}
}

func normalizeRPCError(ctx context.Context, method string, err error) error {
	if _, ok := _errors.As(err); ok {
		return _errors.ToGRPC(err)
	}

	// Existing explicit gRPC statuses are preserved during migration. Their
	// remaining business/usecase emitters are separately inventoried and retired.
	if _, ok := status.FromError(err); ok {
		return err
	}

	logging.WithComponent(ctx, "grpc").ErrorContext(
		ctx,
		"gRPC handler failed",
		slog.String("event_name", "grpc.handler.error"),
		slog.String("grpc.method", method),
		slog.Any("error", err),
	)
	return _errors.ToGRPC(err)
}

func unaryMethod(info *grpc.UnaryServerInfo) string {
	if info == nil {
		return ""
	}
	return info.FullMethod
}
