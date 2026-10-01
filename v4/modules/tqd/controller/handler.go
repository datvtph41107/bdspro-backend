package controller

import (
	"net/http"

	"userv4/modules/tqd/provider"
	"userv4/modules/tqd/services"
	discoveryapp "userv4/modules/tqd/services/discovery/application"
	genreportapp "userv4/modules/tqd/services/generatedreport/application"
	planningapp "userv4/modules/tqd/services/planningclient/application"
	qh_usecase "userv4/modules/tqd/services/qh"
	"userv4/modules/tqd/services/quota"
	relatedapp "userv4/modules/tqd/services/relatedentity/application"
	publiccontentpg "userv4/modules/tqd/postgres/publiccontent"
)

// Handler is the modules/tqd HTTP entry for client, public, admin, and internal routes.
type Handler struct {
	location       *LocationController
	public         *PublicController
	layer          *LayerController
	region         *RegionController
	parcel         *ParcelClientController
	family         *FamilyController
	label          *LabelController
	landUse        *LandUseController
	legend         *LegendController
	planningClient *PlanningClientController
	mapWorkspace   *MapWorkspaceController
	discovery      *DiscoveryController
	parcelUser     *ParcelUserController
	publicExt      *PublicExtendedController
	subscription   *SubscriptionController
	commercial     *CommercialController
	features       *FeatureController
	reports        *ReportController
	properties     *PropertyController

	qhAdminLayer     *QHAdminLayerController
	qhAdminFamily    *QHAdminFamilyController
	qhAdminLabel     *QHAdminLabelController
	qhAdminLandUse   *QHAdminLandUseController
	qhAdminLegend    *QHAdminLegendController
	qhAdminAuthority *QHAdminAuthorityController

	amenity      *AmenityController
	contactLabel *ContactLabelController
	openHour     *OpenHourController
	poiUser      *PoiUserController
	directory    *DirectoryController

	adminLocation    *AdminLocationController
	adminRegionImport *AdminRegionImportController
	adminParcelPoiSub *AdminParcelPoiSubController
	adminReport      *AdminReportController
	adminPlanning    *AdminPlanningController
	adminCommercial  *AdminCommercialController
	internalSEO      *InternalSEOController
}

func NewHandler(
	locationSvc services.LocationService,
	locationUC *services.LocationUsecase,
	publicSvc services.PublicService,
	layerUC services.LayerUsecase,
	regionUC services.RegionUsecase,
	parcelUC services.ParcelUsecase,
	familyUC services.QHLayerFamilyUsecase,
	labelUC services.QHLabelUsecase,
	landUseUC services.QHLayerLandUseGroupUsecase,
	legendUC services.QHLayerLegendUsecase,
	planningClient *planningapp.Service,
	qhProjects qh_usecase.IQHPlanningProjectUsecase,
	qhEvents qh_usecase.IQHPlanningEventUsecase,
	qhDocuments qh_usecase.IQHPlanningDocumentUsecase,
	workspaceUC services.MapWorkspaceUsecase,
	createReport *genreportapp.Service,
	discoverySvc *discoveryapp.Service,
	layerLegalUC *services.LayerLegalUsecase,
	poiUC services.PoiUsecase,
	relatedSvc *relatedapp.Service,
	publicContent *publiccontentpg.Repository,
	crmSeo provider.CrmSeoProvider,
	subUC services.SubscriptionUsecase,
	notifUC services.NotificationUsecase,
	profileUsage *quota.ProfileUsageService,
	featureUC *services.FeatureUsecase,
	reportUC services.ReportUsecase,
	oneHouseUC services.OneHouseUsecase,
	authorityUC services.QHAuthorityIssuringUsecase,
	amenityUC *services.AmenityUsecase,
	contactLabelUC services.ContactLabelUsecase,
	openHourUC services.OpenHourUsecase,
	poiCategoryUC services.PoiCategoryUsecase,
	directorySourceUC *services.DirectorySourceUsecase,
	directorySupplierUC *services.DirectorySupplierUsecase,
	directoryCategoryUC *services.DirectoryCategoryUsecase,
	regionExtendUC services.RegionExtendUsecase,
	importUC services.ImportUsecase,
	aiJobUC qh_usecase.IProAIJobUsecase,
	analytics provider.SearchAnalytics,
) *Handler {
	return &Handler{
		location:       NewLocationController(locationSvc, locationUC),
		public:         NewPublicController(publicSvc),
		layer:          NewLayerController(layerUC),
		region:         NewRegionController(regionUC),
		parcel:         NewParcelClientController(parcelUC),
		family:         NewFamilyController(familyUC),
		label:          NewLabelController(labelUC),
		landUse:        NewLandUseController(landUseUC),
		legend:         NewLegendController(legendUC),
		planningClient: NewPlanningClientController(planningClient, qhProjects, qhEvents, qhDocuments, analytics),
		mapWorkspace:   NewMapWorkspaceController(workspaceUC, createReport),
		discovery:      NewDiscoveryController(discoverySvc, analytics),
		parcelUser:     NewParcelUserController(parcelUC, layerLegalUC),
		publicExt: NewPublicExtendedController(
			parcelUC, locationUC, legendUC, poiUC, planningClient, relatedSvc, publicContent, crmSeo, analytics,
		),
		subscription: NewSubscriptionController(subUC, notifUC, parcelUC),
		commercial:   NewCommercialController(profileUsage),
		features:     NewFeatureController(featureUC),
		reports:      NewReportController(reportUC),
		properties:   NewPropertyController(oneHouseUC),

		qhAdminLayer:     NewQHAdminLayerController(layerUC),
		qhAdminFamily:    NewQHAdminFamilyController(familyUC),
		qhAdminLabel:     NewQHAdminLabelController(labelUC),
		qhAdminLandUse:   NewQHAdminLandUseController(landUseUC),
		qhAdminLegend:    NewQHAdminLegendController(legendUC),
		qhAdminAuthority: NewQHAdminAuthorityController(authorityUC),

		amenity:      NewAmenityController(amenityUC),
		contactLabel: NewContactLabelController(contactLabelUC),
		openHour:     NewOpenHourController(openHourUC),
		poiUser:      NewPoiUserController(poiUC, poiCategoryUC),
		directory:    NewDirectoryController(directorySourceUC, directorySupplierUC, directoryCategoryUC),

		adminLocation:     NewAdminLocationController(locationUC),
		adminRegionImport: NewAdminRegionImportController(regionUC, regionExtendUC, importUC),
		adminParcelPoiSub: NewAdminParcelPoiSubController(parcelUC, poiUC, subUC),
		adminReport:       NewAdminReportController(reportUC),
		adminPlanning:     NewAdminPlanningController(qhProjects, qhDocuments, qhEvents, aiJobUC),
		adminCommercial:   NewAdminCommercialController(),
		internalSEO:       NewInternalSEOController(locationUC, layerUC, parcelUC, qhProjects, qhDocuments),
	}
}

// RegisterHTTP mounts tqd routes with gateway path parity.
func (h *Handler) RegisterHTTP(mux *http.ServeMux) {
	h.location.RegisterHTTP(mux)
	h.public.RegisterHTTP(mux)
	h.layer.RegisterHTTP(mux)
	h.region.RegisterHTTP(mux)
	h.parcel.RegisterHTTP(mux)
	h.family.RegisterHTTP(mux)
	h.label.RegisterHTTP(mux)
	h.landUse.RegisterHTTP(mux)
	h.legend.RegisterHTTP(mux)
	h.planningClient.RegisterHTTP(mux)
	h.mapWorkspace.RegisterHTTP(mux)
	h.discovery.RegisterHTTP(mux)
	h.parcelUser.RegisterHTTP(mux)
	h.publicExt.RegisterHTTP(mux)
	h.subscription.RegisterHTTP(mux)
	h.commercial.RegisterHTTP(mux)
	h.features.RegisterHTTP(mux)
	h.reports.RegisterHTTP(mux)
	h.properties.RegisterHTTP(mux)

	h.qhAdminLayer.RegisterHTTP(mux)
	h.qhAdminFamily.RegisterHTTP(mux)
	h.qhAdminLabel.RegisterHTTP(mux)
	h.qhAdminLandUse.RegisterHTTP(mux)
	h.qhAdminLegend.RegisterHTTP(mux)
	h.qhAdminAuthority.RegisterHTTP(mux)

	h.amenity.RegisterHTTP(mux)
	h.contactLabel.RegisterHTTP(mux)
	h.openHour.RegisterHTTP(mux)
	h.poiUser.RegisterHTTP(mux)
	h.directory.RegisterHTTP(mux)

	h.adminLocation.RegisterHTTP(mux)
	h.adminRegionImport.RegisterHTTP(mux)
	h.adminParcelPoiSub.RegisterHTTP(mux)
	h.adminReport.RegisterHTTP(mux)
	h.adminPlanning.RegisterHTTP(mux)
	h.adminCommercial.RegisterHTTP(mux)
	h.internalSEO.RegisterHTTP(mux)
}
