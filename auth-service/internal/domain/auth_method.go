package domain

import (
	"time"
)

// AuthMethodDomain tương đương với Java Entity
type AuthMethodDomain struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement"`
	Provider    string     `gorm:"size:10"` // enum: email, zalo, facebook, phone
	Username    string     `gorm:"size:50"`
	OAuthID     string     `gorm:"column:oauth_id"`
	Avatar      string     `gorm:"size:255"`
	IsSensitive bool       `gorm:"column:is_sensitive"`
	Name        string     `gorm:"size:255"`
	Email       string     `gorm:"size:255"`
	Phone       string     `gorm:"size:255"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
	UserID      uint64     `gorm:"column:user_id"` // Không cập nhật & xóa, chỉ tạo mới
}

// TableName đặt tên bảng tương đương với Java
func (AuthMethodDomain) TableName() string {
	return "user_auth"
}
