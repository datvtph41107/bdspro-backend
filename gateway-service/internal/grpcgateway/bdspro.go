package grpcgateway

import (
	"context"
	bdspropb "pb/types/bdspro"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterBdsproService(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return registerAll(ctx, mux, conn,
		bdspropb.RegisterBdsproPublicServiceHandler,
		bdspropb.RegisterProductServiceHandler,
		bdspropb.RegisterAssetServiceHandler,
		bdspropb.RegisterPostServiceHandler,
		bdspropb.RegisterCostDocumentServiceHandler,
		bdspropb.RegisterIncomeDocumentServiceHandler,
		bdspropb.RegisterAssetIncomeTypeServiceHandler,
		bdspropb.RegisterAssetExploitationServiceHandler,
		bdspropb.RegisterAssetCostTypeServiceHandler,
		bdspropb.RegisterAssetCostServiceHandler,
		bdspropb.RegisterAssetLegalServiceHandler,
		bdspropb.RegisterAssetSplitHistoryServiceHandler,
		bdspropb.RegisterSharingAccessServiceHandler,
		bdspropb.RegisterTransactionServiceHandler,
		bdspropb.RegisterProjectBuildServiceHandler,
		bdspropb.RegisterAdminProductServiceHandler,
		bdspropb.RegisterAdminAssetServiceHandler,
		bdspropb.RegisterAdminPostServiceHandler,
		bdspropb.RegisterAdminPropertyTypeServiceHandler,
		bdspropb.RegisterAdminProjectServiceHandler,
		bdspropb.RegisterAdminAreaRegionServiceHandler,
		bdspropb.RegisterAdminRegionServiceHandler,
		bdspropb.RegisterAdminTransactionServiceHandler,
		bdspropb.RegisterBdsproDashboardServiceHandler,
		bdspropb.RegisterBdsDomainServiceHandler,
		bdspropb.RegisterPropertyServiceHandler,
		bdspropb.RegisterAdminTxServiceHandler,
		bdspropb.RegisterTxServiceHandler,
		bdspropb.RegisterDealContractServiceHandler,
		bdspropb.RegisterCostTypeServiceHandler,
		bdspropb.RegisterDealCostServiceHandler,
		bdspropb.RegisterDealServiceHandler,
		bdspropb.RegisterDealCommissionServiceHandler,
		bdspropb.RegisterInvestmentServiceHandler,
		bdspropb.RegisterDealMemberServiceHandler,
		bdspropb.RegisterDealMilestoneServiceHandler,
		bdspropb.RegisterDealInternalNoteServiceHandler,
		bdspropb.RegisterDistributionServiceHandler,
		bdspropb.RegisterProductNoteServiceHandler,
		bdspropb.RegisterIdentifierServiceHandler,
	)
}
