package domain

import (
	_models "common/models"
	"crm/internal/enums"
)

type ContactEntity struct {
	_models.BaseEntity
	// OwnerID   uint64  `gorm:"column:owner_id" json:"ownerId"` => createdBy
	Phone        string           `gorm:"column:phone" json:"phone"`
	FullName     string           `gorm:"column:full_name" json:"fullName" binding:"required"`
	Avatar       string           `gorm:"column:avatar" json:"avatar"`
	OwnerID      uint64           `gorm:"column:owner_id" json:"ownerId"`
	OwnerType    enums.EOwnerType `gorm:"column:owner_type;default:10" json:"ownerType" binding:"required"`
	ProfileID    uint64           `gorm:"column:profile_id" json:"profileId"`
	FriendStatus FriendEntity     `gorm:"-" json:"friendStatus" binding:"required"`
}

func (ContactEntity) TableName() string {
	return "db_contact"
}

type ContactItem struct {
	// FriendStatus FriendEntity `gorm:"column:friend_status" json:"friendStatus" binding:"required"`
	ID               uint64           `gorm:"column:id" json:"id"`
	Phone            string           `gorm:"column:phone" json:"phone"`
	FullName         string           `gorm:"column:full_name" json:"fullName" binding:"required"`
	Avatar           string           `gorm:"column:avatar" json:"avatar"`
	OwnerID          uint64           `gorm:"column:owner_id" json:"ownerId"`
	OwnerType        enums.EOwnerType `gorm:"column:owner_type;default:10" json:"ownerType" binding:"required"`
	ProfileID        uint64           `gorm:"column:profile_id" json:"profileId"`
	Following        bool             `gorm:"column:following" json:"following"`
	FriendID         uint64           `gorm:"column:friend_id" json:"friendId"`
	FriendStatus     uint64           `gorm:"column:friend_status" json:"friendStatus"`
	FriendCreatedBy  uint64           `gorm:"column:friend_created_by" json:"friendCreatedBy"`
	FriendReceiverID uint64           `gorm:"column:friend_receiver_id" json:"friendReceiverId"`
	Blocked          bool             `gorm:"column:blocked" json:"blocked"`
}
