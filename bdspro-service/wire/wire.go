//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"bdspro/config"
	bdspro_config "bdspro/config"
	bdspro_infra_client "bdspro/infra/client"
	bdspro_infra_postgre "bdspro/infra/postgre"
	bdspro_infra_postgre_org "bdspro/infra/postgre/org"
	bdspro_infra_postgre_user "bdspro/infra/postgre/user"
	bdspro_infra_providers "bdspro/infra/providers"
	bdspro_infra_repo "bdspro/infra/repo"
	bdspro_infra_repo_user "bdspro/infra/repo/user"
	bdspro_infra_router_group "bdspro/infra/router/group"
	bdspro_infra_router_org "bdspro/infra/router/org"
	bdspro_infra_router_public "bdspro/infra/router/public"
	bdspro_infra_router_user "bdspro/infra/router/user"
	bdspro_infra_services "bdspro/infra/services"
	bdspro_internal_repo "bdspro/internal/repo"
	bdspro_internal_repo_org "bdspro/internal/repo/org"
	bdspro_internal_repo_user "bdspro/internal/repo/user"
	bdspro_internal_usecases "bdspro/internal/usecases"
	bdspro_internal_usecases_group "bdspro/internal/usecases/group"
	bdspro_internal_usecases_interface "bdspro/internal/usecases/interface"
	bdspro_internal_usecases_org "bdspro/internal/usecases/org"
	bdspro_internal_usecases_shared "bdspro/internal/usecases/shared"
	bdspro_internal_usecases_user "bdspro/internal/usecases/user"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	wire.Bind(new(bdspro_internal_usecases_interface.IClientUC), new(*bdspro_infra_client.ServiceClient)),
	wire.Bind(new(bdspro_internal_repo.AmenityRepo), new(*bdspro_infra_postgre.PostgreAmenity)),
	wire.Bind(new(bdspro_internal_repo.ApartmentRepo), new(*bdspro_infra_postgre.PostgreApartment)),
	wire.Bind(new(bdspro_internal_repo_user.UAssetRepo), new(*bdspro_infra_postgre.PostgreAsset)),
	wire.Bind(new(bdspro_internal_repo.AssetShareRepo), new(*bdspro_infra_postgre.AssetShareRepo)),
	wire.Bind(new(bdspro_internal_repo.DocTypeRepo), new(*bdspro_infra_postgre.PostgreDocType)),
	wire.Bind(new(bdspro_internal_repo_org.OrgAssetCostRepo), new(*bdspro_infra_postgre_org.PostgreOrgAssetCost)),
	wire.Bind(new(bdspro_internal_repo_org.OrgAssetExploitationRepo), new(*bdspro_infra_postgre_org.PostgreOrgAssetExploitation)),
	wire.Bind(new(bdspro_internal_repo_org.OrgCostTypeRepo), new(*bdspro_infra_postgre_org.PostgreOrgCostType)),
	wire.Bind(new(bdspro_internal_repo.ProjectRepo), new(*bdspro_infra_postgre.PostgreProject)),
	wire.Bind(new(bdspro_internal_repo.PropertyTypeRepo), new(*bdspro_infra_postgre.PostgrePropertyType)),
	wire.Bind(new(bdspro_internal_repo.RecordHistoryRepo), new(*bdspro_infra_postgre.PostgreReportHistory)),
	wire.Bind(new(bdspro_internal_repo.RegionRepo), new(*bdspro_infra_postgre.PostgreRegion)),
	wire.Bind(new(bdspro_internal_repo_user.UserAssetCostRepo), new(*bdspro_infra_postgre_user.PostgreUserAssetCost)),
	wire.Bind(new(bdspro_internal_repo_user.UserAssetExploitationRepo), new(*bdspro_infra_postgre_user.PostgreUserAssetExploitation)),
	wire.Bind(new(bdspro_internal_repo_user.UserCostTypeRepo), new(*bdspro_infra_postgre_user.PostgreUserCostType)),
	wire.Bind(new(bdspro_internal_usecases_interface.IPlanUC), new(*bdspro_infra_providers.MemberPlanProvider)),
	wire.Bind(new(bdspro_internal_repo.AssetIncomeTypeRepo), new(*bdspro_infra_repo.GormAssetIncomeTypeRepo)),
	wire.Bind(new(bdspro_internal_repo.CostDocumentRepo), new(*bdspro_infra_repo.GormCostDocumentRepo)),
	wire.Bind(new(bdspro_internal_repo.DepositeRepo), new(*bdspro_infra_repo.GormDepositeRepo)),
	wire.Bind(new(bdspro_internal_repo.HouseInfoRepo), new(*bdspro_infra_repo.GormHouseInfoRepo)),
	wire.Bind(new(bdspro_internal_repo.IncomeDocumentRepo), new(*bdspro_infra_repo.GormIncomeDocumentRepo)),
	wire.Bind(new(bdspro_internal_repo.ProductAccessRepo), new(*bdspro_infra_repo.GormProductAccessRepo)),
	wire.Bind(new(bdspro_internal_repo.ProductHistoryRepo), new(*bdspro_infra_repo.PostgreHistoryRepo)),
	wire.Bind(new(bdspro_internal_repo.ProductPriceRepo), new(*bdspro_infra_repo.GormProductPriceRepo)),
	wire.Bind(new(bdspro_internal_repo.ProductPrivateRepo), new(*bdspro_infra_repo.GormProductPrivateRepo)),
	wire.Bind(new(bdspro_internal_repo.UPostRepo), new(*bdspro_infra_repo_user.GormPostRepo)),
	wire.Bind(new(bdspro_internal_repo.UAssetHistoryRepo), new(*bdspro_infra_repo.GormAssetHistoryRepo)),
	wire.Bind(new(bdspro_internal_repo.AssetLegalRepo), new(*bdspro_infra_repo.GormAssetLegalRepo)),
	wire.Bind(new(bdspro_internal_repo.PostMediaRepo), new(*bdspro_infra_repo.GormPostMediaRepo)),
	wire.Bind(new(bdspro_internal_repo_user.UProductRepo), new(*bdspro_infra_repo.GormUserProductRepo)),
	wire.Bind(new(bdspro_internal_repo.UProductChildRepo), new(*bdspro_infra_repo.PostgreProductChild)),
	wire.Bind(new(bdspro_internal_repo.ProductMediaRepo), new(*bdspro_infra_repo.GormProductMediaRepo)),
	wire.Bind(new(bdspro_internal_usecases.PlanUsecase), new(*bdspro_infra_services.MemberPlanUsecase)),
	bdspro_infra_postgre.NewActionRepo,
	bdspro_internal_usecases_shared.NewApartmentUsecase,
	bdspro_infra_services.NewAssetIncomeTypeServer,
	bdspro_internal_usecases.NewAssetIncomeTypeUc,
	bdspro_infra_services.NewAssetLegalServer,
	bdspro_internal_usecases.NewAssetLegalUc,
	bdspro_infra_postgre.NewAssetShareRepo,
	bdspro_infra_router_user.NewAssetShareRouter,
	bdspro_internal_usecases.NewAssetShareUC,
	bdspro_infra_services.NewAssetSplitHistoryServer,
	bdspro_internal_usecases.NewAssetSplitHistoryUc,
	bdspro_infra_services.NewBdsproPublicService,
	bdspro_config.NewConfigApp,
	bdspro_infra_services.NewCostDocumentServer,
	bdspro_internal_usecases.NewCostDocumentUc,
	bdspro_internal_repo.NewDeveloperRepo,
	bdspro_infra_postgre.NewGormAssetCostRepo,
	bdspro_infra_repo.NewGormAssetIncomeTypeRepo,
	bdspro_infra_repo.NewGormAssetLegalRepo,
	bdspro_infra_repo.NewGormAssetSplitHistoryRepo,
	bdspro_infra_repo.NewGormCostDocumentRepo,
	bdspro_infra_repo.NewGormDepositeRepo,
	bdspro_infra_repo.NewGormHouseInfoRepo,
	bdspro_infra_repo.NewGormIncomeDocumentRepo,
	bdspro_infra_repo.NewGormPostMediaRepo,
	bdspro_infra_repo_user.NewGormPostRepo,
	bdspro_infra_repo.NewGormProductMediaRepo,
	bdspro_infra_repo.NewGormProductPriceRepo,
	bdspro_infra_repo.NewGormProductPrivateRepo,
	bdspro_infra_repo.NewGormUserProductRepo,
	bdspro_infra_repo.NewGromProductAccessRepo,
	bdspro_infra_router_group.NewGroupAssetRouter,
	bdspro_internal_usecases_group.NewGroupAssetUsecase,
	bdspro_infra_router_group.NewGroupProductRouter,
	bdspro_internal_usecases_group.NewGroupProductUsecase,
	bdspro_infra_router_group.NewGroupRouter,
	bdspro_infra_services.NewGrpcAssetService,
	bdspro_infra_services.NewGrpcPostService,
	bdspro_infra_services.NewIncomeDocumentServer,
	bdspro_internal_usecases.NewIncomeDocumentUc,
	bdspro_infra_router_user.NewListRouter,
	bdspro_internal_usecases.NewListUsecase,
	bdspro_infra_providers.NewMemberPlanService,
	bdspro_infra_services.NewMemberPlanUsecase,
	bdspro_infra_providers.NewNotiService,
	bdspro_infra_router_org.NewOrgAssetCostRouter,
	bdspro_internal_usecases_org.NewOrgAssetCostUc,
	bdspro_infra_router_org.NewOrgAssetExploitationRouter,
	bdspro_internal_usecases_org.NewOrgAssetExploitationUc,
	bdspro_infra_router_org.NewOrgAssetRouter,
	bdspro_internal_usecases_org.NewOrgAssetUsecase,
	bdspro_infra_router_org.NewOrgCostTypeRouter,
	bdspro_internal_usecases_org.NewOrgCostTypeUc,
	bdspro_infra_router_org.NewOrgPostRoute,
	bdspro_internal_usecases_org.NewOrgPostUsecase,
	bdspro_infra_router_org.NewOrgProductChildRouter,
	bdspro_internal_usecases_org.NewOrgProductChildUsecase,
	bdspro_infra_router_org.NewOrgProductRouter,
	bdspro_internal_usecases_org.NewOrgProductUsecase,
	bdspro_infra_router_org.NewOrgRouterManager,
	bdspro_infra_router_user.NewPostMediaRoute,
	bdspro_infra_postgre.NewPostgreAmenity,
	bdspro_infra_postgre.NewPostgreApartmentRepo,
	bdspro_infra_postgre.NewPostgreAsset,
	bdspro_infra_postgre.NewPostgreDocType,
	bdspro_infra_repo.NewPostgreHistoryRepo,
	bdspro_infra_postgre_org.NewPostgreOrgAssetCost,
	bdspro_infra_postgre_org.NewPostgreOrgAssetExploitation,
	bdspro_infra_postgre_org.NewPostgreOrgCostType,
	bdspro_infra_repo.NewPostgreProductChildRepo,
	bdspro_infra_postgre.NewPostgreProject,
	bdspro_infra_postgre.NewPostgrePropertyType,
	bdspro_infra_postgre.NewPostgreRegion,
	bdspro_infra_postgre.NewPostgreReportHistory,
	bdspro_infra_postgre_user.NewPostgreUserAssetCost,
	bdspro_infra_postgre_user.NewPostgreUserAssetExploitation,
	bdspro_infra_postgre_user.NewPostgreUserCostType,
	bdspro_infra_services.NewProductService,
	bdspro_internal_usecases_shared.NewProductUsecase,
	bdspro_internal_usecases.NewProjectBuildUsecase,
	bdspro_internal_usecases.NewPropertyTypeUsecase,
	bdspro_infra_router_public.NewPublicRouter,
	bdspro_internal_usecases.NewRecordHistoryUsecase,
	bdspro_infra_providers.NewRedisProvider,
	bdspro_infra_router_user.NewRouterManager,
	bdspro_infra_client.NewServiceClient,
	bdspro_internal_usecases_shared.NewSharedAssetUsecase,
	bdspro_internal_usecases_shared.NewSharedPostUsecase,
	bdspro_infra_router_user.NewUAssetLegalRoute,
	bdspro_internal_usecases.NewUPostMediaUsecase,
	bdspro_infra_router_user.NewUProductAccessRouter,
	bdspro_infra_providers.NewUProductProvider,
	bdspro_infra_router_user.NewUserAssetCostRouter,
	bdspro_internal_usecases_user.NewUserAssetCostUc,
	bdspro_infra_router_user.NewUserAssetExploitationRouter,
	bdspro_internal_usecases_user.NewUserAssetExploitationUc,
	bdspro_infra_router_user.NewUserAssetRouter,
	bdspro_internal_usecases_user.NewUserAssetUsecase,
	bdspro_infra_router_user.NewUserCostTypeRouter,
	bdspro_internal_usecases_user.NewUserCostTypeUc,
	bdspro_internal_usecases_user.NewUserHandler,
	bdspro_infra_router_user.NewUserPostRoute,
	bdspro_internal_usecases_user.NewUserPostUsecase,
	bdspro_internal_usecases_user.NewUserProductAccessUsecase,
	bdspro_infra_router_user.NewUserProductChildRouter,
	bdspro_internal_usecases_user.NewUserProductChildUsecase,
	bdspro_infra_router_user.NewUserProductRouter,
	bdspro_internal_usecases_user.NewUserProductUsecase,
	bdspro_infra_router_user.NewUserProjectBuildRouter,
)
func InitializeApp() (*config.ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}