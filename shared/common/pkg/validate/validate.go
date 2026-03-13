// Package validate giải quyết bài toán:
// "Giá trị này có hợp lệ về mặt format không?"
//
// # Bản chất — 2 loại validation tách biệt hoàn toàn
//
//	TYPE 1 — FORMAT (package này):
//	  Pure function, không cần DB, không biết gRPC tồn tại.
//	  Chạy tại Handler, TRƯỚC mapper. Test trivial.
//	  Ví dụ: DateStringOpt, StringMaxLen, Int32Enum
//
//	TYPE 2 — BUSINESS (usecase/repo layer):
//	  Cần DB. Chạy trong usecase sau khi load existing.
//	  Ví dụ: "project_id có tồn tại không?"
//
// # Convention — mọi Rule function
//
//   - Nhận *T: nil = absent = skip (không validate)
//   - Pure: không side effect, không global state
//   - Tên mô tả đúng rule, không mô tả field
//   - Tích lũy lỗi qua *Error — không dừng ở lỗi đầu tiên
//
// # Không có ToGRPC() ở đây
//
// validate package không biết gRPC tồn tại — tránh coupling.
// Transport layer tự convert *validate.Error sang gRPC status.
//
// # Thread safety
//
// *Error KHÔNG thread-safe — dùng 1 instance per request.
// Rule functions là pure — safe for concurrent use.
package validate

import (
	"fmt"
	"strings"
	"time"
)

// ─────────────────────────────────────────────────────────────────────────────
// Error — accumulator
// ─────────────────────────────────────────────────────────────────────────────

// Error tích lũy tất cả lỗi validation.
// Implement error interface — transport layer tự xử lý conversion.
//
// Không dùng Error như map — field có thể lỗi nhiều rule.
// Dùng Add/Merge để tích lũy, HasErrors để check, Error() để lấy message.
type Error struct {
	errs []fieldError
}

type fieldError struct {
	field   string
	message string
}

// New tạo Error mới — zero value cũng dùng được nhưng New() rõ ý định hơn.
func New() *Error { return &Error{} }

// Add thêm lỗi cho field.
func (e *Error) Add(field, message string) {
	e.errs = append(e.errs, fieldError{field: field, message: message})
}

// Addf thêm lỗi với format string.
func (e *Error) Addf(field, format string, args ...any) {
	e.Add(field, fmt.Sprintf(format, args...))
}

// HasErrors trả về true nếu có ít nhất 1 lỗi.
func (e *Error) HasErrors() bool { return len(e.errs) > 0 }

// Merge hợp nhất tất cả lỗi từ other vào e.
// other nil → no-op.
func (e *Error) Merge(other *Error) {
	if other == nil || !other.HasErrors() {
		return
	}
	e.errs = append(e.errs, other.errs...)
}

// Error implement error interface.
// Format: "validation: field1: msg1; field2: msg2"
func (e *Error) Error() string {
	if !e.HasErrors() {
		return ""
	}
	parts := make([]string, 0, len(e.errs))
	for _, f := range e.errs {
		parts = append(parts, f.field+": "+f.message)
	}
	return "validation: " + strings.Join(parts, "; ")
}

// Fields trả về danh sách (field, message) — dùng cho transport conversion.
//
//	for _, f := range e.Fields() {
//	    // build gRPC BadRequest detail
//	}
func (e *Error) Fields() []FieldViolation {
	result := make([]FieldViolation, len(e.errs))
	for i, f := range e.errs {
		result[i] = FieldViolation{Field: f.field, Message: f.message}
	}
	return result
}

// FieldViolation là 1 lỗi cụ thể — transport-agnostic.
type FieldViolation struct {
	Field   string
	Message string
}

// AsError trả về error interface nếu có lỗi, nil nếu không.
// Dùng cho early return pattern:
//
//	if err := e.AsError(); err != nil { return err }
func (e *Error) AsError() error {
	if !e.HasErrors() {
		return nil
	}
	return e
}

// ─────────────────────────────────────────────────────────────────────────────
// Rule functions — pure, không side effect
//
// Convention:
//   - *T nil = absent = skip
//   - Tên = tên rule, không phải tên field
//   - Không return — tích lũy vào *Error
// ─────────────────────────────────────────────────────────────────────────────

// Required: phải present VÀ không rỗng sau khi trim space.
func Required(e *Error, field string, v *string) {
	if v == nil || strings.TrimSpace(*v) == "" {
		e.Add(field, "required")
	}
}

// RequiredID: uint64 ID phải present VÀ > 0.
func RequiredID(e *Error, field string, v *uint64) {
	if v == nil || *v == 0 {
		e.Add(field, "required")
	}
}

// RequiredInt32: int32 phải present (0 là giá trị hợp lệ).
func RequiredInt32(e *Error, field string, v *int32) {
	if v == nil {
		e.Add(field, "required")
	}
}

// ── String rules ──────────────────────────────────────────────────────────────

// StringMaxLen: nil=skip | len(*v) <= max.
func StringMaxLen(e *Error, field string, v *string, max int) {
	if v == nil {
		return
	}
	if len(*v) > max {
		e.Addf(field, "max length is %d, got %d", max, len(*v))
	}
}

// StringMinLen: nil=skip | len(*v) >= min.
func StringMinLen(e *Error, field string, v *string, min int) {
	if v == nil {
		return
	}
	if len(*v) < min {
		e.Addf(field, "min length is %d, got %d", min, len(*v))
	}
}

// StringBetween: nil=skip | min <= len(*v) <= max.
func StringBetween(e *Error, field string, v *string, min, max int) {
	if v == nil {
		return
	}
	n := len(*v)
	if n < min || n > max {
		e.Addf(field, "length must be between %d and %d, got %d", min, max, n)
	}
}

// StringOneOf: nil=skip | *v phải nằm trong allowed.
func StringOneOf(e *Error, field string, v *string, allowed []string) {
	if v == nil {
		return
	}
	for _, a := range allowed {
		if *v == a {
			return
		}
	}
	e.Addf(field, "must be one of [%s], got %q", strings.Join(allowed, ", "), *v)
}

// ── Numeric rules ─────────────────────────────────────────────────────────────

// Int32OneOf: nil=skip | *v phải nằm trong allowed.
func Int32OneOf(e *Error, field string, v *int32, allowed []int32) {
	if v == nil {
		return
	}
	for _, a := range allowed {
		if *v == a {
			return
		}
	}
	strs := make([]string, len(allowed))
	for i, a := range allowed {
		strs[i] = fmt.Sprintf("%d", a)
	}
	e.Addf(field, "must be one of [%s], got %d", strings.Join(strs, ", "), *v)
}

// Int32Between: nil=skip | min <= *v <= max.
func Int32Between(e *Error, field string, v *int32, min, max int32) {
	if v == nil {
		return
	}
	if *v < min || *v > max {
		e.Addf(field, "must be between %d and %d, got %d", min, max, *v)
	}
}

// Float64NonNegative: nil=skip | *v >= 0.
func Float64NonNegative(e *Error, field string, v *float64) {
	if v == nil {
		return
	}
	if *v < 0 {
		e.Addf(field, "must be >= 0, got %v", *v)
	}
}

// Float64Between: nil=skip | min <= *v <= max.
func Float64Between(e *Error, field string, v *float64, min, max float64) {
	if v == nil {
		return
	}
	if *v < min || *v > max {
		e.Addf(field, "must be between %v and %v, got %v", min, max, *v)
	}
}

// Uint32Positive: nil=skip | *v > 0.
func Uint32Positive(e *Error, field string, v *uint32) {
	if v == nil {
		return
	}
	if *v == 0 {
		e.Add(field, "must be > 0")
	}
}

// ── Date rules ────────────────────────────────────────────────────────────────

// DateStringOpt: nil=skip | ""=skip(clear intent hợp lệ) | RFC3339.
//
// Bài toán: DateString dùng "" = clear date → "" không phải lỗi format.
// Chỉ validate khi có giá trị thực sự.
func DateStringOpt(e *Error, field string, v *string) {
	if v == nil || *v == "" {
		return // absent hoặc clear intent
	}
	if _, err := time.Parse(time.RFC3339, *v); err != nil {
		e.Add(field, `invalid date, expected RFC3339 (e.g. "2024-01-15T00:00:00Z")`)
	}
}

// DateStringRequired: phải present, không rỗng, phải RFC3339.
func DateStringRequired(e *Error, field string, v *string) {
	if v == nil || *v == "" {
		e.Add(field, "required")
		return
	}
	if _, err := time.Parse(time.RFC3339, *v); err != nil {
		e.Add(field, `invalid date, expected RFC3339 (e.g. "2024-01-15T00:00:00Z")`)
	}
}

// DateOrder: nil=skip | start phải trước end khi cả 2 đều present và non-empty.
func DateOrder(e *Error, startField, endField string, start, end *string) {
	if start == nil || end == nil || *start == "" || *end == "" {
		return
	}
	s, err1 := time.Parse(time.RFC3339, *start)
	en, err2 := time.Parse(time.RFC3339, *end)
	if err1 != nil || err2 != nil {
		return // format error đã được validate riêng
	}
	if !s.Before(en) {
		e.Addf(startField, "must be before %s", endField)
	}
}

// ── Geo rules ─────────────────────────────────────────────────────────────────

// Latitude: nil=skip | -90 <= *v <= 90.
func Latitude(e *Error, field string, v *float64) {
	Float64Between(e, field, v, -90, 90)
}

// Longitude: nil=skip | -180 <= *v <= 180.
func Longitude(e *Error, field string, v *float64) {
	Float64Between(e, field, v, -180, 180)
}

// ── URL rules ─────────────────────────────────────────────────────────────────

// URLOpt: nil=skip | ""=skip | phải bắt đầu bằng http:// hoặc https://.
func URLOpt(e *Error, field string, v *string) {
	if v == nil || *v == "" {
		return
	}
	if !strings.HasPrefix(*v, "http://") && !strings.HasPrefix(*v, "https://") {
		e.Add(field, "must be a valid URL starting with http:// or https://")
	}
}

// ── Media rules ───────────────────────────────────────────────────────────────

// MediaType: nil=skip | phải là "image" hoặc "video".
func MediaType(e *Error, field string, v *string) {
	StringOneOf(e, field, v, []string{"image", "video"})
}
