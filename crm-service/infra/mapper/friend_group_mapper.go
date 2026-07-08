package mapper

import (
	"crm/internal/domain"

	crmpb "pb/types/crm"
)

func FriendGroupToPb(group *domain.GroupEntity) *crmpb.FriendGroup {
	return &crmpb.FriendGroup{
		Id:   group.ID,
		Name: group.Name,
	}
}

func PbToFriendGroup(pb *crmpb.FriendGroup) *domain.GroupEntity {
	return &domain.GroupEntity{
		ID:   pb.Id,
		Name: pb.Name,
	}
}

func ListFriendGroupToPb(groups []domain.GroupEntity) []*crmpb.FriendGroup {
	pbs := make([]*crmpb.FriendGroup, len(groups))
	for i, group := range groups {
		pbs[i] = FriendGroupToPb(&group)
	}
	return pbs
}