//go:build wireinject
// +build wireinject
	
package wire

import (
	"github.com/google/wire"
	_db "common/db"
	"tqd/initial"
	tqd_infra_client "tqd/infra/client"
	tqd_infra_handler_grpc "tqd/infra/handler/grpc"
	tqd_infra_handler_http "tqd/infra/handler/http"
	tqd_infra_mapper "tqd/infra/mapper"
	tqd_infra_postgres "tqd/infra/postgres"
	tqd_infra_providers "tqd/infra/providers"
	tqd_infra_validator "tqd/infra/validator"
	tqd_initial "tqd/initial"
	tqd_internal_usecase "tqd/internal/usecase"
	tqd_map "tqd/map"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	_db.NewTransactionRepo,
	tqd_infra_handler_grpc.NewAmenityGrpcHandlerImpl,
	tqd_infra_handler_http.NewAmenityHandler,
	tqd_infra_mapper.NewAmenityMapper,
	tqd_infra_postgres.NewAmenityRepo,
	tqd_internal_usecase.NewAmenityUsecase,
	tqd_initial.NewApp,
	tqd_internal_usecase.NewBaseUsecase,
	tqd_infra_providers.NewCacheProvider,
	tqd_infra_handler_grpc.NewContactLabelGrpcHandler,
	tqd_infra_handler_http.NewContactLabelHandler,
	tqd_infra_mapper.NewContactLabelMapper,
	tqd_infra_postgres.NewContactLabelPostgresRepo,
	tqd_internal_usecase.NewContactLabelUsecase,
	tqd_infra_validator.NewContactLabelValidator,
	tqd_infra_handler_grpc.NewDirectoryCategoryGrpcHandler,
	tqd_infra_handler_http.NewDirectoryCategoryHandler,
	tqd_infra_mapper.NewDirectoryCategoryMapper,
	tqd_infra_postgres.NewDirectoryCategoryPostgresRepo,
	tqd_internal_usecase.NewDirectoryCategoryUsecase,
	tqd_infra_validator.NewDirectoryCategoryValidator,
	tqd_infra_handler_grpc.NewDirectorySourceGrpcHandler,
	tqd_infra_handler_http.NewDirectorySourceHandler,
	tqd_infra_mapper.NewDirectorySourceMapper,
	tqd_infra_postgres.NewDirectorySourcePostgres,
	tqd_internal_usecase.NewDirectorySourceUsecase,
	tqd_infra_validator.NewDirectorySourceValidator,
	tqd_infra_handler_grpc.NewDirectorySupplierGrpcHandler,
	tqd_infra_handler_http.NewDirectorySupplierHandler,
	tqd_infra_mapper.NewDirectorySupplierMapper,
	tqd_infra_postgres.NewDirectorySupplierPostgres,
	tqd_internal_usecase.NewDirectorySupplierUsecase,
	tqd_infra_handler_grpc.NewFeatureGrpcHandler,
	tqd_infra_postgres.NewFeatureRepository,
	tqd_internal_usecase.NewFeatureUsecase,
	tqd_infra_providers.NewFileProvider,
	tqd_infra_providers.NewGeoProvider,
	tqd_infra_handler_grpc.NewLocationHandler,
	tqd_infra_mapper.NewLocationMapper,
	tqd_infra_postgres.NewLocationRepository,
	tqd_internal_usecase.NewLocationUsecase,
	tqd_infra_validator.NewLocationValidator,
	tqd_map.NewMapManager,
	tqd_infra_client.NewNotificationClient,
	tqd_infra_handler_grpc.NewOpenHourGrpcHandler,
	tqd_infra_handler_http.NewOpenHourHandler,
	tqd_infra_mapper.NewOpenHourMapper,
	tqd_infra_postgres.NewOpenHourPostgres,
	tqd_internal_usecase.NewOpenHourUsecase,
	tqd_infra_handler_http.NewPOIHTTPHandler,
	tqd_infra_postgres.NewPOIRepo,
	tqd_infra_client.NewPermissionClient,
	tqd_infra_handler_grpc.NewPoiCategoryGrpcHandler,
	tqd_infra_handler_http.NewPoiCategoryHandler,
	tqd_infra_mapper.NewPoiCategoryMapper,
	tqd_infra_postgres.NewPoiCategoryPostgres,
	tqd_internal_usecase.NewPoiCategoryUsecase,
	tqd_infra_handler_grpc.NewPoiGrpcHandler,
	tqd_infra_mapper.NewPoiMapper,
	tqd_internal_usecase.NewPoiUsecase,
	tqd_initial.NewRedisService,
	tqd_infra_providers.NewTransactionProvider,
	tqd_infra_client.NewUserClient,
)
func InitializeApp() (*initial.InitialApp, error) {
	wire.Build(wireSet)
	return nil, nil
}