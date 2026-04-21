// platform/errors/code.go
package apperr

import "net/http"

// Code là error code string — xuất hiện trong JSON response và log.
//
// Tại sao string không phải int?
// - Client code (frontend, mobile) check bằng string match — dễ đọc
// - HTTP status code đã là int — thêm int error code gây confusion
// - Dễ add code mới mà không break existing client (nếu dùng int, có thể trùng)
//
// Convention: SCREAMING_SNAKE_CASE, noun phrase, không có "ERROR" suffix.
type Code string

const (
	// ─── Validation (400) ────────────────────────────────────────────────────
	CodeInvalidArgument Code = "INVALID_ARGUMENT" // Generic validation fail
	CodeInvalidJSON     Code = "INVALID_JSON"     // Malformed JSON body
	CodeRequiredField   Code = "REQUIRED_FIELD"   // Bắt buộc nhưng thiếu
	CodeInvalidFormat   Code = "INVALID_FORMAT"   // Sai format (email, phone, ...)
	CodeOutOfRange      Code = "OUT_OF_RANGE"     // Giá trị ngoài range cho phép
	CodeTooLong         Code = "TOO_LONG"         // String/array quá dài
	CodeTooShort        Code = "TOO_SHORT"        // String/array quá ngắn

	// ─── Business (400/401/403/404/409) ─────────────────────────────────────
	CodeNotFound      Code = "NOT_FOUND"      // Resource không tồn tại
	CodeAlreadyExists Code = "ALREADY_EXISTS" // Resource đã tồn tại (duplicate)
	CodeForbidden     Code = "FORBIDDEN"      // Có auth nhưng thiếu quyền
	CodeUnauthorized  Code = "UNAUTHORIZED"   // Chưa auth
	CodeConflict      Code = "CONFLICT"       // State conflict (optimistic lock)

	// ─── External (503/504/429) ──────────────────────────────────────────────
	CodeUnavailable Code = "UNAVAILABLE" // Downstream service down
	CodeTimeout     Code = "TIMEOUT"     // Downstream timeout
	CodeRateLimit   Code = "RATE_LIMIT"  // Client bị throttle

	// ─── Internal (500) ──────────────────────────────────────────────────────
	CodeInternal Code = "INTERNAL" // Server bug — không expose detail ra client
)

// ToHTTPStatus map Code → HTTP status.
// Đây là NGUỒN SỰ THẬT DUY NHẤT cho mapping này.
// Writer.Error() và bất kỳ nơi nào cần HTTP status đều gọi hàm này.
//
// Tại sao không đặt trực tiếp trong Writer?
// Vì logic mapping thuộc về domain error, không thuộc về transport layer.
// Test error package không cần import transport package.
func (c Code) ToHTTPStatus() int {
	switch c {
	case CodeInvalidArgument, CodeInvalidJSON, CodeRequiredField,
		CodeInvalidFormat, CodeOutOfRange, CodeTooLong, CodeTooShort:
		return http.StatusBadRequest // 400
	case CodeUnauthorized:
		return http.StatusUnauthorized // 401
	case CodeForbidden:
		return http.StatusForbidden // 403
	case CodeNotFound:
		return http.StatusNotFound // 404
	case CodeAlreadyExists, CodeConflict:
		return http.StatusConflict // 409
	case CodeRateLimit:
		return http.StatusTooManyRequests // 429
	case CodeUnavailable:
		return http.StatusServiceUnavailable // 503
	case CodeTimeout:
		return http.StatusGatewayTimeout // 504
	default:
		return http.StatusInternalServerError // 500
	}
}
