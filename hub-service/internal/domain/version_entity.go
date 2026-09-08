package domain

import _models "common/models"

// VersionEntity đại diện cho bảng versions trong database
type VersionEntity struct {
	_models.BaseEntity
	AppName      string `gorm:"size:100;not null" json:"appName"`
	Platform     string `gorm:"size:50;not null" json:"platform"`
	VersionName  string `gorm:"size:50;not null" json:"versionName"`
	BuildNumber  uint64 `gorm:"default:0" json:"buildNumber"`
	ForceUpdate  bool   `gorm:"default:false" json:"forceUpdate"`
	Active       bool   `gorm:"default:true" json:"active"`
	DownloadURL  string `gorm:"size:255" json:"downloadUrl"`
	ReleaseNotes string `gorm:"type:text" json:"releaseNotes"`
	BundleName   string `gorm:"size:255" json:"bundleName"`
	BundleSize   uint64 `gorm:"default:0" json:"bundleSize"`
}

// TableName đặt tên bảng cho GORM
func (VersionEntity) TableName() string {
	return "versions"
}
