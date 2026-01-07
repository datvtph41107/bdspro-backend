package handler

import (
	"bdspro/infra/mapper"
	"bdspro/internal/usecases"
	"context"

	bdspropb "pb/types/bdspro"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PropertyHandler struct {
	bdspropb.UnimplementedPropertyServiceServer
	PropertyUsecase *usecases.PropertyUsecase
	PropertyMapper  *mapper.PropertyMapper
}

func NewPropertyHandler(
	propertyUsecase *usecases.PropertyUsecase,
	propertyMapper *mapper.PropertyMapper,
) *PropertyHandler {
	return &PropertyHandler{
		PropertyUsecase: propertyUsecase,
		PropertyMapper:  propertyMapper,
	}
}

// @Summary Lấy danh sách Property
// @Description Lấy danh sách Property với phân trang và filter
// @Tags Property
// @Accept json
// @Produce json
// @Param page query int false "Trang hiện tại" default(1)
// @Param size query int false "Số lượng item" default(20)
// @Param text query string false "Tìm kiếm theo text"
// @Param propertyTypeId query int false "ID loại BĐS"
// @Param locationId query int false "ID địa điểm"
// @Param projectId query int false "ID dự án"
// @Param recordStatus query string false "Trạng thái bản ghi"
// @Success 200 {object} bdspropb.PropertyListResponse
// @Router /v2/bdspro/v2/property/list [get]
func (h *PropertyHandler) GetList(ctx context.Context, req *bdspropb.PropertySearchRequest) (*bdspropb.PropertyListResponse, error) {
	// Convert proto request to DTO
	dtoReq := h.PropertyMapper.PbToPropertySearchDTO(req)

	// Call usecase
	properties, total, err := h.PropertyUsecase.GetList(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get property list: %v", err)
	}

	// Convert to proto response
	response := &bdspropb.PropertyListResponse{
		Data:  h.PropertyMapper.PropertyListToProto(properties),
		Total: total,
	}

	return response, nil
}

// @Summary Lấy chi tiết Property
// @Description API để lấy chi tiết Property theo ID
// @Tags Property
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Property"
// @Success 200 {object} bdspropb.PropertyResponse "Chi tiết Property"
// @Router /v2/bdspro/v2/property/detail/{id} [get]
func (h *PropertyHandler) GetDetail(ctx context.Context, req *bdspropb.PropertyDetailRequest) (*bdspropb.PropertyResponse, error) {
	property, err := h.PropertyUsecase.Detail(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get property detail: %v", err)
	}
	return h.PropertyMapper.PropertyToProto(property), nil
}

// @Summary Lấy danh sách Asset liên quan đến Property
// @Description API để lấy danh sách Asset liên quan đến Property theo ID
// @Tags Property
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Property"
// @Success 200 {object} bdspropb.PropertyAssetsResponse "Danh sách Asset liên quan đến Property"
// @Router /v2/bdspro/v2/property/{id}/assets [get]
func (h *PropertyHandler) GetAssets(ctx context.Context, req *bdspropb.PropertyAssetsRequest) (*bdspropb.PropertyAssetsResponse, error) {
	return nil, nil
}

func (h *PropertyHandler) GetProducts(ctx context.Context, req *bdspropb.PropertyProductsRequest) (*bdspropb.PropertyProductsResponse, error) {
	return nil, nil
}

// @Summary Tạo Property
// @Description API để tạo một Property mới
// @Tags Property
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body bdspropb.CreatePropertyRequest true "Thông tin Property"
// @Success 201 {object} bdspropb.Response "Property được tạo thành công"
// @Router /v2/bdspro/property/new [post]
func (h *PropertyHandler) CreateProperty(ctx context.Context, req *bdspropb.CreatePropertyRequest) (*bdspropb.Response, error) {
	body := h.PropertyMapper.CreatePropertyPbToDTO(req)

	property, err := h.PropertyUsecase.Create(ctx, body)
	if err != nil {
		return nil, err
	}
	return &bdspropb.Response{
		Id:      property.ID,
		Message: "successfully",
	}, nil
}
