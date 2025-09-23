package client

import (
	"context"
	"pb/clients"
	transactionpb "pb/types/transaction"
	"tx/env"
)

// type BdsproClient interface {
// 	GetProductByIds(ctx context.Context, in *bdspropb.IdRequest) (*bdspropb.ProductInternalResponse, error)
// }

type BdsproClient struct {
	*clients.BdsproGrpcClient
}

func NewBdsproClient() *BdsproClient {
	bdsproClient := &BdsproClient{
		BdsproGrpcClient: clients.MakeBdsproGrpcClient(env.ENV_RUNTIME),
	}
	return bdsproClient
}

func (c *BdsproClient) MapProductToDealPb(ctx context.Context, deals []*transactionpb.DealContractResponse) {
	// ids := make([]uint64, len(deals))
	// for _, deal := range deals {
	// 	if deal.Products != nil {
	// 		for _, product := range deal.Products {
	// 			ids = append(ids, product.Id)
	// 		}
	// 	}
	// }

	// products, err := c.GetProductAttachmentByIds(ctx, ids)
	// if err != nil {
	// 	return
	// }

	// mapProduct := make(map[uint64]*bdspropb.ProductAttachment)
	// for _, product := range products.Data {
	// 	mapProduct[product.Id] = product
	// }

	// for _, deal := range deals {
	// 	if deal.Products != nil {
	// 		for i, product := range deal.Products {
	// 			deal.Products[i] = mapProduct[product.Id]
	// 		}
	// 	}
	// }
}

func (c *BdsproClient) MapDealContractToPb(ctx context.Context, req *transactionpb.DealContractResponse) {
	if req.ProductId == 0 {
		return
	}

	product, err := c.GetProductAttachmentByIds(ctx, []uint64{req.ProductId})
	if err != nil {
		return
	}
	if product != nil && len(product.Data) > 0 {
		req.Product = product.Data[0]
	}
}
