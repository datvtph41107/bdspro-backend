package entity

import (
	"mdn4/modules/tqd/enums"
	"time"

	"gorm.io/datatypes"
)

type UserSubscription struct {
	ID                uint64                   `gorm:"primaryKey" json:"id"`
	UserID            uint64                   `gorm:"column:user_id;not null;index:idx_sub_user_target,priority:1" json:"userId"`
	TargetType        string                   `gorm:"column:target_type;type:varchar(50);not null;index:idx_sub_user_target,priority:2;index:idx_sub_target" json:"targetType"`
	TargetID          uint64                   `gorm:"column:target_id;not null;index:idx_sub_user_target,priority:3;index:idx_sub_target" json:"targetId"`
	SubscriptionScope datatypes.JSON           `gorm:"column:subscription_scope;type:jsonb;not null" json:"subscriptionScope"`
	TriggerConfig     datatypes.JSON           `gorm:"column:trigger_config;type:jsonb" json:"triggerConfig"`
	Status            enums.SubscriptionStatus `gorm:"column:status;index:idx_sub_user_status" json:"status"`
	LastTriggeredAt   *time.Time               `gorm:"column:last_triggered_at" json:"lastTriggeredAt"`
	CreatedAt         time.Time                `gorm:"column:created_at;default:now()" json:"createdAt"`
	UpdatedAt         time.Time                `gorm:"column:updated_at;default:now()" json:"updatedAt"`
	DeletedAt         *time.Time               `gorm:"column:deleted_at;index" json:"deletedAt"`
}

func (UserSubscription) TableName() string { return "user_subscriptions" }
