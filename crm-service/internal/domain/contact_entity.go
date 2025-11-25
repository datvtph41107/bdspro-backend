package domain

import (
	base_enum "base/enum"
	_models "common/models"
	"crm/internal/enums"
	"time"

	"github.com/lib/pq"
)

type ContactEntity struct {
	_models.BaseEntity
	// OwnerID   uint64  `gorm:"column:owner_id" json:"ownerId"` => createdBy
	FullName string     `gorm:"column:full_name" json:"fullName" binding:"required"`
	Phone    string     `gorm:"column:phone" json:"phone"`
	Avatar   string     `gorm:"column:avatar" json:"avatar"`
	Email    string     `gorm:"column:email" json:"email"`
	Zalo     string     `gorm:"column:zalo" json:"zalo"`
	Company  string     `gorm:"column:company" json:"company"`
	Address  string     `json:"address"`
	Birthday *time.Time `gorm:"column:birthday" json:"birthday"`
	Note     string     `gorm:"column:note" json:"note"`

	Visibility   enums.EVisibility `gorm:"column:visibility;default:30" json:"visibility"`
	LeadID       uint64            `gorm:"column:lead_id" json:"leadId"`
	OwnerID      uint64            `gorm:"column:owner_id" json:"ownerId"`
	ProfileID    *uint64           `gorm:"column:profile_id" json:"profileId"`
	FriendStatus *FriendEntity     `gorm:"-" json:"friendStatus" binding:"required"`

	OwnerOf     base_enum.EOwnerOf `gorm:"column:owner_of;default:30" json:"ownerOf" binding:"required"`
	ContactTags pq.Int32Array      `gorm:"type:integer[];column:contact_tags" json:"contactTags"`

	// transient
	Tags []TagEntity `gorm:"-" json:"tags"`
}

func (ContactEntity) TableName() string {
	return "tb_contact"
}

type ContactItem struct {
	// FriendStatus FriendEntity `gorm:"column:friend_status" json:"friendStatus" binding:"required"`
	ID               uint64               `gorm:"column:id" json:"id"`
	Phone            string               `gorm:"column:phone" json:"phone"`
	FullName         string               `gorm:"column:full_name" json:"fullName" binding:"required"`
	Avatar           string               `gorm:"column:avatar" json:"avatar"`
	Email            string               `gorm:"column:email" json:"email"`
	Zalo             string               `gorm:"column:zalo" json:"zalo"`
	Company          string               `gorm:"column:company" json:"company"`
	Note             string               `gorm:"column:note" json:"note"`
	OwnerID          uint64               `gorm:"column:owner_id" json:"ownerId"`
	OwnerOf          enums.EOwnerOf       `gorm:"column:owner_of;default:10" json:"ownerOf" binding:"required"`
	ProfileID        *uint64              `gorm:"column:profile_id" json:"profileId"`
	Following        bool                 `gorm:"column:following" json:"following"`
	FriendID         uint64               `gorm:"column:friend_id" json:"friendId"`
	FriendStatus     uint64               `gorm:"column:friend_status" json:"friendStatus"`
	FriendCreatedBy  uint64               `gorm:"column:friend_created_by" json:"friendCreatedBy"`
	FriendReceiverID uint64               `gorm:"column:friend_receiver_id" json:"friendReceiverId"`
	Blocked          bool                 `gorm:"column:blocked" json:"blocked"`
	Status           enums.EStatusContact `gorm:"column:status" json:"status"`
}
