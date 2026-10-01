package controller

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	_dto "common/domain/dto"
	_errors "common/errors"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	"userv4/modules/tqd/mapper"
	"userv4/modules/tqd/services"
	genreportapp "userv4/modules/tqd/services/generatedreport/application"
)

// MapWorkspaceController owns /v2/tqd/map-workspace/* routes.
type MapWorkspaceController struct {
	workspace       services.MapWorkspaceUsecase
	createReportSvc *genreportapp.Service
	mapper          *mapper.WorkspaceMapper
}

func NewMapWorkspaceController(
	workspace services.MapWorkspaceUsecase,
	createReport *genreportapp.Service,
) *MapWorkspaceController {
	return &MapWorkspaceController{
		workspace:       workspace,
		createReportSvc: createReport,
		mapper:          mapper.NewWorkspaceMapper(),
	}
}

func (c *MapWorkspaceController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/map-workspace/followed-parcels", c.listFollowed)
	httpapi.POST(mux, "/v2/tqd/map-workspace/followed-parcels", c.followParcel)
	httpapi.DELETE(mux, "/v2/tqd/map-workspace/followed-parcels/{followId}", c.removeFollowedByID)
	httpapi.DELETE(mux, "/v2/tqd/map-workspace/followed-parcels/by-parcel/{parcelId}", c.removeFollowedByParcel)

	httpapi.GET(mux, "/v2/tqd/map-workspace/view-history", c.listViewHistory)
	httpapi.POST(mux, "/v2/tqd/map-workspace/view-history", c.addViewHistory)
	httpapi.POST(mux, "/v2/tqd/map-workspace/view-history:track", c.trackViewHistory)
	httpapi.DELETE(mux, "/v2/tqd/map-workspace/view-history/{historyId}", c.removeViewHistory)
	httpapi.DELETE(mux, "/v2/tqd/map-workspace/view-history", c.clearViewHistory)

	httpapi.POST(mux, "/v2/tqd/map-workspace/reports", c.createReport)
	httpapi.GET(mux, "/v2/tqd/map-workspace/reports", c.listReports)
	httpapi.GET(mux, "/v2/tqd/map-workspace/reports/{reportId}", c.getReport)
	httpapi.DELETE(mux, "/v2/tqd/map-workspace/reports/{reportId}", c.removeReport)
	// ServeMux: last segment is "{id}:regenerate" / "{id}:share".
	httpapi.POST(mux, "/v2/tqd/map-workspace/reports/{seg}", c.dispatchReportAction)
}

func (c *MapWorkspaceController) requireWorkspace() error {
	if c.workspace == nil {
		return _errors.ReturnError(501, "map workspace service not configured")
	}
	return nil
}

func workspacePage(req dto.WorkspacePageHTTPRequest) *_dto.Pagable {
	p := &_dto.Pagable{Page: req.Page, Size: req.Size}
	p.Normalize()
	return p
}

func (c *MapWorkspaceController) listFollowed(ctx context.Context, req dto.WorkspacePageHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireWorkspace(); err != nil {
		return nil, err
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	pagable := workspacePage(req)
	result, err := c.workspace.ListFollowedParcels(ctx, userID, pagable)
	if err != nil {
		return nil, err
	}
	items := c.mapper.ToProtoFollowedParcels(nil)
	total := int64(0)
	if result != nil {
		items = c.mapper.ToProtoFollowedParcels(result.Items)
		total = result.Total
	}
	return protoListResult(items, total, pagable.GetPage(), pagable.GetSize())
}

func (c *MapWorkspaceController) followParcel(ctx context.Context, req dto.FollowParcelHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireWorkspace(); err != nil {
		return nil, err
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if req.ParcelID == 0 {
		return nil, _errors.ReturnError(400, "parcelId is required")
	}
	id, err := c.workspace.FollowParcel(ctx, dto.FollowParcelRequestDTO{
		UserID:   userID,
		ParcelID: req.ParcelID,
		Note:     req.Note,
	})
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true, "message": "followed parcel saved", "id": id})
}

func (c *MapWorkspaceController) removeFollowedByID(ctx context.Context, req dto.RemoveFollowedByIDHTTPRequest) (*httpapi.Result, error) {
	return c.removeFollowed(ctx, req.FollowID, 0)
}

func (c *MapWorkspaceController) removeFollowedByParcel(ctx context.Context, req dto.RemoveFollowedByParcelHTTPRequest) (*httpapi.Result, error) {
	return c.removeFollowed(ctx, 0, req.ParcelID)
}

func (c *MapWorkspaceController) removeFollowed(ctx context.Context, followID, parcelID uint64) (*httpapi.Result, error) {
	if err := c.requireWorkspace(); err != nil {
		return nil, err
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if followID == 0 && parcelID == 0 {
		return nil, _errors.ReturnError(400, "followId or parcelId is required")
	}
	if err := c.workspace.RemoveFollowedParcel(ctx, dto.RemoveFollowedParcelRequestDTO{
		UserID:   userID,
		FollowID: followID,
		ParcelID: parcelID,
	}); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true, "message": "followed parcel removed", "id": followID})
}

func (c *MapWorkspaceController) listViewHistory(ctx context.Context, req dto.WorkspacePageHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireWorkspace(); err != nil {
		return nil, err
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	pagable := workspacePage(req)
	result, err := c.workspace.ListViewHistory(ctx, userID, pagable)
	if err != nil {
		return nil, err
	}
	items := c.mapper.ToProtoHistories(nil)
	total := int64(0)
	if result != nil {
		items = c.mapper.ToProtoHistories(result.Items)
		total = result.Total
	}
	return protoListResult(items, total, pagable.GetPage(), pagable.GetSize())
}

func (c *MapWorkspaceController) addViewHistory(ctx context.Context, req dto.AddViewHistoryHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireWorkspace(); err != nil {
		return nil, err
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if req.EntityID == 0 {
		return nil, _errors.ReturnError(400, "entityId is required")
	}
	id, err := c.workspace.AddViewHistory(ctx, dto.AddViewHistoryRequestDTO{
		UserID:     userID,
		EntityType: req.EntityType,
		EntityID:   req.EntityID,
		ParcelID:   req.ParcelID,
		RegionID:   req.RegionID,
		Source:     req.Source,
		Zoom:       req.Zoom,
	})
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true, "message": "view history added", "id": id})
}

func (c *MapWorkspaceController) trackViewHistory(ctx context.Context, req dto.TrackViewHistoryHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireWorkspace(); err != nil {
		return nil, err
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if req.EntityID == 0 {
		return nil, _errors.ReturnError(400, "entityId is required")
	}
	result, err := c.workspace.TrackViewHistory(ctx, dto.TrackViewHistoryRequestDTO{
		UserID:         userID,
		ClientEventID:  req.ClientEventID,
		SessionID:      req.SessionID,
		DeviceID:       req.DeviceID,
		EntityType:     req.EntityType,
		EntityID:       req.EntityID,
		ParcelID:       req.ParcelID,
		RegionID:       req.RegionID,
		Source:         req.Source,
		SourceRef:      req.SourceRef,
		RouteName:      req.RouteName,
		Zoom:           req.Zoom,
		Center:         req.Center,
		ViewportBounds: req.ViewportBounds,
		VisibleMs:      req.VisibleMs,
		CountIntent:    req.CountIntent,
		MetadataJSON:   req.MetadataJSON,
	})
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{
		"success":   true,
		"historyId": result.HistoryID,
		"eventId":   result.EventID,
		"counted":   result.Counted,
		"viewCount": result.ViewCount,
		"message":   "view tracked",
	})
}

func (c *MapWorkspaceController) removeViewHistory(ctx context.Context, req dto.RemoveViewHistoryHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireWorkspace(); err != nil {
		return nil, err
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if req.HistoryID == 0 {
		return nil, _errors.ReturnError(400, "historyId is required")
	}
	if err := c.workspace.RemoveViewHistory(ctx, dto.RemoveViewHistoryRequestDTO{
		UserID:    userID,
		HistoryID: req.HistoryID,
	}); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true, "message": "view history removed", "id": req.HistoryID})
}

func (c *MapWorkspaceController) clearViewHistory(ctx context.Context, _ dto.ClearViewHistoryHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireWorkspace(); err != nil {
		return nil, err
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if err := c.workspace.ClearViewHistory(ctx, userID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true, "message": "view history cleared"})
}

func (c *MapWorkspaceController) createReport(ctx context.Context, req dto.CreateGeneratedReportHTTPRequest) (*httpapi.Result, error) {
	if c.createReportSvc == nil {
		return nil, _errors.ReturnError(501, "generated report runtime is unavailable")
	}
	if _, err := requireProfileID(ctx); err != nil {
		return nil, err
	}
	input := genreportapp.Input{
		ReportType: req.ReportType,
		EntityType: req.EntityType,
		EntityID:   req.EntityID,
		ParcelID:   req.ParcelID,
		RegionID:   req.RegionID,
		Title:      req.Title,
		Subtitle:   req.Subtitle,
	}
	if req.Location != nil {
		input.Location = genreportapp.Location{
			Address:      req.Location.Address,
			Province:     req.Location.Province,
			ProvinceCode: req.Location.ProvinceCode,
			WardCode:     req.Location.WardCode,
		}
	}
	if req.Spatial != nil {
		input.Spatial = genreportapp.Spatial{
			Bounds: genreportapp.Bounds{
				MinLon: req.Spatial.Bounds.MinLon,
				MinLat: req.Spatial.Bounds.MinLat,
				MaxLon: req.Spatial.Bounds.MaxLon,
				MaxLat: req.Spatial.Bounds.MaxLat,
			},
			Centroid: genreportapp.Point{
				Lat: req.Spatial.Centroid.Lat,
				Lon: req.Spatial.Centroid.Lon,
			},
		}
	}
	if req.Comparison != nil {
		input.Comparison = genreportapp.Comparison{
			FromPlanName: req.Comparison.FromPlanName,
			ToPlanName:   req.Comparison.ToPlanName,
			FromYear:     req.Comparison.FromYear,
			ToYear:       req.Comparison.ToYear,
		}
	}
	result, err := c.createReportSvc.CreateGeneratedReport(ctx, input)
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{
		"reportId": result.Report.ID,
		"status":   result.Report.Status,
		"outcome":  result.Outcome,
	})
}

func (c *MapWorkspaceController) listReports(ctx context.Context, req dto.WorkspacePageHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireWorkspace(); err != nil {
		return nil, err
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	pagable := workspacePage(req)
	result, err := c.workspace.ListGeneratedReports(ctx, userID, pagable)
	if err != nil {
		return nil, err
	}
	items := c.mapper.ToProtoReports(nil)
	total := int64(0)
	if result != nil {
		items = c.mapper.ToProtoReports(result.Items)
		total = result.Total
	}
	return protoListResult(items, total, pagable.GetPage(), pagable.GetSize())
}

func (c *MapWorkspaceController) getReport(ctx context.Context, req dto.WorkspaceReportIDHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireWorkspace(); err != nil {
		return nil, err
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if req.ReportID == 0 {
		return nil, _errors.ReturnError(400, "reportId is required")
	}
	item, err := c.workspace.GetGeneratedReport(ctx, userID, req.ReportID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, _errors.ReturnError(404, "report not found")
	}
	return protoResult(c.mapper.ToProtoReport(*item))
}

func (c *MapWorkspaceController) removeReport(ctx context.Context, req dto.WorkspaceReportIDHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireWorkspace(); err != nil {
		return nil, err
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if req.ReportID == 0 {
		return nil, _errors.ReturnError(400, "reportId is required")
	}
	if err := c.workspace.RemoveGeneratedReport(ctx, dto.RemoveGeneratedReportRequestDTO{
		UserID:   userID,
		ReportID: req.ReportID,
	}); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true, "message": "report removed", "id": req.ReportID})
}

func (c *MapWorkspaceController) dispatchReportAction(ctx context.Context, req dto.WorkspaceReportActionHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireWorkspace(); err != nil {
		return nil, err
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	seg := strings.TrimSpace(req.Seg)
	idStr, verb, ok := strings.Cut(seg, ":")
	if !ok || idStr == "" {
		return nil, _errors.ReturnError(404, "not found")
	}
	reportID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || reportID == 0 {
		return nil, _errors.ReturnError(400, "reportId is required")
	}
	switch verb {
	case "regenerate":
		if err := c.workspace.RegenerateReport(ctx, userID, reportID); err != nil {
			return nil, err
		}
		return jsonOK(map[string]any{"success": true, "message": "report regeneration started", "id": reportID})
	case "share":
		result, err := c.workspace.ShareReport(ctx, userID, reportID)
		if err != nil {
			return nil, err
		}
		return jsonOK(map[string]any{
			"success":  result.Success,
			"shareUrl": result.ShareURL,
			"message":  result.Message,
		})
	default:
		return nil, _errors.ReturnError(404, "not found")
	}
}
