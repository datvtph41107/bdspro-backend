package businessoperation

import (
	"errors"
	"strings"
)

const MaxCodeLength = 128

var ErrInvalidCode = errors.New("business operation code is invalid")

/**
 * Code là tên ổn định của một business action.
 *
 * Ví dụ:
 *   workspace.report.generate
 *   parcel.view
 */
type Code string

/**
 * ParseCode kiểm tra và trả về một business operation hợp lệ.
 */
func ParseCode(value string) (Code, error) {
	value = strings.TrimSpace(value)
	code := Code(value)
	if !code.IsValid() {
		return "", ErrInvalidCode
	}
	return code, nil
}

/**
 * IsValid kiểm tra format của business operation.
 *
 * Code chỉ dùng lowercase, số, dấu gạch dưới và dấu chấm.
 * Mỗi code phải có ít nhất hai phần để tránh tên quá chung chung.
 */
func (c Code) IsValid() bool {
	value := string(c)
	if value == "" || len(value) > MaxCodeLength || strings.TrimSpace(value) != value {
		return false
	}

	parts := strings.Split(value, ".")
	if len(parts) < 2 {
		return false
	}

	for _, part := range parts {
		if !validPart(part) {
			return false
		}
	}

	return true
}

func validPart(value string) bool {
	if value == "" || value[0] < 'a' || value[0] > 'z' {
		return false
	}

	for i := 1; i < len(value); i++ {
		c := value[i]
		if (c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') ||
			c == '_' {
			continue
		}
		return false
	}

	return true
}
