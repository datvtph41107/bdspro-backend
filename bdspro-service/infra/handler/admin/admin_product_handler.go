package admin_handler

import (
	"bdspro/infra/client"
	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	admin_usecases "bdspro/internal/usecases/admin"
	_routes "common/routes"
	"context"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
)

type AdminProductHandler struct {
	bdspropb.UnimplementedAdminProductServiceServer
	AdminProductService *admin_usecases.AdminProductUsecase
	ProductMapper       *mapper.ProductMapper
	AuthGrpcClient      *client.AuthClient
	UserClient          *client.UserClient
}

func NewAdminProductHandler(
	uc *admin_usecases.AdminProductUsecase,
	productMapper *mapper.ProductMapper,
	authGrpcClient *client.AuthClient,
	userClient *client.UserClient,
) *AdminProductHandler {
	return &AdminProductHandler{
		AdminProductService: uc,
		ProductMapper:       productMapper,
		AuthGrpcClient:      authGrpcClient,
		UserClient:          userClient,
	}
}

// @Summary Get products
// @Description Get products
// @Tags AdminProduct
// @Accept json
// @Produce json
// @Param page query int true "Page number"
// @Param size query int true "Page size"
// @Success 200 {object} bdspropb.SearchResponse
// @Router /admin/products [get]
func (h *AdminProductHandler) GetProducts(ctx context.Context, req *bdspropb.AdminProductSearchRequest) (*bdspropb.AdminProductSearchResponse, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_SP_XEM"}); err != nil {
		return nil, err
	}

	reqDto := h.ProductMapper.AdminProductSearchRequestToDTO(req)
	result, total, err := h.AdminProductService.GetList(ctx, reqDto)

	if err != nil {
		return nil, err
	}

	// Convert sang AdminProductItem
	adminProducts := h.ProductMapper.MapAdminProductPbList(ctx, result)

	return &bdspropb.AdminProductSearchResponse{
		Data:  adminProducts,
		Total: total,
	}, nil
}

// @Summary Create a new product
// @Description Create a new product
// @Tags AdminProduct
// @Accept json
// @Produce json
// @Param product body bdspropb.ProductSaveRequest true "Product to create"
// @Success 200 {object} bdspropb.ProductResponse
// @Router /admin/products [post]
func (h *AdminProductHandler) CreateProduct(ctx context.Context, req *bdspropb.ProductSaveRequest) (*bdspropb.ProductResponse, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_SP_TAO"}); err != nil {
		return nil, err
	}

	// Validate required fields
	if req.Name == "" {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Tên sản phẩm không được để trống",
		}
	}

	if req.OwnerId == 0 {
		return nil, &_routes.Except{
			Code:    400,
			Message: "OwnerId không được để trống",
		}
	}

	// Convert proto request to DTO
	product := h.ProductMapper.ProductSavePbToDTO(req)
	if product == nil {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Dữ liệu sản phẩm không hợp lệ",
		}
	}

	// Create product using AdminProductService
	createdProduct, err := h.AdminProductService.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	// Return response
	return h.ProductMapper.MapDetailProductPb(ctx, createdProduct), nil
}

// @Summary Update a product
// @Description Update a product
// @Tags AdminProduct
// @Accept json
// @Produce json
// @Param product body bdspropb.ProductSaveRequest true "Product to update"
// @Success 200 {object} bdspropb.ProductResponse
// @Router /admin/products [put]
func (h *AdminProductHandler) UpdateProduct(ctx context.Context, req *bdspropb.ProductSaveRequest) (*bdspropb.ProductResponse, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_SP_SUA"}); err != nil {
		return nil, err
	}
	// return h.AdminProductService.UpdateProduct(ctx, req)
	return nil, nil
}

// @Summary Delete a product
// @Description Delete a product
// @Tags AdminProduct
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} bdspropb.Response
// @Router /admin/products/{id} [delete]
func (h *AdminProductHandler) DeleteProduct(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.Response, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_SP_XOA"}); err != nil {
		return nil, err
	}
	// return h.AdminProductService.DeleteProduct(ctx, req)
	return nil, nil
}

// @Summary Get product by ID
// @Description Get product details by ID for admin
// @Tags AdminProduct
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} bdspropb.ProductResponse
// @Router /admin/products/{id} [get]
func (h *AdminProductHandler) GetProductById(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.ProductResponse, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_SP_XEM"}); err != nil {
		return nil, err
	}
	// return h.AdminProductService.GetProductById(ctx, req)
	product, err := h.AdminProductService.GetDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return h.ProductMapper.MapDetailProductPb(ctx, product), nil
}

// @Summary Get product history
// @Description Get product history by ID for admin
// @Tags AdminProduct
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} bdspropb.ProductHistoryResponse
// @Router /admin/products/{id}/history [get]
func (h *AdminProductHandler) GetProductHistory(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.ProductHistoryResponse, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_SP_XEM"}); err != nil {
		return nil, err
	}

	// Get product history using ProductUsecase
	histories, total, err := h.AdminProductService.ProductUsecase.History(ctx, req.Id, dto.ProductHistorySearch{})
	if err != nil {
		return nil, err
	}

	// Convert to protobuf response
	historyItems := make([]*bdspropb.ProductHistory, len(*histories))
	for i, history := range *histories {
		historyItems[i] = h.ProductMapper.MapProductHistoryPb(ctx, &history)
	}

	return &bdspropb.ProductHistoryResponse{
		Data:  historyItems,
		Total: uint64(total),
	}, nil
}

// @Summary Update product note
// @Description Update product note by ID for admin
// @Tags AdminProduct
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param body body bdspropb.UpdateProductNoteRequest true "Note update request"
// @Success 200 {object} bdspropb.Response
// @Router /admin/products/{id}/note [put]
func (h *AdminProductHandler) UpdateProductNote(ctx context.Context, req *bdspropb.UpdateProductNoteRequest) (*bdspropb.Response, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_SP_SUA"}); err != nil {
		return nil, err
	}

	// Update product note using AdminProductService
	err := h.AdminProductService.UpdateNote(ctx, req.Id, req.Note)
	if err != nil {
		return nil, err
	}

	return &bdspropb.Response{
		Message: "Cập nhật ghi chú thành công",
	}, nil
}
