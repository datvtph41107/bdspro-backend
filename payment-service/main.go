package main

import (
	"common/token"
	"log"
	"os"
	"payment/config"
	_ "payment/docs"
	"payment/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/payOSHQ/payos-lib-golang"
	"github.com/spf13/viper"
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
	"/payment/swagger",
	"/payment/payos/webhook/verify",
}

type ConfigApp struct {
	DB             *gorm.DB
	PaymentService *services.PaymentService
}

var RootApp *ConfigApp

func NewConfigApp(db *gorm.DB, fileService *services.PaymentService) *ConfigApp {
	return &ConfigApp{
		DB:             db,
		PaymentService: fileService,
	}
}

// @title Example API
// @version 1.0
// @description This is a sample API
// @BasePath /payment
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	r := gin.Default()
	config.LoadConfig()

	clientId := viper.GetString("payment.payos.client-id")
	apiKey := viper.GetString("payment.payos.api-key")
	checksumKey := viper.GetString("payment.payos.check-sum")
	payos.Key(clientId, apiKey, checksumKey)

	// Đặt múi giờ mặc định (VD: Asia/Ho_Chi_Minh)
	os.Setenv("TZ", "Europe/London")
	time.Local = time.FixedZone("Europe/London", 7*60*60)

	RootApp, _ = InitializeApp()
	// config.ConfigDB = RootApp.DB

	r.Use(token.JWTAuthMiddleware(publicRoutes))

	log.Printf("MAIN")
	r.GET("/payment/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router := r.Group("/payment/payos")
	router.POST("/webhook/verify", RootApp.PaymentService.VerifyPaymentWebhookData)
	router.POST("/create", RootApp.PaymentService.CreatePaymentLink)
	router.POST("/cancel/:orderId", RootApp.PaymentService.CancelPaymentLink)
	router.GET("/info/:orderId", RootApp.PaymentService.GetPaymentLinkInformation)

	r.Run(":" + viper.GetString("server.port"))
}
