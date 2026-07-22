// platform/observability/logging/logger_test.go
package logging_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	pctx "bdspro/core/context"
	"bdspro/insight/logging"
)

// newObservableLogger tạo Logger với observer core để test.
// observer.New cho phép inspect log entries sau khi log.
// Đây là pattern chuẩn của zap testing — không mock, inspect thật.
func newObservableLogger(t *testing.T) (*logging.Logger, *observer.ObservedLogs) {
	t.Helper()

	core, logs := observer.New(zapcore.InfoLevel)
	logger := logging.NewWithCore(core, "test-service")

	return logger, logs
}

// TestLogger_WithContext_InjectsContextFields kiểm tra F4 và field injection.
func TestLogger_WithContext_InjectsContextFields(t *testing.T) {
	logger, logs := newObservableLogger(t)

	ctx := context.Background()
	ctx = pctx.WithRequestID(ctx, "req-abc")
	ctx = pctx.WithTraceID(ctx, "trace-xyz")
	ctx = pctx.WithUserID(ctx, 42)
	ctx = pctx.WithClientIP(ctx, "1.2.3.4")
	ctx = pctx.WithUserAgent(ctx, "TestAgent/1.0")

	zl := logger.WithContext(ctx)
	zl.Info("test message")

	require.Equal(t, 1, logs.Len())
	entry := logs.All()[0]

	// Kiểm tra từng field được inject từ context
	contextMap := entry.ContextMap()
	assert.Equal(t, "req-abc", contextMap["request_id"], "request_id must come from context")
	assert.Equal(t, "trace-xyz", contextMap["trace_id"], "trace_id must come from context")
	assert.Equal(t, int64(42), contextMap["user_id"], "user_id must come from context")
	assert.Equal(t, "1.2.3.4", contextMap["client_ip"], "client_ip must come from context")
	assert.Equal(t, "TestAgent/1.0", contextMap["user_agent"], "user_agent must come from context")
	assert.Equal(t, "test-service", contextMap["service"], "service name always present")
}

// TestLogger_WithContext_Cache kiểm tra F4 — cache tránh allocation.
func TestLogger_WithContext_Cache(t *testing.T) {
	logger, _ := newObservableLogger(t)

	ctx := context.Background()
	ctx = pctx.WithRequestID(ctx, "req-123")

	// Lần đầu: tạo logger mới
	zl1 := logger.WithContext(ctx)
	// Inject vào context
	ctx = logging.InjectCachedLogger(ctx, zl1)

	// Lần sau: phải return pointer giống hệt
	zl2 := logger.WithContext(ctx)

	// Pointer equality — cùng *zap.Logger object
	assert.Same(t, zl1, zl2, "WithContext should return cached logger on subsequent calls")
}

// TestLogger_WithContext_EmptyContext không panic khi context trống.
func TestLogger_WithContext_EmptyContext(t *testing.T) {
	logger, logs := newObservableLogger(t)

	// Context không có bất kỳ field nào — không panic
	zl := logger.WithContext(context.Background())
	require.NotNil(t, zl)
	zl.Info("empty context log")

	require.Equal(t, 1, logs.Len())
	entry := logs.All()[0]

	// Chỉ có "service" field — không có request_id/trace_id (không có trong context)
	fields := entry.ContextMap()
	assert.Equal(t, "test-service", fields["service"])
	_, hasReqID := fields["request_id"]
	assert.False(t, hasReqID, "empty string fields should not be included")
}

// TestLogger_Sync không trả error với discard core.
func TestLogger_Sync(t *testing.T) {
	logger, _ := newObservableLogger(t)
	err := logger.Sync()
	// observer core không cần sync — không error
	assert.NoError(t, err)
}

// BenchmarkLogger_WithContext_Cached đo allocation với cache.
//
// Expected: 0 alloc/op khi đã cache
func BenchmarkLogger_WithContext_Cached(b *testing.B) {
	core, _ := observer.New(zapcore.InfoLevel)
	logger := logging.NewWithCore(core, "bench")

	ctx := context.Background()
	ctx = pctx.WithRequestID(ctx, "req-bench")
	zl := logger.WithContext(ctx)
	ctx = logging.InjectCachedLogger(ctx, zl) // cache đã có

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = logger.WithContext(ctx)
	}
	// Expected: 0 allocs/op — chỉ read từ context map
}

// BenchmarkLogger_WithContext_NoCache đo allocation không có cache.
func BenchmarkLogger_WithContext_NoCache(b *testing.B) {
	core, _ := observer.New(zapcore.InfoLevel)
	logger := logging.NewWithCore(core, "bench")

	ctx := context.Background()
	ctx = pctx.WithRequestID(ctx, "req-bench")
	ctx = pctx.WithTraceID(ctx, "trace-bench")
	ctx = pctx.WithUserID(ctx, 42)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = logger.WithContext(ctx)
	}
	// Expected: ~3-5 allocs/op — tạo zap.Field slice và *zap.Logger
}
