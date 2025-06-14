package dto

import (
	_dto "common/dto"
	"crm/internal/domain"
)

type ContactSearchDTO struct {
	_dto.Pagable
	Text string `form:"text"`
}

type RelationShipDTO struct {
	FriendStatus *domain.FriendItem
	Following    bool
	Blocked      bool
	NumFollower  int64
	NumFollowing int64
	NumFriend    int64
}
