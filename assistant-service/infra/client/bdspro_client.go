package client

import (
	"pb/clients"
)

// @bind: assistant/internal/interface/provider.BdsproInternalProvider
type BdsproInternalClient struct {
	*clients.BdsproGrpcClient
}

func NewBdsproInternalClient(rpcClient *clients.BdsproGrpcClient) *BdsproInternalClient {
	return &BdsproInternalClient{BdsproGrpcClient: rpcClient}
}

// func (c *BdsproInternalClient) GetSuggest(ctx context.Context, req string) (*_dto.ProductV3DTO, error) {
// 	if c.BdsproGrpcClient == nil {
// 		return nil, errors.New("bdspro grpc client not initialized")
// 	}

// 	response, err := c.GetSuggest(ctx, &bdspropb.SuggestRequest{Content: req})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &_dto.ProductV3DTO{
// 		ID:          response.Id,
// 		Name:        response.Name,
// 		Description: response.Description,
// 	}, nil
// }
