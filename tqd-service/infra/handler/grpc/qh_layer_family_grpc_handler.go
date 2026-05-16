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

// QHLayerFamilyGrpcHandler xử lý gRPC QHLayerFamilyService.
type QHLayerFamilyGrpcHandler struct {
	tqdpb.UnimplementedQHLayerFamilyServiceServer
	uc usecase.QHLayerFamilyUsecase
}

func NewQHLayerFamilyGrpcHandler(uc usecase.QHLayerFamilyUsecase) *QHLayerFamilyGrpcHandler {
	return &QHLayerFamilyGrpcHandler{uc: uc}
}

// CreateLayerFamily — POST /v2/tqd/qh/admin/layer-families
func (h *QHLayerFamilyGrpcHandler) CreateLayerFamily(ctx context.Context, req *tqdpb.CreateLayerFamilyRequest) (*tqdpb.QHLayerFamilyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	row := &qh_domain.QHLayerFamily{Name: strings.TrimSpace(req.Name)}
	out, err := h.uc.Create(ctx, row)
	if err != nil {
		if strings.Contains(err.Error(), "required") {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toQHLayerFamilyPB(out), nil
}

// GetLayerFamily — GET /v2/tqd/qh/admin/layer-families/{id}
func (h *QHLayerFamilyGrpcHandler) GetLayerFamily(ctx context.Context, req *tqdpb.GetLayerFamilyRequest) (*tqdpb.QHLayerFamilyResponse, error) {
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
	return toQHLayerFamilyPB(out), nil
}

// UpdateLayerFamily — PUT /v2/tqd/qh/admin/layer-families/{id}
func (h *QHLayerFamilyGrpcHandler) UpdateLayerFamily(ctx context.Context, req *tqdpb.UpdateLayerFamilyRequest) (*tqdpb.QHLayerFamilyResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	in := &usecase.QHLayerFamilyUpdateInput{}
	if req.Name != nil {
		in.Name = req.Name
	}
	if in.Name == nil {
		return nil, status.Error(codes.InvalidArgument, "at least one field to update is required")
	}
	out, err := h.uc.Update(ctx, req.Id, in)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if strings.Contains(err.Error(), "empty") || strings.Contains(err.Error(), "required") {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toQHLayerFamilyPB(out), nil
}

// DeleteLayerFamily — DELETE /v2/tqd/qh/admin/layer-families/{id}
func (h *QHLayerFamilyGrpcHandler) DeleteLayerFamily(ctx context.Context, req *tqdpb.DeleteLayerFamilyRequest) (*emptypb.Empty, error) {
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

// ListLayerFamilies — GET /v2/tqd/qh/admin/layer-families/list
func (h *QHLayerFamilyGrpcHandler) ListLayerFamilies(ctx context.Context, req *tqdpb.ListLayerFamiliesRequest) (*tqdpb.ListLayerFamiliesResponse, error) {
	page := uint32(0)
	size := uint32(0)
	search := ""
	if req != nil {
		page = req.Page
		size = req.Size
		if req.Search != nil {
			search = *req.Search
		}
	}
	pagable := _dto.NewPagableFromGrpc(&page, &size, nil)
	rows, total, err := h.uc.List(ctx, pagable, search)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	data := make([]*tqdpb.QHLayerFamilyResponse, len(rows))
	for i := range rows {
		data[i] = toQHLayerFamilyPB(&rows[i])
	}
	return &tqdpb.ListLayerFamiliesResponse{
		Data:     data,
		Total:    total,
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}

func toQHLayerFamilyPB(e *qh_domain.QHLayerFamily) *tqdpb.QHLayerFamilyResponse {
	if e == nil {
		return nil
	}
	var createdAt, updatedAt string
	if e.CreatedAt != nil {
		createdAt = _utils.FormatTimeToString(e.CreatedAt)
	}
	if e.UpdatedAt != nil {
		updatedAt = _utils.FormatTimeToString(e.UpdatedAt)
	}
	return &tqdpb.QHLayerFamilyResponse{
		Id:        e.ID,
		Name:      e.Name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
