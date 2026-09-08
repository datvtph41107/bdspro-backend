package domain

import (
	"time"
	"tqd/internal/enums"

	"gorm.io/datatypes"
)

type UserSubscription struct {
	ID                uint64                   `gorm:"primaryKey"`
	UserID            uint64                   `gorm:"column:user_id;not null;index:idx_sub_user_target,priority:1"`
	TargetType        string                   `gorm:"column:target_type;type:varchar(50);not null;index:idx_sub_user_target,priority:2;index:idx_sub_target"`
	TargetID          uint64                   `gorm:"column:target_id;not null;index:idx_sub_user_target,priority:3;index:idx_sub_target"`
	SubscriptionScope datatypes.JSON           `gorm:"column:subscription_scope;type:jsonb;not null"`
	TriggerConfig     datatypes.JSON           `gorm:"column:trigger_config;type:jsonb"`
	Status            enums.SubscriptionStatus `gorm:"column:status;index:idx_sub_user_status"`
	LastTriggeredAt   *time.Time               `gorm:"column:last_triggered_at"`
	CreatedAt         time.Time                `gorm:"column:created_at;default:now()"`
	UpdatedAt         time.Time                `gorm:"column:updated_at;default:now()"`
	DeletedAt         *time.Time               `gorm:"column:deleted_at;index"`
}

func (UserSubscription) TableName() string { return "user_subscriptions" }
