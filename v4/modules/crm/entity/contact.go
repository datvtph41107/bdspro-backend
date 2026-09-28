package entity

import (
	"time"

	base_enum "base/enum"
)

// Contact is a CRM-owned contact row (table: tb_contact on database.crm.dsn).
type Contact struct {
	ID              uint64             `gorm:"column:id;primaryKey"`
	FullName        string             `gorm:"column:full_name"`
	Phone           string             `gorm:"column:phone"`
	Avatar          string             `gorm:"column:avatar"`
	Email           string             `gorm:"column:email"`
	Zalo            string             `gorm:"column:zalo"`
	Company         string             `gorm:"column:company"`
	Note            string             `gorm:"column:note"`
	Birthday        *time.Time         `gorm:"column:birthday"`
	OwnerID         uint64             `gorm:"column:owner_id"`
	OwnerOf         base_enum.EOwnerOf `gorm:"column:owner_of"`
	ProfileID       *uint64            `gorm:"column:profile_id"`
	FriendID        *uint64            `gorm:"column:friend_id"`
	OriginProfileID *uint64            `gorm:"column:origin_profile_id"`
	CreatedBy       *uint64            `gorm:"column:created_by"`
	CreatedAt       *time.Time         `gorm:"column:created_at"`
	UpdatedAt       *time.Time         `gorm:"column:updated_at"`
}

func (Contact) TableName() string { return "tb_contact" }

// ContactListItem is a contact row enriched with friend/follow join fields for list.
type ContactListItem struct {
	ID               uint64             `gorm:"column:id"`
	FullName         string             `gorm:"column:full_name"`
	Phone            string             `gorm:"column:phone"`
	Avatar           string             `gorm:"column:avatar"`
	Email            string             `gorm:"column:email"`
	Zalo             string             `gorm:"column:zalo"`
	Company          string             `gorm:"column:company"`
	Note             string             `gorm:"column:note"`
	OwnerID          uint64             `gorm:"column:owner_id"`
	OwnerOf          base_enum.EOwnerOf `gorm:"column:owner_of"`
	ProfileID        *uint64            `gorm:"column:profile_id"`
	CreatedBy        *uint64            `gorm:"column:created_by"`
	CreatedAt        *time.Time         `gorm:"column:created_at"`
	UpdatedAt        *time.Time         `gorm:"column:updated_at"`
	Birthday         *time.Time         `gorm:"column:birthday"`
	Following        bool               `gorm:"column:following"`
	FriendID         uint64             `gorm:"column:friend_id"`
	FriendStatus     uint64             `gorm:"column:friend_status"`
	FriendCreatedBy  uint64             `gorm:"column:friend_created_by"`
	FriendReceiverID uint64             `gorm:"column:friend_receiver_id"`
}

// FriendRelation is a friendship row between two profiles.
type FriendRelation struct {
	ID         uint64 `gorm:"column:id"`
	ReceiverID uint64 `gorm:"column:receiver_id"`
	CreatedBy  uint64 `gorm:"column:created_by"`
	Status     uint32 `gorm:"column:status"`
}

func (FriendRelation) TableName() string { return "friend" }
