package repo

import (
	"bdspro/internal/domain"
	"context"
)

// ProductAssetRepo interface cho repository ProductAsset
type ProductAssetRepo interface {
	// Link liên kết product với asset
	Link(ctx context.Context, productID, assetID uint64) error

	// Unlink hủy liên kết product với asset
	Unlink(ctx context.Context, productID, assetID uint64) error

	// GetByProductID lấy danh sách asset theo product ID
	GetByProductID(ctx context.Context, productID uint64) ([]*domain.ProductAsset, error)

	// GetByAssetID lấy danh sách product theo asset ID
	GetByAssetID(ctx context.Context, assetID uint64) ([]*domain.ProductAsset, error)

	// CheckExists kiểm tra liên kết đã tồn tại chưa
	CheckExists(ctx context.Context, productID, assetID uint64) (bool, error)

	// GetAssetIDsByProductID lấy danh sách asset ID theo product ID
	GetAssetIDsByProductID(ctx context.Context, productID uint64) ([]uint64, error)

	// GetProductIDsByAssetID lấy danh sách product ID theo asset ID
	GetProductIDsByAssetID(ctx context.Context, assetID uint64) ([]uint64, error)

	// CountAssetsByProductID đếm số asset liên kết với product
	CountAssetsByProductID(ctx context.Context, productID uint64) (int64, error)
}
