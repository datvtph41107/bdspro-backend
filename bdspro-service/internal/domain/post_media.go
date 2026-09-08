package domain

import _models "common/models"

type PostMedia struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	PostID    uint64 `gorm:"not null" json:"postId"`
	MediaURL  string `gorm:"type:VARCHAR(500);not null" json:"mediaUrl"`
	MediaType string `gorm:"type:VARCHAR(20);not null" json:"mediaType"`
	IsMain    bool   `gorm:"default:false" json:"isMain"`
	Order     int    `gorm:"default:0" json:"sort_order,omitempty"`
}

type PostMediaEntity struct {
	_models.BaseEntity
	// ID        uint64 `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	PostID    uint64 `gorm:"not null" json:"postId"`
	MediaURL  string `gorm:"type:VARCHAR(500);not null" json:"mediaUrl"`
	MediaType string `gorm:"type:VARCHAR(20);not null" json:"mediaType"`
	IsMain    bool   `gorm:"default:false" json:"isMain"`
	SortOrder int    `gorm:"default:0" json:"sort_order,omitempty"`
}

func (PostMediaEntity) TableName() string {
	return "post_media"
}
