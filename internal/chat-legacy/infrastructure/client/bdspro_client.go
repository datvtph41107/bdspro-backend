package client

import (
	"chat/models"
	"context"
	"pb/clients"
	bdspropb "pb/types/bdspro"
	chatpb "pb/types/chat"
)

type BdsproClient struct {
	Client *clients.BdsproGrpcClient
}

func NewBdsproClient(rpcClient *clients.BdsproGrpcClient) *BdsproClient {
	return &BdsproClient{Client: rpcClient}
}

func (c *BdsproClient) MapProductPbByIds(ctx context.Context, message []*chatpb.Message) {
	productIdSet := make(map[uint64]struct{})
	for _, product := range message {
		if product.ExtraId != nil {
			productIdSet[*product.ExtraId] = struct{}{}
		}
	}

	productIds := make([]uint64, 0, len(productIdSet))
	for productId := range productIdSet {
		productIds = append(productIds, productId)
	}

	products, err := c.Client.Client.GetProductAttachmentByIds(ctx, &bdspropb.IdRequest{
		Ids: productIds,
	})
	if err != nil {
		return
	}

	productMap := make(map[uint64]*bdspropb.ProductAttachment)
	for _, product := range products.Data {
		productMap[product.Id] = product
	}

	for _, message := range message {
		if message.ContentType == string(models.ContentTypeProduct) &&
			message.ExtraId != nil {
			product, ok := productMap[*message.ExtraId]
			if ok {
				message.Product = product
			}
		}
	}
}
