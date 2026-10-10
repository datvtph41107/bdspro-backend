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

// InvestmentController owns org investment routes (migrated from /v2.1/org/investment).
type InvestmentController struct {
	investment services.InvestmentService
}

func NewInvestmentController(investment services.InvestmentService) *InvestmentController {
	return &InvestmentController{investment: investment}
}

func (c *InvestmentController) RegisterHTTP(mux *http.ServeMux) {
	registerInvestmentRoutes(mux, "/v2.1/org", c)
	registerInvestmentRoutes(mux, "/v4/org", c)
	httpapi.GET(mux, "/v4/deal/investments/list/{dealId}", c.listHTTP)
}

func registerInvestmentRoutes(mux *http.ServeMux, orgBase string, c *InvestmentController) {
	base := orgBase + "/investment"
	httpapi.HandleProto(mux, http.MethodPost, base, c.createProto)
	httpapi.HandleProto(mux, http.MethodPut, base+"/{id}", c.updateProto)
	httpapi.HandleProto(mux, http.MethodDelete, base+"/{id}", c.deleteProto)
	httpapi.HandleProto(mux, http.MethodGet, base+"/{id}", c.detailProto)
	httpapi.HandleProto(mux, http.MethodGet, base+"/member-histories/{dealId}/{memberId}", c.historiesProto)
	httpapi.HandleProto(mux, http.MethodPost, base+"/{id}/approve", c.approveProto)
	httpapi.HandleProto(mux, http.MethodPost, base+"/{id}/reject", c.rejectProto)
	httpapi.HandleProto(mux, http.MethodPost, base+"/{id}/confirm", c.confirmProto)
	httpapi.HandleProto(mux, http.MethodPost, base+"/{id}/unconfirm", c.unconfirmProto)
	listBase := orgBase + "/investments/list/{dealId}"
	httpapi.HandleProto(mux, http.MethodGet, listBase, c.listProto)
}

func (c *InvestmentController) createProto(ctx context.Context, req *bdspropb.SubmitInvestRequest) (*sharepb.IdDTO, error) {
	id, err := c.investment.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &sharepb.IdDTO{Id: id}, nil
}

func (c *InvestmentController) updateProto(ctx context.Context, req *bdspropb.SubmitInvestRequest) (*sharepb.IdDTO, error) {
	id, err := c.investment.Update(ctx, req)
	if err != nil {
		return nil, err
	}
	return &sharepb.IdDTO{Id: id}, nil
}

func (c *InvestmentController) deleteProto(ctx context.Context, req *sharepb.IdDTO) (*sharepb.IdDTO, error) {
	if err := c.investment.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return req, nil
}

func (c *InvestmentController) detailProto(ctx context.Context, req *sharepb.IdDTO) (*bdspropb.Investment, error) {
	return c.investment.Detail(ctx, req.GetId())
}

func (c *InvestmentController) listProto(ctx context.Context, req *bdspropb.InvestmentListRequest) (*bdspropb.ListInvestment, error) {
	return c.investment.List(ctx, req)
}

func (c *InvestmentController) listHTTP(ctx context.Context, req dto.InvestmentListQuery) (*httpapi.Result, error) {
	protoReq := &bdspropb.InvestmentListRequest{
		DealId:     req.DealID,
		InvestType: req.InvestType,
		Status:     req.Status,
		Timestamp:  req.Timestamp,
	}
	if req.MemberID != nil {
		protoReq.MemberId = req.MemberID
	}
	result, err := c.investment.List(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return httpapi.OK([]any{}), nil
	}
	data, err := protoMessages(result.GetData())
	if err != nil {
		return nil, err
	}
	opts := []httpapi.OKOption{
		httpapi.WithPagination(&httpapi.Pagination{Total: int64(result.GetTotal())}),
	}
	return httpapi.OK(data, opts...), nil
}

func (c *InvestmentController) approveProto(ctx context.Context, req *sharepb.IdDTO) (*sharepb.IdDTO, error) {
	id, err := c.investment.Approve(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &sharepb.IdDTO{Id: id}, nil
}

func (c *InvestmentController) rejectProto(ctx context.Context, req *sharepb.IdDTO) (*sharepb.IdDTO, error) {
	id, err := c.investment.Reject(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &sharepb.IdDTO{Id: id}, nil
}

func (c *InvestmentController) confirmProto(ctx context.Context, req *sharepb.IdDTO) (*sharepb.IdDTO, error) {
	id, err := c.investment.Confirm(ctx, req.GetId(), true)
	if err != nil {
		return nil, err
	}
	return &sharepb.IdDTO{Id: id}, nil
}

func (c *InvestmentController) unconfirmProto(ctx context.Context, req *sharepb.IdDTO) (*sharepb.IdDTO, error) {
	id, err := c.investment.Confirm(ctx, req.GetId(), false)
	if err != nil {
		return nil, err
	}
	return &sharepb.IdDTO{Id: id}, nil
}

func (c *InvestmentController) historiesProto(ctx context.Context, req *bdspropb.GetInvestmentHistoriesRequest) (*bdspropb.GetInvestmentHistoriesResponse, error) {
	return c.investment.MemberHistories(ctx, req.GetDealId(), req.GetMemberId())
}
