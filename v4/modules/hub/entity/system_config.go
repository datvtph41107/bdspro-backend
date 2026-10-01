package entity

import (
	"time"

	"gorm.io/gorm"
)

// ESystemConfigGroup identifies the logical grouping of a system config entry.
type ESystemConfigGroup int32

const (
	SystemConfigDesign      ESystemConfigGroup = 10
	SystemConfigFormat      ESystemConfigGroup = 20
	SystemConfigGeography   ESystemConfigGroup = 30
	SystemConfigParameter   ESystemConfigGroup = 40
	SystemConfigPerformance ESystemConfigGroup = 50
	SystemConfigContact     ESystemConfigGroup = 60
	SystemConfigGroupDefault ESystemConfigGroup = 70
	SystemConfigGroupSecret  ESystemConfigGroup = 80
	SystemConfigGroupSystem  ESystemConfigGroup = 90
)

var systemConfigGroupKeyMap = map[string]ESystemConfigGroup{
	"design": SystemConfigDesign, "format": SystemConfigFormat,
	"geography": SystemConfigGeography, "parameter": SystemConfigParameter,
	"performance": SystemConfigPerformance, "contact": SystemConfigContact,
	"default": SystemConfigGroupDefault, "secret": SystemConfigGroupSecret,
	"system": SystemConfigGroupSystem,
}

var systemConfigGroupNameMap = map[ESystemConfigGroup]string{
	SystemConfigDesign: "Thương hiệu & Giao diện", SystemConfigFormat: "Ngôn ngữ & Định dạng",
	SystemConfigGeography: "Phân vùng địa lý", SystemConfigParameter: "Tham số hiển thị",
	SystemConfigPerformance: "Tham số backend", SystemConfigContact: "Thông tin liên hệ",
	SystemConfigGroupDefault: "Cấu hình mặc định", SystemConfigGroupSecret: "Cấu hình bí mật",
	SystemConfigGroupSystem: "Cấu hình hệ thống",
}

var systemConfigGroupEnumToKey = map[ESystemConfigGroup]string{
	SystemConfigDesign: "design", SystemConfigFormat: "format",
	SystemConfigGeography: "geography", SystemConfigParameter: "parameter",
	SystemConfigPerformance: "performance", SystemConfigContact: "contact",
	SystemConfigGroupDefault: "default", SystemConfigGroupSecret: "secret",
	SystemConfigGroupSystem: "system",
}

// ParseSystemConfigGroup returns the enum for the given string key.
func ParseSystemConfigGroup(key string) (ESystemConfigGroup, bool) {
	v, ok := systemConfigGroupKeyMap[key]
	return v, ok
}

// SystemConfigGroupName returns the display name for the group enum.
func SystemConfigGroupName(g ESystemConfigGroup) string {
	if n, ok := systemConfigGroupNameMap[g]; ok {
		return n
	}
	return "Không xác định"
}

// SystemConfigGroupKey returns the string key for the group enum.
func SystemConfigGroupKey(g ESystemConfigGroup) string {
	if k, ok := systemConfigGroupEnumToKey[g]; ok {
		return k
	}
	return ""
}

// SystemConfig is the persistent system-wide configuration entry.
type SystemConfig struct {
	ID          uint64             `gorm:"primaryKey;autoIncrement"`
	Name        string             `gorm:"size:255;not null"`
	Key         string             `gorm:"size:100;uniqueIndex:idx_system_config_key;not null"`
	Value       string             `gorm:"type:text"`
	GroupConfig ESystemConfigGroup `gorm:"column:group_key;default:10"`
	CreatedAt   *time.Time         `gorm:"autoCreateTime"`
	UpdatedAt   *time.Time         `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt     `gorm:"index"`
}

func (SystemConfig) TableName() string { return "system_config" }
