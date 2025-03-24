//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	_db "common/db"
	bdspro_mod_admin "bdspro/mod/admin"
	bdspro_mod_admin_district "bdspro/mod/admin/district"
	bdspro_mod_admin_doc_type "bdspro/mod/admin/doc_type"
	bdspro_mod_admin_property_type "bdspro/mod/admin/property_type"
	bdspro_mod_admin_province "bdspro/mod/admin/province"
	bdspro_mod_admin_ward "bdspro/mod/admin/ward"
	bdspro_mod_global_product_access "bdspro/mod/global/product_access"
	bdspro_mod_user "bdspro/mod/user"
	bdspro_mod_user_asset "bdspro/mod/user/asset"
	bdspro_mod_user_post "bdspro/mod/user/post"
	bdspro_mod_user_product "bdspro/mod/user/product"
	bdspro_mod_user_transaction "bdspro/mod/user/transaction"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	NewConfigApp,
	bdspro_mod_admin.NewAdminRoute,
	bdspro_mod_user_asset.NewAssetRepo,
	bdspro_mod_user_asset.NewAssetService,
	bdspro_mod_admin_district.NewDistrictRepo,
	bdspro_mod_admin_district.NewDistrictRouter,
	bdspro_mod_admin_district.NewDistrictService,
	bdspro_mod_admin_doc_type.NewDocTypeRepo,
	bdspro_mod_admin_doc_type.NewDocTypeRouter,
	bdspro_mod_admin_doc_type.NewDocTypeService,
	bdspro_mod_user_post.NewPostRepo,
	bdspro_mod_user_post.NewPostService,
	bdspro_mod_global_product_access.NewProductAccessRepo,
	bdspro_mod_user_product.NewProductService,
	bdspro_mod_admin_property_type.NewPropertyTypeRepo,
	bdspro_mod_admin_property_type.NewPropertyTypeRouter,
	bdspro_mod_admin_property_type.NewPropertyTypeService,
	bdspro_mod_admin_province.NewProvinceRepo,
	bdspro_mod_admin_province.NewProvinceRouter,
	bdspro_mod_admin_province.NewProvinceService,
	bdspro_mod_user_transaction.NewTransactionRepo,
	bdspro_mod_user_transaction.NewTransactionService,
	bdspro_mod_user_product.NewUserProductRepo,
	bdspro_mod_user_product.NewUserProductRoute,
	bdspro_mod_user.NewUserRoute,
	bdspro_mod_admin_ward.NewWardRepo,
	bdspro_mod_admin_ward.NewWardRouter,
	bdspro_mod_admin_ward.NewWardService,
)
func InitializeApp() (*ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}