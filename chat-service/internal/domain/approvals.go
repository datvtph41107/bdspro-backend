package domain

import (
	"chat/internal/enums"
	"time"
)

type Approval struct {
	ConversationID uint64 `gorm:"primaryKey;autoIncrement:false"`
	InviteByID     uint64
	NewMemberID    uint64
	ApproveByID    *uint64
	InviteType     enums.InviteType `gorm:"type:smallint;not null;index"`
	CreatedAt      time.Time        `gorm:"not null;autoCreateTime"`
	UpdatedAt      time.Time        `gorm:"not null;autoUpdateTime"`
}

func (Approval) TableName() string {
	return "approval_requests"
}
