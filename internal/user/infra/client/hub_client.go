package client

import (
	"context"
	"fmt"
	"pb/clients"
	hubpb "pb/types/hub"
)

// @bind: user/internal/interface/providers.HubProvider
type HubClient struct {
	*clients.HubGrpcClient
}

func NewHubClient(
	rpcClient *clients.HubGrpcClient,
) *HubClient {
	return &HubClient{
		HubGrpcClient: rpcClient,
	}
}

// GetLocationsByIds lấy thông tin province/ward theo IDs
func (c *HubClient) GetLocationsByIds(ctx context.Context, ids []uint64) (map[uint64]*hubpb.Location, error) {
	if len(ids) == 0 {
		return make(map[uint64]*hubpb.Location), nil
	}

	if c.HubGrpcClient.InternalClient == nil {
		return make(map[uint64]*hubpb.Location), nil
	}

	resp, err := c.HubGrpcClient.InternalClient.GetLocationsByIds(ctx, &hubpb.GetLocationsByIdsRequest{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}

	// Convert slice to map for easy lookup
	locationMap := make(map[uint64]*hubpb.Location)
	for _, location := range resp.Data {
		locationMap[location.Id] = location
	}

	return locationMap, nil
}

// GetProvinceAndWardNames lấy tên province và ward
func (c *HubClient) GetProvinceAndWardNames(ctx context.Context, provinceId, wardId *uint64) (provinceName, wardName string) {
	if c.HubGrpcClient.InternalClient == nil {
		return "", ""
	}

	resp, err := c.HubGrpcClient.InternalClient.GetAddressV2ByIds(ctx, &hubpb.GetAddressV2ByIdsRequest{
		ProvinceId: provinceId,
		WardId:     wardId,
	})
	if err != nil || resp == nil || resp.Data == nil {
		return "", ""
	}

	if resp.Data.ProvinceName != nil {
		provinceName = *resp.Data.ProvinceName
	}

	if resp.Data.WardName != nil {
		wardName = *resp.Data.WardName
	}

	return provinceName, wardName
}

// GetUserGuideByKey lấy user guide theo key
func (c *HubClient) GetUserGuideByKey(ctx context.Context, key string) (*hubpb.UserGuideDetail, error) {
	if c.HubGrpcClient.UserGuideClient == nil {
		return nil, fmt.Errorf("user guide client not available")
	}

	resp, err := c.HubGrpcClient.UserGuideClient.GetUserGuideByKey(ctx, &hubpb.GetUserGuideByKeyRequest{
		Key: key,
	})
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}
