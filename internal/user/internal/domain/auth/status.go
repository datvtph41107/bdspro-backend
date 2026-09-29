package auth

import "time"

// UserStatusEntity đại diện cho trạng thái của người dùng
type UserStatusEntity struct {
	IdDomain
	// AuthID      uint64    `gorm:"primaryKey;column:auth_id"`
	Active      bool      `gorm:"column:active"`
	Verified    bool      `gorm:"column:verified"`
	LockedUntil time.Time `gorm:"column:locked_until"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	// DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName đặt tên bảng trong DB
func (UserStatusEntity) TableName() string {
	return "user_status"
}
