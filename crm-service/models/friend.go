package models

import (
	_models "common/models"
	"crm/enums"
	"time"

	"gorm.io/gorm"
)

type FriendEntity struct {
	gorm.Model
	_models.AuditBase
	ID              uint64             `gorm:"primaryKey;autoIncrement"`
	ReceiverID      uint64             `gorm:"column:receiver_id;not null"`
	Status          enums.FriendStatus `gorm:"default:1"` // 1: Pending, 2: Accepted, 3: Declined
	RespondedAt     *time.Time         `gorm:"column:responded_at"`
	GroupID         *uint64            `gorm:"column:group_id"`
	GroupReceiverID *uint64            `gorm:"column:group_receiver_id"`
	// CreatedBy       uint64             `gorm:"column:created_by"`
	// DeletedAt       *time.Time `gorm:"column:deleted_at"`

	GroupUser     GroupEntity `gorm:"foreignKey:GroupID;references:ID"`
	GroupReceiver GroupEntity `gorm:"foreignKey:GroupReceiverID;references:ID"`
}

func (FriendEntity) TableName() string {
	return "friend"
}
