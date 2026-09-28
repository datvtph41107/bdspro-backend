package controller

import (
	"context"
	"net/http"

	bdspropb "pb/types/bdspro"
	"userv4/internal/httpapi"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/services"

	"google.golang.org/protobuf/proto"
)

// DealController owns deal HTTP routes under /v4 (migrated from /v2.1/...).
type DealController struct {
	deal services.DealService
}

func NewDealController(deal services.DealService) *DealController {
	return &DealController{deal: deal}
}

func (c *DealController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v4/deals/me", c.getDealsMe)
	httpapi.GET(mux, "/v4/deal/{id}/ts", c.getTimestamps)
	httpapi.GET(mux, "/v4/deal/{id}/products", c.getProducts)
	httpapi.GET(mux, "/v4/deal/{id}/summary", c.getSummary)
	httpapi.GET(mux, "/v4/deal/{id}/bank-account", c.getBankAccount)
	httpapi.GET(mux, "/v4/deal/{id}/members", c.getMembers)
	httpapi.GET(mux, "/v4/deal/{id}", c.getDetail)
	httpapi.GET(mux, "/v4/deal/invitation/list/{dealId}", c.getInvitations)
	httpapi.GET(mux, "/v4/deal/investments/list/{dealId}", c.getInvestments)
	httpapi.GET(mux, "/v4/transaction/contract/deal/{id}/list", c.getContracts)
	httpapi.GET(mux, "/v4/transaction/contract/deal/{id}/statistic", c.getContractStatistics)
	httpapi.GET(mux, "/v4/transaction/contract/deal/{id}/process", c.getContractProcess)
	httpapi.GET(mux, "/v4/transaction/deal-cost/list/{dealId}", c.getCosts)
}

func (c *DealController) getDealsMe(ctx context.Context, req dto.DealSearchRequest) (*httpapi.Result, error) {
	result, err := c.deal.GetDealsUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return dealsListResult(result)
}

func (c *DealController) getTimestamps(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	result, err := c.deal.GetTimestamps(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return httpapi.OK(map[string]int64{}), nil
	}
	return httpapi.OK(result.Timestamps), nil
}

func (c *DealController) getDetail(ctx context.Context, req dto.DealDetailRequest) (*httpapi.Result, error) {
	result, ts, err := c.deal.GetDetail(ctx, req)
	if err != nil {
		return nil, err
	}
	return protoResult(result, timestampOpt(ts)...)
}

func (c *DealController) getContractProcess(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	result, ts, err := c.deal.GetContractProcess(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return httpapi.OK([]any{}, timestampOpt(ts)...), nil
	}
	data, err := protoMessages(result.GetData())
	if err != nil {
		return nil, err
	}
	return httpapi.OK(data, timestampOpt(ts)...), nil
}

func (c *DealController) getProducts(ctx context.Context, req dto.DealListByIDRequest) (*httpapi.Result, error) {
	result, ts, err := c.deal.GetDealProducts(ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return httpapi.OK([]any{}, timestampOpt(ts)...), nil
	}
	data, err := protoMessages(result.GetData())
	if err != nil {
		return nil, err
	}
	return httpapi.OK(data, timestampOpt(ts)...), nil
}

func (c *DealController) getSummary(ctx context.Context, req dto.DealSummaryRequest) (*httpapi.Result, error) {
	result, ts, err := c.deal.GetSummary(ctx, req)
	if err != nil {
		return nil, err
	}
	return protoResult(result, timestampOpt(ts)...)
}

func (c *DealController) getBankAccount(ctx context.Context, req dto.DealDetailRequest) (*httpapi.Result, error) {
	result, ts, err := c.deal.GetBankAccount(ctx, req)
	if err != nil {
		return nil, err
	}
	return protoResult(result, timestampOpt(ts)...)
}

func (c *DealController) getMembers(ctx context.Context, req dto.DealMembersRequest) (*httpapi.Result, error) {
	result, ts, err := c.deal.GetMembers(ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return httpapi.OK([]any{}, timestampOpt(ts)...), nil
	}
	data, err := protoMessages(result.GetData())
	if err != nil {
		return nil, err
	}
	opts := timestampOpt(ts)
	opts = append(opts, httpapi.WithPagination(&httpapi.Pagination{Total: int64(result.GetTotal())}))
	return httpapi.OK(data, opts...), nil
}

func (c *DealController) getContracts(ctx context.Context, req dto.DealListByIDRequest) (*httpapi.Result, error) {
	result, err := c.deal.GetContracts(ctx, req)
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
	opts = append(opts, timestampOpt(result.GetTimestamp())...)
	return httpapi.OK(data, opts...), nil
}

func (c *DealController) getContractStatistics(ctx context.Context, req dto.DealStatisticsRequest) (*httpapi.Result, error) {
	result, ts, err := c.deal.GetContractStatistics(ctx, req)
	if err != nil {
		return nil, err
	}
	return protoResult(result, timestampOpt(ts)...)
}

func (c *DealController) getInvestments(ctx context.Context, req dto.DealListByIDRequest) (*httpapi.Result, error) {
	result, ts, err := c.deal.GetInvestments(ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return httpapi.OK([]any{}, timestampOpt(ts)...), nil
	}
	data, err := protoMessages(result.GetData())
	if err != nil {
		return nil, err
	}
	opts := timestampOpt(ts)
	opts = append(opts, httpapi.WithPagination(&httpapi.Pagination{Total: int64(result.GetTotal())}))
	return httpapi.OK(data, opts...), nil
}

func (c *DealController) getInvitations(ctx context.Context, req dto.DealListByIDRequest) (*httpapi.Result, error) {
	result, err := c.deal.GetInvitations(ctx, req)
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
		httpapi.WithPagination(&httpapi.Pagination{Total: result.GetTotal()}),
	}
	opts = append(opts, timestampOpt(result.GetTimestamp())...)
	return httpapi.OK(data, opts...), nil
}

func (c *DealController) getCosts(ctx context.Context, req dto.DealListByIDRequest) (*httpapi.Result, error) {
	result, err := c.deal.GetDealCosts(ctx, req)
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
		httpapi.WithPagination(&httpapi.Pagination{Total: int64(result.GetTotalElements())}),
	}
	opts = append(opts, timestampOpt(result.GetTimestamp())...)
	return httpapi.OK(data, opts...), nil
}

func timestampOpt(ts int64) []httpapi.OKOption {
	if ts <= 0 {
		return nil
	}
	return []httpapi.OKOption{httpapi.WithTimestamp(ts)}
}

func protoResult(msg proto.Message, opts ...httpapi.OKOption) (*httpapi.Result, error) {
	data, err := protoToAny(msg)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(data, opts...), nil
}

func dealsListResult(resp *bdspropb.GetDealsResponse) (*httpapi.Result, error) {
	if resp == nil {
		return httpapi.OK([]any{}), nil
	}
	data, err := protoMessages(resp.GetData())
	if err != nil {
		return nil, err
	}
	opts := []httpapi.OKOption{
		httpapi.WithPagination(&httpapi.Pagination{Total: int64(resp.GetTotal())}),
	}
	opts = append(opts, timestampOpt(resp.GetTimestamp())...)
	return httpapi.OK(data, opts...), nil
}
