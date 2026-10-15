package dto

import "userv4/modules/auth/domain/auth"

// AuthParam đại diện cho thông tin xác thực
type AuthParam struct {
	Status  *auth.AuthStatus `json:"status"`
	OTP     *auth.UserOTPEntity    `json:"otp"`
	OTPCode string                 `json:"otpCode"`
	Auth    *auth.AuthMethod       `json:"auth,omitempty"`
}
