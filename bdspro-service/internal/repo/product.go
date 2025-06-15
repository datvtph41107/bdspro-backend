package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"common/case/crud3"
	"context"
)

// ProductRepo định nghĩa các phương thức mà GormUserProductRepo cung cấp
type ProductRepo interface {
	crud3.CrudRepo[domain.Product]
	UpdateSaleStatus(ctx context.Context, id uint64, status enums.EProductStatus) error

	// Cập nhật trạng thái cho thuê
	UpdateRentStatus(ctx context.Context, id uint64, status enums.EProductStatus) error

	// Cập nhật trạng thái hiển thị bán
	UpdateSaleVisibility(ctx context.Context, id uint64, vis enums.EVisibility) error

	// Cập nhật trạng thái hiển thị cho thuê
	UpdateRentVisibility(ctx context.Context, id uint64, vis enums.EVisibility) error

	// Cập nhật diện tích sản phẩm
	UpdateArea(ctx context.Context, id uint64, area float64) error

	// Cập nhật thông tin sản phẩm
	UpdateProduct(c context.Context, id uint64, entity *domain.Product) error

	// // Lấy sản phẩm theo ID
	// GetByID(ctx context.Context, id uint64) (*domain.Product, error)

	// Lấy sản phẩm theo ID (phiên bản khác)
	GetByID2(ctx context.Context, productID uint64) (*domain.Product, error)

	// Cập nhật danh sách tiện ích của sản phẩm
	UpdateProductAmenityIds(productID uint64, amenityIDs []uint64) error

	// Tìm kiếm sản phẩm theo chủ sở hữu
	SearchProductWithOwner(profileId *uint64, dto dto.ProductSearchRequest) ([]domain.Product, int64, error)

	// Tìm kiếm sản phẩm để export excel
	SearchToExport(profileId uint64, dto dto.ProductSearchRequest) ([]domain.Product, error)

	// Tìm kiếm sản phẩm trên thị trường
	SearchProductMarket(profileId uint64, dto dto.ProductSearchRequest) ([]domain.ProductMarket, int64, error)

	// Tìm kiếm thông minh sản phẩm
	// SmartSearchProduct(profileId uint64, dto dto.TotalSearchParser) ([]domain.Product, int64, error)

	// Phát hiện dữ liệu sản phẩm
	DetectProductData(profileId uint64, input string) ([]domain.Product, int64, error)

	// Cập nhật trạng thái lưu trữ sản phẩm
	// UpdateArchived(ctx context.Context, productId uint64, archived bool) error

	// Tạo mới thông tin đặt cọc
	NewDeposite(ctx context.Context, dto *domain.Deposite) error
	CountByAssetID(assetID uint64) (int64, error)
	UpdateAssetID(c context.Context, id uint64, assetID uint64) error
	GetByIDContext(ctx context.Context, id uint64) (*domain.Product, error)
	CancelDeposite(ctx context.Context, id uint64) error

	// SearchProductByOwner(ownerID *uint64, ownerType enums.EOwnerType, dto *dto.ProductSearchRequest) ([]domain.Product, int64, error)
	SearchExport(ownerID *uint64, ownerType enums.EOwnerType, dto *dto.ProductSearchRequest) ([]domain.Product, error)
	SmartSearch(ownerID *uint64, ownerType enums.EOwnerType, dto *dto.TotalSearchParser) ([]domain.Product, int64, error)
	Archived(ctx context.Context, productId *uint64, archived bool) error
	SetRentStatus(ctx context.Context, id *uint64, status enums.EProductStatus) error
	SetSaleStatus(ctx context.Context, id *uint64, status enums.EProductStatus) error
	SetRentVisibility(ctx context.Context, id *uint64, vis enums.EVisibility) error
	SetSaleVisibility(ctx context.Context, id *uint64, vis enums.EVisibility) error
	GetProductByID(ctx context.Context, id *uint64) (*domain.Product, error)
	SearchByProductAccess(ctx context.Context, targetID *uint64, targetType enums.EOwnerType, dto *dto.ProductSearchRequest) ([]domain.Product, int64, error)
}
