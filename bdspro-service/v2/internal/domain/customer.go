package domain

import (
	"time"

	"gorm.io/gorm"
)

// Customer đại diện cho bảng khách hàng

type Customer struct {
	ID              uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Code            string         `gorm:"size:15" json:"code"`
	AccountID       uint64         `json:"account_id"`
	OrganizationID  uint64         `json:"organization_id"`
	Name            string         `gorm:"size:255" json:"name"`
	PhoneNumber     string         `gorm:"size:15" json:"phone_number"`
	ZaloNumber      string         `gorm:"size:15" json:"zalo_number"`
	AvatarURL       string         `gorm:"size:255" json:"avatar_url"`
	Note            string         `json:"note"`
	AssignedStaffID string         `json:"assigned_staff_id"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
