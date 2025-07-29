package client

import (
	"context"
	"pb/clients"
	organizationpb "pb/types/organization"

	"crm/env"
)

type OrganizationClient struct {
	Client *clients.OrganizationMemberGrpcClient
}

func NewOrganizationClient() *OrganizationClient {
	return &OrganizationClient{
		Client: clients.NewOrganizationMemberGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *OrganizationClient) CheckUsersInOrganization(ctx context.Context, organizationID uint32, userIds []uint32) ([]*organizationpb.CheckUserInOrganization, error) {
	return c.Client.CheckUsersInOrganization(ctx, organizationID, userIds)
}
