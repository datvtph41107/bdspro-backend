package models

import (
	_model "common/models"
	"time"
)

// AdminProfile chứa thông tin tài khoản admin
type AdminProfile struct {
	_model.BaseEntity
	AuthID           uint64     `gorm:"primaryKey;autoIncrement;column:auth_id" copier:"-" json:"authId"`
	ProfileID        uint64     `gorm:"column:profile_id" json:"profileId"`
	Username         string     `gorm:"size:50;uniqueIndex" json:"username"`
	Password         string     `gorm:"size:255" json:"password"`
	FullName         string     `gorm:"size:255" json:"fullName"`
	Email            string     `gorm:"size:50;uniqueIndex" json:"email"`
	Phone            string     `gorm:"size:20" json:"phone"`
	Avatar           string     `gorm:"size:255" json:"avatar"`
	JobTitle         string     `gorm:"size:100" json:"jobTitle"`              // Chức danh / Vị trí công việc
	RoleKey          string     `gorm:"size:50" json:"roleKey"`                // Vai trò hệ thống (role)
	Attachments      []string   `gorm:"type:text" json:"attachments"`          // File đính kèm (URLs) - JSON string
	InternalNotes    string     `gorm:"type:text" json:"internalNotes"`        // Ghi chú nội bộ
	SendNotification bool       `gorm:"default:false" json:"sendNotification"` // Gửi thông báo tài khoản qua Email/SMS
	Status           uint8      `gorm:"default:1" json:"status"`               // 1: active, 0: inactive
	LastLoginAt      *time.Time `gorm:"column:last_login_at" json:"lastLoginAt"`
	TotalLogin       uint32     `gorm:"default:0" json:"totalLogin"`
	VerifiedAt       *time.Time `gorm:"column:verified_at" json:"verifiedAt"`
	LockedAt         *time.Time `gorm:"column:locked_at" json:"lockedAt"`
}

// TableName trả về tên bảng trong database
func (AdminProfile) TableName() string {
	return "admin_profiles"
}

// AdminProfileResponse DTO cho response tạo admin
type AdminProfileResponse struct {
	AuthID    uint64 `json:"authId"`
	ProfileID uint64 `json:"profileId"`
	Role      string `json:"role"`
}
