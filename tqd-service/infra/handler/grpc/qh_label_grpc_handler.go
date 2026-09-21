package handler_grpc

import (
	"context"
	"strings"
	"time"

	_dto "common/domain/dto"
	_errors "common/errors"
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
	labelUsecase  usecase.QHLabelUsecase
	legendUsecase usecase.QHLayerLegendUsecase
	SyncProvider  *_utils.SyncUtil
}

func NewQHLabelGrpcHandler(labelUsecase usecase.QHLabelUsecase, legendUsecase usecase.QHLayerLegendUsecase, syncProvider *_utils.SyncUtil) *QHLabelGrpcHandler {
	return &QHLabelGrpcHandler{labelUsecase: labelUsecase, legendUsecase: legendUsecase, SyncProvider: syncProvider}
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
		LayerID:         req.LayerId,
		Name:            req.Name,
		DisplayName:     req.DisplayName,
		Description:     req.Description,
		Color:           req.Color,
		FillOpacity:     req.FillOpacity,
		StrokeColor:     req.StrokeColor,
		StrokeWidth:     int(req.StrokeWidth),
		StrokeDashArray: req.StrokeDashArray,
		DisplayOrder:    int(req.DisplayOrder),
		IsVisible:       req.IsVisible,
		MinZoom:         int(req.MinZoom),
		MaxZoom:         int(req.MaxZoom),
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

	if req.StandardAt != nil {
		s := strings.TrimSpace(*req.StandardAt)
		if s != "" {
			t := _utils.ParseStringToTime(s)
			if t == nil {
				return nil, status.Error(codes.InvalidArgument, "standardAt must be RFC3339")
			}
			label.StandardAt = t
		}
	}

	result, err := h.labelUsecase.Create(ctx, label)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLabelList, 0), t.UnixMilli())
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDLabelDetail, result.ID), t.UnixMilli())

	return toQHLabelResponse(result), nil
}

// BatchCreateLabels - Tạo nhiều label trong một layer; bỏ qua name đã tồn tại hoặc trùng trong cùng request
func (h *QHLabelGrpcHandler) BatchCreateLabels(ctx context.Context, req *tqdpb.BatchCreateLabelsRequest) (*tqdpb.BatchCreateLabelsResponse, error) {
	if req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layerId is required")
	}
	if len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "items is required")
	}

	items := make([]*qh_domain.QHLabel, 0, len(req.Items))
	for _, it := range req.Items {
		if it == nil {
			continue
		}
		l := &qh_domain.QHLabel{
			Name:            it.Name,
			DisplayName:     it.DisplayName,
			Description:     it.Description,
			Color:           it.Color,
			FillOpacity:     it.FillOpacity,
			StrokeColor:     it.StrokeColor,
			StrokeWidth:     int(it.StrokeWidth),
			StrokeDashArray: it.StrokeDashArray,
			DisplayOrder:    int(it.DisplayOrder),
			IsVisible:       it.IsVisible,
			MinZoom:         int(it.MinZoom),
			MaxZoom:         int(it.MaxZoom),
		}
		if it.Status != nil {
			l.Status = enums.LabelStatus(*it.Status)
		}
		if it.StandardAt != nil {
			s := strings.TrimSpace(*it.StandardAt)
			if s != "" {
				t := _utils.ParseStringToTime(s)
				if t == nil {
					return nil, status.Error(codes.InvalidArgument, "standardAt must be RFC3339")
				}
				l.StandardAt = t
			}
		}
		items = append(items, l)
	}

	created, skipped, err := h.labelUsecase.BatchCreate(ctx, req.LayerId, items)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLabelList, 0), t.UnixMilli())

	out := make([]*tqdpb.QHLabelResponse, len(created))
	for i, c := range created {
		out[i] = toQHLabelResponse(c)
	}
	return &tqdpb.BatchCreateLabelsResponse{
		Created:      out,
		SkippedNames: skipped,
	}, nil
}

// GetLabel - Lấy chi tiết label theo ID
func (h *QHLabelGrpcHandler) GetLabel(ctx context.Context, req *tqdpb.GetLabelRequest) (*tqdpb.QHLabelResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	key := h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDLabelDetail, req.Id)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.QHLabelResponse{}, nil
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
	if req.StrokeDashArray != nil {
		label.StrokeDashArray = *req.StrokeDashArray
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
	if req.StandardAt != nil {
		s := strings.TrimSpace(*req.StandardAt)
		if s == "" {
			label.ClearStandardAt = true
		} else {
			t := _utils.ParseStringToTime(s)
			if t == nil {
				return nil, status.Error(codes.InvalidArgument, "standardAt must be RFC3339")
			}
			label.StandardAt = t
		}
	}

	result, err := h.labelUsecase.Update(ctx, req.Id, label)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if req.LayerId != nil && *req.LayerId > 0 {
		if err := h.labelUsecase.UpdateLayerLinkMetadata(ctx, req.Id, *req.LayerId, req.LandUseId, req.LegendId); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLabelList, 0), t.UnixMilli())
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDLabelDetail, req.Id), t.UnixMilli())

	resp := toQHLabelResponse(result)
	if req.LayerId != nil && *req.LayerId > 0 {
		resp.LayerId = *req.LayerId
		resp.LandUseId = req.LandUseId
		resp.LegendId = req.LegendId
	}
	return resp, nil
}

// DeleteLabel - Xóa mềm label
func (h *QHLabelGrpcHandler) DeleteLabel(ctx context.Context, req *tqdpb.DeleteLabelRequest) (*emptypb.Empty, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := h.labelUsecase.Delete(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLabelList, 0), t.UnixMilli())
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDLabelDetail, req.Id), t.UnixMilli())

	return &emptypb.Empty{}, nil
}

// BatchDeleteLabels - Xóa mềm nhiều label theo danh sách id (HTTP POST /v2/tqd/qh/admin/labels/batch-delete)
func (h *QHLabelGrpcHandler) BatchDeleteLabels(ctx context.Context, req *tqdpb.BatchDeleteLabelsRequest) (*emptypb.Empty, error) {
	if len(req.Ids) == 0 {
		return nil, status.Error(codes.InvalidArgument, "ids is required")
	}
	hasNonZero := false
	for _, id := range req.Ids {
		if id != 0 {
			hasNonZero = true
			break
		}
	}
	if !hasNonZero {
		return nil, status.Error(codes.InvalidArgument, "at least one valid id is required")
	}

	if err := h.labelUsecase.BatchDelete(ctx, req.Ids); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLabelList, 0), t.UnixMilli())

	return &emptypb.Empty{}, nil
}

// MergeLabels - Gộp nhiều nhãn thành một nhãn mới (HTTP POST /v2/tqd/qh/admin/labels/merge)
func (h *QHLabelGrpcHandler) MergeLabels(ctx context.Context, req *tqdpb.MergeLabelsRequest) (*tqdpb.QHLabelResponse, error) {
	if req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layerId is required")
	}
	if len(req.SourceLabelIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "sourceLabelIds is required")
	}
	nl := req.NewLabel
	if nl == nil {
		return nil, status.Error(codes.InvalidArgument, "newLabel is required")
	}
	if strings.TrimSpace(nl.Name) == "" {
		return nil, status.Error(codes.InvalidArgument, "newLabel.name is required")
	}

	label := &qh_domain.QHLabel{
		Name:            nl.Name,
		DisplayName:     nl.DisplayName,
		Description:     nl.Description,
		Color:           nl.Color,
		FillOpacity:     nl.FillOpacity,
		StrokeColor:     nl.StrokeColor,
		StrokeWidth:     int(nl.StrokeWidth),
		StrokeDashArray: nl.StrokeDashArray,
		DisplayOrder:    int(nl.DisplayOrder),
		IsVisible:       nl.IsVisible,
		MinZoom:         int(nl.MinZoom),
		MaxZoom:         int(nl.MaxZoom),
	}
	if nl.Status != nil {
		label.Status = enums.LabelStatus(*nl.Status)
	}
	if nl.StandardAt != nil {
		s := strings.TrimSpace(*nl.StandardAt)
		if s != "" {
			t := _utils.ParseStringToTime(s)
			if t == nil {
				return nil, status.Error(codes.InvalidArgument, "standardAt must be RFC3339")
			}
			label.StandardAt = t
		}
	}

	result, err := h.labelUsecase.Merge(ctx, req.LayerId, req.SourceLabelIds, label)
	if err != nil {
		return nil, qhLabelError(err)
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLabelList, 0), t.UnixMilli())
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDLabelDetail, result.ID), t.UnixMilli())

	return toQHLabelResponse(result), nil
}

// ListAdminLabels - Danh sách tất cả label (admin), phân trang (GET /v2/tqd/qh/admin/labels/list)
func (h *QHLabelGrpcHandler) ListAdminLabels(ctx context.Context, req *tqdpb.ListAdminLabelsRequest) (*tqdpb.ListLabelsResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLabelList, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ListLabelsResponse{}, nil
	}

	page := req.GetPage()
	size := req.GetSize()
	pagable := _dto.NewPagableFromGrpc(&page, &size, nil)

	var layerID *uint64
	if req.LayerId != nil {
		lid := req.GetLayerId()
		if lid != 0 {
			v := lid
			layerID = &v
		}
	}

	labels, total, err := h.labelUsecase.ListAdmin(ctx, layerID, pagable, req.GetIncludeInactive())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	data := make([]*tqdpb.QHLabelResponse, len(labels))
	for i := range labels {
		data[i] = toQHLabelResponse(&labels[i])
	}

	return &tqdpb.ListLabelsResponse{
		Data:     data,
		Total:    total,
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}

// ListLabelsByLayer - Danh sách label theo layer (admin)
func (h *QHLabelGrpcHandler) ListLabelsByLayer(ctx context.Context, req *tqdpb.ListLabelsByLayerRequest) (*tqdpb.ListLabelsResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLabelList, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ListLabelsResponse{}, nil
	}

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

// buildLegendMapByLayer trả về map labelID -> legend cho một layer.
func (h *QHLabelGrpcHandler) buildLegendMapByLayer(ctx context.Context, layerID uint64) (map[uint64]*qh_domain.QHLayerLegend, error) {
	legends, err := h.legendUsecase.ListAllByLayer(ctx, layerID)
	if err != nil {
		return nil, err
	}
	m := make(map[uint64]*qh_domain.QHLayerLegend, len(legends))
	for i := range legends {
		m[legends[i].LabelID] = &legends[i]
	}
	return m, nil
}

// =====================================================
// CLIENT APIS
// =====================================================

// ListClientLabelsByLayer - Danh sách label public theo layer
func (h *QHLabelGrpcHandler) ListClientLabelsByLayer(ctx context.Context, req *tqdpb.ListLabelsByLayerRequest) (*tqdpb.ListLabelsResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLabelList, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ListLabelsResponse{}, nil
	}

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

	key := h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDLabelDetail, req.Id)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.QHLabelResponse{}, nil
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

func qhLabelError(err error) error {
	return _errors.ToGRPC(err)
}

func toQHLabelResponse(l *qh_domain.QHLabel) *tqdpb.QHLabelResponse {
	if l == nil {
		return nil
	}

	resp := &tqdpb.QHLabelResponse{
		Id:              l.ID,
		LayerId:         l.LayerID,
		Name:            l.Name,
		DisplayName:     l.DisplayName,
		Description:     l.Description,
		Color:           l.Color,
		FillOpacity:     l.FillOpacity,
		StrokeColor:     l.StrokeColor,
		StrokeWidth:     int32(l.StrokeWidth),
		StrokeDashArray: l.StrokeDashArray,
		DisplayOrder:    int32(l.DisplayOrder),
		IsVisible:       l.IsVisible,
		RegionCount:     l.RegionCount,
		MinZoom:         int32(l.MinZoom),
		MaxZoom:         int32(l.MaxZoom),
		Status:          uint32(l.Status),
		CreatedAt:       _utils.FormatTimeToString(&l.CreatedAt),
		UpdatedAt:       _utils.FormatTimeToString(&l.UpdatedAt),
	}
	resp.LandUseId = l.LandUseID
	resp.LegendId = l.LegendID
	if l.StandardAt != nil {
		resp.StandardAt = _utils.FormatTimeToString(l.StandardAt)
	}
	return resp
}
