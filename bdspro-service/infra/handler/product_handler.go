package handler

import (
	"bdspro/infra/client"
	"bdspro/infra/mapper"
	"bdspro/infra/validator"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	"bdspro/internal/usecases"
	shared_usecase "bdspro/internal/usecases/shared"
	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductHandler struct {
	bdspropb.UnimplementedProductServiceServer
	UC                 *shared_usecase.ProductUsecase
	ProductRepo        repo.ProductRepo
	ProductMapper      *mapper.ProductMapper
	ProductPriceMapper *mapper.ProductPriceMapper
	ProductValidator   *validator.ProductValidator
	UserClient         *client.UserClient
	DepositeMapper     *mapper.DepositeMapper
	// TransactionMapper mapper.TransactionTransformer
	PostMapper     *mapper.PostMapper
	ProductAssetUC *usecases.ProductAssetUsecase
	AssetMapper    *mapper.AssetMapper
	DealMapper     mapper.DealMapper
	SyncProvider   *_utils.SyncUtil
}

func NewProductService(
	uc *shared_usecase.ProductUsecase,
	productRepo repo.ProductRepo,
	productMapper *mapper.ProductMapper,
	productPriceMapper *mapper.ProductPriceMapper,
	productValidator *validator.ProductValidator,
	userClient *client.UserClient,
	depositeMapper *mapper.DepositeMapper,
	// transactionMapper mapper.TransactionTransformer,
	postMapper *mapper.PostMapper,
	productAssetUC *usecases.ProductAssetUsecase,
	assetMapper *mapper.AssetMapper,
	dealMapper mapper.DealMapper,
	SyncProvider *_utils.SyncUtil,
) *ProductHandler {
	return &ProductHandler{
		UC:                 uc,
		ProductRepo:        productRepo,
		ProductMapper:      productMapper,
		ProductPriceMapper: productPriceMapper,
		ProductValidator:   productValidator,
		UserClient:         userClient,
		DepositeMapper:     depositeMapper,
		// TransactionMapper: transactionMapper,
		PostMapper:     postMapper,
		ProductAssetUC: productAssetUC,
		AssetMapper:    assetMapper,
		DealMapper:     dealMapper,
		SyncProvider:   SyncProvider,
	}
}

func (h *ProductHandler) CheckVersionSync(
	ctx context.Context,
	req *bdspropb.CheckVersionSyncRequest,
) (*bdspropb.CheckVersionSyncResponse, error) {
	ownerId := _utils.GetProfileIdWithContext(ctx)
	// Validate request
	if req.Resource == "" {
		return nil, status.Error(codes.InvalidArgument, "resource is required")
	}
	if ownerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "owner_id is required")
	}

	// Convert từ proto sang DTO
	dtoReq := &dto.CheckVersionSyncRequest{
		Resource: req.Resource,
		LastSync: req.LastSync,
		Id:       req.Id,
	}

	// Gọi usecase - CHỈ check Redis
	resp, err := h.UC.CheckVersionSync(ctx, ownerId, dtoReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert từ DTO sang proto
	return &bdspropb.CheckVersionSyncResponse{
		NeedSync:       resp.NeedSync,
		CurrentVersion: resp.CurrentVersion,
		ChangeCount:    int32(resp.ChangeCount),
		ChangedIds:     resp.ChangedIDs,
	}, nil
}

func (h *ProductHandler) SyncProducts(ctx context.Context, req *bdspropb.SyncProductsRequest) (*bdspropb.SyncProductsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request required")
	}

	profileID := req.ProfileId
	if profileID == 0 {
		profileID = _utils.GetProfileIdWithContext(ctx)
		if profileID == 0 {
			return nil, status.Error(codes.Unauthenticated, "profile_id not found")
		}
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	syncReq := &dto.SyncProductsRequest{
		ProfileID:    profileID,
		LastSyncTime: req.LastVersionTimestamp,
		PageSize:     int(pageSize),
		PageToken:    req.PageToken,
		ProductIDs:   req.Ids,
	}

	resp, err := h.UC.SyncProducts(ctx, syncReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "sync failed: %v", err)
	}

	changedProducts := make([]*bdspropb.ChangedProduct, len(resp.ProductIDs))
	for i, p := range resp.ProductIDs {
		changedProducts[i] = &bdspropb.ChangedProduct{
			Id:         p.ID,
			UpdatedAt:  p.UpdatedAt,
			ChangeType: p.ChangeType,
		}
	}

	return &bdspropb.SyncProductsResponse{
		ProductIds:      changedProducts,
		RemovedIds:      resp.RemovedIDs,
		CurrentSyncTime: resp.CurrentSyncTime,
		NextToken:       resp.NextToken,
		HasMore:         resp.HasMore,
	}, nil
}

func uniqueUint64(ids []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(ids))
	out := make([]uint64, 0, len(ids))

	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func (h *ProductHandler) BatchGetProducts(
	ctx context.Context,
	req *bdspropb.BatchGetProductsRequest,
) (*bdspropb.BatchGetProductsResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request required")
	}

	if len(req.Ids) == 0 {
		return &bdspropb.BatchGetProductsResponse{
			Products: map[uint64]*bdspropb.ProductResponse{},
		}, nil
	}

	ids := uniqueUint64(req.Ids)

	if len(ids) > 100 {
		return nil, status.Error(
			codes.InvalidArgument,
			"max 100 ids per request",
		)
	}

	// resolve profile
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId == 0 {
		return nil, status.Error(codes.Unauthenticated, "profile_id not found")
	}

	// build search request
	searchReq := &dto.ProductSearchRequest{
		IDs:           ids,
		RequestUserId: &profileId,
		Pagable: _dto.Pagable{
			Page: 1,
			Size: uint32(len(ids)),
		},
	}

	// call usecase
	result, _, err := h.UC.GetProductOfCurrentUser(ctx, *searchReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to batch get products: %v", err)
	}

	products := h.ProductMapper.MapProductPbListResponse(result)

	respMap := make(map[uint64]*bdspropb.ProductResponse, len(products))
	for _, p := range products {
		respMap[p.Id] = p
	}

	return &bdspropb.BatchGetProductsResponse{
		Products: respMap,
	}, nil
}

func (h *ProductHandler) FlushSyncIds(
	ctx context.Context,
	req *bdspropb.CheckVersionSyncRequest,
) (*bdspropb.Response, error) {
	err := h.UC.FlushSyncIds(ctx, req.Limit)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

func (h *ProductHandler) GetRefreshSyncIds(
	ctx context.Context,
	req *bdspropb.CheckVersionSyncRequest,
) (*bdspropb.CheckVersionSyncResponse, error) {
	resp, err := h.ProductRepo.GetRefreshSyncIds(ctx, &dto.CheckVersionSyncRequest{
		Resource: req.Resource,
		LastSync: req.LastSync,
		Id:       req.Id,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &bdspropb.CheckVersionSyncResponse{
		NeedSync:       true,
		CurrentVersion: time.Now().UnixMilli(),
		ChangeCount:    0,
		ChangedIds:     resp,
		Next:           nil,
	}, nil
}

func (h *ProductHandler) GetVersionData(
	ctx context.Context,
	req *bdspropb.GetVersionDataRequest,
) (*bdspropb.GetVersionDataResponse, error) {
	ownerId := _utils.GetProfileIdWithContext(ctx)
	if ownerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "owner_id is required")
	}
	dtoReq := &dto.CheckVersionSyncRequest{
		Resource: req.Resource,
	}
	_, _, err := h.UC.GetVersionData(ctx, dtoReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return nil, nil
}

// @Summary Lưu trữ sản phẩm
// @Description API này cho phép lưu trữ sản phẩm
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body bdspropb.ArchivedRequest true "Thông tin sản phẩm cần cập nhật"
// @Router /v2/bdspro/v2/product/archived [put]
func (s *ProductHandler) ArchiveProduct(c context.Context, req *bdspropb.ArchivedRequest) (*bdspropb.Response, error) {
	if req.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	err := s.UC.Archived(c, dto.ArchivedRequest{
		ProductId: &req.Id,
		Archived:  req.Archived,
	})
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Chia sản phẩm con
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body bdspropb.SaveChildRequest true "Thông tin sản phẩm con"
// @Router /v2/bdspro/v2/product/child/devide [post]
func (s *ProductHandler) ChildDevideProduct(ctx context.Context, req *bdspropb.DevideChildRequest) (*bdspropb.Response, error) {
	dto := s.ProductMapper.MapDevideChildRequest(req)
	if err := s.ProductValidator.ValidateDevideChildRequest(dto); err != nil {
		return nil, err
	}

	_, err := s.UC.DevideProductChild(ctx, dto)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// // @Summary Lấy lịch sử sản phẩm con
// // @Description API
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "Product ID"
// // @Param body query dto.ProductHistorySearch true "Thông tin sản phẩm con"
// // @Router /v2/bdspro/v2/product/child/history/{id} [get]
// func (s *ProductService) ChildHistories(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.HistoryResponse, error) {
// 	id := req.Id

// 	_body := dto.ProductHistorySearch{}
// 	if err := copier.Copy(&_body, req); err != nil {
// 		return nil, err
// 	}

// 	result, total, err := s.UC.ChildHistory(ctx, id, &_body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	histories := make([]*bdspropb.History, 0)
// 	copier.Copy(&histories, &result)
// 	return &bdspropb.GetHistoryResponse{
// 		Data:  histories,
// 		Total: uint64(total),
// 	}, nil
// }

// @Summary Gộp sản phẩm con
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.MergeProductChildRequest true "Thông tin sản phẩm con"
// @Router /v2/bdspro/v2/product/child/merge [post]
func (s *ProductHandler) ChildMergeProduct(ctx context.Context, req *bdspropb.MergeChildRequest) (*bdspropb.Response, error) {
	_, err := s.UC.MergeProductChild(ctx, &dto.MergeProductChildRequest{
		ParentID: &req.ParentId,
		ChildIDs: &req.ChildIds,
	})
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Gộp toàn bộ sản phẩm con
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Router /v2/bdspro/v2/product/child/merge-all/{id} [post]
func (s *ProductHandler) ChildMergeAll(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.Response, error) {
	_, err := s.UC.MergeAllProductChild(ctx, &req.Id)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Chia sản phẩm con
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Param body body bdspropb.SaveChildRequest true "Thông tin sản phẩm con"
// @Security BearerAuth
// @Router /v2/bdspro/v2/product/child/split [post]
func (s *ProductHandler) ChildSplit(ctx context.Context, req *bdspropb.SaveChildRequest) (*bdspropb.Response, error) {
	dto := s.ProductMapper.SaveProductToEntity(req)
	if err := s.ProductValidator.ValidateSaveProductChild(dto); err != nil {
		return nil, err
	}
	_, err := s.UC.CreateChild(ctx, dto)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Cọc
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body bdspropb.DepositeRequest true "Thông tin sản phẩm cần cập nhật"
// @Router /v2/bdspro/v2/product/deposite [post]
func (s *ProductHandler) ProductDeposite(ctx context.Context, req *bdspropb.DepositeRequest) (*bdspropb.Transaction, error) {
	// dto := s.DepositeMapper.MapDepositePbToDomain(req)

	// if dto.DepositeAmount == nil || *dto.DepositeAmount <= 0 {
	// 	return nil, status.Errorf(codes.InvalidArgument, "deposite amount is required")
	// }

	// if dto.ContactID == nil || *dto.ContactID == 0 {
	// 	return nil, status.Errorf(codes.InvalidArgument, "contactId is required")
	// }

	// if dto.TransactionType != enums.TransactionTypeSale && dto.TransactionType != enums.TransactionTypeRent {
	// 	return nil, status.Errorf(codes.InvalidArgument, "transactionType is invalid")
	// }

	// result, err := s.UC.ProductDeposite(ctx, dto)
	// if err != nil {
	// 	return nil, err
	// }
	// transaction := s.TransactionMapper.EntityToTransactionResponse(result)
	// return transaction, nil
	return nil, nil
}

// @Summary Hủy cọc
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transactionId path int true "Transaction ID"
// @Router /v2/bdspro/v2/product/deposite/{transactionId} [delete]
func (s *ProductHandler) DeleteDeposite(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.Response, error) {
	err := s.UC.CancelDeposite(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Xuất danh sách sản phẩm ra Excel
// @Description API này dùng để xuất dữ liệu sản phẩm ra file Excel
// @Tags User: Sản phẩm
// @Accept json
// @Produce application/octet-stream
// @Security BearerAuth
// @Param filter query string false "Bộ lọc sản phẩm (nếu có)"
// @Success 200 {file} file "Excel file"
// @Failure 400 {object} _routes.ResponseDTO
// @Failure 500 {object} _routes.ResponseDTO
// @Router /v2/bdspro/v2/product/export [post]
func (s *ProductHandler) ExportProduct(context.Context, *bdspropb.IdRequest) (*bdspropb.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExportProduct not implemented")
}

// @Summary Lấy danh sách sản phẩm của user hiện tại
// @Description Returns a list of products based on search filters
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name query string false "Product name"
// @Param code query string false "Product code"
// @Param area query number false "Product area"
// @Param price query integer false "Product price"
// @Param parentId query integer false "null=>all; 0=>parent; >0=>childrent"
// @Param description query string false "Product description"
// @Param note query string false "Additional notes"
// @Param limit query int false "Limit results" default(10)
// @Param skip query int false "Skip results" default(0)
// @Param sort query string false "Sort order" default("price_asc")
// @Router /v2/bdspro/v2/product/me [get]
func (s *ProductHandler) ProductMe(ctx context.Context, req *sharepb.SyncRequest) (*bdspropb.SearchResponse, error) {
	key := s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductMe, req.Id)
	updated := s.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	if !updated && req.Page == 0 {
		return &bdspropb.SearchResponse{}, nil
	}
	body := s.ProductMapper.PbSearchRequestToDTO(&bdspropb.ProductSearch{
		Page: req.Page,
		Size: req.Size,
	})
	body.Sync = true

	result, total, err := s.UC.GetProductOfCurrentUser(ctx, body)
	if err != nil {
		return nil, err
	}

	if req.Page == 0 && len(result) > 0 {
		s.SyncProvider.PutTimestamp(ctx, key, result[0].UpdatedAt.UnixMilli())
	}

	products := s.ProductMapper.MapProductPbListResponse(result)

	// products := make([]*bdspropb.ProductResponse, 0)
	// copier.Copy(&products, &result)
	return &bdspropb.SearchResponse{
		Data:          products,
		TotalElements: total,
	}, nil
}

// @Summary Tạo sản phẩm mới
// @Description API để tạo một sản phẩm mới
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body dto.ProductSaveRequest true "Thông tin sản phẩm"
// @Success 201 {object} dto.CreateResponse "Sản phẩm được tạo thành công"
// @Router /v2/bdspro/v2/product/new [post]
func (s *ProductHandler) ProductNew(ctx context.Context, req *bdspropb.ProductSaveRequest) (*bdspropb.Response, error) {
	body := s.ProductMapper.ProductSavePbToDTO(req)
	if err := s.ProductValidator.ValidateCreateProduct(body); err != nil {
		return nil, err
	}

	result, err := s.UC.CreateProduct(ctx, body)
	if err != nil {
		return nil, err
	}
	id := uint64(0)
	if result.Product != nil {
		id = result.Product.ID
	}
	return &bdspropb.Response{
		Id:      id,
		Message: "success",
	}, nil
}

// @Summary Cập nhật trạng thái thuê sản phẩm
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.StatusUpdate true "Thông tin sản phẩm cần cập nhật"
// @Router /v2/bdspro/v2/product/rent-status [put]
func (s *ProductHandler) RentStatus(ctx context.Context, req *bdspropb.StatusRequest) (*bdspropb.Response, error) {
	if req.Id == 0 || req.Status == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id, status is required")
	}
	body := dto.StatusUpdate{
		ProductID:  &req.Id,
		RentStatus: enums.EProductRentStatus(req.Status),
		ContactID:  &req.ContactId,
		Note:       req.Note,
		Amount:     req.Amount,
		Phone:      req.Phone,
	}

	err := s.UC.RentStatusUpdate(ctx, body)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Cập nhật trạng thái thuê sản phẩm
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.VisibilityUpdate true "Thông tin sản phẩm cần cập nhật"
// @Router /v2/bdspro/v2/product/rent-visibility [put]
func (s *ProductHandler) RentVisibility(ctx context.Context, req *bdspropb.VisibilityRequest) (*bdspropb.Response, error) {
	if req.Id == 0 || req.Visibility == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id, visibility is required")
	}
	body := dto.VisibilityUpdate{
		ProductID:  &req.Id,
		Visibility: enums.EVisibility(req.Visibility),
	}

	err := s.UC.RentVisibilityUpdate(ctx, body)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Cập nhật trạng thái bán sản phẩm
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.StatusUpdate true "Thông tin sản phẩm cần cập nhật"
// @Router /v2/bdspro/v2/product/sale-status [put]
func (s *ProductHandler) SaleStatus(ctx context.Context, req *bdspropb.StatusRequest) (*bdspropb.Response, error) {
	if req.Id == 0 || req.Status == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id, status is required")
	}
	body := dto.StatusUpdate{
		ProductID:  &req.Id,
		SaleStatus: enums.EProductSaleStatus(req.Status),
		ContactID:  &req.ContactId,
		Note:       req.Note,
		Amount:     req.Amount,
		Phone:      req.Phone,
	}

	err := s.UC.SaleStatusUpdate(ctx, body)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Cập nhật trạng thái bán sản phẩm
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.VisibilityUpdate true "Thông tin sản phẩm cần cập nhật"
// @Router /v2/bdspro/v2/product/sale-visibility [put]
func (s *ProductHandler) SaleVisibility(ctx context.Context, req *bdspropb.VisibilityRequest) (*bdspropb.Response, error) {
	if req.Id == 0 || req.Visibility == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id, visibility is required")
	}
	body := dto.VisibilityUpdate{
		ProductID:  &req.Id,
		Visibility: enums.EVisibility(req.Visibility),
	}

	err := s.UC.SaleVisibilityUpdate(ctx, body)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Gợi ý các trường thông tin
// @Description Returns a list of products based on search filters
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name query string false "Product name"
// @Param code query string false "Product code"
// @Param area query number false "Product area"
// @Param price query integer false "Product price"
// @Param description query string false "Product description"
// @Param note query string false "Additional notes"
// @Param limit query int false "Limit results" default(10)
// @Param skip query int false "Skip results" default(0)
// @Param sort query string false "Sort order" default("price_asc")
// @Router /v2/bdspro/v2/product/suggest [post]
func (s *ProductHandler) Suggest(ctx context.Context, req *bdspropb.SuggestRequest) (*bdspropb.ProductResponse, error) {
	// products, err := s.UC.SuggestFields(ctx, req.Content)
	// if err != nil {
	// 	return nil, err
	// }

	// // Tạo ProductResponse từ TotalSearchParser
	// productResponse := s.ProductMapper.MapTotalSearchParserToResponse(products)

	// return productResponse, nil
	return nil, nil
}

// @Summary Gợi ý các trường thông tin sử dụng DeepSeek AI
// @Description Phân tích văn bản mô tả bất động sản bằng DeepSeek AI và trả về các trường thông tin được gợi ý
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body bdspropb.SuggestRequest true "Nội dung mô tả sản phẩm"
// @Success 200 {object} bdspropb.ProductResponse
// @Router /v2/bdspro/v2/product/suggest/deepseek [post]
func (s *ProductHandler) SuggestDeepseek(ctx context.Context, req *bdspropb.SuggestRequest) (*bdspropb.ProductResponse, error) {
	if req.Content == "" {
		return nil, status.Errorf(codes.InvalidArgument, "content is required")
	}

	// Gọi usecase để phân tích bằng DeepSeek
	products, err := s.UC.SuggestFieldsWithDeepseek(ctx, req.Content)
	if err != nil {
		return nil, err
	}

	// Tạo ProductResponse từ TotalSearchParser
	productResponse := s.ProductMapper.MapTotalSearchParserToResponse(products)

	return productResponse, nil
}

// @Summary Lấy chi tiết sản phẩm
// @Description API để lấy chi tiết sản phẩm
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID của sản phẩm"
// @Success 200 {object} bdspropb.IdRequest "Chi tiết sản phẩm"
// @Summary Lấy timestamps các loại data của sản phẩm (chỉ đọc Redis)
// @Description API check sync - lấy timestamp từ Redis cho từng key, FE so sánh với cache và chỉ gọi API detail khi có thay đổi. Không gọi DB.
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} bdspropb.ProductTimestampsResponse "Map key -> timestamp (ms)"
// @Router /v2/bdspro/v2/product/{id}/timestamp [get]
func (s *ProductHandler) GetProductTimestamps(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.ProductTimestampsResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "product id is required")
	}

	redisKeys := []string{
		s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductDetail, req.Id),
		s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductPriceHist, req.Id),
		s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductDistrs, req.Id),
		s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductPosts, req.Id),
		s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductDeals, req.Id),
		s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductNotes, req.Id),
		s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductAppointments, req.Id),
		s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductContacts, req.Id),
		s.SyncProvider.GetKey(ctx, _utils.SyncKeyPropertyHistories, req.Id),
	}

	vals, err := s.SyncProvider.MGet(ctx, redisKeys)
	if err != nil {
		return nil, _errors.InternalServerException("error: %v", err.Error())
	}

	names := []string{
		"detail",
		"price_history",
		"distributions",
		"posts",
		"deals",
		"notes",
		"appointments",
		"contacts",
		"histories",
	}

	timestamps := make(map[string]int64, len(names))

	for i, v := range vals {
		if v == nil {
			timestamps[names[i]] = 0
			continue
		}

		ts, _ := strconv.ParseInt(v.(string), 10, 64)
		timestamps[names[i]] = ts
	}

	return &bdspropb.ProductTimestampsResponse{
		Timestamps: timestamps,
	}, nil
}

func (s *ProductHandler) GetDetail(ctx context.Context, req *sharepb.SyncRequest) (*bdspropb.ProductResponse, error) {
	key := s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductDetail, req.Id)
	updated := s.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	s.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &bdspropb.ProductResponse{}, nil
	}

	product, err := s.UC.Detail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	// var originProfile *sharepb.OriginProfile
	// if product.SourceContactId != nil && *product.SourceContactId > 0 {
	// 	pf, err := s.CrmProvider.GetContactByOriginId(ctx, *product.SourceContactId)
	// 	if err != nil {
	// 		return nil, err
	// 	} else if pf.OriginProfile != nil {
	// 		originProfile = pf.OriginProfile
	// 	}
	// }
	// if originProfile != nil {
	// 	resp.OriginProfileContact = s.ProductMapper.MapOriginProfilePb(originProfile)
	// }

	resp := s.ProductMapper.MapDetailProductPb(ctx, product)
	return resp, nil
}

// // @Summary Lấy danh sách lịch sử của sản phẩm
// // @Description Returns a list of history of a product
// // @Tags User: Sản phẩm
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "Product ID"
// // @Router /v2/bdspro/v2/product/{id}/histories [get]
// func (s *ProductService) GetHistories(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.HistoryResponse, error) {
// 	query := dto.ProductHistorySearch{}
// 	if err := copier.Copy(&query, req); err != nil {
// 		return nil, err
// 	}

// 	histories, total, err := s.UC.History(ctx, req.Id, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	historyResponse := make([]*bdspropb.History, 0)
// 	copier.Copy(&historyResponse, &histories)
// 	return &bdspropb.HistoryResponse{
// 		Data:  historyResponse,
// 		Total: uint64(total),
// 	}, nil
// }

// @Summary Lấy danh sách thành viên của sản phẩm
// @Description Returns a list of members of a product
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Router /v2/bdspro/v2/product/{id}/members [get]
func (s *ProductHandler) GetMembers(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.MemberResponse, error) {
	query := dto.SharingAccessSearch{}
	if err := copier.Copy(&query, req); err != nil {
		return nil, err
	}
	members, err := s.UC.Members(ctx, req.Id, query)
	if err != nil {
		return nil, err
	}
	memberResponse := make([]*bdspropb.ProductAccess, 0)
	copier.Copy(&memberResponse, &members)
	return &bdspropb.MemberResponse{
		Data: memberResponse,
	}, nil
}

// @Summary Lấy danh sách bài viết liên quan
// @Description Returns a list of posts related to a product
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Router /v2/bdspro/v2/product/{id}/posts [get]
func (s *ProductHandler) GetPosts(ctx context.Context, req *sharepb.SyncRequest) (*bdspropb.PostResponse, error) {
	key := s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductPosts, req.Id)
	updated := s.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	s.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	// client mặc định sẽ truyền page = 0 GET /v2/deals?page=0
	if !updated && req.Page == 0 {
		return &bdspropb.PostResponse{}, nil
	}

	query := dto.PostLinkSearch{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		ProductID: &req.Id,
	}
	posts, err := s.UC.Posts(ctx, req.Id, query)
	if err != nil {
		return nil, err
	}
	if req.Page == 0 {
		time := time.Now().UnixMilli() - 24*60*60*1000
		if len(posts) > 0 {
			time = posts[0].ProductUpdatedAt.UnixMilli()
		}
		s.SyncProvider.PutTimestamp(ctx, key, time)
	}
	postResponse := s.PostMapper.PostToPbList(&posts)
	// postResponse := make([]*bdspropb.PostItem, 0)
	// copier.Copy(&postResponse, &posts)
	return &bdspropb.PostResponse{
		Data: postResponse,
	}, nil
}

// @Summary Lấy danh sách thương vụ của sản phẩm
// @Description Returns a list of deals related to a product
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Router /v2/bdspro/v2/product/{id}/deals [get]
func (s *ProductHandler) GetDeals(ctx context.Context, req *sharepb.SyncRequest) (*bdspropb.GetOrganizationDealsResponse, error) {
	key := s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductDeals, req.Id)
	updated := s.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	s.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	// client mặc định sẽ truyền page = 0 GET /v2/deals?page=0
	if !updated && req.Page == 0 {
		return &bdspropb.GetOrganizationDealsResponse{}, nil
	}
	pagable := _dto.Pagable{
		Page: req.Page,
		Size: req.Size,
	}

	deals, total, _, err := s.UC.GetDealsByProduct(ctx, req.Id, pagable)
	if err != nil {
		return nil, err
	}
	// if req.Page == 0 {
	// 	time := time.Now().UnixMilli() - 24*60*60*1000
	// 	if len(deals) > 0 {
	// 		time = deals[0].UpdatedAt.UnixMilli()
	// 	}
	// 	s.SyncProvider.PutTimestamp(ctx, key, time)
	// }
	dealItems := s.DealMapper.EntityToGetOrganizationDealsResponse(deals, uint32(total))

	return dealItems, nil
}

// @Summary Lấy previousId và nextId của sản phẩm
// @Description Trả về ID sản phẩm trước và sau trong danh sách (sắp xếp theo created_at DESC)
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Router /v2/bdspro/v2/product/{id}/navigation [get]
func (s *ProductHandler) GetProductNavigation(ctx context.Context, req *bdspropb.ProductSearch) (*bdspropb.ProductNavigationResponse, error) {
	if req.Id == nil || *req.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	previousID, nextID, err := s.ProductRepo.GetProductNavigation(ctx, *req.Id)
	if err != nil {
		return nil, err
	}

	return &bdspropb.ProductNavigationResponse{
		PreviousId: previousID,
		NextId:     nextID,
	}, nil
}

// @Summary Cập nhật sản phẩm
// @Description API để cập nhật thông tin sản phẩm
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "Product ID"
// @Param body body dto.UpdateProductRequest true "Thông tin sản phẩm cần cập nhật"
// @Success 200 {object} bdspropb.Response "Cập nhật thành công"
// @Router /v2/bdspro/v2/product/edit/{id} [put]
func (s *ProductHandler) UpdateProduct(ctx context.Context, req *bdspropb.UpdateProductRequest) (*bdspropb.Response, error) {
	// Convert proto to UpdateProductRequest DTO
	updateDto := s.ProductMapper.UpdateProductPbToDTO(req)

	// Validate UpdateProductRequest
	if err := s.ProductValidator.ValidateUpdateProduct(updateDto); err != nil {
		return nil, err
	}

	// Call usecase Update
	_, err := s.UC.Update(ctx, req.Id, updateDto)
	if err != nil {
		return nil, err
	}

	return &bdspropb.Response{
		Id:      req.Id,
		Message: "success",
	}, nil
	// id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	// body, err := _utils.ParseBody[dto.ProductSaveRequest](c)
	// if err != nil {
	// 	_routes.RouteResult(c, nil, err)
	// 	return
	// }
	// result, err := route.UC.Update(c, id, body)
	// _routes.RouteResult(c, result, err)
}

// @Summary Xóa sản phẩm
// @Description API
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Param productId path int true "Product ID"
// @Security BearerAuth
// @Router /v2/bdspro/v2/product/{id} [delete]
func (s *ProductHandler) DeleteProduct(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.Response, error) {
	err := s.UC.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	s.SyncProvider.PutTimestamp(ctx, s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductDetail, req.Id), time.Now().UnixMilli())
	s.SyncProvider.PutTimestamp(ctx, s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductMe, req.Id), time.Now().UnixMilli())
	return &bdspropb.Response{
		Message: "ok",
	}, nil
}

// @Summary Lấy thông tin về diện tích sản phẩm
// @Description Returns information about the area of a product
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Router /v2/bdspro/v2/product/area-info/{id} [get]
func (s *ProductHandler) AreaInfo(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.AreaInfoResponse, error) {
	areaInfo, err := s.UC.InfoArea(ctx, &req.Id)
	if err != nil {
		return nil, err
	}
	return &bdspropb.AreaInfoResponse{
		TotalArea:    areaInfo.Total,
		AvaiableArea: areaInfo.Avaiable,
		SplitArea:    areaInfo.Area,
	}, nil
}

func (s *ProductHandler) GetDetailByIds(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.ProductInternalResponse, error) {
	products, err := s.ProductRepo.GetDetailByIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	productResponse := s.ProductMapper.MapProductPbList(products)
	return &bdspropb.ProductInternalResponse{
		Data: productResponse,
	}, nil
}

func (s *ProductHandler) GetProductAttachmentByIds(ctx context.Context, req *bdspropb.IdRequest) (*bdspropb.ProductAttachmentResponse, error) {
	products, err := s.ProductRepo.GetDetailByIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	productResponse := mapper.MapProductAttachmentPbList(ctx, s.ProductMapper.CodeDataUsecase, products)
	return &bdspropb.ProductAttachmentResponse{
		Data: productResponse,
	}, nil
}

// @Summary Get group products
// @Description Get list of products in a group based on sharing access
// @Tags Product
// @Accept json
// @Produce json
// @Param groupId path int true "Group ID"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Param text query string false "Search text"
// @Param parentId query int false "Parent ID"
// @Param propertyTypeId query int false "Property Type ID"
// @Param docTypeId query int false "Document Type ID"
// @Param projectId query int false "Project ID"
// @Param provinceId query int false "Province ID"
// @Param districtId query int false "District ID"
// @Param wardId query int false "Ward ID"
// @Param transactionType query int false "Transaction Type"
// @Param saleStatus query int false "Sale Status"
// @Param rentStatus query int false "Rent Status"
// @Param saleVisibility query int false "Sale Visibility"
// @Param rentVisibility query int false "Rent Visibility"
// @Param sourceType query int false "Source Type"
// @Param archived query bool false "Archived"
// @Success 200 {object} bdspropb.SearchResponse
// @Router /v2/bdspro/v2/product/group/{groupId} [get]
func (s *ProductHandler) GetGroupProducts(ctx context.Context, req *bdspropb.GroupProductSearch) (*bdspropb.SearchResponse, error) {
	// Convert protobuf request to DTO
	searchReq := &dto.GroupProductSearchRequest{
		GroupID:        req.GroupId,
		Text:           req.Text,
		ParentID:       req.ParentId,
		PropertyTypeID: req.PropertyTypeId,
		DocTypeID:      req.DocTypeId,
		ProjectID:      req.ProjectId,
		ProvinceID:     req.ProvinceId,
		// DistrictID:       req.DistrictId,
		WardID:           req.WardId,
		TransactionType:  req.TransactionType,
		SaleStatus:       req.SaleStatus,
		RentStatus:       req.RentStatus,
		SaleVisibility:   req.SaleVisibility,
		RentVisibility:   req.RentVisibility,
		SourceType:       req.SourceType,
		TransactionPrice: req.TransactionPrice,
		SalePrice:        req.SalePrice,
		RentPrice:        req.RentPrice,
		SaleCommission:   req.SaleCommission,
		RentCommission:   req.RentCommission,
		RentPaymentCycle: req.RentPaymentCycle,
		Archived:         req.Archived,
	}

	// Set pagination
	if req.Page > 0 {
		searchReq.Page = req.Page
	}
	if req.Size > 0 {
		searchReq.Size = req.Size
	}

	// Get products from usecase
	products, total, err := s.UC.GetGroupProducts(ctx, searchReq)
	if err != nil {
		return nil, err
	}

	// Convert domain products to protobuf response
	productResponses := s.ProductMapper.MapProductPbList(products)

	return &bdspropb.SearchResponse{
		Data:          productResponses,
		TotalElements: total,
	}, nil
}

func (s *ProductHandler) CurrentOrganizationProducts(ctx context.Context, req *bdspropb.ProductSearch) (*bdspropb.SearchResponse, error) {
	payload := s.ProductMapper.PbSearchRequestToDTO(req)
	products, total, err := s.UC.GetProductOfCurrentOrganization(ctx, payload)
	if err != nil {
		return nil, err
	}

	productResponses := s.ProductMapper.MapProductPbListResponse(products)
	return &bdspropb.SearchResponse{
		Data:          productResponses,
		TotalElements: total,
	}, nil
}

// @Summary Lấy danh sách sản phẩm của user
// @Description Lấy danh sách sản phẩm của một user cụ thể
// @Tags Product
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Param text query string false "Search text"
// @Param parentId query int false "Parent ID"
// @Param propertyTypeId query int false "Property Type ID"
// @Param docTypeId query int false "Document Type ID"
// @Param projectId query int false "Project ID"
// @Param provinceId query int false "Province ID"
// @Param districtId query int false "District ID"
// @Param wardId query int false "Ward ID"
// @Param transactionType query int false "Transaction Type"
// @Param saleStatus query int false "Sale Status"
// @Param rentStatus query int false "Rent Status"
// @Param saleVisibility query int false "Sale Visibility"
// @Param rentVisibility query int false "Rent Visibility"
// @Param sourceType query int false "Source Type"
// @Param archived query bool false "Archived"
// @Success 200 {object} bdspropb.SearchResponse
// @Router /v2/bdspro/v2/product/user/{userId} [get]
func (s *ProductHandler) GetUserProducts(ctx context.Context, req *bdspropb.UserProductSearch) (*bdspropb.SearchResponse, error) {
	// Convert protobuf request to DTO
	searchReq := &dto.UserProductSearchRequest{
		UserID:           req.UserId,
		Text:             req.Text,
		ParentID:         req.ParentId,
		PropertyTypeID:   req.PropertyTypeId,
		DocTypeID:        req.DocTypeId,
		ProjectID:        req.ProjectId,
		ProvinceID:       req.ProvinceId,
		WardID:           req.WardId,
		TransactionType:  req.TransactionType,
		SaleStatus:       req.SaleStatus,
		RentStatus:       req.RentStatus,
		SaleVisibility:   req.SaleVisibility,
		RentVisibility:   req.RentVisibility,
		SourceType:       req.SourceType,
		TransactionPrice: req.TransactionPrice,
		SalePrice:        req.SalePrice,
		RentPrice:        req.RentPrice,
		SaleCommission:   req.SaleCommission,
		RentCommission:   req.RentCommission,
		RentPaymentCycle: req.RentPaymentCycle,
		Archived:         req.Archived,
	}

	// Set pagination
	if req.Page > 0 {
		searchReq.Page = req.Page
	}
	if req.Size > 0 {
		searchReq.Size = req.Size
	}

	// Get products from usecase
	products, total, err := s.UC.GetUserProducts(ctx, searchReq)
	if err != nil {
		return nil, err
	}

	// Convert domain products to protobuf response
	productResponses := s.ProductMapper.MapProductPbList(products)

	return &bdspropb.SearchResponse{
		Data:          productResponses,
		TotalElements: total,
	}, nil
}

func (s *ProductHandler) CreateOrganizationProduct(ctx context.Context, req *bdspropb.ProductSaveRequest) (*bdspropb.Response, error) {
	dto := s.ProductMapper.ProductSavePbToDTO(req)
	result, err := s.UC.CreateOrganizationProduct(ctx, dto)
	if err != nil {
		return nil, err
	}
	id := uint64(0)
	if result.Product != nil {
		id = result.Product.ID
	}
	return &bdspropb.Response{
		Id:      id,
		Message: "success",
	}, nil
}

// @Summary Lấy thống kê tổng hợp sản phẩm của user
// @Description API này cho phép lấy thống kê tổng số sản phẩm, số sản phẩm đang bán và tổng giá trị của user hiện tại với filters
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param propertyTypeIds query string false "Property type IDs (comma-separated)"
// @Param projectIds query string false "Project IDs (comma-separated)"
// @Param transactionTypes query string false "Transaction types (comma-separated)"
// @Param saleStatuses query string false "Sale statuses (comma-separated)"
// @Param rentStatuses query string false "Rent statuses (comma-separated)"
// @Param saleVisibilities query string false "Sale visibilities (comma-separated)"
// @Param rentVisibilities query string false "Rent visibilities (comma-separated)"
// @Param sourceTypeIds query string false "Source type IDs (comma-separated)"
// @Param text query string false "Search text"
// @Param parentId query integer false "null=>all; 0=>parent; >0=>childrent"
// @Param archived query boolean false "Archived status"
// @Success 200 {object} bdspropb.ProductSummaryResponse
// @Router /v2/bdspro/v2/product/summary [get]
func (s *ProductHandler) ProductSummary(ctx context.Context, req *bdspropb.ProductSearch) (*bdspropb.ProductSummaryResponse, error) {
	body := s.ProductMapper.PbSearchRequestToDTO(req)

	summary, err := s.UC.GetProductSummary(ctx, body)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get product summary: %v", err)
	}

	return &bdspropb.ProductSummaryResponse{
		TotalProducts:   summary.TotalProducts,
		SellingProducts: summary.SellingProducts,
		TotalValue:      summary.TotalValue,
	}, nil
}

// @Summary Liên kết nhiều sản phẩm với tài sản
// @Description API này cho phép liên kết nhiều sản phẩm với một tài sản
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body bdspropb.LinkProductAssetRequest true "Thông tin liên kết (nhiều productIds và 1 assetId)"
// @Router /v2/bdspro/v2/product/link-asset [post]
func (s *ProductHandler) LinkProductAsset(ctx context.Context, req *bdspropb.LinkProductAssetRequest) (*bdspropb.Response, error) {
	if len(req.ProductIds) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "productIds không được để trống")
	}
	if req.AssetId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "assetId is required")
	}

	linkReq := &dto.LinkProductAssetRequest{
		ProductIDs: req.ProductIds,
		AssetID:    req.AssetId,
	}

	err := s.ProductAssetUC.LinkProductAsset(ctx, linkReq)
	if err != nil {
		return nil, err
	}

	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Hủy liên kết sản phẩm với tài sản
// @Description API này cho phép hủy liên kết sản phẩm với tài sản
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body bdspropb.UnlinkProductAssetRequest true "Thông tin hủy liên kết"
// @Router /v2/bdspro/v2/product/unlink-asset [post]
func (s *ProductHandler) UnlinkProductAsset(ctx context.Context, req *bdspropb.UnlinkProductAssetRequest) (*bdspropb.Response, error) {
	if req.ProductId == 0 || req.AssetId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "productId and assetId are required")
	}

	unlinkReq := &dto.UnlinkProductAssetRequest{
		ProductID: req.ProductId,
		AssetID:   req.AssetId,
	}

	err := s.ProductAssetUC.UnlinkProductAsset(ctx, unlinkReq)
	if err != nil {
		return nil, err
	}

	return &bdspropb.Response{
		Message: "success",
	}, nil
}

// @Summary Lấy danh sách tài sản theo sản phẩm
// @Description API này cho phép lấy danh sách tài sản được liên kết với sản phẩm
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Router /v2/bdspro/v2/product/{id}/assets [get]
func (s *ProductHandler) GetAssetsByProduct(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.ProductAssetListResponse, error) {
	if req.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	// Lấy danh sách asset đầy đủ từ repo
	assets, err := s.ProductMapper.AssetRepo.GetAssetsByProductID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Map danh sách asset sang proto
	assetList := s.AssetMapper.MapAssetPbList(assets)

	return &bdspropb.ProductAssetListResponse{
		Data: assetList,
	}, nil
}

// @Summary Lấy lịch sử cập nhật giá
// @Description API này cho phép lấy lịch sử cập nhật giá của sản phẩm
// @Tags User: Sản phẩm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param page query int false "Số trang" default(1)
// @Param size query int false "Số lượng bản ghi trên trang" default(20)
// @Router /v2/bdspro/v2/product/{id}/price-history [get]

func (h *ProductHandler) ApplyPriceUpdate(
	ctx context.Context,
	req *bdspropb.ApplyPriceUpdateRequest,
) (*bdspropb.ApplyPriceUpdateResponse, error) {
	dtoReq := h.ProductMapper.ApplyPriceUpdatePbToDTO(req)

	resp, err := h.UC.ApplyPriceUpdate(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return h.ProductMapper.ApplyPriceUpdateDTOToPb(resp), nil
}

func (h *ProductHandler) UpdateProductSource(
	ctx context.Context,
	req *bdspropb.UpdateProductSourceRequest,
) (*bdspropb.Response, error) {
	dtoReq := h.ProductMapper.PbToUpdateProductSourceDTO(req)
	if err := h.UC.UpdateProductSource(ctx, dtoReq); err != nil {
		return nil, err
	}

	return &bdspropb.Response{
		Id:      req.ProductId,
		Message: "successfully",
	}, nil
}

func (s *ProductHandler) GetPriceHistory(
	ctx context.Context,
	req *bdspropb.PriceHistoryRequest,
) (*bdspropb.PriceHistoryResponse, error) {
	if req.ProductId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "product_id is required")
	}
	key := s.SyncProvider.GetKey(ctx, _utils.SyncKeyProductPriceHist, req.ProductId)
	updated := s.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	s.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &bdspropb.PriceHistoryResponse{}, nil
	}

	items, total, _, err := s.UC.GetPriceHistory(ctx, &dto.PriceHistorySearch{
		ProductID: req.ProductId,
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
	})
	if err != nil {
		return nil, err
	}
	// if req.Page == 0 {
	// 	time := time.Now().UnixMilli() - 24*60*60*1000
	// 	if len(items) > 0 {
	// 		time = productUpdatedTime.UnixMilli()
	// 	}
	// 	s.SyncProvider.PutTimestamp(ctx, key, req.Timestamp)
	// }

	return &bdspropb.PriceHistoryResponse{
		Data:  s.ProductPriceMapper.DTOListToPb(items),
		Total: total,
	}, nil
}

func (h *ProductHandler) GetProductPriceDetail(
	ctx context.Context,
	req *bdspropb.ProductPriceDetailRequest,
) (*bdspropb.ProductPriceDetailResponse, error) {

	if req.ProductId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "productId is required")
	}

	result, err := h.UC.GetProductPricedDetail(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}

	return &bdspropb.ProductPriceDetailResponse{
		Data: h.ProductPriceMapper.DTOToPbDetail(result),
	}, nil
}

func (h *ProductHandler) ArchiveProducts(ctx context.Context, req *bdspropb.ArchiveProductsRequest) (*bdspropb.Response, error) {
	dtoReq := h.ProductMapper.PbToArchiveProductsDTO(req)
	err := h.UC.ArchiveProducts(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to archive products: %v", err)
	}

	action := "archived"
	if !req.Archived {
		action = "unarchived"
	}

	return &bdspropb.Response{
		Message: "Products " + action + " successfully",
	}, nil
}
