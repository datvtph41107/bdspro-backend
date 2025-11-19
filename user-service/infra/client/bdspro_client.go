package client

import (
	"context"
	"pb/clients"
	sharepb "pb/types/shared"
	"user/env"
)

type BdsproClient struct {
	*clients.BdsproGrpcClient
}

func NewBdsproClient() *BdsproClient {
	return &BdsproClient{
		BdsproGrpcClient: clients.MakeBdsproGrpcClient(env.ENV_RUNTIME),
	}
}

// GetDashboardByProfileId lấy dashboard stats theo profileId
func (c *BdsproClient) GetDashboardByProfileId(ctx context.Context, profileID uint64) (*sharepb.BDSProDashboardProto, error) {
	if c.BdsproGrpcClient.InternalClient == nil {
		return nil, nil
	}

	resp, err := c.BdsproGrpcClient.GetDashboardByProfileId(ctx, profileID)
	if err != nil {
		return nil, nil
	}

	return resp, nil
}
