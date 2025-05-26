package auth_model

import (
	_db "common/db"
	"time"
)

// UserStatusEntity đại diện cho trạng thái của người dùng
type UserStatusEntity struct {
	UserIdEntity
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

// GetById lấy UserStatusEntity theo authId
func StatusGetById(authId uint64) (*UserStatusEntity, error) {
	var userStatus UserStatusEntity
	err := _db.DB.Where("auth_id = ?", authId).First(&userStatus).Error
	if err != nil {
		return nil, err
	}
	return &userStatus, nil
}
