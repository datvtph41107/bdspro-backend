package logger

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
)

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...any)
	Info(ctx context.Context, msg string, fields ...any)
	Warn(ctx context.Context, msg string, fields ...any)
	Error(ctx context.Context, msg string, fields ...any)
	With(fields ...any) Logger
	Sync() error
	HTTPMiddleware() func(http.Handler) http.Handler
	GRPCInterceptor() grpc.UnaryServerInterceptor
}
