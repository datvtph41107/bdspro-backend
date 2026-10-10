package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	bdspropb "pb/types/bdspro"
	"userv4/modules/bdspro/dto"
)

type stubCostTypeService struct{}

func (stubCostTypeService) Create(context.Context, *bdspropb.CostTypeDTO) (*bdspropb.CostTypeDTO, error) {
	return &bdspropb.CostTypeDTO{Id: 1, TypeName: "test"}, nil
}

func (stubCostTypeService) Update(context.Context, *bdspropb.CostTypeDTO) (*bdspropb.CostTypeDTO, error) {
	return &bdspropb.CostTypeDTO{Id: 1}, nil
}

func (stubCostTypeService) Delete(context.Context, uint64) error { return nil }

func (stubCostTypeService) List(context.Context, dto.CostTypeSearchRequest) (*bdspropb.CostTypeListDTO, error) {
	return &bdspropb.CostTypeListDTO{Data: []*bdspropb.CostTypeDTO{}}, nil
}

func TestCostTypeListRouteCostTypeQuery(t *testing.T) {
	mux := http.NewServeMux()
	NewCostTypeController(stubCostTypeService{}).RegisterHTTP(mux)

	req := httptest.NewRequest(http.MethodGet, "/v2.1/transaction/cost-type?costType=10", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}
