package models

import (
	_enum "common/domain/enum"
	_model "common/models"
	"time"
)

// AdminProfile chứa thông tin tài khoản admin
type AdminProfile struct {
	_model.BaseEntity
	AuthID           uint64            `gorm:"column:auth_id;not null;index" copier:"-" json:"authId"`
	ProfileID        uint64            `gorm:"column:profile_id" json:"profileId"`
	Username         string            `gorm:"-" json:"username"`
	Password         string            `gorm:"-" json:"password"`
	Address          string            `gorm:"size:255" json:"address"`
	Gender           *uint32           `gorm:"default:0" json:"gender"`
	Birth            *time.Time        `gorm:"column:birth" json:"birth"`
	FullName         string            `gorm:"column:full_name;size:255" json:"fullName"`
	Email            string            `gorm:"size:50" json:"email"`
	Phone            string            `gorm:"column:phone;size:20" json:"phone"`
	Avatar           string            `gorm:"column:avatar;size:255" json:"avatar"`
	JobTitle         string            `gorm:"column:job_title;size:100" json:"jobTitle"`     // Chức danh / Vị trí công việc
	WorkAt           *time.Time        `gorm:"column:work_at" json:"workAt"`                  // Thời gian làm việc
	RoleID           *uint64           `gorm:"column:role_id" json:"roleId"`                  // ID vai trò từ auth service
	RoleKey          string            `gorm:"size:50" json:"roleKey"`                        // Vai trò hệ thống (role)
	RoleType         string            `gorm:"size:50" json:"roleType"`                       // Loại vai trò từ auth service
	RoleDescription  string            `gorm:"size:255" json:"roleDescription"`               // Mô tả vai trò
	Attachments      []string          `gorm:"serializer:json;type:jsonb" json:"attachments"` // File đính kèm (URLs) - JSON string
	InternalNotes    string            `gorm:"type:text" json:"internalNotes"`                // Ghi chú nội bộ
	SendNotification bool              `gorm:"default:false" json:"sendNotification"`         // Gửi thông báo tài khoản qua Email/SMS
	Status           _enum.EUserStatus `gorm:"default:10" json:"status"`                      // Trạng thái: 0: Chưa hoạt động, 10: Hoạt động, 20: Khóa tạm thời, 30: Khóa vĩnh viễn, 40: Ngừng hoạt động
	LastLoginAt      *time.Time        `gorm:"column:last_login_at" json:"lastLoginAt"`
	TotalLogin       uint32            `gorm:"default:0" json:"totalLogin"`
	VerifiedAt       *time.Time        `gorm:"column:verified_at" json:"verifiedAt"`
	LockedAt         *time.Time        `gorm:"column:locked_at" json:"lockedAt"`

	// Transient fields (không lưu vào DB)
	Role *RoleDTO `gorm:"-" json:"role,omitempty"`
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
