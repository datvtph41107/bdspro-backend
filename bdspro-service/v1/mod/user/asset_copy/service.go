package u_asset

// import (
// 	"bdspro/internal/enums"
// 	history "bdspro/v1/mod/global/history/asset"
// 	u_product "bdspro/v1/mod/user/product"
// 	u_product_dto "bdspro/v1/mod/user/product/dto"
// 	"bdspro/v1/models"
// 	"common/case/rest"
// 	_jwt "common/jwt"
// 	_routes "common/routes"
// 	_utils "common/utils"
// 	"errors"
// 	"log"

// 	"github.com/gin-gonic/gin"
// 	"github.com/jinzhu/copier"
// )

// // type AssetService struct {
// // 	crud.BaseService[models.Asset]
// // }

// // func //NewAssetService(repo *AssetRepo) *AssetService {
// // 	return &AssetService{
// // 		BaseService: crud.BaseService[models.Asset]{Repo: RestRepo.BaseRepo},
// // 	}
// // }

// type AssetService struct {
// 	rest.RestService[models.Asset, SearchDTO]
// 	AssetRepo        *AssetRepo
// 	ProductUserRepo  *u_product.UserProductRepo
// 	ProductUc        *u_product.ProductService
// 	AssetHistoryRepo *history.AssetHistoryRepo
// }

// func //NewAssetService(assetRepo *AssetRepo,
// 	ProductRepo *u_product.UserProductRepo,
// 	ProductUc *u_product.ProductService,
// 	AssetHistoryRepo *history.AssetHistoryRepo) *AssetService {
// 	return &AssetService{
// 		RestService: rest.RestService[models.Asset, SearchDTO]{
// 			RestRepo: assetRepo,
// 		},
// 		AssetRepo:        assetRepo,
// 		ProductUserRepo:  ProductRepo,
// 		ProductUc:        ProductUc,
// 		AssetHistoryRepo: AssetHistoryRepo,
// 	}
// }

// // UC1: Xem danh sách tổng quan tất cả tài sản với phân trang và profileId
// func (uc *AssetService) GetOwnerAssets(ctx *gin.Context, profileId uint64, limit, offset int) ([]models.Asset, error) {
// 	return uc.AssetRepo.GetAssets(ctx, profileId, limit, offset)
// }

// // UC2: Thực hiện thao tác hàng loạt (cập nhật trạng thái, xóa đồng thời nhiều tài sản)
// func (uc *AssetService) BulkAction(ctx *gin.Context, action string, ids []uint64) error {
// 	if len(ids) == 0 {
// 		return errors.New("no asset IDs provided for bulk action")
// 	}
// 	return uc.AssetRepo.BulkAction(ctx, action, ids)
// }

// func (s *AssetService) RequiredOwner(c *gin.Context, id uint64) (*rest.RestModel[models.Asset], error) {
// 	asset, err := s.RestRepo.GetByID(c, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// profileId := _jwt.GetProfileId(c)
// 	profileId := _utils.GetProfileIdWithContext(c)

// 	// if asset.OwnerUserID == nil {
// 	// 	return nil, &_routes.Except{
// 	// 		Code:    401,
// 	// 		Message: "Bạn không có quyền truy cập",
// 	// 	}
// 	// }
// 	log.Print("*asset.OwnerUserID != profileId", *asset.M.OwnerUserID != profileId)
// 	if asset.M.OwnerUserID == nil || *asset.M.OwnerUserID != profileId {
// 		return nil, &_routes.Except{
// 			Code:    401,
// 			Message: "Bạn không có quyền truy cập",
// 		}
// 	}
// 	return asset, nil
// }

// // UC3: Chỉnh sửa thông tin tài sản trực tiếp từ danh sách
// func (uc *AssetService) UpdateAsset(ctx *gin.Context, id uint64, assetData *models.Asset) error {
// 	if id == 0 {
// 		return errors.New("asset ID is required")
// 	}
// 	asset, err := uc.RequiredOwner(ctx, id)
// 	if err != nil {
// 		return err
// 	}

// 	asset.M.Name = assetData.Name
// 	asset.M.WardID = assetData.WardID
// 	asset.M.DistrictID = assetData.DistrictID
// 	asset.M.ProvinceID = assetData.ProvinceID
// 	// asset.Price = assetData.Price
// 	// asset.BuyedDate = assetData.BuyedDate
// 	// asset.DocStatus = assetData.DocStatus
// 	// asset.ProductID = assetData.ProductID
// 	// asset.OwnerUserID = assetData.OwnerUserID
// 	// asset.OwnerOrgID = assetData.OwnerOrgID
// 	// asset.ParentAssetID = assetData.ParentAssetID
// 	asset.M.SplitMergeStatus = assetData.SplitMergeStatus
// 	asset.M.PurchasePrice = assetData.PurchasePrice
// 	asset.M.PurchaseDate = assetData.PurchaseDate
// 	asset.M.LegalStatus = assetData.LegalStatus
// 	// asset.Archived = assetData.Archived
// 	asset.M.RentStatus = assetData.RentStatus
// 	asset.M.Area = assetData.Area
// 	asset.M.Address = assetData.Address
// 	asset.M.Description = assetData.Description

// 	err = uc.AssetRepo.Update(ctx, id, asset)
// 	if err != nil {
// 		return err
// 	}

// 	afterData := _utils.StructToJSONString(struct{ Data *rest.RestModel[models.Asset] }{Data: asset})
// 	uc.AssetHistoryRepo.CreateHistory(ctx, enums.AssetHistory_Update, &assetData.ID, nil, afterData)

// 	return nil
// }

// // UC4: Xóa nhanh tài sản từ danh sách
// func (uc *AssetService) DeleteAsset(ctx *gin.Context, id uint64) error {
// 	if id == 0 {
// 		return errors.New("asset ID is required")
// 	}
// 	asset, err := uc.RequiredOwner(ctx, id)
// 	if err != nil {
// 		return err
// 	}

// 	if asset.M.RentStatus != enums.AssetStatusOwning {
// 		return errors.New("không thể xóa tài sản do trạng thái tài sản không cho phép")
// 	}
// 	return uc.AssetRepo.Delete(ctx, id)
// }

// // // UC5: Lọc danh sách tài sản theo trạng thái, loại hình tài sản, người sở hữu, diện tích, ngày tạo
// // func (uc *AssetService) FilterAssets(ctx *gin.Context, filters map[string]interface{}) ([]models.Asset, error) {
// // 	return uc.AssetRepo.FilterAssets(ctx, filters)
// // }

// // UC6: Tìm kiếm tài sản theo từ khóa hoặc mã tài sản
// func (uc *AssetService) SearchAssets(ctx *gin.Context, dto SearchDTO) (_routes.ResponseDTO, error) {
// 	profileId := _jwt.GetProfileId(ctx)
// 	assets, total, err := uc.AssetRepo.SearchOwner(ctx, profileId, dto)
// 	if err != nil {
// 		return _routes.ResponseDTO{}, err
// 	}
// 	return _routes.ResponseDTO{
// 		Code:          0,
// 		Data:          assets,
// 		TotalElements: &total,
// 	}, nil
// }

// // UC6: Tìm kiếm tài sản theo từ khóa hoặc mã tài sản
// func (uc *AssetService) SearchAssetInternal(ctx *gin.Context, dto SearchDTO) ([]models.Asset, error) {
// 	profileId := _jwt.GetProfileId(ctx)
// 	return uc.AssetRepo.SearchShare(ctx, profileId, dto)
// }

// // UC7: Tạo mới tài sản bằng nút truy cập nhanh
// func (uc *AssetService) CreateAsset(ctx *gin.Context, assetData *rest.RestModel[models.Asset]) (*rest.RestModel[models.Asset], error) {
// 	if assetData.M.Name == "" || assetData.M.Area <= 0 || assetData.M.WardID == nil {
// 		return nil, errors.New("missing required fields: name, area, or address")
// 	}
// 	id := _jwt.GetProfileId(ctx)
// 	assetData.M.OwnerUserID = &id
// 	assetData.M.Archived = false

// 	err := uc.AssetRepo.Create(ctx, assetData)
// 	afterData := _utils.StructToJSONString(struct {
// 		Data **rest.RestModel[models.Asset]
// 	}{Data: &assetData})

// 	uc.AssetHistoryRepo.CreateHistory(ctx, enums.AssetHistory_New, &assetData.ID, nil, afterData)
// 	return assetData, err
// }

// // UC8: Tạo tài sản từ sản phẩm có sẵn
// func (uc *AssetService) CreateAssetFromProduct(ctx *gin.Context,
// 	productID uint64,
// 	assetData *rest.RestModel[models.Asset]) (*rest.RestModel[models.Asset], error) {
// 	// if req.ProductId == nil {
// 	// 	return nil, fmt.Errorf("product_id is required")
// 	// }

// 	// // Có thể lấy thông tin sản phẩm và tự động điền (nếu cần)
// 	// product, err := s.productRestRepo.GetByID(*req.ProductId)
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("failed to get product: %w", err)
// 	// }

// 	// // Tự động điền các trường cơ bản từ product (nếu cần)
// 	// if req.Name == "" {
// 	// 	req.Name = product.Name
// 	// }
// 	// if req.Area == 0 {
// 	// 	req.Area = product.Area
// 	// }

// 	// return s.saveAsset(ctx, req)

// 	// Lấy thông tin sản phẩm từ ProductService
// 	product, err := uc.ProductUserRepo.GetByID(ctx, productID)
// 	if err != nil {
// 		return nil, errors.New("product not found")
// 	}

// 	// Liên kết thông tin sản phẩm vào tài sản
// 	assetData.M.ProductID = &product.ID
// 	assetData.M.Name = product.Name
// 	assetData.M.Description = product.Description
// 	assetData.M.Archived = false

// 	// Tạo tài sản
// 	err = uc.AssetRepo.Create(ctx, assetData)
// 	return assetData, err
// }

// // UC9: Tạo tài sản và sản phẩm mới
// func (uc *AssetService) CreateAssetAndProduct(ctx *gin.Context,
// 	productData *models.Product,
// 	assetData *rest.RestModel[models.Asset]) (*rest.RestModel[models.Asset], error) {
// 	// Tạo sản phẩm mới
// 	product, err := uc.AssetRepo.CreateProduct(ctx, productData)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Liên kết sản phẩm mới vào tài sản
// 	assetData.M.ProductID = &product.ID
// 	assetData.M.Archived = false

// 	uc.ProductUc.CreateProduct(ctx, &u_product_dto.ProductSaveRequest{
// 		Name:        assetData.M.Name,
// 		Area:        assetData.M.Area,
// 		WardID:      assetData.M.WardID,
// 		ProvinceID:  assetData.M.ProvinceID,
// 		DistrictID:  assetData.M.DistrictID,
// 		Description: assetData.M.Description,
// 		PriceDTO: &models.ProductPrice{
// 			SalePrice: &assetData.M.PurchasePrice,
// 		},
// 	})

// 	// Tạo tài sản
// 	err = uc.AssetRepo.Create(ctx, assetData)
// 	return assetData, err
// }

// // UC10: Chuyển sản phẩm thành tài sản
// func (uc *AssetService) ConvertProductToAsset(ctx *gin.Context,
// 	productID uint64,
// 	assetData *rest.RestModel[models.Asset]) (*rest.RestModel[models.Asset], error) {
// 	// Lấy thông tin sản phẩm từ ProductRepo
// 	product, err := uc.AssetRepo.GetProductByID(ctx, productID)
// 	if err != nil {
// 		return nil, errors.New("product not found")
// 	}

// 	// Liên kết thông tin sản phẩm vào tài sản
// 	assetData.M.ProductID = &product.ID
// 	assetData.M.Name = product.Name
// 	assetData.M.Description = product.Description

// 	// Tạo tài sản
// 	err = uc.AssetRepo.Create(ctx, assetData)
// 	return assetData, err
// }

// // UC11: Xóa tài sản (ẩn vĩnh viễn)
// // func (uc *AssetService) DeleteAssetWithValidation(ctx *gin.Context, assetID uint64) error {
// // 	if assetID == 0 {
// // 		return errors.New("asset ID is required")
// // 	}

// // 	// Lấy thông tin tài sản
// // 	asset, err := uc.AssetRepo.GetAssetDetails(ctx, assetID)
// // 	if err != nil {
// // 		return errors.New("asset not found")
// // 	}

// // 	// Kiểm tra sản phẩm liên kết
// // 	if asset.ProductID != nil {
// // 		// Lưu trữ sản phẩm liên kết
// // 		if err := uc.ProductUserRestRepo.UpdateArchived(ctx, *asset.ProductID, true); err != nil {
// // 			return errors.New("failed to archive linked product")
// // 		}
// // 	}

// // 	// Cập nhật trạng thái tài sản thành "deleted"
// // 	if err := uc.AssetRepo.Delete(ctx, assetID); err != nil {
// // 		return errors.New("failed to delete asset")
// // 	}

// // 	return nil
// // }

// // UC12: Lưu trữ tài sản
// func (uc *AssetService) ArchiveAsset(ctx *gin.Context, assetID uint64, archived bool) error {
// 	if assetID == 0 {
// 		return errors.New("asset ID is required")
// 	}

// 	// Lấy thông tin tài sản
// 	// _, err := uc.AssetRepo.GetAssetDetails(ctx, assetID)
// 	// if err != nil {
// 	// 	return err // errors.New("asset not found")
// 	// }

// 	asset, err := uc.RequiredOwner(ctx, assetID)
// 	if err != nil {
// 		return err
// 	}

// 	// Cập nhật trạng thái tài sản thành "archived"
// 	if err := uc.AssetRepo.UpdateArchived(ctx, assetID, archived); err != nil {
// 		return errors.New("failed to archive asset")
// 	}

// 	// log
// 	beforeData := _utils.StructToJSONString(struct{ Archived bool }{Archived: asset.M.Archived})
// 	afterData := _utils.StructToJSONString(struct{ Archived bool }{Archived: archived})
// 	uc.AssetHistoryRepo.CreateHistory(ctx, enums.AssetHistory_Archived, &assetID, beforeData, afterData)

// 	return nil
// }

// // func (uc *AssetService) ValidateAssetToAction(ctx *gin.Context, e *models.Asset) bool {
// // 	return e.RentStatus != 1
// // }

// // UC13: Gộp tài sản
// func (uc *AssetService) MergeAsset(ctx *gin.Context, dto MergeDTO) (*models.Asset, error) {
// 	// if len(dto.IDs) < 1 {
// 	// 	return nil, &_routes.Except{
// 	// 		Code:    401,
// 	// 		Message: "Bạn phải chọn tối thiểu 2 tài sản",
// 	// 	}
// 	// }

// 	profileId := _jwt.GetProfileId(ctx)
// 	entities, err := uc.AssetRepo.GetByIDsWithValidOwner(ctx, profileId, dto.IDs)
// 	if len(entities) != len(dto.IDs) {
// 		return nil, &_routes.Except{
// 			Code:    401,
// 			Message: "Bạn không sở hữu tài sản hoặc tài sản đang trong giao dịch chưa hoàn tất",
// 		}
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	currentAsset := entities[0]
// 	sumArea := 0.0
// 	parentId := (*uint64)(nil)

// 	for _, e := range entities {
// 		sumArea += e.Area
// 		if e.ParentAssetID != nil {
// 			parentId = e.ParentAssetID
// 		}
// 		if e.ParentAssetID == nil {
// 			return nil, &_routes.Except{
// 				Code:    400,
// 				Message: "Chỉ gộp tài sản con",
// 			}
// 		}
// 		if e.WardId != currentAsset.WardId ||
// 			e.DistrictId != currentAsset.DistrictId ||
// 			e.ProvinceId != currentAsset.ProvinceId {
// 			return nil, &_routes.Except{
// 				Code:    400,
// 				Message: "Danh sách tài sản không liền kề nhau",
// 			}
// 		}
// 	}

// 	if parentId != nil {
// 		parentAsset, err := uc.AssetRepo.GetByID(ctx, *parentId)
// 		if err != nil {
// 			return nil, err
// 		}
// 		beforeData := _utils.StructToJSONString(parentAsset)
// 		parentAsset.M.Area = parentAsset.M.Area + sumArea
// 		if err := uc.AssetRepo.Update(ctx, parentAsset.ID, parentAsset); err != nil {
// 			return nil, err
// 		}

// 		if err = uc.AssetRepo.DeleteBatch(ctx, dto.IDs); err != nil {
// 			return nil, err
// 		}
// 		// uc.AssetHistoryRepo.CreateHistory(ctx, enums.AssetHistory_Merge, parentId, dto.IDs)

// 		afterData := _utils.StructToJSONString(struct {
// 			Parent *rest.RestModel[models.Asset]
// 			Childs []models.Asset
// 		}{
// 			Parent: parentAsset,
// 			Childs: entities,
// 		})
// 		uc.AssetHistoryRepo.CreateHistory(ctx, enums.AssetHistory_Merge, parentId, beforeData, afterData)
// 		return &parentAsset.M, nil
// 	}

// 	// todo: đoạn này có thể không chạy đến do chỉ gộp con chứ không gộp cha

// 	nAsset := &rest.RestModel[models.Asset]{}
// 	copier.Copy(nAsset, &dto)

// 	// profileId := _jwt.GetProfileId(ctx)
// 	nAsset.M.OwnerUserID = &profileId
// 	nAsset.M.WardID = currentAsset.WardId
// 	nAsset.M.DistrictID = currentAsset.DistrictId
// 	nAsset.M.ProvinceID = currentAsset.ProvinceId
// 	nAsset.M.Archived = false
// 	nAsset.M.Area = sumArea
// 	nAsset.M.RentStatus = 1

// 	// Cập nhật trạng thái tài sản thành "archived"
// 	if err := uc.AssetRepo.Create(ctx, nAsset); err != nil {
// 		return nil, err
// 	}

// 	if err = uc.AssetRepo.DeleteBatch(ctx, dto.IDs); err != nil {
// 		return nil, err
// 	}

// 	// todo

// 	return &currentAsset, err
// }

// // UC13: Tách tài sản
// func (uc *AssetService) SplitAsset(ctx *gin.Context, dto SplitDTO) ([]rest.RestModel[models.Asset], error) {
// 	currentAsset, err := uc.RequiredOwner(ctx, dto.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// todo: tài sản liên kết nhiều sp
// 	numProductLink, err := uc.ProductUserRepo.CountByAssetID(currentAsset.ID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if numProductLink > 1 {
// 		return nil, &_routes.Except{
// 			Code:    400,
// 			Message: "Không thể tách do tài sản đang liên kết với nhiều sản phẩm",
// 		}
// 	}

// 	// currentAsset := entities[0]
// 	sumArea := 0.0
// 	for _, e := range dto.Datas {
// 		sumArea += e.M.Area
// 		if e.M.Area < 1 {
// 			return nil, &_routes.Except{
// 				Code:    400,
// 				Message: "Diện tích không hợp lệ",
// 			}
// 		}
// 	}

// 	if sumArea >= currentAsset.M.Area {
// 		return nil, &_routes.Except{
// 			Code:    400,
// 			Message: "Diện tích không hợp lệ",
// 		}
// 	}

// 	profileId := _jwt.GetProfileId(ctx)
// 	for i := range dto.Datas {
// 		nAsset := &dto.Datas[i]
// 		nAsset.M.OwnerUserID = &profileId
// 		nAsset.M.WardId = currentAsset.M.WardID
// 		nAsset.M.DistrictId = currentAsset.M.DistrictID
// 		nAsset.M.ProvinceId = currentAsset.M.ProvinceID
// 		nAsset.M.Archived = false
// 		nAsset.M.RentStatus = 1
// 		nAsset.M.ParentAssetID = &dto.ID
// 		// nAsset.Area = sumArea
// 	}

// 	// Cập nhật trạng thái tài sản thành "archived"
// 	if err := uc.AssetRepo.CreateBatch(ctx, &dto.Datas); err != nil {
// 		return nil, err
// 	}

// 	beforeData := _utils.StructToJSONString(currentAsset)

// 	currentAsset.M.Area = currentAsset.M.Area - sumArea
// 	if err := uc.AssetRepo.Update(ctx, currentAsset.ID, currentAsset); err != nil {
// 		return nil, err
// 	}

// 	afterData := _utils.StructToJSONString(struct {
// 		Parent *rest.RestModel[models.Asset]
// 		Childs []rest.RestModel[models.Asset]
// 	}{
// 		Parent: currentAsset,
// 		Childs: dto.Datas,
// 	})

// 	// var ids []uint64
// 	// for _, r := range dto.Datas {
// 	// 	ids = append(ids, r.ID)
// 	// }
// 	// uc.AssetHistoryRepo.CreateHistory(ctx, enums.AssetHistory_Slit, &dto.ID, ids)
// 	uc.AssetHistoryRepo.CreateHistory(ctx, enums.AssetHistory_Slit, &dto.ID, beforeData, afterData)

// 	// todo

// 	return dto.Datas, nil
// }

// // UC13: Tách tài sản
// func (uc *AssetService) History(ctx *gin.Context, assetID uint64, dto HistoryDTO) (_routes.ResponseDTO, error) {
// 	// profileId := _jwt.GetProfileId(ctx)
// 	results, total, err := uc.AssetRepo.History(ctx, 0, assetID, &dto)
// 	if err != nil {
// 		return _routes.ResponseDTO{}, err
// 	}
// 	return _routes.ResponseDTO{
// 		Code:          0,
// 		Data:          results,
// 		TotalElements: &total,
// 	}, nil
// }

// // UC14: Lịch sử tài sản
// func (uc *AssetService) HistoryOwner(ctx *gin.Context, dto HistoryDTO) (_routes.ResponseDTO, error) {
// 	profileId := _jwt.GetProfileId(ctx)
// 	results, total, err := uc.AssetRepo.History(ctx, profileId, 0, &dto)
// 	if err != nil {
// 		return _routes.ResponseDTO{}, err
// 	}
// 	return _routes.ResponseDTO{
// 		Code:          0,
// 		Data:          results,
// 		TotalElements: &total,
// 	}, nil
// }
