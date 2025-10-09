package client

import (
	"context"
	"pb/clients"
	bdspropb "pb/types/bdspro"

	"appointment/internal/env"
)

type BdsproClient interface {
	GetProductByIds(ctx context.Context, ids []uint64) (*bdspropb.ProductInternalResponse, error)
	GetProductAttachmentByIds(ctx context.Context, ids []uint64) (*bdspropb.ProductAttachmentResponse, error)
}

type bdsproClient struct {
	client *clients.BdsproGrpcClient
}

func NewBdsproClient() BdsproClient {
	return &bdsproClient{
		client: clients.MakeBdsproGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *bdsproClient) GetProductByIds(ctx context.Context, ids []uint64) (*bdspropb.ProductInternalResponse, error) {
	return c.client.GetProductByIds(ctx, ids)
}

func (c *bdsproClient) GetProductAttachmentByIds(ctx context.Context, ids []uint64) (*bdspropb.ProductAttachmentResponse, error) {
	return c.client.GetProductAttachmentByIds(ctx, ids)
}
