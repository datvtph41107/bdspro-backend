//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	_db "common/db"
	noti_bdspro-service_v1_mod_admin "noti/bdspro-service/v1/mod/admin"
	noti_bdspro-service_v1_mod_admin_amenity "noti/bdspro-service/v1/mod/admin/amenity"
	noti_bdspro-service_v1_mod_admin_apartment "noti/bdspro-service/v1/mod/admin/apartment"
	noti_bdspro-service_v1_mod_admin_build_room "noti/bdspro-service/v1/mod/admin/build_room"
	noti_bdspro-service_v1_mod_admin_developer "noti/bdspro-service/v1/mod/admin/developer"
	noti_bdspro-service_v1_mod_admin_district "noti/bdspro-service/v1/mod/admin/district"
	noti_bdspro-service_v1_mod_admin_doc_type "noti/bdspro-service/v1/mod/admin/doc_type"
	noti_bdspro-service_v1_mod_admin_product "noti/bdspro-service/v1/mod/admin/product"
	noti_bdspro-service_v1_mod_admin_product_access "noti/bdspro-service/v1/mod/admin/product_access"
	noti_bdspro-service_v1_mod_admin_project "noti/bdspro-service/v1/mod/admin/project"
	noti_bdspro-service_v1_mod_admin_project_build "noti/bdspro-service/v1/mod/admin/project_build"
	noti_bdspro-service_v1_mod_admin_property_type "noti/bdspro-service/v1/mod/admin/property_type"
	noti_bdspro-service_v1_mod_admin_province "noti/bdspro-service/v1/mod/admin/province"
	noti_bdspro-service_v1_mod_admin_ward "noti/bdspro-service/v1/mod/admin/ward"
	noti_bdspro-service_v1_mod_global_amenity "noti/bdspro-service/v1/mod/global/amenity"
	noti_bdspro-service_v1_mod_global_apartment_attribute "noti/bdspro-service/v1/mod/global/apartment_attribute"
	noti_bdspro-service_v1_mod_global_house_info "noti/bdspro-service/v1/mod/global/house_info"
	noti_bdspro-service_v1_mod_global_member_plan "noti/bdspro-service/v1/mod/global/member_plan"
	noti_bdspro-service_v1_mod_global_product_access "noti/bdspro-service/v1/mod/global/product_access"
	noti_bdspro-service_v1_mod_global_product_media "noti/bdspro-service/v1/mod/global/product_media"
	noti_bdspro-service_v1_mod_global_product_price "noti/bdspro-service/v1/mod/global/product_price"
	noti_bdspro-service_v1_mod_global_product_private "noti/bdspro-service/v1/mod/global/product_private"
	noti_bdspro-service_v1_mod_global_region "noti/bdspro-service/v1/mod/global/region"
	noti_bdspro-service_v1_mod_user "noti/bdspro-service/v1/mod/user"
	noti_bdspro-service_v1_mod_user_asset "noti/bdspro-service/v1/mod/user/asset"
	noti_bdspro-service_v1_mod_user_list "noti/bdspro-service/v1/mod/user/list"
	noti_bdspro-service_v1_mod_user_post "noti/bdspro-service/v1/mod/user/post"
	noti_bdspro-service_v1_mod_user_product "noti/bdspro-service/v1/mod/user/product"
	noti_bdspro-service_v1_mod_user_product_access "noti/bdspro-service/v1/mod/user/product_access"
	noti_bdspro-service_v1_mod_user_product_child "noti/bdspro-service/v1/mod/user/product_child"
	noti_bdspro-service_v1_mod_user_transaction "noti/bdspro-service/v1/mod/user/transaction"
	noti_common_db "noti/common/db"
	noti_common_errors "noti/common/errors"
	noti_common_jwt "noti/common/jwt"
	noti_common_redis "noti/common/redis"
	noti_crm-service_data "noti/crm-service/data"
	noti_crm-service_repo "noti/crm-service/repo"
	noti_crm-service_routes "noti/crm-service/routes"
	noti_crm-service_services "noti/crm-service/services"
	noti_file-service_repositories "noti/file-service/repositories"
	noti_file-service_services "noti/file-service/services"
	noti_gateway-service_pkg "noti/gateway-service/pkg"
	noti_go-chat-service_infrastructure_database_oracle "noti/go-chat-service/infrastructure/database/oracle"
	noti_go-chat-service_infrastructure_database_postgres "noti/go-chat-service/infrastructure/database/postgres"
	noti_go-chat-service_infrastructure_redis "noti/go-chat-service/infrastructure/redis"
	noti_go-chat-service_internal_delivery_errors "noti/go-chat-service/internal/delivery/errors"
	noti_go-chat-service_internal_delivery_handlers "noti/go-chat-service/internal/delivery/handlers"
	noti_go-chat-service_internal_delivery_wshandler "noti/go-chat-service/internal/delivery/wshandler"
	noti_go-chat-service_internal_repository "noti/go-chat-service/internal/repository"
	noti_go-chat-service_internal_repository_orm "noti/go-chat-service/internal/repository/orm"
	noti_go-chat-service_internal_server "noti/go-chat-service/internal/server"
	noti_go-chat-service_internal_usecases "noti/go-chat-service/internal/usecases"
	noti_go-chat-service_proto_types_chat "noti/go-chat-service/proto/types/chat"
	noti_go-chat-service_proto_types_external "noti/go-chat-service/proto/types/external"
	noti_map-service "noti/map-service"
	noti_map-service_repo "noti/map-service/repo"
	noti_map-service_services "noti/map-service/services"
	noti_membership-service_repo "noti/membership-service/repo"
	noti_membership-service_routes "noti/membership-service/routes"
	noti_membership-service_services "noti/membership-service/services"
	noti_notification-service_config "noti/notification-service/config"
	noti_notification-service_data "noti/notification-service/data"
	noti_notification-service_repo "noti/notification-service/repo"
	noti_notification-service_routes "noti/notification-service/routes"
	noti_notification-service_server "noti/notification-service/server"
	noti_notification-service_services "noti/notification-service/services"
	noti_payment-service_repo "noti/payment-service/repo"
	noti_payment-service_services "noti/payment-service/services"
	noti_protobuf_types_bdspro "noti/protobuf/types/bdspro"
	noti_protobuf_types_notification "noti/protobuf/types/notification"
	noti_user-service_data "noti/user-service/data"
	noti_user-service_models "noti/user-service/models"
	noti_user-service_routes "noti/user-service/routes"
	noti_user-service_services "noti/user-service/services"
	noti_user-service_validator "noti/user-service/validator"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	NewConfigApp,
	noti_bdspro-service_v1_mod_admin_product.NewAdminProductRepo,
	noti_bdspro-service_v1_mod_admin_product.NewAdminProductRouter,
	noti_bdspro-service_v1_mod_admin_product.NewAdminProductService,
	noti_bdspro-service_v1_mod_admin.NewAdminRoute,
	noti_bdspro-service_v1_mod_global_amenity.NewAmenityRepo,
	noti_bdspro-service_v1_mod_admin_amenity.NewAmenityRouter,
	noti_bdspro-service_v1_mod_admin_amenity.NewAmenityService,
	noti_bdspro-service_v1_mod_global_apartment_attribute.NewApartmentAttrRepo,
	noti_bdspro-service_v1_mod_admin_apartment.NewApartmentRepo,
	noti_bdspro-service_v1_mod_admin_apartment.NewApartmentRouter,
	noti_bdspro-service_v1_mod_admin_apartment.NewApartmentService,
	noti_payment-service_services.NewApiClient,
	noti_common_errors.NewAppError,
	noti_bdspro-service_v1_mod_user_asset.NewAssetRepo,
	noti_bdspro-service_v1_mod_user_asset.NewAssetService,
	noti_user-service_data.NewAuthLoginResponse,
	noti_user-service_routes.NewAuthRoutes,
	noti_user-service_services.NewAuthService,
	noti_user-service_validator.NewAuthValidator,
	noti_go-chat-service_proto_types_external.NewBackgroundImageServiceClient,
	noti_go-chat-service_internal_usecases.NewBackgroundImageUsecases,
	noti_protobuf_types_bdspro.NewBdsproServiceClient,
	noti_crm-service_repo.NewBlockRepo,
	noti_crm-service_routes.NewBlockRouter,
	noti_crm-service_services.NewBlockService,
	noti_bdspro-service_v1_mod_admin_build_room.NewBuildRoomRepo,
	noti_bdspro-service_v1_mod_admin_build_room.NewBuildRoomRouter,
	noti_bdspro-service_v1_mod_admin_build_room.NewBuildRoomService,
	noti_go-chat-service_internal_delivery_handlers.NewChatHandler,
	noti_go-chat-service_proto_types_chat.NewChatServiceClient,
	noti_map-service.NewConfigApp,
	noti_crm-service_repo.NewContactRepo,
	noti_go-chat-service_internal_usecases.NewConversationUsecases,
	noti_user-service_services.NewCookieService,
	noti_common_db.NewDB,
	noti_go-chat-service_internal_delivery_handlers.NewDelivery,
	noti_bdspro-service_v1_mod_admin_developer.NewDeveloperRepo,
	noti_bdspro-service_v1_mod_admin_developer.NewDeveloperRouter,
	noti_bdspro-service_v1_mod_admin_developer.NewDeveloperService,
	noti_bdspro-service_v1_mod_admin_district.NewDistrictRepo,
	noti_bdspro-service_v1_mod_admin_district.NewDistrictRouter,
	noti_bdspro-service_v1_mod_admin_district.NewDistrictService,
	noti_bdspro-service_v1_mod_admin_doc_type.NewDocTypeRepo,
	noti_bdspro-service_v1_mod_admin_doc_type.NewDocTypeRouter,
	noti_bdspro-service_v1_mod_admin_doc_type.NewDocTypeService,
	noti_go-chat-service_internal_delivery_errors.NewError,
	noti_go-chat-service_internal_delivery_handlers.NewExternalHandler,
	noti_user-service_services.NewFacebookAuthService,
	noti_file-service_repositories.NewFileAccessRepository,
	noti_file-service_repositories.NewFileRepository,
	noti_file-service_services.NewFileService,
	noti_notification-service_config.NewFirebaseConfig,
	noti_notification-service_services.NewFirebaseService,
	noti_crm-service_repo.NewFollowRepo,
	noti_crm-service_routes.NewFollowRouter,
	noti_crm-service_services.NewFollowService,
	noti_crm-service_repo.NewFriendRepo,
	noti_crm-service_routes.NewFriendRouter,
	noti_crm-service_services.NewFriendService,
	noti_gateway-service_pkg.NewGateway,
	noti_user-service_services.NewGoogleAuthService,
	noti_go-chat-service_internal_repository_orm.NewGormRepository,
	noti_crm-service_data.NewGroupDTO,
	noti_crm-service_repo.NewGroupRepo,
	noti_crm-service_routes.NewGroupRouter,
	noti_crm-service_services.NewGroupService,
	noti_go-chat-service_internal_server.NewHTTPServer,
	noti_notification-service_repo.NewHistoryRepository,
	noti_bdspro-service_v1_mod_global_house_info.NewHouseInfoRepo,
	noti_user-service_services.NewHttpClient,
	noti_common_jwt.NewJwtTokenProperties,
	noti_bdspro-service_v1_mod_user_list.NewListRoute,
	noti_map-service_repo.NewMapRepo,
	noti_map-service_services.NewMapService,
	noti_bdspro-service_v1_mod_global_member_plan.NewMemberPlanService,
	noti_membership-service_repo.NewMembershipRepo,
	noti_membership-service_routes.NewMembershipRouter,
	noti_membership-service_services.NewMembershipService,
	noti_go-chat-service_internal_usecases.NewMessageReactionUsecases,
	noti_go-chat-service_internal_usecases.NewMessageUsecases,
	noti_notification-service_repo.NewNotificationRepository,
	noti_notification-service_routes.NewNotificationRouter,
	noti_notification-service_services.NewNotificationService,
	noti_notification-service_services.NewNotificationService2,
	noti_protobuf_types_notification.NewNotificationServiceClient,
	noti_user-service_routes.NewOAuthRouter,
	noti_user-service_services.NewOAuthService,
	noti_go-chat-service_infrastructure_database_oracle.NewOracleDB,
	noti_common_db.NewPage,
	noti_go-chat-service_internal_usecases.NewParticipantUsecases,
	noti_membership-service_services.NewPaymentClient,
	noti_payment-service_repo.NewPaymentRepo,
	noti_payment-service_services.NewPaymentService,
	noti_membership-service_repo.NewPlanRepo,
	noti_membership-service_routes.NewPlanRouter,
	noti_membership-service_services.NewPlanService,
	noti_bdspro-service_v1_mod_user_post.NewPostRepo,
	noti_bdspro-service_v1_mod_user_post.NewPostService,
	noti_go-chat-service_infrastructure_database_postgres.NewPostgresDB,
	noti_bdspro-service_v1_mod_global_product_access.NewProductAccessRepo,
	noti_bdspro-service_v1_mod_admin_product_access.NewProductAccessRouter,
	noti_bdspro-service_v1_mod_admin_product_access.NewProductAccessService,
	noti_bdspro-service_v1_mod_user_product_child.NewProductChildService,
	noti_bdspro-service_v1_mod_global_product_media.NewProductMediaRepo,
	noti_bdspro-service_v1_mod_global_product_price.NewProductPriceRepo,
	noti_bdspro-service_v1_mod_global_product_private.NewProductPrivateRepo,
	noti_bdspro-service_v1_mod_user_product.NewProductService,
	noti_user-service_routes.NewProfileRouter,
	noti_user-service_services.NewProfileService,
	noti_bdspro-service_v1_mod_admin_project_build.NewProjectBuildRepo,
	noti_bdspro-service_v1_mod_admin_project_build.NewProjectBuildRouter,
	noti_bdspro-service_v1_mod_admin_project_build.NewProjectBuildService,
	noti_bdspro-service_v1_mod_admin_project.NewProjectRepo,
	noti_bdspro-service_v1_mod_admin_project.NewProjectRouter,
	noti_bdspro-service_v1_mod_admin_project.NewProjectService,
	noti_bdspro-service_v1_mod_admin_property_type.NewPropertyTypeRepo,
	noti_bdspro-service_v1_mod_admin_property_type.NewPropertyTypeRouter,
	noti_bdspro-service_v1_mod_admin_property_type.NewPropertyTypeService,
	noti_bdspro-service_v1_mod_admin_province.NewProvinceRepo,
	noti_bdspro-service_v1_mod_admin_province.NewProvinceRouter,
	noti_bdspro-service_v1_mod_admin_province.NewProvinceService,
	noti_go-chat-service_internal_usecases.NewReadReceptUsecases,
	noti_go-chat-service_infrastructure_redis.NewRedisClient,
	noti_go-chat-service_internal_delivery_handlers.NewRedisClientWorker,
	noti_common_redis.NewRedisService,
	noti_user-service_data.NewRefreshResponse,
	noti_bdspro-service_v1_mod_global_region.NewRegionRepo,
	noti_go-chat-service_internal_repository.NewRepository,
	noti_user-service_data.NewRequestLoginResponse,
	noti_notification-service_data.NewSendNotiResponse,
	noti_notification-service_server.NewServerHttp,
	noti_file-service_services.NewStorageService,
	noti_go-chat-service_internal_server.NewSwaggerServe,
	noti_bdspro-service_v1_mod_user_transaction.NewTransactionRepo,
	noti_bdspro-service_v1_mod_user_transaction.NewTransactionService,
	noti_user-service_data.NewUserInfoResponse,
	noti_bdspro-service_v1_mod_user_list.NewUserListService,
	noti_user-service_models.NewUserOTPRepository,
	noti_bdspro-service_v1_mod_user_product_access.NewUserProductAccessRoute,
	noti_bdspro-service_v1_mod_user_product_access.NewUserProductAccessService,
	noti_bdspro-service_v1_mod_user_product_child.NewUserProductChildRepo,
	noti_bdspro-service_v1_mod_user_product_child.NewUserProductChildRoute,
	noti_bdspro-service_v1_mod_user_product.NewUserProductRepo,
	noti_bdspro-service_v1_mod_user_product.NewUserProductRoute,
	noti_user-service_models.NewUserProfileEntity,
	noti_user-service_models.NewUserProfileEntityWithRequest,
	noti_bdspro-service_v1_mod_user.NewUserRoute,
	noti_go-chat-service_internal_server.NewWSServer,
	noti_bdspro-service_v1_mod_admin_ward.NewWardRepo,
	noti_bdspro-service_v1_mod_admin_ward.NewWardRouter,
	noti_bdspro-service_v1_mod_admin_ward.NewWardService,
	noti_go-chat-service_internal_delivery_wshandler.NewWebSocketHandler,
	noti_user-service_services.NewZaloAuthService,
)
func InitializeApp() (*ConfigApp, error) {
	wire.Build(wireSet)
	return nil, nil
}