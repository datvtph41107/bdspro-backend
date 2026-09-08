package domain

import _models "common/models"

// UserGuideStepEntity đại diện cho các bước trong user guide
type UserGuideStepEntity struct {
	_models.BaseEntity
	UserGuideID uint64 `gorm:"column:user_guide_id;not null;index" json:"userGuideId"`
	StepOrder   int    `gorm:"column:step_order;not null" json:"stepOrder"` // Thứ tự bước (1, 2, 3...)
	Image       string `gorm:"type:text" json:"image"`                      // URL hình ảnh
	Content     string `gorm:"type:text;not null" json:"content"`           // Nội dung bước
}

// TableName đặt tên bảng
func (UserGuideStepEntity) TableName() string {
	return "user_guide_steps"
}
