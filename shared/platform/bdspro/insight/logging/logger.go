// platform/observability/logging/logger.go
package logging

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	pctx "bdspro/core/context"
)

// Logger là thin wrapper quanh *zap.Logger.
//
// Tại sao không dùng interface Logger?
// - Zap đã là defacto standard trong Go ecosystem
// - Interface thêm indirection không cần thiết cho common case
// - Nếu cần swap logger: thay struct này, không thay interface
//
// Tại sao không embed *zap.Logger trực tiếp?
// - Embedding expose toàn bộ zap API — caller có thể gọi zap.Logger.Sugar() bypass cache
// - Wrapper cho phép kiểm soát API surface
type Logger struct {
	base *zap.Logger // base logger với service field đã gắn
}

// internalCtxKey cho logger cache — unexported, chỉ package logging dùng
// Tách khỏi pctx để không pollute public context API
type internalCtxKey struct{}

var loggerCtxKey = internalCtxKey{}

// New tạo Logger từ Config.
// Trả về error nếu config invalid — caller phải handle, không silent fail.
func New(cfg Config) (*Logger, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	zapCfg := buildZapConfig(cfg)

	base, err := zapCfg.Build(
		zap.AddCaller(),      // thêm file:line vào mỗi log entry
		zap.AddCallerSkip(1), // skip wrapper frame — show caller của Logger.Info, không phải logger.go
	)
	if err != nil {
		return nil, fmt.Errorf("logging: build zap logger: %w", err)
	}

	// Gắn service name vào mọi log entry — không cần caller tự thêm
	base = base.With(zap.String("service", cfg.ServiceName))

	return &Logger{base: base}, nil
}

// buildZapConfig tạo zap.Config phù hợp với environment.
func buildZapConfig(cfg Config) zap.Config {
	var zapCfg zap.Config

	switch cfg.Environment {
	case "production":
		zapCfg = zap.NewProductionConfig()
		// Production: JSON, no color, timestamp ISO8601
		zapCfg.EncoderConfig.TimeKey = "timestamp"
		zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	case "test":
		// Test: discard all output — test không bị noise từ log
		zapCfg = zap.NewProductionConfig()
		zapCfg.OutputPaths = []string{"/dev/null"}
		zapCfg.ErrorOutputPaths = []string{"/dev/null"}

	default: // development
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapCfg.EncoderConfig.TimeKey = "timestamp"
		zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Level override nếu có
	if cfg.Level != "" {
		var lvl zapcore.Level
		if err := lvl.UnmarshalText([]byte(cfg.Level)); err == nil {
			zapCfg.Level = zap.NewAtomicLevelAt(lvl)
		}
	}

	return zapCfg
}

// WithContext trả về *zap.Logger đã có request-scoped fields (request_id, trace_id, user_id...).
//
// Performance (F4):
// - Lần đầu: tạo *zap.Logger mới với các context fields, cache vào context
// - Lần sau (cùng context): đọc từ cache — ZERO allocation
//
// Tại sao return *zap.Logger thay vì Logger?
// - Caller cần gọi zap methods trực tiếp (Info, Error, With, ...)
// - Wrap thêm một lần nữa = thêm boilerplate không cần thiết
// - zap đã được dùng trực tiếp ở middleware — nhất quán
func (l *Logger) WithContext(ctx context.Context) *zap.Logger {
	// Cache lookup (F4) — tránh allocation mỗi log call
	if cached, ok := ctx.Value(loggerCtxKey).(*zap.Logger); ok {
		return cached
	}

	// Build fields từ context values
	// Chỉ add field nếu có value — tránh empty string fields trong log
	fields := make([]zap.Field, 0, 5)

	if v := pctx.RequestID(ctx); v != "" {
		fields = append(fields, zap.String("request_id", v))
	}
	if v := pctx.TraceID(ctx); v != "" {
		fields = append(fields, zap.String("trace_id", v))
	}
	if v := pctx.UserID(ctx); v != 0 {
		fields = append(fields, zap.Uint64("user_id", v))
	}
	if v := pctx.ClientIP(ctx); v != "" {
		fields = append(fields, zap.String("client_ip", v))
	}
	if v := pctx.UserAgent(ctx); v != "" {
		fields = append(fields, zap.String("user_agent", v))
	}

	zl := l.base.With(fields...)

	// Cache KHÔNG thể mutate ctx (context immutable) nên không store back.
	// Cache chỉ work nếu caller store ctx sau khi middleware set.
	// Middleware sẽ store logger vào context sau lần đầu tạo.
	return zl
}

// InjectCachedLogger lưu logger đã build vào context.
// Được gọi bởi Logging middleware sau khi tạo logger lần đầu.
// Các handler sau cùng context sẽ đọc từ cache.
func InjectCachedLogger(ctx context.Context, zl *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, zl)
}

// Base trả về base logger (không có context fields).
// Dùng cho startup/shutdown log — không có request context.
func (l *Logger) Base() *zap.Logger {
	return l.base
}

// Sync flush buffered log entries.
// Phải gọi trước khi process exit (defer pipe.Sync()).
func (l *Logger) Sync() error {
	return l.base.Sync()
}

func NewWithCore(core zapcore.Core, service string) *Logger {
	base := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).
		With(zap.String("service", service))

	return &Logger{base: base}
}
