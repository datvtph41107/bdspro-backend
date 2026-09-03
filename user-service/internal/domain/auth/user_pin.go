package auth

import (
	"time"

	"gorm.io/gorm"
)

// UserPINEntity đại diện cho bảng user_pin trong database
type UserPINEntity struct {
	IdDomain
	PIN          string     `gorm:"size:255" json:"-"`                   // Mã PIN đã được hash (không expose trong JSON)
	PINCheckTime int        `gorm:"default:0" json:"pinCheckTime"`       // Số lần nhập sai PIN
	PINDate      *time.Time `gorm:"type:timestamptz" json:"pinDate"`     // Thời gian tạo/cập nhật PIN
	IsActive     bool       `gorm:"default:true" json:"isActive"`        // PIN có đang hoạt động không
	LockedUntil  *time.Time `gorm:"type:timestamptz" json:"lockedUntil"` // Thời gian khóa tài khoản (nếu nhập sai nhiều lần)
}

// TableName đặt tên bảng
func (UserPINEntity) TableName() string {
	return "user_pin"
}

// BeforeCreate GORM hook
func (u *UserPINEntity) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	u.PINDate = &now
	return
}

// BeforeUpdate GORM hook
func (u *UserPINEntity) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	u.PINDate = &now
	return
}
