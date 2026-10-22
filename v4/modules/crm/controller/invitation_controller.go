package controller

import (
	"context"
	"net/http"

	_errors "common/errors"
	"common/httpapi"
	"mdn4/modules/crm/dto"
	domain "mdn4/modules/crm/entity"
	"mdn4/modules/crm/services"
)

type InvitationController struct {
	uc *services.InvitationInstallUsecase
}

func NewInvitationController(uc *services.InvitationInstallUsecase) *InvitationController {
	return &InvitationController{uc: uc}
}

func (c *InvitationController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/crm/invitation/send", c.send)
	httpapi.GET(mux, "/v2/crm/invitation/list", c.list)
}

func (c *InvitationController) require() error {
	if c == nil || c.uc == nil {
		return _errors.ReturnError(501, "invitation service not configured")
	}
	return nil
}

type invitationSendReq struct {
	IDs []uint64 `json:"ids"`
}
type invitationListReq struct {
	Page uint32 `form:"page"`
	Size uint32 `form:"size"`
}

func (c *InvitationController) send(ctx context.Context, req invitationSendReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.SendInvitation(ctx, req.IDs, &domain.AppInvitedEntity{})
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}
func (c *InvitationController) list(ctx context.Context, req invitationListReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	items, total, err := c.uc.ListInvitation(ctx, &dto.ListInvitationRequest{})
	if err != nil {
		return nil, err
	}
	_ = req
	return httpapi.OK(items, httpapi.WithPagination(&httpapi.Pagination{Total: total})), nil
}
