package client

import (
	"bdspro/env"
	"context"
	"pb/clients"
	sharepb "pb/types/shared"
)

type CrmClient struct {
	*clients.CrmGrpcClient
}

func NewCrmClient() *CrmClient {
	return &CrmClient{
		clients.MakeCrmGrpcClient(env.ENV_RUNTIME),
	}
}

// GetOriginProfileByOriginIds retrieves origin profiles by origin IDs using internal service
func (c *CrmClient) GetOriginProfileByOriginIds(ctx context.Context, originIDs []uint64) ([]*sharepb.OriginProfile, error) {
	if len(originIDs) == 0 {
		return []*sharepb.OriginProfile{}, nil
	}

	if c.CrmInternalServiceClient == nil {
		return []*sharepb.OriginProfile{}, nil
	}

	response, err := c.CrmInternalServiceClient.GetOriginProfileByOriginIds(ctx, &sharepb.IdRequest{
		Ids: originIDs,
	})
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetContactById retrieves contact information by ID using internal service
// TODO: Uncomment sau khi generate proto code với method GetContactById
