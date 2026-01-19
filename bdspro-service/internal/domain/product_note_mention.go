package domain

import (
	_models "common/domain/entity"
)

type ProductNoteMention struct {
	_models.BaseEntity
	ID       uint64 `gorm:"primaryKey"`
	NoteID   uint64 `gorm:"not null"`
	StartPos int32  `gorm:"column:start_pos;"`
	Length   int32  `gorm:"column:length;"`
	UserID   uint64 `gorm:"not null"`
}

func (ProductNoteMention) TableName() string {
	return "product_note_mentions"
}
