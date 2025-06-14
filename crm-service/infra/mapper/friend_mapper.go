package mapper

import (
	_utils "common/utils"
	"crm/data"
	"crm/internal/domain"

	crmpb "pb/types/crm"
)

func FriendToPb(friend *domain.FriendEntity) *crmpb.Friend {
	return &crmpb.Friend{
		Id:              friend.ID,
		CreatedBy:       *friend.CreatedBy,
		ReceiverId:      friend.ReceiverID,
		RespondedAt:     _utils.FormatTimeToString(friend.RespondedAt),
		Status:          int32(friend.Status),
		GroupId:         friend.GroupID,
		GroupReceiverId: friend.GroupReceiverID,
		CreatedAt:       _utils.FormatTimeToString(friend.CreatedAt),
	}
}

func ListFriendToPb(friends []domain.FriendEntity) []*crmpb.Friend {
	pbs := make([]*crmpb.Friend, len(friends))
	for i, friend := range friends {
		pbs[i] = FriendToPb(&friend)
	}
	return pbs
}

func FriendRequestPbToDomain(req *crmpb.FriendRequestDTO) *data.FriendDTO {
	return &data.FriendDTO{
		ReceiverID: req.ReceiverId,
	}
}
