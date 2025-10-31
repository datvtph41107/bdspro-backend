package dto

import "auth/internal/domain"

// RequestLoginResponse đại diện cho response khi request OTP
type RequestLoginResponse struct {
	AuthID     uint64 `json:"authId"`
	Provider   string `json:"provider"`
	OAuthID    string `json:"oauthId"`
	Existed    bool   `json:"existed"`    // true: user đã tồn tại, false: user mới đăng ký
	SmsChannel string `json:"smsChannel"` // "zalo" nếu gửi ZNS thành công, "sms" nếu ZNS lỗi
}

// RefreshResponse đại diện cho response khi refresh token
type RefreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
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
