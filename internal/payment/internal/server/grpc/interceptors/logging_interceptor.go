package interceptors

import (
	"context"
	"log/slog"
	"time"

	_errors "common/errors"
	"common/logging"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryLoggerInterceptor owns only the Payment request-completion projection.
// Error normalization, operational severity and technical-cause evidence belong
// to common/middleware.UnaryErrorInterceptor.
func UnaryLoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()
		resp, err = handler(ctx, req)

		attrs := []any{
			slog.String("method", info.FullMethod),
			slog.Duration("duration", time.Since(start)),
		}
		if err == nil {
			attrs = append(attrs, slog.String("outcome", "success"))
		} else {
			attrs = append(attrs, slog.String("outcome", "error"))
			if application, ok := _errors.As(err); ok {
				attrs = append(attrs,
					slog.Int64("error_code", int64(application.Code())),
					slog.String("error_reason", string(application.Key())),
					slog.String("grpc_code", application.RPCCode().String()),
				)
			} else if grpcStatus, ok := status.FromError(err); ok {
				attrs = append(attrs, slog.String("grpc_code", grpcStatus.Code().String()))
			}
		}

		logging.WithComponent(ctx, "grpc").Info("unary gRPC completed", attrs...)
		return resp, err
	}
}
