package interceptors

import (
	"common/logging"
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

func UnaryLoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)

		logger := logging.WithComponent(ctx, "grpc")
		attrs := []any{
			slog.String("method", info.FullMethod),
			slog.Duration("duration", time.Since(start)),
		}
		if err != nil {
			attrs = append(attrs, slog.Any("error", err))
			logger.Error("unary gRPC completed", attrs...)
		} else {
			logger.Info("unary gRPC completed", attrs...)
		}

		return resp, err
	}
}
