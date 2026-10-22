package provider

import (
	"context"
	"errors"

	_dto "common/domain/dto"
	"pb/clients"
)

// ProductSuggestProvider enriches AI output with catalog ids from bdspro internal GetSuggest.
type ProductSuggestProvider interface {
	GetSuggest(ctx context.Context, content string) (*_dto.ProductV3DTO, error)
}

type bdsproSuggestProvider struct {
	client *clients.BdsproGrpcClient
}

func NewBdsproSuggestProvider(client *clients.BdsproGrpcClient) ProductSuggestProvider {
	return &bdsproSuggestProvider{client: client}
}

func (p *bdsproSuggestProvider) GetSuggest(ctx context.Context, content string) (*_dto.ProductV3DTO, error) {
	if p == nil || p.client == nil {
		return nil, errors.New("bdspro grpc client not available")
	}
	return p.client.GetSuggest(ctx, content)
}
