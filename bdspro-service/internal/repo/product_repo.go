package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"common/case/crud3"
	_enum "common/domain/enum"
	"context"
	"time"

	"gorm.io/gorm"
)

type ProductTimestamp struct {
	ID        uint64
	UpdatedAt time.Time // millisecond
}

// ProductRepo định nghĩa các phương thức mà GormUserProductRepo cung cấp
type ProductRepo interface {
	crud3.CrudRepo[domain.Product]

	GetAllTimestampsForSync(ctx context.Context, profileID uint64, lastID uint64, limit int) ([]domain.ProductUser, error)

	BatchGetTimestamps(ctx context.Context, ids []uint64) (map[uint64]int64, error)

	UpdateViewStats(ctx context.Context, stats *domain.ProductStats, oldLastTime int64) error

	GetPrivateByProductID(ctx context.Context, productID uint64) (*domain.ProductPrivate, error)

	// LastPriceId Product Owner
	GetProductWithCurrentPrice(ctx context.Context, productID uint64) (*domain.Product, error)

	// OriginProfile -> origin_profile, profileId -> user
	GetPriceProductCurrentUser(ctx context.Context, productID uint64, originProfileID uint64) (*dto.PriceRow, error)

	QueryDomain(query *gorm.DB, dto *dto.ProductSearchRequest) *gorm.DB

	UpdateLastPrice(ctx context.Context, productID uint64, priceID uint64) error

	UpdateSourceFields(ctx context.Context, productID uint64, fields map[string]interface{}) error

	UpdateSaleStatus(ctx context.Context, id uint64, status enums.EProductSaleStatus) error

	// Cập nhật trạng thái cho thuê
	UpdateRentStatus(ctx context.Context, id uint64, status enums.EProductSaleStatus) error

	// Cập nhật trạng thái hiển thị bán
	UpdateSaleVisibility(ctx context.Context, id uint64, vis enums.EVisibility) error

	// Cập nhật trạng thái hiển thị cho thuê
	UpdateRentVisibility(ctx context.Context, id uint64, vis enums.EVisibility) error

	// Cập nhật diện tích sản phẩm
	UpdateArea(ctx context.Context, id uint64, area float64) error

	// Cập nhật thông tin sản phẩm
	UpdateProduct(c context.Context, id uint64, entity *domain.Product) error
	UpdateSaleTransactionID(c context.Context, id uint64, transactionID *uint64) error
	UpdateRentTransactionID(c context.Context, id uint64, transactionID *uint64) error
	Save(c context.Context, entity *domain.Product) error

	// // Lấy sản phẩm theo ID
	// GetByID(ctx context.Context, id uint64) (*domain.Product, error)

	// Lấy sản phẩm theo ID (phiên bản khác)
	// GetByID2(ctx context.Context, productID uint64) (*domain.Product, error)

	// Cập nhật danh sách tiện ích của sản phẩm
	UpdateProductAmenityIds(ctx context.Context, productID uint64, amenityIDs []uint64) error

	// Tìm kiếm sản phẩm theo chủ sở hữu
	// deprecated
	// SearchProductWithOwner(profileId *uint64, payload dto.ProductSearchRequest) ([]dto.ProductListResponse, int64, error)

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
	// NewDeposite(ctx context.Context, dto *domain.Deposite) error
	CountByAssetID(assetID uint64) (int64, error)
	UpdateAssetID(c context.Context, id uint64, assetID uint64) error
	UpdateImageID(c context.Context, id uint64, imageID *uint64) error
	GetByIDContext(ctx context.Context, id uint64) (*domain.Product, error)
	// CancelDeposite(ctx context.Context, id uint64) error

	// SearchProductByOwner(ownerID *uint64, ownerType enums.EOwnerOf, dto *dto.ProductSearchRequest) ([]domain.Product, int64, error)
	SearchExport(ownerID *uint64, ownerType enums.EOwnerOf, dto *dto.ProductSearchRequest) ([]domain.Product, error)
	SmartSearch(ownerID *uint64, ownerType enums.EOwnerOf, dto *dto.TotalSearchParser) ([]domain.Product, int64, error)
	Archived(ctx context.Context, productId *uint64, archived bool) error
	SetRentStatus(ctx context.Context, id *uint64, status enums.EProductSaleStatus) error
	SetSaleStatus(ctx context.Context, id *uint64, status enums.EProductSaleStatus) error
	SetRentVisibility(ctx context.Context, id *uint64, vis enums.EVisibility) error
	SetSaleVisibility(ctx context.Context, id *uint64, vis enums.EVisibility) error
	GetProductByID(ctx context.Context, productID *uint64) (*domain.Product, error)
	SearchByProductAccess(ctx context.Context, targetID *uint64, targetType enums.EOwnerOf, dto *dto.ProductSearchRequest) ([]domain.Product, int64, error)
	GetDetailByIds(ctx context.Context, ids []uint64) ([]domain.Product, error)
	SyncData(ctx context.Context) error
	GetGroupProducts(ctx context.Context, req *dto.GroupProductSearchRequest) ([]domain.Product, int64, error)
	GetUserProducts(ctx context.Context, req *dto.UserProductSearchRequest) ([]domain.Product, int64, error)

	GetNextProductId(c context.Context, currentId *uint64) (uint64, error)
	GetProducts(c context.Context, payload dto.ProductSearchRequest) ([]dto.ProductListResponse, int64, error)
	GetProductsByAdmin(c context.Context, dto *dto.ProductSearchRequest) ([]domain.Product, int64, error)

	// Đếm số sản phẩm theo ownerOf và ownerId
	CountByOwner(ctx context.Context, ownerOf _enum.EOwnerOf, ownerId uint64) (uint32, error)

	GetRefreshSyncIds(c context.Context, req *dto.CheckVersionSyncRequest) ([]uint64, error)
	GetUpdatedAt(c context.Context, id uint64) (int64, error)

	// Cập nhật ghi chú sản phẩm
	UpdateNote(ctx context.Context, id uint64, note string) error

	// Đếm tổng số product hiện tại (không bị xóa)
	CountCurrent(ctx context.Context) (int64, error)

	// Lấy ngẫu nhiên 5 sản phẩm
	GetRandomProducts(ctx context.Context, limit int) ([]domain.Product, error)

	// Lấy thống kê tổng hợp sản phẩm của user
	GetProductSummary(ctx context.Context, userID uint64, searchRequest dto.ProductSearchRequest) (*dto.ProductSummaryDTO, error)

	// Cập nhật PostCount cho product
	UpdatePostCount(ctx context.Context, productID uint64, count uint32) error

	// Cập nhật ShareCount cho product
	UpdateShareCount(ctx context.Context, productID uint64, count uint32) error

	// Cập nhật ScheduleCount cho product
	UpdateScheduleCount(ctx context.Context, productID uint64, count uint32) error

	// Cập nhật DealCount cho product
	UpdateDealCount(ctx context.Context, productID uint64, count uint32) error

	// Cập nhật AssetCount cho product
	UpdateAssetCount(ctx context.Context, productID uint64, count uint32) error

	// Cập nhật ContactCount cho product
	UpdateContactCount(ctx context.Context, productID uint64, count uint32) error

	// GetProductNavigation lấy previousId và nextId của product
	GetProductNavigation(ctx context.Context, productID uint64) (previousID uint64, nextID uint64, err error)

	// GetAmenityIDsByProductID lấy danh sách amenity IDs của product
	GetAmenityIDsByProductID(ctx context.Context, productID uint64) ([]uint64, error)
}
