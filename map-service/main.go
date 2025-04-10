package main

import (
	_jwt "common/jwt"
	"log"
	"map/config"
	_ "map/docs"
	"map/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// RegisterDocs registers the Swagger documentation.
// func RegisterDocs() {
// 	// Định nghĩa scheme cho Security
// 	swag.Register("Authorization", map[string]interface{}{
// 		"bearerAuth": map[string]interface{}{
// 			"type": "apiKey",
// 			"in":   "header",
// 			"name": "Authorization",
// 		},
// 	})
// }

var publicRoutes = []string{
	"/map/swagger",
	"/map/locations/nearby",
	"/map/locations/find-by-polygon",
}

type ConfigApp struct {
	DB         *gorm.DB
	MapService *services.MapService
}

var RootApp *ConfigApp

func NewConfigApp(db *gorm.DB, fileService *services.MapService) *ConfigApp {
	return &ConfigApp{
		DB:         db,
		MapService: fileService,
	}
}

// @title Example API
// @version 1.0
// @description This is a sample API
// @BasePath /map
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	r := gin.Default()
	config.LoadConfig()

	// Đặt múi giờ mặc định (VD: Asia/Ho_Chi_Minh)
	// os.Setenv("TZ", "Europe/London")
	// time.Local = time.FixedZone("Europe/London", 7*60*60)

	RootApp, _ = InitializeApp()
	// c_config.ConfigDB = RootApp.DB

	r.Use(_jwt.JWTAuthMiddleware(publicRoutes))

	log.Printf("MAIN")
	r.GET("/map/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router := r.Group("/map/locations")
	router.GET("/nearby", RootApp.MapService.HandleNearbyLocations)
	router.POST("/new", RootApp.MapService.HandleNewLocation)
	router.POST("/find-by-polygon", RootApp.MapService.HandleFindByPolygon)

	r.Run(":8101")
}
