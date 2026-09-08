package client

import (
	"context"
	"pb/clients"
	bdspropb "pb/types/bdspro"
	organizationpb "pb/types/organization"
)

// type BdsproClient interface {
// 	GetProductByIds(ctx context.Context, in *bdspropb.IdRequest) (*bdspropb.ProductInternalResponse, error)
// }

// @bind: organization/internal/interface.BdsproClient
type BdsproClient struct {
	*clients.BdsproGrpcClient
}

func NewBdsproClient(rpcClient *clients.BdsproGrpcClient) *BdsproClient {
	return &BdsproClient{BdsproGrpcClient: rpcClient}
}

func (c *BdsproClient) MapProductToDealPb(ctx context.Context, deals []*organizationpb.GroupDeal) {
	ids := make([]uint64, len(deals))
	for _, deal := range deals {
		if deal.Products != nil {
			for _, product := range deal.Products {
				ids = append(ids, product.Id)
			}
		}
	}

	products, err := c.GetProductAttachmentByIds(ctx, ids)
	if err != nil {
		return
	}

	mapProduct := make(map[uint64]*bdspropb.ProductAttachment)
	for _, product := range products.Data {
		mapProduct[product.Id] = product
	}

	for _, deal := range deals {
		if deal.Products != nil {
			for i, product := range deal.Products {
				deal.Products[i] = mapProduct[product.Id]
			}
		}
	}
}

// func (c *BdsproClient) GetProductAttachmentByIds(ctx context.Context, ids []uint64) ([]*bdspropb.ProductAttachment, error) {
// 	products, err := c.GetProductByIds(ctx, ids)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return products.Data, nil
// }

// func (c *BdsproClient) MapDealContractToPb(ctx context.Context, req *organizationpb.DealContractResponse) {
// 	if req.ProductId == 0 {
// 		return
// 	}

// 	product, err := c.GetProductAttachmentByIds(ctx, []uint64{req.ProductId})
// 	if err != nil {
// 		return
// 	}
// 	if product != nil && len(product.Data) > 0 {
// 		req.Product = product.Data[0]
// 	}
// }

// GetCountByOwner gọi sang BdsproInternalService để lấy số lượng sản phẩm, tài sản, bài viết và dự án
func (c *BdsproClient) GetCountByOwner(ctx context.Context, req *bdspropb.GetCountByOwnerRequest) (*bdspropb.GetCountByOwnerResponse, error) {
	return c.BdsproGrpcClient.GetCountByOwner(ctx, req)
}
