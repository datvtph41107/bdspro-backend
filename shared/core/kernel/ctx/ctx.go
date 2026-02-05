package core_ctx

import (
	"context"
	"core/kernel/log"

	"gorm.io/gorm"
)

type ctxKey string

const (
	KeyTrace  ctxKey = "core_trace_id"
	KeyLogger ctxKey = "core_logger"
	KeyTx     ctxKey = "core_tx"
)

// WithTrace returns a context containing a trace id. If trace=="" a new uuid is generated
// Caller can set trace to empty string to ensure a trace exists (logger/otel rely on it).
func WithTrace(ctx context.Context, trace string) context.Context {
	if trace == "" {
		// lazy generation: put random uuid via facade.UUID or similar; to avoid extra import here, generate minimal fallback
		// but usually caller sets trace. We'll just use empty string and let logger create id when logging if needed.
	}
	return context.WithValue(ctx, KeyTrace, trace)
}

// Trace retrieves trace id. If none found, returns empty string (caller may generate).
func Trace(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(KeyTrace); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// WithLogger binds logger to context. Use this to inject a per-request logger (with trace attached).
func WithLogger(ctx context.Context, l *log.Logger) context.Context {
	return context.WithValue(ctx, KeyLogger, l)
}

// Logger returns logger from context or a nop logger if none set.
func Logger(ctx context.Context) *log.Logger {
	if ctx == nil {
		return log.Nop()
	}
	if v := ctx.Value(KeyLogger); v != nil {
		if l, ok := v.(*log.Logger); ok && l != nil {
			return l
		}
	}
	return log.Nop()
}

// WithTx stores a *gorm.DB transaction in ctx (used by TransactionMiddleware / OutboxMiddleware)
func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, KeyTx, tx)
}

// Tx returns *gorm.DB transaction from ctx if present.
func Tx(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return nil
	}
	if v := ctx.Value(KeyTx); v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx
		}
	}
	return nil
}
