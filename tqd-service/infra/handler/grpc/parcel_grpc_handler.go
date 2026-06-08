package handler_grpc

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	_provider "common/provider"
	_utils "common/utils"
	"context"
	"log"
	"time"

	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
	"tqd/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ParcelGrpcHandler struct {
	tqdpb.UnimplementedParcelServiceServer
	usecase      usecase.ParcelUsecase
	mapper       *mapper.ParcelMapper
	SyncProvider *_provider.SyncProvider
}

func NewParcelGrpcHandler(
	usecase usecase.ParcelUsecase,
	mapper *mapper.ParcelMapper,
	syncProvider *_provider.SyncProvider,
) *ParcelGrpcHandler {
	return &ParcelGrpcHandler{
		usecase:      usecase,
		mapper:       mapper,
		SyncProvider: syncProvider,
	}
}

func regionInfoToProto(v *dto.RegionInfoResponse) *tqdpb.RegionInfoResponse {
	if v == nil {
		return nil
	}

	return &tqdpb.RegionInfoResponse{
		RegionId:  v.RegionID,
		LayerId:   v.LayerID,
		LabelId:   v.LabelID,
		LandUseId: v.LandUseID,
		LegendId:  v.LegendID,

		Name:        v.Name,
		DisplayName: v.DisplayName,
		Description: v.Description,

		LayerName:        v.LayerName,
		LayerDisplayName: v.LayerDisplayName,

		LandUseCode:  v.LandUseCode,
		LandUseName:  v.LandUseName,
		LandUseColor: v.LandUseColor,
		CanBuild:     v.CanBuild,

		LabelName:  v.LabelName,
		LabelColor: v.LabelColor,

		LegendColor:  v.LegendColor,
		LegendType:   v.LegendType,
		GeometryType: v.GeometryType,

		LegalDoc:     v.LegalDoc,
		PlanningName: v.PlanningName,

		AreaSqm:    v.AreaSqm,
		AreaHa:     v.AreaHa,
		PerimeterM: v.PerimeterM,

		CenterLat: v.CenterLat,
		CenterLon: v.CenterLon,

		Bounds: &tqdpb.RegionBounds{
			MinLon: v.MinLon,
			MinLat: v.MinLat,
			MaxLon: v.MaxLon,
			MaxLat: v.MaxLat,
		},

		GeoJson: v.GeoJSON,

		Province:     v.Province,
		ProvinceCode: v.ProvinceCode,
		WardCode:     v.WardCode,
	}
}

// =====================================================
// PARCEL METADATA — thông tin cơ bản thửa đất (Free)
// =====================================================

func (h *ParcelGrpcHandler) GetParcelInfo(ctx context.Context, req *tqdpb.GetParcelDetailRequest) (*tqdpb.ParcelInfoResponse, error) {
	if req == nil || req.ParcelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcel_id is required")
	}

	key := h.SyncProvider.GetKeyWithoutMe(ctx, _provider.SyncKeyTQDParcelDetail, req.ParcelId)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ParcelInfoResponse{}, nil
	}

	data, err := h.usecase.GetParcelInfo(ctx, req.ParcelId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if data == nil {
		return nil, status.Error(codes.NotFound, "parcel not found")
	}

	return h.mapper.ToProtoParcelInfoDetail(data), nil
}

// =====================================================
// OVERVIEW — Tổng quan quy hoạch (Free)
// Layers, regions, geometry, basic risk, conclude_label
// =====================================================

func (h *ParcelGrpcHandler) GetParcelOverview(ctx context.Context, req *tqdpb.GetParcelDetailRequest) (*tqdpb.ParcelOverviewResponse, error) {
	if req == nil || req.ParcelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcel_id is required")
	}

	key := h.SyncProvider.GetKeyWithoutMe(ctx, _provider.SyncKeyTQDParcelLayers, req.ParcelId)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ParcelOverviewResponse{}, nil
	}

	data, err := h.usecase.GetParcelOverview(ctx, req.ParcelId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if data == nil {
		return nil, status.Error(codes.NotFound, "parcel not found")
	}

	return h.mapper.ToProtoParcelOverview(data), nil
}

// =====================================================
// ASSESSMENT — Đánh giá chuyên sâu (VIP)
// Conflicts, landUseGroups, assessment, 3 modes
// =====================================================

func (h *ParcelGrpcHandler) GetParcelAssessment(ctx context.Context, req *tqdpb.GetParcelAssessmentRequest) (*tqdpb.ParcelAssessmentResponse, error) {
	if req == nil || req.ParcelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcel_id is required")
	}

	dtoReq := &dto.GetParcelAssessmentRequest{
		ParcelID: req.ParcelId,
		Mode:     req.Mode,
		Preset:   req.Preset,
		AsOfDate: req.AsOfDate,
	}

	result, err := h.usecase.GetParcelAssessment(ctx, dtoReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.mapper.ResolveResultToProtoAssessment(result, req.ParcelId, req.Mode), nil
}

// GetParcelTimeline - Lịch sử biến động quy hoạch của thửa đất
func (h *ParcelGrpcHandler) GetParcelTimeline(ctx context.Context, req *tqdpb.GetParcelDetailRequest) (*tqdpb.ParcelTimelineResponse, error) {
	if req == nil || req.ParcelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcel_id is required")
	}

	timeline, err := h.usecase.GetParcelTimeline(ctx, req.ParcelId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	events := make([]*tqdpb.TimelineEvent, 0, len(timeline.Events))
	for _, evt := range timeline.Events {
		events = append(events, &tqdpb.TimelineEvent{
			Id:          evt.ID,
			Date:        evt.Date.Format("2006-01-02"),
			EventType:   string(evt.EventType),
			LayerId:     evt.LayerID,
			LayerName:   evt.LayerName,
			Description: evt.Description,
		})
	}

	var summary *tqdpb.TimelineSummary
	if timeline.Summary != nil {
		summary = &tqdpb.TimelineSummary{
			TotalEvents:  int32(timeline.Summary.TotalEvents),
			MajorChanges: timeline.Summary.MajorChanges,
		}
		if !timeline.Summary.FirstEventAt.IsZero() {
			summary.FirstEventAt = timeline.Summary.FirstEventAt.Format("2006-01-02")
		}
		if !timeline.Summary.LastEventAt.IsZero() {
			summary.LastEventAt = timeline.Summary.LastEventAt.Format("2006-01-02")
		}
	}

	return &tqdpb.ParcelTimelineResponse{
		ParcelId: req.ParcelId,
		Events:   events,
		Summary:  summary,
	}, nil
}

// GetLayerHistory - Lịch sử phiên bản của 1 layer
func (h *ParcelGrpcHandler) GetLayerHistory(ctx context.Context, req *tqdpb.GetLayerHistoryRequest) (*tqdpb.LayerHistoryResponse, error) {
	if req == nil || req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layer_id is required")
	}

	lineage, err := h.usecase.GetLayerHistory(ctx, req.LayerId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	versions := make([]*tqdpb.LayerVersionInfo, 0, len(lineage.Versions))
	for _, v := range lineage.Versions {
		vi := &tqdpb.LayerVersionInfo{
			Version:   int32(v.Version),
			LayerId:   v.LayerID,
			State:     string(v.State),
			CreatedAt: v.CreatedAt.Format("2006-01-02"),
			CreatedBy: v.CreatedBy,
			IsLatest:  v.IsLatest,
		}
		for _, t := range v.Transitions {
			vi.Transitions = append(vi.Transitions, &tqdpb.LifecycleTransitionInfo{
				FromState: string(t.From),
				ToState:   string(t.To),
				Timestamp: t.Timestamp.Format("2006-01-02"),
				Actor:     t.Actor,
				Reason:    t.Reason,
				LegalDoc:  t.LegalDoc,
			})
		}
		versions = append(versions, vi)
	}

	return &tqdpb.LayerHistoryResponse{
		LayerId:   req.LayerId,
		FamilyId:  lineage.FamilyID,
		Versions:  versions,
		CurrentId: lineage.CurrentID,
	}, nil
}

// =====================================================
// Batch Operations
// =====================================================

// CreateBatchParcels - Tạo batch parcels
func (h *ParcelGrpcHandler) CreateBatchParcels(ctx context.Context, req *tqdpb.ParcelListDTO) (*tqdpb.ParcelListDTO, error) {
	// Get user ID from context
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	// Validate request
	if req == nil || len(req.Data) == 0 {
		return nil, status.Error(codes.InvalidArgument, "request data is required")
	}

	log.Printf("[CreateBatchParcels] UserID: %d, Batch size: %d", userID, len(req.Data))

	// Convert proto to request DTO
	createReq := h.mapper.ProtoToDomains(req)

	// Call usecase
	result, err := h.usecase.CreateBatch(ctx, &dto.ParcelList{
		Data: createReq,
	})
	if err != nil {
		log.Printf("[CreateBatchParcels] Error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create parcel: %v", err)
	}

	log.Printf("[CreateBatchParcels] Successfully created %d parcels", len(result.Data))

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _provider.SyncKeyTQDParcelDetail, 0), t.UnixMilli())

	return &tqdpb.ParcelListDTO{
		Data: h.mapper.ToProtoFromDomainList(result),
	}, nil
}

// =====================================================
// Location-based Queries
// =====================================================

// FindParcelByLocation - Tìm parcel theo tọa độ
func (h *ParcelGrpcHandler) FindParcelByLocation(ctx context.Context, req *tqdpb.FindParcelByLocationRequest) (*tqdpb.Parcel, error) {
	key := h.SyncProvider.GetKey(ctx, _provider.SyncKeyTQDParcelDetail, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.Parcel{}, nil
	}

	// Validate input
	if req.Latitude == 0 && req.Longitude == 0 {
		return nil, status.Error(codes.InvalidArgument, "latitude and longitude are required")
	}

	// Validate coordinate ranges
	if req.Latitude < -90 || req.Latitude > 90 {
		return nil, status.Error(codes.InvalidArgument, "latitude must be between -90 and 90")
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		return nil, status.Error(codes.InvalidArgument, "longitude must be between -180 and 180")
	}

	log.Printf("[FindParcelByLocation] Lat: %f, Lng: %f", req.Latitude, req.Longitude)

	// Call usecase
	result, err := h.usecase.FindByLocation(ctx, req.Latitude, req.Longitude)
	if err != nil {
		log.Printf("[FindParcelByLocation] Error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to find parcel by location: %v", err)
	}

	if result == nil {
		return nil, status.Error(codes.NotFound, "parcel not found")
	}

	return h.mapper.DomainToProto(result), nil
}

// ResolveMapTargetByLocation - Resolve target theo tọa độ: ưu tiên parcel, fallback region.
func (h *ParcelGrpcHandler) ResolveMapTargetByLocation(
	ctx context.Context,
	req *tqdpb.ResolveMapTargetByLocationRequest,
) (*tqdpb.ResolveMapTargetByLocationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	if req.GetLatitude() == 0 && req.GetLongitude() == 0 {
		return nil, status.Error(codes.InvalidArgument, "latitude and longitude are required")
	}
	if req.GetLatitude() < -90 || req.GetLatitude() > 90 {
		return nil, status.Error(codes.InvalidArgument, "latitude must be between -90 and 90")
	}
	if req.GetLongitude() < -180 || req.GetLongitude() > 180 {
		return nil, status.Error(codes.InvalidArgument, "longitude must be between -180 and 180")
	}

	result, err := h.usecase.ResolveMapTargetByLocation(ctx, dto.ResolveMapTargetRequestDTO{
		Latitude:  req.GetLatitude(),
		Longitude: req.GetLongitude(),
		Radius:    req.GetRadius(),
	})
	if err != nil {
		if err.Error() == "map target not found" {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &tqdpb.ResolveMapTargetByLocationResponse{
		TargetType: result.TargetType,
		TargetId:   result.TargetID,
	}

	if result.Parcel != nil {
		resp.Parcel = h.mapper.ToProtoParcelInfoDetail(result.Parcel)
	}
	if result.Region != nil {
		resp.Region = regionInfoToProto(result.Region)
	}

	return resp, nil
}

// =====================================================
// Polygon-based Queries
// =====================================================

// FindParcelsByPolygon - Tìm parcels theo polygon
func (h *ParcelGrpcHandler) FindParcelsByPolygon(
	ctx context.Context,
	req *tqdpb.FindParcelsByPolygonRequest,
) (*tqdpb.FindParcelsByPolygonResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _provider.SyncKeyTQDParcelDetail, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.FindParcelsByPolygonResponse{}, nil
	}

	// Validate coordinates
	var points []repo.Point
	if len(req.Coordinates) == 0 {
		return nil, status.Error(codes.InvalidArgument, "coordinates are required")
	}

	// Build points list
	points = make([]repo.Point, len(req.Coordinates))
	for i, pt := range req.Coordinates {
		// Validate each point
		if pt.Lat < -90 || pt.Lat > 90 {
			return nil, status.Error(codes.InvalidArgument, "latitude must be between -90 and 90")
		}
		if pt.Lng < -180 || pt.Lng > 180 {
			return nil, status.Error(codes.InvalidArgument, "longitude must be between -180 and 180")
		}
		points[i] = repo.Point{Lat: pt.Lat, Lng: pt.Lng}
	}

	// Validate polygon has at least 3 points
	if len(points) < 3 {
		return nil, status.Error(codes.InvalidArgument, "polygon must have at least 3 points")
	}

	// Build pagable
	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)

	log.Printf("[FindParcelsByPolygon] Points count: %d, Intersect: %v, MaxAreaKm2: %f, Page: %d, Size: %d",
		len(points), req.Intersect, req.MaxAreaKm2, pagable.GetPage(), pagable.GetSize())

	page := int(pagable.GetPage())
	pageSize := int(pagable.GetSize())

	parcels, total, err := h.usecase.FindParcelsByPolygon(
		ctx,
		points,
		req.Intersect,
		req.MaxAreaKm2,
		page,
		pageSize,
	)
	if err != nil {
		log.Printf("[FindParcelsByPolygon] Usecase error: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Printf("[FindParcelsByPolygon] Found %d parcels, Total: %d", len(parcels), total)

	// Convert domain to proto
	data := make([]*tqdpb.Parcel, len(parcels))
	for i, p := range parcels {
		data[i] = h.mapper.DomainToProto(&p)
	}

	return &tqdpb.FindParcelsByPolygonResponse{
		Data:     data,
		Total:    int32(total),
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}

func (h *ParcelGrpcHandler) SearchTxtClientParcels(
	ctx context.Context,
	req *sharepb.RequestV3Proto,
) (*tqdpb.ListParcelResponse, error) {
	// Infer map number and land number from search text
	searchInfo := _utils.InferParcelSearch(req.Text)
	if searchInfo.CleanedText == "" && searchInfo.MapNumber == "" && searchInfo.LandNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "search text is required")
	}

	userID := _utils.GetOriginIdFromContext(ctx)
	log.Printf("[SearchTxtClientParcels] UserID: %d, OriginalText: %s, CleanedText: %s, Map: %s, Land: %s, Page: %d, Size: %d",
		userID, req.Text, searchInfo.CleanedText, searchInfo.MapNumber, searchInfo.LandNumber, req.Page, req.Size)

	pagable := _dto.Pagable{
		Page: req.Page,
		Size: req.Size,
	}
	pagable.Normalize()

	// Gọi SearchParcelsByText với các thông tin đã bóc tách
	parcels, total, err := h.usecase.SearchParcelsByText(ctx, searchInfo.CleanedText, searchInfo.MapNumber, searchInfo.LandNumber, int(pagable.GetPage()), int(pagable.GetSize()))
	if err != nil {
		log.Printf("[SearchTxtClientParcels] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert kết quả sang proto
	results := make([]*tqdpb.ParcelInfo, 0, len(parcels))
	for _, p := range parcels {
		results = append(results, h.mapper.ToProtoParcelInfo(p))
	}

	return &tqdpb.ListParcelResponse{
		Data:  results,
		Total: uint32(total),
	}, nil
}

// PreviewParcelsById - Xem trước thông tin parcel theo ID
func (h *ParcelGrpcHandler) PreviewParcelsById(
	ctx context.Context,
	req *sharepb.RequestV3Proto,
) (*tqdpb.ParcelInfo, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	userID := _utils.GetOriginIdFromContext(ctx)
	log.Printf("[PreviewParcelsById] UserID: %d, ParcelID: %d", userID, req.Id)

	// Gọi đúng method GetParcelInfoByID (không phải GetParcelDetail)
	info, err := h.usecase.GetParcelInfoByID(ctx, req.Id)
	if err != nil {
		log.Printf("[PreviewParcelsById] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	if info == nil {
		return nil, status.Error(codes.NotFound, "parcel not found")
	}

	// Convert sang proto
	return h.mapper.ToProtoParcelInfo(info), nil
}

func (h *ParcelGrpcHandler) PublicParcelsById(
	ctx context.Context,
	req *sharepb.RequestV3Proto,
) (*tqdpb.ParcelInfo, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	userID := _utils.GetOriginIdFromContext(ctx)
	log.Printf("[PreviewParcelsById] UserID: %d, ParcelID: %d", userID, req.Id)

	// Gọi đúng method GetParcelInfoByID (không phải GetParcelDetail)
	info, err := h.usecase.GetParcelInfoByID(ctx, req.Id)
	if err != nil {
		log.Printf("[PreviewParcelsById] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	if info == nil {
		return nil, status.Error(codes.NotFound, "parcel not found")
	}

	// Convert sang proto
	return h.mapper.ToProtoParcelInfo(info), nil
}
