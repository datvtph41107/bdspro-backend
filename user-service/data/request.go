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

// UserInfoRequest đại diện cho yêu cầu cập nhật thông tin người dùng
// @swagger:model
type UserInfoRequest struct {
	Phone        string     `json:"phone,omitempty" validate:"omitempty,min=10,max=15"`
	Email        string     `json:"email,omitempty" validate:"omitempty,email,max=100"`
	TaxCode      string     `json:"taxCode,omitempty" validate:"omitempty,min=2,max=50"`
	Fullname     string     `json:"fullname,omitempty" validate:"omitempty,min=2,max=50"`
	Address      string     `json:"address,omitempty" validate:"omitempty,min=2,max=100"`
	Avatar       string     `json:"avatar,omitempty" validate:"omitempty,url,max=255"`
	Position     string     `json:"position,omitempty" validate:"omitempty,min=2,max=50"`
	DepartmentID uint64     `json:"department_id,omitempty" validate:"omitempty,gte=1"`
	Birth        *time.Time `json:"birth,omitempty" validate:"omitempty"`
}

// Hàm validate để kiểm tra dữ liệu đầu vào
func ValidateStruct(s interface{}) error {
	validate := validator.New()
	return validate.Struct(s)
}
