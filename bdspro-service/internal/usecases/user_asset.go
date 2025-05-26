package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"errors"
	"log"
)

type UAssetUsecase struct {
	AssetRepo        repo.AssetRepo
	ProductUserRepo  repo.UProductRepo
	ProductUc        *UProductUsecase
	AssetHistoryRepo repo.UAssetHistoryRepo
}

func NewUAssetUsecase(assetRepo repo.AssetRepo,
	ProductRepo repo.UProductRepo,
	// ProductUc *UProductUsecase,
	AssetHistoryRepo repo.UAssetHistoryRepo,
) *UAssetUsecase {
	return &UAssetUsecase{
		AssetRepo:       assetRepo,
		ProductUserRepo: ProductRepo,
		// ProductUc:       ProductUc,
	}
}

// UC1: Xem danh sách tổng quan tất cả tài sản với phân trang và profileId
func (uc *UAssetUsecase) GetOwnerAssets(ctx context.Context, profileId uint64, limit, offset int) ([]domain.Asset, error) {
	return uc.AssetRepo.GetAssets(ctx, profileId, limit, offset)
}

// UC2: Thực hiện thao tác hàng loạt (cập nhật trạng thái, xóa đồng thời nhiều tài sản)
func (uc *UAssetUsecase) BulkAction(ctx context.Context, action string, ids []uint64) error {
	if len(ids) == 0 {
		return errors.New("no asset IDs provided for bulk action")
	}
	return uc.AssetRepo.BulkAction(ctx, action, ids)
}

func (s *UAssetUsecase) RequiredOwner(c context.Context, id uint64) (*domain.Asset, error) {
	asset, err := s.AssetRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// profileId := _jwt.GetProfileId(c)
	profileId := _utils.GetProfileIdWithContext(c)

	// if asset.OwnerUserID == nil {
	// 	return nil, &_routes.Except{
	// 		Code:    401,
	// 		Message: "Bạn không có quyền truy cập",
	// 	}
	// }
	log.Print("*asset.OwnerUserID != profileId", *asset.OwnerUserID != profileId)
	if asset.OwnerUserID == nil || *asset.OwnerUserID != profileId {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Bạn không có quyền truy cập",
		}
	}
	return asset, nil
}

// UC3: Chỉnh sửa thông tin tài sản trực tiếp từ danh sách
func (uc *UAssetUsecase) UpdateAsset(ctx context.Context, id uint64, assetData *domain.Asset) (*domain.Asset, error) {
	if id == 0 {
		return nil, errors.New("asset ID is required")
	}
	asset, err := uc.RequiredOwner(ctx, id)
	if err != nil {
		return nil, err
	}

	asset.ID = id
	asset.Name = assetData.Name
	asset.WardID = assetData.WardID
	asset.DistrictID = assetData.DistrictID
	asset.ProvinceID = assetData.ProvinceID
	// asset.Price = assetData.Price
	// asset.BuyedDate = assetData.BuyedDate
	// asset.DocStatus = assetData.DocStatus
	// asset.ProductID = assetData.ProductID
	// asset.OwnerUserID = assetData.OwnerUserID
	// asset.OwnerOrgID = assetData.OwnerOrgID
	// asset.ParentAssetID = assetData.ParentAssetID
	asset.SplitMergeStatus = assetData.SplitMergeStatus
	asset.PurchasePrice = assetData.PurchasePrice
	asset.PurchaseDate = assetData.PurchaseDate
	asset.LegalStatus = assetData.LegalStatus
	// asset.Archived = assetData.Archived
	asset.RentStatus = assetData.RentStatus
	asset.Area = assetData.Area
	asset.Address = assetData.Address
	asset.Description = assetData.Description

	err = uc.AssetRepo.Update(ctx, id, asset)
	if err != nil {
		return nil, err
	}

	afterData := _utils.StructToJSONString(struct{ Data *domain.Asset }{Data: asset})
	uc.AssetHistoryRepo.CreateHistory(ctx, enums.AssetHistory_Update, &assetData.ID, nil, afterData)

	return asset, nil
}

// UC4: Xóa nhanh tài sản từ danh sách
func (uc *UAssetUsecase) DeleteAsset(ctx context.Context, id uint64) error {
	if id == 0 {
		return errors.New("asset ID is required")
	}
	asset, err := uc.RequiredOwner(ctx, id)
	if err != nil {
		return err
	}

	if asset.RentStatus != enums.AssetStatusOwning {
		return errors.New("không thể xóa tài sản do trạng thái tài sản không cho phép")
	}
	return uc.AssetRepo.Delete(ctx, id)
}

// UC5: Lọc danh sách tài sản theo trạng thái, loại hình tài sản, người sở hữu, diện tích, ngày tạo
// func (uc *UAssetUsecase) FilterAssets(ctx context.Context, filters map[string]interface{}) ([]domain.Asset, error) {
// 	return uc.AssetRepo.FilterAssets(ctx, filters)
// }

// UC6: Tìm kiếm tài sản theo từ khóa hoặc mã tài sản
// func (uc *UAssetUsecase) SearchAssets(ctx context.Context, keyword string) ([]domain.Asset, error) {
// 	if keyword == "" {
// 		return nil, errors.New("search keyword is required")
// 	}
// 	filters := map[string]interface{}{
// 		"name LIKE ?": "%" + keyword + "%",
// 		"id LIKE ?":   "%" + keyword + "%",
// 	}
// 	return uc.AssetRepo.FilterAssets(ctx, filters)
// }

// UC7: Tạo mới tài sản bằng nút truy cập nhanh
func (uc *UAssetUsecase) CreateAsset(ctx context.Context, assetData *domain.Asset) (*domain.Asset, error) {
	if assetData.Name == "" || assetData.Area <= 0 || assetData.Address == "" {
		return nil, errors.New("missing required fields: name, area, or address")
	}
	id := ctx.Value("profileId").(uint64)
	assetData.OwnerUserID = &id
	assetData.Archived = false

	err := uc.AssetRepo.Create(ctx, assetData)
	return assetData, err
}

// UC8: Tạo tài sản từ sản phẩm có sẵn
func (uc *UAssetUsecase) CreateAssetFromProduct(ctx context.Context, productID uint64, assetData *domain.Asset) (*domain.Asset, error) {
	// if req.ProductId == nil {
	// 	return nil, fmt.Errorf("product_id is required")
	// }

	// // Có thể lấy thông tin sản phẩm và tự động điền (nếu cần)
	// product, err := s.productRepo.GetByID(*req.ProductId)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get product: %w", err)
	// }

	// // Tự động điền các trường cơ bản từ product (nếu cần)
	// if req.Name == "" {
	// 	req.Name = product.Name
	// }
	// if req.Area == 0 {
	// 	req.Area = product.Area
	// }

	// return s.saveAsset(ctx, req)

	// Lấy thông tin sản phẩm từ ProductService
	product, err := uc.ProductUserRepo.GetByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Liên kết thông tin sản phẩm vào tài sản
	assetData.ProductID = &product.ID
	assetData.Name = product.Name
	assetData.Description = product.Description
	assetData.Archived = false

	// Tạo tài sản
	err = uc.AssetRepo.Create(ctx, assetData)
	return assetData, err
}

// UC9: Tạo tài sản và sản phẩm mới
func (uc *UAssetUsecase) CreateAssetAndProduct(ctx context.Context, productData *domain.Product, assetData *domain.Asset) (*domain.Asset, error) {
	// Tạo sản phẩm mới
	product, err := uc.AssetRepo.CreateProduct(ctx, productData)
	if err != nil {
		return nil, err
	}

	// Liên kết sản phẩm mới vào tài sản
	assetData.ProductID = &product.ID
	assetData.Archived = false

	uc.ProductUc.CreateProduct(ctx, &dto.ProductSaveRequest{
		Name:        assetData.Name,
		Area:        assetData.Area,
		WardID:      assetData.WardID,
		ProvinceID:  assetData.ProvinceID,
		DistrictID:  assetData.DistrictID,
		Description: assetData.Description,
		PriceDTO: &domain.ProductPrice{
			SalePrice: &assetData.PurchasePrice,
		},
	})

	// Tạo tài sản
	err = uc.AssetRepo.Create(ctx, assetData)
	return assetData, err
}

// UC10: Chuyển sản phẩm thành tài sản
func (uc *UAssetUsecase) ConvertProductToAsset(ctx context.Context, productID uint64, assetData *domain.Asset) (*domain.Asset, error) {
	// Lấy thông tin sản phẩm từ ProductRepo
	product, err := uc.AssetRepo.GetProductByID(ctx, &productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Liên kết thông tin sản phẩm vào tài sản
	assetData.ProductID = &product.ID
	assetData.Name = product.Name
	assetData.Description = product.Description

	// Tạo tài sản
	err = uc.AssetRepo.Create(ctx, assetData)
	return assetData, err
}

// UC11: Xóa tài sản (ẩn vĩnh viễn)
// func (uc *AssetUsecase) DeleteAssetWithValidation(ctx context.Context, assetID uint64) error {
// 	if assetID == 0 {
// 		return errors.New("asset ID is required")
// 	}

// 	// Lấy thông tin tài sản
// 	asset, err := uc.AssetRepo.GetAssetDetails(ctx, assetID)
// 	if err != nil {
// 		return errors.New("asset not found")
// 	}

// 	// Kiểm tra sản phẩm liên kết
// 	if asset.ProductID != nil {
// 		// Lưu trữ sản phẩm liên kết
// 		if err := uc.ProductUserRepo.UpdateArchived(ctx, *asset.ProductID, true); err != nil {
// 			return errors.New("failed to archive linked product")
// 		}
// 	}

// 	// Cập nhật trạng thái tài sản thành "deleted"
// 	if err := uc.AssetRepo.Delete(ctx, assetID); err != nil {
// 		return errors.New("failed to delete asset")
// 	}

// 	return nil
// }

// UC12: Lưu trữ tài sản
func (uc *UAssetUsecase) ArchiveAsset(ctx context.Context, assetID uint64, archived bool) error {
	if assetID == 0 {
		return errors.New("asset ID is required")
	}

	// Lấy thông tin tài sản
	_, err := uc.AssetRepo.GetByID(assetID)
	if err != nil {
		return err // errors.New("asset not found")
	}

	if _, err = uc.RequiredOwner(ctx, assetID); err != nil {
		return err
	}

	// Cập nhật trạng thái tài sản thành "archived"
	if err := uc.AssetRepo.UpdateArchived(ctx, assetID, archived); err != nil {
		return errors.New("failed to archive asset")
	}

	return nil
}
