package dto

import "user/internal/domain"

// AuthParam đại diện cho thông tin xác thực
type AuthParam struct {
	Status  *domain.UserStatusEntity `json:"status"`
	OTP     *domain.UserOTPEntity    `json:"otp"`
	OTPCode string                   `json:"otpCode"`
	Auth    *domain.AuthMethod       `json:"auth,omitempty"`
}
