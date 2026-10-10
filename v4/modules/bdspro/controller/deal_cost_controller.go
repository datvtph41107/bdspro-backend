package controller

import (
	"context"
	"net/http"

	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
	"userv4/internal/httpapi"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/services"

	_dto "common/domain/dto"
)

// DealCostController owns /v2.1 and /v4 transaction/deal-cost routes.
type DealCostController struct {
	dealCost services.DealCostService
}

func NewDealCostController(dealCost services.DealCostService) *DealCostController {
	return &DealCostController{dealCost: dealCost}
}

func (c *DealCostController) RegisterHTTP(mux *http.ServeMux) {
	registerDealCostRoutes(mux, "/v2.1/transaction/deal-cost", c)
	registerDealCostRoutes(mux, "/v4/transaction/deal-cost", c)
}

func registerDealCostRoutes(mux *http.ServeMux, base string, c *DealCostController) {
	httpapi.HandleProto(mux, http.MethodGet, base+"/list/{dealId}", c.searchProto)
	httpapi.HandleProto(mux, http.MethodGet, base+"/detail/{id}", c.detailProto)
	httpapi.HandleProto(mux, http.MethodPost, base, c.createProto)
	httpapi.HandleProto(mux, http.MethodPut, base+"/{id}", c.updateProto)
	httpapi.HandleProto(mux, http.MethodDelete, base+"/{id}", c.deleteProto)
}

func (c *DealCostController) searchProto(ctx context.Context, req *bdspropb.DealCostSearchDTO) (*bdspropb.DealCostListDTO, error) {
	search := dto.DealCostSearchRequest{
		Pagable:   _dto.Pagable{Page: req.GetPage(), Size: req.GetSize()},
		DealID:    req.GetDealId(),
		Timestamp: req.GetTimestamp(),
	}
	return c.dealCost.Search(ctx, search)
}

func (c *DealCostController) detailProto(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.DealCostDTO, error) {
	return c.dealCost.Detail(ctx, req.GetId())
}

func (c *DealCostController) createProto(ctx context.Context, req *bdspropb.DealCostDTO) (*bdspropb.DealCostDTO, error) {
	return c.dealCost.Create(ctx, req)
}

func (c *DealCostController) updateProto(ctx context.Context, req *bdspropb.DealCostDTO) (*bdspropb.DealCostDTO, error) {
	return c.dealCost.Update(ctx, req)
}

func (c *DealCostController) deleteProto(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Empty, error) {
	if err := c.dealCost.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &sharepb.Empty{}, nil
}
