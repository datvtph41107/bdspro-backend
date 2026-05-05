package handler_grpc

import (
	"context"

	_dto "common/domain/dto"
	_utils "common/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	tqdpb "pb/types/tqd"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
	"tqd/internal/usecase"
)

// QHLabelGrpcHandler xử lý gRPC cho QHLabelService
type QHLabelGrpcHandler struct {
	tqdpb.UnimplementedQHLabelServiceServer
	labelUsecase usecase.QHLabelUsecase
}

func NewQHLabelGrpcHandler(labelUsecase usecase.QHLabelUsecase) *QHLabelGrpcHandler {
	return &QHLabelGrpcHandler{labelUsecase: labelUsecase}
}

// =====================================================
// ADMIN APIS
// =====================================================

// CreateLabel - Tạo label mới
func (h *QHLabelGrpcHandler) CreateLabel(ctx context.Context, req *tqdpb.CreateLabelRequest) (*tqdpb.QHLabelResponse, error) {
	if req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layerId is required")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	label := &qh_domain.QHLabel{
		LayerID:      req.LayerId,
		Name:         req.Name,
		DisplayName:  req.DisplayName,
		Description:  req.Description,
		Color:        req.Color,
		FillOpacity:  req.FillOpacity,
		StrokeColor:  req.StrokeColor,
		StrokeWidth:  int(req.StrokeWidth),
		DisplayOrder: int(req.DisplayOrder),
		IsVisible:    req.IsVisible,
		MinZoom:      int(req.MinZoom),
		MaxZoom:      int(req.MaxZoom),
	}

	// Set status from request or default to active
	if req.Status != nil {
		label.Status = enums.LabelStatus(*req.Status)
	} else {
		label.Status = 10 // active
	}

	// Set defaults if empty
	if label.Color == "" {
		label.Color = "#CCCCCC"
	}
	if label.FillOpacity == 0 {
		label.FillOpacity = 0.6
	}
	if label.StrokeColor == "" {
		label.StrokeColor = "#000000"
	}
	if label.StrokeWidth == 0 {
		label.StrokeWidth = 1
	}

	result, err := h.labelUsecase.Create(ctx, label)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toQHLabelResponse(result), nil
}

// GetLabel - Lấy chi tiết label theo ID
func (h *QHLabelGrpcHandler) GetLabel(ctx context.Context, req *tqdpb.GetLabelRequest) (*tqdpb.QHLabelResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	result, err := h.labelUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if result == nil {
		return nil, status.Error(codes.NotFound, "label not found")
	}

	return toQHLabelResponse(result), nil
}

// UpdateLabel - Cập nhật label
func (h *QHLabelGrpcHandler) UpdateLabel(ctx context.Context, req *tqdpb.UpdateLabelRequest) (*tqdpb.QHLabelResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	label := &qh_domain.QHLabel{ID: req.Id}

	if req.Name != nil {
		label.Name = *req.Name
	}
	if req.DisplayName != nil {
		label.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		label.Description = *req.Description
	}
	if req.Color != nil {
		label.Color = *req.Color
	}
	if req.FillOpacity != nil {
		label.FillOpacity = *req.FillOpacity
	}
	if req.StrokeColor != nil {
		label.StrokeColor = *req.StrokeColor
	}
	if req.StrokeWidth != nil {
		label.StrokeWidth = int(*req.StrokeWidth)
	}
	if req.DisplayOrder != nil {
		label.DisplayOrder = int(*req.DisplayOrder)
	}
	if req.IsVisible != nil {
		label.IsVisible = *req.IsVisible
	}
	if req.MinZoom != nil {
		label.MinZoom = int(*req.MinZoom)
	}
	if req.MaxZoom != nil {
		label.MaxZoom = int(*req.MaxZoom)
	}
	if req.Status != nil {
		label.Status = enums.LabelStatus(*req.Status)
	}

	result, err := h.labelUsecase.Update(ctx, req.Id, label)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toQHLabelResponse(result), nil
}

// DeleteLabel - Xóa mềm label
func (h *QHLabelGrpcHandler) DeleteLabel(ctx context.Context, req *tqdpb.DeleteLabelRequest) (*emptypb.Empty, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := h.labelUsecase.Delete(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

// ListLabelsByLayer - Danh sách label theo layer (admin)
func (h *QHLabelGrpcHandler) ListLabelsByLayer(ctx context.Context, req *tqdpb.ListLabelsByLayerRequest) (*tqdpb.ListLabelsResponse, error) {
	if req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layerId is required")
	}

	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)

	labels, total, err := h.labelUsecase.ListByLayerID(ctx, req.LayerId, pagable, req.IncludeInactive)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	data := make([]*tqdpb.QHLabelResponse, len(labels))
	for i, l := range labels {
		data[i] = toQHLabelResponse(&l)
	}

	return &tqdpb.ListLabelsResponse{
		Data:     data,
		Total:    total,
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}

// =====================================================
// CLIENT APIS
// =====================================================

// ListClientLabelsByLayer - Danh sách label public theo layer
func (h *QHLabelGrpcHandler) ListClientLabelsByLayer(ctx context.Context, req *tqdpb.ListLabelsByLayerRequest) (*tqdpb.ListLabelsResponse, error) {
	if req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layerId is required")
	}

	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)

	labels, total, err := h.labelUsecase.ListClientByLayerID(ctx, req.LayerId, pagable)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	data := make([]*tqdpb.QHLabelResponse, len(labels))
	for i, l := range labels {
		data[i] = toQHLabelResponse(&l)
	}

	return &tqdpb.ListLabelsResponse{
		Data:     data,
		Total:    total,
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}

// GetClientLabel - Chi tiết label public
func (h *QHLabelGrpcHandler) GetClientLabel(ctx context.Context, req *tqdpb.GetLabelRequest) (*tqdpb.QHLabelResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	label, err := h.labelUsecase.GetClientByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if label == nil {
		return nil, status.Error(codes.NotFound, "label not found")
	}

	return toQHLabelResponse(label), nil
}

// =====================================================
// HELPER FUNCTIONS
// =====================================================

func toQHLabelResponse(l *qh_domain.QHLabel) *tqdpb.QHLabelResponse {
	if l == nil {
		return nil
	}

	return &tqdpb.QHLabelResponse{
		Id:           l.ID,
		LayerId:      l.LayerID,
		Name:         l.Name,
		DisplayName:  l.DisplayName,
		Description:  l.Description,
		Color:        l.Color,
		FillOpacity:  l.FillOpacity,
		StrokeColor:  l.StrokeColor,
		StrokeWidth:  int32(l.StrokeWidth),
		DisplayOrder: int32(l.DisplayOrder),
		IsVisible:    l.IsVisible,
		MinZoom:      int32(l.MinZoom),
		MaxZoom:      int32(l.MaxZoom),
		Status:       uint32(l.Status),
		CreatedAt:    _utils.FormatTimeToString(&l.CreatedAt),
		UpdatedAt:    _utils.FormatTimeToString(&l.UpdatedAt),
	}
}
