package handler_grpc

import (
	"bytes"
	"context"
	"errors"
	"log"
	"strings"
	"time"

	_dto "common/domain/dto"
	"common/fault"
	_utils "common/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"
	"tqd/internal/usecase"
)

type RegionGrpcHandler struct {
	tqdpb.UnimplementedRegionServiceServer
	regionUsecase usecase.RegionUsecase
	labelUsecase  usecase.QHLabelUsecase
	regionMapper  *mapper.RegionMapper
	SyncProvider  *_utils.SyncUtil
}

func NewRegionGrpcHandler(
	regionUsecase usecase.RegionUsecase,
	labelUsecase usecase.QHLabelUsecase,
	syncProvider *_utils.SyncUtil,
) *RegionGrpcHandler {
	return &RegionGrpcHandler{
		regionUsecase: regionUsecase,
		labelUsecase:  labelUsecase,
		regionMapper:  mapper.NewRegionMapper(),
		SyncProvider:  syncProvider,
	}
}

func mapRegionError(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := fault.As(err); ok {
		return fault.ToGRPC(err)
	}

	return fault.ToGRPC(fault.Wrap(
		err,
		fault.KindInternal,
		"tqd.region.internal",
		"region operation failed",
	))
}

func buildAdminRegionFilter(
	pagable *_dto.Pagable,
	layerID *uint64,
	labelID *uint64,
	status *uint32,
	processingStatus *uint32,
	search *string,
	orderBy *string,
	orderDir *string,
) *repo.RegionFilter {
	filter := &repo.RegionFilter{Pagable: pagable}
	if layerID != nil {
		filter.LayerID = layerID
	}
	if labelID != nil {
		filter.LabelID = labelID
	}
	if status != nil {
		s := *status
		filter.Status = &s
	}
	if processingStatus != nil {
		ps := *processingStatus
		filter.ProcessingStatus = &ps
	}
	if search != nil && *search != "" {
		filter.Search = search
	}
	if orderBy != nil {
		filter.OrderBy = *orderBy
	}
	if orderDir != nil {
		filter.OrderDir = *orderDir
	}
	return filter
}

func (h *RegionGrpcHandler) buildListRegionsAdminResponse(
	ctx context.Context,
	regions []qh_domain.QHRegion,
	total int64,
	pagable *_dto.Pagable,
) *tqdpb.ListRegionsResponse {
	protoRegions := make([]*tqdpb.RegionResponse, 0, len(regions))
	for i := range regions {
		region := &regions[i]
		var label *qh_domain.QHLabel
		if region.LabelID != nil && region.LabelName != "" {
			label = &qh_domain.QHLabel{
				ID:    *region.LabelID,
				Name:  region.LabelName,
				Color: region.LabelColor,
			}
		}
		protoRegions = append(protoRegions, h.regionMapper.ToProtoResponse(region, label))
	}
	return &tqdpb.ListRegionsResponse{
		Data:     protoRegions,
		Total:    total,
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}
}

// =====================================================
// ADMIN APIS
// =====================================================

// CreateRegion — POST /v2/tqd/admin/regions — Tạo region (layerId, name, geometry GeoJSON: Polygon, MultiPolygon hoặc Feature).
func (h *RegionGrpcHandler) CreateRegion(ctx context.Context, req *tqdpb.CreateRegionRequest) (*tqdpb.RegionResponse, error) {
	if req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layer_id is required")
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if len(bytes.TrimSpace(req.Geometry)) == 0 {
		return nil, status.Error(codes.InvalidArgument, "geometry is required")
	}

	log.Printf("[CreateRegion] LayerID: %d, Name: %s", req.LayerId, req.Name)

	region, err := h.regionUsecase.Create(
		ctx,
		req.LayerId,
		req.Name,
		req.Geometry,
		req.LabelId,
	)
	if err != nil {
		return nil, mapRegionError(err)
	}
	if region == nil {
		return nil, status.Error(codes.Internal, "create region returned nil")
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDRegionList, 0), t.UnixMilli())
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDRegionDetail, region.ID), t.UnixMilli())

	var label *qh_domain.QHLabel
	if region.LabelID != nil {
		label, _ = h.regionUsecase.GetLabelByRegion(ctx, region.ID)
	}
	return h.regionMapper.ToProtoResponse(region, label), nil
}

// UpdateRegion — PATCH /v2/tqd/admin/regions/{id} — Cập nhật metadata (không đổi geometry).
func (h *RegionGrpcHandler) UpdateRegion(ctx context.Context, req *tqdpb.UpdateRegionRequest) (*tqdpb.RegionResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	patch := usecase.RegionUpdatePatch{RegionID: req.Id}
	if req.Name != nil {
		patch.Name = req.Name
	}
	if req.DisplayName != nil {
		patch.DisplayName = req.DisplayName
	}
	if req.Description != nil {
		patch.Description = req.Description
	}
	if req.Status != nil {
		st := enums.RegionStatus(*req.Status)
		patch.Status = &st
	}
	if req.LabelId != nil {
		lid := *req.LabelId
		patch.LabelID = &lid
	}

	log.Printf("[UpdateRegion] RegionID: %d", req.Id)

	region, err := h.regionUsecase.Update(ctx, patch)
	if err != nil {
		return nil, mapRegionError(err)
	}
	if region == nil {
		return nil, status.Error(codes.Internal, "update region returned nil")
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDRegionList, 0), t.UnixMilli())
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDRegionDetail, req.Id), t.UnixMilli())

	var label *qh_domain.QHLabel
	if region.LabelID != nil {
		label, _ = h.regionUsecase.GetLabelByRegion(ctx, region.ID)
	}
	return h.regionMapper.ToProtoResponse(region, label), nil
}

func (h *RegionGrpcHandler) GetRegion(ctx context.Context, req *tqdpb.GetRegionRequest) (*tqdpb.RegionResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	key := h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDRegionDetail, req.Id)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.RegionResponse{}, nil
	}

	log.Printf("[GetRegion] RegionID: %d", req.Id)

	region, err := h.regionUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if region == nil {
		return nil, status.Error(codes.NotFound, "region not found")
	}

	label, _ := h.regionUsecase.GetLabelByRegion(ctx, req.Id)

	return h.regionMapper.ToProtoResponse(region, label), nil
}

func (h *RegionGrpcHandler) ListRegions(ctx context.Context, req *tqdpb.ListRegionsRequest) (*tqdpb.ListRegionsResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDRegionList, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ListRegionsResponse{}, nil
	}

	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, req.OrderBy)
	filter := buildAdminRegionFilter(
		pagable,
		req.LayerId,
		req.LabelId,
		req.Status,
		req.ProcessingStatus,
		req.Search,
		req.OrderBy,
		req.OrderDir,
	)

	regions, total, err := h.regionUsecase.List(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.buildListRegionsAdminResponse(ctx, regions, total, pagable), nil
}

// ListAdminRegionsByLayer - Danh sách region theo layer (admin), kiểm tra layer tồn tại
func (h *RegionGrpcHandler) ListAdminRegionsByLayer(ctx context.Context, req *tqdpb.ListAdminRegionsByLayerRequest) (*tqdpb.ListRegionsResponse, error) {
	if req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layer_id is required")
	}

	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDRegionList, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ListRegionsResponse{}, nil
	}

	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, req.OrderBy)
	layerID := req.LayerId
	layerPtr := &layerID
	filter := buildAdminRegionFilter(
		pagable,
		layerPtr,
		req.LabelId,
		req.Status,
		req.ProcessingStatus,
		req.Search,
		req.OrderBy,
		req.OrderDir,
	)

	regions, total, err := h.regionUsecase.ListForAdminByLayer(ctx, req.LayerId, filter)
	if err != nil {
		if errors.Is(err, usecase.ErrLayerNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.buildListRegionsAdminResponse(ctx, regions, total, pagable), nil
}

func (h *RegionGrpcHandler) DeleteRegion(ctx context.Context, req *tqdpb.DeleteRegionRequest) (*emptypb.Empty, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	log.Printf("[DeleteRegion] RegionID: %d", req.Id)

	if err := h.regionUsecase.Delete(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDRegionList, 0), t.UnixMilli())
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDRegionDetail, req.Id), t.UnixMilli())

	return &emptypb.Empty{}, nil
}

// SyncRegion — POST /v2/tqd/admin/regions/sync-region — Đồng bộ landUseId/legendId cho regions theo layerId và labelId.
func (h *RegionGrpcHandler) SyncRegion(ctx context.Context, req *tqdpb.SyncRegionRequest) (*tqdpb.SyncRegionResponse, error) {
	if req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layerId is required")
	}
	if req.LabelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "labelId is required")
	}

	log.Printf("[SyncRegion] LayerID: %d, LabelID: %d", req.LayerId, req.LabelId)

	result, err := h.regionUsecase.SyncRegion(
		ctx,
		req.LayerId,
		req.LabelId,
	)
	if err != nil {
		return nil, mapRegionError(err)
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDRegionList, 0), t.UnixMilli())

	return &tqdpb.SyncRegionResponse{
		LayerId:   result.LayerID,
		LabelId:   result.LabelID,
		LandUseId: result.LandUseID,
		LegendId:  result.LegendID,
		Total:     result.SyncedCount,
	}, nil
}

func (h *RegionGrpcHandler) UpdateRegionProcessingStatus(ctx context.Context, req *tqdpb.UpdateProcessingStatusRequest) (*tqdpb.RegionResponse, error) {
	if req.RegionId == 0 {
		return nil, status.Error(codes.InvalidArgument, "regionId is required")
	}

	log.Printf("[UpdateRegionProcessingStatus] RegionID: %d, Status: %d", req.RegionId, req.ProcessingStatus)

	userID := _utils.GetOriginIdFromContext(ctx)
	notes := ""
	if req.Notes != nil {
		notes = *req.Notes
	}

	err := h.regionUsecase.UpdateProcessingStatus(ctx, req.RegionId, enums.ProcessingStatus(req.ProcessingStatus), userID, notes)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDRegionList, 0), t.UnixMilli())
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDRegionDetail, req.RegionId), t.UnixMilli())

	region, err := h.regionUsecase.GetByID(ctx, req.RegionId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	label, _ := h.regionUsecase.GetLabelByRegion(ctx, req.RegionId)

	return h.regionMapper.ToProtoResponse(region, label), nil
}

// =====================================================
// CLIENT APIS
// =====================================================

func (h *RegionGrpcHandler) ListClientRegions(ctx context.Context, req *tqdpb.ListClientRegionsRequest) (*tqdpb.ListRegionsResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDRegionList, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ListRegionsResponse{}, nil
	}

	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)

	var bBox *usecase.BBox
	if req.MinLng != nil && req.MinLat != nil && req.MaxLng != nil && req.MaxLat != nil {
		bBox = &usecase.BBox{
			MinLng: *req.MinLng,
			MinLat: *req.MinLat,
			MaxLng: *req.MaxLng,
			MaxLat: *req.MaxLat,
		}
	}

	var layerID uint64
	if req.LayerId != nil {
		layerID = *req.LayerId
	}

	regions, total, err := h.regionUsecase.ListClient(ctx, layerID, bBox, pagable.GetPage(), pagable.GetSize())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoRegions := make([]*tqdpb.RegionResponse, 0, len(regions))
	for _, region := range regions {
		protoRegions = append(protoRegions, h.regionMapper.ToProtoClientResponse(&region))
	}

	return &tqdpb.ListRegionsResponse{
		Data:     protoRegions,
		Total:    total,
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}

func (h *RegionGrpcHandler) FindRegionByPoint(ctx context.Context, req *tqdpb.FindRegionByPointRequest) (*tqdpb.RegionResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDRegionList, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.RegionResponse{}, nil
	}

	if req.Lat < -90 || req.Lat > 90 || req.Lng < -180 || req.Lng > 180 {
		return nil, status.Error(codes.InvalidArgument, "invalid coordinates")
	}

	var layerID *uint64
	if req.LayerId != nil {
		val := uint64(*req.LayerId)
		layerID = &val
	}

	region, err := h.regionUsecase.FindByPoint(ctx, req.Lat, req.Lng, layerID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if region == nil {
		return nil, status.Error(codes.NotFound, "no region found")
	}

	return h.regionMapper.ToProtoClientResponse(region), nil
}

func (h *RegionGrpcHandler) FindRegionsByBBox(ctx context.Context, req *tqdpb.FindRegionsByBBoxRequest) (*tqdpb.ListRegionsResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDRegionList, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ListRegionsResponse{}, nil
	}

	if req.Bbox == nil || req.Bbox.MinLng > req.Bbox.MaxLng || req.Bbox.MinLat > req.Bbox.MaxLat {
		return nil, status.Error(codes.InvalidArgument, "invalid bbox")
	}

	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)

	var layerID, labelID *uint64
	if req.LayerId != nil {
		val := uint64(*req.LayerId)
		layerID = &val
	}
	if req.LabelId != nil {
		val := uint64(*req.LabelId)
		labelID = &val
	}

	regions, total, err := h.regionUsecase.FindByBBox(
		ctx,
		req.Bbox.MinLng, req.Bbox.MinLat,
		req.Bbox.MaxLng, req.Bbox.MaxLat,
		layerID, labelID,
		int(pagable.GetPage()), int(pagable.GetSize()),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoRegions := make([]*tqdpb.RegionResponse, 0, len(regions))
	for _, region := range regions {
		protoRegions = append(protoRegions, h.regionMapper.ToProtoClientResponse(&region))
	}

	return &tqdpb.ListRegionsResponse{
		Data:     protoRegions,
		Total:    total,
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}
