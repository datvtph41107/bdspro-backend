package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"common/case/crud3"
	_enum "common/domain/enum"
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
	GetAssetByProductID(ctx context.Context, id *uint64) (*domain.Asset, error)
	SplitAsset(ctx context.Context, id uint64, splitData map[string]interface{}) ([]domain.Asset, error)
	MergeAssets(ctx context.Context, ids []uint64) (*domain.Asset, error)
	// Các phương thức mới
	GetProductByID(ctx context.Context, productID *uint64) (*domain.Product, error)
	CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetByIDsWithValidOwner(ctx context.Context, profileId uint64, ids []uint64) ([]domain.Asset, error)
	DeleteBatch(c context.Context, ids []uint64) error
	CreateBatch(c context.Context, entity []domain.Asset) error
	// History(
	// 	c context.Context,
	// 	ownerID *uint64,
	// 	ownerType uint32,
	// 	assetId uint64,
	// 	dto *dto.AssetHistoryDTO) (*[]domain.AssetHistory, int64, error)
	SearchShare(ctx context.Context, ownerID uint64, ownerType uint32, dto dto.AssetSearchRequest) ([]domain.AssetList, int64, error)
	SearchOwner(ctx context.Context, ownerID uint64, ownerType uint32, dto dto.AssetSearchRequest) ([]domain.AssetList, int64, error)
	GetAssetsWithChildCount(ctx context.Context, ownerID uint64, dto dto.AssetSearchRequest) ([]dto.AssetWithChildCountDTO, int64, error)
	GetAssetPublish(profileId uint64) ([]domain.Asset, int64, error)
	SearchByAssetShare(ctx context.Context, targetID *uint64, targetType enums.EOwnerOf, dto *dto.ProductSearchRequest) ([]domain.Product, int64, error)

	ListAssets(ctx context.Context, dto *dto.AssetSearchRequest) ([]domain.Asset, int64, error)

	// Đếm số tài sản theo ownerOf và ownerId
	CountByOwner(ctx context.Context, ownerOf _enum.EOwnerOf, ownerId uint64) (uint32, error)

	// Đếm tổng số asset hiện tại (không bị xóa)
	CountCurrent(ctx context.Context) (int64, error)

	// GetAssetsByProductID lấy danh sách assets liên kết với product (thông tin cơ bản)
	GetAssetsByProductID(ctx context.Context, productID uint64) ([]domain.Asset, error)
}
