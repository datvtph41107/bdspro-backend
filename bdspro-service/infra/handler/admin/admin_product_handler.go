package admin_handler

import (
	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	admin_usecases "bdspro/internal/usecases/admin"
	_dto "common/dto"
	_routes "common/routes"
	"context"
	"pb/clients"
	pb_bdspro "pb/types/bdspro"
	sharepb "pb/types/shared"
)

type AdminProductHandler struct {
	pb_bdspro.UnimplementedAdminProductServiceServer
	AdminProductService *admin_usecases.AdminProductUsecase
	ProductMapper       *mapper.ProductMapper
	AuthGrpcClient      *clients.AuthGrpcClient
}

func NewAdminProductHandler(
	uc *admin_usecases.AdminProductUsecase,
	productMapper *mapper.ProductMapper,
	authGrpcClient *clients.AuthGrpcClient,
) *AdminProductHandler {
	return &AdminProductHandler{
		AdminProductService: uc,
		ProductMapper:       productMapper,
		AuthGrpcClient:      authGrpcClient,
	}
}

// @Summary Get products
// @Description Get products
// @Tags AdminProduct
// @Accept json
// @Produce json
// @Param page query int true "Page number"
// @Param size query int true "Page size"
// @Success 200 {object} pb_bdspro.SearchResponse
// @Router /admin/products [get]
func (h *AdminProductHandler) GetProducts(ctx context.Context, req *sharepb.IdRequest) (*pb_bdspro.SearchResponse, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_SP_XEM"}); err != nil {
		return nil, err
	}

	result, total, err := h.AdminProductService.GetList(ctx,
		&dto.ProductSearchRequest{
			Pagable: _dto.Pagable{
				Page: req.Page,
				Size: req.Size,
			},
		})
	if err != nil {
		return nil, err
	}

	products := h.ProductMapper.MapProductPbList(result)
	return &pb_bdspro.SearchResponse{
		Data:          products,
		TotalElements: total,
	}, nil
}

// @Summary Create a new product
// @Description Create a new product
// @Tags AdminProduct
// @Accept json
// @Produce json
// @Param product body pb_bdspro.ProductSaveRequest true "Product to create"
// @Success 200 {object} pb_bdspro.ProductResponse
// @Router /admin/products [post]
func (h *AdminProductHandler) CreateProduct(ctx context.Context, req *pb_bdspro.ProductSaveRequest) (*pb_bdspro.ProductResponse, error) {
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
	return h.ProductMapper.MapDetailProductPb(createdProduct), nil
}

// @Summary Update a product
// @Description Update a product
// @Tags AdminProduct
// @Accept json
// @Produce json
// @Param product body pb_bdspro.ProductSaveRequest true "Product to update"
// @Success 200 {object} pb_bdspro.ProductResponse
// @Router /admin/products [put]
func (h *AdminProductHandler) UpdateProduct(ctx context.Context, req *pb_bdspro.ProductSaveRequest) (*pb_bdspro.ProductResponse, error) {
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
// @Success 200 {object} pb_bdspro.Response
// @Router /admin/products/{id} [delete]
func (h *AdminProductHandler) DeleteProduct(ctx context.Context, req *sharepb.IdRequest) (*pb_bdspro.Response, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_SP_XOA"}); err != nil {
		return nil, err
	}
	// return h.AdminProductService.DeleteProduct(ctx, req)
	return nil, nil
}

func (h *AdminProductHandler) GetProductById(ctx context.Context, req *sharepb.IdRequest) (*pb_bdspro.ProductResponse, error) {
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_SP_XEM"}); err != nil {
		return nil, err
	}
	// return h.AdminProductService.GetProductById(ctx, req)
	product, err := h.AdminProductService.GetDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return h.ProductMapper.MapDetailProductPb(product), nil
}
