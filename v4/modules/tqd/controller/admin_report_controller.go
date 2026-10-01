package controller

import (
	"context"
	"net/http"

	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	"userv4/modules/tqd/repository"
	"userv4/modules/tqd/services"
)

// AdminReportController owns /v2/tqd/admin/reports*.
type AdminReportController struct {
	uc services.ReportUsecase
}

func NewAdminReportController(uc services.ReportUsecase) *AdminReportController {
	return &AdminReportController{uc: uc}
}

func (c *AdminReportController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/admin/reports", c.list)
	httpapi.GET(mux, "/v2/tqd/admin/reports/summary", c.summary)
	httpapi.GET(mux, "/v2/tqd/admin/reports/queue", c.queue)
	httpapi.GET(mux, "/v2/tqd/admin/reports/queue/summary", c.queueSummary)
	httpapi.GET(mux, "/v2/tqd/admin/reports/events", c.listEvents)
	httpapi.GET(mux, "/v2/tqd/admin/reports/{id}/events", c.getEvent)
	httpapi.GET(mux, "/v2/tqd/admin/reports/{id}/detail", c.getDetail)
	// httpapi.GET(mux, "/v2/tqd/admin/reports/{id}/events", c.reportEvents)
	httpapi.POST(mux, "/v2/tqd/admin/reports/internal", c.createInternal)
	httpapi.POST(mux, "/v2/tqd/admin/reports/{id}/assign", c.assign)
	httpapi.POST(mux, "/v2/tqd/admin/reports/{id}/close", c.close)
	httpapi.POST(mux, "/v2/tqd/admin/reports/{id}/images", c.images)
	httpapi.POST(mux, "/v2/tqd/admin/reports/{id}/linkage", c.linkage)
	httpapi.PATCH(mux, "/v2/tqd/admin/reports/{id}/qa-status", c.qaStatus)
	httpapi.PATCH(mux, "/v2/tqd/admin/reports/{id}/severity", c.severity)
}

type adminReportListReq struct {
	Page   uint32  `form:"page"`
	Size   uint32  `form:"size"`
	Status *uint32 `form:"status"`
}

func (c *AdminReportController) list(ctx context.Context, req adminReportListReq) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.uc.AdminList(ctx, repository.AdminReportFilter{
		Status: req.Status, Page: int(pagable.GetPage()), Limit: int(pagable.GetSize()),
	})
	if err != nil {
		return nil, err
	}
	return jsonListResult(rows, total, pagable.GetPage(), pagable.GetSize())
}

func (c *AdminReportController) summary(ctx context.Context, _ dto.EmptyJSON) (*httpapi.Result, error) {
	out, err := c.uc.AdminSummary(ctx)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminReportController) queue(ctx context.Context, req struct {
	Page     uint32  `form:"page"`
	Size     uint32  `form:"size"`
	Bucket   string  `form:"bucket"`
	Severity *uint32 `form:"severity"`
}) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	actorID := _utils.GetProfileIdWithContext(ctx)
	rows, total, err := c.uc.AdminQueue(ctx, req.Bucket, req.Severity, int(pagable.GetPage()), int(pagable.GetSize()), actorID)
	if err != nil {
		return nil, err
	}
	return jsonListResult(rows, total, pagable.GetPage(), pagable.GetSize())
}

func (c *AdminReportController) queueSummary(ctx context.Context, _ dto.EmptyJSON) (*httpapi.Result, error) {
	actorID := _utils.GetProfileIdWithContext(ctx)
	out, err := c.uc.AdminQueueSummary(ctx, actorID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminReportController) getDetail(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.uc.AdminGet(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminReportController) listEvents(ctx context.Context, req adminReportListReq) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.uc.AdminListEvents(ctx, repository.AdminEventFilter{Page: int(pagable.GetPage()), Limit: int(pagable.GetSize())})
	if err != nil {
		return nil, err
	}
	return jsonListResult(rows, total, pagable.GetPage(), pagable.GetSize())
}

func (c *AdminReportController) getEvent(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	rows, _, err := c.uc.AdminListEvents(ctx, repository.AdminEventFilter{Page: 1, Limit: 1})
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		if row.ID == req.ID {
			return jsonOK(row)
		}
	}
	return nil, _errors.ReturnError(404, "event not found")
}

func (c *AdminReportController) reportEvents(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	rows, err := c.uc.AdminListReportEvents(ctx, req.ID, 100)
	if err != nil {
		return nil, err
	}
	return jsonListResult(rows, int64(len(rows)), 1, uint32(len(rows)))
}

func (c *AdminReportController) createInternal(ctx context.Context, req dto.AdminCreateReportRequest) (*httpapi.Result, error) {
	actorID := _utils.GetProfileIdWithContext(ctx)
	out, err := c.uc.AdminCreateInternal(ctx, actorID, req)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminReportController) assign(ctx context.Context, req dto.AdminAssignReportRequest) (*httpapi.Result, error) {
	if err := c.uc.AdminAssign(ctx, req.ID, req.AssigneeID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminReportController) close(ctx context.Context, req struct {
	ID     uint64 `path:"id" json:"-"`
	Reject bool   `json:"reject"`
	Note   string `json:"note"`
}) (*httpapi.Result, error) {
	if err := c.uc.AdminClose(ctx, req.ID, req.Reject, req.Note); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminReportController) images(ctx context.Context, req struct {
	ID     uint64   `path:"id" json:"-"`
	Images []string `json:"images"`
}) (*httpapi.Result, error) {
	if err := c.uc.AdminUpdateImages(ctx, req.ID, req.Images); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminReportController) linkage(ctx context.Context, req struct {
	ID      uint64                 `path:"id" json:"-"`
	Linkage map[string]interface{} `json:"linkage"`
}) (*httpapi.Result, error) {
	if err := c.uc.AdminUpdateLinkage(ctx, req.ID, req.Linkage); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminReportController) qaStatus(ctx context.Context, req dto.AdminQaStatusRequest) (*httpapi.Result, error) {
	if err := c.uc.AdminUpdateQaStatus(ctx, req.ID, req.QaStatus, req.Note); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminReportController) severity(ctx context.Context, req dto.AdminSeverityRequest) (*httpapi.Result, error) {
	if err := c.uc.AdminUpdateSeverity(ctx, req.ID, req.Severity, req.Reason); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}
