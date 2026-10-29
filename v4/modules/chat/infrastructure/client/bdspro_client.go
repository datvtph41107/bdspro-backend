package client

import (
	"context"
	"mdn4/modules/chat/models"
	bdspropb "pb/types/bdspro"
	chatpb "pb/types/chat"
)

type ProductAttachmentsFunc func(context.Context, []uint64) (map[uint64]*bdspropb.ProductAttachment, error)

type BdsproClient struct {
	lookup ProductAttachmentsFunc
}

func NewBdsproClient(lookup ProductAttachmentsFunc) *BdsproClient {
	return &BdsproClient{lookup: lookup}
}

func (c *BdsproClient) MapProductPbByIds(ctx context.Context, messages []*chatpb.Message) {
	if c == nil || c.lookup == nil {
		return
	}
	ids := make([]uint64, 0)
	seen := make(map[uint64]struct{})
	for _, message := range messages {
		if message == nil || message.ContentType != string(models.ContentTypeProduct) || message.ExtraId == nil {
			continue
		}
		if _, exists := seen[*message.ExtraId]; !exists {
			seen[*message.ExtraId] = struct{}{}
			ids = append(ids, *message.ExtraId)
		}
	}
	if len(ids) == 0 {
		return
	}
	products, err := c.lookup(ctx, ids)
	if err != nil {
		return
	}
	for _, message := range messages {
		if message == nil || message.ContentType != string(models.ContentTypeProduct) || message.ExtraId == nil {
			continue
		}
		message.Product = products[*message.ExtraId]
	}
}
