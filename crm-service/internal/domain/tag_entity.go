package domain

import (
	_models "common/models"
)

// TagEntity biểu diễn thẻ tùy chỉnh cho contact.
type TagEntity struct {
	_models.BaseEntity
	Name      string  `gorm:"column:name;type:varchar(255);not null" json:"name"`
	IsDefault bool    `gorm:"column:is_default;default:false" json:"isDefault"`
	IsActive  bool    `gorm:"column:is_active;type:boolean;not null;default:true" json:"isActive"`
	OwnerID   *uint64 `gorm:"column:owner_id" json:"ownerId"`
	OwnerOf   *int32  `gorm:"column:owner_of" json:"ownerOf"`
}

func (TagEntity) TableName() string {
	return "tb_contact_tag"
}