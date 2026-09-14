package handler_grpc

import (
	"context"
	"errors"
	"strconv"
	"strings"

	_dto "common/domain/dto"
	"common/fault"
	_utils "common/utils"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/usecase"
)

// QHLayerFamilyGrpcHandler xử lý gRPC QHLayerFamilyService.
type QHLayerFamilyGrpcHandler struct {
	tqdpb.UnimplementedQHLayerFamilyServiceServer
	uc           usecase.QHLayerFamilyUsecase
	SyncProvider *_utils.SyncUtil
}

func NewQHLayerFamilyGrpcHandler(uc usecase.QHLayerFamilyUsecase, syncProvider *_utils.SyncUtil) *QHLayerFamilyGrpcHandler {
	return &QHLayerFamilyGrpcHandler{uc: uc, SyncProvider: syncProvider}
}

// CreateLayerFamily — POST /v2/tqd/qh/admin/layer-families
func (h *QHLayerFamilyGrpcHandler) CreateLayerFamily(ctx context.Context, req *tqdpb.CreateLayerFamilyRequest) (*tqdpb.QHLayerFamilyResponse, error) {
	if req == nil {
		return nil, qhLayerFamilyValidation(
			"tqd.qh_layer_family.request_required",
			"request is required",
			"",
		)
	}
	row := &qh_domain.QHLayerFamily{
		Name:       strings.TrimSpace(req.Name),
		SortNumber: req.SortNumber,
	}
	out, err := h.uc.Create(ctx, row)
	if err != nil {
		return nil, qhLayerFamilyError(err)
	}
	return mapper.ToProtoQHLayerFamilyClient(out, 6, 16), nil
}

// GetLayerFamily — GET /v2/tqd/qh/admin/layer-families/{id}
func (h *QHLayerFamilyGrpcHandler) GetLayerFamily(ctx context.Context, req *tqdpb.GetLayerFamilyRequest) (*tqdpb.QHLayerFamilyResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, qhLayerFamilyValidation(
			"tqd.qh_layer_family.id_required",
			"id is required",
			"id",
		)
	}
	out, err := h.uc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, qhLayerFamilyError(err)
	}
	minZoom, maxZoom, err := h.uc.GetZoomRangeByFamilyID(ctx, out.ID)
	if err != nil {
		minZoom, maxZoom = 6, 16
	}
	return mapper.ToProtoQHLayerFamilyClient(out, minZoom, maxZoom), nil
}

// UpdateLayerFamily — PUT /v2/tqd/qh/admin/layer-families/{id}
func (h *QHLayerFamilyGrpcHandler) UpdateLayerFamily(ctx context.Context, req *tqdpb.UpdateLayerFamilyRequest) (*tqdpb.QHLayerFamilyResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, qhLayerFamilyValidation(
			"tqd.qh_layer_family.id_required",
			"id is required",
			"id",
		)
	}
	in := &usecase.QHLayerFamilyUpdateInput{}
	if req.Name != nil {
		in.Name = req.Name
	}
	if req.SortNumber != nil {
		in.SortNumber = req.SortNumber
	}
	if in.Name == nil && in.SortNumber == nil {
		return nil, qhLayerFamilyValidation(
			"tqd.qh_layer_family.update_fields_required",
			"at least one field to update is required",
			"",
		)
	}
	out, err := h.uc.Update(ctx, req.Id, in)
	if err != nil {
		return nil, qhLayerFamilyError(err)
	}
	minZoom, maxZoom, err := h.uc.GetZoomRangeByFamilyID(ctx, out.ID)
	if err != nil {
		minZoom, maxZoom = 6, 16
	}
	return mapper.ToProtoQHLayerFamilyClient(out, minZoom, maxZoom), nil
}

// DeleteLayerFamily — DELETE /v2/tqd/qh/admin/layer-families/{id}
func (h *QHLayerFamilyGrpcHandler) DeleteLayerFamily(ctx context.Context, req *tqdpb.DeleteLayerFamilyRequest) (*emptypb.Empty, error) {
	if req == nil || req.Id == 0 {
		return nil, qhLayerFamilyValidation(
			"tqd.qh_layer_family.id_required",
			"id is required",
			"id",
		)
	}
	if err := h.uc.Delete(ctx, req.Id); err != nil {
		return nil, qhLayerFamilyError(err)
	}
	return &emptypb.Empty{}, nil
}

// ListLayerFamilies — GET /v2/tqd/qh/admin/layer-families/list
func (h *QHLayerFamilyGrpcHandler) ListLayerFamilies(ctx context.Context, req *tqdpb.ListLayerFamiliesRequest) (*tqdpb.ListLayerFamiliesResponse, error) {
	page := uint32(0)
	size := uint32(0)
	search := ""
	var sort *string
	if req != nil {
		page = req.Page
		size = req.Size
		if req.Search != nil {
			search = *req.Search
		}
		if req.Sort != nil {
			sort = req.Sort
		}
	}
	pagable := _dto.NewPagableFromGrpc(&page, &size, sort)
	rows, total, err := h.uc.List(ctx, pagable, search)
	if err != nil {
		return nil, qhLayerFamilyError(err)
	}
	data := make([]*tqdpb.QHLayerFamilyResponse, len(rows))
	for i := range rows {
		family := &rows[i]
		minZoom, maxZoom, err := h.uc.GetZoomRangeByFamilyID(ctx, family.ID)
		if err != nil {
			minZoom, maxZoom = 6, 16
		}
		data[i] = mapper.ToProtoQHLayerFamilyClient(family, minZoom, maxZoom)
	}
	return &tqdpb.ListLayerFamiliesResponse{
		Data:     data,
		Total:    total,
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}

// GetClientLayerFamilies — GET /v2/tqd/client/qh/layer-families/list
func (h *QHLayerFamilyGrpcHandler) GetClientLayerFamilies(ctx context.Context, req *tqdpb.ListLayerFamiliesRequest) (*tqdpb.ListLayerFamiliesResponse, error) {
	return h.ListLayerFamilies(ctx, req)
}

// ListClientFamilies — GET /v2/tqd/client/families
func (h *QHLayerFamilyGrpcHandler) ListClientFamilies(ctx context.Context, req *tqdpb.ListClientFamiliesRequest) (*tqdpb.ListLayerFamiliesResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLayerFamilyList, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.GetTimestamp())
	h.SyncProvider.PutTimeRequest(ctx, key, req.GetTimestamp())
	if !updated {
		return &tqdpb.ListLayerFamiliesResponse{}, nil
	}

	page := req.GetPage()
	size := req.GetSize()
	pagable := _dto.NewPagableFromGrpc(&page, &size, nil)
	search := ""
	if req.Search != nil {
		search = *req.Search
	}

	rows, total, err := h.uc.ListClient(ctx, pagable, search)
	if err != nil {
		return nil, qhLayerFamilyError(err)
	}

	data := make([]*tqdpb.QHLayerFamilyResponse, 0, len(rows))
	for i := range rows {
		family := &rows[i]
		minZoom, maxZoom, err := h.uc.GetZoomRangeByFamilyID(ctx, family.ID)
		if err != nil {
			minZoom, maxZoom = 6, 16
		}
		data = append(data, mapper.ToProtoQHLayerFamilyClient(family, minZoom, maxZoom))
	}

	return &tqdpb.ListLayerFamiliesResponse{
		Data:     data,
		Total:    total,
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}

// BuildFamilyPMTiles — POST /v2/tqd/client/families/{familyId}/build-tile
func (h *QHLayerFamilyGrpcHandler) BuildFamilyPMTiles(req *tqdpb.BuildFamilyPMTilesRequest, stream grpc.ServerStreamingServer[tqdpb.BuildFamilyPMTilesProgress]) error {
	ctx := stream.Context()
	if req == nil || req.FamilyId == 0 {
		return qhLayerFamilyValidation(
			"tqd.qh_layer_family.family_id_required",
			"family_id is required",
			"family_id",
		)
	}

	minZoom := req.MinZoom
	maxZoom := req.MaxZoom
	if maxZoom > 16 {
		maxZoom = 16
	}
	if minZoom > 0 && maxZoom > 0 && minZoom > maxZoom {
		minZoom = maxZoom
	}

	outputDir := "./files/tiles"
	if req.OutputDir != nil && *req.OutputDir != "" {
		outputDir = *req.OutputDir
	}

	progressFn := func(p usecase.BuildProgress) {
		_ = stream.Send(&tqdpb.BuildFamilyPMTilesProgress{
			Status:       p.Status,
			Message:      p.Message,
			CurrentZoom:  p.CurrentZoom,
			TilesWritten: p.TilesWritten,
			TotalTiles:   p.TotalTiles,
			OutputPath:   p.OutputPath,
		})
	}

	err := h.uc.BuildFamilyPMTiles(ctx, req.FamilyId, minZoom, maxZoom, outputDir, progressFn)
	if err == nil {
		return nil
	}

	var buildErr *usecase.ErrFamilyBuildInProgress
	if errors.As(err, &buildErr) {
		_ = stream.Send(&tqdpb.BuildFamilyPMTilesProgress{
			Status:       buildErr.Progress.Status,
			Message:      "another family tile build is already in progress",
			CurrentZoom:  buildErr.Progress.CurrentZoom,
			TilesWritten: buildErr.Progress.TilesWritten,
			TotalTiles:   buildErr.Progress.TotalTiles,
			OutputPath:   buildErr.Progress.OutputPath,
		})
		return fault.ToGRPC(fault.New(
			fault.KindResourceExhausted,
			"tqd.qh_layer_family.build_in_progress",
			"family tile build is already in progress",
		).WithMetadata(map[string]string{
			"family_id": strconv.FormatUint(req.FamilyId, 10),
		}))
	}

	if failure, ok := fault.As(err); ok {
		_ = stream.Send(&tqdpb.BuildFamilyPMTilesProgress{
			Status:  "error",
			Message: failure.PublicMessage(),
		})
		return fault.ToGRPC(err)
	}

	_ = stream.Send(&tqdpb.BuildFamilyPMTilesProgress{
		Status:  "error",
		Message: "family tile build failed",
	})
	return fault.ToGRPC(fault.Wrap(
		err,
		fault.KindInternal,
		"tqd.qh_layer_family.build_failed",
		"family tile build failed",
	))
}

func qhLayerFamilyValidation(code, message, field string) error {
	violations := make([]fault.FieldViolation, 0, 1)
	if field != "" {
		violations = append(violations, fault.FieldViolation{Field: field, Description: message})
	}
	return fault.ToGRPC(fault.Validation(code, message, violations...))
}

func qhLayerFamilyError(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := fault.As(err); ok {
		return fault.ToGRPC(err)
	}
	return fault.ToGRPC(fault.Wrap(
		err,
		fault.KindInternal,
		"tqd.qh_layer_family.internal",
		"layer family operation failed",
	))
}
