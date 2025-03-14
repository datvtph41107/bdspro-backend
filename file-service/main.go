package main

import (
	c_config "common/config"
	cjwt "common/jwt"
	_ "file/docs"
	"file/services"
	"log"
	"os"
	"time"

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
	"/file/swagger",
	"/file/upload",
	"/file/load",
}

type ConfigApp struct {
	DB          *gorm.DB
	FileService *services.FileService
}

var RootApp *ConfigApp

func NewConfigApp(db *gorm.DB, fileService *services.FileService) *ConfigApp {
	return &ConfigApp{
		DB:          db,
		FileService: fileService,
	}
}

// @title Example API
// @version 1.0
// @description This is a sample API
// @BasePath /file
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	r := gin.Default()

	// Đặt múi giờ mặc định (VD: Asia/Ho_Chi_Minh)
	os.Setenv("TZ", "Europe/London")
	time.Local = time.FixedZone("Europe/London", 7*60*60)

	RootApp, _ = InitializeApp()
	c_config.ConfigDB = RootApp.DB

	r.Use(cjwt.JWTAuthMiddleware(publicRoutes))

	log.Printf("MAIN")
	r.GET("/file/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/file/load/:p/:s", RootApp.FileService.GetFile)
	r.POST("/file/upload", RootApp.FileService.Upload)
	r.PUT("/file/signature", RootApp.FileService.GenerateSignature)

	r.Run(":8002")
}
