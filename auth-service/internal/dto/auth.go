package dto

import "auth/internal/domain"

// Hàm khởi tạo AuthLoginResponse từ UserAuthEntity
func MakeAuthLoginResponse(accessToken string, authEntity *domain.AuthMethodDomain, role string) *AuthLoginResponse {
	return &AuthLoginResponse{
		AccessToken: accessToken,
		// RefreshToken: refreshToken,
		Role:      role,
		Username:  authEntity.Username,
		AuthID:    authEntity.ID,
		ProfileID: authEntity.UserID,
	}
}

// OtpRequest đại diện cho yêu cầu gửi OTP
type OtpRequest struct {
	Phone    string `json:"phone" binding:"required" validate:"required,e164"`
	Fullname string `json:"fullname,omitempty"`
}

// OtpVerifyRequest mở rộng từ OtpRequest để xác minh OTP
type OtpVerifyRequest struct {
	OtpRequest
	Otp      string `json:"otp" binding:"required" validate:"required"`
	AuthID   uint64 `json:"authId,omitempty"`
	Version  string `json:"version,omitempty"`
	Platform string `json:"platform,omitempty"`
	OS       string `json:"os,omitempty"`
	Fullname string `json:"fullname,omitempty"`
}

// RequestOTPResponse đại diện cho response gửi OTP
type RequestOTPResponse struct {
	AuthId  uint64 `json:"authId"`
	Message string `json:"message"`
}

// VerifyOTPResponse đại diện cho response xác minh OTP
type VerifyOTPResponse struct {
	AuthId       uint64 `json:"authId"`
	ProfileId    uint64 `json:"profileId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	FullName     string `json:"fullName"`
	Avatar       string `json:"avatar"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
}

// ResendOTPRequest đại diện cho yêu cầu gửi lại OTP
type ResendOTPRequest struct {
	Phone    string `json:"phone" binding:"required" validate:"required,e164"`
	Fullname string `json:"fullname,omitempty"`
}

// ResendOTPResponse đại diện cho response gửi lại OTP
type ResendOTPResponse struct {
	AuthId  uint64 `json:"authId"`
	Message string `json:"message"`
}

// RefreshTokenRequest đại diện cho yêu cầu làm mới token
type RefreshTokenRequest struct {
	AccessToken  string `json:"accessToken" binding:"required" validate:"required"`
	RefreshToken string `json:"refreshToken" binding:"required" validate:"required"`
}

// RefreshTokenResponse đại diện cho response làm mới token
type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// LogoutRequest đại diện cho yêu cầu đăng xuất
type LogoutRequest struct {
	// Empty request
}

// LogoutResponse đại diện cho response đăng xuất
type LogoutResponse struct {
	Message string `json:"message"`
}

// DeleteAccountRequest đại diện cho yêu cầu xóa tài khoản
type DeleteAccountRequest struct {
	// Empty request
}

// DeleteAccountResponse đại diện cho response xóa tài khoản
type DeleteAccountResponse struct {
	Message string `json:"message"`
}

// RestorePhoneRequest đại diện cho yêu cầu khôi phục số điện thoại
type RestorePhoneRequest struct {
	Phone string `json:"phone" binding:"required" validate:"required,e164"`
}

// RestoreAccountResponse đại diện cho response khôi phục tài khoản
type RestoreAccountResponse struct {
	AuthId   uint64 `json:"authId"`
	Provider string `json:"provider"`
	OAuthId  string `json:"oauthId"`
}

// SessionQRRequest đại diện cho yêu cầu tạo QR session
type SessionQRRequest struct {
	ClientID   string `json:"clientId"`
	UserAgent  string `json:"userAgent"`
	CsrfToken  string `json:"csrfToken"`
	Expiration int64  `json:"expiration"`
	SessionKey string `json:"sessionKey"`
	SessionId  uint64 `json:"sessionId"`
}

// SessionQRResponse đại diện cho response QR session
type SessionQRResponse struct {
	SessionKey string `json:"sessionKey"`
	SessionId  uint64 `json:"sessionId"`
}

// SessionQRConfirm đại diện cho yêu cầu xác nhận QR
type SessionQRConfirm struct {
	SessionId  uint64 `json:"sessionId"`
	SessionKey string `json:"sessionKey"`
}

// QRConfirmResponse đại diện cho response xác nhận QR
type QRConfirmResponse struct {
	Message string `json:"message"`
}

// QRStatusRequest đại diện cho yêu cầu kiểm tra trạng thái QR
type QRStatusRequest struct {
	SessionKey string `json:"sessionKey"`
	SessionId  uint64 `json:"sessionId"`
}

// QRStatusResponse đại diện cho response trạng thái QR
type QRStatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// AuthLoginResponse đại diện cho response đăng nhập
type AuthLoginResponse struct {
	AccessToken    string `json:"accessToken"`
	AuthID         uint64 `json:"authId"`
	Role           string `json:"role,omitempty"`
	Username       string `json:"username,omitempty"`
	ProfileID      uint64 `json:"profileId"`
	RefreshToken   string `json:"refreshToken"`
	OrganizationID uint64 `json:"organizationId,omitempty"`
	// FullName       string `json:"fullName"`
	// Avatar         string `json:"avatar"`
	// Phone          string `json:"phone"`
	// Email          string `json:"email"`
}

// OAuthLoginRequest đại diện cho yêu cầu đăng nhập OAuth
type OAuthLoginRequest struct {
	Provider string `json:"provider" binding:"required" validate:"required"`
	Code     string `json:"code" binding:"required" validate:"required"`
	State    string `json:"state"`
	Platform string `json:"platform"`
}

// OAuthLoginResponse đại diện cho response đăng nhập OAuth
type OAuthLoginResponse struct {
	AuthID       uint64 `json:"authId"`
	ProfileID    uint64 `json:"profileId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	FullName     string `json:"fullName"`
	Avatar       string `json:"avatar"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Provider     string `json:"provider"`
}

// GetLoginURLRequest đại diện cho yêu cầu lấy URL đăng nhập OAuth
type GetLoginURLRequest struct {
	Provider string `json:"provider" binding:"required" validate:"required"`
}

// GetLoginURLResponse đại diện cho response URL đăng nhập OAuth
type GetLoginURLResponse struct {
	URL   string `json:"url"`
	State string `json:"state"`
}

// HandleCallbackRequest đại diện cho yêu cầu xử lý OAuth callback
type HandleCallbackRequest struct {
	Provider string `json:"provider" binding:"required" validate:"required"`
	Code     string `json:"code" binding:"required" validate:"required"`
	State    string `json:"state"`
	Platform string `json:"platform"`
}

// HandleCallbackResponse đại diện cho response xử lý OAuth callback
type HandleCallbackResponse struct {
	AuthId       uint64 `json:"authId"`
	ProfileId    uint64 `json:"profileId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	FullName     string `json:"fullName"`
	Avatar       string `json:"avatar"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Provider     string `json:"provider"`
}
