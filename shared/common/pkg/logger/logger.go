package logger

// // shared/pkg/logger/logger.go

// import (
// 	"context"
// 	"log/slog"
// 	"os"

// 	sharedconfig "common/pkg/config"
// )

// // Lý do dùng slog (Go 1.21 standard library) thay vì zap/zerolog:
// // - Không có thêm dependency
// // - Interface giống nhau, có thể swap backend
// // - slog.Handler là interface, dễ wrap thêm context fields
// // Khi cần performance cao hơn (>100k req/s), swap sang zap backend mà không sửa call site

// type contextKey struct{}

// func New(cfg sharedconfig.ObservabilityConfig) *slog.Logger {
// 	var handler slog.Handler
// 	opts := &slog.HandlerOptions{
// 		Level: parseLevel(cfg.LogLevel),
// 		// AddSource: chỉ dev — chi phí cao, thêm file:line vào mỗi log entry
// 		AddSource: cfg.LogLevel == "debug",
// 	}

// 	if cfg.LogFormat == "json" {
// 		handler = slog.NewJSONHandler(os.Stdout, opts)
// 	} else {
// 		handler = slog.NewTextHandler(os.Stdout, opts)
// 	}

// 	return slog.New(handler)
// }

// // WithContext extract requestID, traceID từ context và thêm vào logger
// // Dùng pattern này thay vì global logger để mỗi request có logger riêng với ID
// func WithContext(ctx context.Context, base *slog.Logger) *slog.Logger {
// 	attrs := []any{}

// 	if rid := RequestIDFromContext(ctx); rid != "" {
// 		attrs = append(attrs, slog.String("request_id", rid))
// 	}
// 	if tid := TraceIDFromContext(ctx); tid != "" {
// 		attrs = append(attrs, slog.String("trace_id", tid))
// 	}
// 	if uid := ProfileIDFromContext(ctx); uid != 0 {
// 		attrs = append(attrs, slog.Uint64("profile_id", uid))
// 	}

// 	return base.With(attrs...)
// }
