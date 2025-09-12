package transformer

import (
	organizationpb "pb/types/organization"

	"organization/internal/domain/entity"
	"organization/internal/dto"
)

type OrganizationMemberTransformer interface {
	CreateOrganizationMemberRequestToEntity(request *organizationpb.CreateOrganizationMemberRequest) *entity.OrganizationMember
	EntityToCreateOrganizationMemberResponse(organizationMember *entity.OrganizationMember) *organizationpb.CreateOrganizationMemberResponse
	// CreateOrganizationMemberBatchRequestToEntities(request *organizationpb.CreateOrganizationMemberBatchRequest) []*entity.OrganizationMember
	EntitiesToCreateOrganizationMemberBatchResponse(organizationMembers []*entity.OrganizationMember) *organizationpb.CreateOrganizationMemberBatchResponse
	UpdateOrganizationMemberRequestToEntity(request *organizationpb.UpdateOrganizationMemberRequest) *entity.OrganizationMember
	EntityToUpdateOrganizationMemberResponse(organizationMember *entity.OrganizationMember) *organizationpb.UpdateOrganizationMemberResponse
	DeleteOrganizationMemberRequestToEntity(request *organizationpb.DeleteOrganizationMemberRequest) *entity.OrganizationMember
	EntityToDeleteOrganizationMemberResponse(organizationMember *entity.OrganizationMember) *organizationpb.DeleteOrganizationMemberResponse
	EntityToPbsResponse(organizationMembers []*dto.OrganizationMemberWithProfile, total uint32, hidePhone bool) *organizationpb.GetOrganizationMembersResponse
	// EntityToGetOrganizationOfMemberResponse(organization *entity.Organization) *organizationpb.GetOrganizationMembersResponse
}

type organizationMemberTransformer struct{}

func NewOrganizationMemberTransformer() OrganizationMemberTransformer {
	return &organizationMemberTransformer{}
}

func (o *organizationMemberTransformer) CreateOrganizationMemberRequestToEntity(request *organizationpb.CreateOrganizationMemberRequest) *entity.OrganizationMember {
	return &entity.OrganizationMember{
		OrganizationID: request.OrganizationId,
		UserID:         request.UserId,
		RoleId:         uint64(request.Role),
	}
}

func (o *organizationMemberTransformer) EntityToCreateOrganizationMemberResponse(organizationMember *entity.OrganizationMember) *organizationpb.CreateOrganizationMemberResponse {
	return &organizationpb.CreateOrganizationMemberResponse{
		Id: organizationMember.ID,
	}
}

// func (o *organizationMemberTransformer) CreateOrganizationMemberBatchRequestToEntities(request *organizationpb.CreateOrganizationMemberBatchRequest) []*entity.OrganizationMember {
// 	entities := make([]*entity.OrganizationMember, len(request.Members))
// 	for i, memberRequest := range request.Members {
// 		entities[i] = o.CreateOrganizationMemberRequestToEntity(memberRequest)
// 	}
// 	return entities
// }

func (o *organizationMemberTransformer) EntitiesToCreateOrganizationMemberBatchResponse(organizationMembers []*entity.OrganizationMember) *organizationpb.CreateOrganizationMemberBatchResponse {
	ids := make([]uint32, len(organizationMembers))
	for i, member := range organizationMembers {
		ids[i] = member.ID
	}
	return &organizationpb.CreateOrganizationMemberBatchResponse{
		Ids: ids,
	}
}

func (o *organizationMemberTransformer) UpdateOrganizationMemberRequestToEntity(request *organizationpb.UpdateOrganizationMemberRequest) *entity.OrganizationMember {
	return &entity.OrganizationMember{
		ID:      request.Id,
		Status:  entity.OrganizationMemberStatus(request.Status),
		RoleKey: request.RoleKey,
	}
}

func (o *organizationMemberTransformer) EntityToUpdateOrganizationMemberResponse(organizationMember *entity.OrganizationMember) *organizationpb.UpdateOrganizationMemberResponse {
	return &organizationpb.UpdateOrganizationMemberResponse{
		Id: organizationMember.ID,
	}
}

func (o *organizationMemberTransformer) DeleteOrganizationMemberRequestToEntity(request *organizationpb.DeleteOrganizationMemberRequest) *entity.OrganizationMember {
	return &entity.OrganizationMember{
		ID: request.Id,
	}
}

func (o *organizationMemberTransformer) EntityToDeleteOrganizationMemberResponse(organizationMember *entity.OrganizationMember) *organizationpb.DeleteOrganizationMemberResponse {
	return &organizationpb.DeleteOrganizationMemberResponse{
		Id: organizationMember.ID,
	}
}

func (o *organizationMemberTransformer) EntityToPbsResponse(members []*dto.OrganizationMemberWithProfile, total uint32, hidePhone bool) *organizationpb.GetOrganizationMembersResponse {
	membersResponse := make([]*organizationpb.OrganizationMember, len(members))
	for i, member := range members {
		membersResponse[i] = &organizationpb.OrganizationMember{
			Id:       member.OrganizationMember.ID,
			UserId:   member.OrganizationMember.UserID,
			RoleId:   member.OrganizationMember.RoleId, // Add roleId field
			RoleKey:  member.OrganizationMember.RoleKey,
			Status:   uint32(member.OrganizationMember.Status),
			RoleName: member.OrganizationMember.RoleName,
		}
		profile := member.Profile
		if profile != nil {
			if hidePhone {
				profile.Phone = ""
			}
			membersResponse[i].Profile = profile
		}
		if member.OrganizationMember.DealMember != nil {
			membersResponse[i].DealMember = &organizationpb.DealMember{
				MemberId: member.OrganizationMember.DealMember.MemberID,
				RoleId:   &member.OrganizationMember.DealMember.RoleID,
				// RespondedAt: _utils.FormatTimeToString(&member.OrganizationMember.DealMember.CreatedAt),
				Status:  uint32(member.OrganizationMember.DealMember.Status),
				RoleKey: uint32(member.OrganizationMember.DealMember.RoleKey),
			}
		}
	}
	return &organizationpb.GetOrganizationMembersResponse{
		Data:  membersResponse,
		Total: total,
	}
}

// func (o *organizationMemberTransformer) EntityToGetOrganizationOfMemberResponse(organization *entity.Organization) *organizationpb.GetOrganizationMembersResponse {
// 	return &organizationpb.GetOrganizationMembersResponse{
// 		OrganizationMembers: []*organizationpb.OrganizationMember{
// 			{
// 				Id:     organization.ID,
// 				UserId: organization.CreatedBy,
// 				Role:   organization.CreatedBy,
// 				Status: getStatusString(entity.OrganizationMemberStatusActive),
// 			},
// 		},
// 	}
// }

func getStatusString(status entity.OrganizationMemberStatus) string {
	switch status {
	case entity.OrganizationMemberStatusActive:
		return "active"
	case entity.OrganizationMemberStatusInactive:
		return "inactive"
	case entity.OrganizationMemberStatusPending:
		return "pending"
	default:
		return "unknown"
	}
}
