package client

import (
	"bdspro/internal/enums"
	"context"
	"pb/clients"
	bdspropb "pb/types/bdspro"
)

type OrganizationClient struct {
	*clients.OrganizationClient
}

func NewOrganizationClient(
	rpcClient *clients.OrganizationClient,
) *OrganizationClient {
	return &OrganizationClient{
		OrganizationClient: rpcClient,
	}
}

func (c *OrganizationClient) MapOrganizationToSharingAccess(
	ctx context.Context,
	sharingAccesses []*bdspropb.SharingAccess) {
	if c.Client == nil {
		return
	}

	organizationIDSet := make(map[uint64]struct{})
	for _, sharingAccess := range sharingAccesses {
		if sharingAccess.ToType == enums.EOwnerOfOrganization {
			organizationIDSet[sharingAccess.ToId] = struct{}{}
		}
	}
	organizations, err := c.OrganizationClient.GetOrganizationMapByIDs(ctx, organizationIDSet)
	if err != nil {
		return
	}

	for _, sharingAccess := range sharingAccesses {
		if sharingAccess.ToType == enums.EOwnerOfOrganization && organizations[sharingAccess.ToId] != nil {
			sharingAccess.Avatar = organizations[sharingAccess.ToId].Avatar
			sharingAccess.ToName = organizations[sharingAccess.ToId].Name
		}
	}
}

func (c *OrganizationClient) MapGroupToSharingAccess(
	ctx context.Context,
	sharingAccesses []*bdspropb.SharingAccess) {
	if c.Client == nil {
		return
	}

	groupIDSet := make(map[uint64]struct{})
	for _, sharingAccess := range sharingAccesses {
		if sharingAccess.ToType == enums.EOwnerOfGroup {
			groupIDSet[sharingAccess.ToId] = struct{}{}
		}
	}
	groups, err := c.OrganizationClient.GetGroupMapByIDs(ctx, groupIDSet)
	if err != nil {
		return
	}

	for _, sharingAccess := range sharingAccesses {
		if sharingAccess.ToType == enums.EOwnerOfGroup && groups[sharingAccess.ToId] != nil {
			sharingAccess.Avatar = groups[sharingAccess.ToId].Avatar
			sharingAccess.ToName = groups[sharingAccess.ToId].Name
		}
	}
}
