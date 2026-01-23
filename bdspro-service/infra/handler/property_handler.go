package handler

import (
	"bdspro/infra/mapper"
	"bdspro/internal/dto"
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
	AssetMapper     *mapper.AssetMapper
}

func NewPropertyHandler(
	propertyUsecase *usecases.PropertyUsecase,
	propertyMapper *mapper.PropertyMapper,
	assetMapper *mapper.AssetMapper,
) *PropertyHandler {
	return &PropertyHandler{
		PropertyUsecase: propertyUsecase,
		PropertyMapper:  propertyMapper,
		AssetMapper:     assetMapper,
	}
}

func (h *PropertyHandler) SearchTags(
	ctx context.Context,
	req *bdspropb.TagSearchRequest,
) (*bdspropb.TagListResponse, error) {
	dtoReq := h.PropertyMapper.PbToTagSearchDTO(req)
	items, total, err := h.PropertyUsecase.SearchTags(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &bdspropb.TagListResponse{
		Data:  h.PropertyMapper.TagListToProto(items),
		Total: total,
	}, nil
}

func (h *PropertyHandler) GetPropertyAmenities(
	ctx context.Context,
	req *bdspropb.GetPropertyAmenitiesRequest,
) (*bdspropb.GetPropertyAmenitiesResponse, error) {

	items, err := h.PropertyUsecase.GetPropertyAmenities(ctx, req.PropertyId)
	if err != nil {
		return nil, err
	}

	result := make([]*bdspropb.PropertyAmenityItem, 0, len(items))
	for _, item := range items {
		result = append(result, &bdspropb.PropertyAmenityItem{
			Id:      item.ID,
			Name:    item.Name,
			Active:  item.Active,
			Checked: item.Checked,
		})
	}

	return &bdspropb.GetPropertyAmenitiesResponse{
		Data: result,
	}, nil
}

func (h *PropertyHandler) UpdatePropertyAmenities(
	ctx context.Context,
	req *bdspropb.UpdatePropertyAmenitiesRequest,
) (*bdspropb.Response, error) {
	err := h.PropertyUsecase.UpdatePropertyAmenities(ctx, &dto.UpdatePropertyAmenitiesRequest{
		PropertyID: req.PropertyId,
		AmenityIDs: req.AmenityIds,
	})
	if err != nil {
		return nil, err
	}

	return &bdspropb.Response{
		Id:      req.PropertyId,
		Message: "successfully",
	}, nil
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
		return nil, err
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

// @Summary Lấy danh sách Asset và Product liên quan đến Property
// @Description API để lấy danh sách Asset và Product liên quan đến Property theo ID
// @Tags Property
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Property"
// @Success 200 {object} bdspropb.PropertyRelationResponse "Danh sách Asset và Product liên quan đến Property"
// @Router /v2/bdspro/v2/property/{id}/relation [get]
func (h *PropertyHandler) GetRelation(ctx context.Context, req *bdspropb.PropertyRelationRequest) (*bdspropb.PropertyRelationResponse, error) {
	relations, err := h.PropertyUsecase.GetRelation(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get property relation: %v", err)
	}

	// Convert relations to proto
	relationList := h.PropertyMapper.PropertyRelationListToProto(relations)

	return &bdspropb.PropertyRelationResponse{
		Data:  relationList,
		Total: int64(len(relationList)),
	}, nil
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

// @Summary Cập nhật Property
// @Description API để cập nhật thông tin Property
// @Tags Property
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Property"
// @Param body body bdspropb.UpdatePropertyRequest true "Thông tin Property cần cập nhật"
// @Success 200 {object} bdspropb.Response "Property được cập nhật thành công"
// @Router /v2/bdspro/v2/property/{id} [put]
func (h *PropertyHandler) UpdateProperty(ctx context.Context, req *bdspropb.UpdatePropertyRequest) (*bdspropb.Response, error) {
	property, landInfo, buildingInfo := h.PropertyMapper.UpdatePropertyPbToDTO(req)

	err := h.PropertyUsecase.Update(ctx, req.Id, property, landInfo, buildingInfo, req.AmenityIds)
	if err != nil {
		return nil, err
	}

	return &bdspropb.Response{
		Id:      req.Id,
		Message: "successfully",
	}, nil
}

func (h *PropertyHandler) DeleteProperties(ctx context.Context, req *bdspropb.DeletePropertiesRequest) (*bdspropb.Response, error) {
	err := h.PropertyUsecase.DeleteProperties(ctx, req.Ids)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete properties: %v", err)
	}

	return &bdspropb.Response{
		Message: "Properties deleted successfully",
	}, nil
}

func (h *PropertyHandler) ArchiveProperties(ctx context.Context, req *bdspropb.ArchivePropertiesRequest) (*bdspropb.Response, error) {
	dtoReq := h.PropertyMapper.PbToArchivePropertiesDTO(req)
	err := h.PropertyUsecase.ArchiveProperties(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to archive properties: %v", err)
	}

	action := "archived"
	if !req.Archived {
		action = "unarchived"
	}

	return &bdspropb.Response{
		Message: "Properties " + action + " successfully",
	}, nil
}
