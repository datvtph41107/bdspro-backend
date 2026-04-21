package handler_grpc

import (
	_errors "common/errors"
	"context"
	"fmt"
	"time"

	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal/domain"
	"tqd/internal/enums"
	"tqd/internal/usecase"
)

type LayerGrpcHandler struct {
	tqdpb.UnimplementedLayerServiceServer
	layerUsecase usecase.LayerUsecase
}

func NewLayerGrpcHandler(layerUsecase usecase.LayerUsecase) *LayerGrpcHandler {
	return &LayerGrpcHandler{
		layerUsecase: layerUsecase,
	}
}

// CreateLayer thêm mới layer
func (h *LayerGrpcHandler) CreateLayer(ctx context.Context, req *tqdpb.CreateLayerRequest) (*tqdpb.LayerResponse, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, _errors.ReturnError(400, "name is required")
	}
	if req.DisplayName == "" {
		return nil, _errors.ReturnError(400, "displayName is required")
	}
	if req.Type == 0 {
		return nil, _errors.ReturnError(400, "type is required, must be 1,2,3,4,5,6")
	}

	// Validate type
	layerType := enums.LayerType(req.Type)
	if !layerType.IsValid() {
		return nil, _errors.ReturnError(400, fmt.Sprintf("invalid type: %d, must be 1,2,3,4,5,6", req.Type))
	}

	// Parse status from request (uint32)
	status := enums.LayerStatusActive
	if req.Status != nil && *req.Status > 0 {
		status = enums.LayerStatus(*req.Status)
		if !status.IsValid() {
			return nil, _errors.ReturnError(400, fmt.Sprintf("invalid status: %d, must be 10,20,30", *req.Status))
		}
	}

	// Tạo domain object
	layerDomain := &domain.Layer{
		Name:         req.Name,
		DisplayName:  req.DisplayName,
		Type:         layerType,
		Status:       status,
		DisplayOrder: int(req.Order),
		MinZoom:      int(req.MinZoom),
		MaxZoom:      int(req.MaxZoom),
		Visible:      req.Visible,
		LegalDoc:     req.LegalDoc,
	}

	// Handle effectiveDate
	if req.EffectiveDate != nil && *req.EffectiveDate != "" {
		parsedTime, err := time.Parse("2006-01-02", *req.EffectiveDate)
		if err != nil {
			return nil, _errors.ReturnError(400, "invalid effectiveDate format, expected YYYY-MM-DD")
		}
		layerDomain.EffectiveDate = &parsedTime
	}

	// Gọi usecase
	createdDomain, err := h.layerUsecase.Create(ctx, layerDomain)
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}

	// Convert response
	protoLayer, err := mapper.ToProtoLayer(createdDomain)
	if err != nil {
		return nil, _errors.ReturnError(500, "internal error: "+err.Error())
	}

	return protoLayer, nil
}

// UpdateLayer cập nhật layer
func (h *LayerGrpcHandler) UpdateLayer(ctx context.Context, req *tqdpb.UpdateLayerRequest) (*tqdpb.LayerResponse, error) {
	if req.Id == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}

	// Convert to domain
	layerDomain := &domain.Layer{
		ID: uint(req.Id),
	}

	if req.Name != nil {
		layerDomain.Name = *req.Name
	}
	if req.DisplayName != nil {
		layerDomain.DisplayName = *req.DisplayName
	}
	if req.Type != nil && *req.Type > 0 {
		layerType := enums.LayerType(*req.Type)
		if !layerType.IsValid() {
			return nil, _errors.ReturnError(400, fmt.Sprintf("invalid type: %d", *req.Type))
		}
		layerDomain.Type = layerType
	}
	if req.Status != nil && *req.Status > 0 {
		status := enums.LayerStatus(*req.Status)
		if !status.IsValid() {
			return nil, _errors.ReturnError(400, fmt.Sprintf("invalid status: %d, must be 10,20,30", *req.Status))
		}
		layerDomain.Status = status
	}
	if req.Order != nil {
		layerDomain.DisplayOrder = int(*req.Order)
	}
	if req.MinZoom != nil {
		layerDomain.MinZoom = int(*req.MinZoom)
	}
	if req.MaxZoom != nil {
		layerDomain.MaxZoom = int(*req.MaxZoom)
	}
	if req.Visible != nil {
		layerDomain.Visible = *req.Visible
	}
	if req.LegalDoc != nil {
		layerDomain.LegalDoc = *req.LegalDoc
	}
	if req.EffectiveDate != nil && *req.EffectiveDate != "" {
		parsedTime, err := time.Parse("2006-01-02", *req.EffectiveDate)
		if err != nil {
			return nil, _errors.ReturnError(400, "invalid effectiveDate format, expected YYYY-MM-DD")
		}
		layerDomain.EffectiveDate = &parsedTime
	}

	// Call usecase
	updatedDomain, err := h.layerUsecase.Update(ctx, uint(req.Id), layerDomain)
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}

	if updatedDomain == nil {
		return nil, _errors.ReturnError(404, "layer not found")
	}

	// Convert response
	protoLayer, err := mapper.ToProtoLayer(updatedDomain)
	if err != nil {
		return nil, _errors.ReturnError(500, "internal error: "+err.Error())
	}

	return protoLayer, nil
}

// DeleteLayer xóa layer
func (h *LayerGrpcHandler) DeleteLayer(ctx context.Context, req *tqdpb.DeleteLayerRequest) (*sharepb.Empty, error) {
	if req.Id == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}

	err := h.layerUsecase.Delete(ctx, uint(req.Id))
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}

	return &sharepb.Empty{}, nil
}

// GetLayer lấy chi tiết layer
func (h *LayerGrpcHandler) GetLayer(ctx context.Context, req *tqdpb.GetLayerRequest) (*tqdpb.LayerResponse, error) {
	if req.Id == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}

	layerDomain, err := h.layerUsecase.GetByID(ctx, uint(req.Id))
	if err != nil {
		return nil, _errors.ReturnError(404, err.Error())
	}

	if layerDomain == nil {
		return nil, _errors.ReturnError(404, "layer not found")
	}

	protoLayer, err := mapper.ToProtoLayer(layerDomain)
	if err != nil {
		return nil, _errors.ReturnError(500, "internal error: "+err.Error())
	}

	return protoLayer, nil
}

// ListLayers danh sách layers
func (h *LayerGrpcHandler) ListLayers(ctx context.Context, req *tqdpb.ListLayersRequest) (*tqdpb.ListLayersResponse, error) {
	// Set default values
	page := int32(1)
	limit := int32(20)

	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}

	// Get status filter
	var status *uint32
	if req.Status != nil && *req.Status > 0 {
		status = req.Status
	}

	// Get type filter
	var layerType *uint32
	if req.Type != nil && *req.Type > 0 {
		layerType = req.Type
	}

	search := ""
	if req.Search != nil {
		search = *req.Search
	}

	// Call usecase
	layers, total, err := h.layerUsecase.List(ctx, status, layerType, search, int(page), int(limit))
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}

	// Convert response
	protoLayers := make([]*tqdpb.LayerResponse, 0, len(layers))
	for _, layer := range layers {
		protoLayer, err := mapper.ToProtoLayer(&layer)
		if err != nil {
			continue
		}
		protoLayers = append(protoLayers, protoLayer)
	}

	return &tqdpb.ListLayersResponse{
		Data:  protoLayers,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}
