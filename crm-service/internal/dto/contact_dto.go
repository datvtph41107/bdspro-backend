package dto

import (
	base_enum "base/enum"
	_dto "common/domain/dto"
	"crm/internal/enums"
	"time"
)

type ContactSearchDTO struct {
	_dto.Pagable
	Text         string  `form:"text"`
	FriendStatus []int32 `form:"friendStatus"`
	Following    *bool   `form:"following"`
	Roles        []int32 `form:"roles"`
	Leaded       *bool   `form:"leaded"`
	Connected    *bool   `form:"connected"`
	AppInstalled *bool   `form:"appInstalled"`

	OwnerId uint64             `form:"ownerId"`
	OwnerOf base_enum.EOwnerOf `form:"ownerOf"`
}

type PinRequestDTO struct {
	Target    enums.PinTarget
	TargetIDs map[string]uint64
	NameSpace string
	Action    string
}

type ContactProductInterestedDTO struct {
	ContactID   *uint64   `gorm:"column:contact_id"`
	OriginID    *uint64   `gorm:"column:origin_profile_id"`
	DisplayName string    `gorm:"column:display_name"`
	Phone       string    `gorm:"column:phone"`
	Avatar      string    `gorm:"column:avatar"`
	Status      uint32    `gorm:"column:status"`
	Pinned      bool      `gorm:"column:pinned"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	HasContact  bool
}

type ContactProductInterestedFilter struct {
	_dto.Pagable
	// OwnerID   uint64
	// OwnerOf   base_enum.EOwnerOf
	ProductID uint64
}
