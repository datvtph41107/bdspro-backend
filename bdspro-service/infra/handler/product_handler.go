package handler

import (
	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	shared_usecase "bdspro/internal/usecases/shared"
	"context"
	pb_bdspro "pb/types/bdspro"
)

type ProductHandler struct {
	pb_bdspro.UnimplementedProductServiceServer
	ProductUsecase *shared_usecase.ProductUsecase
	ProductMapper  *mapper.ProductMapper
}

func NewProductHandler(
	productUsecase *shared_usecase.ProductUsecase,
	productMapper *mapper.ProductMapper,
) *ProductHandler {
	return &ProductHandler{
		ProductUsecase: productUsecase,
		ProductMapper:  productMapper,
	}
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
// @Success 200 {object} pb_bdspro.SearchResponse
// @Router /v2/bdspro/v2/product/group/{groupId} [get]
func (h *ProductHandler) GetGroupProducts(ctx context.Context, req *pb_bdspro.GroupProductSearch) (*pb_bdspro.SearchResponse, error) {
	// Convert protobuf request to DTO
	searchReq := &dto.GroupProductSearchRequest{
		GroupID:          req.GroupId,
		Text:             req.Text,
		ParentID:         req.ParentId,
		PropertyTypeID:   &req.PropertyTypeId,
		DocTypeID:        &req.DocTypeId,
		ProjectID:        &req.ProjectId,
		ProvinceID:       &req.ProvinceId,
		DistrictID:       &req.DistrictId,
		WardID:           &req.WardId,
		TransactionType:  &req.TransactionType,
		SaleStatus:       &req.SaleStatus,
		RentStatus:       &req.RentStatus,
		SaleVisibility:   &req.SaleVisibility,
		RentVisibility:   &req.RentVisibility,
		SourceType:       &req.SourceType,
		TransactionPrice: &req.TransactionPrice,
		SalePrice:        &req.SalePrice,
		RentPrice:        &req.RentPrice,
		SaleCommission:   &req.SaleCommission,
		RentCommission:   &req.RentCommission,
		RentPaymentCycle: &req.RentPaymentCycle,
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
	products, total, err := h.ProductUsecase.GetGroupProducts(ctx, searchReq)
	if err != nil {
		return nil, err
	}

	// Convert domain products to protobuf response
	var productResponses []*pb_bdspro.ProductResponse
	for _, product := range products {
		productResponse := h.ProductMapper.MapProductPb(&product)
		productResponses = append(productResponses, productResponse)
	}

	return &pb_bdspro.SearchResponse{
		Data:          productResponses,
		TotalElements: total,
	}, nil
}
