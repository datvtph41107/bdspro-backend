package transformer

import (
	grouppb "pb/types/group"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/interface/dto"
)

type GroupMemberTransformer interface {
	CreateGroupMemberRequestToEntity(request *grouppb.CreateGroupMemberRequest) *entity.GroupMember
	UpdateGroupMemberRequestToEntity(request *grouppb.UpdateGroupMemberRequest) *entity.GroupMember
	EntityToCreateGroupMemberResponse(member *entity.GroupMember) *grouppb.CreateGroupMemberResponse
	EntityToUpdateGroupMemberResponse(member *entity.GroupMember) *grouppb.UpdateGroupMemberResponse
	EntityToGetGroupMemberResponse(member *entity.GroupMember) *grouppb.GetGroupMemberResponse
	EntityToGetGroupMembersResponse(members []*dto.GroupMemberWithProfile) *grouppb.GetGroupMembersResponse
}

type groupMemberTransformer struct{}

func NewGroupMemberTransformer() GroupMemberTransformer {
	return &groupMemberTransformer{}
}

func (t *groupMemberTransformer) CreateGroupMemberRequestToEntity(request *grouppb.CreateGroupMemberRequest) *entity.GroupMember {
	return &entity.GroupMember{
		GroupId: request.GroupId,
		UserId:  uint32(request.UserId),
		Role:    entity.GroupMemberRole(request.Role),
	}
}

func (t *groupMemberTransformer) UpdateGroupMemberRequestToEntity(request *grouppb.UpdateGroupMemberRequest) *entity.GroupMember {
	return &entity.GroupMember{
		Id:   request.Id,
		Role: entity.GroupMemberRole(request.Role),
	}
}

func (t *groupMemberTransformer) EntityToCreateGroupMemberResponse(member *entity.GroupMember) *grouppb.CreateGroupMemberResponse {
	return &grouppb.CreateGroupMemberResponse{
		Id: member.Id,
	}
}

func (t *groupMemberTransformer) EntityToUpdateGroupMemberResponse(member *entity.GroupMember) *grouppb.UpdateGroupMemberResponse {
	return &grouppb.UpdateGroupMemberResponse{
		Id: member.Id,
	}
}

func (t *groupMemberTransformer) EntityToGetGroupMemberResponse(member *entity.GroupMember) *grouppb.GetGroupMemberResponse {
	return &grouppb.GetGroupMemberResponse{
		Member: &grouppb.GroupMember{
			Id:      member.Id,
			GroupId: member.GroupId,
			UserId:  member.UserId,
			Role:    string(member.Role),
		},
	}
}

func (t *groupMemberTransformer) EntityToGetGroupMembersResponse(members []*dto.GroupMemberWithProfile) *grouppb.GetGroupMembersResponse {
	response := &grouppb.GetGroupMembersResponse{
		Members: make([]*grouppb.GroupMember, 0, len(members)),
	}

	for _, member := range members {
		profile := member.Profile
		profile.Phone = ""
		response.Members = append(response.Members, &grouppb.GroupMember{
			Id:      member.GroupMember.Id,
			GroupId: member.GroupMember.GroupId,
			UserId:  member.GroupMember.UserId,
			Role:    string(member.GroupMember.Role),
			Profile: profile,
		})
	}

	return response
}
