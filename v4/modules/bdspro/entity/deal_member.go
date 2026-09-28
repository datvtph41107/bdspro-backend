package entity

import (
	"time"

	"gorm.io/gorm"
)

// DealMember maps to deal_members.
type DealMember struct {
	ID              uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt       *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`
	DealID          uint64         `gorm:"column:deal_id;index;not null"`
	MemberID        uint64         `gorm:"column:member_id;index;not null"`
	MemberType      uint32         `gorm:"column:member_type;default:10"`
	Status          uint32         `gorm:"column:status;default:10"`
	RoleID          uint64         `gorm:"column:role_id;not null;default:0"`
	AmountCommit    float64        `gorm:"column:amount_commit;default:0"`
	CommissionValue float64        `gorm:"column:commission_value;default:0"`
	CommissionType  uint32         `gorm:"column:commission_type;default:10"`
	Note            string         `gorm:"column:note;type:text"`
	IsUnilateral    bool           `gorm:"column:is_unilateral;default:false"`
	DoneInvestment  bool           `gorm:"column:done_investment;default:false"`
	RoleKey         uint32         `gorm:"column:role_key;default:430"`
	Message         string         `gorm:"column:message;type:text"`
	InvitedAt       *time.Time     `gorm:"column:invited_at"`
	RespondedAt     *time.Time     `gorm:"column:responded_at"`
	WithdrawnAt     *time.Time     `gorm:"column:withdrawn_at"`
	InviterID       uint64         `gorm:"column:inviter_id"`
	ColorID         *uint32        `gorm:"column:color_id"`
	IsOwner         bool           `gorm:"column:is_owner;default:false"`
}

func (DealMember) TableName() string { return "deal_members" }
