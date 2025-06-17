package transformer

import (
	grouppb "pb/types/organization"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type GroupTransformer interface {
	CreateGroupRequestToEntity(request *grouppb.CreateGroupRequest) *entity.Group
	UpdateGroupRequestToEntity(request *grouppb.UpdateGroupRequest) *entity.Group
	EntityToCreateGroupResponse(group *entity.Group) *grouppb.CreateGroupResponse
	EntityToUpdateGroupResponse(group *entity.Group) *grouppb.UpdateGroupResponse
	EntityToGetGroupResponse(group *entity.Group) *grouppb.GetGroupResponse
	EntityToGetGroupsResponse(groups []*entity.Group, total uint32) *grouppb.GetGroupsResponse
}

type groupTransformer struct{}

func NewGroupTransformer() GroupTransformer {
	return &groupTransformer{}
}

func (t *groupTransformer) CreateGroupRequestToEntity(request *grouppb.CreateGroupRequest) *entity.Group {
	return &entity.Group{
		Name:        request.Name,
		Description: request.Description,
		AvatarUrl:   request.AvatarUrl,
	}
}

func (t *groupTransformer) UpdateGroupRequestToEntity(request *grouppb.UpdateGroupRequest) *entity.Group {
	return &entity.Group{
		Id:          request.Id,
		Name:        request.Name,
		Description: request.Description,
	}
}

func (t *groupTransformer) EntityToCreateGroupResponse(group *entity.Group) *grouppb.CreateGroupResponse {
	return &grouppb.CreateGroupResponse{
		Id: group.Id,
	}
}

func (t *groupTransformer) EntityToUpdateGroupResponse(group *entity.Group) *grouppb.UpdateGroupResponse {
	return &grouppb.UpdateGroupResponse{
		Id: group.Id,
	}
}

func (t *groupTransformer) EntityToGetGroupResponse(group *entity.Group) *grouppb.GetGroupResponse {
	return &grouppb.GetGroupResponse{
		Group: &grouppb.Group{
			Id:          group.Id,
			Name:        group.Name,
			Description: group.Description,
			AvatarUrl:   group.AvatarUrl,
		},
	}
}

func (t *groupTransformer) EntityToGetGroupsResponse(groups []*entity.Group, total uint32) *grouppb.GetGroupsResponse {
	groupsResponse := make([]*grouppb.Group, len(groups))
	for i, group := range groups {
		groupsResponse[i] = &grouppb.Group{
			Id:          group.Id,
			Name:        group.Name,
			Description: group.Description,
			AvatarUrl:   group.AvatarUrl,
		}
	}
	return &grouppb.GetGroupsResponse{
		Groups: groupsResponse,
		Total:  total,
	}
}
