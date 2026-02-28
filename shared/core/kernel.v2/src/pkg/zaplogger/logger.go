package zaplogger

import (
	"context"
	"os"

	"core/kernel.v2/src/contracts/logger"
	"core/kernel.v2/src/internal/ctxutil"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	zap *zap.Logger
}

func New(cfg logger.Config) (*ZapLogger, error) {
	var zapCfg zap.Config
	if cfg.Development {
		zapCfg = zap.NewDevelopmentConfig()
	} else {
		zapCfg = zap.NewProductionConfig()
	}

	lvl, err := zapcore.ParseLevel(string(cfg.Level))
	if err != nil {
		lvl = zapcore.InfoLevel
	}
	zapCfg.Level = zap.NewAtomicLevelAt(lvl)

	if cfg.Format == "console" {
		zapCfg.Encoding = "console"
	} else {
		zapCfg.Encoding = "json"
	}

	if cfg.OutputPath != "" && cfg.OutputPath != "stdout" {
		zapCfg.OutputPaths = []string{cfg.OutputPath}
	} else {
		zapCfg.OutputPaths = []string{"stdout"}
	}
	zapCfg.ErrorOutputPaths = []string{"stderr"}

	z, err := zapCfg.Build()
	if err != nil {
		return nil, err
	}
	return &ZapLogger{zap: z}, nil
}

func toZapFields(ctx context.Context, fields []logger.Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields)+1)
	if traceID := ctxutil.TraceID(ctx); traceID != "" {
		zapFields = append(zapFields, zap.String("trace_id", traceID))
	}
	if spanID := ctxutil.SpanID(ctx); spanID != "" {
		zapFields = append(zapFields, zap.String("span_id", spanID))
	}
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}

func (l *ZapLogger) Debug(ctx context.Context, msg string, fields ...logger.Field) {
	l.zap.Debug(msg, toZapFields(ctx, fields)...)
}
func (l *ZapLogger) Info(ctx context.Context, msg string, fields ...logger.Field) {
	l.zap.Info(msg, toZapFields(ctx, fields)...)
}
func (l *ZapLogger) Warn(ctx context.Context, msg string, fields ...logger.Field) {
	l.zap.Warn(msg, toZapFields(ctx, fields)...)
}
func (l *ZapLogger) Error(ctx context.Context, msg string, fields ...logger.Field) {
	l.zap.Error(msg, toZapFields(ctx, fields)...)
}
func (l *ZapLogger) Fatal(ctx context.Context, msg string, fields ...logger.Field) {
	l.zap.Fatal(msg, toZapFields(ctx, fields)...)
	os.Exit(1)
}

func (l *ZapLogger) With(fields ...logger.Field) logger.Logger {
	return &ZapLogger{
		zap: l.zap.With(toZapFields(context.Background(), fields)...),
	}
}

func (l *ZapLogger) Sync() error {
	return l.zap.Sync()
}
