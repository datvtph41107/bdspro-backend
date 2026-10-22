package entity

import (
	_models "common/domain/entity"
	"mdn4/modules/tqd/enums"
)

// Tag nhãn phân loại dùng chung, gắn với planning project qua bảng nối.
type Tag struct {
	_models.BaseEntity
	Name string        `gorm:"type:varchar(255);not null" json:"name"`
	Code string        `gorm:"type:varchar(100);index" json:"code,omitempty"`
	Type enums.TagType `gorm:"type:int4;not null;default:10;index" json:"type"`
}

func (Tag) TableName() string {
	return "tags"
}
