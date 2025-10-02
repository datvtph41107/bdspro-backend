package initial

import (
	_redis "common/redis"
	handler_http "tqd/infra/handler/http"
)

// InitialApp holds all application dependencies
type InitialApp struct {
	// Handlers - chỉ giữ handlers vì cần để đăng ký routes
	ContactLabelHandler      *handler_http.ContactLabelHandler
	DirectoryCategoryHandler *handler_http.DirectoryCategoryHandler
	DirectorySourceHandler   *handler_http.DirectorySourceHandler
	DirectorySupplierHandler *handler_http.DirectorySupplierHandler
	PoiCategoryHandler       *handler_http.PoiCategoryHandler
	POIHandler               *handler_http.POIHTTPHandler
	OpenHourHandler          *handler_http.OpenHourHandler
	AmenityHandler           *handler_http.AmenityHandler
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
) *InitialApp {
	return &InitialApp{
		ContactLabelHandler:      contactLabelHandler,
		DirectoryCategoryHandler: directoryCategoryHandler,
		DirectorySourceHandler:   directorySourceHandler,
		DirectorySupplierHandler: directorySupplierHandler,
		PoiCategoryHandler:       poiCategoryHandler,
		POIHandler:               poiHandler,
		OpenHourHandler:          openHourHandler,
		AmenityHandler:           amenityHandler,
	}
}

func NewRedisService() *_redis.RedisService {
	return _redis.NewRedisService()
}
