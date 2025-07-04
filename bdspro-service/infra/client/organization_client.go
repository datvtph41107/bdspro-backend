package client

import (
	"bdspro/env"
	"bdspro/internal/enums"
	"context"
	"pb/clients"
	pb_bdspro "pb/types/bdspro"
)

type OrganizationClient struct {
	Client *clients.OrganizationClient
}

func NewOrganizationClient() *OrganizationClient {
	return &OrganizationClient{
		Client: clients.MakeOrganizationClient(env.ENV_RELEASE),
	}
}

func (c *OrganizationClient) MapOrganizationToSharingAccess(
	ctx context.Context,
	sharingAccesses []*pb_bdspro.SharingAccess) {
	if c.Client == nil {
		return
	}

	organizationIDSet := make(map[uint64]struct{})
	for _, sharingAccess := range sharingAccesses {
		if sharingAccess.ToType == enums.EOwnerTypeOrgnization {
			organizationIDSet[sharingAccess.ToId] = struct{}{}
		}
	}
	organizations, err := c.Client.GetOrganizationMapByIDs(ctx, organizationIDSet)
	if err != nil {
		return
	}

	for _, sharingAccess := range sharingAccesses {
		if sharingAccess.ToType == enums.EOwnerTypeOrgnization && organizations[sharingAccess.ToId] != nil {
			sharingAccess.Avatar = organizations[sharingAccess.ToId].Avatar
			sharingAccess.ToName = organizations[sharingAccess.ToId].Name
		}
	}
}

func (c *OrganizationClient) MapGroupToSharingAccess(
	ctx context.Context,
	sharingAccesses []*pb_bdspro.SharingAccess) {
	if c.Client == nil {
		return
	}

	groupIDSet := make(map[uint64]struct{})
	for _, sharingAccess := range sharingAccesses {
		if sharingAccess.ToType == enums.EOwnerTypeGroup {
			groupIDSet[sharingAccess.ToId] = struct{}{}
		}
	}
	groups, err := c.Client.GetGroupMapByIDs(ctx, groupIDSet)
	if err != nil {
		return
	}

	for _, sharingAccess := range sharingAccesses {
		if sharingAccess.ToType == enums.EOwnerTypeGroup && groups[sharingAccess.ToId] != nil {
			sharingAccess.Avatar = groups[sharingAccess.ToId].Avatar
			sharingAccess.ToName = groups[sharingAccess.ToId].Name
		}
	}
}
