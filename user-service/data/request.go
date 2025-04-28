package data

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// OtpRequest đại diện cho yêu cầu gửi OTP
// @swagger:model
type OtpRequest struct {
	Phone    string `json:"phone" binding:"required" validate:"required,e164" example:"0901230456"`
	Fullname string `json:"fullname,omitempty"`
}

// OtpVerifyRequest mở rộng từ OtpRequest để xác minh OTP
// @swagger:model
type OtpVerifyRequest struct {
	OtpRequest
	Otp      string `json:"otp" binding:"required" validate:"required"`
	AuthID   int64  `json:"authId,omitempty"`
	Version  string `json:"version,omitempty"`
	Platform string `json:"platform,omitempty"`
	OS       string `json:"os,omitempty"`
	Fullname string `json:"fullname,omitempty"`
}

// RefreshTokenRequest đại diện cho yêu cầu làm mới token
// @swagger:model
type RefreshTokenRequest struct {
	AccessToken  string `json:"accessToken" binding:"required" validate:"required"`
	RefreshToken string `json:"refreshToken" binding:"required" validate:"required"`
}

// RestorePhoneRequest đại diện cho yêu cầu khôi phục số điện thoại
// @swagger:model
type RestorePhoneRequest struct {
	Phone string `json:"phone" binding:"required" validate:"required,e164"`
}

type SessionQRRequest struct {
	ClientID   string `json:"clientId"`   // ID của trình duyệt hoặc phiên web
	UserAgent  string `json:"userAgent"`  // Chuỗi nhận diện trình duyệt
	CsrfToken  string `json:"csrfToken"`  // Token chống CSRF (nếu cần)
	Expiration int    `json:"expiration"` // Thời gian hết hạn QR code (tính bằng giây)
	SessionKey string `json:"sessionKey"`
	SessionId  uint64 `json:"sessionId"`
}

type SessionQRResponse struct {
	SessionKey string `json:"sessionKey"`
	SessionId  uint64 `json:"sessionId"`
}
type SessionQRConfirm struct {
	SessionId uint64 `json:"sessionId"` // ID của trình duyệt hoặc phiên web
}

// UserInfoRequest đại diện cho yêu cầu cập nhật thông tin người dùng
// @swagger:model
type UserInfoRequest struct {
	Phone        string     `json:"phone,omitempty"`
	Email        string     `json:"email,omitempty"`
	TaxCode      string     `json:"taxCode,omitempty"`
	FullName     string     `json:"fullname,omitempty"`
	Address      string     `json:"address,omitempty"`
	Avatar       string     `json:"avatar,omitempty"`
	ReferralBy   string     `json:"referralBy,omitempty"`
	Position     string     `json:"position,omitempty"`
	DepartmentID uint64     `json:"departmentId,omitempty"`
	Birth        *time.Time `json:"birth,omitempty" time_format:"2006-01-02T15:04:05"`
	StartDate    *time.Time `json:"startDate,omitempty" time_format:"2006-01-02T15:04:05"`
}

// Hàm validate để kiểm tra dữ liệu đầu vào
func ValidateStruct(s interface{}) error {
	validate := validator.New()
	return validate.Struct(s)
}

type ProfileIDsRequest struct {
	IDs []uint64 `form:"ids"`
}
