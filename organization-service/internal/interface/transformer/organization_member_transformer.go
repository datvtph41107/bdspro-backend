package transformer

import (
	organizationpb "pb/types/organization"

	"organization/internal/domain/entity"
	"organization/internal/interface/dto"
)

type OrganizationMemberTransformer interface {
	CreateOrganizationMemberRequestToEntity(request *organizationpb.CreateOrganizationMemberRequest) *entity.OrganizationMember
	EntityToCreateOrganizationMemberResponse(organizationMember *entity.OrganizationMember) *organizationpb.CreateOrganizationMemberResponse
	UpdateOrganizationMemberRequestToEntity(request *organizationpb.UpdateOrganizationMemberRequest) *entity.OrganizationMember
	EntityToUpdateOrganizationMemberResponse(organizationMember *entity.OrganizationMember) *organizationpb.UpdateOrganizationMemberResponse
	DeleteOrganizationMemberRequestToEntity(request *organizationpb.DeleteOrganizationMemberRequest) *entity.OrganizationMember
	EntityToDeleteOrganizationMemberResponse(organizationMember *entity.OrganizationMember) *organizationpb.DeleteOrganizationMemberResponse
	EntityToGetOrganizationMembersResponse(organizationMembers []*dto.OrganizationMemberWithProfile, total uint32) *organizationpb.GetOrganizationMembersResponse
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
		RoleId:         request.Role,
	}
}

func (o *organizationMemberTransformer) EntityToCreateOrganizationMemberResponse(organizationMember *entity.OrganizationMember) *organizationpb.CreateOrganizationMemberResponse {
	return &organizationpb.CreateOrganizationMemberResponse{
		Id: organizationMember.ID,
	}
}

func (o *organizationMemberTransformer) UpdateOrganizationMemberRequestToEntity(request *organizationpb.UpdateOrganizationMemberRequest) *entity.OrganizationMember {
	return &entity.OrganizationMember{
		ID:     request.Id,
		RoleId: request.Role,
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

func (o *organizationMemberTransformer) EntityToGetOrganizationMembersResponse(members []*dto.OrganizationMemberWithProfile, total uint32) *organizationpb.GetOrganizationMembersResponse {
	membersResponse := make([]*organizationpb.OrganizationMember, len(members))
	for i, member := range members {
		membersResponse[i] = &organizationpb.OrganizationMember{
			Id:     member.OrganizationMember.ID,
			UserId: member.OrganizationMember.UserID,
			Role:   member.OrganizationMember.RoleId,
			Status: getStatusString(member.OrganizationMember.Status),
			// Profile: profile,
		}
		profile := member.Profile
		if profile != nil {
			profile.Phone = ""
			membersResponse[i].Profile = profile
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
