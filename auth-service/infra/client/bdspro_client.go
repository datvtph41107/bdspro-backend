package client

import (
	"auth/env"
	"auth/internal/dto"
	"auth/internal/interface/provider"
	"context"
	"errors"
	"pb/clients"
)

type BdsproClient struct {
	*clients.BdsproGrpcClient
}

func NewBdsproClient() provider.BdsproProvider {
	return &BdsproClient{
		BdsproGrpcClient: clients.MakeBdsproGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *BdsproClient) GetDealMember(ctx context.Context, dealId uint64, userId uint64) (*dto.DealMember, error) {
	if c.BdsproGrpcClient == nil {
		return nil, errors.New("bdspro grpc client not avaiable")
	}

	resp, err := c.BdsproGrpcClient.GetDealMember(ctx, dealId, userId)
	if err != nil {
		return nil, err
	}
	return &dto.DealMember{
		DealId:  resp.DealId,
		UserId:  resp.MemberId,
		RoleId:  resp.RoleId,
		RoleKey: resp.RoleKey,
	}, nil
}
