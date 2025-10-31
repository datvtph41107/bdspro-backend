package domain

import (
	_models "common/domain/entity"
)

type UserGuideEntity struct {
	_models.BaseEntity
	Title       string `gorm:"column:title;not null" json:"title"`
	Description string `gorm:"column:description;type:text;not null" json:"description"`
	GroupKey    string `gorm:"column:group_key;not null;index" json:"groupKey"`
}

func (UserGuideEntity) TableName() string {
	return "tb_user_guide"
}
