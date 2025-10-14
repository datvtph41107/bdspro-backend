package models

import (
	_models "common/domain/entity"
	"time"
)

type Membership struct {
	_models.BaseEntity
	ProfileID    uint64     `gorm:"column:profile_id;unique;not null;index" json:"profileId"`
	PlanID       uint64     `gorm:"column:plan_id;not null" json:"planId"`
	FinishedDate *time.Time `gorm:"column:finished_date" json:"finishedDate,omitempty"`
	OrderID      uint64     `gorm:"column:order_id" json:"orderId"`
	PaidAt       *time.Time `gorm:"column:paid_at" json:"paidAt,omitempty"`
}

// TableName overrides the table name used by GORM
func (Membership) TableName() string {
	return "memberships"
}
