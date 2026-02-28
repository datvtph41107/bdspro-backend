package ctxutil

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

type ctxKey string

const (
	keyTraceID ctxKey = "trace_id"
	keySpanID  ctxKey = "span_id"
)

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, keyTraceID, traceID)
}

func TraceID(ctx context.Context) string {
	val := ctx.Value(keyTraceID)
	if val == nil {
		return ""
	}
	return val.(string)
}

func WithSpanID(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, keySpanID, spanID)
}

func SpanID(ctx context.Context) string {
	val := ctx.Value(keySpanID)
	if val == nil {
		return ""
	}
	return val.(string)
}

func GenerateTraceID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
