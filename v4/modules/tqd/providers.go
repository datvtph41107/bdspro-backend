package tqd

import (
	_db "common/db"

	"github.com/google/wire"
	infra_postgres "userv4/infra/postgres"
	"userv4/modules/tqd/controller"
	"userv4/modules/tqd/mapper"
	tqdpostgres "userv4/modules/tqd/postgres"
	discoverypg "userv4/modules/tqd/postgres/discovery"
	planningclientpg "userv4/modules/tqd/postgres/planningclient"
	publiccontentpg "userv4/modules/tqd/postgres/publiccontent"
	relatedpg "userv4/modules/tqd/postgres/relatedentity"
	"userv4/modules/tqd/provider"
	repo "userv4/modules/tqd/repository"
	"userv4/modules/tqd/services"
	discoveryapp "userv4/modules/tqd/services/discovery/application"
	discoveryports "userv4/modules/tqd/services/discovery/ports"
	genreportapp "userv4/modules/tqd/services/generatedreport/application"
	planningapp "userv4/modules/tqd/services/planningclient/application"
	qh_usecase "userv4/modules/tqd/services/qh"
	"userv4/modules/tqd/services/quota"
	relatedapp "userv4/modules/tqd/services/relatedentity/application"
	relatedports "userv4/modules/tqd/services/relatedentity/ports"
	"userv4/modules/tqd/services/resolver/layer_resolver/rule"
)

// ProviderSet wires the TQD client HTTP stack (repos + usecases + NewHandler).
// Postgres constructors take *gorm.DB; planning repos take *_db.TransactionRepo.
// Both are derived from the TQD-owned pool so they do not collide with the user DSN.
var ProviderSet = wire.NewSet(
	ProvideLocationRepository,
	ProvideLayerRepository,
	ProvideRegionRepository,
	ProvideRegionExtendRepository,
	ProvideQHLabelRepository,
	ProvideQHLayerLegendRepository,
	ProvideQHLandUseRepository,
	ProvideLayerLegalRepo,
	ProvideParcelRepo,
	ProvideLegalDocumentRepo,
	ProvideAuditRepository,
	ProvideLifecycleLayerRepository,
	ProvideQHLayerFamilyRepository,
	ProvidePlanningClientRepository,
	ProvidePublicContentRepository,
	ProvideQHPlanningProjectRepo,
	ProvideQHPlanningDocumentRepo,
	ProvideQHPlanningEventRepo,
	ProvideProAIJobRepo,
	ProvideMapWorkspaceRepo,
	ProvideDiscoveryRepository,
	ProvideRelatedEntityReader,
	ProvideRelatedEntityRepository,
	ProvideSubscriptionRepository,
	ProvideNotificationRepository,
	ProvideFeatureRepository,
	ProvideReportRepository,
	ProvideOneHouseRepo,
	ProvidePOIRepo,
	ProvideAmenityRepo,
	ProvideContactLabelRepo,
	ProvideOpenHourRepo,
	ProvidePoiCategoryRepo,
	ProvideDirectorySourceRepo,
	ProvideDirectorySupplierRepo,
	ProvideDirectoryCategoryRepo,
	ProvideQHAuthorityIssuringRepo,
	ProvideGeneratedReportService,
	ProvideProfileUsageService,

	provider.NewNoopUserProvider,
	provider.NewAllowAllPermissionProvider,
	mapper.NewPoiMapper,
	mapper.NewContactLabelMapper,
	mapper.NewOpenHourMapper,
	mapper.NewPoiCategoryMapper,
	mapper.NewDirectorySupplierMapper,
	services.NewLocalLocationService,
	wire.Bind(new(services.LocationService), new(*services.LocalLocationService)),
	services.NewLocationUsecase,
	services.NewLocalPublicService,
	wire.Bind(new(services.PublicService), new(*services.LocalPublicService)),
	services.NewLayerUsecase,
	services.NewRegionUsecase,
	services.NewRegionExtendUsecase,
	services.NewImportUsecase,
	services.NewParcelUsecase,
	services.NewQHLayerFamilyUsecase,
	services.NewQHLabelUsecase,
	services.NewQHLayerLandUseGroupUsecase,
	services.NewQHLayerLegendUsecase,
	services.NewQHAuthorityIssuringUsecase,
	services.NewMapWorkspaceUsecase,
	services.NewLayerLegalUsecase,
	services.NewSubscriptionUsecase,
	services.NewNotificationUsecase,
	services.NewFeatureUsecase,
	ProvideReportUsecase,
	services.NewOneHouseUsecase,
	services.NewPoiUsecase,
	services.NewAmenityUsecase,
	services.NewContactLabelUsecase,
	services.NewOpenHourUsecase,
	services.NewPoiCategoryUsecase,
	services.NewDirectorySourceUsecase,
	services.NewDirectorySupplierUsecase,
	services.NewDirectoryCategoryUsecase,
	discoveryapp.NewService,
	relatedapp.NewService,
	planningapp.NewService,
	wire.Bind(new(planningapp.PlanningRepository), new(*planningclientpg.Repository)),
	wire.Bind(new(planningapp.ProjectionRepository), new(*publiccontentpg.Repository)),
	qh_usecase.NewQHPlanningProjectUsecase,
	qh_usecase.NewQHPlanningEventUsecase,
	qh_usecase.NewQHPlanningDocumentUsecase,
	qh_usecase.NewProAIJobUsecase,
	controller.NewHandler,
)

func ProvideLocationRepository(tqd *infra_postgres.TqdDB) repo.LocationRepository {
	return tqdpostgres.NewLocationRepository(tqd.DB)
}

func ProvideLayerRepository(tqd *infra_postgres.TqdDB) repo.LayerRepository {
	return tqdpostgres.NewLayerRepository(tqd.DB)
}

func ProvideRegionRepository(tqd *infra_postgres.TqdDB) repo.RegionRepository {
	return tqdpostgres.NewRegionRepository(tqd.DB)
}

func ProvideRegionExtendRepository(tqd *infra_postgres.TqdDB) repo.RegionExtendRepository {
	return tqdpostgres.NewRegionExtendRepository(tqd.DB)
}

func ProvideQHLabelRepository(tqd *infra_postgres.TqdDB) repo.QHLabelRepository {
	return tqdpostgres.NewQHLabelRepository(tqd.DB)
}

func ProvideQHLayerLegendRepository(tqd *infra_postgres.TqdDB) repo.QHLayerLegendRepository {
	return tqdpostgres.NewQHLayerLegendRepository(tqd.DB)
}

func ProvideQHLandUseRepository(tqd *infra_postgres.TqdDB) repo.QHLandUseRepository {
	return tqdpostgres.NewQHLayerLandUseRepository(tqd.DB)
}

func ProvideLayerLegalRepo(tqd *infra_postgres.TqdDB) repo.ILayerLegalRepo {
	return tqdpostgres.NewLayerLegalPostgres(tqd.DB)
}

func ProvideParcelRepo(tqd *infra_postgres.TqdDB) repo.IParcelRepo {
	return tqdpostgres.NewParcelPostgres(tqd.DB)
}

func ProvideLegalDocumentRepo(tqd *infra_postgres.TqdDB) repo.ILegalDocumentRepo {
	return tqdpostgres.NewLegalDocumentPostgres(tqd.DB)
}

func ProvideAuditRepository(tqd *infra_postgres.TqdDB) rule.AuditRepository {
	return tqdpostgres.NewQHAuditPostgres(tqd.DB)
}

func ProvideLifecycleLayerRepository(tqd *infra_postgres.TqdDB) rule.LayerRepository {
	return tqdpostgres.NewQHLayerLifecyclePostgres(tqd.DB)
}

func ProvideQHLayerFamilyRepository(tqd *infra_postgres.TqdDB) repo.QHLayerFamilyRepository {
	return tqdpostgres.NewQHLayerFamilyRepository(tqd.DB)
}

func ProvidePlanningClientRepository(tqd *infra_postgres.TqdDB) *planningclientpg.Repository {
	return planningclientpg.NewRepository(tqd.DB)
}

func ProvidePublicContentRepository(tqd *infra_postgres.TqdDB) *publiccontentpg.Repository {
	return publiccontentpg.NewRepository(tqd.DB)
}

func ProvideTqdTransactionRepo(tqd *infra_postgres.TqdDB) *_db.TransactionRepo {
	return _db.NewTransactionRepo(tqd.DB)
}

func ProvideQHPlanningProjectRepo(tqd *infra_postgres.TqdDB) repo.IQHPlanningProjectRepo {
	return tqdpostgres.NewQHPlanningProjectRepo(ProvideTqdTransactionRepo(tqd))
}

func ProvideQHPlanningDocumentRepo(tqd *infra_postgres.TqdDB) repo.IQHPlanningDocumentRepo {
	return tqdpostgres.NewQHPlanningDocumentRepo(ProvideTqdTransactionRepo(tqd))
}

func ProvideQHPlanningEventRepo(tqd *infra_postgres.TqdDB) repo.IQHPlanningEventRepo {
	return tqdpostgres.NewQHPlanningEventRepo(ProvideTqdTransactionRepo(tqd))
}

func ProvideProAIJobRepo(tqd *infra_postgres.TqdDB) repo.IProAIJobRepo {
	return tqdpostgres.NewProAIJobRepo(ProvideTqdTransactionRepo(tqd))
}

func ProvideMapWorkspaceRepo(tqd *infra_postgres.TqdDB) repo.IMapWorkspaceRepo {
	return tqdpostgres.NewMapWorkspacePostgres(tqd.DB)
}

func ProvideDiscoveryRepository(tqd *infra_postgres.TqdDB) discoveryports.Repository {
	return discoverypg.NewDiscoveryRepository(tqd.DB)
}

func ProvideRelatedEntityReader(discoveryRepo discoveryports.Repository) relatedports.EntityReader {
	return relatedpg.NewEntityReader(discoveryRepo)
}

func ProvideRelatedEntityRepository(tqd *infra_postgres.TqdDB, reader relatedports.EntityReader) relatedports.Repository {
	return relatedpg.NewRepository(tqd.DB, reader)
}

func ProvideSubscriptionRepository(tqd *infra_postgres.TqdDB) repo.SubscriptionRepository {
	return tqdpostgres.NewSubscriptionRepository(tqd.DB)
}

func ProvideNotificationRepository(tqd *infra_postgres.TqdDB) repo.NotificationRepository {
	return tqdpostgres.NewNotificationRepository(tqd.DB)
}

func ProvideFeatureRepository(tqd *infra_postgres.TqdDB) repo.FeatureRepository {
	return tqdpostgres.NewFeatureRepository(tqd.DB)
}

func ProvideReportRepository(tqd *infra_postgres.TqdDB) repo.ReportRepository {
	return tqdpostgres.NewReportRepository(tqd.DB)
}

func ProvideOneHouseRepo(tqd *infra_postgres.TqdDB) repo.OneHouseRepo {
	return tqdpostgres.NewOneHouseRepo(tqd.DB)
}

func ProvidePOIRepo(tqd *infra_postgres.TqdDB) repo.IPOIRepo {
	return tqdpostgres.NewPOIRepo(ProvideTqdTransactionRepo(tqd))
}

func ProvideAmenityRepo(tqd *infra_postgres.TqdDB) repo.IAmenityRepo {
	return tqdpostgres.NewAmenityRepo(ProvideTqdTransactionRepo(tqd))
}

func ProvideContactLabelRepo(tqd *infra_postgres.TqdDB) repo.ContactLabelRepository {
	return tqdpostgres.NewContactLabelPostgresRepo(tqd.DB)
}

func ProvideOpenHourRepo(tqd *infra_postgres.TqdDB) repo.IOpenHourRepo {
	return tqdpostgres.NewOpenHourPostgres(ProvideTqdTransactionRepo(tqd))
}

func ProvidePoiCategoryRepo(tqd *infra_postgres.TqdDB) repo.IPoiCategoryRepo {
	return tqdpostgres.NewPoiCategoryPostgres(ProvideTqdTransactionRepo(tqd))
}

func ProvideDirectorySourceRepo(tqd *infra_postgres.TqdDB) repo.DirectorySourceRepository {
	return tqdpostgres.NewDirectorySourcePostgres(tqd.DB)
}

func ProvideDirectorySupplierRepo(tqd *infra_postgres.TqdDB) repo.DirectorySupplierRepository {
	return tqdpostgres.NewDirectorySupplierPostgres(tqd.DB)
}

func ProvideDirectoryCategoryRepo(tqd *infra_postgres.TqdDB) repo.DirectoryCategoryRepository {
	return tqdpostgres.NewDirectoryCategoryPostgresRepo(tqd.DB)
}

func ProvideQHAuthorityIssuringRepo(tqd *infra_postgres.TqdDB) repo.QHAuthorityIssuringRepository {
	return tqdpostgres.NewQHAuthorityIssuringRepository(tqd.DB)
}

func ProvideReportUsecase(reportRepo repo.ReportRepository) services.ReportUsecase {
	return services.NewReportUsecase(reportRepo, nil, nil, nil)
}

// ProvideGeneratedReportService is optional until full report runtime is wired.
func ProvideGeneratedReportService() *genreportapp.Service {
	return nil
}

// ProvideProfileUsageService is optional until access/durable/runtime stores are wired.
func ProvideProfileUsageService() *quota.ProfileUsageService {
	return nil
}
