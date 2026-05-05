package handler_grpc

import (
	"context"
	"log"

	_dto "common/domain/dto"
	_utils "common/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"
	"tqd/internal/usecase"
)

// RegionGrpcHandler xử lý gRPC cho RegionService
type RegionGrpcHandler struct {
	tqdpb.UnimplementedRegionServiceServer
	regionUsecase usecase.RegionUsecase
	labelUsecase  usecase.QHLabelUsecase
	regionMapper  *mapper.RegionMapper
}

func NewRegionGrpcHandler(
	regionUsecase usecase.RegionUsecase,
	labelUsecase usecase.QHLabelUsecase,
) *RegionGrpcHandler {
	return &RegionGrpcHandler{
		regionUsecase: regionUsecase,
		labelUsecase:  labelUsecase,
		regionMapper:  mapper.NewRegionMapper(),
	}
}

// =====================================================
// ADMIN APIS
// =====================================================

// GetRegion - Lấy chi tiết region
func (h *RegionGrpcHandler) GetRegion(ctx context.Context, req *tqdpb.GetRegionRequest) (*tqdpb.RegionResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	log.Printf("[GetRegion] RegionID: %d", req.Id)

	region, err := h.regionUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if region == nil {
		return nil, status.Error(codes.NotFound, "region not found")
	}

	labels, _ := h.regionUsecase.GetLabelsByRegion(ctx, req.Id)

	return h.regionMapper.ToProtoResponse(region, labels), nil
}

// ListRegions - Danh sách regions với filter
func (h *RegionGrpcHandler) ListRegions(ctx context.Context, req *tqdpb.ListRegionsRequest) (*tqdpb.ListRegionsResponse, error) {
	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, req.OrderBy)

	filter := &repo.RegionFilter{Pagable: pagable}

	if req.LayerId != nil {
		filter.LayerID = req.LayerId
	}
	if req.LabelId != nil {
		filter.LabelID = req.LabelId
	}
	if req.Status != nil {
		statusVal := *req.Status
		filter.Status = &statusVal
	}
	if req.ProcessingStatus != nil {
		processingStatus := *req.ProcessingStatus
		filter.ProcessingStatus = &processingStatus
	}
	if req.Search != nil && *req.Search != "" {
		filter.Search = req.Search
	}
	if req.OrderBy != nil {
		filter.OrderBy = *req.OrderBy
	}
	if req.OrderDir != nil {
		filter.OrderDir = *req.OrderDir
	}

	log.Printf("[ListRegions] Filter: layerId=%v, labelId=%v, status=%v, processingStatus=%v, search=%v, page=%d, size=%d",
		filter.LayerID, filter.LabelID, filter.Status, filter.ProcessingStatus, filter.Search, pagable.GetPage(), pagable.GetSize())

	regions, total, err := h.regionUsecase.List(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoRegions := make([]*tqdpb.RegionResponse, 0, len(regions))
	for _, region := range regions {
		labels, _ := h.regionUsecase.GetLabelsByRegion(ctx, region.ID)
		protoRegions = append(protoRegions, h.regionMapper.ToProtoResponse(&region, labels))
	}

	return &tqdpb.ListRegionsResponse{
		Data:     protoRegions,
		Total:    total,
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}

// DeleteRegion - Xóa mềm region
func (h *RegionGrpcHandler) DeleteRegion(ctx context.Context, req *tqdpb.DeleteRegionRequest) (*emptypb.Empty, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	log.Printf("[DeleteRegion] RegionID: %d", req.Id)

	if err := h.regionUsecase.Delete(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

// UpdateRegionProcessingStatus - Cập nhật trạng thái xử lý của region
func (h *RegionGrpcHandler) UpdateRegionProcessingStatus(ctx context.Context, req *tqdpb.UpdateProcessingStatusRequest) (*tqdpb.RegionResponse, error) {
	if req.RegionId == 0 {
		return nil, status.Error(codes.InvalidArgument, "regionId is required")
	}

	log.Printf("[UpdateRegionProcessingStatus] RegionID: %d, Status: %d", req.RegionId, req.ProcessingStatus)

	userID := _utils.GetOriginIdFromContext(ctx)
	err := h.regionUsecase.UpdateProcessingStatus(ctx, req.RegionId, enums.ProcessingStatus(req.ProcessingStatus), userID, *req.Notes)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	region, err := h.regionUsecase.GetByID(ctx, req.RegionId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	labels, _ := h.regionUsecase.GetLabelsByRegion(ctx, req.RegionId)
	return h.regionMapper.ToProtoResponse(region, labels), nil
}

// =====================================================
// CLIENT APIS
// =====================================================

// ListClientRegions - Danh sách regions public
func (h *RegionGrpcHandler) ListClientRegions(ctx context.Context, req *tqdpb.ListClientRegionsRequest) (*tqdpb.ListRegionsResponse, error) {
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

// FindRegionByPoint - Tìm region theo tọa độ điểm
func (h *RegionGrpcHandler) FindRegionByPoint(ctx context.Context, req *tqdpb.FindRegionByPointRequest) (*tqdpb.RegionResponse, error) {
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

// FindRegionsByBBox - Tìm regions theo bounding box
func (h *RegionGrpcHandler) FindRegionsByBBox(ctx context.Context, req *tqdpb.FindRegionsByBBoxRequest) (*tqdpb.ListRegionsResponse, error) {
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
