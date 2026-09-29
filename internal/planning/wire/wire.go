//go:build wireinject
// +build wireinject

package wire

import (
	_db "common/db"
	common_injection "common/injection"
	sharedredis "common/redis"
	common_utils "common/utils"
	"github.com/google/wire"
	redisv9 "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pb/clients"
	assistantpb "pb/types/assistant"
	userpb "pb/types/user"
	tqd_config "tqd/config"
	tqd_infra_client "tqd/infra/client"
	tqd_infra_client_file_generatedreport "tqd/infra/client/file/generatedreport"
	tqd_infra_client_user_generatedreport "tqd/infra/client/user/generatedreport"
	tqd_infra_handler_grpc "tqd/infra/handler/grpc"
	tqd_infra_handler_grpc_generatedreport "tqd/infra/handler/grpc/generatedreport"
	tqd_infra_handler_http "tqd/infra/handler/http"
	tqd_infra_handler_http_discovery "tqd/infra/handler/http/discovery"
	tqd_infra_handler_http_publiccontent "tqd/infra/handler/http/publiccontent"
	tqd_infra_mapper "tqd/infra/mapper"
	tqd_infra_postgres "tqd/infra/postgres"
	tqd_infra_postgres_discovery "tqd/infra/postgres/discovery"
	tqd_infra_postgres_generatedreport_job "tqd/infra/postgres/generatedreport/job"
	tqd_infra_postgres_generatedreport_render "tqd/infra/postgres/generatedreport/render"
	tqd_infra_postgres_generatedreport_report "tqd/infra/postgres/generatedreport/report"
	tqd_infra_postgres_mappoint "tqd/infra/postgres/mappoint"
	tqd_infra_postgres_planningclient "tqd/infra/postgres/planningclient"
	tqd_infra_postgres_publiccontent "tqd/infra/postgres/publiccontent"
	tqd_infra_postgres_quota "tqd/infra/postgres/quota"
	tqd_infra_postgres_relatedentity "tqd/infra/postgres/relatedentity"
	tqd_infra_postgres_usage "tqd/infra/postgres/usage"
	tqd_infra_providers "tqd/infra/providers"
	tqd_infra_redis_quota "tqd/infra/redis/quota"
	tqd_infra_validator "tqd/infra/validator"
	tqd_infra_worker "tqd/infra/worker"
	tqd_infra_worker_generatedreport_generator "tqd/infra/worker/generatedreport/generator"
	tqd_infra_worker_generatedreport_weasyprintpdf "tqd/infra/worker/generatedreport/weasyprintpdf"
	tqd_initial "tqd/initial"
	tqd_internal_domain_discovery_model "tqd/internal/domain/discovery/model"
	tqd_internal_domain_usage "tqd/internal/domain/usage"
	tqd_internal_dto "tqd/internal/dto"
	tqd_internal_usecase "tqd/internal/usecase"
	tqd_internal_usecase_discovery_application "tqd/internal/usecase/discovery/application"
	tqd_internal_usecase_generatedreport_admin "tqd/internal/usecase/generatedreport/admin"
	tqd_internal_usecase_generatedreport_application "tqd/internal/usecase/generatedreport/application"
	tqd_internal_usecase_generatedreport_memoryjob "tqd/internal/usecase/generatedreport/memoryjob"
	tqd_internal_usecase_generatedreport_processing "tqd/internal/usecase/generatedreport/processing"
	tqd_internal_usecase_generatedreport_rendering "tqd/internal/usecase/generatedreport/rendering"
	tqd_internal_usecase_mappoint "tqd/internal/usecase/mappoint"
	tqd_internal_usecase_planningclient_application "tqd/internal/usecase/planningclient/application"
	tqd_internal_usecase_qh "tqd/internal/usecase/qh"
	tqd_internal_usecase_quota "tqd/internal/usecase/quota"
	tqd_internal_usecase_quota_memory "tqd/internal/usecase/quota/memory"
	tqd_internal_usecase_quota_reconciliation_reservationcleanup "tqd/internal/usecase/quota/reconciliation/reservationcleanup"
	tqd_internal_usecase_quota_reconciliation_usageprojection "tqd/internal/usecase/quota/reconciliation/usageprojection"
	tqd_internal_usecase_relatedentity_application "tqd/internal/usecase/relatedentity/application"
	tqd_internal_usecase_resolver_layer_resolver "tqd/internal/usecase/resolver/layer_resolver"
	tqd_internal_usecase_resolver_layer_resolver_builder "tqd/internal/usecase/resolver/layer_resolver/builder"
	tqd_internal_usecase_resolver_layer_resolver_db "tqd/internal/usecase/resolver/layer_resolver/db"
	tqd_internal_usecase_resolver_layer_resolver_engine "tqd/internal/usecase/resolver/layer_resolver/engine"
	tqd_internal_usecase_resolver_layer_resolver_engine_parcel_engine "tqd/internal/usecase/resolver/layer_resolver/engine/parcel_engine"
	tqd_internal_usecase_resolver_layer_resolver_rule "tqd/internal/usecase/resolver/layer_resolver/rule"
	tqd_internal_usecase_usage "tqd/internal/usecase/usage"
	tqd_internal_usecase_usage_memory "tqd/internal/usecase/usage/memory"
	tqd_map "tqd/map"
)

func fileConfigFromRuntime(cfg *tqd_config.RuntimeConfig) tqd_infra_client.FileConfig {
	if cfg == nil {
		return tqd_infra_client.FileConfig{}
	}
	return tqd_infra_client.FileConfig{BaseURL: cfg.File.BaseURL, ServiceAuthKey: cfg.File.ServiceAuthKey, Timeout: cfg.File.Timeout}
}
func classifyPolicyFromRuntime(cfg *tqd_config.RuntimeConfig) tqd_internal_usecase_qh.ClassifyPolicy {
	if cfg == nil {
		return tqd_internal_usecase_qh.ClassifyPolicy{}
	}
	return tqd_internal_usecase_qh.ClassifyPolicy{ReadContent: cfg.Classify.ReadContent, MaxInlineBytes: cfg.Classify.MaxInlineBytes}
}
func classifyWorkerConfigFromRuntime(cfg *tqd_config.RuntimeConfig) tqd_infra_worker.ClassifyWorkerConfig {
	if cfg == nil {
		return tqd_infra_worker.ClassifyWorkerConfig{}
	}
	return tqd_infra_worker.ClassifyWorkerConfig{Enabled: cfg.Classify.Enabled, Interval: cfg.Classify.Interval, BatchSize: cfg.Classify.BatchSize}
}
func internalAccessClient(rpcClient *clients.UserGrpcClient) userpb.InternalAccessServiceClient {
	if rpcClient == nil {
		return nil
	}
	return rpcClient.AccessClient
}
func redisClientFromService(service *sharedredis.RedisService) *redisv9.Client {
	if service == nil {
		return nil
	}
	return service.Client
}
func usagePermissionAuthorizer(rpcClient *clients.AuthGrpcClient) tqd_infra_handler_grpc.UsagePermissionAuthorizer {
	return rpcClient
}

// Inject dependencies
var wireSet = wire.NewSet(
	fileConfigFromRuntime,
	classifyPolicyFromRuntime,
	classifyWorkerConfigFromRuntime,
	internalAccessClient,
	redisClientFromService,
	usagePermissionAuthorizer,
	_db.NewTransactionRepo,
	wire.Bind(new(tqd_internal_usecase_planningclient_application.PlanningRepository), new(*tqd_infra_postgres_planningclient.Repository)),
	wire.Bind(new(tqd_internal_usecase_planningclient_application.ProjectionRepository), new(*tqd_infra_postgres_publiccontent.Repository)),
	wire.Bind(new(tqd_internal_usecase_quota.AdminUsageRepository), new(*tqd_infra_postgres_quota.AdminUsageStore)),
	wire.Bind(new(tqd_internal_usecase_quota.RuntimeUsageReader), new(*tqd_infra_redis_quota.Store)),
	wire.Bind(new(tqd_internal_usecase_quota.ProfileAccessReader), new(*tqd_infra_client_user_generatedreport.Client)),
	wire.Bind(new(tqd_internal_usecase_quota.DurableUsageReader), new(*tqd_infra_postgres_usage.Store)),
	wire.Bind(new(tqd_internal_usecase_usage.Store), new(*tqd_infra_postgres_usage.Store)),
	wire.Bind(new(tqd_internal_usecase_generatedreport_admin.Repository), new(*tqd_infra_postgres_generatedreport_report.ReportStore)),
	wire.Bind(new(tqd_internal_usecase_quota_reconciliation_usageprojection.QuotaStore), new(*tqd_infra_redis_quota.Store)),
	common_injection.NewHttpClient,
	common_utils.NewSyncUtil,
	tqd_infra_client_file_generatedreport.New,
	tqd_infra_worker_generatedreport_weasyprintpdf.New,
	tqd_infra_handler_grpc_generatedreport.NewAdapter,
	tqd_infra_handler_grpc.NewAdminUsageGrpcHandler,
	tqd_internal_usecase_generatedreport_admin.NewService,
	tqd_internal_usecase_quota.NewAdminUsageService,
	tqd_infra_handler_grpc.NewProfileUsageGrpcHandler,
	tqd_internal_usecase_quota.NewProfileUsageService,
	tqd_infra_postgres_quota.NewAdminUsageStore,
	tqd_infra_handler_grpc.NewAmenityGrpcHandlerImpl,
	tqd_infra_handler_http.NewAmenityHandler,
	tqd_infra_mapper.NewAmenityMapper,
	tqd_infra_postgres.NewAmenityRepo,
	tqd_internal_usecase.NewAmenityUsecase,
	tqd_initial.NewApp,
	tqd_internal_usecase_resolver_layer_resolver_rule.NewAssessEngine,
	tqd_infra_client.NewAssistantClient,
	tqd_internal_usecase_resolver_layer_resolver_rule.NewAuditEngine,
	tqd_internal_usecase.NewBaseUsecase,
	tqd_infra_handler_grpc.NewBdsproPublicService,
	tqd_internal_usecase_resolver_layer_resolver_builder.NewCandidateBuilder,
	tqd_infra_client_user_generatedreport.NewClient,
	tqd_internal_usecase_resolver_layer_resolver_rule.NewCompareEngine,
	tqd_internal_usecase_resolver_layer_resolver_db.NewConfigLoader,
	tqd_internal_usecase_resolver_layer_resolver_db.NewConfigPostgres,
	tqd_internal_usecase_resolver_layer_resolver_rule.NewConflictClassifier,
	tqd_internal_usecase_resolver_layer_resolver_rule.NewConflictEngine,
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
	tqd_infra_handler_grpc.NewDiscoveryGrpcHandler,
	tqd_infra_mapper.NewDiscoveryMapper,
	tqd_infra_postgres_discovery.NewDiscoveryRepository,
	tqd_infra_postgres_relatedentity.NewEntityReader,
	tqd_internal_domain_discovery_model.NewEntityRef,
	tqd_internal_domain_usage.NewEvent,
	tqd_infra_handler_grpc.NewFeatureGrpcHandler,
	tqd_infra_postgres.NewFeatureRepository,
	tqd_internal_usecase.NewFeatureUsecase,
	tqd_infra_client.NewFileClient,
	tqd_infra_providers.NewGeoProvider,
	tqd_internal_dto.NewGeometryParser,
	tqd_infra_handler_grpc.NewGisReportAdminGrpcHandler,
	tqd_infra_handler_http_discovery.NewHandler,
	tqd_infra_handler_http_publiccontent.NewHandler,
	tqd_internal_usecase_resolver_layer_resolver_rule.NewHistoricalRiskEngine,
	tqd_infra_handler_grpc.NewImportGrpcHandler,
	tqd_infra_handler_http.NewImportHTTPHandler,
	tqd_infra_mapper.NewImportMapper,
	tqd_internal_usecase.NewImportUsecase,
	tqd_infra_mapper.NewLabelMapper,
	tqd_infra_handler_grpc.NewLayerGrpcHandler,
	tqd_infra_handler_grpc.NewLayerLegalGrpcHandler,
	tqd_infra_postgres.NewLayerLegalPostgres,
	tqd_internal_usecase.NewLayerLegalUsecase,
	tqd_infra_postgres.NewLayerRepository,
	tqd_config.NewLayerResolver,
	tqd_infra_postgres.NewLayerResolverConfigPostgres,
	tqd_internal_usecase.NewLayerUsecase,
	tqd_infra_postgres.NewLegalDocumentPostgres,
	tqd_internal_usecase_resolver_layer_resolver_rule.NewLegalStatusResolver,
	tqd_internal_usecase_resolver_layer_resolver_rule.NewLifecycleEngine,
	tqd_infra_handler_grpc.NewLocationHandler,
	tqd_infra_mapper.NewLocationMapper,
	tqd_infra_postgres.NewLocationRepository,
	tqd_internal_usecase.NewLocationUsecase,
	tqd_infra_validator.NewLocationValidator,
	tqd_map.NewMapManager,
	tqd_infra_handler_grpc.NewMapWorkspaceGrpcHandler,
	tqd_infra_postgres.NewMapWorkspacePostgres,
	tqd_internal_usecase.NewMapWorkspaceUsecase,
	tqd_infra_postgres.NewNotificationRepository,
	tqd_internal_usecase.NewNotificationUsecase,
	tqd_infra_handler_grpc.NewOneHouseGrpcHandler,
	tqd_infra_postgres.NewOneHouseRepo,
	tqd_internal_usecase.NewOneHouseUsecase,
	tqd_infra_handler_grpc.NewOpenHourGrpcHandler,
	tqd_infra_handler_http.NewOpenHourHandler,
	tqd_infra_mapper.NewOpenHourMapper,
	tqd_infra_postgres.NewOpenHourPostgres,
	tqd_internal_usecase.NewOpenHourUsecase,
	tqd_infra_handler_http.NewPOIHTTPHandler,
	tqd_infra_postgres.NewPOIRepo,
	tqd_internal_usecase_resolver_layer_resolver_engine_parcel_engine.NewParcelEngine,
	tqd_infra_handler_grpc.NewParcelGrpcHandler,
	tqd_infra_mapper.NewParcelMapper,
	tqd_infra_postgres.NewParcelPostgres,
	tqd_internal_usecase.NewParcelUsecase,
	tqd_infra_client.NewPermissionClient,
	tqd_infra_handler_grpc.NewPlanningClientGrpcHandler,
	tqd_infra_handler_grpc.NewPoiCategoryGrpcHandler,
	tqd_infra_handler_http.NewPoiCategoryHandler,
	tqd_infra_mapper.NewPoiCategoryMapper,
	tqd_infra_postgres.NewPoiCategoryPostgres,
	tqd_internal_usecase.NewPoiCategoryUsecase,
	tqd_infra_handler_grpc.NewPoiGrpcHandler,
	tqd_infra_mapper.NewPoiMapper,
	tqd_internal_usecase.NewPoiUsecase,
	tqd_internal_usecase_resolver_layer_resolver_rule.NewPriorityEngine,
	tqd_infra_postgres.NewProAIJobRepo,
	tqd_internal_usecase_qh.NewProAIJobUsecase,
	tqd_infra_postgres.NewQHAuditPostgres,
	tqd_infra_handler_grpc.NewQHAuthorityIssuringGrpcHandler,
	tqd_infra_postgres.NewQHAuthorityIssuringRepository,
	tqd_internal_usecase.NewQHAuthorityIssuringUsecase,
	tqd_infra_handler_grpc.NewQHLabelGrpcHandler,
	tqd_infra_handler_http.NewQHLabelHandler,
	tqd_infra_postgres.NewQHLabelRepository,
	tqd_internal_usecase.NewQHLabelUsecase,
	tqd_infra_handler_grpc.NewQHLandUseGrpcHandler,
	tqd_infra_handler_grpc.NewQHLayerFamilyGrpcHandler,
	tqd_infra_postgres.NewQHLayerFamilyRepository,
	tqd_internal_usecase.NewQHLayerFamilyUsecase,
	tqd_internal_usecase.NewQHLayerLandUseGroupUsecase,
	tqd_infra_postgres.NewQHLayerLandUseRepository,
	tqd_infra_handler_grpc.NewQHLayerLegendGrpcHandler,
	tqd_infra_postgres.NewQHLayerLegendRepository,
	tqd_internal_usecase.NewQHLayerLegendUsecase,
	tqd_infra_postgres.NewQHLayerLifecyclePostgres,
	tqd_internal_usecase_qh.NewQHPlanningClassifyUsecase,
	tqd_infra_worker.NewQHPlanningClassifyWorker,
	tqd_infra_postgres.NewQHPlanningDocumentRepo,
	tqd_internal_usecase_qh.NewQHPlanningDocumentUsecase,
	tqd_infra_postgres.NewQHPlanningEventRepo,
	tqd_internal_usecase_qh.NewQHPlanningEventUsecase,
	tqd_internal_usecase_qh.NewQHPlanningFolderUsecase,
	tqd_infra_handler_grpc.NewQHPlanningGrpcHandler,
	tqd_infra_mapper.NewQHPlanningMapper,
	tqd_infra_postgres.NewQHPlanningProjectRepo,
	tqd_internal_usecase_qh.NewQHPlanningProjectUsecase,
	tqd_infra_validator.NewQHPlanningValidator,
	tqd_infra_handler_grpc.NewRegionExtendGrpcHandler,
	tqd_infra_postgres.NewRegionExtendRepository,
	tqd_internal_usecase.NewRegionExtendUsecase,
	tqd_infra_handler_grpc.NewRegionGrpcHandler,
	tqd_infra_handler_http.NewRegionHttpHandler,
	tqd_infra_mapper.NewRegionMapper,
	tqd_infra_postgres.NewRegionRepository,
	tqd_internal_usecase.NewRegionUsecase,
	tqd_infra_handler_grpc.NewRelatedEntityGrpcHandler,
	tqd_infra_handler_grpc.NewReportGrpcHandler,
	tqd_infra_postgres.NewReportRepository,
	tqd_infra_postgres_generatedreport_report.NewReportStore,
	tqd_internal_usecase.NewReportUsecase,
	tqd_infra_postgres_planningclient.NewRepository,
	tqd_infra_postgres_publiccontent.NewRepository,
	tqd_infra_postgres_relatedentity.NewRepository,
	tqd_internal_usecase_resolver_layer_resolver_rule.NewResolutionPipeline,
	tqd_internal_usecase_resolver_layer_resolver.NewResolver,
	tqd_internal_usecase_resolver_layer_resolver_engine.NewResolverEngine,
	tqd_internal_usecase_resolver_layer_resolver.NewResolverWithDB,
	tqd_internal_usecase_generatedreport_processing.NewRunner,
	tqd_internal_usecase_quota_reconciliation_reservationcleanup.NewRunner,
	tqd_internal_usecase_quota_reconciliation_usageprojection.NewRunner,
	tqd_infra_handler_grpc.NewSeoProjectionGrpcHandler,
	tqd_internal_usecase_discovery_application.NewService,
	tqd_internal_usecase_generatedreport_application.NewService,
	tqd_internal_usecase_generatedreport_rendering.NewService,
	tqd_internal_usecase_planningclient_application.NewService,
	tqd_internal_usecase_quota.NewService,
	tqd_internal_usecase_quota_reconciliation_reservationcleanup.NewService,
	tqd_internal_usecase_quota_reconciliation_usageprojection.NewService,
	tqd_internal_usecase_relatedentity_application.NewService,
	tqd_infra_worker_generatedreport_generator.NewSmokeGenerator,
	tqd_infra_postgres_generatedreport_render.NewSource,
	tqd_infra_postgres_generatedreport_report.NewSourceStore,
	tqd_infra_postgres_generatedreport_job.NewStore,
	tqd_infra_postgres_usage.NewStore,
	tqd_infra_redis_quota.NewStore,
	tqd_internal_usecase_generatedreport_memoryjob.NewStore,
	tqd_internal_usecase_quota_memory.NewStore,
	tqd_internal_usecase_usage_memory.NewStore,
	tqd_infra_handler_grpc.NewSubscriptionGrpcHandler,
	tqd_infra_handler_grpc.NewMapPointGrpcHandler,
	tqd_infra_postgres_mappoint.NewRepository,
	tqd_internal_usecase_mappoint.NewService,
	wire.Bind(new(tqd_internal_usecase_mappoint.Repository), new(*tqd_infra_postgres_mappoint.Repository)),
	tqd_infra_postgres.NewSubscriptionRepository,
	tqd_internal_usecase.NewSubscriptionUsecase,
	tqd_internal_usecase_resolver_layer_resolver_rule.NewTimelineEngine,
	tqd_infra_providers.NewTransactionProvider,
	tqd_infra_client.NewUserClient,
	tqd_internal_usecase_generatedreport_processing.NewWorker,
	tqd_infra_mapper.NewWorkspaceMapper,
)

func InitializeApp(
	database *gorm.DB,
	redisService *sharedredis.RedisService,
	userRPC *clients.UserGrpcClient,
	authRPC *clients.AuthGrpcClient,
	assistantRPC assistantpb.AssistantServiceClient,
	runtimeConfig *tqd_config.RuntimeConfig,
) (*tqd_initial.InitialApp, func(), error) {
	wire.Build(wireSet)
	return nil, nil, nil
}
