package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	bdspropb "pb/types/bdspro"
	"userv4/modules/bdspro/dto"
)

type routeStubDeal struct{}

func (routeStubDeal) GetDealsUser(context.Context, dto.DealSearchRequest) (*bdspropb.GetDealsResponse, error) {
	return &bdspropb.GetDealsResponse{}, nil
}

func (routeStubDeal) GetDealsByProduct(context.Context, uint64, dto.SyncByIDRequest) (*bdspropb.GetOrganizationDealsResponse, error) {
	return &bdspropb.GetOrganizationDealsResponse{}, nil
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

func (routeStubDeal) GetDealProducts(context.Context, dto.DealListByIDRequest) (*bdspropb.GetDealProductsResponse, int64, error) {
	return &bdspropb.GetDealProductsResponse{}, 1, nil
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
	return 0, nil
}

func (routeStubDeal) SeedDemoRelated(context.Context, uint64, uint64, uint64, uint32) error { return nil }

func (routeStubDeal) DeleteGroupDeal(context.Context, uint64) error { return nil }

func TestDealRoutes_RegisterNoConflict(t *testing.T) {
	mux := http.NewServeMux()
	NewDealController(routeStubDeal{}).RegisterHTTP(mux)

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
		"/v4/transaction/contract/deal/14/process",
		"/v4/transaction/deal-cost/list/14",
	} {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, path, nil))
		if rr.Code == http.StatusNotFound {
			t.Fatalf("%s => %d", path, rr.Code)
		}
	}
}
