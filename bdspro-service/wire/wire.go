//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"bdspro/config"
	bdspro_config "bdspro/config"
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
	bdspro_v1_mod_global_house_info "bdspro/v1/mod/global/house_info"
	bdspro_v1_mod_global_member_plan "bdspro/v1/mod/global/member_plan"
	bdspro_v1_mod_global_product_access "bdspro/v1/mod/global/product_access"
	bdspro_v1_mod_global_product_media "bdspro/v1/mod/global/product_media"
	bdspro_v1_mod_global_product_price "bdspro/v1/mod/global/product_price"
	bdspro_v1_mod_global_product_private "bdspro/v1/mod/global/product_private"
	bdspro_v1_mod_global_region "bdspro/v1/mod/global/region"
	bdspro_v1_mod_global_transmit "bdspro/v1/mod/global/transmit"
	bdspro_v1_mod_user "bdspro/v1/mod/user"
	bdspro_v1_mod_user_action "bdspro/v1/mod/user/action"
	bdspro_v1_mod_user_asset "bdspro/v1/mod/user/asset"
	bdspro_v1_mod_user_list "bdspro/v1/mod/user/list"
	bdspro_v1_mod_user_post "bdspro/v1/mod/user/post"
	bdspro_v1_mod_user_product "bdspro/v1/mod/user/product"
	bdspro_v1_mod_user_product_access "bdspro/v1/mod/user/product_access"
	bdspro_v1_mod_user_product_child "bdspro/v1/mod/user/product_child"
	bdspro_v1_mod_user_transaction "bdspro/v1/mod/user/transaction"
	bdspro_v2_external_handler "bdspro/v2/external/handler"
	bdspro_v2_external_repo "bdspro/v2/external/repo"
	bdspro_v2_external_service "bdspro/v2/external/service"
	bdspro_v2_repo "bdspro/v2/repo"
	bdspro_v2_usecases "bdspro/v2/usecases"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	wire.Bind(new(bdspro_v2_repo.AmenityRepo), new(*bdspro_v2_external_repo.GormAmenityRepo)),
	wire.Bind(new(bdspro_v2_repo.AssetRepo), new(*bdspro_v2_external_repo.GormAssetRepo)),
	wire.Bind(new(bdspro_v2_repo.AssetHistoryRepo), new(*bdspro_v2_external_repo.GormAssetHistoryRepo)),
	wire.Bind(new(bdspro_v2_repo.DocTypeRepo), new(*bdspro_v2_external_repo.GormDocTypeRepo)),
	wire.Bind(new(bdspro_v2_repo.HouseInfoRepo), new(*bdspro_v2_external_repo.GormHouseInfoRepo)),
	wire.Bind(new(bdspro_v2_repo.ProductAccessRepo), new(*bdspro_v2_external_repo.GormProductAccessRepo)),
	wire.Bind(new(bdspro_v2_repo.ProductMediaRepo), new(*bdspro_v2_external_repo.GormProductMediaRepo)),
	wire.Bind(new(bdspro_v2_repo.ProductPriceRepo), new(*bdspro_v2_external_repo.GormProductPriceRepo)),
	wire.Bind(new(bdspro_v2_repo.ProductPrivateRepo), new(*bdspro_v2_external_repo.GormProductPrivateRepo)),
	wire.Bind(new(bdspro_v2_repo.ProductUserRepo), new(*bdspro_v2_external_repo.GormUserProductRepo)),
	wire.Bind(new(bdspro_v2_repo.PropertyTypeRepo), new(*bdspro_v2_external_repo.GormPropertyTypeRepo)),
	wire.Bind(new(bdspro_v2_repo.RegionRepo), new(*bdspro_v2_external_repo.GormRegionRepo)),
	wire.Bind(new(bdspro_v2_usecases.PlanUsecase), new(*bdspro_v2_external_service.MemberPlanUsecase)),
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
	bdspro_v1_mod_user_asset.NewAssetRepo,
	bdspro_v1_mod_user_asset.NewAssetService,
	bdspro_v2_usecases.NewAssetUsecase,
	bdspro_v1_mod_admin_build_room.NewBuildRoomRepo,
	bdspro_v1_mod_admin_build_room.NewBuildRoomRouter,
	bdspro_v1_mod_admin_build_room.NewBuildRoomService,
	bdspro_config.NewConfigApp,
	bdspro_v1_mod_admin_developer.NewDeveloperRepo,
	bdspro_v1_mod_admin_developer.NewDeveloperRouter,
	bdspro_v1_mod_admin_developer.NewDeveloperService,
	bdspro_v1_mod_admin_district.NewDistrictRepo,
	bdspro_v1_mod_admin_district.NewDistrictRouter,
	bdspro_v1_mod_admin_district.NewDistrictService,
	bdspro_v1_mod_admin_doc_type.NewDocTypeRepo,
	bdspro_v1_mod_admin_doc_type.NewDocTypeRouter,
	bdspro_v1_mod_admin_doc_type.NewDocTypeService,
	bdspro_v2_external_repo.NewGormAmenityRepo,
	bdspro_v2_external_repo.NewGormAssetRepo,
	bdspro_v2_external_repo.NewGormDocTypeRepo,
	bdspro_v2_external_repo.NewGormHistoryRepo,
	bdspro_v2_external_repo.NewGormHouseInfoRepo,
	bdspro_v2_external_repo.NewGormProductMediaRepo,
	bdspro_v2_external_repo.NewGormProductPriceRepo,
	bdspro_v2_external_repo.NewGormProductPrivateRepo,
	bdspro_v2_external_repo.NewGormPropertyTypeRepo,
	bdspro_v2_external_repo.NewGormRegionRepo,
	bdspro_v2_external_repo.NewGormUserProductRepo,
	bdspro_v2_external_repo.NewGromProductAccessRepo,
	bdspro_v2_external_handler.NewGrpcService,
	bdspro_v2_usecases.NewHistoryUsecase,
	bdspro_v1_mod_global_house_info.NewHouseInfoRepo,
	bdspro_v1_mod_user_list.NewListRoute,
	bdspro_v1_mod_global_member_plan.NewMemberPlanService,
	bdspro_v2_external_service.NewMemberPlanUsecase,
	bdspro_v1_mod_global_transmit.NewNotiService,
	bdspro_v1_mod_user_post.NewPostRepo,
	bdspro_v1_mod_user_post.NewPostRouter,
	bdspro_v1_mod_user_post.NewPostService,
	bdspro_v1_mod_global_product_access.NewProductAccessRepo,
	bdspro_v1_mod_admin_product_access.NewProductAccessRouter,
	bdspro_v1_mod_admin_product_access.NewProductAccessService,
	bdspro_v1_mod_user_product_child.NewProductChildService,
	bdspro_v1_mod_global_product_media.NewProductMediaRepo,
	bdspro_v1_mod_global_product_price.NewProductPriceRepo,
	bdspro_v1_mod_global_product_private.NewProductPrivateRepo,
	bdspro_v1_mod_user_product.NewProductService,
	bdspro_v2_usecases.NewProductUserUsecase,
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
	bdspro_v1_mod_global_region.NewRegionRepo,
	bdspro_v1_mod_user_transaction.NewTransactionRepo,
	bdspro_v1_mod_user_transaction.NewTransactionService,
	bdspro_v1_mod_user_list.NewUserListService,
	bdspro_v1_mod_user_product_access.NewUserProductAccessRoute,
	bdspro_v1_mod_user_product_access.NewUserProductAccessService,
	bdspro_v1_mod_user_product_child.NewUserProductChildRepo,
	bdspro_v1_mod_user_product_child.NewUserProductChildRoute,
	bdspro_v1_mod_user_product.NewUserProductRepo,
	bdspro_v1_mod_user_product.NewUserProductRoute,
	bdspro_v1_mod_user.NewUserRoute,
	bdspro_v1_mod_admin_ward.NewWardRepo,
	bdspro_v1_mod_admin_ward.NewWardRouter,
	bdspro_v1_mod_admin_ward.NewWardService,
)
func InitializeApp() (*config.ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}