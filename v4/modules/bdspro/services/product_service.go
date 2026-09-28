package services

import (
	"context"

	bdspropb "pb/types/bdspro"

	"userv4/modules/bdspro/dto"
)

// ProductService owns product list/detail/create/delete/timestamps flows for modules/bdspro.
type ProductService interface {
	ProductMe(ctx context.Context, req dto.SyncListRequest) (*bdspropb.SearchResponse, error)
	// GetDetail returns product payload and sync timestamp (ms) for HTTP envelope root.
	GetDetail(ctx context.Context, req dto.SyncListRequest) (*bdspropb.ProductResponse, int64, error)
	ProductNew(ctx context.Context, req *bdspropb.ProductSaveRequest) (*bdspropb.Response, error)
	DeleteProduct(ctx context.Context, id uint64) error
	GetTimestamps(ctx context.Context, id uint64) (*bdspropb.ProductTimestampsResponse, error)
}
