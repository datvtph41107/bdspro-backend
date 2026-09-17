// Package validation mô tả lỗi dữ liệu đầu vào, độc lập HTTP/gRPC và logger.
package validation

import (
	"errors"
	"fmt"
	"runtime"
)

// Error chỉ chứa thông báo an toàn cho người dùng; không bọc lỗi SQL/network vào đây.
// Source dành cho log nội bộ tại boundary, không được trả ra response.
type Error struct {
	Field   string
	Reason  string
	Message string
	Source  string
}

func (e *Error) Error() string { return e.Message }

// New ghi nhận vị trí quy tắc bị vi phạm ngay tại owner, trước khi lỗi đi qua adapter.
func New(field, reason, message string) error {
	_, file, line, _ := runtime.Caller(1)
	return &Error{Field: field, Reason: reason, Message: message, Source: fmt.Sprintf("%s:%d", file, line)}
}

func As(err error) (*Error, bool) {
	var value *Error
	ok := errors.As(err, &value)
	return value, ok
}
