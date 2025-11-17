package models

import (
	_models "common/models"
	"user/enums"
)

// TagEntity biểu diễn thẻ tùy chỉnh của người dùng.
type TagEntity struct {
	_models.BaseEntity
	Name      string        `gorm:"column:name;type:varchar(255);not null" json:"name"`
	IsDefault bool          `gorm:"column:is_default;default:false" json:"isDefault"`
	IsActive  bool          `gorm:"column:is_active;type:boolean;not null" json:"isActive"`
	UserID    *uint64       `gorm:"column:user_id" json:"userId"`
	TagType   enums.TagType `gorm:"column:tag_type;not null" json:"tagType"`
}

func (TagEntity) TableName() string {
	return "tags"
}
