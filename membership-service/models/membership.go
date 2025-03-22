package models

import (
	_models "common/models"
	"time"
)

type Membership struct {
	_models.BaseEntity
	ProfileId    uint64     `gorm:"unique" json:"profileId"`
	PlanId       uint64     `json:"planId"`
	FinishedDate *time.Time `json:"finishedDate,omitempty"`
	OrderId      uint64     `json:"orderId"`
	PaidAt       *time.Time `json:"paidAt,omitempty"`
}

// TableName overrides the table name used by GORM
func (Membership) TableName() string {
	return "memberships"
}
