package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	_errors "common/errors"
	bdspropb "pb/types/bdspro"
	"userv4/internal/httpapi"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/enums"
	"userv4/modules/bdspro/services"

	"google.golang.org/protobuf/proto"
)

// DealController owns deal HTTP routes under /v4 (migrated from /v2.1/...).
type DealController struct {
	deal       services.DealService
	commission services.DealCommissionService
}

func NewDealController(deal services.DealService, commission services.DealCommissionService) *DealController {
	return &DealController{deal: deal, commission: commission}
}

func (c *DealController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.HandleProto(mux, http.MethodPost, "/v4/deal", c.createProto)
	httpapi.HandleProto(mux, http.MethodPut, "/v4/deal/{id}", c.updateProto)
	httpapi.HandleProto(mux, http.MethodPut, "/v4/deal/{dealId}/setting", c.updateSettingProto)
	httpapi.HandleProto(mux, http.MethodDelete, "/v4/deal/{id}", c.deleteProto)
	httpapi.GET(mux, "/v4/deals/me", c.getDealsMe)
	httpapi.GET(mux, "/v4/deal/{id}/ts", c.getTimestamps)
	httpapi.GET(mux, "/v4/deal/{id}/products", c.getProducts)
	httpapi.POST(mux, "/v4/deal/{id}/products", c.saveProducts)
	httpapi.GET(mux, "/v4/deal/{id}/summary", c.getSummary)
	httpapi.GET(mux, "/v4/deal/{id}/bank-account", c.getBankAccount)
	httpapi.GET(mux, "/v4/deal/{id}/members", c.getMembers)
	httpapi.GET(mux, "/v4/deal/{id}/histories", c.getHistories)
	httpapi.POST(mux, "/v4/deal/{id}/commission", c.updateMembersCommission)
	httpapi.GET(mux, "/v4/deal/{id}", c.getDetail)
}

func (c *DealController) updateMembersCommission(ctx context.Context, req dto.UpdateDealMembersCommissionRequest) (*httpapi.Result, error) {
	if c.commission == nil {
		return nil, _errors.InternalServerException("deal commission service unavailable")
	}
	result, err := c.commission.UpdateMembersCommission(ctx, req)
	if err != nil {
		return nil, err
	}
	return protoResult(result)
}

func (c *DealController) createProto(ctx context.Context, req *bdspropb.CreateGroupDealRequest) (*bdspropb.CreateGroupDealResponse, error) {
	id, err := c.deal.CreateGroupDeal(ctx, req)
	if err != nil {
		return nil, err
	}
	return &bdspropb.CreateGroupDealResponse{Id: uint32(id)}, nil
}

func (c *DealController) updateProto(ctx context.Context, req *bdspropb.UpdateGroupDealRequest) (*bdspropb.UpdateGroupDealResponse, error) {
	id, err := c.deal.UpdateGroupDeal(ctx, req)
	if err != nil {
		return nil, err
	}
	return &bdspropb.UpdateGroupDealResponse{Id: uint32(id)}, nil
}

func (c *DealController) updateSettingProto(ctx context.Context, req *bdspropb.UpdateDealSettingRequest) (*bdspropb.UpdateDealSettingResponse, error) {
	dealID, err := c.deal.UpdateDealSetting(ctx, req)
	if err != nil {
		return nil, err
	}
	return &bdspropb.UpdateDealSettingResponse{DealId: dealID}, nil
}

func (c *DealController) deleteProto(ctx context.Context, req *bdspropb.DeleteGroupDealRequest) (*bdspropb.DeleteGroupDealResponse, error) {
	if err := c.deal.DeleteGroupDeal(ctx, uint64(req.GetId())); err != nil {
		return nil, err
	}
	return &bdspropb.DeleteGroupDealResponse{Id: req.GetId()}, nil
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

func (c *DealController) getProducts(ctx context.Context, req dto.DealListByIDRequest) (*httpapi.Result, error) {
	items, ts, err := c.deal.GetDealProducts(ctx, req)
	if err != nil {
		return nil, err
	}
	data := make([]any, 0, len(items))
	for _, item := range items {
		row := map[string]any{
			"id":        item.ID,
			"dealId":    item.DealID,
			"productId": item.ProductID,
		}
		if item.CreatedAt != "" {
			row["createdAt"] = item.CreatedAt
		}
		if item.UpdatedAt != "" {
			row["updatedAt"] = item.UpdatedAt
		}
		if item.CreatedBy != nil {
			row["createdBy"] = *item.CreatedBy
		}
		if item.UpdatedBy != nil {
			row["updatedBy"] = *item.UpdatedBy
		}
		if item.Product != nil {
			if msg, ok := item.Product.(proto.Message); ok {
				product, err := protoToAny(msg)
				if err != nil {
					return nil, err
				}
				row["product"] = product
			} else {
				row["product"] = item.Product
			}
		} else {
			row["product"] = map[string]any{}
		}
		data = append(data, row)
	}
	return httpapi.OK(data, timestampOpt(ts)...), nil
}

func (c *DealController) saveProducts(ctx context.Context, req dto.SaveDealProductsRequest) (*httpapi.Result, error) {
	ts, err := c.deal.SaveDealProducts(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(map[string]any{"status": 1}, timestampOpt(ts)...), nil
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

func (c *DealController) getHistories(ctx context.Context, req dto.DealListByIDRequest) (*httpapi.Result, error) {
	result, ts, err := c.deal.GetHistories(ctx, req)
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
	opts := []httpapi.OKOption{
		httpapi.WithPagination(&httpapi.Pagination{Total: int64(result.GetTotal())}),
	}
	opts = append(opts, timestampOpt(ts)...)
	return httpapi.OK(data, opts...), nil
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
	for i, row := range data {
		if m, ok := row.(map[string]any); ok {
			compactDealInvitationRole(m)
			data[i] = m
		}
	}
	opts := []httpapi.OKOption{
		httpapi.WithPagination(&httpapi.Pagination{Total: int64(resp.GetTotal())}),
	}
	opts = append(opts, timestampOpt(resp.GetTimestamp())...)
	return httpapi.OK(data, opts...), nil
}

// compactDealInvitationRole shapes invitation.role as {key, name} for deals/me.
func compactDealInvitationRole(deal map[string]any) {
	inv, ok := deal["invitation"].(map[string]any)
	if !ok || inv == nil {
		return
	}
	var roleKey uint32
	if role, ok := inv["role"].(map[string]any); ok && role != nil {
		roleKey = jsonUint32(role["roleKey"])
		if roleKey == 0 {
			roleKey = jsonUint32(role["key"])
		}
	}
	if roleKey == 0 {
		roleKey = jsonUint32(inv["roleKey"])
	}
	if roleKey == 0 {
		return
	}
	name := enums.RoleKeyMap[enums.RoleKey(roleKey)]
	inv["role"] = map[string]any{
		"key":  roleKey,
		"name": name,
	}
}

func jsonUint32(v any) uint32 {
	switch n := v.(type) {
	case float64:
		return uint32(n)
	case json.Number:
		u, _ := n.Int64()
		return uint32(u)
	case string:
		u, _ := strconv.ParseUint(n, 10, 32)
		return uint32(u)
	default:
		return 0
	}
}
