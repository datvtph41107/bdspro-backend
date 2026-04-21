// platform/errors/app_error.go
package apperr

import (
	"runtime/debug"
	"strings"
)

// AppError là error type duy nhất của platform.
//
// Thiết kế immutable-by-convention: mọi With* method trả về bản copy mới.
// Lý do: error thường được tạo dưới dạng package-level var rồi được "customize".
// Nếu mutate in-place: ErrNotFound.WithField("id", ...) sẽ ảnh hưởng global var.
//
// Tại sao struct (không phải interface)?
// - Concrete type → compiler check đầy đủ
// - Không cần mock AppError trong test
// - errors.As(*AppError) hoạt động tự nhiên
type AppError struct {
	kind       Kind
	code       Code
	message    string       // User-facing message
	fields     []FieldError // Validation field errors
	cause      error        // Wrapped original error
	op         string       // Operation chain: "CreateUser:ValidateEmail"
	retryable  bool         // Client có thể retry?
	retryAfter int          // Retry sau bao nhiêu giây
	stack      []byte       // Stack trace — chỉ cho Internal errors
}

// New tạo AppError cơ bản. Thường không dùng trực tiếp — dùng constructors.
func New(kind Kind, code Code, message string) *AppError {
	return &AppError{
		kind:    kind,
		code:    code,
		message: message,
	}
}

// Error implement error interface.
// Format: "internal | INTERNAL | op=CreateUser:DB | cause=connection refused"
// Mục đích: dễ grep trong log khi chưa có structured logging.
func (e *AppError) Error() string {
	parts := make([]string, 0, 5)
	parts = append(parts, e.kind.String())
	parts = append(parts, string(e.code))
	if e.op != "" {
		parts = append(parts, "op="+e.op)
	}
	if e.message != "" {
		parts = append(parts, e.message)
	}
	if e.cause != nil {
		parts = append(parts, "cause="+e.cause.Error())
	}
	return strings.Join(parts, " | ")
}

// Unwrap implement errors.Unwrap — cho phép errors.Is/As hoạt động với chain.
func (e *AppError) Unwrap() error { return e.cause }

// Is implement errors.Is — chỉ match khi kind VÀ code đều khớp.
//
// Tại sao không chỉ check code?
// Code "INTERNAL" thuộc KindInternal — không thể Is với KindExternal dù code giống.
// Pattern an toàn: errors.Is(err, ErrNotFound) → chỉ true nếu là business not found.
//
// Tại sao không check pointer equality?
// Vì AppError dùng immutable clone pattern — mỗi With* call tạo pointer mới.
// errors.Is(Internal("op", err), ErrInternal) phải return true.
func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok || e == nil || t == nil {
		return false
	}
	// Chỉ match kind+code, không check message/op/cause
	// → cho phép customize message mà vẫn Is match
	return e.kind == t.kind && e.code == t.code
}

// ─── Getters (read-only) ─────────────────────────────────────────────────────
// Không expose fields trực tiếp — caller không nên mutate AppError.

func (e *AppError) Kind() Kind      { return e.kind }
func (e *AppError) Code() Code      { return e.code }
func (e *AppError) Message() string { return e.message }
func (e *AppError) Op() string      { return e.op }
func (e *AppError) Retryable() bool { return e.retryable }
func (e *AppError) RetryAfter() int { return e.retryAfter }
func (e *AppError) Cause() error    { return e.cause }
func (e *AppError) Stack() []byte   { return e.stack }
func (e *AppError) HasFields() bool { return len(e.fields) > 0 }

// Fields trả về bản sao của field errors với Value bị strip.
//
// Tại sao strip Value?
// Fields() được dùng trong HTTP response — không được leak user input ra ngoài.
// Internal log dùng LogFields() — có đầy đủ thông tin kể cả Value.
func (e *AppError) Fields() []FieldError {
	if len(e.fields) == 0 {
		return nil
	}
	out := make([]FieldError, len(e.fields))
	for i, f := range e.fields {
		out[i] = FieldError{
			Field:   f.Field,
			Code:    f.Code,
			Message: f.Message,
			// Value intentionally omitted — security: không leak input ra response
		}
	}
	return out
}

// ─── Immutable builders ──────────────────────────────────────────────────────
// Tất cả With* method return con trỏ mới — không mutate receiver.

// clone tạo shallow copy của AppError.
// slice field cần copy riêng để tránh shared backing array.
func (e *AppError) clone() *AppError {
	cp := *e // value copy
	if len(e.fields) > 0 {
		cp.fields = make([]FieldError, len(e.fields))
		copy(cp.fields, e.fields)
	}
	// stack không cần copy — []byte readonly sau debug.Stack()
	return &cp
}

// WithStack capture call stack tại thời điểm gọi.
// Chỉ Internal errors cần stack — validation/business errors không phải bug.
func (e *AppError) WithStack() *AppError {
	cp := e.clone()
	cp.stack = debug.Stack() // runtime/debug.Stack() — không require cgo
	return cp
}

// WithOp thêm operation vào chain.
// Chain format: "A:B:C" — đọc từ trái sang phải là call hierarchy.
// Ví dụ: Internal("DB.Query", err).WithOp("UserRepo.Create").WithOp("UserService.Register")
// → op = "UserService.Register:UserRepo.Create:DB.Query"
func (e *AppError) WithOp(op string) *AppError {
	cp := e.clone()
	if cp.op == "" {
		cp.op = op
	} else {
		cp.op = op + ":" + cp.op // prepend để đọc outermost first
	}
	return cp
}

// WithCause wraps original error — preserve root cause để debug.
func (e *AppError) WithCause(cause error) *AppError {
	cp := e.clone()
	cp.cause = cause
	return cp
}

// WithMessage override user-facing message.
func (e *AppError) WithMessage(msg string) *AppError {
	cp := e.clone()
	cp.message = msg
	return cp
}

// WithRetryable đánh dấu error có thể retry.
func (e *AppError) WithRetryable(v bool) *AppError {
	cp := e.clone()
	cp.retryable = v
	return cp
}

// WithRetryAfter set retry delay và tự động mark retryable=true.
func (e *AppError) WithRetryAfter(sec int) *AppError {
	cp := e.clone()
	cp.retryable = true
	cp.retryAfter = sec
	return cp
}

// WithField thêm một field validation error.
func (e *AppError) WithField(field, code, message string, value any) *AppError {
	cp := e.clone()
	cp.fields = append(cp.fields, NewFieldError(field, code, message, value))
	return cp
}

// WithFields thêm nhiều field validation errors cùng lúc.
// Trả về e nếu fields rỗng — không tạo copy không cần thiết.
func (e *AppError) WithFields(fields ...FieldError) *AppError {
	if len(fields) == 0 {
		return e
	}
	cp := e.clone()
	cp.fields = append(cp.fields, fields...)
	return cp
}

// ─── Helper functions ────────────────────────────────────────────────────────

// HTTPStatus lấy HTTP status code từ error.
// Hàm này là NGUỒN SỰ THẬT cho HTTP mapping — Writer.Error() chỉ gọi hàm này.
func HTTPStatus(err error) int {
	ae, ok := As(err)
	if !ok {
		return 500
	}
	return ae.Code().ToHTTPStatus()
}
