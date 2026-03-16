//go:build wireinject
// +build wireinject

package wire

import (
	_db "common/db"
	tqd_infra_client "tqd/infra/client"
	tqd_infra_handler_grpc "tqd/infra/handler/grpc"
	tqd_infra_handler_http "tqd/infra/handler/http"
	tqd_infra_mapper "tqd/infra/mapper"
	tqd_infra_postgres "tqd/infra/postgres"
	tqd_infra_providers "tqd/infra/providers"
	tqd_infra_validator "tqd/infra/validator"
	"tqd/initial"
	tqd_initial "tqd/initial"
	tqd_internal_usecase "tqd/internal/usecase"
	tqd_map "tqd/map"

	"github.com/google/wire"
)

// Inject dependencies
var wireSet = wire.NewSet(
	_db.NewDB,
	_db.NewTransactionRepo,
	tqd_infra_providers.NewCacheProvider,
	tqd_infra_providers.NewFileProvider,
	tqd_infra_providers.NewGeoProvider,
	tqd_infra_providers.NewTransactionProvider,
	tqd_map.NewMapManager,
	tqd_infra_client.NewUserClient,
	tqd_infra_client.NewPermissionClient,
	tqd_infra_client.NewNotificationClient,
	
	// Amenity
	tqd_infra_postgres.NewAmenityRepo,
	tqd_infra_mapper.NewAmenityMapper,
	tqd_internal_usecase.NewAmenityUsecase,
	tqd_infra_handler_http.NewAmenityHandler,
	tqd_infra_handler_grpc.NewAmenityGrpcHandler,
	
	// ContactLabel
	tqd_infra_postgres.NewContactLabelPostgresRepo,
	tqd_infra_mapper.NewContactLabelMapper,
	tqd_infra_validator.NewContactLabelValidator,
	tqd_internal_usecase.NewContactLabelUsecase,
	tqd_infra_handler_http.NewContactLabelHandler,
	tqd_infra_handler_grpc.NewContactLabelGrpcHandler,
	
	// DirectoryCategory
	tqd_infra_postgres.NewDirectoryCategoryPostgresRepo,
	tqd_infra_mapper.NewDirectoryCategoryMapper,
	tqd_infra_validator.NewDirectoryCategoryValidator,
	tqd_internal_usecase.NewDirectoryCategoryUsecase,
	tqd_infra_handler_http.NewDirectoryCategoryHandler,
	tqd_infra_handler_grpc.NewDirectoryCategoryGrpcHandler,
	
	// DirectorySource
	tqd_infra_postgres.NewDirectorySourcePostgres,
	tqd_infra_mapper.NewDirectorySourceMapper,
	tqd_infra_validator.NewDirectorySourceValidator,
	tqd_internal_usecase.NewDirectorySourceUsecase,
	tqd_infra_handler_http.NewDirectorySourceHandler,
	tqd_infra_handler_grpc.NewDirectorySourceGrpcHandler,
	
	// DirectorySupplier
	tqd_infra_postgres.NewDirectorySupplierPostgres,
	tqd_infra_mapper.NewDirectorySupplierMapper,
	tqd_internal_usecase.NewDirectorySupplierUsecase,
	tqd_infra_handler_http.NewDirectorySupplierHandler,
	tqd_infra_handler_grpc.NewDirectorySupplierGrpcHandler,
	
	// OpenHour
	tqd_infra_postgres.NewOpenHourPostgres,
	tqd_infra_mapper.NewOpenHourMapper,
	tqd_internal_usecase.NewOpenHourUsecase,
	tqd_infra_handler_http.NewOpenHourHandler,
	tqd_infra_handler_grpc.NewOpenHourGrpcHandler,
	
	// PoiCategory
	tqd_infra_postgres.NewPoiCategoryPostgres,
	tqd_infra_mapper.NewPoiCategoryMapper,
	tqd_internal_usecase.NewPoiCategoryUsecase,
	tqd_infra_handler_http.NewPoiCategoryHandler,
	tqd_infra_handler_grpc.NewPoiCategoryGrpcHandler,
	
	// POI
	tqd_infra_postgres.NewPOIRepo,
	tqd_infra_mapper.NewPoiMapper,
	tqd_internal_usecase.NewPOIUsecase,
	tqd_infra_handler_http.NewPOIHTTPHandler,
	tqd_infra_handler_grpc.NewPoiGrpcHandler,
	
	// Base & App
	tqd_internal_usecase.NewBaseUsecase,
	tqd_initial.NewRedisService,
	tqd_initial.NewApp,
)
func InitializeApp() (*initial.InitialApp, error) {
	wire.Build(wireSet)
	return nil, nil
}