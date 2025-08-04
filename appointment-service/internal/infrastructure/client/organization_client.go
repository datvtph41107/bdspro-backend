package client

import (
	"context"
	"pb/clients"
	organizationpb "pb/types/organization"

	"github.com/tuandq2112/appointment-service/internal/env"
)

type OrganizationClient interface {
	GetDealById(ctx context.Context, dealId uint64) (*organizationpb.GroupDeal, error)
	GetDealTransactionsByIds(ctx context.Context, dealIds []uint64) ([]*organizationpb.DealTransactionResponse, error)
	GetDealInvitationsByIds(ctx context.Context, dealId uint64, ids []uint64) (map[uint64]*organizationpb.DealInvitationMember, error)
}

type organizationClient struct {
	client *clients.OrganizationClient
}

func NewOrganizationClient() OrganizationClient {
	return &organizationClient{
		client: clients.MakeOrganizationClient(env.ENV_RUNTIME),
	}
}

func (c *organizationClient) GetDealById(ctx context.Context, dealId uint64) (*organizationpb.GroupDeal, error) {
	return c.client.GetDealById(ctx, dealId)
}

func (c *organizationClient) GetDealTransactionsByIds(ctx context.Context, dealIds []uint64) ([]*organizationpb.DealTransactionResponse, error) {
	return c.client.GetDealTransactionsByIds(ctx, dealIds)
}

func (c *organizationClient) GetDealInvitationsByIds(ctx context.Context, dealId uint64, ids []uint64) (map[uint64]*organizationpb.DealInvitationMember, error) {
	return c.client.GetDealMemberByIds(ctx, dealId, ids)
}

