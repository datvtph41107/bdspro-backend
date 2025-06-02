//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"bdspro/config"
	bdspro_config "bdspro/config"
	bdspro_external_handler "bdspro/external/handler"
	bdspro_external_providers "bdspro/external/providers"
	bdspro_external_repo "bdspro/external/repo"
	bdspro_external_repo_user "bdspro/external/repo/user"
	bdspro_external_router_group "bdspro/external/router/group"
	bdspro_external_router_org "bdspro/external/router/org"
	bdspro_external_router_public "bdspro/external/router/public"
	bdspro_external_router_user "bdspro/external/router/user"
	bdspro_external_services "bdspro/external/services"
	bdspro_internal_repo "bdspro/internal/repo"
	bdspro_internal_usecases "bdspro/internal/usecases"
	bdspro_internal_usecases_group "bdspro/internal/usecases/group"
	bdspro_internal_usecases_org "bdspro/internal/usecases/org"
	bdspro_v1_mod_admin "bdspro/v1/mod/admin"
	bdspro_v1_mod_admin_amenity "bdspro/v1/mod/admin/amenity"
	bdspro_v1_mod_admin_apartment "bdspro/v1/mod/admin/apartment"
	bdspro_v1_mod_admin_build_room "bdspro/v1/mod/admin/build_room"
	bdspro_v1_mod_admin_developer "bdspro/v1/mod/admin/developer"
	bdspro_v1_mod_admin_district "bdspro/v1/mod/admin/district"
	bdspro_v1_mod_admin_doc_type "bdspro/v1/mod/admin/doc_type"
	bdspro_v1_mod_admin_product "bdspro/v1/mod/admin/product"
	bdspro_v1_mod_admin_product_access "bdspro/v1/mod/admin/product_access"
	bdspro_v1_mod_admin_project "bdspro/v1/mod/admin/project"
	bdspro_v1_mod_admin_project_build "bdspro/v1/mod/admin/project_build"
	bdspro_v1_mod_admin_property_type "bdspro/v1/mod/admin/property_type"
	bdspro_v1_mod_admin_province "bdspro/v1/mod/admin/province"
	bdspro_v1_mod_admin_ward "bdspro/v1/mod/admin/ward"
	bdspro_v1_mod_global_amenity "bdspro/v1/mod/global/amenity"
	bdspro_v1_mod_global_apartment_attribute "bdspro/v1/mod/global/apartment_attribute"
	bdspro_v1_mod_global_history_asset "bdspro/v1/mod/global/history/asset"
	bdspro_v1_mod_global_house_info "bdspro/v1/mod/global/house_info"
	bdspro_v1_mod_global_product_access "bdspro/v1/mod/global/product_access"
	bdspro_v1_mod_global_product_history "bdspro/v1/mod/global/product_history"
	bdspro_v1_mod_global_product_media "bdspro/v1/mod/global/product_media"
	bdspro_v1_mod_global_product_price "bdspro/v1/mod/global/product_price"
	bdspro_v1_mod_global_product_private "bdspro/v1/mod/global/product_private"
	bdspro_v1_mod_global_region "bdspro/v1/mod/global/region"
	bdspro_v1_mod_user "bdspro/v1/mod/user"
	bdspro_v1_mod_user_action "bdspro/v1/mod/user/action"
	bdspro_v1_mod_user_apartment "bdspro/v1/mod/user/apartment"
	bdspro_v1_mod_user_asset "bdspro/v1/mod/user/asset"
	bdspro_v1_mod_user_asset2 "bdspro/v1/mod/user/asset2"
	bdspro_v1_mod_user_asset_cost "bdspro/v1/mod/user/asset_cost"
	bdspro_v1_mod_user_asset_cost_type "bdspro/v1/mod/user/asset_cost_type"
	bdspro_v1_mod_user_asset_exploitation "bdspro/v1/mod/user/asset_exploitation"
	bdspro_v1_mod_user_asset_income_type "bdspro/v1/mod/user/asset_income_type"
	bdspro_v1_mod_user_asset_legal "bdspro/v1/mod/user/asset_legal"
	bdspro_v1_mod_user_asset_share "bdspro/v1/mod/user/asset_share"
	bdspro_v1_mod_user_build_room "bdspro/v1/mod/user/build_room"
	bdspro_v1_mod_user_like "bdspro/v1/mod/user/like"
	bdspro_v1_mod_user_list "bdspro/v1/mod/user/list"
	bdspro_v1_mod_user_post "bdspro/v1/mod/user/post"
	bdspro_v1_mod_user_product "bdspro/v1/mod/user/product"
	bdspro_v1_mod_user_product_access "bdspro/v1/mod/user/product_access"
	bdspro_v1_mod_user_product_child "bdspro/v1/mod/user/product_child"
	bdspro_v1_mod_user_product_interface "bdspro/v1/mod/user/product/interface"
	bdspro_v1_mod_user_project_build "bdspro/v1/mod/user/project_build"
	bdspro_v1_mod_user_transaction "bdspro/v1/mod/user/transaction"
	bdspro_v1_providers "bdspro/v1/providers"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	wire.Bind(new(bdspro_internal_repo.AmenityRepo), new(*bdspro_external_repo.GormAmenityRepo)),
	wire.Bind(new(bdspro_internal_repo.AssetCostRepo), new(*bdspro_external_repo.GormAssetCostRepo)),
	wire.Bind(new(bdspro_internal_repo.AssetCostTypeRepo), new(*bdspro_external_repo.GormAssetCostTypeRepo)),
	wire.Bind(new(bdspro_internal_repo.AssetExploitationRepo), new(*bdspro_external_repo.GormAssetExploitationRepo)),
	wire.Bind(new(bdspro_internal_repo.AssetIncomeTypeRepo), new(*bdspro_external_repo.GormAssetIncomeTypeRepo)),
	wire.Bind(new(bdspro_internal_repo.CostDocumentRepo), new(*bdspro_external_repo.GormCostDocumentRepo)),
	wire.Bind(new(bdspro_internal_repo.DepositeRepo), new(*bdspro_external_repo.GormDepositeRepo)),
	wire.Bind(new(bdspro_internal_repo.DocTypeRepo), new(*bdspro_external_repo.GormDocTypeRepo)),
	wire.Bind(new(bdspro_internal_repo.HouseInfoRepo), new(*bdspro_external_repo.GormHouseInfoRepo)),
	wire.Bind(new(bdspro_internal_repo.IncomeDocumentRepo), new(*bdspro_external_repo.GormIncomeDocumentRepo)),
	wire.Bind(new(bdspro_internal_repo.ProductAccessRepo), new(*bdspro_external_repo.GormProductAccessRepo)),
	wire.Bind(new(bdspro_internal_repo.ProductHistoryRepo), new(*bdspro_external_repo.PostgreHistoryRepo)),
	wire.Bind(new(bdspro_internal_repo.ProductPriceRepo), new(*bdspro_external_repo.GormProductPriceRepo)),
	wire.Bind(new(bdspro_internal_repo.ProductPrivateRepo), new(*bdspro_external_repo.GormProductPrivateRepo)),
	wire.Bind(new(bdspro_internal_repo.PropertyTypeRepo), new(*bdspro_external_repo.GormPropertyTypeRepo)),
	wire.Bind(new(bdspro_internal_repo.RegionRepo), new(*bdspro_external_repo.GormRegionRepo)),
	wire.Bind(new(bdspro_internal_repo.AssetRepo), new(*bdspro_external_repo_user.GormAssetRepo)),
	wire.Bind(new(bdspro_internal_repo.UPostRepo), new(*bdspro_external_repo_user.GormPostRepo)),
	wire.Bind(new(bdspro_internal_repo.UAssetHistoryRepo), new(*bdspro_external_repo.GormAssetHistoryRepo)),
	wire.Bind(new(bdspro_internal_repo.AssetLegalRepo), new(*bdspro_external_repo.GormAssetLegalRepo)),
	wire.Bind(new(bdspro_internal_repo.PostMediaRepo), new(*bdspro_external_repo.GormPostMediaRepo)),
	wire.Bind(new(bdspro_internal_repo.UProductRepo), new(*bdspro_external_repo.GormUserProductRepo)),
	wire.Bind(new(bdspro_internal_repo.UProductChildRepo), new(*bdspro_external_repo.GormUProductChildRepo)),
	wire.Bind(new(bdspro_internal_repo.ProductMediaRepo), new(*bdspro_external_repo.GormProductMediaRepo)),
	wire.Bind(new(bdspro_internal_usecases.PlanUsecase), new(*bdspro_external_services.MemberPlanUsecase)),
	wire.Bind(new(bdspro_v1_mod_user_product_interface.IUserProductService), new(*bdspro_v1_mod_user_product.UserProductService)),
	bdspro_v1_mod_user_action.NewActionRepo,
	bdspro_v1_mod_user_action.NewActionRouter,
	bdspro_v1_mod_user_action.NewActionService,
	bdspro_v1_mod_admin_product.NewAdminProductRepo,
	bdspro_v1_mod_admin_product.NewAdminProductRouter,
	bdspro_v1_mod_admin_product.NewAdminProductService,
	bdspro_v1_mod_admin.NewAdminRoute,
	bdspro_v1_mod_global_amenity.NewAmenityRepo,
	bdspro_v1_mod_admin_amenity.NewAmenityRouter,
	bdspro_v1_mod_admin_amenity.NewAmenityService,
	bdspro_v1_mod_global_apartment_attribute.NewApartmentAttrRepo,
	bdspro_v1_mod_admin_apartment.NewApartmentRepo,
	bdspro_v1_mod_admin_apartment.NewApartmentRouter,
	bdspro_v1_mod_admin_apartment.NewApartmentService,
	bdspro_v1_mod_user_asset_cost.NewAssetCostRepo,
	bdspro_v1_mod_user_asset_cost.NewAssetCostRouter,
	bdspro_external_services.NewAssetCostServer,
	bdspro_v1_mod_user_asset_cost.NewAssetCostService,
	bdspro_v1_mod_user_asset_cost_type.NewAssetCostTypeRepo,
	bdspro_v1_mod_user_asset_cost_type.NewAssetCostTypeRouter,
	bdspro_external_services.NewAssetCostTypeServer,
	bdspro_v1_mod_user_asset_cost_type.NewAssetCostTypeService,
	bdspro_internal_usecases.NewAssetCostTypeUc,
	bdspro_internal_usecases.NewAssetCostUc,
	bdspro_v1_mod_user_asset_exploitation.NewAssetExploitationRepo,
	bdspro_v1_mod_user_asset_exploitation.NewAssetExploitationRouter,
	bdspro_external_services.NewAssetExploitationServer,
	bdspro_v1_mod_user_asset_exploitation.NewAssetExploitationService,
	bdspro_internal_usecases.NewAssetExploitationUc,
	bdspro_v1_mod_global_history_asset.NewAssetHistoryRepo,
	bdspro_v1_mod_user_asset_income_type.NewAssetIncomeTypeRepo,
	bdspro_v1_mod_user_asset_income_type.NewAssetIncomeTypeRouter,
	bdspro_external_services.NewAssetIncomeTypeServer,
	bdspro_v1_mod_user_asset_income_type.NewAssetIncomeTypeService,
	bdspro_internal_usecases.NewAssetIncomeTypeUc,
	bdspro_v1_mod_user_asset_legal.NewAssetLegalRepo,
	bdspro_v1_mod_user_asset_legal.NewAssetLegalRouter,
	bdspro_external_services.NewAssetLegalServer,
	bdspro_v1_mod_user_asset_legal.NewAssetLegalService,
	bdspro_internal_usecases.NewAssetLegalUc,
	bdspro_v1_mod_user_asset.NewAssetRepo,
	bdspro_v1_mod_user_asset2.NewAssetRepo2,
	bdspro_v1_mod_user_asset.NewAssetRouter,
	bdspro_v1_mod_user_asset2.NewAssetRouter2,
	bdspro_v1_mod_user_asset2.NewAssetService2,
	bdspro_v1_mod_user_asset_share.NewAssetShareRepo,
	bdspro_v1_mod_user_asset_share.NewAssetShareRouter,
	bdspro_v1_mod_user_asset_share.NewAssetShareService,
	bdspro_external_services.NewAssetSplitHistoryServer,
	bdspro_internal_usecases.NewAssetSplitHistoryUc,
	bdspro_v1_mod_admin_build_room.NewBuildRoomRepo,
	bdspro_v1_mod_admin_build_room.NewBuildRoomRouter,
	bdspro_v1_mod_admin_build_room.NewBuildRoomService,
	bdspro_config.NewConfigApp,
	bdspro_external_services.NewCostDocumentServer,
	bdspro_internal_usecases.NewCostDocumentUc,
	bdspro_v1_mod_admin_developer.NewDeveloperRepo,
	bdspro_v1_mod_admin_developer.NewDeveloperRouter,
	bdspro_v1_mod_admin_developer.NewDeveloperService,
	bdspro_v1_mod_admin_district.NewDistrictRepo,
	bdspro_v1_mod_admin_district.NewDistrictRouter,
	bdspro_v1_mod_admin_district.NewDistrictService,
	bdspro_v1_mod_admin_doc_type.NewDocTypeRepo,
	bdspro_v1_mod_admin_doc_type.NewDocTypeRouter,
	bdspro_v1_mod_admin_doc_type.NewDocTypeService,
	bdspro_external_handler.NewGenericService,
	bdspro_external_repo.NewGormAmenityRepo,
	bdspro_external_repo.NewGormAssetCostRepo,
	bdspro_external_repo.NewGormAssetCostTypeRepo,
	bdspro_external_repo.NewGormAssetExploitationRepo,
	bdspro_external_repo.NewGormAssetIncomeTypeRepo,
	bdspro_external_repo.NewGormAssetLegalRepo,
	bdspro_external_repo_user.NewGormAssetRepo,
	bdspro_external_repo.NewGormAssetSplitHistoryRepo,
	bdspro_external_repo.NewGormCostDocumentRepo,
	bdspro_external_repo.NewGormDepositeRepo,
	bdspro_external_repo.NewGormDocTypeRepo,
	bdspro_external_repo.NewGormHouseInfoRepo,
	bdspro_external_repo.NewGormIncomeDocumentRepo,
	bdspro_external_repo.NewGormPostMediaRepo,
	bdspro_external_repo_user.NewGormPostRepo,
	bdspro_external_repo.NewGormProductMediaRepo,
	bdspro_external_repo.NewGormProductPriceRepo,
	bdspro_external_repo.NewGormProductPrivateRepo,
	bdspro_external_repo.NewGormPropertyTypeRepo,
	bdspro_external_repo.NewGormRegionRepo,
	bdspro_external_repo.NewGormUProductChildRepo,
	bdspro_external_repo.NewGormUserProductRepo,
	bdspro_external_repo.NewGromProductAccessRepo,
	bdspro_external_router_group.NewGroupProductRouter,
	bdspro_internal_usecases_group.NewGroupProductUsecase,
	bdspro_external_router_group.NewGroupRouter,
	bdspro_external_handler.NewGrpcService,
	bdspro_v1_mod_global_house_info.NewHouseInfoRepo,
	bdspro_external_services.NewIncomeDocumentServer,
	bdspro_internal_usecases.NewIncomeDocumentUc,
	bdspro_v1_mod_user_like.NewLikeRepo,
	bdspro_v1_mod_user_like.NewLikeRouter,
	bdspro_v1_mod_user_like.NewLikeService,
	bdspro_v1_mod_user_list.NewListRoute,
	bdspro_v1_providers.NewMemberPlanService,
	bdspro_external_services.NewMemberPlanUsecase,
	bdspro_v1_providers.NewNotiService,
	bdspro_external_router_org.NewOrgAccessRoute,
	bdspro_internal_usecases.NewOrgMemAccessUsecase,
	bdspro_external_router_org.NewOrgProductRouter,
	bdspro_internal_usecases_org.NewOrgProductUsecase,
	bdspro_external_router_org.NewOrgRouterManager,
	bdspro_v1_mod_user_post.NewPostRepo,
	bdspro_v1_mod_user_post.NewPostRouter,
	bdspro_external_repo.NewPostgreHistoryRepo,
	bdspro_v1_mod_global_product_access.NewProductAccessRepo,
	bdspro_v1_mod_admin_product_access.NewProductAccessRouter,
	bdspro_v1_mod_admin_product_access.NewProductAccessService,
	bdspro_v1_mod_user_product_child.NewProductChildService,
	bdspro_v1_mod_global_product_history.NewProductHistoryRepo,
	bdspro_v1_mod_global_product_media.NewProductMediaRepo,
	bdspro_v1_mod_global_product_price.NewProductPriceRepo,
	bdspro_v1_mod_global_product_private.NewProductPrivateRepo,
	bdspro_v1_mod_admin_project_build.NewProjectBuildRepo,
	bdspro_v1_mod_admin_project_build.NewProjectBuildRouter,
	bdspro_v1_mod_admin_project_build.NewProjectBuildService,
	bdspro_v1_mod_admin_project.NewProjectRepo,
	bdspro_v1_mod_admin_project.NewProjectRouter,
	bdspro_v1_mod_admin_project.NewProjectService,
	bdspro_v1_mod_admin_property_type.NewPropertyTypeRepo,
	bdspro_v1_mod_admin_property_type.NewPropertyTypeRouter,
	bdspro_v1_mod_admin_property_type.NewPropertyTypeService,
	bdspro_v1_mod_admin_province.NewProvinceRepo,
	bdspro_v1_mod_admin_province.NewProvinceRouter,
	bdspro_v1_mod_admin_province.NewProvinceService,
	bdspro_external_router_public.NewPublicRouter,
	bdspro_v1_providers.NewRedisProvider,
	bdspro_v1_mod_global_region.NewRegionRepo,
	bdspro_external_router_user.NewRouterManager,
	bdspro_v1_mod_user_transaction.NewTransactionRepo,
	bdspro_v1_mod_user_transaction.NewTransactionService,
	bdspro_external_providers.NewTransmitProvider,
	bdspro_external_router_user.NewUAssetLegalRoute,
	bdspro_external_router_user.NewUAssetRoute,
	bdspro_internal_usecases.NewUAssetUsecase,
	bdspro_external_router_user.NewUPostMediaRoute,
	bdspro_internal_usecases.NewUPostMediaUsecase,
	bdspro_internal_usecases.NewUPostUsecase,
	bdspro_external_router_user.NewUProductAccessRoute,
	bdspro_internal_usecases.NewUProductAccessUsecase,
	bdspro_external_router_user.NewUProductChildRoute,
	bdspro_internal_usecases.NewUProductChildUsecase,
	bdspro_external_providers.NewUProductProvider,
	bdspro_external_router_user.NewUProductRoute,
	bdspro_internal_usecases.NewUProductUsecase,
	bdspro_v1_mod_user_apartment.NewUserApartmentRepo,
	bdspro_v1_mod_user_apartment.NewUserApartmentRouter,
	bdspro_v1_mod_user_apartment.NewUserApartmentService,
	bdspro_v1_mod_user_asset.NewUserAssetService,
	bdspro_v1_mod_user_build_room.NewUserBuildRoomRepo,
	bdspro_v1_mod_user_build_room.NewUserBuildRoomRouter,
	bdspro_v1_mod_user_build_room.NewUserBuildRoomService,
	bdspro_v1_mod_user_list.NewUserListService,
	bdspro_external_router_user.NewUserPostRoute,
	bdspro_v1_mod_user_post.NewUserPostService,
	bdspro_v1_mod_user_product_access.NewUserProductAccessRoute,
	bdspro_v1_mod_user_product_access.NewUserProductAccessService,
	bdspro_v1_mod_user_product_child.NewUserProductChildRepo,
	bdspro_v1_mod_user_product_child.NewUserProductChildRoute,
	bdspro_v1_mod_user_product.NewUserProductRepo,
	bdspro_v1_mod_user_product.NewUserProductRoute,
	bdspro_v1_mod_user_product.NewUserProductService,
	bdspro_v1_mod_user_project_build.NewUserProjectBuildRepo,
	bdspro_v1_mod_user_project_build.NewUserProjectBuildRouter,
	bdspro_v1_mod_user_project_build.NewUserProjectBuildService,
	bdspro_v1_mod_user.NewUserRoute,
	bdspro_v1_mod_admin_ward.NewWardRepo,
	bdspro_v1_mod_admin_ward.NewWardRouter,
	bdspro_v1_mod_admin_ward.NewWardService,
)
func InitializeApp() (*config.ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}