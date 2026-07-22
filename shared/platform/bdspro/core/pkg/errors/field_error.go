// platform/errors/field_error.go
package apperr

// FieldError mô tả lỗi validation của một field cụ thể.
//
// Tại sao có Value field nhưng Fields() không trả về nó?
// - Value dùng để log nội bộ (debug) — biết giá trị nào bị reject
// - Response ra client không cần Value — tránh leak data user input
// - Xem Fields() trong app_error.go để thấy cách strip Value
type FieldError struct {
	Field   string // Tên field: "email", "items[0].price"
	Code    string // Error code: "REQUIRED_FIELD", "INVALID_FORMAT"
	Message string // Human-readable: "Email không đúng định dạng"
	Value   any    // Giá trị bị lỗi — chỉ dùng nội bộ, không serialize ra client
}

// NewFieldError tạo FieldError với đủ thông tin.
func NewFieldError(field, code, message string, value any) FieldError {
	return FieldError{
		Field:   field,
		Code:    code,
		Message: message,
		Value:   value,
	}
}
