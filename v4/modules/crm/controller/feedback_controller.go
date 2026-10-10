package controller

import (
	"context"
	"net/http"

	_dto "common/domain/dto"
	_utils "common/utils"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	"userv4/internal/httpapi"
	"userv4/modules/crm/dto"
	domain "userv4/modules/crm/entity"
	"userv4/modules/crm/mapper"
	"userv4/modules/crm/services"
)

// FeedbackController owns /v2/feedback/* (rate, report, report-reason, registry),
// ported from crm-service Rate/Report/ReportAdmin/ReportReason/Registry handlers.
type FeedbackController struct {
	rate         *services.RateUsecase
	report       *services.ReportUsecase
	reportReason *services.ReportReasonUsecase
	registry     *services.RegistryUsecase
	rateMapper   *mapper.RateMapper
	reportMapper *mapper.ReportMapper
}

func NewFeedbackController(
	rate *services.RateUsecase,
	report *services.ReportUsecase,
	reportReason *services.ReportReasonUsecase,
	registry *services.RegistryUsecase,
	rateMapper *mapper.RateMapper,
	reportMapper *mapper.ReportMapper,
) *FeedbackController {
	return &FeedbackController{
		rate:         rate,
		report:       report,
		reportReason: reportReason,
		registry:     registry,
		rateMapper:   rateMapper,
		reportMapper: reportMapper,
	}
}

// RegisterHTTP mounts the paths from shared/protobuf/schema/crm/{rate,report,
// report_admin,report_reason,registry}.proto.
func (c *FeedbackController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.HandleProto(mux, http.MethodPost, "/v2/feedback/rate/{registry}", c.CreateRate)
	httpapi.HandleProto(mux, http.MethodGet, "/v2/feedback/rate/{registry}/list/{ownerId}", c.GetRateList)
	httpapi.HandleProto(mux, http.MethodGet, "/v2/feedback/rate/{registry}/stats/{ownerId}", c.GetRateStats)
	httpapi.HandleProto(mux, http.MethodPut, "/v2/feedback/rate/{id}", c.UpdateRate)
	httpapi.HandleProto(mux, http.MethodDelete, "/v2/feedback/rate/{id}", c.DeleteRate)
	httpapi.HandleProto(mux, http.MethodGet, "/v2/feedback/rate/detail/{id}", c.GetRate)
	httpapi.HandleProto(mux, http.MethodPut, "/v2/feedback/rate/{registry}/hidden/{id}", c.UpdateHidden)

	httpapi.HandleProto(mux, http.MethodGet, "/v2/feedback/registry", c.ListRegistries)

	httpapi.HandleProto(mux, http.MethodPost, "/v2/feedback/report/{registry}", c.CreateReport)
	httpapi.HandleProto(mux, http.MethodGet, "/v2/feedback/report/{registry}/list/{ownerId}", c.GetReportList)
	httpapi.HandleProto(mux, http.MethodGet, "/v2/feedback/report/reasons", c.GetReportReasons)
	httpapi.HandleProto(mux, http.MethodPut, "/v2/feedback/report/{id}", c.UpdateReport)
	httpapi.HandleProto(mux, http.MethodDelete, "/v2/feedback/report/{id}", c.DeleteReport)

	httpapi.HandleProto(mux, http.MethodGet, "/v2/feedback/admin/report/{registry}", c.GetAdminReportList)
	httpapi.HandleProto(mux, http.MethodPost, "/v2/feedback/admin/report/{id}/approve", c.ApproveReport)
	httpapi.HandleProto(mux, http.MethodPost, "/v2/feedback/admin/report/{id}/reject", c.RejectReport)
	httpapi.HandleProto(mux, http.MethodDelete, "/v2/feedback/admin/report/{id}/remove", c.RemoveReport)

	httpapi.HandleProto(mux, http.MethodPost, "/v2/feedback/report-reason", c.CreateReportReason)
	httpapi.HandleProto(mux, http.MethodGet, "/v2/feedback/report-reason", c.GetReportReasonList)
	httpapi.HandleProto(mux, http.MethodGet, "/v2/feedback/report-reason/{id}", c.GetReportReason)
	httpapi.HandleProto(mux, http.MethodPut, "/v2/feedback/report-reason/{id}", c.UpdateReportReason)
	httpapi.HandleProto(mux, http.MethodDelete, "/v2/feedback/report-reason/{id}", c.DeleteReportReason)
}

// ── Rate ─────────────────────────────────────────────────────────────────────

func (c *FeedbackController) CreateRate(ctx context.Context, req *crmpb.CreateRateRequest) (*crmpb.RateResponse, error) {
	rate := rateCreateRequestToDomain(&dto.RateCreateRequest{
		Anonymous: req.Anonymous,
		OwnerID:   req.OwnerId,
		OwnerOf:   req.OwnerOf,
		Score:     uint8(req.Score),
		Comment:   req.Comment,
		ParentID:  req.ParentId,
		Attachs:   req.Attachs,
	})
	created, err := c.rate.CreateRate(ctx, req.Registry, rate)
	if err != nil {
		return nil, err
	}
	resp := c.rateMapper.MapRateToPb(ctx, created)
	c.rateMapper.MapUserToPbSingle(ctx, resp)
	return resp, nil
}

func (c *FeedbackController) GetRateList(ctx context.Context, req *crmpb.GetRateListRequest) (*crmpb.GetRateListResponse, error) {
	rates, total, err := c.rate.GetRateList(ctx, req.OwnerId, &dto.RateSearchDTO{
		Pagable:  _dto.Pagable{Page: req.Page, Size: req.Size},
		OwnerId:  req.OwnerId,
		ParentId: req.ParentId,
		Registry: req.Registry,
	})
	if err != nil {
		return nil, err
	}
	return &crmpb.GetRateListResponse{Data: c.rateMapper.MapRateToPbs(ctx, rates), Total: total}, nil
}

func (c *FeedbackController) GetRateStats(ctx context.Context, req *crmpb.GetRateStatsRequest) (*crmpb.RateStatsResponse, error) {
	stats, err := c.rate.GetRateStats(ctx, req.OwnerId, req.Registry)
	if err != nil {
		return nil, err
	}
	return &crmpb.RateStatsResponse{
		TotalReviews: stats.TotalReviews,
		AverageScore: stats.AverageScore,
		Star1Count:   stats.Star1Count,
		Star2Count:   stats.Star2Count,
		Star3Count:   stats.Star3Count,
		Star4Count:   stats.Star4Count,
		Star5Count:   stats.Star5Count,
	}, nil
}

func (c *FeedbackController) UpdateRate(ctx context.Context, req *crmpb.UpdateRateRequest) (*crmpb.UpdateRateResponse, error) {
	rate := rateUpdateRequestToDomain(&dto.RateUpdateRequest{
		Score:   uint8(req.Score),
		Comment: req.Comment,
		Attachs: req.Attachs,
	})
	if err := c.rate.UpdateRate(ctx, req.Id, rate); err != nil {
		return nil, err
	}
	return &crmpb.UpdateRateResponse{Success: true}, nil
}

func (c *FeedbackController) DeleteRate(ctx context.Context, req *crmpb.RegistryIdRequest) (*sharepb.SubmitResponse, error) {
	if err := c.rate.DeleteRate(ctx, req.Id); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "Rate deleted successfully"}, nil
}

func (c *FeedbackController) GetRate(ctx context.Context, req *crmpb.RegistryIdRequest) (*crmpb.RateResponse, error) {
	rate, err := c.rate.GetRate(ctx, req.Id, req.Registry)
	if err != nil {
		return nil, err
	}
	resp := c.rateMapper.MapRateToPb(ctx, rate)
	c.rateMapper.MapUserToPbSingle(ctx, resp)
	return resp, nil
}

func (c *FeedbackController) UpdateHidden(ctx context.Context, req *crmpb.UpdateHiddenRequest) (*sharepb.SubmitResponse, error) {
	if err := c.rate.UpdateHidden(ctx, req.Id, req.Hidden, req.Registry); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "Rate hidden status updated successfully"}, nil
}

func rateCreateRequestToDomain(req *dto.RateCreateRequest) *domain.Rate {
	rate := &domain.Rate{
		Anonymous: req.Anonymous,
		OwnerID:   req.OwnerID,
		OwnerOf:   req.OwnerOf,
		Score:     req.Score,
		Comment:   req.Comment,
		ParentID:  req.ParentID,
	}
	rate.Attachs = imageAttachs(req.Attachs)
	return rate
}

func rateUpdateRequestToDomain(req *dto.RateUpdateRequest) *domain.Rate {
	rate := &domain.Rate{Score: req.Score, Comment: req.Comment}
	// A non-nil empty list clears all attachments.
	if req.Attachs != nil {
		rate.Attachs = imageAttachs(req.Attachs)
		if rate.Attachs == nil {
			rate.Attachs = []domain.FeedbackAttachEntity{}
		}
	}
	return rate
}

func imageAttachs(urls []string) []domain.FeedbackAttachEntity {
	var out []domain.FeedbackAttachEntity
	for _, u := range urls {
		out = append(out, domain.FeedbackAttachEntity{FileURL: u, FileType: "image", FileName: "attachment"})
	}
	return out
}

// ── Registry ─────────────────────────────────────────────────────────────────

func (c *FeedbackController) ListRegistries(ctx context.Context, _ *sharepb.Empty) (*crmpb.GetRegistryListResponse, error) {
	registries, err := c.registry.GetAllRegistries(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*crmpb.FeedbackRegistry, len(registries))
	for i, r := range registries {
		out[i] = &crmpb.FeedbackRegistry{
			RegistryKey:  r.RegistryKey,
			RegistryName: r.RegistryName,
			Port:         int32(r.Port),
			Service:      r.Service,
		}
	}
	return &crmpb.GetRegistryListResponse{Data: out}, nil
}

// ── Report ───────────────────────────────────────────────────────────────────

func (c *FeedbackController) CreateReport(ctx context.Context, req *crmpb.CreateReportRequest) (*crmpb.ReportResponse, error) {
	created, err := c.report.CreateReport(ctx, &dto.ReportCreateRequest{
		ReasonID:     req.ReasonId,
		UserID:       req.UserId,
		OwnerID:      req.OwnerId,
		OwnerOf:      req.OwnerOf,
		RegistryName: req.Registry,
		Content:      req.Content,
		ProofDocs:    req.ProofDocs,
	})
	if err != nil {
		return nil, err
	}
	return c.reportMapper.ReportToPb(created), nil
}

func (c *FeedbackController) GetReportList(ctx context.Context, req *crmpb.GetReportListRequest) (*crmpb.GetReportListResponse, error) {
	ownerID := req.OwnerId
	reports, total, err := c.report.GetReportList(ctx, &dto.ReportListRequest{
		Pagable:  _dto.Pagable{Page: req.Page, Size: req.Size},
		OwnerId:  &ownerID,
		Status:   req.Status,
		ReasonID: req.ReasonId,
		Registry: req.Registry,
	})
	if err != nil {
		return nil, err
	}
	return &crmpb.GetReportListResponse{Data: c.reportsToPb(reports), Total: total}, nil
}

func (c *FeedbackController) GetReportReasons(ctx context.Context, _ *crmpb.GetReportReasonsRequest) (*crmpb.GetReportReasonsResponse, error) {
	reasons, err := c.report.GetReportReasons(ctx)
	if err != nil {
		return nil, err
	}
	return &crmpb.GetReportReasonsResponse{
		Data:  c.reportMapper.ReportReasonListToReportReasonResponseList(reasons),
		Total: int64(len(reasons)),
	}, nil
}

func (c *FeedbackController) UpdateReport(ctx context.Context, req *crmpb.UpdateReportRequest) (*crmpb.UpdateReportResponse, error) {
	updated, err := c.report.UpdateReport(ctx, req.Id, &dto.ReportUpdateRequest{ProofDocs: req.ProofDocs})
	if err != nil {
		return nil, err
	}
	return &crmpb.UpdateReportResponse{Report: c.reportMapper.ReportToPb(updated)}, nil
}

func (c *FeedbackController) DeleteReport(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if err := c.report.DeleteReport(ctx, req.Id); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "Report deleted successfully"}, nil
}

// ── Report (admin) ───────────────────────────────────────────────────────────

func (c *FeedbackController) GetAdminReportList(ctx context.Context, req *crmpb.GetAdminReportListRequest) (*crmpb.GetAdminReportListResponse, error) {
	search := &dto.ReportListRequest{
		Pagable:  _dto.Pagable{Page: req.Page, Size: req.Size},
		Registry: req.Registry,
		Status:   req.Status,
		ReasonID: req.ReasonId,
	}
	if req.UserId != nil {
		id := *req.UserId
		search.OwnerId = &id
	}
	if req.OwnerId != nil {
		id := *req.OwnerId
		search.OwnerId = &id
	}
	reports, total, err := c.report.GetAdminReportList(ctx, search)
	if err != nil {
		return nil, err
	}
	return &crmpb.GetAdminReportListResponse{Data: c.reportsToPb(reports), Total: total}, nil
}

func (c *FeedbackController) ApproveReport(ctx context.Context, req *crmpb.ApproveReportRequest) (*sharepb.SubmitResponse, error) {
	if err := c.report.ApproveReport(ctx, req.Id, adminReportAction(req.Response, req.AdminNote, req.SendNotification)); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "Báo cáo đã được duyệt"}, nil
}

func (c *FeedbackController) RejectReport(ctx context.Context, req *crmpb.RejectReportRequest) (*sharepb.SubmitResponse, error) {
	if err := c.report.RejectReport(ctx, req.Id, adminReportAction(req.Response, req.AdminNote, req.SendNotification)); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "Báo cáo đã bị từ chối"}, nil
}

func (c *FeedbackController) RemoveReport(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if err := c.report.RemoveReport(ctx, req.Id); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "Báo cáo đã được gỡ khỏi hệ thống"}, nil
}

func adminReportAction(response, adminNote *string, sendNotification *bool) *dto.AdminReportActionRequest {
	action := &dto.AdminReportActionRequest{Response: response, AdminNote: adminNote}
	if sendNotification != nil {
		action.SendNotification = *sendNotification
	}
	return action
}

func (c *FeedbackController) reportsToPb(reports []domain.Report) []*crmpb.ReportResponse {
	out := make([]*crmpb.ReportResponse, 0, len(reports))
	for i := range reports {
		out = append(out, c.reportMapper.ReportToPb(&reports[i]))
	}
	return out
}

// ── Report reason ────────────────────────────────────────────────────────────

func (c *FeedbackController) CreateReportReason(ctx context.Context, req *crmpb.CreateReportReasonRequest) (*crmpb.CreateReportReasonResponse, error) {
	resp, err := c.reportReason.CreateReportReason(ctx, &dto.ReportReasonCreateRequest{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	})
	if err != nil {
		return nil, err
	}
	return &crmpb.CreateReportReasonResponse{Reason: reportReasonToPb(resp)}, nil
}

func (c *FeedbackController) GetReportReasonList(ctx context.Context, req *crmpb.GetReportReasonListRequest) (*crmpb.GetReportReasonListResponse, error) {
	resp, err := c.reportReason.GetReportReasonList(ctx, &dto.ReportReasonListRequest{
		Pagable:  _dto.Pagable{Page: req.Page, Size: req.Size},
		IsActive: req.IsActive,
	})
	if err != nil {
		return nil, err
	}
	reasons := make([]*crmpb.ReportReason, len(resp.Data))
	for i := range resp.Data {
		reasons[i] = reportReasonToPb(&resp.Data[i])
	}
	return &crmpb.GetReportReasonListResponse{Data: reasons, Total: resp.Total}, nil
}

func (c *FeedbackController) GetReportReason(ctx context.Context, req *sharepb.IdRequest) (*crmpb.GetReportReasonResponse, error) {
	resp, err := c.reportReason.GetReportReasonByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &crmpb.GetReportReasonResponse{Reason: reportReasonToPb(resp)}, nil
}

func (c *FeedbackController) UpdateReportReason(ctx context.Context, req *crmpb.UpdateReportReasonRequest) (*crmpb.UpdateReportReasonResponse, error) {
	resp, err := c.reportReason.UpdateReportReason(ctx, &dto.ReportReasonUpdateRequest{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	})
	if err != nil {
		return nil, err
	}
	return &crmpb.UpdateReportReasonResponse{Reason: reportReasonToPb(resp)}, nil
}

func (c *FeedbackController) DeleteReportReason(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if err := c.reportReason.DeleteReportReason(ctx, req.Id); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "Report reason deleted successfully"}, nil
}

func reportReasonToPb(reason *dto.ReportReasonResponse) *crmpb.ReportReason {
	if reason == nil {
		return nil
	}
	return &crmpb.ReportReason{
		Id:          reason.ID,
		Name:        reason.Name,
		Description: reason.Description,
		IsActive:    reason.IsActive,
		CreatedAt:   _utils.FormatTimeToString(reason.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(reason.UpdatedAt),
	}
}
