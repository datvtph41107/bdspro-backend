package transformer

import (
	organizationpb "pb/types/organization"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type OrganizationMemberTransformer interface {
	CreateOrganizationMemberRequestToEntity(request *organizationpb.CreateOrganizationMemberRequest) (*entity.OrganizationMember)
	EntityToCreateOrganizationMemberResponse(organizationMember *entity.OrganizationMember) *organizationpb.CreateOrganizationMemberResponse
	UpdateOrganizationMemberRequestToEntity(request *organizationpb.UpdateOrganizationMemberRequest) *entity.OrganizationMember
	EntityToUpdateOrganizationMemberResponse(organizationMember *entity.OrganizationMember) *organizationpb.UpdateOrganizationMemberResponse
	DeleteOrganizationMemberRequestToEntity(request *organizationpb.DeleteOrganizationMemberRequest) *entity.OrganizationMember
	EntityToDeleteOrganizationMemberResponse(organizationMember *entity.OrganizationMember) *organizationpb.DeleteOrganizationMemberResponse
	EntityToGetOrganizationMembersResponse(organizationMembers []*entity.OrganizationMember) *organizationpb.GetOrganizationMembersResponse
}

type organizationMemberTransformer struct{}

func NewOrganizationMemberTransformer() OrganizationMemberTransformer {
	return &organizationMemberTransformer{}
}

func (o *organizationMemberTransformer) CreateOrganizationMemberRequestToEntity(request *organizationpb.CreateOrganizationMemberRequest) (*entity.OrganizationMember,) {
	return &entity.OrganizationMember{
		OrganizationID: request.OrganizationId,
		UserID: request.UserId,
		RoleId: request.Role,
	}
}
func (o *organizationMemberTransformer) EntityToCreateOrganizationMemberResponse(organizationMember *entity.OrganizationMember) *organizationpb.CreateOrganizationMemberResponse {
	return &organizationpb.CreateOrganizationMemberResponse{
		Id: organizationMember.ID,
	}
}

func (o *organizationMemberTransformer) UpdateOrganizationMemberRequestToEntity(request *organizationpb.UpdateOrganizationMemberRequest) *entity.OrganizationMember {
	return &entity.OrganizationMember{
		ID: request.Id,
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

func (o *organizationMemberTransformer) EntityToGetOrganizationMembersResponse(organizationMembers []*entity.OrganizationMember) *organizationpb.GetOrganizationMembersResponse {
	organizationMembersResponse := make([]*organizationpb.OrganizationMember, len(organizationMembers))
	for i, organizationMember := range organizationMembers {
		organizationMembersResponse[i] = &organizationpb.OrganizationMember{
			Id: organizationMember.ID,
			UserId: organizationMember.UserID,
			Role: organizationMember.RoleId,
		}
	}
	return &organizationpb.GetOrganizationMembersResponse{
		OrganizationMembers: organizationMembersResponse,
	}
}