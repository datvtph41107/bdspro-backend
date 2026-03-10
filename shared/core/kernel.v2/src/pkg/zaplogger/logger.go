package zaplogger

import (
	"context"
	"net/http"
	"time"

	"core/kernel.v2/src/contracts/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

type ZapLogger struct {
	logger *zap.Logger
}

func New(level, format string) (*ZapLogger, error) {
	var cfg zap.Config
	if format == "console" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	lvl, err := zapcore.ParseLevel(level)
	if err != nil {
		lvl = zapcore.InfoLevel
	}
	cfg.Level = zap.NewAtomicLevelAt(lvl)
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &ZapLogger{logger: logger}, nil
}

func (l *ZapLogger) Debug(ctx context.Context, msg string, fields ...any) {
	l.logger.Debug(msg, l.toFields(fields...)...)
}
func (l *ZapLogger) Info(ctx context.Context, msg string, fields ...any) {
	l.logger.Info(msg, l.toFields(fields...)...)
}
func (l *ZapLogger) Warn(ctx context.Context, msg string, fields ...any) {
	l.logger.Warn(msg, l.toFields(fields...)...)
}
func (l *ZapLogger) Error(ctx context.Context, msg string, fields ...any) {
	l.logger.Error(msg, l.toFields(fields...)...)
}
func (l *ZapLogger) With(fields ...any) logger.Logger {
	newLogger := l.logger.With(l.toFields(fields...)...)
	return &ZapLogger{logger: newLogger}
}
func (l *ZapLogger) Sync() error {
	return l.logger.Sync()
}
func (l *ZapLogger) toFields(fields ...any) []zap.Field {
	// chuyển đổi []any thành []zap.Field
	// giả sử fields là cặp key-value
	if len(fields)%2 != 0 {
		return []zap.Field{zap.Any("fields", fields)}
	}
	zfields := make([]zap.Field, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		zfields = append(zfields, zap.Any(key, fields[i+1]))
	}
	return zfields
}

func (l *ZapLogger) HTTPMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			l.Info(r.Context(), "http request",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start),
			)
		})
	}
}

func (l *ZapLogger) GRPCInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		l.Info(ctx, "grpc call",
			"method", info.FullMethod,
			"duration", time.Since(start),
			"error", err,
		)
		return resp, err
	}
}
