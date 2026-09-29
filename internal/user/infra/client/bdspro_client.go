package client

import (
	"context"
	"errors"
	"pb/clients"
	sharepb "pb/types/shared"
	"user/internal/dto"
)

// @bind: user/internal/interface/providers.BdsproProvider
type BdsproClient struct {
	*clients.BdsproGrpcClient
}

func NewBdsproClient(
	rpcClient *clients.BdsproGrpcClient,
) *BdsproClient {
	return &BdsproClient{
		BdsproGrpcClient: rpcClient,
	}
}

func (c *BdsproClient) GetDealMember(ctx context.Context, dealId uint64, userId uint64) (*dto.DealMember, error) {
	if c.BdsproGrpcClient == nil {
		return nil, errors.New("bdspro grpc client not available")
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

// GetDashboardByProfileId lấy dashboard stats theo profileId (dùng bởi profile handler).
func (c *BdsproClient) GetDashboardByProfileId(ctx context.Context, profileID uint64) (*sharepb.BDSProDashboardProto, error) {
	if c.BdsproGrpcClient == nil || c.BdsproGrpcClient.InternalClient == nil {
		return nil, nil
	}

	resp, err := c.BdsproGrpcClient.GetDashboardByProfileId(ctx, profileID)
	if err != nil {
		return nil, nil
	}
	return resp, nil
}
