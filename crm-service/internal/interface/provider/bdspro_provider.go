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
	SearchProjects(ctx context.Context, text string, page uint32, size uint32) ([]*bdspropb.Project, int64, error)
	SearchRegions(ctx context.Context, text string, page uint32, size uint32) ([]*bdspropb.Region, int64, error)
	GetRegionByID(ctx context.Context, id uint64) (*bdspropb.Region, error)
	ListAreaRegions(ctx context.Context) ([]*bdspropb.AreaRegion, error)
	GetAreaRegionByID(ctx context.Context, id uint64) (*bdspropb.AreaRegion, error)
}