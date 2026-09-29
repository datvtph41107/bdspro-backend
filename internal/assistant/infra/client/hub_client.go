package client

import (
	"context"
	"pb/clients"
	hubpb "pb/types/hub"
)

// @bind: assistant/internal/interface/provider.HubProvider
type HubClient struct {
	*clients.HubGrpcClient
}

func NewHubClient(rpcClient *clients.HubGrpcClient) *HubClient {
	return &HubClient{HubGrpcClient: rpcClient}
}

func (c *HubClient) SearchLocationV2(ctx context.Context, keyword string, limit int32) ([]*hubpb.LocationSearchResultV2, error) {
	if keyword == "" {
		return []*hubpb.LocationSearchResultV2{}, nil
	}

	if c == nil || c.HubGrpcClient == nil || c.LocationV2Client == nil {
		return []*hubpb.LocationSearchResultV2{}, nil
	}

	if limit <= 0 {
		limit = 5
	}

	resp, err := c.LocationV2Client.SearchLocationV2(ctx, &hubpb.SearchLocationV2Request{
		Keyword: keyword,
		Page:    1,
		Size:    limit,
	})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return []*hubpb.LocationSearchResultV2{}, nil
	}

	return resp.GetData(), nil
}
