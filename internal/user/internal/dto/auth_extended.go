package dto

import "user/internal/domain/auth"

// AuthParam đại diện cho thông tin xác thực
type AuthParam struct {
	Status  *auth.UserStatusEntity `json:"status"`
	OTP     *auth.UserOTPEntity    `json:"otp"`
	OTPCode string                 `json:"otpCode"`
	Auth    *auth.AuthMethod       `json:"auth,omitempty"`
}
