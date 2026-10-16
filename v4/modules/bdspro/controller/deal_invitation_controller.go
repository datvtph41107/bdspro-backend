package controller

import (
	"context"
	"net/http"

	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
	"userv4/internal/httpapi"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/services"
)

// DealInvitationController owns /v4/deals/{id}/invitation* routes.
type DealInvitationController struct {
	invitations services.DealInvitationService
}

func NewDealInvitationController(invitations services.DealInvitationService) *DealInvitationController {
	return &DealInvitationController{invitations: invitations}
}

func (c *DealInvitationController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v4/deal/{id}/invitations", c.list)
	httpapi.POST(mux, "/v4/deal/{id}/invitation", c.send)
	httpapi.POST(mux, "/v4/deal/{id}/invitation/accept", c.accept)
	httpapi.POST(mux, "/v4/deal/{id}/invitation/reject", c.reject)
	httpapi.POST(mux, "/v4/deals/invitation/remove", c.remove)
	httpapi.PUT(mux, "/v4/deals/invitation/role", c.updateRole)
}

func (c *DealInvitationController) list(ctx context.Context, req dto.DealInvitationListRequest) (*httpapi.Result, error) {
	result, err := c.invitations.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return dealInvitationsListResult(result)
}

func (c *DealInvitationController) send(ctx context.Context, req dto.SendDealInvitationRequest) (*httpapi.Result, error) {
	result, err := c.invitations.Send(ctx, req)
	if err != nil {
		return nil, err
	}
	return protoResult(result)
}

func (c *DealInvitationController) accept(ctx context.Context, req dto.AcceptDealInvitationRequest) (*httpapi.Result, error) {
	result, err := c.invitations.Accept(ctx, req)
	if err != nil {
		return nil, err
	}
	return protoResult(result)
}

func (c *DealInvitationController) reject(ctx context.Context, req dto.RejectDealInvitationRequest) (*httpapi.Result, error) {
	result, err := c.invitations.Reject(ctx, req)
	if err != nil {
		return nil, err
	}
	return protoResult(result)
}

func (c *DealInvitationController) remove(ctx context.Context, req dto.RemoveFromDealRequest) (*httpapi.Result, error) {
	result, err := c.invitations.Remove(ctx, req)
	if err != nil {
		return nil, err
	}
	return protoResult(result)
}

func (c *DealInvitationController) updateRole(ctx context.Context, req dto.UpdateDealMemberRoleRequest) (*httpapi.Result, error) {
	result, err := c.invitations.UpdateMemberRole(ctx, req)
	if err != nil {
		return nil, err
	}
	return protoResult(result)
}

func dealInvitationsListResult(result *bdspropb.GetDealInvitationsByDealResponse) (*httpapi.Result, error) {
	ts := int64(0)
	total := int64(0)
	var rows []*sharepb.DealMember
	if result != nil {
		ts = result.GetTimestamp()
		total = result.GetTotal()
		rows = result.GetData()
	}
	data, err := protoMessages(rows)
	if err != nil {
		return nil, err
	}
	opts := []httpapi.OKOption{
		httpapi.WithPagination(&httpapi.Pagination{Total: total}),
	}
	opts = append(opts, timestampOpt(ts)...)
	return httpapi.OK(data, opts...), nil
}
