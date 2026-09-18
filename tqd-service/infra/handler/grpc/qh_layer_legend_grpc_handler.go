package handler_grpc

import (
	"context"
	"tqd/infra/handler/errorresponse"

	_dto "common/domain/dto"
	_utils "common/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/usecase"
)

type QHLayerLegendGrpcHandler struct {
	tqdpb.UnimplementedQHLayerLegendServiceServer
	uc usecase.QHLayerLegendUsecase
}

func NewQHLayerLegendGrpcHandler(uc usecase.QHLayerLegendUsecase) *QHLayerLegendGrpcHandler {
	return &QHLayerLegendGrpcHandler{uc: uc}
}

func (h *QHLayerLegendGrpcHandler) CreateLayerLegend(ctx context.Context, req *tqdpb.CreateLayerLegendRequest) (*tqdpb.QHLayerLegendResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	row := &qh_domain.QHLayerLegend{
		LayerID:         req.LayerId,
		LabelID:         req.LabelId,
		Color:           req.Color,
		Description:     req.Description,
		Note:            req.Note,
		DisplayOrder:    int(req.DisplayOrder),
		IsVisible:       req.IsVisible,
		LegendType:      req.LegendType,
		GeometryType:    req.GeometryType,
		StrokeDashArray: req.StrokeDashArray,
		IconURL:         req.IconUrl,
		ImageURL:        req.ImageUrl,
	}
	if req.LandUseGroupId != nil {
		row.LandUseID = req.LandUseGroupId
	}
	out, err := h.uc.Create(ctx, row)
	if err != nil {
		return nil, errorresponse.GRPC(ctx, err)
	}
	return toQHLayerLegendPB(out), nil
}

func (h *QHLayerLegendGrpcHandler) BatchCreateLayerLegends(ctx context.Context, req *tqdpb.BatchCreateLayerLegendsRequest) (*tqdpb.BatchCreateLayerLegendsResponse, error) {
	if req == nil || req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layerId is required")
	}
	items := make([]qh_domain.QHLayerLegend, len(req.Items))
	for i, it := range req.Items {
		items[i] = qh_domain.QHLayerLegend{
			LabelID:         it.LabelId,
			Note:            it.Note,
			DisplayOrder:    int(it.DisplayOrder),
			IsVisible:       it.IsVisible,
			LegendType:      it.LegendType,
			GeometryType:    it.GeometryType,
			StrokeDashArray: it.StrokeDashArray,
			IconURL:         it.IconUrl,
			ImageURL:        it.ImageUrl,
			Description:     it.Description,
			Color:           it.Color,
		}
		if it.LandUseGroupId != nil {
			items[i].LandUseID = it.LandUseGroupId
		}
	}
	created, skipped, err := h.uc.BatchCreate(ctx, req.LayerId, items)
	if err != nil {
		return nil, errorresponse.GRPC(ctx, err)
	}
	data := make([]*tqdpb.QHLayerLegendResponse, len(created))
	for i := range created {
		data[i] = toQHLayerLegendPB(&created[i])
	}
	return &tqdpb.BatchCreateLayerLegendsResponse{
		Created:         data,
		SkippedLabelIds: skipped,
	}, nil
}

func (h *QHLayerLegendGrpcHandler) GetLayerLegend(ctx context.Context, req *tqdpb.GetLayerLegendRequest) (*tqdpb.QHLayerLegendResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	out, err := h.uc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, errorresponse.GRPC(ctx, err)
	}
	return toQHLayerLegendPB(out), nil
}

func (h *QHLayerLegendGrpcHandler) UpdateLayerLegend(ctx context.Context, req *tqdpb.UpdateLayerLegendRequest) (*tqdpb.QHLayerLegendResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	in := &usecase.QHLayerLegendUpdateInput{}
	hasUpdate := false
	if req.LayerId != nil {
		in.LayerID = req.LayerId
		hasUpdate = true
	}
	if req.LabelId != nil {
		in.LabelID = req.LabelId
		hasUpdate = true
	}
	if req.Note != nil {
		in.Note = req.Note
		hasUpdate = true
	}
	if req.DisplayOrder != nil {
		v := int(*req.DisplayOrder)
		in.DisplayOrder = &v
		hasUpdate = true
	}
	if req.IsVisible != nil {
		in.IsVisible = req.IsVisible
		hasUpdate = true
	}
	if req.LegendType != nil {
		in.LegendType = req.LegendType
		hasUpdate = true
	}
	if req.GeometryType != nil {
		in.GeometryType = req.GeometryType
		hasUpdate = true
	}
	if req.StrokeDashArray != nil {
		in.StrokeDashArray = req.StrokeDashArray
		hasUpdate = true
	}
	if req.IconUrl != nil {
		in.IconURL = req.IconUrl
		hasUpdate = true
	}
	if req.ImageUrl != nil {
		in.ImageURL = req.ImageUrl
		hasUpdate = true
	}
	if req.LandUseGroupId != nil {
		in.LandUseID = req.LandUseGroupId
		hasUpdate = true
	}
	if req.Description != "" {
		in.Description = &req.Description
		hasUpdate = true
	}
	if req.Color != "" {
		in.Color = &req.Color
		hasUpdate = true
	}
	if !hasUpdate {
		return nil, status.Error(codes.InvalidArgument, "at least one field to update is required")
	}
	out, err := h.uc.Update(ctx, req.Id, in)
	if err != nil {
		return nil, errorresponse.GRPC(ctx, err)
	}
	return toQHLayerLegendPB(out), nil
}

func (h *QHLayerLegendGrpcHandler) DeleteLayerLegend(ctx context.Context, req *tqdpb.DeleteLayerLegendRequest) (*emptypb.Empty, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if err := h.uc.Delete(ctx, req.Id); err != nil {
		return nil, errorresponse.GRPC(ctx, err)
	}
	return &emptypb.Empty{}, nil
}

func (h *QHLayerLegendGrpcHandler) ListLayerLegends(ctx context.Context, req *tqdpb.ListLayerLegendsRequest) (*tqdpb.ListLayerLegendsResponse, error) {
	page := uint32(0)
	size := uint32(0)
	var layerID, labelID *uint64
	var legendType *string
	var landUseGroupID *uint64
	if req != nil {
		page = req.Page
		size = req.Size
		if req.LayerId != nil {
			layerID = req.LayerId
		}
		if req.LabelId != nil {
			labelID = req.LabelId
		}
		if req.LegendType != nil {
			legendType = req.LegendType
		}
		if req.LandUseGroupId != nil {
			landUseGroupID = req.LandUseGroupId
		}
	}
	pagable := _dto.NewPagableFromGrpc(&page, &size, nil)
	rows, total, err := h.uc.List(ctx, pagable, layerID, labelID, legendType, landUseGroupID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	data := make([]*tqdpb.QHLayerLegendResponse, len(rows))
	for i := range rows {
		data[i] = toQHLayerLegendPB(&rows[i])
	}
	return &tqdpb.ListLayerLegendsResponse{
		Data:     data,
		Total:    total,
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}

func (h *QHLayerLegendGrpcHandler) ListClientLayerLegends(ctx context.Context, req *tqdpb.ListClientLayerLegendsRequest) (*tqdpb.ListLayerLegendsResponse, error) {
	if req == nil || req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layerId is required")
	}
	page := req.Page
	size := req.Size
	pagable := _dto.NewPagableFromGrpc(&page, &size, nil)
	rows, total, err := h.uc.ListClientByLayer(ctx, req.LayerId, pagable)
	if err != nil {
		return nil, errorresponse.GRPC(ctx, err)
	}
	data := make([]*tqdpb.QHLayerLegendResponse, len(rows))
	for i := range rows {
		data[i] = toQHLayerLegendPB(&rows[i])
	}
	return &tqdpb.ListLayerLegendsResponse{
		Data:     data,
		Total:    total,
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}

// AdminCreateLayerLegend tạo chú thích layer.
// @Summary Tạo chú thích layer
// @Tags Layer Legend Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tqdpb.CreateLayerLegendRequest true "Chú thích layer"
// @Success 200 {object} tqdpb.QHLayerLegendResponse
// @Router /v2/tqd/admin/legends [post]
func (h *QHLayerLegendGrpcHandler) AdminCreateLayerLegend(ctx context.Context, req *tqdpb.CreateLayerLegendRequest) (*tqdpb.QHLayerLegendResponse, error) {
	if req == nil || req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layerId is required")
	}
	return h.CreateLayerLegend(ctx, req)
}

// AdminGetLayerLegend lấy chi tiết chú thích.
// @Summary Chi tiết chú thích layer
// @Tags Layer Legend Admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "Legend ID"
// @Success 200 {object} tqdpb.QHLayerLegendResponse
// @Router /v2/tqd/admin/legends/{id} [get]
func (h *QHLayerLegendGrpcHandler) AdminGetLayerLegend(ctx context.Context, req *tqdpb.GetLayerLegendRequest) (*tqdpb.QHLayerLegendResponse, error) {
	return h.GetLayerLegend(ctx, req)
}

// AdminUpdateLayerLegend cập nhật chú thích.
// @Summary Cập nhật chú thích layer
// @Tags Layer Legend Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Legend ID"
// @Param request body tqdpb.UpdateLayerLegendRequest true "Chú thích layer"
// @Success 200 {object} tqdpb.QHLayerLegendResponse
// @Router /v2/tqd/admin/legends/{id} [put]
func (h *QHLayerLegendGrpcHandler) AdminUpdateLayerLegend(ctx context.Context, req *tqdpb.UpdateLayerLegendRequest) (*tqdpb.QHLayerLegendResponse, error) {
	return h.UpdateLayerLegend(ctx, req)
}

// AdminDeleteLayerLegend xóa mềm chú thích.
// @Summary Xóa chú thích layer
// @Tags Layer Legend Admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "Legend ID"
// @Success 200 {object} sharepb.SubmitResponse
// @Router /v2/tqd/admin/legends/{id} [delete]
func (h *QHLayerLegendGrpcHandler) AdminDeleteLayerLegend(ctx context.Context, req *tqdpb.DeleteLayerLegendRequest) (*sharepb.SubmitResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if err := h.uc.Delete(ctx, req.Id); err != nil {
		return nil, errorresponse.GRPC(ctx, err)
	}
	return &sharepb.SubmitResponse{Message: "legend deleted successfully"}, nil
}

// AdminListLayerLegends danh sách chú thích (lọc theo layerId).
// @Summary Danh sách chú thích layer
// @Tags Layer Legend Admin
// @Produce json
// @Security BearerAuth
// @Param layerId query int false "Layer ID"
// @Param page query int false "Trang"
// @Param size query int false "Số bản ghi"
// @Success 200 {object} tqdpb.ListLayerLegendsResponse
// @Router /v2/tqd/admin/legends [get]
func (h *QHLayerLegendGrpcHandler) AdminListLayerLegends(ctx context.Context, req *tqdpb.ListLayerLegendsRequest) (*tqdpb.ListLayerLegendsResponse, error) {
	return h.ListLayerLegends(ctx, req)
}

// AdminListLayerLegendsByLayer danh sách chú thích theo layer.
// @Summary Danh sách chú thích theo layer
// @Tags Layer Legend Admin
// @Produce json
// @Security BearerAuth
// @Param layerId path int true "Layer ID"
// @Param page query int false "Trang"
// @Param size query int false "Số bản ghi"
// @Success 200 {object} tqdpb.ListLayerLegendsResponse
// @Router /v2/tqd/admin/layers/{layerId}/legends [get]
func (h *QHLayerLegendGrpcHandler) AdminListLayerLegendsByLayer(ctx context.Context, req *tqdpb.ListClientLayerLegendsRequest) (*tqdpb.ListLayerLegendsResponse, error) {
	if req == nil || req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layerId is required")
	}
	layerID := req.LayerId
	return h.ListLayerLegends(ctx, &tqdpb.ListLayerLegendsRequest{
		Page:    req.Page,
		Size:    req.Size,
		LayerId: &layerID,
	})
}

func (h *QHLayerLegendGrpcHandler) GetAllLayerLegends(ctx context.Context, req *tqdpb.GetAllLayerLegendsRequest) (*tqdpb.ListLayerLegendsResponse, error) {
	rows, total, err := h.uc.GetAll(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	data := make([]*tqdpb.QHLayerLegendResponse, len(rows))
	for i := range rows {
		data[i] = toQHLayerLegendPB(&rows[i])
	}
	return &tqdpb.ListLayerLegendsResponse{
		Data:     data,
		Total:    total,
		Page:     1,
		PageSize: 200,
	}, nil
}

func toQHLayerLegendPB(e *qh_domain.QHLayerLegend) *tqdpb.QHLayerLegendResponse {
	if e == nil {
		return nil
	}
	resp := &tqdpb.QHLayerLegendResponse{
		Id:              e.ID,
		LayerId:         e.LayerID,
		LabelId:         e.LabelID,
		Note:            e.Note,
		DisplayOrder:    int32(e.DisplayOrder),
		IsVisible:       e.IsVisible,
		LegendType:      e.GetLegendType(),
		GeometryType:    e.GetGeometryType(),
		StrokeDashArray: e.GetEffectiveStrokeDash(),
		IconUrl:         e.IconURL,
		ImageUrl:        e.ImageURL,
		Description:     e.Description,
		Color:           e.GetDisplayColor(),
		LandUseGroupId:  e.LandUseID,
		CreatedAt:       _utils.FormatTimeToString(e.CreatedAt),
		UpdatedAt:       _utils.FormatTimeToString(e.UpdatedAt),
	}
	if e.Label != nil {
		resp.LabelStyle = &tqdpb.LegendLabelStyleInfo{
			Id:              e.Label.ID,
			Name:            e.Label.Name,
			DisplayName:     e.Label.DisplayName,
			Color:           e.Label.Color,
			FillOpacity:     e.Label.FillOpacity,
			StrokeColor:     e.Label.StrokeColor,
			StrokeWidth:     int32(e.Label.StrokeWidth),
			StrokeDashArray: e.Label.StrokeDashArray,
		}
	}
	if e.LandUse != nil {
		resp.LandUseGroupInfo = &tqdpb.LandUseGroupInfo{
			Id:          e.LandUse.ID,
			Code:        e.LandUse.Code,
			Name:        e.LandUse.Name,
			Description: e.LandUse.Description,
			Color:       e.LandUse.Color,
			CanBuild:    e.LandUse.CanBuild,
			Priority:    int32(e.LandUse.Priority),
		}
	}
	return resp
}
