package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	bdspropb "pb/types/bdspro"
	"userv4/modules/bdspro/dto"
)

type routeStubDeal struct{}

func (routeStubDeal) GetDealsUser(context.Context, dto.DealSearchRequest) (*bdspropb.GetDealsResponse, error) {
	return &bdspropb.GetDealsResponse{}, nil
}

func (routeStubDeal) GetDealsByProduct(context.Context, uint64, dto.SyncByIDRequest) ([]dto.ProductDealItem, int64, error) {
	return []dto.ProductDealItem{}, 1, nil
}

func (routeStubDeal) GetTimestamps(context.Context, uint64) (*bdspropb.DealTimestampsResponse, error) {
	return &bdspropb.DealTimestampsResponse{Timestamps: map[string]int64{}}, nil
}

func (routeStubDeal) GetDetail(context.Context, dto.DealDetailRequest) (*bdspropb.GroupDeal, int64, error) {
	return &bdspropb.GroupDeal{Id: 14}, 123, nil
}

func (routeStubDeal) GetContractProcess(context.Context, uint64) (*bdspropb.ContractProcessResponse, int64, error) {
	return &bdspropb.ContractProcessResponse{}, 456, nil
}

func (routeStubDeal) GetDealProducts(context.Context, dto.DealListByIDRequest) ([]dto.DealProductItem, int64, error) {
	return []dto.DealProductItem{}, 1, nil
}

func (routeStubDeal) GetContracts(context.Context, dto.DealListByIDRequest) (*bdspropb.DealContractListResponse, error) {
	return &bdspropb.DealContractListResponse{}, nil
}

func (routeStubDeal) GetInvestments(context.Context, dto.DealListByIDRequest) (*bdspropb.ListInvestment, int64, error) {
	return &bdspropb.ListInvestment{}, 1, nil
}

func (routeStubDeal) GetInvitations(context.Context, dto.DealListByIDRequest) (*bdspropb.GetDealInvitationsByDealResponse, error) {
	return &bdspropb.GetDealInvitationsByDealResponse{}, nil
}

func (routeStubDeal) GetMembers(context.Context, dto.DealMembersRequest) (*bdspropb.GetDealMembersResponse, int64, error) {
	return &bdspropb.GetDealMembersResponse{}, 1, nil
}

func (routeStubDeal) GetDealCosts(context.Context, dto.DealListByIDRequest) (*bdspropb.DealCostListDTO, error) {
	return &bdspropb.DealCostListDTO{}, nil
}

func (routeStubDeal) GetSummary(context.Context, dto.DealSummaryRequest) (*bdspropb.SummaryResponse, int64, error) {
	return &bdspropb.SummaryResponse{AmountTarget: 1}, 1, nil
}

func (routeStubDeal) GetBankAccount(context.Context, dto.DealDetailRequest) (*bdspropb.BankAccountItem, int64, error) {
	return &bdspropb.BankAccountItem{BankName: "VCB"}, 1, nil
}

func (routeStubDeal) GetContractStatistics(context.Context, dto.DealStatisticsRequest) (*bdspropb.ContractDealStatisticsDTO, int64, error) {
	return &bdspropb.ContractDealStatisticsDTO{DealId: 13}, 1, nil
}

func (routeStubDeal) CreateGroupDeal(context.Context, *bdspropb.CreateGroupDealRequest) (uint64, error) {
	return 14, nil
}

func (routeStubDeal) UpdateGroupDeal(context.Context, *bdspropb.UpdateGroupDealRequest) (uint64, error) {
	return 14, nil
}

func (routeStubDeal) UpdateDealSetting(context.Context, *bdspropb.UpdateDealSettingRequest) (uint64, error) {
	return 14, nil
}

func (routeStubDeal) SeedDemoRelated(context.Context, uint64, uint64, uint64, uint32) error { return nil }

func (routeStubDeal) DeleteGroupDeal(context.Context, uint64) error { return nil }

type stubDealCost struct{}

func (stubDealCost) Search(context.Context, dto.DealCostSearchRequest) (*bdspropb.DealCostListDTO, error) {
	return &bdspropb.DealCostListDTO{}, nil
}
func (stubDealCost) Detail(context.Context, uint64) (*bdspropb.DealCostDTO, error) {
	return &bdspropb.DealCostDTO{}, nil
}
func (stubDealCost) Create(context.Context, *bdspropb.DealCostDTO) (*bdspropb.DealCostDTO, error) {
	return &bdspropb.DealCostDTO{}, nil
}
func (stubDealCost) Update(context.Context, *bdspropb.DealCostDTO) (*bdspropb.DealCostDTO, error) {
	return &bdspropb.DealCostDTO{}, nil
}
func (stubDealCost) Delete(context.Context, uint64) error { return nil }

type stubContractDeal struct{}

func (stubContractDeal) GetContract(context.Context, uint64) (*bdspropb.DealContractResponse, error) {
	return &bdspropb.DealContractResponse{}, nil
}
func (stubContractDeal) DeleteContract(context.Context, uint64) (uint64, error) { return 1, nil }
func (stubContractDeal) CreateContractPayment(context.Context, *bdspropb.DealContractPaymentRequest) (*bdspropb.DealContractResponse, error) {
	return &bdspropb.DealContractResponse{}, nil
}
func (stubContractDeal) CreateActionPayment(context.Context, *bdspropb.PaymentActionRequest) (*bdspropb.ActionResponse, error) {
	return &bdspropb.ActionResponse{}, nil
}
func (stubContractDeal) GetMyContracts(context.Context, uint32, uint32) (*bdspropb.DealContractListResponse, error) {
	return &bdspropb.DealContractListResponse{}, nil
}

type stubInvestment struct{}

func (stubInvestment) Create(context.Context, *bdspropb.SubmitInvestRequest) (uint64, error) {
	return 1, nil
}
func (stubInvestment) Update(context.Context, *bdspropb.SubmitInvestRequest) (uint64, error) {
	return 1, nil
}
func (stubInvestment) Delete(context.Context, uint64) error { return nil }
func (stubInvestment) Detail(context.Context, uint64) (*bdspropb.Investment, error) {
	return &bdspropb.Investment{}, nil
}
func (stubInvestment) List(context.Context, *bdspropb.InvestmentListRequest) (*bdspropb.ListInvestment, error) {
	return &bdspropb.ListInvestment{}, nil
}
func (stubInvestment) Approve(context.Context, uint64) (uint64, error)  { return 1, nil }
func (stubInvestment) Reject(context.Context, uint64) (uint64, error)   { return 1, nil }
func (stubInvestment) Confirm(context.Context, uint64, bool) (uint64, error) { return 1, nil }
func (stubInvestment) MemberHistories(context.Context, uint64, uint64) (*bdspropb.GetInvestmentHistoriesResponse, error) {
	return &bdspropb.GetInvestmentHistoriesResponse{}, nil
}

func TestDealRoutes_RegisterNoConflict(t *testing.T) {
	mux := http.NewServeMux()
	NewDealController(routeStubDeal{}).RegisterHTTP(mux)
	NewDealCostController(stubDealCost{}).RegisterHTTP(mux)
	NewDealContractController(routeStubDeal{}, stubContractDeal{}).RegisterHTTP(mux)
	NewInvestmentController(stubInvestment{}).RegisterHTTP(mux)

	for _, path := range []string{
		"/v4/deals/me",
		"/v4/deal/14/ts",
		"/v4/deal/14/products",
		"/v4/deal/14/summary",
		"/v4/deal/14/bank-account",
		"/v4/deal/14/members",
		"/v4/deal/14",
		"/v4/deal/invitation/list/14",
		"/v4/deal/investments/list/14",
		"/v4/transaction/contract/deal/14/list",
		"/v4/transaction/contract/deal/13/statistic",
		"/v4/transaction/contract/deal/14/statistics?payment=true&cost=true",
		"/v4/transaction/contract/deal/14/process",
		"/v2.1/transaction/contract/deal/list/14",
		"/v2.1/transaction/contract/deal/statistics/13",
		"/v2.1/transaction/contract/deal/14/process",
		"/v2.1/transaction/deal-cost/list/14",
		"/v4/transaction/deal-cost/list/14",
		"/v2.1/org/investment/1",
		"/v2.1/org/investments/list/14",
		"/v4/org/investment/1",
		"/v4/deal/investments/list/14",
	} {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, path, nil))
		if rr.Code == http.StatusNotFound {
			t.Fatalf("%s => %d", path, rr.Code)
		}
	}

	for _, tc := range []struct {
		method string
		path   string
		body   string
	}{
		{http.MethodPost, "/v4/deal", `{"name":"TV demo"}`},
		{http.MethodPut, "/v4/deal/14", `{"name":"TV updated"}`},
		{http.MethodPut, "/v4/deal/14/setting", `{"allowSharing":true}`},
		{http.MethodDelete, "/v4/deal/14", ""},
	} {
		var req *http.Request
		if tc.body != "" {
			req = httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
		} else {
			req = httptest.NewRequest(tc.method, tc.path, nil)
		}
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code == http.StatusNotFound {
			t.Fatalf("%s %s => %d", tc.method, tc.path, rr.Code)
		}
	}
}
