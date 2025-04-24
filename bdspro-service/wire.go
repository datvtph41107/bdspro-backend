//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	_db "common/db"
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
	bdspro_v1_mod_global_product_access "bdspro/v1/mod/global/product_access"
	bdspro_v1_mod_global_product_media "bdspro/v1/mod/global/product_media"
	bdspro_v1_mod_global_product_price "bdspro/v1/mod/global/product_price"
	bdspro_v1_mod_global_product_private "bdspro/v1/mod/global/product_private"
	bdspro_v1_mod_global_region "bdspro/v1/mod/global/region"
	bdspro_v1_mod_user "bdspro/v1/mod/user"
	bdspro_v1_mod_user_asset "bdspro/v1/mod/user/asset"
	bdspro_v1_mod_user_list "bdspro/v1/mod/user/list"
	bdspro_v1_mod_user_post "bdspro/v1/mod/user/post"
	bdspro_v1_mod_user_product "bdspro/v1/mod/user/product"
	bdspro_v1_mod_user_transaction "bdspro/v1/mod/user/transaction"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	NewConfigApp,
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
	bdspro_v1_mod_admin_build_room.NewBuildRoomRepo,
	bdspro_v1_mod_admin_build_room.NewBuildRoomRouter,
	bdspro_v1_mod_admin_build_room.NewBuildRoomService,
	bdspro_v1_mod_admin_developer.NewDeveloperRepo,
	bdspro_v1_mod_admin_developer.NewDeveloperRouter,
	bdspro_v1_mod_admin_developer.NewDeveloperService,
	bdspro_v1_mod_admin_district.NewDistrictRepo,
	bdspro_v1_mod_admin_district.NewDistrictRouter,
	bdspro_v1_mod_admin_district.NewDistrictService,
	bdspro_v1_mod_admin_doc_type.NewDocTypeRepo,
	bdspro_v1_mod_admin_doc_type.NewDocTypeRouter,
	bdspro_v1_mod_admin_doc_type.NewDocTypeService,
	bdspro_v1_mod_global_house_info.NewHouseInfoRepo,
	bdspro_v1_mod_user_list.NewListRoute,
	bdspro_v1_mod_user_post.NewPostRepo,
	bdspro_v1_mod_user_post.NewPostService,
	bdspro_v1_mod_global_product_access.NewProductAccessRepo,
	bdspro_v1_mod_admin_product_access.NewProductAccessRouter,
	bdspro_v1_mod_admin_product_access.NewProductAccessService,
	bdspro_v1_mod_global_product_media.NewProductMediaRepo,
	bdspro_v1_mod_global_product_price.NewProductPriceRepo,
	bdspro_v1_mod_global_product_private.NewProductPrivateRepo,
	bdspro_v1_mod_user_product.NewProductService,
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
	bdspro_v1_mod_user_product.NewUserProductRepo,
	bdspro_v1_mod_user_product.NewUserProductRoute,
	bdspro_v1_mod_user.NewUserRoute,
	bdspro_v1_mod_admin_ward.NewWardRepo,
	bdspro_v1_mod_admin_ward.NewWardRouter,
	bdspro_v1_mod_admin_ward.NewWardService,
)
func InitializeApp() (*ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}