package models

import (
	_models "common/models"
	"time"
	"user/internal/enums"
)

type FriendEntity struct {
	_models.BaseEntity
	ReceiverID      uint64              `gorm:"column:receiver_id;not null" json:"receiverId"`
	Status          enums.EFriendStatus `gorm:"default:1" json:"status"` // 1: Pending, 2: Accepted, 3: Declined
	RespondedAt     *time.Time          `gorm:"column:responded_at" json:"respondedAt,omitempty"`
	GroupID         *uint64             `gorm:"column:group_id" json:"groupId,omitempty"`
	GroupReceiverID *uint64             `gorm:"column:group_receiver_id" json:"groupReceiverId,omitempty"`

	GroupUser     *GroupEntity `gorm:"foreignKey:GroupID;references:ID" json:"groupUser,omitempty"`
	GroupReceiver *GroupEntity `gorm:"foreignKey:GroupReceiverID;references:ID" json:"groupReceiver,omitempty"`
}

func (FriendEntity) TableName() string {
	return "friend"
}

// type FriendItem struct {
// 	ID         uint64              `gorm:"column:id;not null" json:"id"`
// 	ReceiverID uint64              `gorm:"column:receiver_id;not null" json:"receiverId"`
// 	CreatedBy  uint64              `gorm:"column:created_by;not null" json:"createdBy"`
// 	Status     enums.EFriendStatus `gorm:"default:10" json:"status"` // 1: Pending, 2: Accepted, 3: Declined
// }

// func (FriendItem) TableName() string {
// 	return "friend"
// }
