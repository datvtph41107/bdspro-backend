package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
)

// BdsDomainRepo interface cho repository BDSDomain
type BdsDomainRepo interface {
	// Create tạo mới bất động sản
	Create(ctx context.Context, bds *domain.BDSDomain) error

	// GetList lấy danh sách bất động sản với phân trang và bộ lọc
	GetList(ctx context.Context, request *dto.BdsDomainListRequest) ([]domain.BDSDomain, int64, error)

	// GetByID lấy chi tiết bất động sản theo ID
	GetByID(ctx context.Context, id uint64) (*domain.BDSDomain, error)

	// GetByProductID lấy bất động sản theo Product ID
	GetByProductID(ctx context.Context, productID uint64) (*domain.BDSDomain, error)

	// GetByAssetID lấy bất động sản theo Asset ID
	GetByAssetID(ctx context.Context, assetID uint64) (*domain.BDSDomain, error)
}
