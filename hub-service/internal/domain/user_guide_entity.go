package domain

import (
	_models "common/domain/entity"
)

type UserGuideEntity struct {
	_models.BaseEntity
	Title       string `gorm:"column:title;not null" json:"title"`
	Description string `gorm:"column:description;type:text;not null" json:"description"`
	GroupKey    string `gorm:"column:group_key;not null;index" json:"groupKey"`
	Key         string `gorm:"column:key;uniqueIndex" json:"key"`
	Mode        int32  `gorm:"column:mode;not null;default:10" json:"mode"` // 10: chế độ đơn, 20: các bước
}

func (UserGuideEntity) TableName() string {
	return "tb_user_guide"
}
