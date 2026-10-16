package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	bdspropb "pb/types/bdspro"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/services"
)

type stubCostTypeService struct{}

func (stubCostTypeService) Create(context.Context, *bdspropb.CostTypeDTO) (*bdspropb.CostTypeDTO, error) {
	return &bdspropb.CostTypeDTO{Id: 1, TypeName: "test"}, nil
}

func (stubCostTypeService) Update(context.Context, *bdspropb.CostTypeDTO) (*bdspropb.CostTypeDTO, error) {
	return &bdspropb.CostTypeDTO{Id: 1}, nil
}

func (stubCostTypeService) Delete(context.Context, uint64) error { return nil }

func (stubCostTypeService) List(context.Context, dto.CostTypeSearchRequest) (*services.CostTypeListOutcome, error) {
	return &services.CostTypeListOutcome{
		List: &bdspropb.CostTypeListDTO{Data: []*bdspropb.CostTypeDTO{}},
	}, nil
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

	reqV4 := httptest.NewRequest(http.MethodGet, "/v4/transaction/cost-type?costType=10", nil)
	recV4 := httptest.NewRecorder()
	mux.ServeHTTP(recV4, reqV4)
	if recV4.Code != http.StatusOK {
		t.Fatalf("v4 status = %d, body = %s", recV4.Code, recV4.Body.String())
	}
}
