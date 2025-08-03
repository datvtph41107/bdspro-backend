package transformer

import (
	"organization/internal/domain/entity"
	organizationpb "pb/types/organization"
)

type InternalOrganizationMapper struct {
}

func NewInternalOrganizationMapper() *InternalOrganizationMapper {
	return &InternalOrganizationMapper{}
}

func (m *InternalOrganizationMapper) OrganizationToPb(organization *entity.Organization) *organizationpb.OrganizationItem {
	result := &organizationpb.OrganizationItem{
		Id:   uint64(organization.ID),
		Name: organization.Name,
	}

	if organization.LogoUrl != nil {
		result.Avatar = *organization.LogoUrl
	}

	return result
}

func (m *InternalOrganizationMapper) OrganizationsToPb(organizations []*entity.Organization) []*organizationpb.OrganizationItem {
	pbOrganizations := make([]*organizationpb.OrganizationItem, len(organizations))
	for i, organization := range organizations {
		pbOrganizations[i] = m.OrganizationToPb(organization)
	}
	return pbOrganizations
}

func (m *InternalOrganizationMapper) GroupToPb(group *entity.Group) *organizationpb.GroupItem {
	result := &organizationpb.GroupItem{
		Id:     uint64(group.Id),
		Name:   group.Name,
		Avatar: group.AvatarUrl,
	}

	return result
}

func (m *InternalOrganizationMapper) GroupsToPb(groups []*entity.Group) []*organizationpb.GroupItem {
	pbGroups := make([]*organizationpb.GroupItem, len(groups))
	for i, group := range groups {
		pbGroups[i] = m.GroupToPb(group)
	}
	return pbGroups
}

func (m *InternalOrganizationMapper) DealInvitationsToPb(dealInvitations []*entity.DealInvitation) []*organizationpb.DealInvitationMember {
	pbDealInvitations := make([]*organizationpb.DealInvitationMember, len(dealInvitations))
	for i, dealInvitation := range dealInvitations {
		pbDealInvitations[i] = m.DealInvitationToPb(dealInvitation)
	}
	return pbDealInvitations
}

func (m *InternalOrganizationMapper) DealInvitationToPb(dealInvitation *entity.DealInvitation) *organizationpb.DealInvitationMember {
	return &organizationpb.DealInvitationMember{
		UserId: uint64(dealInvitation.InviteeID),
		Status: uint32(dealInvitation.Status),
		// JoinedAt: dealInvitation.RespondedAt,
	}
}
