package services

import (
	"context"

	bdspropb "pb/types/bdspro"

	"userv4/modules/bdspro/dto"
)

// DealService owns deal list/detail/timestamps/process flows for modules/bdspro.
type DealService interface {
	GetDealsUser(ctx context.Context, req dto.DealSearchRequest) (*bdspropb.GetDealsResponse, error)
	GetDealsByProduct(ctx context.Context, productID uint64, req dto.SyncByIDRequest) (*bdspropb.GetOrganizationDealsResponse, error)
	GetTimestamps(ctx context.Context, id uint64) (*bdspropb.DealTimestampsResponse, error)
	GetDetail(ctx context.Context, req dto.DealDetailRequest) (*bdspropb.GroupDeal, int64, error)
	GetContractProcess(ctx context.Context, dealID uint64) (*bdspropb.ContractProcessResponse, int64, error)
	GetDealProducts(ctx context.Context, req dto.DealListByIDRequest) (*bdspropb.GetDealProductsResponse, int64, error)
	GetContracts(ctx context.Context, req dto.DealListByIDRequest) (*bdspropb.DealContractListResponse, error)
	GetInvestments(ctx context.Context, req dto.DealListByIDRequest) (*bdspropb.ListInvestment, int64, error)
	GetInvitations(ctx context.Context, req dto.DealListByIDRequest) (*bdspropb.GetDealInvitationsByDealResponse, error)
	GetMembers(ctx context.Context, req dto.DealMembersRequest) (*bdspropb.GetDealMembersResponse, int64, error)
	GetDealCosts(ctx context.Context, req dto.DealListByIDRequest) (*bdspropb.DealCostListDTO, error)
	GetSummary(ctx context.Context, req dto.DealSummaryRequest) (*bdspropb.SummaryResponse, int64, error)
	GetBankAccount(ctx context.Context, req dto.DealDetailRequest) (*bdspropb.BankAccountItem, int64, error)
	GetContractStatistics(ctx context.Context, req dto.DealStatisticsRequest) (*bdspropb.ContractDealStatisticsDTO, int64, error)
	CreateGroupDeal(ctx context.Context, req *bdspropb.CreateGroupDealRequest) (uint64, error)
	SeedDemoRelated(ctx context.Context, dealID, profileID, customerID uint64, dealType uint32) error
	DeleteGroupDeal(ctx context.Context, id uint64) error
}
