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

// CostTypeController owns transaction cost-type routes (migrated from /v2.1/transaction/cost-type).
type CostTypeController struct {
	costType services.CostTypeService
}

func NewCostTypeController(costType services.CostTypeService) *CostTypeController {
	return &CostTypeController{costType: costType}
}

func (c *CostTypeController) RegisterHTTP(mux *http.ServeMux) {
	registerCostTypeRoutes(mux, "/v2.1/transaction/cost-type", c)
	registerCostTypeRoutes(mux, "/v4/transaction/cost-type", c)
}

func registerCostTypeRoutes(mux *http.ServeMux, base string, c *CostTypeController) {
	httpapi.HandleProto(mux, http.MethodGet, base, c.listProto)
	httpapi.HandleProto(mux, http.MethodPost, base, c.createProto)
	httpapi.HandleProto(mux, http.MethodPut, base+"/{id}", c.updateProto)
	httpapi.HandleProto(mux, http.MethodDelete, base+"/{id}", c.deleteProto)
}

func (c *CostTypeController) listProto(ctx context.Context, req *bdspropb.CostTypeSearchDTO) (*bdspropb.CostTypeListDTO, error) {
	search := dto.CostTypeSearchRequest{
		TypeName: req.GetTypeName(),
		CostType: req.GetCostType(),
	}
	search.Page = req.GetPage()
	search.Size = req.GetSize()
	return c.costType.List(ctx, search)
}

func (c *CostTypeController) createProto(ctx context.Context, req *bdspropb.CostTypeDTO) (*bdspropb.CostTypeDTO, error) {
	return c.costType.Create(ctx, req)
}

func (c *CostTypeController) updateProto(ctx context.Context, req *bdspropb.CostTypeDTO) (*bdspropb.CostTypeDTO, error) {
	return c.costType.Update(ctx, req)
}

func (c *CostTypeController) deleteProto(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Empty, error) {
	if err := c.costType.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &sharepb.Empty{}, nil
}
