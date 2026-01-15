package domain

import _models "common/models"

type ProductNote struct {
	_models.BaseEntity
	ProductID uint64 `gorm:"not null;index"`
	AuthorID  uint64 `gorm:"not null;index"`
	Content   string `gorm:"type:text;not null"`
	IsPinned  bool   `gorm:"column:is_pinned;default:false"`
}

func (ProductNote) TableName() string {
	return "product_notes"
}
