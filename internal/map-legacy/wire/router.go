package wire

import (
	_jwt "common/jwt"
	"map/infra/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var publicRoutes = []string{
	"/map/swagger",
	"/map/locations/nearby",
	"/map/locations/find-by-polygon",
}

var tempRoutes = []string{}

// NewRouter creates a new router with all routes configured
func NewRouter(mapPointHandler *handler.MapPointHandler) *gin.Engine {
	r := gin.Default()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	// JWT middleware
	r.Use(_jwt.JWTAuthMiddleware(publicRoutes, tempRoutes))

	// Swagger documentation
	r.GET("/map/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/map/locations")
	{
		api.GET("/nearby", mapPointHandler.FindNearbyLocations)
		api.POST("/new", mapPointHandler.CreateMapPoint)
		api.POST("/find-by-polygon", mapPointHandler.FindLocationsByPolygon)
		api.GET("/:id", mapPointHandler.GetMapPointByID)
		api.PUT("/:id", mapPointHandler.UpdateMapPoint)
		api.DELETE("/:id", mapPointHandler.DeleteMapPoint)
	}

	return r
}
