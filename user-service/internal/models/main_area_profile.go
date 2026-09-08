package models

import "time"

// MainAreaProfileEntity là bảng nối giữa user (profile) và nhãn khu vực hoạt động.
type MainAreaProfileEntity struct {
	ProfileID  uint64    `gorm:"column:profile_id;primaryKey" json:"profileId"`
	MainAreaID uint64    `gorm:"column:main_area_id;primaryKey" json:"mainAreaId"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (MainAreaProfileEntity) TableName() string {
	return "main_area_profile"
}
