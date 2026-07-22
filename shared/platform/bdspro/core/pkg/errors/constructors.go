// platform/errors/constructors.go
//
// Package-level error variables (sentinel errors) và constructor functions.
//
// Nguyên tắc:
// 1. Sentinel errors là TEMPLATE — không dùng trực tiếp, customize bằng With*
// 2. Constructors có signature rõ ràng — không dùng variadic interface{}
// 3. Internal() LUÔN có stack trace (F10)
package apperr

import stderrors "errors"

// ─── Sentinel errors (templates) ─────────────────────────────────────────────
// Dùng với errors.Is():  errors.Is(err, ErrNotFound)
// Customize:             ErrNotFound.WithField("id", ...).WithOp("GetUser")

var (
	// Validation — 400
	ErrInvalidArgument = New(KindValidation, CodeInvalidArgument, "invalid argument")
	ErrInvalidJSON     = New(KindValidation, CodeInvalidJSON, "invalid json body")
	ErrRequiredField   = New(KindValidation, CodeRequiredField, "field is required")
	ErrInvalidFormat   = New(KindValidation, CodeInvalidFormat, "invalid format")
	ErrOutOfRange      = New(KindValidation, CodeOutOfRange, "value out of range")
	ErrTooLong         = New(KindValidation, CodeTooLong, "value too long")
	ErrTooShort        = New(KindValidation, CodeTooShort, "value too short")

	// Business — 400/401/403/404/409
	ErrNotFound      = New(KindBusiness, CodeNotFound, "resource not found")
	ErrAlreadyExists = New(KindBusiness, CodeAlreadyExists, "resource already exists")
	ErrForbidden     = New(KindBusiness, CodeForbidden, "forbidden")
	ErrUnauthorized  = New(KindBusiness, CodeUnauthorized, "unauthorized")
	ErrConflict      = New(KindBusiness, CodeConflict, "conflict")

	// External — 429/503/504
	ErrUnavailable = New(KindExternal, CodeUnavailable, "service unavailable").WithRetryable(true)
	ErrTimeout     = New(KindExternal, CodeTimeout, "timeout").WithRetryable(true)
	ErrRateLimit   = New(KindExternal, CodeRateLimit, "rate limit exceeded").WithRetryable(true)

	// Internal — 500
	ErrInternal = New(KindInternal, CodeInternal, "internal server error")
)

// ─── Constructor functions ────────────────────────────────────────────────────

// Validation tạo validation error với optional field errors.
// Signature rõ ràng hơn variadic interface{} — caller biết gì được pass.
//
// Ví dụ:
//
//	apperr.Validation("validation failed",
//	    apperr.NewFieldError("email", "INVALID_FORMAT", "Email không hợp lệ", req.Email),
//	)
func Validation(message string, fields ...FieldError) *AppError {
	e := New(KindValidation, CodeInvalidArgument, message)
	if len(fields) > 0 {
		e = e.WithFields(fields...)
	}
	return e
}

// Business tạo business error với custom code và message.
func Business(code Code, message string) *AppError {
	return New(KindBusiness, code, message)
}

// Internal tạo internal server error với op chain và root cause.
//
// LUÔN capture stack trace (F10) — không có ngoại lệ.
// Lý do: Internal error là server bug. Không có stack = không thể debug production.
//
// Ví dụ:
//
//	if err := db.QueryRow(...).Scan(&user); err != nil {
//	    return apperr.Internal("UserRepo.GetByID", err)
//	}
func Internal(op string, cause error) *AppError {
	if cause == nil {
		// caller có thể vô tình pass nil — tạo synthetic cause để không mất info
		cause = stderrors.New("internal error (nil cause)")
	}
	return New(KindInternal, CodeInternal, "internal server error").
		WithOp(op).
		WithCause(cause).
		WithStack() // F10: MANDATORY — always capture stack
}

// External tạo external service error.
// RetryAfter mặc định 0 — caller set nếu biết retry delay.
func External(code Code, message string, cause error) *AppError {
	if cause == nil {
		cause = stderrors.New(string(code))
	}
	return New(KindExternal, code, message).
		WithCause(cause).
		WithRetryable(true)
}
