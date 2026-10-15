package access

import (
	"time"
)

// AdminAccessDomain đại diện cho quyền truy cập hệ thống nội bộ
type AdminAccessDomain struct {
	ID            uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64     `gorm:"column:user_id;not null" json:"userId"`
	IPAddress     string     `gorm:"column:ip_address;size:45" json:"ipAddress"` // IPv4/IPv6
	IPRange       string     `gorm:"column:ip_range;size:50" json:"ipRange"`     // CIDR notation
	DeviceID      string     `gorm:"column:device_id;size:255" json:"deviceId"`
	DeviceName    string     `gorm:"column:device_name;size:255" json:"deviceName"`
	DeviceType    string     `gorm:"column:device_type;size:50" json:"deviceType"`         // mobile, desktop, tablet
	AuthType      string     `gorm:"column:auth_type;size:20" json:"authType"`             // IP, DEVICE, OTP, 2FA
	Status        string     `gorm:"column:status;size:20;default:'active'" json:"status"` // active, blocked, expired
	EffectiveFrom *time.Time `gorm:"column:effective_from" json:"effectiveFrom"`
	EffectiveTo   *time.Time `gorm:"column:effective_to" json:"effectiveTo"`
	MaxDevices    int        `gorm:"column:max_devices;default:3" json:"maxDevices"`
	Require2FA    bool       `gorm:"column:require_2fa;default:true" json:"require2fa"`
	CreatedBy     uint64     `gorm:"column:created_by" json:"createdBy"`
	UpdatedBy     uint64     `gorm:"column:updated_by" json:"updatedBy"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt     *time.Time `gorm:"column:deleted_at" json:"deletedAt"`
}

// TableName đặt tên bảng
func (AdminAccessDomain) TableName() string {
	return "admin_access_control"
}

// AdminAccessLogDomain đại diện cho log truy cập admin
type AdminAccessLogDomain struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint64    `gorm:"column:user_id;not null" json:"userId"`
	IPAddress  string    `gorm:"column:ip_address;size:45" json:"ipAddress"`
	DeviceID   string    `gorm:"column:device_id;size:255" json:"deviceId"`
	DeviceName string    `gorm:"column:device_name;size:255" json:"deviceName"`
	Action     string    `gorm:"column:action;size:50" json:"action"` // login, logout, access_denied
	Status     string    `gorm:"column:status;size:20" json:"status"` // success, failed, blocked
	Reason     string    `gorm:"column:reason;size:500" json:"reason"`
	UserAgent  string    `gorm:"column:user_agent;size:500" json:"userAgent"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

// TableName đặt tên bảng
func (AdminAccessLogDomain) TableName() string {
	return "admin_access_logs"
}

// AdminAccessValidationResult kết quả validation truy cập
type AdminAccessValidationResult struct {
	IsAllowed  bool   `json:"isAllowed"`
	Reason     string `json:"reason"`
	Require2FA bool   `json:"require2fa"`
	AccessID   uint64 `json:"accessId"`
	DeviceID   string `json:"deviceId"`
	IPAddress  string `json:"ipAddress"`
}
