package auth

import (
	"time"

	_models "common/domain/entity"
)

// DeviceEntity đại diện cho thông tin thiết bị ứng dụng của người dùng
type DeviceEntity struct {
	_models.BaseEntity

	DeviceID       string     `gorm:"column:device_id;size:255;uniqueIndex:idx_device_unique" json:"deviceId"`
	DeviceName     string     `gorm:"column:device_name;size:255" json:"deviceName"`
	DeviceType     string     `gorm:"column:device_type;size:50" json:"deviceType"`
	Platform       string     `gorm:"column:platform;size:100" json:"platform"`
	OSVersion      string     `gorm:"column:os_version;size:100" json:"osVersion"`
	AppVersion     string     `gorm:"column:app_version;size:50" json:"appVersion"`
	BuildNumber    string     `gorm:"column:build_number;size:50" json:"buildNumber"`
	Manufacturer   string     `gorm:"column:manufacturer;size:100" json:"manufacturer"`
	Model          string     `gorm:"column:model;size:100" json:"model"`
	Locale         string     `gorm:"column:locale;size:20" json:"locale"`
	Timezone       string     `gorm:"column:timezone;size:50" json:"timezone"`
	PushToken      string     `gorm:"column:push_token;size:255" json:"pushToken"`
	LastSeenAt     *time.Time `gorm:"column:last_seen_at" json:"lastSeenAt"`
	AuthID         *uint64    `gorm:"column:auth_id" json:"authId,omitempty"`
	ProfileID      *uint64    `gorm:"column:profile_id" json:"profileId,omitempty"`
	OrganizationID *uint64    `gorm:"column:organization_id" json:"organizationId,omitempty"`
	IPAddress      string     `gorm:"column:ip_address;size:64" json:"ipAddress"`
	UserAgent      string     `gorm:"column:user_agent;size:512" json:"userAgent"`
}

// TableName trả về tên bảng trong cơ sở dữ liệu
func (DeviceEntity) TableName() string {
	return "auth_devices"
}
