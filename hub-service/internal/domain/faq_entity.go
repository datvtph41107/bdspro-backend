package domain

import _models "common/models"

// FAQEntity đại diện cho bảng faqs trong database
type FAQEntity struct {
	_models.BaseEntity
	Question string `gorm:"type:text;not null" json:"question"`
	Answer   string `gorm:"type:text;not null" json:"answer"`
	GroupKey string `gorm:"size:100;index" json:"groupKey"`
}

// TableName đặt tên bảng
func (FAQEntity) TableName() string {
	return "faqs"
}
