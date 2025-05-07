package repo

import (
	"bdspro/enums"
	"bdspro/v2/internal/domain"
	"bdspro/v2/internal/dto"
	"common/case/crud3"
	_utils "common/utils"
	"context"
)

// ProductUserRepo định nghĩa các phương thức mà GormUserProductRepo cung cấp
type ProductUserRepo interface {
	crud3.CrudRepo[domain.Product]
	UpdateSaleStatus(ctx context.Context, id uint64, status enums.ProductStatus) error

	// Cập nhật trạng thái cho thuê
	UpdateRentStatus(ctx context.Context, id uint64, status enums.ProductStatus) error

	// Cập nhật trạng thái hiển thị bán
	UpdateSaleVisibility(ctx context.Context, id uint64, vis enums.VisibilityStatus) error

	// Cập nhật trạng thái hiển thị cho thuê
	UpdateRentVisibility(ctx context.Context, id uint64, vis enums.VisibilityStatus) error

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
	SearchProductWithOwner(profileId uint64, dto dto.ProductSearchRequest, pageable _utils.Page) ([]domain.Product, int64, error)

	// Tìm kiếm sản phẩm trên thị trường
	SearchProductMarket(profileId uint64, dto dto.ProductSearchRequest, pageable _utils.Page) ([]domain.Product, int64, error)

	// Tìm kiếm thông minh sản phẩm
	SmartSearchProduct(profileId uint64, dto dto.TotalSearchParser, pageable _utils.Page) ([]domain.Product, int64, error)

	// Phát hiện dữ liệu sản phẩm
	DetectProductData(profileId uint64, input string) ([]domain.Product, int64, error)

	// Cập nhật trạng thái lưu trữ sản phẩm
	UpdateArchived(ctx context.Context, productId uint64, archived bool) error

	// Tạo mới thông tin đặt cọc
	NewDeposite(ctx context.Context, dto *domain.Deposite) error
}
