package initial

import (
	_redis "common/redis"
	handler_grpc "tqd/infra/handler/grpc"
	handler_http "tqd/infra/handler/http"
)

// InitialApp holds all application dependencies
type InitialApp struct {
	// HTTP Handlers
	ContactLabelHandler      *handler_http.ContactLabelHandler
	DirectoryCategoryHandler *handler_http.DirectoryCategoryHandler
	DirectorySourceHandler   *handler_http.DirectorySourceHandler
	DirectorySupplierHandler *handler_http.DirectorySupplierHandler
	PoiCategoryHandler       *handler_http.PoiCategoryHandler
	POIHandler               *handler_http.POIHTTPHandler
	OpenHourHandler          *handler_http.OpenHourHandler
	AmenityHandler           *handler_http.AmenityHandler
	QHLabelHandler           *handler_http.QHLabelHandler
	RegionHttpHandler        *handler_http.RegionHttpHandler
	ImportHTTPHandler        *handler_http.ImportHTTPHandler

	// gRPC Handlers
	ContactLabelGrpcHandler        *handler_grpc.ContactLabelGrpcHandler
	DirectoryCategoryGrpcHandler   *handler_grpc.DirectoryCategoryGrpcHandler
	DirectorySourceGrpcHandler     *handler_grpc.DirectorySourceGrpcHandler
	DirectorySupplierGrpcHandler   *handler_grpc.DirectorySupplierGrpcHandler
	PoiCategoryGrpcHandler         *handler_grpc.PoiCategoryGrpcHandler
	POIGrpcHandler                 *handler_grpc.PoiGrpcHandler
	OpenHourGrpcHandler            *handler_grpc.OpenHourGrpcHandler
	AmenityGrpcHandler             *handler_grpc.AmenityGrpcHandlerImpl
	LocationGrpcHandler            *handler_grpc.LocationHandler
	FeatureGrpcHandler             *handler_grpc.FeatureGrpcHandler
	ParcelGrpcHandler              *handler_grpc.ParcelGrpcHandler
	OneHouseGrpcHandler            *handler_grpc.OneHouseGrpcHandler
	LayerGrpcHandler               *handler_grpc.LayerGrpcHandler
	QHLabelGrpcHandler             *handler_grpc.QHLabelGrpcHandler
	QHAuthorityIssuringGrpcHandler *handler_grpc.QHAuthorityIssuringGrpcHandler
	QHLayerFamilyGrpcHandler       *handler_grpc.QHLayerFamilyGrpcHandler
	QHLayerLegendGrpcHandler       *handler_grpc.QHLayerLegendGrpcHandler
	RegionGrpcHandler              *handler_grpc.RegionGrpcHandler
	ReportGrpcHandler              *handler_grpc.ReportGrpcHandler
	SubscriptionGrpcHandler        *handler_grpc.SubscriptionGrpcHandler
	TqdPublicService               *handler_grpc.TqdPublicService
	ImportGrpcHandler              *handler_grpc.ImportGrpcHandler
	LayerLegalGrpcHandler          *handler_grpc.LayerLegalGrpcHandler
}

// NewApp creates a new App instance
func NewApp(
	contactLabelHandler *handler_http.ContactLabelHandler,
	directoryCategoryHandler *handler_http.DirectoryCategoryHandler,
	directorySourceHandler *handler_http.DirectorySourceHandler,
	directorySupplierHandler *handler_http.DirectorySupplierHandler,
	poiCategoryHandler *handler_http.PoiCategoryHandler,
	poiHandler *handler_http.POIHTTPHandler,
	openHourHandler *handler_http.OpenHourHandler,
	amenityHandler *handler_http.AmenityHandler,
	qhLabelHandler *handler_http.QHLabelHandler,
	regionHttpHandler *handler_http.RegionHttpHandler,
	importHTTPHandler *handler_http.ImportHTTPHandler,
	contactLabelGrpcHandler *handler_grpc.ContactLabelGrpcHandler,
	directoryCategoryGrpcHandler *handler_grpc.DirectoryCategoryGrpcHandler,
	directorySourceGrpcHandler *handler_grpc.DirectorySourceGrpcHandler,
	directorySupplierGrpcHandler *handler_grpc.DirectorySupplierGrpcHandler,
	poiCategoryGrpcHandler *handler_grpc.PoiCategoryGrpcHandler,
	poiGrpcHandler *handler_grpc.PoiGrpcHandler,
	openHourGrpcHandler *handler_grpc.OpenHourGrpcHandler,
	amenityGrpcHandler *handler_grpc.AmenityGrpcHandlerImpl,
	locationGrpcHandler *handler_grpc.LocationHandler,
	featureGrpcHandler *handler_grpc.FeatureGrpcHandler,
	oneHouseGrpcHandler *handler_grpc.OneHouseGrpcHandler,
	layerGrpcHandler *handler_grpc.LayerGrpcHandler,
	qhLabelGrpcHandler *handler_grpc.QHLabelGrpcHandler,
	qhAuthorityIssuringGrpcHandler *handler_grpc.QHAuthorityIssuringGrpcHandler,
	qhLayerFamilyGrpcHandler *handler_grpc.QHLayerFamilyGrpcHandler,
	qhLayerLegendGrpcHandler *handler_grpc.QHLayerLegendGrpcHandler,
	pa *handler_grpc.ParcelGrpcHandler,
	regionGrpcHandler *handler_grpc.RegionGrpcHandler,
	reportGrpcHandler *handler_grpc.ReportGrpcHandler,
	subscriptionGrpcHandler *handler_grpc.SubscriptionGrpcHandler,
	tqdPublicService *handler_grpc.TqdPublicService,
	importGrpcHandler *handler_grpc.ImportGrpcHandler,
	layerLegalGrpcHandler *handler_grpc.LayerLegalGrpcHandler,
) *InitialApp {
	return &InitialApp{
		ContactLabelHandler:            contactLabelHandler,
		DirectoryCategoryHandler:       directoryCategoryHandler,
		DirectorySourceHandler:         directorySourceHandler,
		DirectorySupplierHandler:       directorySupplierHandler,
		PoiCategoryHandler:             poiCategoryHandler,
		POIHandler:                     poiHandler,
		OpenHourHandler:                openHourHandler,
		AmenityHandler:                 amenityHandler,
		QHLabelHandler:                 qhLabelHandler,
		RegionHttpHandler:              regionHttpHandler,
		ImportHTTPHandler:              importHTTPHandler,
		ContactLabelGrpcHandler:        contactLabelGrpcHandler,
		DirectoryCategoryGrpcHandler:   directoryCategoryGrpcHandler,
		DirectorySourceGrpcHandler:     directorySourceGrpcHandler,
		DirectorySupplierGrpcHandler:   directorySupplierGrpcHandler,
		PoiCategoryGrpcHandler:         poiCategoryGrpcHandler,
		POIGrpcHandler:                 poiGrpcHandler,
		OpenHourGrpcHandler:            openHourGrpcHandler,
		AmenityGrpcHandler:             amenityGrpcHandler,
		LocationGrpcHandler:            locationGrpcHandler,
		FeatureGrpcHandler:             featureGrpcHandler,
		ParcelGrpcHandler:              pa,
		OneHouseGrpcHandler:            oneHouseGrpcHandler,
		LayerGrpcHandler:               layerGrpcHandler,
		QHLabelGrpcHandler:             qhLabelGrpcHandler,
		QHAuthorityIssuringGrpcHandler: qhAuthorityIssuringGrpcHandler,
		QHLayerFamilyGrpcHandler:       qhLayerFamilyGrpcHandler,
		QHLayerLegendGrpcHandler:       qhLayerLegendGrpcHandler,
		RegionGrpcHandler:              regionGrpcHandler,
		ReportGrpcHandler:              reportGrpcHandler,
		SubscriptionGrpcHandler:        subscriptionGrpcHandler,
		TqdPublicService:               tqdPublicService,
		ImportGrpcHandler:              importGrpcHandler,
		LayerLegalGrpcHandler:          layerLegalGrpcHandler,
	}
}

func NewRedisService() *_redis.RedisService {
	return _redis.NewRedisService()
}
