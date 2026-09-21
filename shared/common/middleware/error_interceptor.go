package _middleware

import (
	"context"
	"errors"
	"log/slog"

	_errors "common/errors"
	"common/logging"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryErrorInterceptor is the canonical outer application-error normalization
// boundary for unary RPC handlers. Expected application failures are projected
// without being promoted to operational ERROR. Canonical operational failures
// and unknown Go errors are logged exactly once with request context before the
// public-safe gRPC projection is returned.
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
	if application, ok := _errors.As(err); ok {
		logCanonicalOperationalFailure(ctx, method, application)
		return _errors.ToGRPC(err)
	}

	// Existing explicit gRPC statuses are preserved during migration. Their
	// remaining business/usecase emitters are separately inventoried and retired.
	// Do not reinterpret codes.Internal here because legacy business callers still
	// use it as a compatibility transport code.
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

func logCanonicalOperationalFailure(ctx context.Context, method string, application *_errors.Error) {
	if application == nil || !isOperationalRPCCode(application.RPCCode()) {
		return
	}

	attrs := []slog.Attr{
		slog.String("event_name", "grpc.application.error"),
		slog.String("grpc.method", method),
		slog.String("grpc_code", application.RPCCode().String()),
		slog.Int64("error_code", int64(application.Code())),
		slog.String("error_reason", string(application.Key())),
	}
	if cause := errors.Unwrap(application); cause != nil {
		attrs = append(attrs, slog.Any("error", cause))
	}

	logging.WithComponent(ctx, "grpc").LogAttrs(
		ctx,
		slog.LevelError,
		"canonical application failure",
		attrs...,
	)
}

func isOperationalRPCCode(code codes.Code) bool {
	switch code {
	case codes.Internal, codes.Unknown, codes.Unavailable, codes.DataLoss, codes.DeadlineExceeded:
		return true
	default:
		return false
	}
}

func unaryMethod(info *grpc.UnaryServerInfo) string {
	if info == nil {
		return ""
	}
	return info.FullMethod
}
