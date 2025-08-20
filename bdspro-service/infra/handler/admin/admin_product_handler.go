package admin_handler

import (
	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	admin_usecases "bdspro/internal/usecases/admin"
	_dto "common/dto"
	"context"
	pb_bdspro "pb/types/bdspro"
	sharepb "pb/types/shared"
)

type AdminProductHandler struct {
	pb_bdspro.UnimplementedAdminProductServiceServer
	AdminProductService *admin_usecases.AdminProductUsecase
	ProductMapper       *mapper.ProductMapper
}

func NewAdminProductHandler(
	uc *admin_usecases.AdminProductUsecase,
	productMapper *mapper.ProductMapper,
) *AdminProductHandler {
	return &AdminProductHandler{
		AdminProductService: uc,
		ProductMapper:       productMapper,
	}
}

func (h *AdminProductHandler) GetProducts(ctx context.Context, req *sharepb.IdRequest) (*pb_bdspro.SearchResponse, error) {
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

func (h *AdminProductHandler) CreateProduct(ctx context.Context, req *pb_bdspro.ProductSaveRequest) (*pb_bdspro.ProductResponse, error) {
	// product := h.ProductMapper.ProductSavePbToDTO(req)
	// if req.Name == "" {
	// 	return nil, _errors.ReturnError(400, "Tên sản phẩm không được để trống")
	// }
	// product, err := h.AdminProductService.Create(ctx, product)
	// if err != nil {
	// 	return nil, err
	// }
	// return h.ProductMapper.MapProductPb(product), nil
	return nil, nil
}

func (h *AdminProductHandler) UpdateProduct(ctx context.Context, req *pb_bdspro.ProductSaveRequest) (*pb_bdspro.ProductResponse, error) {
	// return h.AdminProductService.UpdateProduct(ctx, req)
	return nil, nil
}

func (h *AdminProductHandler) DeleteProduct(ctx context.Context, req *sharepb.IdRequest) (*pb_bdspro.Response, error) {
	// return h.AdminProductService.DeleteProduct(ctx, req)
	return nil, nil
}
