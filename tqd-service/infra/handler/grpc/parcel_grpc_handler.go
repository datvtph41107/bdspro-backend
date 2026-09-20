package handler_grpc

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
	"log/slog"
	"strings"
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
	SyncProvider *_utils.SyncUtil
}

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
	slog.InfoContext(ctx, fmt.Sprintf("[CreateBatchParcels] UserID: %d, Batch size: %d", userID, len(req.Data)))

	// Convert proto to request DTO
	createReq := h.mapper.ProtoToDomains(req)

	// Call usecase
	result, err := h.usecase.CreateBatch(ctx, &dto.ParcelList{
		Data: createReq,
	})
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[CreateBatchParcels] Error: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to create parcel: %v", err)
	}
	slog.InfoContext(ctx, fmt.Sprintf("[CreateBatchParcels] Successfully created %d parcels", len(result.Data)))

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
	slog.InfoContext(ctx, fmt.Sprintf("[FindParcelByLocation] Lat: %f, Lng: %f", req.Latitude, req.Longitude))

	// Call usecase
	result, err := h.usecase.FindByLocation(ctx, req.Latitude, req.Longitude)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[FindParcelByLocation] Error: %v", err))
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
		slog.InfoContext(ctx, fmt.Sprintf("[FindParcelsByPolygon] Points=%d, Intersect=%v, MaxAreaKm2=%f, Page=%d, Size=%d, Z=<nil>, Mode=legacy-parcel, HasFilter=%v, IncludeAnalysis=%v",
			len(points),
			req.Intersect,
			req.MaxAreaHa,
			page,
			pageSize,
			filter != nil,
			includeAnalysis),
		)
	} else {
		slog.InfoContext(ctx, fmt.Sprintf("[FindParcelsByPolygon] Points=%d, Intersect=%v, MaxAreaKm2=%f, Page=%d, Size=%d, Z=%d, HasFilter=%v, IncludeAnalysis=%v",
			len(points),
			req.Intersect,
			req.MaxAreaHa,
			page,
			pageSize,
			*z,
			filter != nil,
			includeAnalysis),
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
		slog.ErrorContext(ctx, fmt.Sprintf("[FindParcelsByPolygon] Usecase error: %v", err))
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
	slog.InfoContext(ctx, fmt.Sprintf("[FindParcelsByPolygon] ResultType=%d, Z=%d, Parcels=%d, Regions=%d, Total=%d",
		result.ResultType,
		result.Z,
		len(result.Parcels),
		len(result.Regions),
		result.Total),
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
	slog.InfoContext(ctx, fmt.Sprintf("[CheckingPolygon] Points=%d, Intersect=%v, MaxAreaKm2=%f, Page=%d, Size=%d, Z=%v, HasFilter=%v, IncludeAnalysis=%v",
		len(points),
		req.Intersect,
		req.MaxAreaHa,
		page,
		pageSize,
		z,
		filter != nil,
		includeAnalysis),
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
		slog.ErrorContext(ctx, fmt.Sprintf("[CheckingPolygon] Usecase error: %v", err))
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
	slog.InfoContext(ctx, fmt.Sprintf("[CheckingPolygon] ResultType=%d, Z=%d, Parcels=%d, Regions=%d, Total=%d, handlerTotal=%v",
		result.ResultType,
		result.Z,
		len(result.Parcels),
		len(result.Regions),
		result.Total,
		time.Since(handlerStart)),
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
	slog.InfoContext(ctx, fmt.Sprintf("[SearchTxtClientParcels] UserID: %d, OriginalText: %s, CleanedText: %s, Map: %s, Land: %s, Page: %d, Size: %d",
		userID, req.Text, searchInfo.CleanedText, searchInfo.MapNumber, searchInfo.LandNumber, req.Page, req.Size))

	pagable := _dto.Pagable{
		Page: req.Page,
		Size: req.Size,
	}
	pagable.Normalize()

	// Gọi SearchParcelsByText với các thông tin đã bóc tách
	parcels, total, err := h.usecase.SearchParcelsByText(ctx, searchInfo.CleanedText, searchInfo.MapNumber, searchInfo.LandNumber, int(pagable.GetPage()), int(pagable.GetSize()))
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[SearchTxtClientParcels] Error: %v", err))
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
	slog.InfoContext(ctx, fmt.Sprintf("[SearchPublicParcels] Page: %d, Size: %d", req.Page, req.Size))

	pagable := _dto.Pagable{
		Page: req.Page,
		Size: req.Size,
	}
	pagable.Normalize()

	parcels, total, err := h.usecase.SearchPublicParcels(ctx, int(pagable.GetPage()), int(pagable.GetSize()))
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[SearchPublicParcels] Error: %v", err))
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
	slog.InfoContext(ctx, fmt.Sprintf("[PreviewParcelsById] UserID: %d, ParcelID: %d", userID, req.Id))

	// Gọi đúng method GetParcelInfoByID (không phải GetParcelDetail)
	info, err := h.usecase.GetParcelInfoByID(ctx, req.Id)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[PreviewParcelsById] Error: %v", err))
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
	slog.InfoContext(ctx, fmt.Sprintf("[PreviewParcelsById] UserID: %d, ParcelID: %d", userID, req.Id))

	// Gọi đúng method GetParcelInfoByID (không phải GetParcelDetail)
	info, err := h.usecase.GetParcelInfoByID(ctx, req.Id)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[PreviewParcelsById] Error: %v", err))
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

func (h *ParcelGrpcHandler) GetParcelDetail(ctx context.Context, req *tqdpb.GetParcelDetailRequest) (*tqdpb.ParcelDetailResponse, error) {
	slog.InfoContext(ctx, fmt.Sprintf("[GetParcelDetail] ParcelID: %d", req.ParcelId))

	if req == nil || req.ParcelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcel_id is required")
	}

	data, err := h.usecase.GetParcelDetail(ctx, req.ParcelId)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[GetParcelDetail] Error: %v", err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if data == nil {
		return nil, status.Error(codes.NotFound, "parcel not found")
	}

	return h.mapper.ToProtoParcelDetail(data), nil
}

func (h *ParcelGrpcHandler) GetParcelLayers(ctx context.Context, req *tqdpb.GetParcelLayersRequest) (*tqdpb.ParcelLayersResponse, error) {
	slog.InfoContext(ctx, fmt.Sprintf("[GetParcelLayers] ParcelID: %d, Page: %d, Limit: %d", req.ParcelId, req.Page, req.Limit))

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
		slog.ErrorContext(ctx, fmt.Sprintf("[GetParcelLayers] Error: %v", err))
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
	slog.InfoContext(ctx, fmt.Sprintf("[GetZoneGeometry] ParcelID: %d, LayerID: %d, ZoneID: %d", req.ParcelId, req.LayerId, req.ZoneId))

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
		slog.ErrorContext(ctx, fmt.Sprintf("[GetZoneGeometry] Error: %v", err))
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
	slog.InfoContext(ctx, fmt.Sprintf("[GetParcelLegalDocuments] ParcelID: %d, Page: %d, Limit: %d", req.ParcelId, req.Page, req.Limit))

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
		slog.ErrorContext(ctx, fmt.Sprintf("[GetParcelLegalDocuments] Error: %v", err))
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
