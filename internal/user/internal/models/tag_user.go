package models

import "time"

// TagUserEntity là bảng nối giữa user (profile) và tag.
type TagUserEntity struct {
	ProfileID uint64    `gorm:"column:profile_id;primaryKey" json:"profileId"`
	TagID     uint64    `gorm:"column:tag_id;primaryKey" json:"tagId"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (TagUserEntity) TableName() string {
	return "tag_user"
}
