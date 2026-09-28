package entity

import (
	"time"

	"gorm.io/gorm"
)

type Version struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt    *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedBy    *uint64        `gorm:"column:created_by"`
	UpdatedBy    *uint64        `gorm:"column:updated_by"`
	AppName      string         `gorm:"column:app_name;size:100;not null"`
	Platform     string         `gorm:"column:platform;size:50;not null"`
	VersionName  string         `gorm:"column:version_name;size:50;not null"`
	BuildNumber  uint64         `gorm:"column:build_number;default:0"`
	ForceUpdate  bool           `gorm:"column:force_update;default:false"`
	Active       bool           `gorm:"column:active;default:true"`
	DownloadURL  string         `gorm:"column:download_url;size:255"`
	ReleaseNotes string         `gorm:"column:release_notes;type:text"`
	BundleName   string         `gorm:"column:bundle_name;size:255"`
	BundleSize   uint64         `gorm:"column:bundle_size;default:0"`
}

func (Version) TableName() string { return "versions" }

type APIKey struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt   *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
	Name        string         `gorm:"column:name;size:100;index;not null"`
	AppName     string         `gorm:"column:app_name;size:100;index;not null"`
	Description string         `gorm:"column:description;size:255"`
	APIKey      string         `gorm:"column:api_key;size:255;index;not null"`
	ExpiredAt   *time.Time     `gorm:"column:expired_at"`
}

func (APIKey) TableName() string { return "api_keys" }
