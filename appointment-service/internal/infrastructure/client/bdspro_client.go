package client

import (
	"context"
	"pb/clients"
	pb_bdspro "pb/types/bdspro"

	"github.com/tuandq2112/appointment-service/internal/env"
)

type BdsproClient interface {
	GetProductByIds(ctx context.Context, ids []uint64) (*pb_bdspro.ProductInternalResponse, error)
	GetProductAttachmentByIds(ctx context.Context, ids []uint64) (*pb_bdspro.ProductAttachmentResponse, error)
}

type bdsproClient struct {
	client *clients.BdsproGrpcClient
}

func NewBdsproClient() BdsproClient {
	return &bdsproClient{
		client: clients.MakeBdsproGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *bdsproClient) GetProductByIds(ctx context.Context, ids []uint64) (*pb_bdspro.ProductInternalResponse, error) {
	return c.client.GetProductByIds(ctx, ids)
}

func (c *bdsproClient) GetProductAttachmentByIds(ctx context.Context, ids []uint64) (*pb_bdspro.ProductAttachmentResponse, error) {
	return c.client.GetProductAttachmentByIds(ctx, ids)
}