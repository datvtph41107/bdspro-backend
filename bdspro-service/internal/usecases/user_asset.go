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

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
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
		AssetRepo:        assetRepo,
		ProductUserRepo:  ProductRepo,
		AssetHistoryRepo: AssetHistoryRepo,
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
func (uc *UAssetUsecase) CreateAssetFromProduct(ctx *gin.Context, assetData *domain.Asset) (*domain.Asset, error) {
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
	product, err := uc.ProductUserRepo.GetByID(*assetData.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Liên kết thông tin sản phẩm vào tài sản
	assetData.ProductID = &product.ID
	id := _utils.GetProfileIdWithContext(ctx)
	assetData.OwnerUserID = &id
	assetData.Archived = false

	// Tạo tài sản
	err = uc.AssetRepo.Create(ctx, assetData)

	if assetData.ProductID != nil {
		uc.ProductUserRepo.UpdateAssetID(ctx, *assetData.ProductID, assetData.ID)
	}

	afterData := _utils.StructToJSONString(struct{ Data *domain.Asset }{Data: assetData})
	uc.AssetHistoryRepo.CreateHistory(ctx, enums.AssetHistory_New, &assetData.ID, nil, afterData)

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

// UC13: Gộp tài sản
func (uc *UAssetUsecase) MergeAsset(ctx context.Context, dto dto.MergeAssetDTO) (*domain.Asset, error) {
	// if len(dto.IDs) < 1 {
	// 	return nil, &_routes.Except{
	// 		Code:    401,
	// 		Message: "Bạn phải chọn tối thiểu 2 tài sản",
	// 	}
	// }

	profileId := _utils.GetProfileIdWithContext(ctx)
	entities, err := uc.AssetRepo.GetByIDsWithValidOwner(ctx, profileId, dto.IDs)
	if len(entities) != len(dto.IDs) {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Bạn không sở hữu tài sản hoặc tài sản đang trong giao dịch chưa hoàn tất",
		}
	}

	if err != nil {
		return nil, err
	}

	currentAsset := entities[0]
	sumArea := 0.0
	parentId := (*uint64)(nil)

	for _, e := range entities {
		sumArea += e.Area
		if e.ParentAssetID != nil {
			parentId = e.ParentAssetID
		}
		if e.ParentAssetID == nil {
			return nil, &_routes.Except{
				Code:    400,
				Message: "Chỉ gộp tài sản con",
			}
		}
		// if e.WardID != currentAsset.WardID ||
		// 	e.DistrictID != currentAsset.DistrictID ||
		// 	e.ProvinceID != currentAsset.ProvinceID {
		// 	return nil, &_routes.Except{
		// 		Code:    400,
		// 		Message: "Danh sách tài sản không liền kề nhau",
		// 	}
		// }
	}

	if parentId != nil {
		parentAsset, err := uc.AssetRepo.GetByID(*parentId)
		if err != nil {
			return nil, err
		}
		beforeData := _utils.StructToJSONString(parentAsset)
		parentAsset.Area = parentAsset.Area + sumArea
		if err := uc.AssetRepo.Update(ctx, parentAsset.ID, parentAsset); err != nil {
			return nil, err
		}

		if err = uc.AssetRepo.DeleteBatch(ctx, dto.IDs); err != nil {
			return nil, err
		}
		// uc.AssetHistoryRepo.CreateHistory(ctx, enums.AssetHistory_Merge, parentId, dto.IDs)

		afterData := _utils.StructToJSONString(struct {
			Parent *domain.Asset
			Childs []domain.Asset
		}{
			Parent: parentAsset,
			Childs: entities,
		})
		uc.AssetHistoryRepo.CreateHistory(ctx, enums.AssetHistory_Merge, parentId, beforeData, afterData)
		return parentAsset, nil
	}

	// todo: đoạn này có thể không chạy đến do chỉ gộp con chứ không gộp cha

	nAsset := &domain.Asset{}
	copier.Copy(nAsset, &dto)

	// profileId := _jwt.GetProfileId(ctx)
	nAsset.OwnerUserID = &profileId
	nAsset.WardID = currentAsset.WardID
	nAsset.DistrictID = currentAsset.DistrictID
	nAsset.ProvinceID = currentAsset.ProvinceID
	nAsset.Archived = false
	nAsset.Area = sumArea
	nAsset.RentStatus = 1

	// Cập nhật trạng thái tài sản thành "archived"
	if err := uc.AssetRepo.Create(ctx, nAsset); err != nil {
		return nil, err
	}

	if err = uc.AssetRepo.DeleteBatch(ctx, dto.IDs); err != nil {
		return nil, err
	}

	// todo

	return &currentAsset, err
}

// UC13: Tách tài sản
func (uc *UAssetUsecase) SplitAsset(ctx context.Context, dto dto.SplitAssetDTO) ([]domain.Asset, error) {
	currentAsset, err := uc.RequiredOwner(ctx, dto.ID)
	if err != nil {
		return nil, err
	}

	// todo: tài sản liên kết nhiều sp
	numProductLink, err := uc.ProductUserRepo.CountByAssetID(currentAsset.ID)
	if err != nil {
		return nil, err
	}
	if numProductLink > 1 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không thể tách do tài sản đang liên kết với nhiều sản phẩm",
		}
	}

	// currentAsset := entities[0]
	sumArea := 0.0
	for _, e := range dto.Datas {
		sumArea += e.Area
		if e.Area < 1 {
			return nil, &_routes.Except{
				Code:    400,
				Message: "Diện tích không hợp lệ",
			}
		}
	}

	if sumArea >= currentAsset.Area {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Diện tích không hợp lệ",
		}
	}

	profileId := _utils.GetProfileIdWithContext(ctx)
	for i := range dto.Datas {
		nAsset := &dto.Datas[i]
		nAsset.OwnerUserID = &profileId
		nAsset.WardID = currentAsset.WardID
		nAsset.DistrictID = currentAsset.DistrictID
		nAsset.ProvinceID = currentAsset.ProvinceID
		nAsset.Archived = false
		nAsset.RentStatus = 1
		nAsset.ParentAssetID = &dto.ID
		// nAsset.Area = sumArea
	}

	// Cập nhật trạng thái tài sản thành "archived"
	if err := uc.AssetRepo.CreateBatch(ctx, dto.Datas); err != nil {
		return nil, err
	}

	beforeData := _utils.StructToJSONString(currentAsset)

	currentAsset.Area = currentAsset.Area - sumArea
	if err := uc.AssetRepo.Update(ctx, currentAsset.ID, currentAsset); err != nil {
		return nil, err
	}

	afterData := _utils.StructToJSONString(struct {
		Parent *domain.Asset
		Childs []domain.Asset
	}{
		Parent: currentAsset,
		Childs: dto.Datas,
	})

	// var ids []uint64
	// for _, r := range dto.Datas {
	// 	ids = append(ids, r.ID)
	// }
	// uc.AssetHistoryRepo.CreateHistory(ctx, enums.AssetHistory_Slit, &dto.ID, ids)
	uc.AssetHistoryRepo.CreateHistory(ctx, enums.AssetHistory_Slit, &dto.ID, beforeData, afterData)

	// todo

	return dto.Datas, nil
}

func (uc *UAssetUsecase) SearchAssetMe(ctx context.Context, dto dto.AssetSearchDTO) (_routes.ResponseDTO, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	assets, total, err := uc.AssetRepo.SearchOwner(ctx, profileId, dto)
	if err != nil {
		return _routes.ResponseDTO{}, err
	}
	return _routes.ResponseDTO{
		Code:          0,
		Data:          assets,
		TotalElements: &total,
	}, nil
}

// UC6: Tìm kiếm tài sản theo từ khóa hoặc mã tài sản
func (uc *UAssetUsecase) SearchAssetShare(ctx context.Context, dto dto.AssetSearchDTO) (_routes.ResponseDTO, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	results, total, err := uc.AssetRepo.SearchShare(ctx, profileId, dto)
	if err != nil {
		return _routes.ResponseDTO{}, err
	}
	return _routes.ResponseDTO{
		Data:          results,
		TotalElements: &total,
	}, nil
}

// UC13: Tách tài sản
func (uc *UAssetUsecase) History(ctx context.Context, assetID uint64, dto dto.AssetHistoryDTO) (_routes.ResponseDTO, error) {
	// profileId := _jwt.GetProfileId(ctx)
	results, total, err := uc.AssetRepo.History(ctx, 0, assetID, &dto)
	if err != nil {
		return _routes.ResponseDTO{}, err
	}
	return _routes.ResponseDTO{
		Code:          0,
		Data:          results,
		TotalElements: &total,
	}, nil
}

// UC14: Lịch sử tài sản
func (uc *UAssetUsecase) HistoryOwner(ctx context.Context, dto dto.AssetHistoryDTO) (_routes.ResponseDTO, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	results, total, err := uc.AssetRepo.History(ctx, profileId, 0, &dto)
	if err != nil {
		return _routes.ResponseDTO{}, err
	}
	return _routes.ResponseDTO{
		Code:          0,
		Data:          results,
		TotalElements: &total,
	}, nil
}

func (uc *UAssetUsecase) GetAssetsWithChildCount(ctx context.Context, searchDto dto.AssetSearchDTO) ([]dto.AssetWithChildCountDTO, int64, error) {
	// Validate input if needed
	profileId := _utils.GetProfileIdWithContext(ctx)

	// Call repository to get data
	assets, total, err := uc.AssetRepo.GetAssetsWithChildCount(ctx, profileId, searchDto)
	if err != nil {
		return nil, 0, err
	}

	return assets, total, nil
}
