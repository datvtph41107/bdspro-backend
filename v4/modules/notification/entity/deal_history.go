package entity

import (
	"time"

	"gorm.io/datatypes"
)

// DealHistory is deal timeline entries.
type DealHistory struct {
	ID          uint64            `gorm:"primaryKey;autoIncrement"`
	DealID      uint64            `gorm:"column:deal_id;index;not null"`
	ActorID     uint64            `gorm:"column:actor_id;index;not null"`
	ActorName   string            `gorm:"column:actor_name;size:255"`
	ActorAvatar string            `gorm:"column:actor_avatar;size:500"`
	ActionType  string            `gorm:"column:action_type;size:50;index;not null"`
	ActionName  string            `gorm:"column:action_name;size:255;not null"`
	Content     string            `gorm:"column:content;type:text;not null"`
	Metadata    datatypes.JSONMap `gorm:"column:metadata;type:jsonb;default:'{}'"`
	CreatedAt   time.Time         `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time         `gorm:"column:updated_at;autoUpdateTime"`
}

func (DealHistory) TableName() string { return "deal_history" }
