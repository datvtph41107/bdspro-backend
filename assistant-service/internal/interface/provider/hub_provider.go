package provider

import (
	_dto "common/domain/dto"
	"context"

	hubpb "pb/types/hub"
)

// HubProvider - Interface để tương tác với Hub Service
type HubProvider interface {
	// SearchLocationV2 tìm kiếm ward/province theo từ khóa
	SearchLocationV2(ctx context.Context, keyword string, limit int32) ([]*hubpb.LocationSearchResultV2, error)
	InferAddressFromText(ctx context.Context, text string) (*_dto.AddressV3DTO, error)
}
