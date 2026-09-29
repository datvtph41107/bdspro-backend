package mapper

import (
	_utils "common/utils"
	"crm/data"
	"crm/internal/domain"

	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
)

type FriendMapper struct {
}

func NewFriendMapper() *FriendMapper {
	return &FriendMapper{}
}

func (m *FriendMapper) FriendItemToPb(friend *domain.FriendItem) *sharepb.FriendItem {
	if friend == nil || friend.ID == 0 {
		return nil
	}
	return &sharepb.FriendItem{
		Id:         friend.ID,
		ReceiverId: friend.ReceiverID,
		Status:     uint32(friend.Status),
		CreatedBy:  &friend.CreatedBy,
	}
}

func (m *FriendMapper) FriendToPb(friend *domain.FriendEntity) *crmpb.Friend {
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

func (m *FriendMapper) ListFriendToPb(friends []domain.FriendEntity) []*crmpb.Friend {
	pbs := make([]*crmpb.Friend, len(friends))
	for i, friend := range friends {
		pbs[i] = m.FriendToPb(&friend)
	}
	return pbs
}

func (m *FriendMapper) FriendRequestPbToDomain(req *crmpb.FriendRequestDTO) *data.FriendDTO {
	return &data.FriendDTO{
		ReceiverID: req.ReceiverId,
	}
}

func (m *FriendMapper) FriendItemToPbs(friends []domain.FriendItem) []*sharepb.FriendItem {
	pbs := make([]*sharepb.FriendItem, len(friends))
	for i, friend := range friends {
		pbs[i] = m.FriendItemToPb(&friend)
	}
	return pbs
}