package shared_usecase

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/provider"
	"bdspro/internal/repo"
	"bdspro/internal/utils"
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	"common/logging"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/jinzhu/copier"
)

type IAssetUsecase interface {
	GetOwnerID(c context.Context) (*uint64, error)
	HasPermission(c context.Context, assetId *uint64, permission enums.EPermission) (bool, error)
	GetOwnerType(c context.Context) enums.EOwnerOf
}

type AssetUsecase struct {
	AssetRepo             repo.AssetRepo
	ProductUserRepo       repo.ProductRepo
	NotificationClient    provider.NotificationProvider
	AssetLegalRepo        repo.AssetLegalRepo
	AssetOrganizationRepo repo.AssetOrganizationRepo
	AssetUserRepo         repo.AssetUserRepo
	AssetUsecase          IAssetUsecase
	RequestUsecase        *RequestUsecase
}

func NewSharedAssetUsecase(
	assetRepo repo.AssetRepo,
	ProductRepo repo.ProductRepo,
	notificationClient provider.NotificationProvider,
	AssetLegalRepo repo.AssetLegalRepo,
	AssetOrganizationRepo repo.AssetOrganizationRepo,
	AssetUserRepo repo.AssetUserRepo,
	requestUsecase *RequestUsecase,
) *AssetUsecase {
	return &AssetUsecase{
		AssetRepo:             assetRepo,
		ProductUserRepo:       ProductRepo,
		NotificationClient:    notificationClient,
		AssetLegalRepo:        AssetLegalRepo,
		AssetOrganizationRepo: AssetOrganizationRepo,
		AssetUserRepo:         AssetUserRepo,
		RequestUsecase:        requestUsecase,
		// ProductUc:       ProductUc,
	}
}

// UC1: Xem danh sách tổng quan tất cả tài sản với phân trang và profileId
func (uc *AssetUsecase) GetOwnerAssets(ctx context.Context, profileId uint64, limit, offset int) ([]domain.Asset, error) {
	return uc.AssetRepo.GetAssets(ctx, profileId, limit, offset)
}

// UC2: Thực hiện thao tác hàng loạt (cập nhật trạng thái, xóa đồng thời nhiều tài sản)
func (uc *AssetUsecase) BulkAction(ctx context.Context, action string, ids []uint64) error {
	if len(ids) == 0 {
		return errors.New("no asset IDs provided for bulk action")
	}
	return uc.AssetRepo.BulkAction(ctx, action, ids)
}

func (s *AssetUsecase) RequiredOwner(c context.Context, id uint64) (*domain.Asset, error) {
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
	logging.FromContext(c).Debug(
		"asset owner mismatch check",
		slog.Bool("owner_mismatch", *asset.OwnerID != profileId),
	)
	if asset.OwnerID == nil || *asset.OwnerID != profileId {
		return nil, &_routes.Except{
			Code:    401,
			Message: "Bạn không có quyền truy cập",
		}
	}
	return asset, nil
}

// UC3: Chỉnh sửa thông tin tài sản trực tiếp từ danh sách
func (uc *AssetUsecase) UpdateAsset(ctx context.Context, id uint64, asset *domain.Asset) (*domain.Asset, error) {
	if id == 0 {
		return nil, errors.New("asset ID is required")
	}
	legalItems := []domain.AssetLegal{}
	if asset.LegalItems != nil {
		legalItems = asset.LegalItems
	}
	asset.LegalItems = nil

	_, err := uc.RequiredOwner(ctx, id)
	if err != nil {
		return nil, err
	}

	err = uc.AssetRepo.Update(ctx, id, asset)
	if err != nil {
		return nil, err
	}

	if len(legalItems) > 0 {
		uc.AssetLegalRepo.UpdateLegalItems(ctx, asset.ID, legalItems)
		// Lấy legal item đầu tiên và gán vào ImageID
		legalList, _, err := uc.AssetLegalRepo.GetData(ctx, dto.AssetLegalGetDTO{
			AssetID: asset.ID,
		})
		if err == nil && len(legalList) > 0 {
			asset.ImageID = &legalList[0].ID
			// Cập nhật lại ImageID vào database
			uc.AssetRepo.Update(ctx, id, asset)
		}
	} else {
		// Nếu không có legal items, xóa ImageID
		asset.ImageID = nil
		uc.AssetRepo.Update(ctx, id, asset)
	}

	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		// afterData := _utils.StructToJSONString(struct{ Data *domain.Asset }{Data: asset})
		uc.NotificationClient.CreateHistory(cloneCtx, &_dto.HistoryDTO{
			ActionType: _enum.HistoryAssetUpdate,
			TargetId:   asset.ID,
			OwnerID:    &id,
			Note:       []string{"Cập nhật tài sản ", asset.Name, " thành công"},
			Title:      "Cập nhật tài sản thành công",
		})
	}()

	return asset, nil
}

// UC4: Xóa nhanh tài sản từ danh sách
func (uc *AssetUsecase) DeleteAsset(ctx context.Context, id uint64) error {
	if id == 0 {
		return errors.New("asset ID is required")
	}
	asset, err := uc.RequiredOwner(ctx, id)
	if err != nil {
		return err
	}

	if asset.RentStatus != enums.EAssetStatusOwning {
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
func (uc *AssetUsecase) CreateAsset(ctx context.Context, asset *domain.Asset) (*domain.Asset, error) {
	if asset.Name == "" || asset.Area <= 0 || asset.ProvinceID == nil || asset.OwnerOf == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "missing required fields: name, area, provinceId, ownerType",
		}
	}

	legalItems := []domain.AssetLegal{}
	if asset.LegalItems != nil {
		legalItems = asset.LegalItems
	}
	asset.LegalItems = nil

	ownerID, ownerType, err := uc.RequestUsecase.GetOwnerIDAndType(ctx)
	if err != nil {
		return nil, err
	}
	asset.OwnerID = ownerID
	asset.OwnerOf = ownerType
	asset.Archived = false

	// asset := &domain.Asset{}
	// copier.Copy(asset, assetData)

	err = uc.AssetRepo.Create(ctx, asset)

	if asset.ProductID != nil {
		uc.ProductUserRepo.UpdateAssetID(ctx, *asset.ProductID, asset.ID)
	}

	if len(legalItems) > 0 {
		uc.AssetLegalRepo.UpdateLegalItems(ctx, asset.ID, legalItems)
		// Lấy legal item đầu tiên và gán vào ImageID
		legalList, _, err := uc.AssetLegalRepo.GetData(ctx, dto.AssetLegalGetDTO{
			AssetID: asset.ID,
		})
		if err == nil && len(legalList) > 0 {
			asset.ImageID = &legalList[0].ID
			// Cập nhật lại ImageID vào database
			uc.AssetRepo.Update(ctx, asset.ID, asset)
		}
	}

	if asset.RequestOrganizationID != nil {
		uc.AssetOrganizationRepo.CreateOwner(ctx, asset.ID, *asset.RequestOrganizationID)
	} else {
		profileId := _utils.GetProfileIdWithContext(ctx)
		uc.AssetUserRepo.CreateOwner(ctx, asset.ID, profileId)
	}

	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		// afterData := _utils.StructToJSONString(struct{ Data *domain.Asset }{Data: asset})
		uc.NotificationClient.CreateHistory(cloneCtx, &_dto.HistoryDTO{
			ActionType: _enum.HistoryAssetCreate,
			TargetId:   asset.ID,
			OwnerID:    &asset.ID,
			Note:       []string{"Tạo tài sản ", asset.Name, " thành công"},
			Title:      "Tạo tài sản thành công",
		})
	}()

	return asset, err
}

// UC8: Tạo tài sản từ sản phẩm có sẵn
// func (uc *AssetUsecase) CreateAssetFromProduct(ctx *gin.Context, assetData *dto.AssetSaveRequest) (*domain.Asset, error) {
// 	// if req.ProductId == nil {
// 	// 	return nil, fmt.Errorf("product_id is required")
// 	// }

// 	// // Có thể lấy thông tin sản phẩm và tự động điền (nếu cần)
// 	// product, err := s.productRepo.GetByID(*req.ProductId)
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
// 	product, err := uc.ProductUserRepo.GetByID(*assetData.ProductID)
// 	if err != nil {
// 		return nil, errors.New("product not found")
// 	}

// 	// Liên kết thông tin sản phẩm vào tài sản
// 	assetData.ProductID = &product.ID
// 	ownerID, ownerType, err := uc.RequestUsecase.GetOwnerIDAndType(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	assetData.OwnerID = ownerID
// 	assetData.OwnerType = ownerType
// 	assetData.Archived = false

// 	// Tạo tài sản
// 	err = uc.AssetRepo.Create(ctx, &assetData.Asset)

// 	if assetData.ProductID != nil {
// 		uc.ProductUserRepo.UpdateAssetID(ctx, *assetData.ProductID, assetData.ID)
// 	}
// 	uc.AssetLegalRepo.UpdateLegalItems(ctx, assetData.ID, assetData.LegalItems)

// 	afterData := _utils.StructToJSONString(struct{ Data *domain.Asset }{Data: &assetData.Asset})
// 	uc.RecordHistoryRepo.CreateHistory(ctx, enums.EActionAssetCreate, &assetData.ID, nil, []uint64{}, "", *afterData, assetData.ID)

// 	return &assetData.Asset, err
// }

// // UC9: Tạo tài sản và sản phẩm mới
// func (uc *AssetUsecase) CreateAssetAndProduct(ctx context.Context, productData *domain.Product, assetData *domain.Asset) (*domain.Asset, error) {
// 	// Tạo sản phẩm mới
// 	product, err := uc.AssetRepo.CreateProduct(ctx, productData)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Liên kết sản phẩm mới vào tài sản
// 	assetData.ProductID = &product.ID
// 	assetData.Archived = false

// 	// uc.ProductUc.CreateProduct(ctx, &dto.ProductSaveRequest{
// 	// 	Name:        assetData.Name,
// 	// 	Area:        assetData.Area,
// 	// 	WardID:      assetData.WardID,
// 	// 	ProvinceID:  assetData.ProvinceID,
// 	// 	DistrictID:  assetData.DistrictID,
// 	// 	Description: assetData.Description,
// 	// 	PriceDTO: &domain.ProductPrice{
// 	// 		SalePrice: &assetData.PurchasePrice,
// 	// 	},
// 	// })

// 	// Tạo tài sản
// 	err = uc.AssetRepo.Create(ctx, assetData)
// 	return assetData, err
// }

// UC10: Chuyển sản phẩm thành tài sản
func (uc *AssetUsecase) ConvertProductToAsset(ctx context.Context, productID uint64, assetData *domain.Asset) (*domain.Asset, error) {
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
func (uc *AssetUsecase) ArchiveAsset(ctx context.Context, assetID uint64, archived bool) error {
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
func (uc *AssetUsecase) MergeAsset(ctx context.Context, dto dto.MergeAssetDTO) (*domain.Asset, error) {
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

	var result *domain.Asset
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
		// todo: tạo lịch sử gộp tài sản
		// beforeData := _utils.StructToJSONString(parentAsset)
		parentAsset.Area = parentAsset.Area + sumArea
		if err := uc.AssetRepo.Update(ctx, parentAsset.ID, parentAsset); err != nil {
			return nil, err
		}

		if err = uc.AssetRepo.DeleteBatch(ctx, dto.IDs); err != nil {
			return nil, err
		}
		// uc.AssetHistoryRepo.CreateHistory(ctx, enums.AssetHistory_Merge, parentId, dto.IDs)

		// afterData := _utils.StructToJSONString(struct {
		// 	Parent *domain.Asset
		// 	Childs []domain.Asset
		// }{
		// 	Parent: parentAsset,
		// 	Childs: entities,
		// })
		// uc.RecordHistoryRepo.CreateHistory(ctx, enums.EActionAssetMerge, parentId, beforeData, afterData, nil, "")
		// return parentAsset, nil
		result = parentAsset
	} else {
		// todo: đoạn này có thể không chạy đến do chỉ gộp con chứ không gộp cha

		nAsset := &domain.Asset{}
		copier.Copy(nAsset, &dto)

		ownerID, ownerType, err := uc.RequestUsecase.GetOwnerIDAndType(ctx)
		if err != nil {
			return nil, err
		}
		nAsset.OwnerID = ownerID
		nAsset.OwnerOf = ownerType
		nAsset.WardID = currentAsset.WardID
		// nAsset.DistrictID = currentAsset.DistrictID
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

		result = nAsset
	}

	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		var arr []string
		for _, r := range dto.IDs {
			arr = append(arr, utils.ConvertToBDS(r))
		}
		childCodes := strings.Join(arr, ", ")
		uc.NotificationClient.CreateHistory(cloneCtx, &_dto.HistoryDTO{
			ActionType: _enum.HistoryAssetMerge,
			TargetId:   result.ID,
			OwnerID:    &profileId,
			Title:      "Gộp tài sản đã chọn thành công!",
			Note:       []string{"Gộp tài sản ", childCodes, " thành tài sản ", utils.ConvertToBDS(result.ID)},
		})
	}()

	// todo: tạo thông báo gộp tài sản

	return result, err
}

// UC13: Tách tài sản
func (uc *AssetUsecase) SplitAsset(ctx context.Context, dto dto.SplitAssetDTO) ([]domain.Asset, error) {
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

	for i := range dto.Datas {
		nAsset := &dto.Datas[i]
		ownerID, ownerType, err := uc.RequestUsecase.GetOwnerIDAndType(ctx)
		if err != nil {
			return nil, err
		}
		nAsset.OwnerID = ownerID
		nAsset.OwnerOf = ownerType
		nAsset.WardID = currentAsset.WardID
		// nAsset.DistrictID = currentAsset.DistrictID
		nAsset.ProvinceID = currentAsset.ProvinceID
		nAsset.Archived = false
		nAsset.RentStatus = enums.EAssetStatusOwning
		nAsset.ParentAssetID = &dto.ID
		nAsset.PurchasePrice = dto.Datas[i].PurchasePrice
		// nAsset.Area = sumArea
	}

	// Cập nhật trạng thái tài sản thành "archived"
	if err := uc.AssetRepo.CreateBatch(ctx, dto.Datas); err != nil {
		return nil, err
	}

	// beforeData := _utils.StructToJSONString(currentAsset)

	currentAsset.Area = currentAsset.Area - sumArea
	if err := uc.AssetRepo.Update(ctx, currentAsset.ID, currentAsset); err != nil {
		return nil, err
	}

	// afterData := _utils.StructToJSONString(struct {
	// 	Parent *domain.Asset
	// 	Childs []domain.Asset
	// }{
	// 	Parent: currentAsset,
	// 	Childs: dto.Datas,
	// })

	// var ids []uint64
	// for _, r := range dto.Datas {
	// 	ids = append(ids, r.ID)
	// }
	// uc.AssetHistoryRepo.CreateHistory(ctx, enums.AssetHistory_Slit, &dto.ID, ids)
	// todo: tạo lịch sử tách tài sản
	go func() {
		cloneCtx := _utils.CloneContext(ctx)

		var arr []string
		for _, r := range dto.Datas {
			arr = append(arr, utils.ConvertToBDS(r.ID))
		}
		childCodes := strings.Join(arr, ", ")
		uc.NotificationClient.CreateHistory(cloneCtx, &_dto.HistoryDTO{
			ActionType: _enum.HistoryAssetSplit,
			TargetId:   dto.ID,
			OwnerID:    &dto.ID,
			Note:       []string{"Tách tài sản ", utils.ConvertToBDS(currentAsset.ID), " thành các tài sản con :", childCodes},
			Title:      "Tách tài sản thành công!",
		})
	}()

	// todo: tạo thông báo tách tài sản

	return dto.Datas, nil
}

func (uc *AssetUsecase) SearchAssetMe(ctx context.Context, dto dto.AssetSearchRequest) ([]domain.AssetList, int64, error) {
	ownerID, ownerType, err := uc.RequestUsecase.GetOwnerIDAndType(ctx)
	if err != nil {
		return nil, 0, err
	}
	assets, total, err := uc.AssetRepo.SearchOwner(ctx, *ownerID, ownerType, dto)
	if err != nil {
		return nil, 0, err
	}
	return assets, total, nil
}

// UC6: Tìm kiếm tài sản theo từ khóa hoặc mã tài sản
func (uc *AssetUsecase) SearchAssetShare(ctx context.Context, dto dto.AssetSearchRequest) ([]domain.AssetList, int64, error) {
	ownerID, ownerType, err := uc.RequestUsecase.GetOwnerIDAndType(ctx)
	if err != nil {
		return nil, 0, err
	}
	results, total, err := uc.AssetRepo.SearchShare(ctx, *ownerID, ownerType, dto)
	if err != nil {
		return nil, 0, err
	}
	return results, total, nil
}

// // UC13: Tách tài sản
// func (uc *AssetUsecase) History(ctx context.Context, assetID uint64, dto dto.AssetHistoryDTO) (*[]domain.AssetHistory, int64, error) {
// 	// profileId := _jwt.GetProfileId(ctx)
// 	ownerID, ownerType, err := uc.RequestUsecase.GetOwnerIDAndType(ctx)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	results, total, err := uc.AssetRepo.History(ctx, ownerID, ownerType, assetID, &dto)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	return results, total, nil
// }

// // UC14: Lịch sử tài sản
// func (uc *AssetUsecase) HistoryOwner(ctx context.Context, dto dto.AssetHistoryDTO) (_routes.ResponseDTO, error) {
// 	ownerID, ownerType, err := uc.RequestUsecase.GetOwnerIDAndType(ctx)
// 	if err != nil {
// 		return _routes.ResponseDTO{}, err
// 	}
// 	results, total, err := uc.AssetRepo.History(ctx, ownerID, ownerType, 0, &dto)
// 	if err != nil {
// 		return _routes.ResponseDTO{}, err
// 	}
// 	return _routes.ResponseDTO{
// 		Code:          0,
// 		Data:          results,
// 		TotalElements: &total,
// 	}, nil
// }

func (uc *AssetUsecase) GetAssetsWithChildCount(ctx context.Context, searchDto dto.AssetSearchRequest) ([]dto.AssetWithChildCountDTO, int64, error) {
	// Validate input if needed
	profileId := _utils.GetProfileIdWithContext(ctx)

	// Call repository to get data
	assets, total, err := uc.AssetRepo.GetAssetsWithChildCount(ctx, profileId, searchDto)
	if err != nil {
		return nil, 0, err
	}

	return assets, total, nil
}

func (uc *AssetUsecase) GetAssetDetail(ctx context.Context, assetID uint64) (*domain.Asset, error) {
	asset, err := uc.AssetRepo.GetByID(assetID)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func (uc *AssetUsecase) GetPublish(ctx context.Context, profileId uint64) ([]domain.Asset, error) {
	assets, _, err := uc.AssetRepo.GetAssetPublish(profileId)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func (s *AssetUsecase) SearchByAssetShare(c context.Context, targetID *uint64, targetType enums.EOwnerOf, dto *dto.ProductSearchRequest) ([]domain.Product, int64, error) {
	return s.AssetRepo.SearchByAssetShare(c, targetID, targetType, dto)
}

// CreateAssetOrganization tạo tài sản cho tổ chức
func (s *AssetUsecase) CreateOrganizationAsset(c context.Context, asset *domain.Asset) (*domain.Asset, error) {
	organizationId := _utils.GetOrganizationIdFromContext(c)

	asset.RequestOrganizationID = &organizationId
	return s.CreateAsset(c, asset)
}

func (s *AssetUsecase) GetAssetOfCurrentUser(c context.Context, dto *dto.AssetSearchRequest) ([]domain.Asset, int64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	dto.RequestUserId = &profileId
	return s.searchAsset(c, dto)
}

func (s *AssetUsecase) GetAssetOfCurrentOrganization(c context.Context, dto *dto.AssetSearchRequest) ([]domain.Asset, int64, error) {
	organizationId := _utils.GetOrganizationIdFromContext(c)
	dto.RequestOrganizationId = &organizationId
	return s.searchAsset(c, dto)
}

func (s *AssetUsecase) searchAsset(c context.Context, dto *dto.AssetSearchRequest) ([]domain.Asset, int64, error) {
	results, total, err := s.AssetRepo.ListAssets(c, dto)
	if err != nil {
		return nil, 0, err
	}
	// for i := 0; i < len(results); i++ {
	// 	s.MapAssetStatus(c, &results[i])
	// }
	return results, total, nil
}
