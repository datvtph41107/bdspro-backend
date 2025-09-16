package dto

import (
	base_enum "base/enum"
	_dto "common/dto"
	"crm/internal/domain"
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

type RelationShipDTO struct {
	FriendStatus *domain.FriendItem
	Following    bool
	Blocked      bool
	NumFollower  int64
	NumFollowing int64
	NumFriend    int64
}
