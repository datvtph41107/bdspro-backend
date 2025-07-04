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
