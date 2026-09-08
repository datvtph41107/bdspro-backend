package enums

import "fmt"

// AuthCodeEnum định nghĩa các mã trạng thái xác thực
type AuthCodeEnum int

const (
	LOCKED             AuthCodeEnum = 1001
	LIMIT_OTP_SEND     AuthCodeEnum = 1002
	LIMIT_OTP_ENTER    AuthCodeEnum = 1003
	CHECK_OTP          AuthCodeEnum = 1004
	EXPIRED            AuthCodeEnum = 1005
	LIMIT_NEXT_TIME    AuthCodeEnum = 1006
	LIMIT_REQUEST_TIME AuthCodeEnum = 1017
	ACTIVATED          AuthCodeEnum = 1016
	INACTIVE_ACCOUNT   AuthCodeEnum = 1030
	NOT_FOUND_ACCOUNT  AuthCodeEnum = 1044
	FULLNAME_REQUIRED  AuthCodeEnum = 1045
)

// authCodeMap ánh xạ mã số với tên enum
var authCodeMap = map[int]AuthCodeEnum{
	1001: LOCKED,
	1002: LIMIT_OTP_SEND,
	1003: LIMIT_OTP_ENTER,
	1004: CHECK_OTP,
	1005: EXPIRED,
	1006: LIMIT_NEXT_TIME,
	1017: LIMIT_REQUEST_TIME,
	1016: ACTIVATED,
	1030: INACTIVE_ACCOUNT,
	1044: NOT_FOUND_ACCOUNT,
}

// FromCode tìm enum tương ứng với mã số
func FromCode(code int) (AuthCodeEnum, error) {
	if status, exists := authCodeMap[code]; exists {
		return status, nil
	}
	return 0, fmt.Errorf("No enum constant with code %d", code)
}
