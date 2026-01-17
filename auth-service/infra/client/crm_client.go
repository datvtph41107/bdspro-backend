package client

import (
	"auth/env"
	"context"
	"log"
	"pb/clients"
	sharepb "pb/types/shared"
)

// @bind: auth/internal/interface/provider.CrmProvider
type CrmClient struct {
	*clients.CrmGrpcClient
}

func NewCrmClient() *CrmClient {
	return &CrmClient{
		CrmGrpcClient: clients.MakeCrmGrpcClient(env.ENV_RUNTIME),
	}
}

// GetOriginProfile lấy origin profile từ CRM service theo profileID
func (c *CrmClient) GetOriginID(ctx context.Context, profileID uint64) (uint64, error) {
	if c.CrmGrpcClient == nil || c.CrmGrpcClient.CrmInternalServiceClient == nil {
		log.Printf("CRM client is nil")
		return 0, nil
	}

	resp, err := c.CrmGrpcClient.CrmInternalServiceClient.GetOriginProfileByProfileID(ctx, &sharepb.IdRequest{
		Id: profileID,
	})
	if err != nil {
		log.Printf("Failed to get origin profile: %v", err)
		return 0, err
	}

	if resp == nil {
		return 0, nil
	}

	// Trả về origin profile đầu tiên
	return resp.OriginId, nil
}
