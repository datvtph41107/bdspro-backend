package dto

import "auth/internal/domain"

// RequestLoginResponse đại diện cho response khi request OTP
type RequestLoginResponse struct {
	AuthID   uint64 `json:"authId"`
	Provider string `json:"provider"`
	OAuthID  string `json:"oauthId"`
}

// RefreshResponse đại diện cho response khi refresh token
type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
}

// AuthParam đại diện cho thông tin xác thực
type AuthParam struct {
	Status  *domain.UserStatusEntity `json:"status"`
	OTP     *domain.UserOTPEntity    `json:"otp"`
	OTPCode string                   `json:"otpCode"`
}

// PublicUserInfo đại diện cho thông tin user công khai
type PublicUserInfo struct {
	ProfileID uint64 `json:"profileId"`
	FullName  string `json:"fullName"`
	Avatar    string `json:"avatar"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
