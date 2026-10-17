package entity

import (
	"time"

	"userv4/modules/bdspro/enums"
)

// DealInternalNote maps to deal_internal_notes.
type DealInternalNote struct {
	ID         uint64               `gorm:"primaryKey;autoIncrement"`
	DealID     uint64               `gorm:"column:deal_id;not null;index"`
	ActorID    uint64               `gorm:"column:actor_id;not null"`
	Content    string               `gorm:"column:content;type:text;not null"`
	ActionType enums.DealActionType `gorm:"column:action_type;not null;index"`
	CreatedAt  time.Time            `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time            `gorm:"column:updated_at;autoUpdateTime"`
}

func (DealInternalNote) TableName() string { return "deal_internal_notes" }
