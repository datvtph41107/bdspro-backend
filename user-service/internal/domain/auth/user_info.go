package auth

import (
	_entity "common/domain/entity"
	"time"
	"user/internal/domain/access"

	"gorm.io/gorm"
)

// Account Status enum
type AccountStatus uint8

const (
	StatusActive            AccountStatus = 10 // Tài khoản hoạt động bình thường
	StatusInactive          AccountStatus = 40 // Tài khoản bị vô hiệu hóa
	StatusTemporarilyLocked AccountStatus = 20 // Tài khoản bị khóa tạm thời
	StatusPermanentlyLocked AccountStatus = 30 // Tài khoản bị khóa vĩnh viễn
)

func (s AccountStatus) Str() string {
	return map[AccountStatus]string{
		StatusActive:            "Đang hoạt động",
		StatusInactive:          "Chưa kích hoạt",
		StatusTemporarilyLocked: "Khóa tạm thời",
		StatusPermanentlyLocked: "Khóa vĩnh viễn",
	}[s]
}

// AuthUser chứa thông tin lock của user theo profileId
type AuthUser struct {
	_entity.BaseEntityNotId
	ProfileID   uint64        `gorm:"primaryKey;column:profile_id" json:"profileId"`
	Status      AccountStatus `gorm:"column:status;default:10" json:"status"` // 10: active, 40: inactive, 20: temporarily_locked, 30: permanently_locked
	LockedAt    *time.Time    `gorm:"column:locked_at" json:"lockedAt"`
	LockedUntil *time.Time    `gorm:"column:locked_until" json:"lockedUntil"`
	LockReason  string        `gorm:"column:lock_reason;size:500" json:"lockReason"`
	LockedBy    *uint64       `gorm:"column:locked_by" json:"lockedBy"`
	Role        *access.Role  `gorm:"foreignKey:RoleID;references:ID" json:"role"`
	RoleID      *uint64       `gorm:"column:role_id" json:"roleId"`
	// CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	// UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName đặt tên bảng
func (AuthUser) TableName() string {
	return "user_info"
}

// IsLocked kiểm tra user có bị khóa không
func (u *AuthUser) IsLocked() bool {
	return u.Status == StatusTemporarilyLocked || u.Status == StatusPermanentlyLocked
}

// IsTemporarilyLocked kiểm tra user có bị khóa tạm thời không
func (u *AuthUser) IsTemporarilyLocked() bool {
	return u.Status == StatusTemporarilyLocked
}

// IsPermanentlyLocked kiểm tra user có bị khóa vĩnh viễn không
func (u *AuthUser) IsPermanentlyLocked() bool {
	return u.Status == StatusPermanentlyLocked
}

// IsLockExpired kiểm tra khóa tạm thời có hết hạn không
func (u *AuthUser) IsLockExpired() bool {
	if !u.IsTemporarilyLocked() || u.LockedUntil == nil {
		return false
	}
	return time.Now().After(*u.LockedUntil)
}

// CanLogin kiểm tra user có thể đăng nhập không
func (u *AuthUser) CanLogin() bool {
	if u.Status == StatusInactive {
		return false
	}
	if u.IsPermanentlyLocked() {
		return false
	}
	if u.IsTemporarilyLocked() && !u.IsLockExpired() {
		return false
	}
	return true
}

// Validate kiểm tra tính hợp lệ của AuthUser
func (u *AuthUser) Validate() error {
	// ProfileID sẽ được set khi tạo mới, không cần validate ở đây
	// vì nó được set từ request
	return nil
}

// GORM hook để validate trước khi tạo
func (u *AuthUser) BeforeCreate(tx *gorm.DB) error {
	return u.Validate()
}

// GORM hook để validate trước khi cập nhật
func (u *AuthUser) BeforeUpdate(tx *gorm.DB) error {
	return u.Validate()
}
