package handler_grpc

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/infra/validator"
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
	SyncProvider *_utils.SyncUtil
}

var _ tqdpb.ParcelServiceServer = (*ParcelGrpcHandler)(nil)

func NewParcelGrpcHandler(
	usecase usecase.ParcelUsecase,
	mapper *mapper.ParcelMapper,
	syncProvider *_utils.SyncUtil,
) *ParcelGrpcHandler {
	return &ParcelGrpcHandler{
		usecase:      usecase,
		mapper:       mapper,
		SyncProvider: syncProvider,
	}
}

// =====================================================
// PARCEL METADATA — thông tin cơ bản thửa đất (Free)
// =====================================================

func (h *ParcelGrpcHandler) GetParcelInfo(ctx context.Context, req *tqdpb.GetParcelDetailRequest) (*tqdpb.ParcelInfoResponse, error) {
	if req == nil || req.ParcelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcel_id is required")
	}

	key := h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDParcelDetail, req.ParcelId)
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

// @Summary Thửa đất quanh thửa hiện tại
// @Description Trả khoảng 5 thửa gần nhất (ưu tiên tiếp giáp)
// @Tags Parcel
// @Accept json
// @Produce json
// @Param parcelId path int true "Parcel ID"
// @Param limit query int false "Limit (default 5)"
// @Success 200 {object} tqdpb.GetParcelsAroundResponse
// @Router /v2/tqd/parcels/around/{parcelId} [get]
func (h *ParcelGrpcHandler) GetParcelsAround(ctx context.Context, req *tqdpb.GetParcelsAroundRequest) (*tqdpb.GetParcelsAroundResponse, error) {
	if req == nil || req.ParcelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcel_id is required")
	}

	limit := int(req.GetLimit())
	if limit <= 0 {
		limit = 5
	}

	rows, err := h.usecase.GetParcelsAround(ctx, req.ParcelId, limit)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.GetParcelsAroundResponse{
		Data:  h.mapper.ToProtoNearbyParcelItems(rows),
		Total: uint32(len(rows)),
	}, nil
}

func responseMetaToProto(meta *dto.ResponseMeta) *tqdpb.ResponseMetadataInfo {
	if meta == nil {
		return nil
	}

	return &tqdpb.ResponseMetadataInfo{
		ResponseTimeMs: meta.ResponseTimeMs,
		DataVersion:    meta.DataSource,
		IsCached:       meta.Cached,
	}
}

func (h *ParcelGrpcHandler) GetParcelPlanning(ctx context.Context, req *tqdpb.GetParcelPlanningRequest) (*tqdpb.ParcelPlanningResponse, error) {
	if req == nil || req.ParcelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcel_id is required")
	}

	data, err := h.usecase.GetParcelPlanning(ctx, &dto.GetParcelPlanningRequest{
		ParcelID:          req.ParcelId,
		View:              req.View,
		Mode:              req.Mode,
		Preset:            req.Preset,
		AsOfDate:          req.AsOfDate,
		IncludeOverview:   req.IncludeOverview,
		IncludeAssessment: req.IncludeAssessment,
		Timestamp:         req.Timestamp,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if data == nil {
		return nil, status.Error(codes.NotFound, "parcel planning not found")
	}

	resp := &tqdpb.ParcelPlanningResponse{
		ParcelId: data.ParcelID,
		View:     data.View,
		Metadata: responseMetaToProto(data.Metadata),
	}
	if data.Overview != nil {
		resp.Overview = h.mapper.ToProtoParcelOverview(data.Overview)
	}
	if data.Assessment != nil {
		resp.Assessment = h.mapper.ResolveResultToProtoAssessment(data.Assessment, data.ParcelID, data.Mode)
	}

	return resp, nil
}

// @Summary Đánh giá quy hoạch thửa đất (tạm thời)
// @Description Trả mức cần kiểm tra, đánh giá chính và gợi ý kiểm tra
// @Tags Parcel
// @Accept json
// @Produce json
// @Param parcelId path int true "Parcel ID"
// @Success 200 {object} tqdpb.ParcelPlanningEvalResponse
// @Router /v2/tqd/parcels/{parcelId}/planning-eval [get]
func (h *ParcelGrpcHandler) GetParcelPlanningEval(
	ctx context.Context,
	req *tqdpb.GetParcelDetailRequest,
) (*tqdpb.ParcelPlanningEvalResponse, error) {
	if req == nil || req.ParcelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcel_id is required")
	}

	// TODO: thay bằng logic đánh giá thật — hiện trả dữ liệu mẫu tạm thời
	return &tqdpb.ParcelPlanningEvalResponse{
		ParcelId:   req.ParcelId,
		CheckLevel: tqdpb.ParcelCheckLevel_PARCEL_CHECK_LEVEL_HIGH,
		MainAssessment: []string{
			"Thửa đất nằm trong vùng quy hoạch sử dụng đất giao thông (đường vành đai).",
			"Có lớp QH phân khu H2-4 chồng lấn, cần đối chiếu chi tiết.",
			"Diện tích đất ở còn lại sau quy hoạch ước tính khoảng 85 m².",
		},
		CheckSuggestions: []string{
			"Đối chiếu bản đồ quy hoạch giao thông (đường vành đai) với ranh thửa.",
			"Kiểm tra mức độ chồng lấn với quy hoạch phân khu H2-4.",
			"Xác nhận diện tích đất ở còn lại và hạn chế xây dựng liên quan.",
		},
	}, nil
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

// GetParcelHistory - Lịch sử thửa đất (shape giống Đồ án quy hoạch / ListPlanningEvents)
// @Summary Lấy lịch sử thửa đất
// @Description Trả danh sách sự kiện lịch sử của thửa theo parcelId (phân trang)
// @Tags Parcel
// @Accept json
// @Produce json
// @Param parcelId path int true "Parcel ID"
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Success 200 {object} tqdpb.GetParcelHistoryResponse
// @Router /v2/tqd/parcels/{parcelId}/history [get]
func (h *ParcelGrpcHandler) GetParcelHistory(ctx context.Context, req *tqdpb.GetParcelHistoryRequest) (*tqdpb.GetParcelHistoryResponse, error) {
	if req == nil || req.ParcelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcel_id is required")
	}

	pagable := _dto.Pagable{
		Page: req.Page,
		Size: req.Size,
	}
	pagable.Normalize()

	rows, total, err := h.usecase.GetParcelHistory(ctx, req.ParcelId, pagable)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.GetParcelHistoryResponse{
		Data:  h.mapper.ToProtoParcelHistoryItems(rows),
		Total: total,
		Page:  pagable.GetPage(),
		Size:  uint32(pagable.GetLimit()),
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
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDParcelDetail, 0), t.UnixMilli())

	return &tqdpb.ParcelListDTO{
		Data: h.mapper.ToProtoFromDomainList(result),
	}, nil
}

// =====================================================
// Location-based Queries
// =====================================================

// FindParcelByLocation - Tìm parcel theo tọa độ
func (h *ParcelGrpcHandler) FindParcelByLocation(ctx context.Context, req *tqdpb.FindParcelByLocationRequest) (*tqdpb.Parcel, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDParcelDetail, 0)
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
// Lưu ý nghiệp vụ:
// - Click map giữ logic parcel-first.
// - Không phụ thuộc zoom.
// - Zoom chỉ dùng cho polygon draw selection.
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
		resp.Region = h.mapper.RegionInfoToProto(result.Region)
	}

	return resp, nil
}

// ResolveMapTargetByPoi - Resolve target theo POI: lấy tọa độ POI rồi ưu tiên parcel, fallback region.
// @Summary Resolve map target by POI
// @Description Lấy chi tiết thửa đất hoặc region liên quan tới POI (cùng contract với by-location)
// @Tags Parcel
// @Produce json
// @Param poi_id path uint64 true "POI ID"
// @Success 200 {object} tqdpb.ResolveMapTargetByLocationResponse
// @Router /v2/tqd/map-target/by-poi/{poi_id} [get]
// @bind: internal/usecase.ParcelUsecase
func (h *ParcelGrpcHandler) ResolveMapTargetByPoi(
	ctx context.Context,
	req *tqdpb.ResolveMapTargetByPoiRequest,
) (*tqdpb.ResolveMapTargetByLocationResponse, error) {
	if req == nil || req.GetPoiId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "poi_id is required")
	}

	result, err := h.usecase.ResolveMapTargetByPoi(ctx, req.GetPoiId())
	if err != nil {
		msg := err.Error()
		if msg == "map target not found" ||
			msg == "poi not found or missing coordinates" ||
			msg == "poi id is required" {
			return nil, status.Error(codes.NotFound, msg)
		}
		return nil, status.Error(codes.Internal, msg)
	}

	resp := &tqdpb.ResolveMapTargetByLocationResponse{
		TargetType: result.TargetType,
		TargetId:   result.TargetID,
	}
	if result.Parcel != nil {
		resp.Parcel = h.mapper.ToProtoParcelInfoDetail(result.Parcel)
	}
	if result.Region != nil {
		resp.Region = h.mapper.RegionInfoToProto(result.Region)
	}
	return resp, nil
}

func getOptionalPolygonZoom(req *tqdpb.FindParcelsByPolygonRequest) *uint32 {
	if req == nil {
		return nil
	}

	// Với proto3 optional uint32 z = 8, Go generated thường sinh field Z *uint32.
	// Nếu generated code của bạn chưa có req.Z, cần generate lại protobuf sau khi sửa proto.
	if req.Z == nil {
		return nil
	}

	z := req.GetZ()
	return &z
}

func uniqueUint64List(values []uint64) []uint64 {
	if len(values) == 0 {
		return nil
	}

	seen := make(map[uint64]struct{}, len(values))
	result := make([]uint64, 0, len(values))

	for _, value := range values {
		if value == 0 {
			continue
		}

		if _, ok := seen[value]; ok {
			continue
		}

		seen[value] = struct{}{}
		result = append(result, value)
	}

	return result
}

func uniqueStringList(values []string) []string {
	if len(values) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))

	for _, value := range values {
		normalized := strings.ToLower(strings.TrimSpace(value))
		if normalized == "" {
			continue
		}

		if _, ok := seen[normalized]; ok {
			continue
		}

		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}

	return result
}

func getPolygonTargetFilter(req *tqdpb.FindParcelsByPolygonRequest) *dto.PolygonTargetFilter {
	if req == nil {
		return nil
	}

	filter := &dto.PolygonTargetFilter{
		LayerIDs:     uniqueUint64List(req.LayerIds),
		LabelIDs:     uniqueUint64List(req.LabelIds),
		LandUseIDs:   uniqueUint64List(req.LandUseIds),
		LandUseCodes: uniqueStringList(req.LandUseCodes),
		Keyword:      strings.TrimSpace(req.Keyword),
	}

	if req.CanBuild != nil {
		canBuild := req.GetCanBuild()
		filter.CanBuild = &canBuild
	}

	if !filter.HasPlanningFilters() {
		return nil
	}

	return filter
}

func validatePolygonCoordinates(req *tqdpb.FindParcelsByPolygonRequest) ([]repo.Point, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	if len(req.Coordinates) == 0 {
		return nil, status.Error(codes.InvalidArgument, "coordinates are required")
	}

	points := make([]repo.Point, 0, len(req.Coordinates))

	for _, pt := range req.Coordinates {
		if pt == nil {
			return nil, status.Error(codes.InvalidArgument, "coordinate item is required")
		}

		if pt.Lat < -90 || pt.Lat > 90 {
			return nil, status.Error(codes.InvalidArgument, "latitude must be between -90 and 90")
		}

		if pt.Lng < -180 || pt.Lng > 180 {
			return nil, status.Error(codes.InvalidArgument, "longitude must be between -180 and 180")
		}

		points = append(points, repo.Point{
			Lat: pt.Lat,
			Lng: pt.Lng,
		})
	}

	if len(points) < 3 {
		return nil, status.Error(codes.InvalidArgument, "polygon must have at least 3 points")
	}

	return points, nil
}

// FindParcelsByPolygon - Resolve targets theo polygon.
//
// Legacy compatibility:
//   - Nếu client không truyền z:
//     giữ hành vi cũ, trả parcel.
//
// New behavior:
//   - Nếu client truyền z và z >= ParcelPolygonMinZ (hiện tại 15):
//     trả parcel.
//   - Nếu client truyền z và z < ParcelPolygonMinZ (14 trở xuống):
//     trả region theo layer/label zoom.
//
// Lưu ý:
// - Click location không đi qua function này.
// - Click location vẫn giữ parcel-first fallback region.
func (h *ParcelGrpcHandler) FindParcelsByPolygon(
	ctx context.Context,
	req *tqdpb.FindParcelsByPolygonRequest,
) (*tqdpb.FindParcelsByPolygonResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDParcelDetail, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)

	if !updated {
		return &tqdpb.FindParcelsByPolygonResponse{}, nil
	}

	points, err := validatePolygonCoordinates(req)
	if err != nil {
		return nil, err
	}

	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)

	page := int(pagable.GetPage())
	pageSize := int(pagable.GetSize())

	z := getOptionalPolygonZoom(req)
	filter := getPolygonTargetFilter(req)
	includeAnalysis := req.GetIncludeAnalysis()

	if z == nil {
		log.Printf(
			"[FindParcelsByPolygon] Points=%d, Intersect=%v, MaxAreaKm2=%f, Page=%d, Size=%d, Z=<nil>, Mode=legacy-parcel, HasFilter=%v, IncludeAnalysis=%v",
			len(points),
			req.Intersect,
			req.MaxAreaHa,
			page,
			pageSize,
			filter != nil,
			includeAnalysis,
		)
	} else {
		log.Printf(
			"[FindParcelsByPolygon] Points=%d, Intersect=%v, MaxAreaKm2=%f, Page=%d, Size=%d, Z=%d, HasFilter=%v, IncludeAnalysis=%v",
			len(points),
			req.Intersect,
			req.MaxAreaHa,
			page,
			pageSize,
			*z,
			filter != nil,
			includeAnalysis,
		)
	}

	result, err := h.usecase.FindMapTargetsByPolygon(
		ctx,
		points,
		req.Intersect,
		req.MaxAreaHa,
		page,
		pageSize,
		z,
		filter,
		includeAnalysis,
	)
	if err != nil {
		log.Printf("[FindParcelsByPolygon] Usecase error: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if result == nil {
		return &tqdpb.FindParcelsByPolygonResponse{
			Parcels:    []*tqdpb.Parcel{},
			Regions:    []*tqdpb.RegionInfoResponse{},
			Total:      0,
			Page:       int32(pagable.GetPage()),
			PageSize:   int32(pagable.GetSize()),
			ResultType: mapper.WorkspaceEntityTypeParcel,
		}, nil
	}

	log.Printf(
		"[FindParcelsByPolygon] ResultType=%d, Z=%d, Parcels=%d, Regions=%d, Total=%d",
		result.ResultType,
		result.Z,
		len(result.Parcels),
		len(result.Regions),
		result.Total,
	)

	return h.mapper.PolygonMapTargetsToProto(result), nil
}

func (h *ParcelGrpcHandler) CheckingPolygon(
	ctx context.Context,
	req *tqdpb.FindParcelsByPolygonRequest,
) (*tqdpb.FindParcelsByPolygonResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	// tối đa 1km2
	req.MaxAreaHa = 100

	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDParcelDetail, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)

	if !updated {
		return &tqdpb.FindParcelsByPolygonResponse{}, nil
	}

	points, err := validatePolygonCoordinates(req)
	if err != nil {
		return nil, err
	}

	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)

	page := int(pagable.GetPage())
	pageSize := int(pagable.GetSize())

	z := getOptionalPolygonZoom(req)
	filter := getPolygonTargetFilter(req)
	includeAnalysis := req.GetIncludeAnalysis()

	log.Printf(
		"[CheckingPolygon] Points=%d, Intersect=%v, MaxAreaKm2=%f, Page=%d, Size=%d, Z=%v, HasFilter=%v, IncludeAnalysis=%v",
		len(points),
		req.Intersect,
		req.MaxAreaHa,
		page,
		pageSize,
		z,
		filter != nil,
		includeAnalysis,
	)

	handlerStart := time.Now()
	result, err := h.usecase.CheckingPolygon(
		ctx,
		points,
		req.Intersect,
		req.MaxAreaHa,
		page,
		pageSize,
		z,
		filter,
		includeAnalysis,
	)
	if err != nil {
		log.Printf("[CheckingPolygon] Usecase error: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if result == nil {
		return &tqdpb.FindParcelsByPolygonResponse{
			Parcels:    []*tqdpb.Parcel{},
			Regions:    []*tqdpb.RegionInfoResponse{},
			Total:      0,
			Page:       int32(pagable.GetPage()),
			PageSize:   int32(pagable.GetSize()),
			ResultType: mapper.WorkspaceEntityTypeParcel,
		}, nil
	}

	log.Printf(
		"[CheckingPolygon] ResultType=%d, Z=%d, Parcels=%d, Regions=%d, Total=%d, handlerTotal=%v",
		result.ResultType,
		result.Z,
		len(result.Parcels),
		len(result.Regions),
		result.Total,
		time.Since(handlerStart),
	)

	return h.mapper.PolygonMapTargetsToProto(result), nil
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

func (h *ParcelGrpcHandler) SearchPublicParcels(
	ctx context.Context,
	req *sharepb.RequestV3Proto,
) (*tqdpb.ListParcelResponse, error) {
	log.Printf("[SearchPublicParcels] Page: %d, Size: %d", req.Page, req.Size)

	pagable := _dto.Pagable{
		Page: req.Page,
		Size: req.Size,
	}
	pagable.Normalize()

	parcels, total, err := h.usecase.SearchPublicParcels(ctx, int(pagable.GetPage()), int(pagable.GetSize()))
	if err != nil {
		log.Printf("[SearchPublicParcels] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

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

func (h *ParcelGrpcHandler) GetParcelSeoSource(ctx context.Context, req *tqdpb.GetParcelDetailRequest) (*tqdpb.ParcelSeoSourceResponse, error) {
	if req.ParcelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcelId is required")
	}

	source, err := h.usecase.GetParcelSeoSource(ctx, req.ParcelId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if source == nil {
		return nil, status.Error(codes.NotFound, "parcel not found")
	}

	return parcelSeoSourceToProto(source), nil
}

func (h *ParcelGrpcHandler) ListParcelSeoSourcesForGenerate(ctx context.Context, req *tqdpb.ListParcelSeoSourcesRequest) (*tqdpb.ListParcelSeoSourcesResponse, error) {
	sources, err := h.usecase.ListParcelSeoSourcesForGenerate(ctx, req.GetLimit())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	items := make([]*tqdpb.ParcelSeoSourceResponse, 0, len(sources))
	for i := range sources {
		items = append(items, parcelSeoSourceToProto(&sources[i]))
	}
	return &tqdpb.ListParcelSeoSourcesResponse{Data: items}, nil
}

func (h *ParcelGrpcHandler) UpdateParcelSeoID(ctx context.Context, req *tqdpb.UpdateParcelSeoIDRequest) (*sharepb.SubmitResponse, error) {
	if req.ParcelId == 0 || req.SeoId == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcelId và seoId là bắt buộc")
	}

	if err := h.usecase.UpdateParcelSeoID(ctx, req.ParcelId, req.SeoId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &sharepb.SubmitResponse{
		Id:      req.ParcelId,
		Message: "parcel seo_id updated successfully",
	}, nil
}

// AdminGenParcelSeo tạo seo_domain trên CRM cho danh sách parcel.
// @Summary Admin tạo SEO domain từ parcel
// @Description Nhận danh sách parcelIds, gọi CRM sinh bản ghi seo_domain tương ứng.
// @Tags Parcel Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body tqdpb.AdminGenParcelSeoRequest true "Danh sách parcelIds"
// @Router /v2/tqd/admin/parcels/gen-seo [post]
func (h *ParcelGrpcHandler) AdminGenParcelSeo(ctx context.Context, req *tqdpb.AdminGenParcelSeoRequest) (*tqdpb.AdminGenParcelSeoResponse, error) {
	if req == nil || len(req.ParcelIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcelIds is required")
	}

	log.Printf("[AdminGenParcelSeo] parcelIds=%v", req.ParcelIds)

	result, err := h.usecase.AdminGenParcelSeo(ctx, req.ParcelIds)
	if err != nil {
		log.Printf("[AdminGenParcelSeo] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	if result == nil {
		return &tqdpb.AdminGenParcelSeoResponse{}, nil
	}

	items := make([]*tqdpb.AdminGenParcelSeoItem, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, &tqdpb.AdminGenParcelSeoItem{
			ParcelId:     item.ParcelID,
			SeoId:        item.SeoID,
			CanonicalUrl: item.CanonicalURL,
			Slug:         item.Slug,
			Created:      item.Created,
			Message:      item.Message,
		})
	}

	return &tqdpb.AdminGenParcelSeoResponse{
		Data:    items,
		Total:   result.Total,
		Created: result.Created,
		Skipped: result.Skipped,
		Failed:  result.Failed,
		Errors:  result.Errors,
	}, nil
}

// AdminListParcels — danh sách thửa đất cho admin (phân trang + followerCount).
// @Summary Admin danh sách thửa đất
// @Description Trả danh sách parcels có phân trang, kèm followerCount
// @Tags Parcel Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Param search query string false "Search keyword"
// @Success 200 {object} tqdpb.AdminListParcelsResponse
// @Router /v2/tqd/admin/parcels [get]
func (h *ParcelGrpcHandler) AdminListParcels(ctx context.Context, req *tqdpb.AdminListParcelsRequest) (*tqdpb.AdminListParcelsResponse, error) {
	pagable := _dto.Pagable{
		Page: req.GetPage(),
		Size: req.GetSize(),
	}
	pagable.Normalize()

	rows, total, err := h.usecase.AdminListParcels(ctx, req.GetSearch(), int(pagable.GetPage()), int(pagable.GetSize()))
	if err != nil {
		log.Printf("[AdminListParcels] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.AdminListParcelsResponse{
		Data:  h.mapper.ToProtoAdminParcelListItems(rows),
		Total: total,
	}, nil
}

// AdminGetParcel — chi tiết thửa đất cho admin.
// @Summary Admin chi tiết thửa đất
// @Description Trả thông tin chi tiết thửa đất cho màn admin
// @Tags Parcel Admin
// @Produce json
// @Security BearerAuth
// @Param parcelId path int true "Parcel ID"
// @Success 200 {object} tqdpb.AdminParcelDetail
// @Router /v2/tqd/admin/parcels/{parcelId} [get]
// @bind: internal/usecase.ParcelUsecase
func (h *ParcelGrpcHandler) AdminGetParcel(ctx context.Context, req *tqdpb.AdminGetParcelRequest) (*tqdpb.AdminParcelDetail, error) {
	if req == nil || req.GetParcelId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcelId is required")
	}

	detail, err := h.usecase.AdminGetParcel(ctx, req.GetParcelId())
	if err != nil {
		log.Printf("[AdminGetParcel] Error: %v", err)
		if strings.Contains(err.Error(), "not found") {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.mapper.ToProtoAdminParcelDetail(detail), nil
}

// AdminGetParcelNeighbors — thửa lân cận theo geometry.
// @Summary Admin thửa lân cận
// @Description Quét geometry quanh thửa hiện tại, trả tối đa 10 thửa (ưu tiên tiếp giáp, rồi gần nhất)
// @Tags Parcel Admin
// @Produce json
// @Security BearerAuth
// @Param parcelId path int true "Parcel ID"
// @Success 200 {object} tqdpb.AdminGetParcelNeighborsResponse
// @Router /v2/tqd/admin/parcels/{parcelId}/neighbors [get]
// @bind: internal/usecase.ParcelUsecase
func (h *ParcelGrpcHandler) AdminGetParcelNeighbors(ctx context.Context, req *tqdpb.AdminGetParcelNeighborsRequest) (*tqdpb.AdminGetParcelNeighborsResponse, error) {
	if req == nil || req.GetParcelId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcelId is required")
	}

	rows, err := h.usecase.AdminGetParcelNeighbors(ctx, req.GetParcelId())
	if err != nil {
		log.Printf("[AdminGetParcelNeighbors] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.AdminGetParcelNeighborsResponse{
		Data:  h.mapper.ToProtoAdminParcelNeighbors(rows),
		Total: int64(len(rows)),
	}, nil
}

// AdminUpdateParcel — sửa thông tin thửa đất cho admin.
// @Summary Admin sửa thửa đất
// @Description Cập nhật thông tin metadata thửa đất (địa chỉ, tờ/thửa, tọa độ, diện tích)
// @Tags Parcel Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param parcelId path int true "Parcel ID"
// @Param body body tqdpb.AdminUpdateParcelRequest true "AdminUpdateParcelRequest"
// @Success 200 {object} tqdpb.AdminParcelDetail
// @Router /v2/tqd/admin/parcels/{parcelId} [put]
// @bind: internal/usecase.ParcelUsecase
func (h *ParcelGrpcHandler) AdminUpdateParcel(ctx context.Context, req *tqdpb.AdminUpdateParcelRequest) (*tqdpb.AdminParcelDetail, error) {
	dtoReq := h.mapper.ToAdminUpdateParcelDTO(req)
	if err := validator.NewParcelValidator().ValidateAdminUpdate(dtoReq); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	detail, err := h.usecase.AdminUpdateParcel(ctx, dtoReq)
	if err != nil {
		log.Printf("[AdminUpdateParcel] Error: %v", err)
		if strings.Contains(err.Error(), "not found") {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.mapper.ToProtoAdminParcelDetail(detail), nil
}

func parcelSeoSourceToProto(source *dto.ParcelSeoSource) *tqdpb.ParcelSeoSourceResponse {
	if source == nil {
		return nil
	}
	var seoID uint64
	if source.SeoID != nil {
		seoID = *source.SeoID
	}
	return &tqdpb.ParcelSeoSourceResponse{
		ParcelId:  source.ParcelID,
		AdrSearch: source.AdrSearch,
		SeoId:     seoID,
	}
}

func (h *ParcelGrpcHandler) GetParcelQuickLayers(
	ctx context.Context,
	req *tqdpb.GetParcelPlanningRequest,
) (*tqdpb.LayerQuickInfoResponse, error) {
	if req.ParcelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcelId is required")
	}

	quickInfo, err := h.usecase.GetParcelQuickInfo(ctx, req.ParcelId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.mapper.BuildLayerQuickInfoFromRows(quickInfo), nil
}

// GetWardQuickLayers — GET /v2/tqd/client/administrative-units/{compositeId}/quick-layers
// @Summary Lấy quick-layers theo geometry đơn vị hành chính
// @Description Tìm giao geometry tỉnh/phường với qh_regions, trả về danh sách layer và land use info.
// @Tags Administrative Unit Client
// @Produce json
// @Param compositeId path int true "Composite ID ({sourceId}{1=tỉnh|2=phường})"
// @Success 200 {object} tqdpb.LayerQuickInfoResponse
// @Router /v2/tqd/client/administrative-units/{compositeId}/quick-layers [get]
func (h *ParcelGrpcHandler) GetWardQuickLayers(
	ctx context.Context,
	req *tqdpb.GetWardQuickLayersRequest,
) (*tqdpb.LayerQuickInfoResponse, error) {
	if req == nil || req.CompositeId == 0 {
		return nil, status.Error(codes.InvalidArgument, "compositeId is required")
	}

	quickInfo, err := h.usecase.GetAdministrativeUnitQuickInfo(ctx, req.CompositeId)
	if err != nil {
		if strings.Contains(err.Error(), "decode compositeId") || strings.Contains(err.Error(), "invalid compositeId kind") {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.mapper.BuildLayerQuickInfoFromRows(quickInfo), nil
}

func (h *ParcelGrpcHandler) GetParcelDetail(ctx context.Context, req *tqdpb.GetParcelDetailRequest) (*tqdpb.ParcelDetailResponse, error) {
	log.Printf("[GetParcelDetail] ParcelID: %d", req.ParcelId)

	if req == nil || req.ParcelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcel_id is required")
	}

	data, err := h.usecase.GetParcelDetail(ctx, req.ParcelId)
	if err != nil {
		log.Printf("[GetParcelDetail] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	if data == nil {
		return nil, status.Error(codes.NotFound, "parcel not found")
	}

	return h.mapper.ToProtoParcelDetail(data), nil
}

func (h *ParcelGrpcHandler) GetParcelLayers(ctx context.Context, req *tqdpb.GetParcelLayersRequest) (*tqdpb.ParcelLayersResponse, error) {
	log.Printf("[GetParcelLayers] ParcelID: %d, Page: %d, Limit: %d", req.ParcelId, req.Page, req.Limit)

	if req == nil || req.ParcelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcel_id is required")
	}

	pagable := _dto.Pagable{
		Page: req.Page,
		Size: req.Limit,
	}
	pagable.Normalize()

	data, total, err := h.usecase.GetParcelLayers(ctx, req.ParcelId, pagable)
	if err != nil {
		log.Printf("[GetParcelLayers] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &tqdpb.ParcelLayersResponse{
		Data: h.mapper.ToProtoLayerDetailItems(data, req.ParcelId),
		Pagination: &tqdpb.Pagination{ // ✅ BỎ COMMENT
			Page:    pagable.GetPage(),
			Limit:   pagable.GetSize(),
			Total:   uint32(total),
			HasMore: total > int64(pagable.GetOffset()+pagable.GetLimit()),
		},
		Meta: &tqdpb.LayersMeta{
			PartialData:  false,
			FailedLayers: []*tqdpb.FailedLayer{},
		},
		Links: &tqdpb.Links{
			Self: fmt.Sprintf("/v2/tqd/parcels/%d/layers?page=%d&limit=%d", req.ParcelId, pagable.GetPage(), pagable.GetSize()),
		},
	}

	return resp, nil
}

func (h *ParcelGrpcHandler) GetZoneGeometry(ctx context.Context, req *tqdpb.GetZoneGeometryRequest) (*tqdpb.ZoneGeometryResponse, error) {
	log.Printf("[GetZoneGeometry] ParcelID: %d, LayerID: %d, ZoneID: %d", req.ParcelId, req.LayerId, req.ZoneId)

	if req == nil || req.ParcelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcel_id is required")
	}
	if req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layer_id is required")
	}
	if req.ZoneId == 0 {
		return nil, status.Error(codes.InvalidArgument, "zone_id is required")
	}

	geoJSON, props, err := h.usecase.GetZoneGeometry(ctx, req.ParcelId, req.LayerId, req.ZoneId)
	if err != nil {
		log.Printf("[GetZoneGeometry] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	if geoJSON == "" {
		return nil, status.Error(codes.NotFound, "zone geometry not found")
	}

	return &tqdpb.ZoneGeometryResponse{
		Type: "Feature",
		Properties: &tqdpb.ZoneGeometryProperties{
			ZoneId:           req.ZoneId,
			ZoneName:         props.ZoneName,
			LandUseCode:      props.LandUseCode,
			LandUseName:      props.LandUseName,
			LandUseColor:     props.LandUseColor,
			BuildStatus:      uint32(props.BuildStatus),
			AlertLevel:       uint32(props.AlertLevel),
			LegalDocumentIds: props.LegalDocumentIds,
		},
		Geometry: geoJSON,
	}, nil
}

func (h *ParcelGrpcHandler) GetParcelLegalDocuments(ctx context.Context, req *tqdpb.GetParcelLegalDocumentsRequest) (*tqdpb.ParcelLegalDocumentsResponse, error) {
	log.Printf("[GetParcelLegalDocuments] ParcelID: %d, Page: %d, Limit: %d", req.ParcelId, req.Page, req.Limit)

	if req == nil || req.ParcelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcel_id is required")
	}

	pagable := _dto.Pagable{
		Page: req.Page,
		Size: req.Limit,
	}
	pagable.Normalize()

	docs, total, err := h.usecase.GetParcelLegalDocuments(ctx, req.ParcelId, req.Type, req.Status, pagable)
	if err != nil {
		log.Printf("[GetParcelLegalDocuments] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	items := h.mapper.ToProtoLegalDocumentItems(docs)
	summary := h.mapper.ToProtoLegalDocumentSummary(docs)

	return &tqdpb.ParcelLegalDocumentsResponse{
		Data: items,
		Summary: &tqdpb.LegalDocumentSummary{
			Total:             uint32(total),
			ByType:            summary.ByType,
			ByStatus:          summary.ByStatus,
			HasSupersededDocs: summary.HasSupersededDocs,
		},
		Filters: &tqdpb.LegalFilters{
			Applied: &tqdpb.LegalFiltersApplied{
				Type:   req.Type,
				Status: req.Status,
			},
			AvailableTypes:    []string{"PLANNING_DECISION", "LAND_CERTIFICATE", "LAW", "DECREE", "REGULATION", "CIRCULAR"}, // ✅ BỎ COMMENT
			AvailableStatuses: []string{"ACTIVE", "SUPERSEDED", "EXPIRED", "DRAFT"},                                         // ✅ BỎ COMMENT
		},
		Meta: &tqdpb.ResponseMetadata{
			DataVersion: "2.1.0",
			DataQuality: "VERIFIED",
			QueriedAt:   time.Now().UTC().Format(time.RFC3339),
			PartialData: false,
		},
		Links: &tqdpb.Links{
			Self: fmt.Sprintf("/v2/tqd/parcels/%d/legal-documents", req.ParcelId),
		},
	}, nil
}
