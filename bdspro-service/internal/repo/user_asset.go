package repo

import (
	"bdspro/internal/domain"
	"common/case/crud3"
	"context"
)

type AssetRepo interface {
	crud3.CrudRepo[domain.Asset]
	// Create(ctx context.Context, asset *domain.Asset) (*domain.Asset, error)
	GetAssets(ctx context.Context, profileId uint64, limit, offset int) ([]domain.Asset, error)
	// FilterAssets(ctx context.Context, filter map[string]interface{}) ([]domain.Asset, error)
	BulkAction(ctx context.Context, action string, ids []uint64) error
	// Update(ctx context.Context, id uint64, asset *domain.Asset) error
	UpdateArchived(ctx context.Context, id uint64, archived bool) error
	// Delete(ctx context.Context, id uint64) error
	// GetAssetDetails(ctx context.Context, id uint64) (*domain.Asset, error)
	ShareAsset(ctx context.Context, id uint64, userID uint64) error
	SplitAsset(ctx context.Context, id uint64, splitData map[string]interface{}) ([]domain.Asset, error)
	MergeAssets(ctx context.Context, ids []uint64) (*domain.Asset, error)
	// Các phương thức mới
	GetProductByID(ctx context.Context, productID *uint64) (*domain.Product, error)
	CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
}
