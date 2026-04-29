package handler_grpc

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"encoding/json"
	"log"

	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
	"tqd/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RegionGrpcHandler struct {
	tqdpb.UnimplementedRegionServiceServer
	regionUsecase usecase.RegionUsecase
}

func NewRegionGrpcHandler(regionUsecase usecase.RegionUsecase) *RegionGrpcHandler {
	return &RegionGrpcHandler{
		regionUsecase: regionUsecase,
	}
}

func (h *RegionGrpcHandler) CreateRegion(ctx context.Context, req *tqdpb.CreateRegionRequest) (*tqdpb.RegionResponse, error) {
	// Validate
	if req.Name == "" {
		return nil, _errors.ReturnError(400, "name is required")
	}
	if req.DisplayName == "" {
		return nil, _errors.ReturnError(400, "displayName is required")
	}
	if req.Geometry == "" {
		return nil, _errors.ReturnError(400, "geometry is required")
	}
	if req.Type < 1 || req.Type > 4 {
		return nil, _errors.ReturnError(400, "type must be 1-4")
	}
	if req.Status < 10 || req.Status > 30 {
		return nil, _errors.ReturnError(400, "status must be 10,20,30")
	}

	createReq := h.regionMapper.ToCreateRegionDTO(req)

	region, err := h.regionUsecase.Create(ctx, createReq)
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}

	return h.regionMapper.ToProtoFromDomain(region), nil
}

func (h *RegionGrpcHandler) CreateBatchRegion(ctx context.Context, req *tqdpb.RegionBatchRequest) (*tqdpb.ListRegionsResponse, error) {
	createReq := dto.CreateRegionBatchRequest{
		Items: []dto.CreateRegionBatchItem{},
	}
	for i := 0; i < len(req.Geometry); i++ {
		it := req.Geometry[i]
		createReq.Items = append(createReq.Items, dto.CreateRegionBatchItem{
			Geometry: it,
			LayerID:  req.LayerId,
		})
	}

	region, err := h.regionUsecase.CreateBatch(ctx, &createReq)
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}

	items := make([]*tqdpb.RegionResponse, len(region.Items))
	// for i := 0; i < len(req.Geometry); i++ {
	// it := region.Items[i]
	// items[i] = &tqdpb.RegionResponse{
	// 	Geometry: it.Region.Geometry,
	// 	Id:       uint32(it.Region.ID),
	// }
	// }
	return &tqdpb.ListRegionsResponse{
		Data: items,
	}, nil
}

func (h *RegionGrpcHandler) UpdateRegion(ctx context.Context, req *tqdpb.UpdateRegionRequest) (*tqdpb.RegionResponse, error) {
	if req.Id == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}

	// Convert to DTO using mapper
	updateReq := h.regionMapper.ToUpdateRegionDTO(req)

	region, err := h.regionUsecase.Update(ctx, uint64(req.Id), updateReq)
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}

	return h.regionMapper.ToProtoFromDomain(region), nil
}

func (h *RegionGrpcHandler) DeleteRegion(ctx context.Context, req *tqdpb.DeleteRegionRequest) (*sharepb.Empty, error) {
	if req.Id == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}

	if err := h.regionUsecase.Delete(ctx, uint64(req.Id)); err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}
	return &sharepb.Empty{}, nil
}

func (h *RegionGrpcHandler) GetRegion(ctx context.Context, req *tqdpb.GetRegionRequest) (*tqdpb.RegionResponse, error) {
	if req.Id == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}

	region, err := h.regionUsecase.GetByID(ctx, uint64(req.Id))
	if err != nil {
		return nil, _errors.ReturnError(404, err.Error())
	}
	if region == nil {
		return nil, _errors.ReturnError(404, "region not found")
	}

	return h.regionMapper.ToProtoFromDTO(region), nil
}

// ListRegions - Admin: Danh sách regions với filter
func (h *RegionGrpcHandler) ListRegions(ctx context.Context, req *tqdpb.ListRegionsRequest) (*tqdpb.ListRegionsResponse, error) {
	// Build pagable
	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, req.OrderBy)

	// Build filter
	filter := &repo.RegionFilter{
		Pagable: pagable,
	}

	if req.LayerId != nil {
		filter.LayerID = req.LayerId
	}
	if req.LandUseCode != nil {
		filter.LandUseCode = req.LandUseCode
	}
	if req.Status != nil {
		statusVal := int(*req.Status)
		filter.Status = &statusVal
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

	log.Printf("[ListRegions] Filter: layer_id=%v, land_use=%v, page=%d, size=%d",
		filter.LayerID, filter.LandUseCode, pagable.GetPage(), pagable.GetSize())

	// Call usecase
	regions, total, err := h.regionUsecase.List(ctx, filter)
	if err != nil {
		log.Printf("[ListRegions] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert to proto
	protoRegions := make([]*tqdpb.RegionResponse, 0, len(regions))
	for _, region := range regions {
		protoRegions = append(protoRegions, h.toProtoRegion(&region))
	}

	log.Printf("[ListRegions] Found %d regions, total: %d", len(protoRegions), total)

	return &tqdpb.ListRegionsResponse{
		Data:  protoRegions,
		Total: total,
		Page:  int32(pagable.GetPage()),
		Size:  int32(pagable.GetSize()),
	}, nil
}

// GetRegion - Admin: Chi tiết region
func (h *RegionGrpcHandler) GetRegion(ctx context.Context, req *tqdpb.GetRegionRequest) (*tqdpb.RegionResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	log.Printf("[GetRegion] RegionID: %d", req.Id)

	region, err := h.regionUsecase.GetByID(ctx, req.Id)
	if err != nil {
		log.Printf("[GetRegion] Error: %v", err)
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if region == nil {
		return nil, status.Error(codes.NotFound, "region not found")
	}

	return h.toProtoRegion(region), nil
}

// GetRegionStats - Admin: Thống kê regions theo layer
func (h *RegionGrpcHandler) GetRegionStats(ctx context.Context, req *tqdpb.GetRegionStatsRequest) (*tqdpb.RegionStatsResponse, error) {
	if req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layer_id is required")
	}

	log.Printf("[GetRegionStats] LayerID: %d", req.LayerId)

	stats, err := h.regionUsecase.GetStatsByLayer(ctx, req.LayerId)
	if err != nil {
		log.Printf("[GetRegionStats] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Build response
	protoStats := &tqdpb.RegionStatsResponse{
		LayerId:      stats.LayerID,
		TotalRegions: stats.TotalRegions,
		TotalAreaHa:  stats.TotalAreaHa,
		AvgAreaHa:    stats.AvgAreaHa,
		MinAreaHa:    stats.MinAreaHa,
		MaxAreaHa:    stats.MaxAreaHa,
	}

	for _, ls := range stats.LandUseStats {
		protoStats.LandUseStats = append(protoStats.LandUseStats, &tqdpb.LandUseStat{
			LandUseCode: ls.LandUseCode,
			LandUseName: ls.LandUseName,
			Count:       ls.Count,
			TotalAreaHa: ls.TotalAreaHa,
		})
	}

	log.Printf("[GetRegionStats] Total regions: %d, Total area: %.2f ha", stats.TotalRegions, stats.TotalAreaHa)

	return protoStats, nil
}

// =====================================================
// CLIENT APIS (Public)
// =====================================================

// ListClientRegions - Client: Danh sách regions (chỉ active)
func (h *RegionGrpcHandler) ListClientRegions(ctx context.Context, req *tqdpb.ListClientRegionsRequest) (*tqdpb.ListRegionsResponse, error) {
	// Build pagable
	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)

	// Build filter
	filter := &repo.ClientRegionFilter{
		Pagable: pagable,
	}

	if req.LayerId != nil {
		filter.LayerID = req.LayerId
	}
	if req.LandUseCode != nil {
		filter.LandUseCode = req.LandUseCode
	}
	if req.Bbox != nil {
		filter.BBox = &repo.BBox{
			MinLng: req.Bbox.MinLng,
			MinLat: req.Bbox.MinLat,
			MaxLng: req.Bbox.MaxLng,
			MaxLat: req.Bbox.MaxLat,
		}
	}

	log.Printf("[ListClientRegions] Filter: layer_id=%v, land_use=%v, bbox=%v, page=%d, size=%d",
		filter.LayerID, filter.LandUseCode, filter.BBox != nil, pagable.GetPage(), pagable.GetSize())

	// Call usecase
	regions, total, err := h.regionUsecase.ListClient(ctx, filter)
	if err != nil {
		log.Printf("[ListClientRegions] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert to proto (client version)
	protoRegions := make([]*tqdpb.RegionResponse, 0, len(regions))
	for _, region := range regions {
		protoRegions = append(protoRegions, h.toProtoRegionClient(&region))
	}

	log.Printf("[ListClientRegions] Found %d regions, total: %d", len(protoRegions), total)

	return &tqdpb.ListRegionsResponse{
		Data:  protoRegions,
		Total: total,
		Page:  int32(pagable.GetPage()),
		Size:  int32(pagable.GetSize()),
	}, nil
}

// GetClientRegion - Client: Chi tiết region (chỉ active)
func (h *RegionGrpcHandler) GetClientRegion(ctx context.Context, req *tqdpb.GetClientRegionRequest) (*tqdpb.RegionResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	log.Printf("[GetClientRegion] RegionID: %d", req.Id)

	region, err := h.regionUsecase.GetClientByID(ctx, req.Id)
	if err != nil {
		log.Printf("[GetClientRegion] Error: %v", err)
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if region == nil {
		return nil, status.Error(codes.NotFound, "region not found")
	}

	return h.toProtoRegionClient(region), nil
}

// FindRegionByPoint - Client: Tìm region theo tọa độ điểm
func (h *RegionGrpcHandler) FindRegionByPoint(ctx context.Context, req *tqdpb.FindRegionByPointRequest) (*tqdpb.RegionResponse, error) {
	// Validate coordinates
	if req.Lat < -90 || req.Lat > 90 {
		return nil, status.Error(codes.InvalidArgument, "latitude must be between -90 and 90")
	}
	if req.Lng < -180 || req.Lng > 180 {
		return nil, status.Error(codes.InvalidArgument, "longitude must be between -180 and 180")
	}

	log.Printf("[FindRegionByPoint] Lat: %f, Lng: %f", req.Lat, req.Lng)

	region, err := h.regionUsecase.FindByPoint(ctx, req.Lat, req.Lng)
	if err != nil {
		log.Printf("[FindRegionByPoint] Error: %v", err)
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if region == nil {
		return nil, status.Error(codes.NotFound, "no region found at this location")
	}

	return h.toProtoRegionClient(region), nil
}

// FindRegionsByBBox - Client: Tìm regions theo bounding box
func (h *RegionGrpcHandler) FindRegionsByBBox(ctx context.Context, req *tqdpb.FindRegionsByBBoxRequest) (*tqdpb.ListRegionsResponse, error) {
	// Validate bbox
	if req.Bbox == nil {
		return nil, status.Error(codes.InvalidArgument, "bbox is required")
	}
	if req.Bbox.MinLng > req.Bbox.MaxLng || req.Bbox.MinLat > req.Bbox.MaxLat {
		return nil, status.Error(codes.InvalidArgument, "invalid bbox coordinates")
	}

	// Build pagable
	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)

	// Build filter
	filter := &repo.RegionFilter{
		Pagable: pagable,
	}
	if req.LandUseCode != nil {
		filter.LandUseCode = req.LandUseCode
	}

	log.Printf("[FindRegionsByBBox] BBox: [%f,%f,%f,%f], land_use=%v, page=%d, size=%d",
		req.Bbox.MinLng, req.Bbox.MinLat, req.Bbox.MaxLng, req.Bbox.MaxLat,
		filter.LandUseCode, pagable.GetPage(), pagable.GetSize())

	// Call usecase
	regions, total, err := h.regionUsecase.FindByBBox(
		ctx,
		req.Bbox.MinLng, req.Bbox.MinLat,
		req.Bbox.MaxLng, req.Bbox.MaxLat,
		filter,
	)
	if err != nil {
		log.Printf("[FindRegionsByBBox] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert to proto (client version)
	protoRegions := make([]*tqdpb.RegionResponse, 0, len(regions))
	for _, region := range regions {
		protoRegions = append(protoRegions, h.toProtoRegionClient(&region))
	}

	log.Printf("[FindRegionsByBBox] Found %d regions, total: %d", len(protoRegions), total)

	return &tqdpb.ListRegionsResponse{
		Data:  protoRegions,
		Total: total,
		Page:  int32(pagable.GetPage()),
		Size:  int32(pagable.GetSize()),
	}, nil
}

// =====================================================
// Helper Functions
// =====================================================

// toProtoRegion - Convert domain to proto (full version for admin)
func (h *RegionGrpcHandler) toProtoRegion(region *qh_domain.QHRegion) *tqdpb.RegionResponse {
	if region == nil {
		return nil
	}

	return &tqdpb.RegionResponse{
		Id:          region.ID,
		LayerId:     region.LayerID,
		Name:        region.Name,
		DisplayName: region.DisplayName,
		Description: region.Description,
		LandUseCode: region.LandUseCode,
		LandUseName: region.LandUseName,
		AreaHa:      region.AreaHa,
		ColorRed:    uint32(region.ColorRed),
		ColorGreen:  uint32(region.ColorGreen),
		ColorBlue:   uint32(region.ColorBlue),
		ColorHex:    region.CalculateColorHex(),
		Status:      uint32(region.Status),
		Version:     uint32(region.Version),
		IsLatest:    region.IsLatest,
		CreatedAt:   _utils.FormatTimeToString(&region.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(&region.UpdatedAt),
	}
}

// toProtoRegionClient - Convert domain to proto (simplified version for client)
func (h *RegionGrpcHandler) toProtoRegionClient(region *qh_domain.QHRegion) *tqdpb.RegionResponse {
	if region == nil {
		return nil
	}

	resp := h.toProtoRegion(region)

	// Add style config for client
	styleConfig := region.GetStyleConfig()
	styleJSON, _ := json.Marshal(styleConfig)
	resp.StyleConfig = string(styleJSON)

	// Client doesn't need status and version
	resp.Status = 0
	resp.Version = 0

	return resp
}
