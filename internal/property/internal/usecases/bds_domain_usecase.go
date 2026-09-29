package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	"context"
)

type BdsDomainUsecase struct {
	BdsDomainRepo repo.BdsDomainRepo
}

func NewBdsDomainUsecase(bdsDomainRepo repo.BdsDomainRepo) *BdsDomainUsecase {
	return &BdsDomainUsecase{
		BdsDomainRepo: bdsDomainRepo,
	}
}

// GetList lấy danh sách bất động sản với phân trang và bộ lọc
func (uc *BdsDomainUsecase) GetList(ctx context.Context, request *dto.BdsDomainListRequest) ([]domain.BDSDomain, int64, error) {
	// Validate pagination
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Size < 1 {
		request.Size = 20
	}
	if request.Size > 100 {
		request.Size = 100
	}

	return uc.BdsDomainRepo.GetList(ctx, request)
}

// GetByID lấy chi tiết bất động sản theo ID
func (uc *BdsDomainUsecase) GetByID(ctx context.Context, id uint64) (*domain.BDSDomain, error) {
	return uc.BdsDomainRepo.GetByID(ctx, id)
}

// GetByProductID lấy bất động sản theo Product ID
func (uc *BdsDomainUsecase) GetByProductID(ctx context.Context, productID uint64) (*domain.BDSDomain, error) {
	return uc.BdsDomainRepo.GetByProductID(ctx, productID)
}

// GetByAssetID lấy bất động sản theo Asset ID
func (uc *BdsDomainUsecase) GetByAssetID(ctx context.Context, assetID uint64) (*domain.BDSDomain, error) {
	return uc.BdsDomainRepo.GetByAssetID(ctx, assetID)
}

// CreateBdsDomainFromProduct tạo BDSDomain từ Product và HouseInfo
func (uc *BdsDomainUsecase) CreateBdsDomainFromProduct(ctx context.Context, product *domain.Product, houseInfo *domain.HouseInfo, profileID uint64) error {
	bds := &domain.BDSDomain{
		ProductID:  &product.ID,
		AssetID:    product.AssetId,
		ProvinceID: product.ProvinceID,
		WardID:     product.WardID,
		Address:    product.Address,
		Note:       product.Note,
		ProfileID:  &profileID,
	}

	// Copy diện tích từ Product
	if product.Area > 0 {
		area := product.Area
		bds.Area = &area
	}

	// Copy thông tin từ HouseInfo nếu có
	if houseInfo != nil {
		if houseInfo.FrontWidth != nil {
			frontWidth := float64(*houseInfo.FrontWidth)
			bds.FrontWidth = &frontWidth
		}
		if houseInfo.BackWidth != nil {
			backWidth := float64(*houseInfo.BackWidth)
			bds.BackWidth = &backWidth
		}
		if houseInfo.Width != nil {
			width := float64(*houseInfo.Width)
			bds.Width = &width
		}
		if houseInfo.Height != nil {
			height := float64(*houseInfo.Height)
			bds.Length = &height // Height trong HouseInfo là chiều dài
		}
		if houseInfo.Orientation != nil {
			bds.Orientation = houseInfo.Orientation
			bds.MainDirection = houseInfo.Orientation
		}
		if houseInfo.RoadWidth != nil {
			roadTypeStr := ""
			// Convert uint32 to string if needed, or use a mapping
			// For now, just set as empty, can be enhanced later
			bds.RoadType = &roadTypeStr
		}
	}

	return uc.BdsDomainRepo.Create(ctx, bds)
}
