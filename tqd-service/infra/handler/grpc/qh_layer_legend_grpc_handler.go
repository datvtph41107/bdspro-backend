package handler_grpc

import (
	"context"
	"strings"

	_dto "common/domain/dto"
	_utils "common/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	tqdpb "pb/types/tqd"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/usecase"
)

// QHLayerLegendGrpcHandler xử lý gRPC QHLayerLegendService.
type QHLayerLegendGrpcHandler struct {
	tqdpb.UnimplementedQHLayerLegendServiceServer
	uc usecase.QHLayerLegendUsecase
}

func NewQHLayerLegendGrpcHandler(uc usecase.QHLayerLegendUsecase) *QHLayerLegendGrpcHandler {
	return &QHLayerLegendGrpcHandler{uc: uc}
}

// CreateLayerLegend — POST /v2/tqd/qh/admin/layer-legends
func (h *QHLayerLegendGrpcHandler) CreateLayerLegend(ctx context.Context, req *tqdpb.CreateLayerLegendRequest) (*tqdpb.QHLayerLegendResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	row := &qh_domain.QHLayerLegend{
		LayerID: req.LayerId,
		LabelID: req.LabelId,
		Note:    req.Note,
	}
	out, err := h.uc.Create(ctx, row)
	if err != nil {
		if strings.Contains(err.Error(), "already has legend") {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "required") {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toQHLayerLegendPB(out), nil
}

// GetLayerLegend — GET /v2/tqd/qh/admin/layer-legends/{id}
func (h *QHLayerLegendGrpcHandler) GetLayerLegend(ctx context.Context, req *tqdpb.GetLayerLegendRequest) (*tqdpb.QHLayerLegendResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	out, err := h.uc.GetByID(ctx, req.Id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if strings.Contains(err.Error(), "required") {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toQHLayerLegendPB(out), nil
}

// UpdateLayerLegend — PUT /v2/tqd/qh/admin/layer-legends/{id}
func (h *QHLayerLegendGrpcHandler) UpdateLayerLegend(ctx context.Context, req *tqdpb.UpdateLayerLegendRequest) (*tqdpb.QHLayerLegendResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	in := &usecase.QHLayerLegendUpdateInput{}
	if req.LayerId != nil {
		in.LayerID = req.LayerId
	}
	if req.LabelId != nil {
		in.LabelID = req.LabelId
	}
	if req.Note != nil {
		in.Note = req.Note
	}
	if in.LayerID == nil && in.LabelID == nil && in.Note == nil {
		return nil, status.Error(codes.InvalidArgument, "at least one field to update is required")
	}
	out, err := h.uc.Update(ctx, req.Id, in)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if strings.Contains(err.Error(), "already has legend") {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		if strings.Contains(err.Error(), "required") {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toQHLayerLegendPB(out), nil
}

// DeleteLayerLegend — DELETE /v2/tqd/qh/admin/layer-legends/{id}
func (h *QHLayerLegendGrpcHandler) DeleteLayerLegend(ctx context.Context, req *tqdpb.DeleteLayerLegendRequest) (*emptypb.Empty, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if err := h.uc.Delete(ctx, req.Id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

// ListLayerLegends — GET /v2/tqd/qh/admin/layer-legends/list
func (h *QHLayerLegendGrpcHandler) ListLayerLegends(ctx context.Context, req *tqdpb.ListLayerLegendsRequest) (*tqdpb.ListLayerLegendsResponse, error) {
	page := uint32(0)
	size := uint32(0)
	var layerID, labelID *uint64
	if req != nil {
		page = req.Page
		size = req.Size
		if req.LayerId != nil {
			layerID = req.LayerId
		}
		if req.LabelId != nil {
			labelID = req.LabelId
		}
	}
	pagable := _dto.NewPagableFromGrpc(&page, &size, nil)
	rows, total, err := h.uc.List(ctx, pagable, layerID, labelID)
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

func toQHLayerLegendPB(e *qh_domain.QHLayerLegend) *tqdpb.QHLayerLegendResponse {
	if e == nil {
		return nil
	}
	return &tqdpb.QHLayerLegendResponse{
		Id:        e.ID,
		LayerId:   e.LayerID,
		LabelId:   e.LabelID,
		Note:      e.Note,
		CreatedAt: _utils.FormatTimeToString(e.CreatedAt),
		UpdatedAt: _utils.FormatTimeToString(e.UpdatedAt),
	}
}
