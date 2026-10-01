package controller

import (
	"context"
	"encoding/json"
	"net/http"

	_dto "common/domain/dto"
	_errors "common/errors"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	"userv4/modules/tqd/enums"
	"userv4/modules/tqd/mapper"
	"userv4/modules/tqd/services"
	"userv4/modules/tqd/services/quota"
)

// SubscriptionController owns /v2/tqd/subscriptions/* and /v2/tqd/notifications/*.
type SubscriptionController struct {
	subs    services.SubscriptionUsecase
	notifs  services.NotificationUsecase
	parcels services.ParcelUsecase
}

func NewSubscriptionController(
	subs services.SubscriptionUsecase,
	notifs services.NotificationUsecase,
	parcels services.ParcelUsecase,
) *SubscriptionController {
	return &SubscriptionController{subs: subs, notifs: notifs, parcels: parcels}
}

func (c *SubscriptionController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/subscriptions/parcels", c.listParcelSubs)
	httpapi.GET(mux, "/v2/tqd/subscriptions/parcels/{parcelId}/status", c.parcelStatus)
	httpapi.POST(mux, "/v2/tqd/subscriptions/parcels/{parcelId}", c.createParcelSub)
	httpapi.DELETE(mux, "/v2/tqd/subscriptions/parcels/{parcelId}", c.deleteParcelSub)
	httpapi.GET(mux, "/v2/tqd/subscriptions/regions", c.listRegionSubs)
	httpapi.POST(mux, "/v2/tqd/subscriptions/regions", c.createRegionSub)
	httpapi.DELETE(mux, "/v2/tqd/subscriptions/regions/{subscriptionId}", c.deleteRegionSub)

	httpapi.GET(mux, "/v2/tqd/notifications", c.listNotifications)
	httpapi.PUT(mux, "/v2/tqd/notifications/{notificationId}/read", c.markRead)
	httpapi.PUT(mux, "/v2/tqd/notifications/mark-all-read", c.markAllRead)
}

func (c *SubscriptionController) listParcelSubs(ctx context.Context, req dto.SubscriptionListHTTPRequest) (*httpapi.Result, error) {
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.subs.ListParcelSubscriptions(ctx, userID, req.Status, int(pagable.GetPage()), int(pagable.GetSize()))
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"data": rows, "total": total, "page": pagable.GetPage(), "size": pagable.GetSize()})
}

func (c *SubscriptionController) parcelStatus(ctx context.Context, req dto.SubscriptionParcelIDHTTPRequest) (*httpapi.Result, error) {
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	subscribed, sub, err := c.subs.GetParcelStatus(ctx, userID, req.ParcelID)
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"subscribed": subscribed, "subscription": sub})
}

func (c *SubscriptionController) createParcelSub(ctx context.Context, req dto.CreateParcelSubscriptionHTTPRequest) (*httpapi.Result, error) {
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	sub, err := c.subs.CreateParcelSubscription(ctx, userID, req.ParcelID, req.Scope, req.TriggerConfig)
	if err != nil {
		return nil, err
	}
	return jsonOK(sub)
}

func (c *SubscriptionController) deleteParcelSub(ctx context.Context, req dto.SubscriptionParcelIDHTTPRequest) (*httpapi.Result, error) {
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if err := c.subs.DeleteParcelSubscription(ctx, userID, req.ParcelID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *SubscriptionController) listRegionSubs(ctx context.Context, req dto.SubscriptionListHTTPRequest) (*httpapi.Result, error) {
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	subs, regions, total, err := c.subs.ListUserRegionSubscriptions(ctx, userID, int(pagable.GetPage()), int(pagable.GetSize()))
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"data": subs, "regions": regions, "total": total, "page": pagable.GetPage(), "size": pagable.GetSize()})
}

func (c *SubscriptionController) createRegionSub(ctx context.Context, req dto.CreateRegionSubscriptionHTTPRequest) (*httpapi.Result, error) {
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	sub, err := c.subs.CreateRegionSubscription(ctx, userID, req.RegionID, req.Scope, req.TriggerConfig)
	if err != nil {
		return nil, err
	}
	return jsonOK(sub)
}

func (c *SubscriptionController) deleteRegionSub(ctx context.Context, req dto.SubscriptionIDHTTPRequest) (*httpapi.Result, error) {
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if err := c.subs.DeleteRegionSubscription(ctx, userID, req.SubscriptionID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *SubscriptionController) listNotifications(ctx context.Context, req dto.NotificationListHTTPRequest) (*httpapi.Result, error) {
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, unread, err := c.notifs.List(ctx, userID, req.Type, req.Read, int(pagable.GetPage()), int(pagable.GetSize()))
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{
		"data":   rows,
		"total":  total,
		"unread": unread,
		"page":   pagable.GetPage(),
		"size":   pagable.GetSize(),
	})
}

func (c *SubscriptionController) markRead(ctx context.Context, req dto.NotificationIDHTTPRequest) (*httpapi.Result, error) {
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if err := c.notifs.MarkRead(ctx, userID, req.NotificationID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *SubscriptionController) markAllRead(ctx context.Context, _ dto.MarkAllNotificationsReadHTTPRequest) (*httpapi.Result, error) {
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if err := c.notifs.MarkAllRead(ctx, userID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

// CommercialController owns /v2/tqd/commercial/* profile usage routes.
type CommercialController struct {
	usage *quota.ProfileUsageService
}

func NewCommercialController(usage *quota.ProfileUsageService) *CommercialController {
	return &CommercialController{usage: usage}
}

func (c *CommercialController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/commercial/generated-report-usage", c.generatedReportUsage)
	httpapi.GET(mux, "/v2/tqd/commercial/usage", c.listUsage)
}

func (c *CommercialController) generatedReportUsage(ctx context.Context, _ dto.CommercialUsageHTTPRequest) (*httpapi.Result, error) {
	if c.usage == nil {
		return nil, _errors.ReturnError(501, "profile usage projection is not configured")
	}
	if _, err := requireProfileID(ctx); err != nil {
		return nil, err
	}
	projection, err := c.usage.GetGeneratedReportUsage(ctx)
	if err != nil {
		return nil, err
	}
	return jsonOK(projection)
}

func (c *CommercialController) listUsage(ctx context.Context, _ dto.CommercialUsageHTTPRequest) (*httpapi.Result, error) {
	if c.usage == nil {
		return nil, _errors.ReturnError(501, "profile usage projection is not configured")
	}
	if _, err := requireProfileID(ctx); err != nil {
		return nil, err
	}
	rows, err := c.usage.ListUsage(ctx)
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"usage": rows})
}

// FeatureController owns /v2/tqd/features/*.
type FeatureController struct {
	features *services.FeatureUsecase
}

func NewFeatureController(features *services.FeatureUsecase) *FeatureController {
	return &FeatureController{features: features}
}

func (c *FeatureController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/features/nearby", c.nearby)
	httpapi.POST(mux, "/v2/tqd/features/polygon", c.byPolygon)
}

func (c *FeatureController) nearby(ctx context.Context, req dto.FeaturesNearbyHTTPRequest) (*httpapi.Result, error) {
	if c.features == nil {
		return nil, _errors.ReturnError(501, "feature service not configured")
	}
	rows, total, err := c.features.GetFeaturesByPointRadius(ctx, &services.PointRadiusRequest{
		Lat: req.Lat, Lng: req.Lng, Radius: req.Radius, Limit: req.Limit, Offset: req.Offset,
	})
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"data": rows, "total": total})
}

func (c *FeatureController) byPolygon(ctx context.Context, req dto.FeaturesPolygonHTTPRequest) (*httpapi.Result, error) {
	if c.features == nil {
		return nil, _errors.ReturnError(501, "feature service not configured")
	}
	raw, err := json.Marshal(req.Polygon)
	if err != nil {
		return nil, _errors.ReturnError(400, "invalid polygon")
	}
	rows, total, err := c.features.GetFeaturesByPolygon(ctx, raw, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"data": rows, "total": total})
}

// ReportController owns /v2/tqd/reports/*.
type ReportController struct {
	reports services.ReportUsecase
}

func NewReportController(reports services.ReportUsecase) *ReportController {
	return &ReportController{reports: reports}
}

func (c *ReportController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/tqd/reports", c.create)
	httpapi.GET(mux, "/v2/tqd/reports/{reportId}/status", c.status)
	httpapi.GET(mux, "/v2/tqd/reports/{reportId}/download", c.download)
	httpapi.GET(mux, "/v2/tqd/reports", c.list)
	httpapi.DELETE(mux, "/v2/tqd/reports/{reportId}", c.delete)
}

func (c *ReportController) create(ctx context.Context, req dto.CreateReportHTTPRequest) (*httpapi.Result, error) {
	if c.reports == nil {
		return nil, _errors.ReturnError(501, "report service not configured")
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	id, err := c.reports.CreateAsync(ctx, userID, dto.ReportDTO{
		ProblemReport: enums.ProblemReport(req.ProblemReport),
		ReportType:    enums.ReportType(req.ReportType),
		Description:   req.Description,
		TargetId:      req.TargetID,
	})
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"reportId": id})
}

func (c *ReportController) status(ctx context.Context, req dto.ReportIDHTTPRequest) (*httpapi.Result, error) {
	if c.reports == nil {
		return nil, _errors.ReturnError(501, "report service not configured")
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	report, err := c.reports.GetStatus(ctx, userID, req.ReportID)
	if err != nil {
		return nil, err
	}
	return jsonOK(report)
}

func (c *ReportController) download(ctx context.Context, req dto.ReportIDHTTPRequest) (*httpapi.Result, error) {
	if c.reports == nil {
		return nil, _errors.ReturnError(501, "report service not configured")
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	report, err := c.reports.Download(ctx, userID, req.ReportID)
	if err != nil {
		return nil, err
	}
	return jsonOK(report)
}

func (c *ReportController) list(ctx context.Context, req dto.ListReportsHTTPRequest) (*httpapi.Result, error) {
	if c.reports == nil {
		return nil, _errors.ReturnError(501, "report service not configured")
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.reports.ListUser(ctx, userID, req.ReportType, req.Status, int(pagable.GetPage()), int(pagable.GetSize()))
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"data": rows, "total": total, "page": pagable.GetPage(), "size": pagable.GetSize()})
}

func (c *ReportController) delete(ctx context.Context, req dto.ReportIDHTTPRequest) (*httpapi.Result, error) {
	if c.reports == nil {
		return nil, _errors.ReturnError(501, "report service not configured")
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if err := c.reports.Delete(ctx, userID, req.ReportID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

// PropertyController owns /v2/tqd/properties/*.
type PropertyController struct {
	oneHouse services.OneHouseUsecase
}

func NewPropertyController(oneHouse services.OneHouseUsecase) *PropertyController {
	return &PropertyController{oneHouse: oneHouse}
}

func (c *PropertyController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/properties/{id}", c.get)
}

func (c *PropertyController) get(ctx context.Context, req dto.PropertyIDHTTPRequest) (*httpapi.Result, error) {
	if c.oneHouse == nil {
		return nil, _errors.ReturnError(501, "property service not configured")
	}
	oh, err := c.oneHouse.GetDetail(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	proto, err := mapper.ToProtoOneHouse(oh)
	if err != nil {
		return nil, err
	}
	return protoResult(proto)
}
