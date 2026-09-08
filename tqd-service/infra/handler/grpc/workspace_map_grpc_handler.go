package handler_grpc

import (
	"context"
	"errors"
	"log"
	"strings"

	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"

	tqdpb "pb/types/tqd"
	"tqd/infra/handler/grpc/generatedreport"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type MapWorkspaceGrpcHandler struct {
	tqdpb.UnimplementedMapWorkspaceServiceServer

	usecase      usecase.MapWorkspaceUsecase
	mapper       *mapper.WorkspaceMapper
	SyncProvider *_utils.SyncUtil
	createReport *grpcadapter.Adapter
}

func NewMapWorkspaceGrpcHandler(
	usecase usecase.MapWorkspaceUsecase,
	mapper *mapper.WorkspaceMapper,
	syncProvider *_utils.SyncUtil,
) *MapWorkspaceGrpcHandler {
	return &MapWorkspaceGrpcHandler{
		usecase:      usecase,
		mapper:       mapper,
		SyncProvider: syncProvider,
	}
}

// UseCreateReportAdapter attaches the only production owner of
// CreateGeneratedReport. Absence means capability unavailable, not fallback.
func (h *MapWorkspaceGrpcHandler) UseCreateReportAdapter(adapter *grpcadapter.Adapter) {
	h.createReport = adapter
}

func (h *MapWorkspaceGrpcHandler) currentUserID(ctx context.Context) (uint64, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return 0, _errors.UnauthorizedException()
	}

	return userID, nil
}

func workspaceStatusError(err error) error {
	if err == nil {
		return nil
	}

	msg := err.Error()

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return status.Error(codes.NotFound, "record not found")
	case strings.Contains(msg, "not found"):
		return status.Error(codes.NotFound, msg)
	case strings.Contains(msg, "is required"), strings.Contains(msg, "invalid"):
		return status.Error(codes.InvalidArgument, msg)
	case strings.Contains(msg, "cannot regenerate"),
		strings.Contains(msg, "not ready to share"),
		strings.Contains(msg, "does not support"):
		return status.Error(codes.FailedPrecondition, msg)
	default:
		return status.Error(codes.Internal, msg)
	}
}

func protoPoint(v dto.SpatialPointDTO) *tqdpb.SpatialPoint {
	return &tqdpb.SpatialPoint{
		Lat: v.Lat,
		Lon: v.Lon,
	}
}

func protoBounds(v dto.SpatialBoundsDTO) *tqdpb.SpatialBounds {
	return &tqdpb.SpatialBounds{
		MinLon: v.MinLon,
		MinLat: v.MinLat,
		MaxLon: v.MaxLon,
		MaxLat: v.MaxLat,
	}
}

func protoGeometry(v dto.GeometryPreviewDTO) *tqdpb.GeometryPreview {
	return &tqdpb.GeometryPreview{
		Type:     v.Type,
		GeoJson:  v.GeoJSON,
		Centroid: protoPoint(v.Centroid),
		Bounds:   protoBounds(v.Bounds),
		Srid:     v.SRID,
	}
}

func protoPreviewHint(v dto.PreviewHintDTO) *tqdpb.PreviewHint {
	return &tqdpb.PreviewHint{
		Mode:            v.Mode,
		Bounds:          protoBounds(v.Bounds),
		Centroid:        protoPoint(v.Centroid),
		FitPaddingRatio: v.FitPaddingRatio,
		ThumbnailUrl:    v.ThumbnailURL,
		HighlightLayer:  v.HighlightLayer,
	}
}

func protoFocusTarget(v dto.SpatialFocusTargetDTO) *tqdpb.SpatialFocusTarget {
	return &tqdpb.SpatialFocusTarget{
		Type:     v.Type,
		Bounds:   protoBounds(v.Bounds),
		Centroid: protoPoint(v.Centroid),
		GeoJson:  v.GeoJSON,
		Srid:     v.SRID,
	}
}

func protoLocation(v dto.WorkspaceLocationDTO) *tqdpb.WorkspaceLocation {
	return &tqdpb.WorkspaceLocation{
		Address:      v.Address,
		Province:     v.Province,
		ProvinceCode: v.ProvinceCode,
		WardCode:     v.WardCode,
	}
}

func protoParcelMeta(v dto.WorkspaceParcelMetaDTO) *tqdpb.WorkspaceParcelMeta {
	return &tqdpb.WorkspaceParcelMeta{
		MapNumber:  v.MapNumber,
		LandNumber: v.LandNumber,
		AreaSqm:    v.AreaSqm,
	}
}

func followedDTOToProto(item dto.FollowedParcelPreviewDTO) *tqdpb.FollowedParcelPreview {
	return &tqdpb.FollowedParcelPreview{
		Id:              item.ID,
		Type:            item.Type,
		FollowId:        item.FollowID,
		UserId:          item.UserID,
		ParcelId:        item.ParcelID,
		Title:           item.Title,
		Subtitle:        item.Subtitle,
		Parcel:          protoParcelMeta(item.Parcel),
		Location:        protoLocation(item.Location),
		Lat:             item.Lat,
		Lon:             item.Lon,
		GeometryPreview: protoGeometry(item.GeometryPreview),
		Preview:         protoPreviewHint(item.Preview),
		FocusTarget:     protoFocusTarget(item.FocusTarget),
		FollowedAt:      item.FollowedAt,
		Actions: &tqdpb.FollowedParcelActions{
			Focus:        item.Actions.Focus,
			Remove:       item.Actions.Remove,
			Compare:      item.Actions.Compare,
			CreateReport: item.Actions.CreateReport,
		},
	}
}

func historyDTOToProto(item dto.ViewHistoryPreviewDTO) *tqdpb.ViewHistoryPreview {
	return &tqdpb.ViewHistoryPreview{
		Id:              item.ID,
		Type:            item.Type,
		HistoryId:       item.HistoryID,
		UserId:          item.UserID,
		EntityType:      item.EntityType,
		EntityId:        item.EntityID,
		ParcelId:        item.ParcelID,
		RegionId:        item.RegionID,
		Title:           item.Title,
		Subtitle:        item.Subtitle,
		Parcel:          protoParcelMeta(item.Parcel),
		Location:        protoLocation(item.Location),
		Lat:             item.Lat,
		Lon:             item.Lon,
		GeometryPreview: protoGeometry(item.GeometryPreview),
		Preview:         protoPreviewHint(item.Preview),
		FocusTarget:     protoFocusTarget(item.FocusTarget),
		ViewedAt:        item.ViewedAt,
		ViewContext: &tqdpb.ViewHistoryContext{
			Source: item.ViewContext.Source,
			Zoom:   item.ViewContext.Zoom,
		},
		ViewCount:        item.ViewCount,
		CountedViewCount: item.CountedViewCount,
		Actions: &tqdpb.ViewHistoryActions{
			Focus:  item.Actions.Focus,
			Remove: item.Actions.Remove,
			Follow: item.Actions.Follow,
		},
	}
}

func reportDTOToProto(item dto.GeneratedReportPreviewDTO) *tqdpb.GeneratedReportPreview {
	return &tqdpb.GeneratedReportPreview{
		Id:         item.ID,
		Type:       item.Type,
		ReportId:   item.ReportID,
		UserId:     item.UserID,
		ReportType: item.ReportType,
		Status:     item.Status,
		Title:      item.Title,
		Subtitle:   item.Subtitle,

		Location: protoLocation(item.Location),
		Spatial: &tqdpb.ReportSpatialPreview{
			Bounds:   protoBounds(item.Spatial.Bounds),
			Centroid: protoPoint(item.Spatial.Centroid),
		},
		Assets: &tqdpb.ReportAssets{
			ThumbnailUrl: item.Assets.ThumbnailURL,
			ImageUrl:     item.Assets.ImageURL,
			PdfUrl:       item.Assets.PDFURL,
			ShareUrl:     item.Assets.ShareURL,
		},
		ReportMeta: &tqdpb.ReportMeta{
			CreatedAt: item.Meta.CreatedAt,
			UpdatedAt: item.Meta.UpdatedAt,
			ExpiresAt: item.Meta.ExpiresAt,
			FileSize:  item.Meta.FileSize,
			Format:    item.Meta.Format,
		},
		Comparison: &tqdpb.ReportComparison{
			FromPlanName: item.Comparison.FromPlanName,
			ToPlanName:   item.Comparison.ToPlanName,
			FromYear:     item.Comparison.FromYear,
			ToYear:       item.Comparison.ToYear,
		},
		Actions: &tqdpb.ReportActions{
			Focus:         item.Actions.Focus,
			DownloadImage: item.Actions.DownloadImage,
			DownloadPdf:   item.Actions.DownloadPDF,
			Share:         item.Actions.Share,
			Regenerate:    item.Actions.Regenerate,
			Remove:        item.Actions.Remove,
		},

		EntityType: item.EntityType,
		EntityId:   item.EntityID,
		ParcelId:   item.ParcelID,
		RegionId:   item.RegionID,
	}
}

func workspaceLocationFromProto(v *tqdpb.WorkspaceLocation) dto.WorkspaceLocationDTO {
	if v == nil {
		return dto.WorkspaceLocationDTO{}
	}

	return dto.WorkspaceLocationDTO{
		Address:      v.GetAddress(),
		Province:     v.GetProvince(),
		ProvinceCode: v.GetProvinceCode(),
		WardCode:     v.GetWardCode(),
	}
}

func reportSpatialFromProto(v *tqdpb.ReportSpatialPreview) dto.ReportSpatialPreviewDTO {
	if v == nil {
		return dto.ReportSpatialPreviewDTO{}
	}

	return dto.ReportSpatialPreviewDTO{
		Bounds: dto.SpatialBoundsDTO{
			MinLon: v.GetBounds().GetMinLon(),
			MinLat: v.GetBounds().GetMinLat(),
			MaxLon: v.GetBounds().GetMaxLon(),
			MaxLat: v.GetBounds().GetMaxLat(),
		},
		Centroid: dto.SpatialPointDTO{
			Lat: v.GetCentroid().GetLat(),
			Lon: v.GetCentroid().GetLon(),
		},
	}
}

func reportComparisonFromProto(v *tqdpb.ReportComparison) dto.ReportComparisonDTO {
	if v == nil {
		return dto.ReportComparisonDTO{}
	}

	return dto.ReportComparisonDTO{
		FromPlanName: v.GetFromPlanName(),
		ToPlanName:   v.GetToPlanName(),
		FromYear:     v.GetFromYear(),
		ToYear:       v.GetToYear(),
	}
}

func (h *MapWorkspaceGrpcHandler) ListFollowedParcels(
	ctx context.Context,
	req *tqdpb.WorkspacePageRequest,
) (*tqdpb.ListFollowedParcelsResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)

	result, err := h.usecase.ListFollowedParcels(ctx, userID, pagable)
	if err != nil {
		log.Printf("[DEBUG][Handler][ListFollowedParcels][ERROR] userID=%d err=%v", userID, err)
		return nil, workspaceStatusError(err)
	}

	items := make([]*tqdpb.FollowedParcelPreview, 0)
	total := int64(0)
	nextCursor := ""

	if result != nil {
		total = result.Total
		nextCursor = result.NextCursor

		items = make([]*tqdpb.FollowedParcelPreview, 0, len(result.Items))
		for _, item := range result.Items {
			items = append(items, followedDTOToProto(item))
		}
	} else {
		log.Printf(
			"[DEBUG][Handler][ListFollowedParcels][USECASE_RESULT] userID=%d resultNil=true",
			userID,
		)
	}

	return &tqdpb.ListFollowedParcelsResponse{
		Data:       items,
		NextCursor: nextCursor,
		Total:      total,
		Page:       int32(pagable.GetPage()),
		Size:       int32(pagable.GetSize()),
	}, nil
}

func (h *MapWorkspaceGrpcHandler) FollowParcel(
	ctx context.Context,
	req *tqdpb.FollowParcelRequest,
) (*tqdpb.WorkspaceMutationResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetParcelId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcelId is required")
	}

	log.Printf("[FollowParcel] UserID: %d, ParcelID: %d", userID, req.GetParcelId())

	id, err := h.usecase.FollowParcel(ctx, dto.FollowParcelRequestDTO{
		UserID:   userID,
		ParcelID: req.GetParcelId(),
		Note:     req.GetNote(),
	})
	if err != nil {
		log.Printf("[FollowParcel] Error: %v", err)
		return nil, workspaceStatusError(err)
	}

	return &tqdpb.WorkspaceMutationResponse{
		Success: true,
		Message: "followed parcel saved",
		Id:      id,
	}, nil
}

func (h *MapWorkspaceGrpcHandler) RemoveFollowedParcel(
	ctx context.Context,
	req *tqdpb.RemoveFollowedParcelRequest,
) (*tqdpb.WorkspaceMutationResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetFollowId() == 0 && req.GetParcelId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "followId or parcelId is required")
	}

	log.Printf(
		"[RemoveFollowedParcel] UserID: %d, FollowID: %d, ParcelID: %d",
		userID,
		req.GetFollowId(),
		req.GetParcelId(),
	)

	err = h.usecase.RemoveFollowedParcel(ctx, dto.RemoveFollowedParcelRequestDTO{
		UserID:   userID,
		FollowID: req.GetFollowId(),
		ParcelID: req.GetParcelId(),
	})
	if err != nil {
		log.Printf("[RemoveFollowedParcel] Error: %v", err)
		return nil, workspaceStatusError(err)
	}

	return &tqdpb.WorkspaceMutationResponse{
		Success: true,
		Message: "followed parcel removed",
		Id:      req.GetFollowId(),
	}, nil
}

func (h *MapWorkspaceGrpcHandler) ListViewHistory(
	ctx context.Context,
	req *tqdpb.WorkspacePageRequest,
) (*tqdpb.ListViewHistoryResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)

	log.Printf(
		"[ListViewHistory] UserID: %d, Page: %d, Size: %d, Limit: %d, Offset: %d",
		userID,
		pagable.GetPage(),
		pagable.GetSize(),
		pagable.GetLimit(),
		pagable.GetOffset(),
	)

	result, err := h.usecase.ListViewHistory(ctx, userID, pagable)
	if err != nil {
		log.Printf("[ListViewHistory] Error: %v", err)
		return nil, workspaceStatusError(err)
	}

	items := make([]*tqdpb.ViewHistoryPreview, 0)
	total := int64(0)
	nextCursor := ""

	if result != nil {
		total = result.Total
		nextCursor = result.NextCursor

		items = make([]*tqdpb.ViewHistoryPreview, 0, len(result.Items))
		for _, item := range result.Items {
			items = append(items, historyDTOToProto(item))
		}
	}

	log.Printf(
		"[ListViewHistory] UserID: %d, Total: %d, Items: %d",
		userID,
		total,
		len(items),
	)

	return &tqdpb.ListViewHistoryResponse{
		Data:       items,
		NextCursor: nextCursor,
		Total:      total,
		Page:       int32(pagable.GetPage()),
		Size:       int32(pagable.GetSize()),
	}, nil
}

func (h *MapWorkspaceGrpcHandler) AddViewHistory(
	ctx context.Context,
	req *tqdpb.AddViewHistoryRequest,
) (*tqdpb.WorkspaceMutationResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetEntityId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "entityId is required")
	}

	log.Printf(
		"[AddViewHistory] UserID: %d, EntityType: %d, EntityID: %d, ParcelID: %d, RegionID: %d",
		userID,
		req.GetEntityType(),
		req.GetEntityId(),
		req.GetParcelId(),
		req.GetRegionId(),
	)

	id, err := h.usecase.AddViewHistory(ctx, dto.AddViewHistoryRequestDTO{
		UserID:     userID,
		EntityType: req.GetEntityType(),
		EntityID:   req.GetEntityId(),
		ParcelID:   req.GetParcelId(),
		RegionID:   req.GetRegionId(),
		Source:     req.GetSource(),
		Zoom:       req.GetZoom(),
	})
	if err != nil {
		log.Printf("[AddViewHistory] Error: %v", err)
		return nil, workspaceStatusError(err)
	}

	return &tqdpb.WorkspaceMutationResponse{
		Success: true,
		Message: "view history added",
		Id:      id,
	}, nil
}

func (h *MapWorkspaceGrpcHandler) TrackViewHistory(
	ctx context.Context,
	req *tqdpb.TrackViewHistoryRequest,
) (*tqdpb.TrackViewHistoryResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetEntityId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "entityId is required")
	}

	log.Printf(
		"[TrackViewHistory] UserID: %d, EntityType: %d, EntityID: %d, ParcelID: %d, RegionID: %d, CountIntent: %v, VisibleMs: %d",
		userID,
		req.GetEntityType(),
		req.GetEntityId(),
		req.GetParcelId(),
		req.GetRegionId(),
		req.GetCountIntent(),
		req.GetVisibleMs(),
	)

	result, err := h.usecase.TrackViewHistory(ctx, dto.TrackViewHistoryRequestDTO{
		UserID:        userID,
		ClientEventID: req.GetClientEventId(),
		SessionID:     req.GetSessionId(),
		DeviceID:      req.GetDeviceId(),
		EntityType:    req.GetEntityType(),
		EntityID:      req.GetEntityId(),
		ParcelID:      req.GetParcelId(),
		RegionID:      req.GetRegionId(),
		Source:        req.GetSource(),
		SourceRef:     req.GetSourceRef(),
		RouteName:     req.GetRouteName(),
		Zoom:          req.GetZoom(),
		Center: dto.SpatialPointDTO{
			Lat: req.GetCenter().GetLat(),
			Lon: req.GetCenter().GetLon(),
		},
		ViewportBounds: dto.SpatialBoundsDTO{
			MinLon: req.GetViewportBounds().GetMinLon(),
			MinLat: req.GetViewportBounds().GetMinLat(),
			MaxLon: req.GetViewportBounds().GetMaxLon(),
			MaxLat: req.GetViewportBounds().GetMaxLat(),
		},
		VisibleMs:    req.GetVisibleMs(),
		CountIntent:  req.GetCountIntent(),
		MetadataJSON: req.GetMetadataJson(),
	})
	if err != nil {
		log.Printf("[TrackViewHistory] Error: %v", err)
		return nil, workspaceStatusError(err)
	}

	return &tqdpb.TrackViewHistoryResponse{
		Success:   true,
		HistoryId: result.HistoryID,
		EventId:   result.EventID,
		Counted:   result.Counted,
		ViewCount: result.ViewCount,
		Message:   "view tracked",
	}, nil
}

func (h *MapWorkspaceGrpcHandler) RemoveViewHistory(
	ctx context.Context,
	req *tqdpb.RemoveViewHistoryRequest,
) (*tqdpb.WorkspaceMutationResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetHistoryId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "historyId is required")
	}

	log.Printf("[RemoveViewHistory] UserID: %d, HistoryID: %d", userID, req.GetHistoryId())

	err = h.usecase.RemoveViewHistory(ctx, dto.RemoveViewHistoryRequestDTO{
		UserID:    userID,
		HistoryID: req.GetHistoryId(),
	})
	if err != nil {
		log.Printf("[RemoveViewHistory] Error: %v", err)
		return nil, workspaceStatusError(err)
	}

	return &tqdpb.WorkspaceMutationResponse{
		Success: true,
		Message: "view history removed",
		Id:      req.GetHistoryId(),
	}, nil
}

func (h *MapWorkspaceGrpcHandler) ClearViewHistory(
	ctx context.Context,
	req *tqdpb.ClearViewHistoryRequest,
) (*tqdpb.WorkspaceMutationResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("[ClearViewHistory] UserID: %d", userID)

	if err := h.usecase.ClearViewHistory(ctx, userID); err != nil {
		log.Printf("[ClearViewHistory] Error: %v", err)
		return nil, workspaceStatusError(err)
	}

	return &tqdpb.WorkspaceMutationResponse{
		Success: true,
		Message: "view history cleared",
	}, nil
}

func (h *MapWorkspaceGrpcHandler) ListGeneratedReports(
	ctx context.Context,
	req *tqdpb.WorkspacePageRequest,
) (*tqdpb.ListGeneratedReportsResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)

	log.Printf(
		"[ListGeneratedReports] UserID: %d, Page: %d, Size: %d, Limit: %d, Offset: %d",
		userID,
		pagable.GetPage(),
		pagable.GetSize(),
		pagable.GetLimit(),
		pagable.GetOffset(),
	)

	result, err := h.usecase.ListGeneratedReports(ctx, userID, pagable)
	if err != nil {
		log.Printf("[ListGeneratedReports] Error: %v", err)
		return nil, workspaceStatusError(err)
	}

	items := make([]*tqdpb.GeneratedReportPreview, 0)
	total := int64(0)
	nextCursor := ""

	if result != nil {
		total = result.Total
		nextCursor = result.NextCursor

		items = make([]*tqdpb.GeneratedReportPreview, 0, len(result.Items))
		for _, item := range result.Items {
			items = append(items, reportDTOToProto(item))
		}
	}

	log.Printf(
		"[ListGeneratedReports] UserID: %d, Total: %d, Items: %d",
		userID,
		total,
		len(items),
	)

	return &tqdpb.ListGeneratedReportsResponse{
		Data:       items,
		NextCursor: nextCursor,
		Total:      total,
		Page:       int32(pagable.GetPage()),
		Size:       int32(pagable.GetSize()),
	}, nil
}

func (h *MapWorkspaceGrpcHandler) CreateGeneratedReport(
	ctx context.Context,
	req *tqdpb.CreateGeneratedReportRequest,
) (*tqdpb.CreateGeneratedReportResponse, error) {
	if h.createReport == nil {
		return nil, status.Error(
			codes.Unavailable,
			"generated report runtime is unavailable",
		)
	}
	return h.createReport.CreateGeneratedReport(ctx, req)
}

func (h *MapWorkspaceGrpcHandler) GetGeneratedReport(
	ctx context.Context,
	req *tqdpb.GetGeneratedReportRequest,
) (*tqdpb.GeneratedReportPreview, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetReportId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "reportId is required")
	}

	log.Printf("[GetGeneratedReport] UserID: %d, ReportID: %d", userID, req.GetReportId())

	item, err := h.usecase.GetGeneratedReport(ctx, userID, req.GetReportId())
	if err != nil {
		log.Printf("[GetGeneratedReport] Error: %v", err)
		return nil, workspaceStatusError(err)
	}

	if item == nil {
		return nil, status.Error(codes.NotFound, "report not found")
	}

	return reportDTOToProto(*item), nil
}

func (h *MapWorkspaceGrpcHandler) RemoveGeneratedReport(
	ctx context.Context,
	req *tqdpb.RemoveGeneratedReportRequest,
) (*tqdpb.WorkspaceMutationResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetReportId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "reportId is required")
	}

	log.Printf("[RemoveGeneratedReport] UserID: %d, ReportID: %d", userID, req.GetReportId())

	err = h.usecase.RemoveGeneratedReport(ctx, dto.RemoveGeneratedReportRequestDTO{
		UserID:   userID,
		ReportID: req.GetReportId(),
	})
	if err != nil {
		log.Printf("[RemoveGeneratedReport] Error: %v", err)
		return nil, workspaceStatusError(err)
	}

	return &tqdpb.WorkspaceMutationResponse{
		Success: true,
		Message: "report removed",
		Id:      req.GetReportId(),
	}, nil
}

func (h *MapWorkspaceGrpcHandler) RegenerateReport(
	ctx context.Context,
	req *tqdpb.RegenerateReportRequest,
) (*tqdpb.WorkspaceMutationResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetReportId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "reportId is required")
	}

	log.Printf("[RegenerateReport] UserID: %d, ReportID: %d", userID, req.GetReportId())

	err = h.usecase.RegenerateReport(ctx, userID, req.GetReportId())
	if err != nil {
		log.Printf("[RegenerateReport] Error: %v", err)
		return nil, workspaceStatusError(err)
	}

	return &tqdpb.WorkspaceMutationResponse{
		Success: true,
		Message: "report regeneration started",
		Id:      req.GetReportId(),
	}, nil
}

func (h *MapWorkspaceGrpcHandler) ShareReport(
	ctx context.Context,
	req *tqdpb.ShareReportRequest,
) (*tqdpb.ShareReportResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetReportId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "reportId is required")
	}

	log.Printf("[ShareReport] UserID: %d, ReportID: %d", userID, req.GetReportId())

	result, err := h.usecase.ShareReport(ctx, userID, req.GetReportId())
	if err != nil {
		log.Printf("[ShareReport] Error: %v", err)
		return nil, workspaceStatusError(err)
	}

	return &tqdpb.ShareReportResponse{
		Success:  result.Success,
		ShareUrl: result.ShareURL,
		Message:  result.Message,
	}, nil
}
