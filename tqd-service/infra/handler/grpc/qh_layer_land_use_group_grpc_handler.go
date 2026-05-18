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

type QHLayerLandUseGroupGrpcHandler struct {
	tqdpb.UnimplementedQHLayerLandUseGroupServiceServer
	uc usecase.QHLayerLandUseGroupUsecase
}

func NewQHLayerLandUseGroupGrpcHandler(uc usecase.QHLayerLandUseGroupUsecase) *QHLayerLandUseGroupGrpcHandler {
	return &QHLayerLandUseGroupGrpcHandler{uc: uc}
}

func (h *QHLayerLandUseGroupGrpcHandler) CreateLayerLandUseGroup(ctx context.Context, req *tqdpb.CreateLayerLandUseGroupRequest) (*tqdpb.QHLayerLandUseGroupResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	row := &qh_domain.QHLayerLandUseGroup{
		LayerID:        req.LayerId,
		LandUseGroupID: req.LandUseGroupId,
		Note:           req.Note,
		DisplayOrder:   int(req.DisplayOrder),
		IsVisible:      req.IsVisible,
		Color:          req.Color,
	}
	out, err := h.uc.Create(ctx, row)
	if err != nil {
		return nil, mapLandUseGroupError(err)
	}
	return toQHLayerLandUseGroupPB(out), nil
}

func (h *QHLayerLandUseGroupGrpcHandler) BatchCreateLayerLandUseGroups(ctx context.Context, req *tqdpb.BatchCreateLayerLandUseGroupsRequest) (*tqdpb.BatchCreateLayerLandUseGroupsResponse, error) {
	if req == nil || req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layerId is required")
	}
	items := make([]qh_domain.QHLayerLandUseGroup, len(req.Items))
	for i, it := range req.Items {
		items[i] = qh_domain.QHLayerLandUseGroup{
			LandUseGroupID: it.LandUseGroupId,
			Note:           it.Note,
			DisplayOrder:   int(it.DisplayOrder),
			IsVisible:      it.IsVisible,
			Color:          it.Color,
		}
	}
	created, skipped, err := h.uc.BatchCreate(ctx, req.LayerId, items)
	if err != nil {
		return nil, mapLandUseGroupError(err)
	}
	data := make([]*tqdpb.QHLayerLandUseGroupResponse, len(created))
	for i := range created {
		data[i] = toQHLayerLandUseGroupPB(&created[i])
	}
	return &tqdpb.BatchCreateLayerLandUseGroupsResponse{
		Created:         data,
		SkippedGroupIds: skipped,
	}, nil
}

func (h *QHLayerLandUseGroupGrpcHandler) GetLayerLandUseGroup(ctx context.Context, req *tqdpb.GetLayerLandUseGroupRequest) (*tqdpb.QHLayerLandUseGroupResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	out, err := h.uc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, mapLandUseGroupError(err)
	}
	return toQHLayerLandUseGroupPB(out), nil
}

func (h *QHLayerLandUseGroupGrpcHandler) UpdateLayerLandUseGroup(ctx context.Context, req *tqdpb.UpdateLayerLandUseGroupRequest) (*tqdpb.QHLayerLandUseGroupResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	in := &usecase.QHLayerLandUseGroupUpdateInput{}
	hasUpdate := false
	if req.LayerId != nil {
		in.LayerID = req.LayerId
		hasUpdate = true
	}
	if req.LandUseGroupId != nil {
		in.LandUseGroupID = req.LandUseGroupId
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
	if req.Color != nil {
		in.Color = req.Color
		hasUpdate = true
	}
	if !hasUpdate {
		return nil, status.Error(codes.InvalidArgument, "at least one field to update is required")
	}
	out, err := h.uc.Update(ctx, req.Id, in)
	if err != nil {
		return nil, mapLandUseGroupError(err)
	}
	return toQHLayerLandUseGroupPB(out), nil
}

func (h *QHLayerLandUseGroupGrpcHandler) DeleteLayerLandUseGroup(ctx context.Context, req *tqdpb.DeleteLayerLandUseGroupRequest) (*emptypb.Empty, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if err := h.uc.Delete(ctx, req.Id); err != nil {
		return nil, mapLandUseGroupError(err)
	}
	return &emptypb.Empty{}, nil
}

func (h *QHLayerLandUseGroupGrpcHandler) ListLayerLandUseGroups(ctx context.Context, req *tqdpb.ListLayerLandUseGroupsRequest) (*tqdpb.ListLayerLandUseGroupsResponse, error) {
	page := uint32(0)
	size := uint32(0)
	var layerID, groupID *uint64
	if req != nil {
		page = req.Page
		size = req.Size
		if req.LayerId != nil {
			layerID = req.LayerId
		}
		if req.LandUseGroupId != nil {
			groupID = req.LandUseGroupId
		}
	}
	pagable := _dto.NewPagableFromGrpc(&page, &size, nil)
	rows, total, err := h.uc.List(ctx, pagable, layerID, groupID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	data := make([]*tqdpb.QHLayerLandUseGroupResponse, len(rows))
	for i := range rows {
		data[i] = toQHLayerLandUseGroupPB(&rows[i])
	}
	return &tqdpb.ListLayerLandUseGroupsResponse{
		Data:     data,
		Total:    total,
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}

func (h *QHLayerLandUseGroupGrpcHandler) ListClientLayerLandUseGroups(ctx context.Context, req *tqdpb.ListClientLayerLandUseGroupsRequest) (*tqdpb.ListLayerLandUseGroupsResponse, error) {
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
	data := make([]*tqdpb.QHLayerLandUseGroupResponse, len(rows))
	for i := range rows {
		data[i] = toQHLayerLandUseGroupPB(&rows[i])
	}
	return &tqdpb.ListLayerLandUseGroupsResponse{
		Data:     data,
		Total:    total,
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}

func mapLandUseGroupError(err error) error {
	msg := err.Error()
	if strings.Contains(msg, "already has land use group") {
		return status.Error(codes.AlreadyExists, msg)
	}
	if strings.Contains(msg, "not found") {
		return status.Error(codes.NotFound, msg)
	}
	if strings.Contains(msg, "required") {
		return status.Error(codes.InvalidArgument, msg)
	}
	return status.Error(codes.Internal, msg)
}

func toQHLayerLandUseGroupPB(e *qh_domain.QHLayerLandUseGroup) *tqdpb.QHLayerLandUseGroupResponse {
	if e == nil {
		return nil
	}
	resp := &tqdpb.QHLayerLandUseGroupResponse{
		Id:             e.ID,
		LayerId:        e.LayerID,
		LandUseGroupId: e.LandUseGroupID,
		Note:           e.Note,
		DisplayOrder:   int32(e.DisplayOrder),
		IsVisible:      e.IsVisible,
		Color:          e.GetDisplayColor(),
		CreatedAt:      _utils.FormatTimeToString(e.CreatedAt),
		UpdatedAt:      _utils.FormatTimeToString(e.UpdatedAt),
	}
	if e.LandUseGroup != nil {
		resp.LandUseGroupInfo = &tqdpb.LandUseGroupInfo{
			Id:          e.LandUseGroup.ID,
			Code:        e.LandUseGroup.Code,
			Name:        e.LandUseGroup.Name,
			Description: e.LandUseGroup.Description,
			Color:       e.LandUseGroup.Color,
			CanBuild:    e.LandUseGroup.CanBuild,
			Priority:    int32(e.LandUseGroup.Priority),
		}
	}
	return resp
}
