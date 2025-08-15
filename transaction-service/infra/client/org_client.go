package client

import (
	"context"
	"pb/clients"
	organizationpb "pb/types/organization"
	"tx/env"
	"tx/internal/dto"
)

// @bind: organization/internal/interface.IOrganizationClient
type OrganizationClient struct {
	*clients.OrganizationClient
}

func NewOrganizationClient() *OrganizationClient {
	return &OrganizationClient{
		OrganizationClient: clients.MakeOrganizationClient(env.ENV_RELEASE),
	}
}

func (c *OrganizationClient) AddProductToDeal(ctx context.Context, req *dto.AddProductToDealRequest) error {
	_, err := c.OrganizationClient.Client.AddProductToDeal(ctx, &organizationpb.AddProductToDealRequest{
		DealId:    req.DealID,
		ProductId: req.ProductID,
	})
	if err != nil {
		return err
	}
	return nil
}
