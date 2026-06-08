package handler_grpc

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	_provider "common/provider"
	_utils "common/utils"
	"context"
	"errors"
	"strings"

	tqdpb "pb/types/tqd"
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
	SyncProvider *_provider.SyncProvider
}

func NewMapWorkspaceGrpcHandler(
	usecase usecase.MapWorkspaceUsecase,
	mapper *mapper.WorkspaceMapper,
	syncProvider *_provider.SyncProvider,
) *MapWorkspaceGrpcHandler {
	return &MapWorkspaceGrpcHandler{
		usecase:      usecase,
		mapper:       mapper,
		SyncProvider: syncProvider,
	}
}

func (h *MapWorkspaceGrpcHandler) currentUserID(ctx context.Context) (uint64, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return 0, _errors.UnauthorizedException()
	}

	return userID, nil
}

func pagableFromWorkspaceRequest(req *tqdpb.WorkspacePageRequest) _dto.IPagable {
	page := req.GetPage()
	size := req.GetSize()

	return _dto.NewPagableFromGrpc(&page, &size, nil)
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
	case strings.Contains(msg, "cannot regenerate"), strings.Contains(msg, "not ready to share"), strings.Contains(msg, "does not support"):
		return status.Error(codes.FailedPrecondition, msg)
	default:
		return status.Error(codes.Internal, msg)
	}
}

func protoPoint(v dto.SpatialPointDTO) *tqdpb.SpatialPoint {
	return &tqdpb.SpatialPoint{Lat: v.Lat, Lon: v.Lon}
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
		Location:   protoLocation(item.Location),
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

	result, err := h.usecase.ListFollowedParcels(ctx, userID, pagableFromWorkspaceRequest(req))
	if err != nil {
		return nil, workspaceStatusError(err)
	}

	items := make([]*tqdpb.FollowedParcelPreview, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, followedDTOToProto(item))
	}

	return &tqdpb.ListFollowedParcelsResponse{Items: items, NextCursor: result.NextCursor, Total: result.Total}, nil
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

	id, err := h.usecase.FollowParcel(ctx, dto.FollowParcelRequestDTO{UserID: userID, ParcelID: req.GetParcelId(), Note: req.GetNote()})
	if err != nil {
		return nil, workspaceStatusError(err)
	}

	return &tqdpb.WorkspaceMutationResponse{Success: true, Message: "followed parcel saved", Id: id}, nil
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

	err = h.usecase.RemoveFollowedParcel(ctx, dto.RemoveFollowedParcelRequestDTO{UserID: userID, FollowID: req.GetFollowId(), ParcelID: req.GetParcelId()})
	if err != nil {
		return nil, workspaceStatusError(err)
	}

	return &tqdpb.WorkspaceMutationResponse{Success: true, Message: "followed parcel removed", Id: req.GetFollowId()}, nil
}

func (h *MapWorkspaceGrpcHandler) ListViewHistory(
	ctx context.Context,
	req *tqdpb.WorkspacePageRequest,
) (*tqdpb.ListViewHistoryResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	result, err := h.usecase.ListViewHistory(ctx, userID, pagableFromWorkspaceRequest(req))
	if err != nil {
		return nil, workspaceStatusError(err)
	}

	items := make([]*tqdpb.ViewHistoryPreview, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, historyDTOToProto(item))
	}

	return &tqdpb.ListViewHistoryResponse{Items: items, NextCursor: result.NextCursor, Total: result.Total}, nil
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
		return nil, workspaceStatusError(err)
	}

	return &tqdpb.WorkspaceMutationResponse{Success: true, Message: "view history added", Id: id}, nil
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

	err = h.usecase.RemoveViewHistory(ctx, dto.RemoveViewHistoryRequestDTO{UserID: userID, HistoryID: req.GetHistoryId()})
	if err != nil {
		return nil, workspaceStatusError(err)
	}

	return &tqdpb.WorkspaceMutationResponse{Success: true, Message: "view history removed", Id: req.GetHistoryId()}, nil
}

func (h *MapWorkspaceGrpcHandler) ClearViewHistory(
	ctx context.Context,
	req *tqdpb.ClearViewHistoryRequest,
) (*tqdpb.WorkspaceMutationResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.usecase.ClearViewHistory(ctx, userID); err != nil {
		return nil, workspaceStatusError(err)
	}

	return &tqdpb.WorkspaceMutationResponse{Success: true, Message: "view history cleared"}, nil
}

func (h *MapWorkspaceGrpcHandler) ListGeneratedReports(
	ctx context.Context,
	req *tqdpb.WorkspacePageRequest,
) (*tqdpb.ListGeneratedReportsResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	result, err := h.usecase.ListGeneratedReports(ctx, userID, pagableFromWorkspaceRequest(req))
	if err != nil {
		return nil, workspaceStatusError(err)
	}

	items := make([]*tqdpb.GeneratedReportPreview, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, reportDTOToProto(item))
	}

	return &tqdpb.ListGeneratedReportsResponse{Items: items, NextCursor: result.NextCursor, Total: result.Total}, nil
}

func (h *MapWorkspaceGrpcHandler) CreateGeneratedReport(
	ctx context.Context,
	req *tqdpb.CreateGeneratedReportRequest,
) (*tqdpb.CreateGeneratedReportResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	if req.GetReportType() == 0 {
		return nil, status.Error(codes.InvalidArgument, "reportType is required")
	}

	item, err := h.usecase.CreateGeneratedReport(ctx, dto.CreateGeneratedReportRequestDTO{
		UserID:          userID,
		ReportType:      req.GetReportType(),
		EntityType:      req.GetEntityType(),
		EntityID:        req.GetEntityId(),
		ParcelID:        req.GetParcelId(),
		RegionID:        req.GetRegionId(),
		Title:           req.GetTitle(),
		Subtitle:        req.GetSubtitle(),
		Location:        workspaceLocationFromProto(req.GetLocation()),
		Spatial:         reportSpatialFromProto(req.GetSpatial()),
		Comparison:      reportComparisonFromProto(req.GetComparison()),
		ClientRequestID: req.GetClientRequestId(),
		MetadataJSON:    req.GetMetadataJson(),
	})
	if err != nil {
		return nil, workspaceStatusError(err)
	}

	return &tqdpb.CreateGeneratedReportResponse{Success: true, Message: "report generation started", Item: reportDTOToProto(*item)}, nil
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

	item, err := h.usecase.GetGeneratedReport(ctx, userID, req.GetReportId())
	if err != nil {
		return nil, workspaceStatusError(err)
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

	err = h.usecase.RemoveGeneratedReport(ctx, dto.RemoveGeneratedReportRequestDTO{UserID: userID, ReportID: req.GetReportId()})
	if err != nil {
		return nil, workspaceStatusError(err)
	}

	return &tqdpb.WorkspaceMutationResponse{Success: true, Message: "report removed", Id: req.GetReportId()}, nil
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

	err = h.usecase.RegenerateReport(ctx, userID, req.GetReportId())
	if err != nil {
		return nil, workspaceStatusError(err)
	}

	return &tqdpb.WorkspaceMutationResponse{Success: true, Message: "report regeneration started", Id: req.GetReportId()}, nil
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

	result, err := h.usecase.ShareReport(ctx, userID, req.GetReportId())
	if err != nil {
		return nil, workspaceStatusError(err)
	}

	return &tqdpb.ShareReportResponse{Success: result.Success, ShareUrl: result.ShareURL, Message: result.Message}, nil
}
