package provider

import (
	"context"
	bdspropb "pb/types/bdspro"
)

type BdsproProvider interface {
	GetDealById(ctx context.Context, dealId uint64) (*bdspropb.GroupDeal, error)
	GetDealsByIds(ctx context.Context, dealIds []uint64) ([]*bdspropb.GroupDeal, error)
	GetProductAttachmentByIds(ctx context.Context, ids []uint64) (*bdspropb.ProductAttachmentResponse, error)
	GetDealContractsByIds(ctx context.Context, dealContractIds []uint64) ([]*bdspropb.DealContractItem, error)
	UpdateScheduleCount(ctx context.Context, productID uint64, count uint32) error
	UpdateContactCount(ctx context.Context, productID uint64, count uint32) error
}
