package request

import "strings"

const (
	RequestIDHeader      = "X-Request-ID"
	RequestIDMetadataKey = "x-request-id"
	MaxRequestIDLength   = 128
)

// RequestIDSource canonical source
type RequestIDSource string

const (
	RequestIDSourceClient    RequestIDSource = "client"
	RequestIDSourceGenerated RequestIDSource = "generated"
)

// RequestIDReason xác định request input
/**
 * Description: Có bốn khả năng
 * 1. Gửi một ID hợp lệ
 * 2. Không gửi X-Request-ID
 * 3. Gửi một ID không an toàn
 * 4. Gửi nhiều giá trị X-Request-ID
 */
type RequestIDReason string

const (
	RequestIDReasonValid    RequestIDReason = "valid"
	RequestIDReasonMissing  RequestIDReason = "missing"  // tạo mới_req_id
	RequestIDReasonInvalid  RequestIDReason = "invalid"  // tạo mới_req_id
	RequestIDReasonMultiple RequestIDReason = "multiple" // tạo mới req_id
)

// RequestIDDecision kết quả canonicalization tiêu chuẩn
/**
 * ID là giá trị duy nhất được phép truyền tiếp
 * Source và Reason dùng cho logs
 */
type RequestIDDecision struct {
	ID     string
	Source RequestIDSource // client gửi hoặc được tạo ở server
	Reason RequestIDReason
}

// NormalizeRequestID remove whitespace hai đầu
// Opaque identifier Không lowercase, Không uppercase, Không thay đổi ký tự bên trong
// AbC123-XyZ9
func NormalizeRequestID(raw string) string {
	return strings.TrimSpace(raw)
}

// IsValidRequestID kiểm tra request ID
// khi đi qua HTTP, gRPC metadata và structured logs không.
func IsValidRequestID(value string) bool {
	if value == "" {
		return false
	}

	if len(value) > MaxRequestIDLength {
		return false
	}

	for index := 0; index < len(value); index++ {
		if !isAllowedRequestIDCharacter(value[index]) {
			return false
		}
	}

	return true
}

/**
 * Description: Canonical ID
 * không chứa: newline, space, /
 */
func isAllowedRequestIDCharacter(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '-' ||
		c == '_' ||
		c == '.' ||
		c == ':'
}

// ResolveRequestID chuyển zero, one hoặc multiple values
// thành một canonical decision.

/**
 * Note: lấy header
 * gRPC metadata -> values := md.Get(RequestIDMetadataKey)
 * http -> values := r.Header.Values(RequestIDHeader)
 * đều trả về []string
 */
func ResolveRequestID(values []string) RequestIDDecision {
	return resolveRequestID(values, NewRequestID)
}

func resolveRequestID(
	values []string,
	generate func() string,
) RequestIDDecision {
	// note: len(nil slice) == 0 - ResolveRequestID(nil) -> []string(nil)
	// nil slice == empty slice
	switch len(values) {
	case 0:
		return generatedRequestIDDecision(
			RequestIDReasonMissing,
			generate,
		)
	case 1:
		normalized := NormalizeRequestID(values[0])

		if IsValidRequestID(normalized) {
			return RequestIDDecision{
				ID:     normalized,
				Source: RequestIDSourceClient,
				Reason: RequestIDReasonValid,
			}
		}

		return generatedRequestIDDecision(
			RequestIDReasonInvalid,
			generate,
		)

	default:
		return generatedRequestIDDecision(
			RequestIDReasonMultiple,
			generate,
		)
	}

}

func generatedRequestIDDecision(
	reason RequestIDReason,
	generate func() string,
) RequestIDDecision {
	requestID := generate()

	if !IsValidRequestID(requestID) {
		panic(
			"request: generator returned invalid request ID",
		)
	}

	return RequestIDDecision{
		ID:     requestID,
		Source: RequestIDSourceGenerated,
		Reason: reason,
	}
}
