package tqd

import (
	_db "common/db"

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

	"github.com/google/wire"
	"gorm.io/gorm"
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

func ProvideLocationRepository(db *gorm.DB) repo.LocationRepository {
	return tqdpostgres.NewLocationRepository(db)
}

func ProvideLayerRepository(db *gorm.DB) repo.LayerRepository {
	return tqdpostgres.NewLayerRepository(db)
}

func ProvideRegionRepository(db *gorm.DB) repo.RegionRepository {
	return tqdpostgres.NewRegionRepository(db)
}

func ProvideRegionExtendRepository(db *gorm.DB) repo.RegionExtendRepository {
	return tqdpostgres.NewRegionExtendRepository(db)
}

func ProvideQHLabelRepository(db *gorm.DB) repo.QHLabelRepository {
	return tqdpostgres.NewQHLabelRepository(db)
}

func ProvideQHLayerLegendRepository(db *gorm.DB) repo.QHLayerLegendRepository {
	return tqdpostgres.NewQHLayerLegendRepository(db)
}

func ProvideQHLandUseRepository(db *gorm.DB) repo.QHLandUseRepository {
	return tqdpostgres.NewQHLayerLandUseRepository(db)
}

func ProvideLayerLegalRepo(db *gorm.DB) repo.ILayerLegalRepo {
	return tqdpostgres.NewLayerLegalPostgres(db)
}

func ProvideParcelRepo(db *gorm.DB) repo.IParcelRepo {
	return tqdpostgres.NewParcelPostgres(db)
}

func ProvideLegalDocumentRepo(db *gorm.DB) repo.ILegalDocumentRepo {
	return tqdpostgres.NewLegalDocumentPostgres(db)
}

func ProvideAuditRepository(db *gorm.DB) rule.AuditRepository {
	return tqdpostgres.NewQHAuditPostgres(db)
}

func ProvideLifecycleLayerRepository(db *gorm.DB) rule.LayerRepository {
	return tqdpostgres.NewQHLayerLifecyclePostgres(db)
}

func ProvideQHLayerFamilyRepository(db *gorm.DB) repo.QHLayerFamilyRepository {
	return tqdpostgres.NewQHLayerFamilyRepository(db)
}

func ProvidePlanningClientRepository(db *gorm.DB) *planningclientpg.Repository {
	return planningclientpg.NewRepository(db)
}

func ProvidePublicContentRepository(db *gorm.DB) *publiccontentpg.Repository {
	return publiccontentpg.NewRepository(db)
}

func ProvideTqdTransactionRepo(db *gorm.DB) *_db.TransactionRepo {
	return _db.NewTransactionRepo(db)
}

func ProvideQHPlanningProjectRepo(db *gorm.DB) repo.IQHPlanningProjectRepo {
	return tqdpostgres.NewQHPlanningProjectRepo(ProvideTqdTransactionRepo(db))
}

func ProvideQHPlanningDocumentRepo(db *gorm.DB) repo.IQHPlanningDocumentRepo {
	return tqdpostgres.NewQHPlanningDocumentRepo(ProvideTqdTransactionRepo(db))
}

func ProvideQHPlanningEventRepo(db *gorm.DB) repo.IQHPlanningEventRepo {
	return tqdpostgres.NewQHPlanningEventRepo(ProvideTqdTransactionRepo(db))
}

func ProvideProAIJobRepo(db *gorm.DB) repo.IProAIJobRepo {
	return tqdpostgres.NewProAIJobRepo(ProvideTqdTransactionRepo(db))
}

func ProvideMapWorkspaceRepo(db *gorm.DB) repo.IMapWorkspaceRepo {
	return tqdpostgres.NewMapWorkspacePostgres(db)
}

func ProvideDiscoveryRepository(db *gorm.DB) discoveryports.Repository {
	return discoverypg.NewDiscoveryRepository(db)
}

func ProvideRelatedEntityReader(discoveryRepo discoveryports.Repository) relatedports.EntityReader {
	return relatedpg.NewEntityReader(discoveryRepo)
}

func ProvideRelatedEntityRepository(db *gorm.DB, reader relatedports.EntityReader) relatedports.Repository {
	return relatedpg.NewRepository(db, reader)
}

func ProvideSubscriptionRepository(db *gorm.DB) repo.SubscriptionRepository {
	return tqdpostgres.NewSubscriptionRepository(db)
}

func ProvideNotificationRepository(db *gorm.DB) repo.NotificationRepository {
	return tqdpostgres.NewNotificationRepository(db)
}

func ProvideFeatureRepository(db *gorm.DB) repo.FeatureRepository {
	return tqdpostgres.NewFeatureRepository(db)
}

func ProvideReportRepository(db *gorm.DB) repo.ReportRepository {
	return tqdpostgres.NewReportRepository(db)
}

func ProvideOneHouseRepo(db *gorm.DB) repo.OneHouseRepo {
	return tqdpostgres.NewOneHouseRepo(db)
}

func ProvidePOIRepo(db *gorm.DB) repo.IPOIRepo {
	return tqdpostgres.NewPOIRepo(ProvideTqdTransactionRepo(db))
}

func ProvideAmenityRepo(db *gorm.DB) repo.IAmenityRepo {
	return tqdpostgres.NewAmenityRepo(ProvideTqdTransactionRepo(db))
}

func ProvideContactLabelRepo(db *gorm.DB) repo.ContactLabelRepository {
	return tqdpostgres.NewContactLabelPostgresRepo(db)
}

func ProvideOpenHourRepo(db *gorm.DB) repo.IOpenHourRepo {
	return tqdpostgres.NewOpenHourPostgres(ProvideTqdTransactionRepo(db))
}

func ProvidePoiCategoryRepo(db *gorm.DB) repo.IPoiCategoryRepo {
	return tqdpostgres.NewPoiCategoryPostgres(ProvideTqdTransactionRepo(db))
}

func ProvideDirectorySourceRepo(db *gorm.DB) repo.DirectorySourceRepository {
	return tqdpostgres.NewDirectorySourcePostgres(db)
}

func ProvideDirectorySupplierRepo(db *gorm.DB) repo.DirectorySupplierRepository {
	return tqdpostgres.NewDirectorySupplierPostgres(db)
}

func ProvideDirectoryCategoryRepo(db *gorm.DB) repo.DirectoryCategoryRepository {
	return tqdpostgres.NewDirectoryCategoryPostgresRepo(db)
}

func ProvideQHAuthorityIssuringRepo(db *gorm.DB) repo.QHAuthorityIssuringRepository {
	return tqdpostgres.NewQHAuthorityIssuringRepository(db)
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
