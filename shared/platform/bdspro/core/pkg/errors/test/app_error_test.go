// platform/errors/app_error_test.go
package apperr_test

import (
	apperr "bdspro/core/pkg/errors"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAppError_Is kiểm tra Is() semantics — quan trọng nhất vì errors.Is phụ thuộc vào đây.
//
// Kịch bản cần test:
// 1. Same kind+code → true
// 2. Different kind → false (F9: không match nhầm)
// 3. Different code → false
// 4. Wrapped error → true (Unwrap chain)
// 5. Non-AppError target → false
func TestAppError_Is(t *testing.T) {
	t.Run("same kind and code returns true", func(t *testing.T) {
		err := apperr.New(apperr.KindBusiness, apperr.CodeNotFound, "not found")
		assert.True(t, errors.Is(err, apperr.ErrNotFound))
	})

	t.Run("different kind returns false even if code same", func(t *testing.T) {
		// Không thể tạo Internal với CodeNotFound qua constructor chuẩn
		// Nhưng test đảm bảo kind isolation
		errA := apperr.New(apperr.KindInternal, apperr.CodeNotFound, "msg")
		errB := apperr.New(apperr.KindBusiness, apperr.CodeNotFound, "msg")
		assert.False(t, errors.Is(errA, errB), "different kind should not match")
	})

	t.Run("different code returns false", func(t *testing.T) {
		assert.False(t, errors.Is(apperr.ErrNotFound, apperr.ErrAlreadyExists))
	})

	t.Run("F9: Internal does not match Business even if both are AppError", func(t *testing.T) {
		// Đây chính là F9 — trước khi fix, một số implementation match bất kỳ AppError nào
		assert.False(t, errors.Is(apperr.ErrInternal, apperr.ErrNotFound))
		assert.False(t, errors.Is(apperr.ErrNotFound, apperr.ErrInternal))
	})

	t.Run("customized error still matches sentinel", func(t *testing.T) {
		// ErrNotFound được customize nhưng vẫn Is(ErrNotFound)
		err := apperr.ErrNotFound.
			WithField("id", "NOT_FOUND", "User không tồn tại", 123).
			WithOp("UserService.Get")
		assert.True(t, errors.Is(err, apperr.ErrNotFound))
	})

	t.Run("wrapped in fmt.Errorf chain", func(t *testing.T) {
		// errors.Is phải traverse Unwrap chain
		inner := apperr.Internal("DB", fmt.Errorf("timeout"))
		wrapped := fmt.Errorf("handler: %w", inner)
		assert.True(t, errors.Is(wrapped, apperr.ErrInternal))
	})

	t.Run("nil target returns false", func(t *testing.T) {
		assert.False(t, errors.Is(apperr.ErrNotFound, (*apperr.AppError)(nil)))
	})
}

// TestInternal_AlwaysHasStack kiểm tra F10 — Internal() phải capture stack.
func TestInternal_AlwaysHasStack(t *testing.T) {
	t.Run("Internal constructor captures stack", func(t *testing.T) {
		err := apperr.Internal("TestOp", fmt.Errorf("db error"))
		require.NotEmpty(t, err.Stack(), "Internal error MUST have stack trace (F10)")
	})

	t.Run("stack contains caller function name", func(t *testing.T) {
		err := apperr.Internal("TestOp.SubOp", fmt.Errorf("test cause"))
		stackStr := string(err.Stack())
		// Stack phải chứa goroutine info
		assert.Contains(t, stackStr, "goroutine")
	})

	t.Run("Internal with nil cause creates synthetic cause", func(t *testing.T) {
		err := apperr.Internal("op", nil)
		// Không panic, có cause
		require.NotNil(t, err.Cause())
		assert.NotEmpty(t, err.Stack())
	})

	t.Run("Validation does not capture stack", func(t *testing.T) {
		// Validation không phải bug — không cần stack
		err := apperr.Validation("bad input")
		assert.Empty(t, err.Stack(), "Validation error should NOT have stack")
	})

	t.Run("External does not capture stack", func(t *testing.T) {
		err := apperr.External(apperr.CodeTimeout, "timeout", fmt.Errorf("ctx deadline"))
		assert.Empty(t, err.Stack(), "External error should NOT have stack")
	})
}

// TestAppError_OpChain kiểm tra op chain building.
// Op chain: outermost caller đầu tiên, root cause cuối cùng.
func TestAppError_OpChain(t *testing.T) {
	err := apperr.Internal("DB.Query", fmt.Errorf("timeout")).
		WithOp("UserRepo.GetByID").
		WithOp("UserService.Get")

	assert.Equal(t, "UserService.Get:UserRepo.GetByID:DB.Query", err.Op())
}

// TestAppError_ImmutableClone đảm bảo With* không mutate original.
// Quan trọng vì sentinel errors là package-level vars — nếu mutate thì global state bị ô nhiễm.
func TestAppError_ImmutableClone(t *testing.T) {
	original := apperr.ErrNotFound

	// Customize — tạo bản mới
	customized := apperr.ErrNotFound.
		WithField("id", "NOT_FOUND", "not found", 42).
		WithOp("GetUser")

	// Original không bị thay đổi
	assert.False(t, original.HasFields(), "original should not have fields")
	assert.Equal(t, "", original.Op(), "original op should be empty")

	// Customized có fields mới
	assert.True(t, customized.HasFields())
	assert.Equal(t, "GetUser", customized.Op())
}

// TestFields_StripValue đảm bảo Fields() không trả về Value (security).
func TestFields_StripValue(t *testing.T) {
	err := apperr.ErrInvalidArgument.WithField(
		"email", "INVALID_FORMAT", "Email không hợp lệ", "user@badformat",
	)

	fields := err.Fields()
	require.Len(t, fields, 1)

	f := fields[0]
	assert.Equal(t, "email", f.Field)
	assert.Equal(t, "INVALID_FORMAT", f.Code)
	assert.Equal(t, "Email không hợp lệ", f.Message)
	assert.Nil(t, f.Value, "Fields() MUST strip Value — security requirement")
}

// TestHTTPStatus kiểm tra mapping — đây là nguồn sự thật duy nhất.
func TestHTTPStatus(t *testing.T) {
	cases := []struct {
		err      error
		wantHTTP int
	}{
		{apperr.ErrInvalidArgument, http.StatusBadRequest},        // 400
		{apperr.ErrInvalidJSON, http.StatusBadRequest},            // 400
		{apperr.ErrUnauthorized, http.StatusUnauthorized},         // 401
		{apperr.ErrForbidden, http.StatusForbidden},               // 403
		{apperr.ErrNotFound, http.StatusNotFound},                 // 404
		{apperr.ErrAlreadyExists, http.StatusConflict},            // 409
		{apperr.ErrRateLimit, http.StatusTooManyRequests},         // 429
		{apperr.ErrUnavailable, http.StatusServiceUnavailable},    // 503
		{apperr.ErrTimeout, http.StatusGatewayTimeout},            // 504
		{apperr.ErrInternal, http.StatusInternalServerError},      // 500
		{fmt.Errorf("raw error"), http.StatusInternalServerError}, // 500 fallback
	}

	for _, tc := range cases {
		t.Run(tc.err.Error(), func(t *testing.T) {
			assert.Equal(t, tc.wantHTTP, apperr.HTTPStatus(tc.err))
		})
	}
}

// TestLogFields kiểm tra output của LogFields — dùng để verify log structure.
func TestLogFields(t *testing.T) {
	t.Run("internal error includes stack", func(t *testing.T) {
		err := apperr.Internal("op", fmt.Errorf("db fail"))
		fields := apperr.LogFields(err)

		assert.Equal(t, "internal", fields["kind"])
		assert.Equal(t, "INTERNAL", fields["code"])
		assert.Equal(t, "op", fields["op"])
		assert.Equal(t, "db fail", fields["cause"])

		stack, hasStack := fields["stack"]
		assert.True(t, hasStack, "Internal LogFields MUST include stack")
		assert.NotEmpty(t, stack)
	})

	t.Run("validation error no stack", func(t *testing.T) {
		err := apperr.Validation("bad", apperr.NewFieldError("name", "REQUIRED_FIELD", "required", nil))
		fields := apperr.LogFields(err)

		assert.Equal(t, "validation", fields["kind"])
		_, hasStack := fields["stack"]
		assert.False(t, hasStack, "Validation LogFields MUST NOT include stack")

		// Fields có value vì là internal log
		logFields := fields["fields"]
		assert.NotNil(t, logFields)
	})

	t.Run("non-AppError returns kind=unknown", func(t *testing.T) {
		err := fmt.Errorf("some raw error")
		fields := apperr.LogFields(err)

		assert.Equal(t, "unknown", fields["kind"])
		assert.Equal(t, "UNKNOWN", fields["code"])
		assert.Equal(t, "some raw error", fields["error"])
	})

	t.Run("nil error returns nil", func(t *testing.T) {
		assert.Nil(t, apperr.LogFields(nil))
	})
}

// TestKind_Methods kiểm tra IsClientFault/IsSystemFault.
func TestKind_Methods(t *testing.T) {
	assert.True(t, apperr.KindValidation.IsClientFault())
	assert.True(t, apperr.KindBusiness.IsClientFault())
	assert.False(t, apperr.KindInternal.IsClientFault())
	assert.False(t, apperr.KindExternal.IsClientFault())

	assert.True(t, apperr.KindInternal.IsSystemFault())
	assert.False(t, apperr.KindValidation.IsSystemFault())
}

// TestAppError_Error kiểm tra Error() string format cho grep-ability.
func TestAppError_Error(t *testing.T) {
	err := apperr.Internal("UserService.Get:DB.Query", fmt.Errorf("connection refused"))
	msg := err.Error()

	assert.Contains(t, msg, "internal")
	assert.Contains(t, msg, "INTERNAL")
	assert.Contains(t, msg, "op=UserService.Get")
	assert.Contains(t, msg, "cause=connection refused")

	// Verify | separator
	parts := strings.Split(msg, " | ")
	assert.GreaterOrEqual(t, len(parts), 3)
}
