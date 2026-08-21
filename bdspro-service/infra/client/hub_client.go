package client

import (
	"context"
	"pb/clients"
	hubpb "pb/types/hub"

	admin_usecases "bdspro/internal/usecases/admin"
)

// @bind: bdspro/internal/provider.HubProvider
type HubClient struct {
	*clients.HubGrpcClient
	areaRegionUsecase *admin_usecases.AreaRegionUsecase
}

func NewHubClient(
	rpcClient *clients.HubGrpcClient,
	areaRegionUsecase *admin_usecases.AreaRegionUsecase,
) *HubClient {
	return &HubClient{
		HubGrpcClient:     rpcClient,
		areaRegionUsecase: areaRegionUsecase,
	}
}

func (c *HubClient) GetRegionsByIdsMap(ctx context.Context, ids []uint64) (map[uint64]*hubpb.Region, error) {
	if len(ids) == 0 {
		return map[uint64]*hubpb.Region{}, nil
	}
	// Use local AreaRegion (bdspro) instead of calling hub
	if c.areaRegionUsecase != nil {
		regions, err := c.areaRegionUsecase.GetByIDs(ctx, ids)
		if err != nil {
			return nil, err
		}
		result := make(map[uint64]*hubpb.Region)
		for i := range regions {
			r := &regions[i]
			result[r.ID] = &hubpb.Region{
				Id:     r.ID,
				Name:   r.Name,
				Code:   r.Code,
				Active: r.Active,
			}
		}
		return result, nil
	}
	resp, err := c.InternalClient.GetByIDs(ctx, &hubpb.GetRegionsRequest{Ids: ids})
	if err != nil {
		return nil, err
	}
	result := make(map[uint64]*hubpb.Region)
	for _, r := range resp.Regions {
		result[r.Id] = r
	}
	return result, nil
}

func (c *HubClient) GetProductStatsViewClient(
	ctx context.Context,
	productID uint64,
	fromTime int64,
) (*hubpb.ProductStatsViewResponse, error) {
	resp, err := c.HubGrpcClient.InternalClient.GetProductViewStats(ctx, &hubpb.ProductStatsViewRequest{
		ProductId: productID,
		FromTime:  fromTime,
	})
	if err != nil {
		return nil, err
	}

	if resp.ViewCount == 0 {
		return nil, nil
	}

	return &hubpb.ProductStatsViewResponse{
		ViewCount:     resp.ViewCount,
		TotalDuration: resp.TotalDuration,
		MaxEventTime:  resp.MaxEventTime,
	}, nil
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

// func (c *HubClient) PutUpdate(ctx context.Context, resource string, resourceId uint64, ownerIds []uint64, updatedAt int64) error {

// 	return c.PutUpdate(ctx, resource, resourceId, )
// }
// func (c *HubClient) DelUpdate(ctx context.Context, resource string, id uint64, ownerId uint64) error {
// 	return nil
// }
