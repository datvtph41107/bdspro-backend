package client

import (
	"bdspro/env"
	"context"
	"pb/clients"
	hubpb "pb/types/hub"
)

// @bind: bdspro/internal/provider.HubProvider

type HubClient struct {
	*clients.HubGrpcClient
}

func NewHubClient() *HubClient {
	return &HubClient{
		HubGrpcClient: clients.MakeHubGrpcClient(env.ENV_RUNTIME),
	}
}

// GetLocationsByIds retrieves locations (provinces, districts, wards) by IDs
func (c *HubClient) GetLocationsByIds(ctx context.Context, ids []uint64) ([]*hubpb.Location, error) {
	if len(ids) == 0 {
		return []*hubpb.Location{}, nil
	}

	response, err := c.InternalClient.GetLocationsByIds(ctx, &hubpb.GetLocationsByIdsRequest{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetLocationsByIdsMap retrieves locations and returns as a map[id]Location
func (c *HubClient) GetLocationsByIdsMap(ctx context.Context, ids []uint64) (map[uint64]*hubpb.Location, error) {
	locations, err := c.GetLocationsByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	result := make(map[uint64]*hubpb.Location)
	for _, loc := range locations {
		result[loc.Id] = loc
	}

	return result, nil
}
