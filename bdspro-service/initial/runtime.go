package initial

import (
	admin_handler "bdspro/infra/handler/admin"
	group_router "bdspro/infra/router/group"
	org_router "bdspro/infra/router/org"
	public_router "bdspro/infra/router/public"
	user_router "bdspro/infra/router/user"
	service "bdspro/infra/services"
	shared_usecase "bdspro/internal/usecases/shared"
)

type InitialApp struct {
	// GrpcServer        *server.BdsproGrpcServer
	// GrpcGenericServer *server.GenericGrpcServer
	PostService    *service.PostService
	ProductService *service.ProductService
	AssetService   *service.AssetService
	UserRouter     *user_router.RouterManager
	OrgRouter      *org_router.RouterManager
	// AdminRouter *admin_router.RouterManager
	// UserRouter        *user_router.RouterManager
	GroupRouter              *group_router.GroupManager
	EnumRoute                *public_router.PublicRouter
	BdsproPublicService      *service.BdsproPublicService
	AssetExploitationService *service.AssetExploitationService
	AssetCostService         *service.AssetCostService
	AssetCostTypeService     *service.AssetCostTypeService
	SharingAccessService     *service.SharingAccessService
	AssetLegalService        *service.AssetLegalService
	ProductUsecase           *shared_usecase.ProductUsecase
	ProjectBuildService      *service.ProjectBuildService
	TransactionService       *service.TransactionHandler
	AdminProductHandler      *admin_handler.AdminProductHandler
	AdminAssetHandler        *admin_handler.AdminAssetHandler
	AdminPostHandler         *admin_handler.AdminPostHandler
	AdminPropertyTypeHandler *admin_handler.AdminPropertyTypeHandler
	AdminProjectHandler      *admin_handler.AdminProjectHandler

	// @bind: register_server_config
	// CostDocumentServer *service.CostDocumentServer
	// IncomeDocumentServer    *service.IncomeDocumentServer
	// AssetIncomeTypeServer   *service.AssetIncomeTypeServer
	// AssetExploitationServer *service.AssetExploitationServer
	// AssetCostTypeServer     *service.AssetCostTypeServer
	// AssetCostServer         *service.AssetCostServer
	// AssetLegalServer        *service.AssetLegalServer
	// AssetSplitHistoryServer *service.AssetSplitHistoryServer
}

func NewInitialApp(
	// GrpcGenericServer *server.GenericGrpcServer,
	userRouter *user_router.RouterManager,
	orgRouter *org_router.RouterManager,
	groupRouter *group_router.GroupManager,
	enumRouter *public_router.PublicRouter,
	postService *service.PostService,
	productService *service.ProductService,
	bdsproPublicService *service.BdsproPublicService,
	assetService *service.AssetService,
	assetExploitationServer *service.AssetExploitationService,
	assetCostService *service.AssetCostService,
	assetCostTypeService *service.AssetCostTypeService,
	sharingAccessService *service.SharingAccessService,
	productUsecase *shared_usecase.ProductUsecase,
	// adminRoute *admin_router.RouterManager,
	// @bind: register_server_param
	// CostDocumentServer *service.CostDocumentServer,
	// IncomeDocumentServer *service.IncomeDocumentServer,
	// AssetIncomeTypeServer *service.AssetIncomeTypeServer,
	// AssetExploitationServer *service.AssetExploitationServer,
	// AssetCostTypeServer *service.AssetCostTypeServer,
	// AssetCostServer *service.AssetCostServer,
	// AssetLegalServer *service.AssetLegalServer,
	// AssetSplitHistoryServer *service.AssetSplitHistoryServer,
	assetLegalService *service.AssetLegalService,
	projectBuildService *service.ProjectBuildService,
	transactionService *service.TransactionHandler,
	adminProductHandler *admin_handler.AdminProductHandler,
	adminAssetHandler *admin_handler.AdminAssetHandler,
	adminPostHandler *admin_handler.AdminPostHandler,
	adminPropertyTypeHandler *admin_handler.AdminPropertyTypeHandler,
	adminProjectHandler *admin_handler.AdminProjectHandler,
	// @bind: register_server_new
	// CostDocumentServer: CostDocumentServer,
	// IncomeDocumentServer:    IncomeDocumentServer,
	// AssetIncomeTypeServer:   AssetIncomeTypeServer,
	// AssetExploitationServer: AssetExploitationServer,
	// AssetCostTypeServer:     AssetCostTypeServer,
	// AssetCostServer:         AssetCostServer,
	// AssetLegalServer:        AssetLegalServer,
	// AssetSplitHistoryServer: AssetSplitHistoryServer,
) *InitialApp {
	return &InitialApp{
		// GrpcServer:        GrpcServer,
		// GrpcGenericServer: GrpcGenericServer,
		PostService:              postService,
		ProductService:           productService,
		UserRouter:               userRouter,
		OrgRouter:                orgRouter,
		GroupRouter:              groupRouter,
		EnumRoute:                enumRouter,
		BdsproPublicService:      bdsproPublicService,
		AssetService:             assetService,
		AssetExploitationService: assetExploitationServer,
		AssetCostService:         assetCostService,
		AssetCostTypeService:     assetCostTypeService,
		SharingAccessService:     sharingAccessService,
		AssetLegalService:        assetLegalService,
		ProductUsecase:           productUsecase,
		ProjectBuildService:      projectBuildService,
		TransactionService:       transactionService,
		AdminProductHandler:      adminProductHandler,
		AdminAssetHandler:        adminAssetHandler,
		AdminPostHandler:         adminPostHandler,
		AdminPropertyTypeHandler: adminPropertyTypeHandler,
		AdminProjectHandler:      adminProjectHandler,
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
