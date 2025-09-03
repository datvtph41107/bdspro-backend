package transformer

import (
	organizationpb "pb/types/organization"

	"organization/internal/domain/entity"
	"organization/internal/dto"
)

type GroupMemberTransformer interface {
	CreateGroupMemberRequestToEntity(request *organizationpb.CreateGroupMemberRequest) *entity.GroupMember
	UpdateGroupMemberRequestToEntity(request *organizationpb.UpdateGroupMemberRequest) *entity.GroupMember
	EntityToCreateGroupMemberResponse(member *entity.GroupMember) *organizationpb.CreateGroupMemberResponse
	EntityToUpdateGroupMemberResponse(member *entity.GroupMember) *organizationpb.UpdateGroupMemberResponse
	EntityToGetGroupMemberResponse(member *entity.GroupMember) *organizationpb.GetGroupMemberResponse
	EntityToGetGroupMembersResponse(members []*dto.GroupMemberWithProfile) *organizationpb.GetGroupMembersResponse
	EntitiesToCreateGroupMemberBatchResponse(members []*entity.GroupMember) *organizationpb.CreateGroupMemberBatchResponse
}

type groupMemberTransformer struct{}

func NewGroupMemberTransformer() GroupMemberTransformer {
	return &groupMemberTransformer{}
}

func (t *groupMemberTransformer) CreateGroupMemberRequestToEntity(request *organizationpb.CreateGroupMemberRequest) *entity.GroupMember {
	return &entity.GroupMember{
		GroupId: request.GroupId,
		UserId:  uint32(request.UserId),
		Role:    entity.GroupMemberRole(request.Role),
		RoleId:  request.RoleId,
	}
}

func (t *groupMemberTransformer) UpdateGroupMemberRequestToEntity(request *organizationpb.UpdateGroupMemberRequest) *entity.GroupMember {
	return &entity.GroupMember{
		Id:     request.Id,
		Role:   entity.GroupMemberRole(request.Role),
		RoleId: request.RoleId,
	}
}

func (t *groupMemberTransformer) EntityToCreateGroupMemberResponse(member *entity.GroupMember) *organizationpb.CreateGroupMemberResponse {
	return &organizationpb.CreateGroupMemberResponse{
		Id: member.Id,
	}
}

func (t *groupMemberTransformer) EntityToUpdateGroupMemberResponse(member *entity.GroupMember) *organizationpb.UpdateGroupMemberResponse {
	return &organizationpb.UpdateGroupMemberResponse{
		Id: member.Id,
	}
}

func (t *groupMemberTransformer) EntityToGetGroupMemberResponse(member *entity.GroupMember) *organizationpb.GetGroupMemberResponse {
	return &organizationpb.GetGroupMemberResponse{
		Member: &organizationpb.GroupMember{
			Id:      member.Id,
			GroupId: member.GroupId,
			UserId:  member.UserId,
			Role:    string(member.Role),
			RoleId:  member.RoleId,
		},
	}
}

func (t *groupMemberTransformer) EntityToGetGroupMembersResponse(members []*dto.GroupMemberWithProfile) *organizationpb.GetGroupMembersResponse {
	response := &organizationpb.GetGroupMembersResponse{
		Data: make([]*organizationpb.GroupMember, 0, len(members)),
	}

	for _, member := range members {
		profile := member.Profile
		profile.Phone = ""
		response.Data = append(response.Data, &organizationpb.GroupMember{
			Id:      member.GroupMember.Id,
			GroupId: member.GroupMember.GroupId,
			UserId:  member.GroupMember.UserId,
			Role:    string(member.GroupMember.Role),
			RoleId:  member.GroupMember.RoleId,
			Profile: profile,
		})
	}

	return response
}

func (t *groupMemberTransformer) EntitiesToCreateGroupMemberBatchResponse(members []*entity.GroupMember) *organizationpb.CreateGroupMemberBatchResponse {
	response := &organizationpb.CreateGroupMemberBatchResponse{
		Ids: make([]uint32, 0, len(members)),
	}

	return response
}
