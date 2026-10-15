package auth

import (
	"time"

	"gorm.io/gorm"
)

// UserOTPEntity đại diện cho bảng user_otp trong database
type UserOTPEntity struct {
	IdDomain
	// ID           uint           `gorm:"primaryKey" json:"id"`
	// AuthID uint `gorm:"column:auth_id" json:"authId"`
	OTP          string     `gorm:"size:10" json:"otp"`
	OTPSendTime  int        `gorm:"default:0" json:"otpSendTime"`
	OTPCheckTime int        `gorm:"default:0" json:"otpCheckTime"`
	OTPDate      *time.Time `gorm:"type:timestamptz" json:"otpDate"` // Bỏ autoUpdateTime để tránh update nhầm khi check OTP
	ExpiredTime  *time.Time `gorm:"type:timestamptz" json:"expiredTime"`
	Activate     bool       `gorm:"default:false" json:"activate"`
	// CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	// UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	// DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName đặt tên bảng tương đương với Java
func (UserOTPEntity) TableName() string {
	return "user_otp"
}

// BeforeCreate GORM hook để đặt expiredTime trước khi tạo bản ghi mới
func (u *UserOTPEntity) BeforeCreate(tx *gorm.DB) (err error) {
	// Thiết lập thời gian hết hạn mặc định nếu chưa có
	// t := time.Now().Add(time.Second * time.Duration(120)) // Thay đổi 120 giây thành giá trị phù hợp
	// u.ExpiredTime = &t
	return
}

// BeforeCreate GORM hook để đặt expiredTime trước khi tạo bản ghi mới
func (u *UserOTPEntity) BeforeUpdate(tx *gorm.DB) (err error) {
	// Thiết lập thời gian hết hạn mặc định nếu chưa có
	t := time.Now().Add(time.Second * time.Duration(120)) // Thay đổi 120 giây thành giá trị phù hợp
	u.ExpiredTime = &t
	return
}
