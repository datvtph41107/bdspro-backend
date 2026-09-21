package handler_grpc

import (
	"context"

	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	tqdpb "pb/types/tqd"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/usecase"
)

type QHLandUseGrpcHandler struct {
	tqdpb.UnimplementedQHLandUseServiceServer
	uc usecase.QHLayerLandUseGroupUsecase
}

func NewQHLandUseGrpcHandler(uc usecase.QHLayerLandUseGroupUsecase) *QHLandUseGrpcHandler {
	return &QHLandUseGrpcHandler{uc: uc}
}

func (h *QHLandUseGrpcHandler) CreateLandUse(ctx context.Context, req *tqdpb.CreateLandUseRequest) (*tqdpb.QHLandUseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	row := &qh_domain.QHLandUse{
		Note:         req.Note,
		DisplayOrder: int(req.DisplayOrder),
		IsVisible:    req.IsVisible,
		Color:        req.Color,
		Code:         req.Code,
	}
	row.Name = req.GetName()
	if req.GetLandUseId() > 0 {
		row.ID = req.GetLandUseId()
	}
	if req.GetLayerId() > 0 {
		row.Layers = []*qh_domain.QHLayer{{ID: req.GetLayerId()}}
	}
	out, err := h.uc.Create(ctx, row)
	if err != nil {
		return nil, mapLandUseGroupError(err)
	}
	return toQHLandUsePB(out), nil
}

func (h *QHLandUseGrpcHandler) CreateClientLandUse(ctx context.Context, req *tqdpb.CreateClientLandUseRequest) (*tqdpb.QHLandUseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	out, err := h.uc.CreateLayerLandUse(ctx, req.LayerId, req.LandUseId)
	if err != nil {
		return nil, mapLandUseGroupError(err)
	}
	return toQHLandUsePB(out), nil
}

// func (h *QHLandUseGrpcHandler) BatchCreateLandUses(ctx context.Context, req *tqdpb.BatchCreateLandUsesRequest) (*tqdpb.BatchCreateLandUsesResponse, error) {
// 	if req == nil || req.LayerId == 0 {
// 		return nil, status.Error(codes.InvalidArgument, "layerId is required")
// 	}
// 	items := make([]qh_domain.QHLandUse, len(req.Items))
// 	for i, it := range req.Items {
// 		items[i] = qh_domain.QHLandUse{
// 			Note:         it.Note,
// 			DisplayOrder: int(it.DisplayOrder),
// 			IsVisible:    it.IsVisible,
// 			Color:        it.Color,
// 			Code:         it.Code,
// 			Name:         it.Name,
// 		}
// 	}
// 	created, skipped, err := h.uc.BatchCreate(ctx, req.LayerId, items)
// 	if err != nil {
// 		return nil, mapLandUseGroupError(err)
// 	}
// 	data := make([]*tqdpb.QHLandUseResponse, len(created))
// 	for i := range created {
// 		data[i] = toQHLandUsePB(&created[i])
// 	}
// 	return &tqdpb.BatchCreateLandUsesResponse{
// 		Created:         data,
// 		SkippedGroupIds: skipped,
// 	}, nil
// }

func (h *QHLandUseGrpcHandler) GetLandUse(ctx context.Context, req *tqdpb.GetLandUseRequest) (*tqdpb.QHLandUseResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	out, err := h.uc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, mapLandUseGroupError(err)
	}
	return toQHLandUsePB(out), nil
}

func (h *QHLandUseGrpcHandler) UpdateLandUse(ctx context.Context, req *tqdpb.UpdateLandUseRequest) (*tqdpb.QHLandUseResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	in := &qh_domain.QHLandUse{
		// ID:    req.Id,
		Note:  req.Note,
		Color: req.Color,
		Code:  req.Code,
		Name:  req.Name,
	}

	out, err := h.uc.Update(ctx, req.Id, in)
	if err != nil {
		return nil, mapLandUseGroupError(err)
	}
	return toQHLandUsePB(out), nil
}

func (h *QHLandUseGrpcHandler) DeleteLandUse(ctx context.Context, req *tqdpb.DeleteLandUseRequest) (*emptypb.Empty, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if err := h.uc.Delete(ctx, req.Id); err != nil {
		return nil, mapLandUseGroupError(err)
	}
	return &emptypb.Empty{}, nil
}

func (h *QHLandUseGrpcHandler) ListLandUses(ctx context.Context, req *tqdpb.ListLandUseRequest) (*tqdpb.ListLandUseResponse, error) {
	var layerID *uint64
	if req != nil && req.LayerId != nil {
		layerID = req.LayerId
	}
	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)
	rows, total, err := h.uc.List(ctx, pagable, layerID, nil)
	if err != nil {
		return nil, mapLandUseGroupError(err)
	}
	data := make([]*tqdpb.QHLandUseResponse, len(rows))
	for i := range rows {
		data[i] = toQHLandUsePB(&rows[i])
	}
	return &tqdpb.ListLandUseResponse{
		Data:     data,
		Total:    total,
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}

func (h *QHLandUseGrpcHandler) ListClientLandUse(ctx context.Context, req *tqdpb.ListClientLandUseRequest) (*tqdpb.ListLandUseResponse, error) {
	if req == nil || req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layerId is required")
	}
	page := req.Page
	size := req.Size
	pagable := _dto.NewPagableFromGrpc(&page, &size, nil)
	rows, total, err := h.uc.ListClientByLayer(ctx, req.LayerId, pagable)
	if err != nil {
		return nil, mapLandUseGroupError(err)
	}
	data := make([]*tqdpb.QHLandUseResponse, len(rows))
	for i := range rows {
		data[i] = toQHLandUsePB(&rows[i])
	}
	return &tqdpb.ListLandUseResponse{
		Data:     data,
		Total:    total,
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}

func mapLandUseGroupError(err error) error {
	return _errors.ToGRPC(err)
}

func toQHLandUsePB(e *qh_domain.QHLandUse) *tqdpb.QHLandUseResponse {
	if e == nil {
		return nil
	}
	note := e.Note
	color := e.Color
	code := e.Code
	name := e.Name
	// if e.LandUse != nil {
	// 	if note == "" {
	// 		note = e.LandUse.Note
	// 	}
	// 	if color == "" {
	// 		color = e.LandUse.Color
	// 	}
	// 	if code == "" {
	// 		code = e.LandUse.Code
	// 	}
	// 	if name == "" {
	// 		name = e.LandUse.Name
	// 	}
	// }
	return &tqdpb.QHLandUseResponse{
		Id:        e.ID,
		Note:      note,
		Color:     color,
		Code:      code,
		Name:      name,
		CreatedAt: _utils.FormatTimeToString(e.CreatedAt),
		UpdatedAt: _utils.FormatTimeToString(e.UpdatedAt),
	}
}
