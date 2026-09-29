package _utils

import (
	"regexp"
	"strconv"
	"strings"
)

// ParseInt32Array parse string "1,2,3" thành []int32{1, 2, 3}
// Dùng cho các query params như friendStatus, roles, etc.
func ParseInt32Array(str string) []int32 {
	if str == "" {
		return nil
	}

	parts := strings.Split(str, ",")
	result := make([]int32, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		num, err := strconv.ParseInt(part, 10, 32)
		if err == nil {
			result = append(result, int32(num))
		}
	}

	return result
}

// ParseInt64Array parse string "1,2,3" thành []int64{1, 2, 3}
func ParseInt64Array(str string) []int64 {
	if str == "" {
		return nil
	}

	parts := strings.Split(str, ",")
	result := make([]int64, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		num, err := strconv.ParseInt(part, 10, 64)
		if err == nil {
			result = append(result, num)
		}
	}

	return result
}

// ParseUint64Array parse string "1,2,3" thành []uint64{1, 2, 3}
func ParseUint64Array(str string) []uint64 {
	if str == "" {
		return nil
	}

	parts := strings.Split(str, ",")
	result := make([]uint64, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		num, err := strconv.ParseUint(part, 10, 64)
		if err == nil {
			result = append(result, num)
		}
	}

	return result
}

// ParseStringArray parse string "a,b,c" thành []string{"a", "b", "c"}
// Tự động trim space cho mỗi phần tử
func ParseStringArray(str string) []string {
	if str == "" {
		return nil
	}

	parts := strings.Split(str, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}

	return result
}

// ParseBool parse string "true", "1", "yes" thành true
// "false", "0", "no", "" thành false
func ParseBool(str string) bool {
	str = strings.ToLower(strings.TrimSpace(str))
	switch str {
	case "true", "1", "yes", "y":
		return true
	default:
		return false
	}
}

// ParseFloat64 parse string thành float64, trả về 0 nếu parse lỗi
func ParseFloat64(str string) float64 {
	str = strings.TrimSpace(str)
	if str == "" {
		return 0
	}

	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return val
}

// ParseInt parse string thành int, trả về 0 nếu parse lỗi
func ParseInt(str string) int {
	str = strings.TrimSpace(str)
	if str == "" {
		return 0
	}

	val, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return val
}

// ParseInt64 parse string thành int64, trả về 0 nếu parse lỗi
func ParseInt64(str string) int64 {
	str = strings.TrimSpace(str)
	if str == "" {
		return 0
	}

	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return val
}

// ParseUint64 parse string thành uint64, trả về 0 nếu parse lỗi
func ParseUint64(str string) uint64 {
	str = strings.TrimSpace(str)
	if str == "" {
		return 0
	}

	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0
	}
	return val
}

// IsValidEmail kiểm tra email có hợp lệ không
func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}

	// Regex pattern cho email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidPhone kiểm tra số điện thoại có hợp lệ không (Việt Nam)
func IsValidPhone(phone string) bool {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return false
	}

	// Regex cho số điện thoại Việt Nam (10-11 số, bắt đầu bằng 0 hoặc +84)
	phoneRegex := regexp.MustCompile(`^(\+84|0)[0-9]{9,10}$`)
	return phoneRegex.MatchString(phone)
}
