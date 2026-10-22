package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/services"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
)

// InvestmentController owns org investment routes (migrated from /v2.1/org/investment).
type InvestmentController struct {
	investment services.InvestmentService
}

func NewInvestmentController(investment services.InvestmentService) *InvestmentController {
	return &InvestmentController{investment: investment}
}

func (c *InvestmentController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v4/deal/{dealId}/investments", c.listHTTP)
	httpapi.HandleProto(mux, http.MethodPost, "/v4/deals/investment", c.createProto)
	httpapi.HandleProto(mux, http.MethodPut, "/v4/deals/investment/{id}", c.updateProto)
	httpapi.HandleProto(mux, http.MethodDelete, "/v4/deals/investment/{id}", c.deleteProto)
	httpapi.HandleProto(mux, http.MethodGet, "/v4/deals/investment/{id}", c.detailProto)
	httpapi.HandleProto(mux, http.MethodGet, "/v4/deals/investment/member-histories/{dealId}/{memberId}", c.historiesProto)
	httpapi.HandleProto(mux, http.MethodPost, "/v4/deals/investment/{id}/approve", c.approveProto)
	httpapi.HandleProto(mux, http.MethodPost, "/v4/deals/investment/{id}/reject", c.rejectProto)
	httpapi.HandleProto(mux, http.MethodPost, "/v4/deals/investment/{id}/confirm", c.confirmProto)
	httpapi.HandleProto(mux, http.MethodPost, "/v4/deals/investment/{id}/unconfirm", c.unconfirmProto)
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

func (c *InvestmentController) listHTTP(ctx context.Context, req dto.InvestmentListQuery) (*httpapi.Result, error) {
	result, total, ts, err := c.investment.List(ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return httpapi.OK([]any{}), nil
	}
	opts := []httpapi.OKOption{
		httpapi.WithPagination(&httpapi.Pagination{Total: total}),
	}
	opts = append(opts, httpapi.WithTimestamp(ts))
	return httpapi.OK(result, opts...), nil
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
