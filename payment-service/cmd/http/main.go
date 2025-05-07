package cmd_http

import (
	_db "common/db"
	_jwt "common/jwt"
	"log"
	"os"
	"payment/config"
	migrate "payment/db"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/payOSHQ/payos-lib-golang"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var HttpCmd = &cobra.Command{
	Use:   "http",
	Short: "Run the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()
		config.LoadConfig()
		_db.NewDB()

		clientId := viper.GetString("payment.payos.client-id")
		apiKey := viper.GetString("payment.payos.api-key")
		checksumKey := viper.GetString("payment.payos.check-sum")

		payos.Key(clientId, apiKey, checksumKey)

		// Đặt múi giờ mặc định (VD: Asia/Ho_Chi_Minh)
		os.Setenv("TZ", "Europe/London")
		time.Local = time.FixedZone("Europe/London", 7*60*60)

		// RootApp, _ := InitializeApp()
		// _db.DB.AutoMigrate(&models.Order{})
		// config.ConfigDB = RootApp.DB
		migrate.MigrateDomain()

		r.Use(_jwt.JWTAuthMiddleware(config.PublicRoutes))

		log.Printf("MAIN")
		// r.GET("/payment/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		// router := r.Group("/payment/payos")
		// router.POST("/webhook/verify", RootApp.PaymentService.VerifyPaymentWebhookData)
		// router.POST("/create", RootApp.PaymentService.CreatePaymentLink)
		// router.POST("/cancel/:orderId", RootApp.PaymentService.CancelPaymentLink)
		// router.GET("/info/:orderId", RootApp.PaymentService.GetPaymentLinkInformation)

		// group := r.Group("/payment/wallet")
		// RootApp.AccountRouter.RegisterRouter(group)
		// RootApp.TransactionRouter.RegisterRouter(group)
		r.Run(":" + viper.GetString("server.port"))

	},
}
