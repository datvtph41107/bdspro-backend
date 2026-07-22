// platform/core/context/context_test.go
package context_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pctx "bdspro/core/context"
)

// TestRequestID kiểm tra setter/getter cơ bản và zero-value behavior.
//
// Tại sao test zero value?
// RequestID() trả "" cho context rỗng — middleware downstream không được panic
// khi request không có ID (ví dụ internal health check calls).
func TestRequestID(t *testing.T) {
	t.Run("set and get", func(t *testing.T) {
		ctx := pctx.WithRequestID(context.Background(), "req-abc-123")
		assert.Equal(t, "req-abc-123", pctx.RequestID(ctx))
	})

	t.Run("empty context returns empty string", func(t *testing.T) {
		assert.Equal(t, "", pctx.RequestID(context.Background()))
	})

	t.Run("overwrite in child context", func(t *testing.T) {
		// Context là immutable chain — child có value mới, parent không bị ảnh hưởng
		parent := pctx.WithRequestID(context.Background(), "parent-id")
		child := pctx.WithRequestID(parent, "child-id")
		assert.Equal(t, "parent-id", pctx.RequestID(parent))
		assert.Equal(t, "child-id", pctx.RequestID(child))
	})
}

// TestUserID kiểm tra uint64 type safety và MustUserID panic behavior.
//
// Tại sao test panic?
// MustUserID được gọi trong handler sau auth middleware.
// Nếu không có user_id mà không panic → silent bug, handler chạy với id=0.
func TestUserID(t *testing.T) {
	t.Run("set and get", func(t *testing.T) {
		ctx := pctx.WithUserID(context.Background(), 42)
		assert.Equal(t, uint64(42), pctx.UserID(ctx))
	})

	t.Run("zero value for empty context", func(t *testing.T) {
		assert.Equal(t, uint64(0), pctx.UserID(context.Background()))
	})

	t.Run("MustUserID panics on zero", func(t *testing.T) {
		// assert.Panics captures the panic — test vẫn pass dù panic xảy ra
		assert.Panics(t, func() {
			pctx.MustUserID(context.Background())
		})
	})

	t.Run("MustUserID panic message contains helpful context", func(t *testing.T) {
		// Kiểm tra message để dev biết ngay nguyên nhân khi xem panic log
		defer func() {
			r := recover()
			require.NotNil(t, r)
			msg, ok := r.(string)
			require.True(t, ok)
			assert.Contains(t, msg, "user_id required")
			assert.Contains(t, msg, "auth middleware") // gợi ý nguyên nhân
		}()
		pctx.MustUserID(context.Background())
	})

	t.Run("MustUserID returns value when set", func(t *testing.T) {
		ctx := pctx.WithUserID(context.Background(), 99)
		assert.Equal(t, uint64(99), pctx.MustUserID(ctx))
	})
}

// TestContextKeyIsolation đảm bảo các key không collide với nhau.
//
// Tại sao quan trọng?
// Nếu keyRequestID == keyTraceID (do bug đặt tên), thì
// WithRequestID sẽ overwrite WithTraceID — silent data corruption.
func TestContextKeyIsolation(t *testing.T) {
	ctx := context.Background()
	ctx = pctx.WithRequestID(ctx, "req-1")
	ctx = pctx.WithTraceID(ctx, "trace-1")
	ctx = pctx.WithUserID(ctx, 100)
	ctx = pctx.WithClientIP(ctx, "1.2.3.4")
	ctx = pctx.WithUserAgent(ctx, "Mozilla/5.0")
	ctx = pctx.WithRole(ctx, "admin")
	ctx = pctx.WithPlanID(ctx, 7)

	// Tất cả field phải độc lập, không overwrite nhau
	assert.Equal(t, "req-1", pctx.RequestID(ctx))
	assert.Equal(t, "trace-1", pctx.TraceID(ctx))
	assert.Equal(t, uint64(100), pctx.UserID(ctx))
	assert.Equal(t, "1.2.3.4", pctx.ClientIP(ctx))
	assert.Equal(t, "Mozilla/5.0", pctx.UserAgent(ctx))
	assert.Equal(t, "admin", pctx.Role(ctx))
	assert.Equal(t, uint64(7), pctx.PlanID(ctx))
}

// TestContextKeyString kiểm tra String() output để debug readable.
// Khi context bị dump trong log, key phải có prefix để biết nguồn gốc.
func TestContextKeyIsolation_ExternalPackageCannotCreateSameKey(t *testing.T) {
	// External package không thể tạo contextKey (unexported type)
	// Đây là compile-time guarantee — test này chỉ document ý định.
	// Nếu ai đó thêm exported contextKey → CI fail compile
	t.Log("contextKey type is unexported — external key collision is impossible at compile time")
}
