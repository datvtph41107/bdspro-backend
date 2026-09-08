package initial

import (
	"bdspro/infra/handler"
	admin_handler "bdspro/infra/handler/admin"
	property_handler "bdspro/infra/handler/property"
	redis_cli "bdspro/infra/redis"
	"bdspro/internal/job"
	shared_usecase "bdspro/internal/usecases/shared"
)

type InitialApp struct {
	// GrpcServer        *server.BdsproGrpcServer
	// GrpcGenericServer *server.GenericGrpcServer
	PostService    *handler.PostService
	ProductService *handler.ProductHandler
	AssetService   *handler.AssetService
	// AdminRouter *admin_router.RouterManager
	// UserRouter        *user_router.RouterManager
	BdsproPublicService      *handler.BdsproPublicService
	AssetExploitationService *handler.AssetExploitationService
	AssetCostService         *handler.AssetCostService
	AssetCostTypeService     *handler.AssetCostTypeService
	SharingAccessService     *handler.SharingAccessService
	AssetLegalService        *handler.AssetLegalService
	ProductUsecase           *shared_usecase.ProductUsecase
	RedisClient              *redis_cli.RedisClient
	ProjectBuildService      *handler.ProjectBuildService
	TransactionService       *handler.TransactionHandler
	AdminProductHandler      *admin_handler.AdminProductHandler
	AdminAssetHandler        *admin_handler.AdminAssetHandler
	AdminPostHandler         *admin_handler.AdminPostHandler
	AdminPropertyTypeHandler *admin_handler.AdminPropertyTypeHandler
	AdminProjectHandler      *admin_handler.AdminProjectHandler
	AdminAreaRegionHandler   *handler.AdminAreaRegionHandler
	AdminRegionHandler       *handler.AdminRegionHandler
	InternalHandler          *handler.InternalHandler
	BdsproDashboardHandler   *handler.BdsproDashboardHandler
	BdsDomainHandler         *handler.BdsDomainHandler
	AssetStatsJob            *job.AssetDashboardStatsJob
	ProductStatsJob          *job.ProductDashboardStatsJob
	PostStatsJob             *job.PostDashboardStatsJob

	DealHandler             *handler.DealHandler
	DealCommissionHandler   *handler.DealCommissionHandler
	DealInvestmentHandler   *handler.DealInvestmentHandler
	DealMilestoneHandler    *handler.DealMilestoneHandler
	DealInvitationHandler   *handler.DealInvitationHandler
	DealInternalNoteHandler *handler.DealInternalNoteHandler

	TxInternalHandler     *handler.InternalTransactionHandler
	TxHandler             *handler.TxHandler
	TxContractDealHandler *handler.TxContractDealHandler
	TxCostTypeHandler     *handler.TxCostTypeHandler
	TxDealCostHandler     *handler.TxDealCostHandler
	AdminTxHandler        *admin_handler.AdminTxHandler
	PropertyHandler       *property_handler.PropertyHandler
	DistributeHandler     *handler.DistributionHandler
	ProductNoteHandler    *handler.ProductNoteHandler
	PropertyIdentifier    *handler.PropertyIdentifierHandler
	// @bind: register_server_config
	// CostDocumentServer *handler.CostDocumentServer
	// IncomeDocumentServer    *handler.IncomeDocumentServer
	// AssetIncomeTypeServer   *handler.AssetIncomeTypeServer
	// AssetExploitationServer *handler.AssetExploitationServer
	// AssetCostTypeServer     *handler.AssetCostTypeServer
	// AssetCostServer         *handler.AssetCostServer
	// AssetLegalServer        *handler.AssetLegalServer
	// AssetSplitHistoryServer *handler.AssetSplitHistoryServer
}

func NewInitialApp(
	// GrpcGenericServer *server.GenericGrpcServer,
	postService *handler.PostService,
	productService *handler.ProductHandler,
	bdsproPublicService *handler.BdsproPublicService,
	assetService *handler.AssetService,
	assetExploitationServer *handler.AssetExploitationService,
	assetCostService *handler.AssetCostService,
	assetCostTypeService *handler.AssetCostTypeService,
	sharingAccessService *handler.SharingAccessService,
	productUsecase *shared_usecase.ProductUsecase,
	// adminRoute *admin_router.RouterManager,
	// @bind: register_server_param
	// CostDocumentServer *handler.CostDocumentServer,
	// IncomeDocumentServer *handler.IncomeDocumentServer,
	// AssetIncomeTypeServer *handler.AssetIncomeTypeServer,
	// AssetExploitationServer *handler.AssetExploitationServer,
	// AssetCostTypeServer *handler.AssetCostTypeServer,
	// AssetCostServer *handler.AssetCostServer,
	// AssetLegalServer *handler.AssetLegalServer,
	// AssetSplitHistoryServer *handler.AssetSplitHistoryServer,
	assetLegalService *handler.AssetLegalService,
	projectBuildService *handler.ProjectBuildService,
	transactionService *handler.TransactionHandler,
	adminProductHandler *admin_handler.AdminProductHandler,
	adminAssetHandler *admin_handler.AdminAssetHandler,
	adminPostHandler *admin_handler.AdminPostHandler,
	adminPropertyTypeHandler *admin_handler.AdminPropertyTypeHandler,
	adminProjectHandler *admin_handler.AdminProjectHandler,
	adminAreaRegionHandler *handler.AdminAreaRegionHandler,
	adminRegionHandler *handler.AdminRegionHandler,
	internalHandler *handler.InternalHandler,
	bdsproDashboardHandler *handler.BdsproDashboardHandler,
	bdsDomainHandler *handler.BdsDomainHandler,
	// @bind: register_server_new
	assetStatsJob *job.AssetDashboardStatsJob,
	productStatsJob *job.ProductDashboardStatsJob,
	postStatsJob *job.PostDashboardStatsJob,

	dealHandler *handler.DealHandler,
	dealCommissionHandler *handler.DealCommissionHandler,
	dealInvestmentHandler *handler.DealInvestmentHandler,
	dealMilestoneHandler *handler.DealMilestoneHandler,
	dealInvitationHandler *handler.DealInvitationHandler,
	dealInternalNoteHandler *handler.DealInternalNoteHandler,
	txInternalHandler *handler.InternalTransactionHandler,
	txHandler *handler.TxHandler,
	txContractDealHandler *handler.TxContractDealHandler,
	txCostTypeHandler *handler.TxCostTypeHandler,
	txDealCostHandler *handler.TxDealCostHandler,
	adminTxHandler *admin_handler.AdminTxHandler,
	// CostDocumentServer: CostDocumentServer,
	// IncomeDocumentServer:    IncomeDocumentServer,
	// AssetIncomeTypeServer:   AssetIncomeTypeServer,
	// AssetExploitationServer: AssetExploitationServer,
	// AssetCostTypeServer:     AssetCostTypeServer,
	// AssetCostServer:         AssetCostServer,
	// AssetLegalServer:        AssetLegalServer,
	// AssetSplitHistoryServer: AssetSplitHistoryServer,
	propertyHandler *property_handler.PropertyHandler,
	distributeHandler *handler.DistributionHandler,
	productNoteHandler *handler.ProductNoteHandler,
	propertyIdentifier *handler.PropertyIdentifierHandler,
	redisClient *redis_cli.RedisClient,
) *InitialApp {
	return &InitialApp{
		// GrpcServer:        GrpcServer,
		// GrpcGenericServer: GrpcGenericServer,
		PostService:              postService,
		ProductService:           productService,
		BdsproPublicService:      bdsproPublicService,
		AssetService:             assetService,
		AssetExploitationService: assetExploitationServer,
		AssetCostService:         assetCostService,
		AssetCostTypeService:     assetCostTypeService,
		SharingAccessService:     sharingAccessService,
		AssetLegalService:        assetLegalService,
		ProductUsecase:           productUsecase,
		RedisClient:              redisClient,
		ProjectBuildService:      projectBuildService,
		TransactionService:       transactionService,
		AdminProductHandler:      adminProductHandler,
		AdminAssetHandler:        adminAssetHandler,
		AdminPostHandler:         adminPostHandler,
		AdminPropertyTypeHandler: adminPropertyTypeHandler,
		AdminProjectHandler:      adminProjectHandler,
		AdminAreaRegionHandler:   adminAreaRegionHandler,
		AdminRegionHandler:       adminRegionHandler,
		InternalHandler:          internalHandler,
		BdsproDashboardHandler:   bdsproDashboardHandler,
		BdsDomainHandler:         bdsDomainHandler,

		DealHandler:             dealHandler,
		DealCommissionHandler:   dealCommissionHandler,
		DealInvestmentHandler:   dealInvestmentHandler,
		DealMilestoneHandler:    dealMilestoneHandler,
		DealInvitationHandler:   dealInvitationHandler,
		DealInternalNoteHandler: dealInternalNoteHandler,

		TxInternalHandler:     txInternalHandler,
		TxHandler:             txHandler,
		TxContractDealHandler: txContractDealHandler,
		TxCostTypeHandler:     txCostTypeHandler,
		TxDealCostHandler:     txDealCostHandler,
		AdminTxHandler:        adminTxHandler,

		PropertyHandler:    propertyHandler,
		DistributeHandler:  distributeHandler,
		AssetStatsJob:      assetStatsJob,
		ProductStatsJob:    productStatsJob,
		PostStatsJob:       postStatsJob,
		ProductNoteHandler: productNoteHandler,
		PropertyIdentifier: propertyIdentifier,
		// AdminRouter: adminRoute,
		// @bind: register_server_new
		// CostDocumentServer: CostDocumentServer,
		// IncomeDocumentServer:    IncomeDocumentServer,
		// AssetIncomeTypeServer:   AssetIncomeServer,
		// AssetExploitationServer: AssetExploitationServer,
		// AssetCostTypeServer:     AssetCostTypeServer,
		// AssetCostServer:         AssetCostServer,
		// AssetLegalServer:        AssetLegalServer,
		// AssetSplitHistoryServer: AssetSplitHistoryServer,
	}
}
