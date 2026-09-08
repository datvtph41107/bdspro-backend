package domain

import (
	"bdspro/internal/enums"
	"time"
)

type DealInternalNote struct {
	ID         uint64               `gorm:"primaryKey;autoIncrement"`
	DealID     uint64               `gorm:"not null;index"`
	ActorID    uint64               `gorm:"not null"` // ID người tạo ghi chú
	Content    string               `gorm:"type:text;not null"`
	ActionType enums.DealActionType `gorm:"not null;index"`
	CreatedAt  time.Time            `gorm:"autoCreateTime"`
	UpdatedAt  time.Time            `gorm:"autoUpdateTime"`

	// Relations
	Deal  *Deal       `gorm:"foreignKey:DealID;references:ID"`
	Actor interface{} `gorm:"-"` // Sẽ được populate từ User service
}

func (DealInternalNote) TableName() string {
	return "deal_internal_notes"
}
