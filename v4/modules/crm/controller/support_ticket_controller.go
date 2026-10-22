package controller

import (
	"context"
	"net/http"

	_errors "common/errors"
	"common/httpapi"
	"mdn4/modules/crm/dto"
	"mdn4/modules/crm/services"
)

type SupportTicketController struct {
	uc *services.SupportTicketUsecase
}

func NewSupportTicketController(uc *services.SupportTicketUsecase) *SupportTicketController {
	return &SupportTicketController{uc: uc}
}

func (c *SupportTicketController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/crm/admin/support/tickets", c.list)
	httpapi.GET(mux, "/v2/crm/admin/support/ticket-summary", c.summary)
	httpapi.GET(mux, "/v2/crm/admin/support/tickets/{id}", c.detail)
	httpapi.POST(mux, "/v2/crm/admin/support/tickets", c.create)
	httpapi.PATCH(mux, "/v2/crm/admin/support/tickets/{id}", c.update)
	httpapi.POST(mux, "/v2/crm/admin/support/tickets/{id}/assign", c.assign)
	httpapi.PATCH(mux, "/v2/crm/admin/support/tickets/{id}/priority", c.priority)
	httpapi.PATCH(mux, "/v2/crm/admin/support/tickets/{id}/status", c.status)
	httpapi.POST(mux, "/v2/crm/admin/support/tickets/{id}/notes", c.note)
	httpapi.POST(mux, "/v2/crm/admin/support/tickets/{id}/close", c.close)
	httpapi.POST(mux, "/v2/crm/admin/support/tickets/{id}/images", c.images)
	httpapi.POST(mux, "/v2/crm/admin/support/tickets/{id}/transfer", c.transfer)
}

func (c *SupportTicketController) require() error {
	if c == nil || c.uc == nil {
		return _errors.ReturnError(501, "support ticket service not configured")
	}
	return nil
}

type supportPath struct {
	ID uint64 `path:"id"`
}
type supportUpdateReq struct {
	ID uint64 `path:"id"`
	dto.SupportTicketUpdateRequest
}
type supportAssignReq struct {
	ID uint64 `path:"id"`
	dto.SupportTicketAssignRequest
}
type supportPriorityReq struct {
	ID uint64 `path:"id"`
	dto.SupportTicketPriorityRequest
}
type supportStatusReq struct {
	ID uint64 `path:"id"`
	dto.SupportTicketStatusRequest
}
type supportNoteReq struct {
	ID uint64 `path:"id"`
	dto.SupportTicketNoteRequest
}
type supportCloseReq struct {
	ID uint64 `path:"id"`
	dto.SupportTicketCloseRequest
}
type supportImagesReq struct {
	ID uint64 `path:"id"`
	dto.SupportTicketImagesRequest
}
type supportTransferReq struct {
	ID uint64 `path:"id"`
	dto.SupportTicketTransferRequest
}

func (c *SupportTicketController) list(ctx context.Context, req dto.SupportTicketListRequest) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.List(ctx, &req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}
func (c *SupportTicketController) summary(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.Summary(ctx)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}
func (c *SupportTicketController) detail(ctx context.Context, req supportPath) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}
func (c *SupportTicketController) create(ctx context.Context, req dto.SupportTicketCreateRequest) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.Create(ctx, &req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}
func (c *SupportTicketController) update(ctx context.Context, req supportUpdateReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.Update(ctx, req.ID, &req.SupportTicketUpdateRequest); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "updated")
}
func (c *SupportTicketController) assign(ctx context.Context, req supportAssignReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.Assign(ctx, req.ID, &req.SupportTicketAssignRequest); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "assigned")
}
func (c *SupportTicketController) priority(ctx context.Context, req supportPriorityReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.UpdatePriority(ctx, req.ID, &req.SupportTicketPriorityRequest); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "updated")
}
func (c *SupportTicketController) status(ctx context.Context, req supportStatusReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.UpdateStatus(ctx, req.ID, &req.SupportTicketStatusRequest); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "updated")
}
func (c *SupportTicketController) note(ctx context.Context, req supportNoteReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.AddNote(ctx, req.ID, &req.SupportTicketNoteRequest); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "noted")
}
func (c *SupportTicketController) close(ctx context.Context, req supportCloseReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.Close(ctx, req.ID, &req.SupportTicketCloseRequest); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "closed")
}
func (c *SupportTicketController) images(ctx context.Context, req supportImagesReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.UpdateImages(ctx, req.ID, req.Images); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "updated")
}
func (c *SupportTicketController) transfer(ctx context.Context, req supportTransferReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.Transfer(ctx, req.ID, &req.SupportTicketTransferRequest); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "transferred")
}
