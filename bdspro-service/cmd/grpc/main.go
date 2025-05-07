package cmd_grpc

import (
	"bdspro/config"
	"bdspro/v2/db"
	w "bdspro/wire"
	"common"
	_db "common/db"
	"fmt"
	"log"
	"net"
	"os"
	pb_bdspro "pb/types/bdspro"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var GrpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run the GRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadConfig()
		_db.NewDB()

		// Đặt múi giờ mặc định (VD: Asia/Ho_Chi_Minh)
		os.Setenv("TZ", "Europe/London")
		time.Local = time.FixedZone("Europe/London", 7*60*60)

		app, _ := w.InitializeApp()
		// _db.DB.AutoMigrate(&models.Order{})
		// config.ConfigDB = RootApp.DB
		db.MigrateDomain()

		// r.Use(_jwt.JWTAuthMiddleware(publicRoutes))

		// log.Printf("MAIN")
		// r.GET("/payment/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		// router := r.Group("/payment/payos")
		// router.POST("/webhook/verify", app.PaymentService.VerifyPaymentWebhookData)
		// router.POST("/create", app.PaymentService.CreatePaymentLink)
		// router.POST("/cancel/:orderId", app.PaymentService.CancelPaymentLink)
		// router.GET("/info/:orderId", app.PaymentService.GetPaymentLinkInformation)
		// pb_noti.RegisterNotificationServiceServer(s, app.ServerHttp)
		// group := r.Group("/payment/wallet")
		// RootApp.AccountRouter.RegisterRouter(group)
		// RootApp.TransactionRouter.RegisterRouter(group)

		port := fmt.Sprintf(":%s", viper.GetString("server.tcp_port"))
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer(
			grpc.UnaryInterceptor(common.ProfileIDInterceptor),
		)
		pb_bdspro.RegisterBdsproServiceServer(s, app.GrpcServer)
		log.Printf("Listen: %v", port)

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}
