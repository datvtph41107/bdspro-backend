package cmd_grpc

import (
	_middleware "common/middleware"
	"fmt"
	"log"
	"net"
	"os"
	"time"
	"utility/config"
	"utility/infra/db"
	"utility/wire"

	utilitypb "pb/types/utility"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var GrpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run the GRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadConfig()

		// Đặt múi giờ mặc định (VD: Asia/Ho_Chi_Minh)
		os.Setenv("TZ", "Europe/London")
		time.Local = time.FixedZone("Europe/London", 7*60*60)

		app, _ := wire.InitializeApp()
		db.AutoMigrate()

		port := fmt.Sprintf(":%s", viper.GetString("server.tcp_port"))
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		interceptors := grpc.ChainUnaryInterceptor(
			_middleware.ParseGrpcMetadataContextMiddleware,
			_middleware.UnaryRecoveryInterceptor(app.Logger),
		)
		s := grpc.NewServer(interceptors)
		utilitypb.RegisterEventQueueServiceServer(s, app.EventQueueService)
		utilitypb.RegisterUtilityInternalServiceServer(s, app.UtilityInternalService)

		log.Printf("Listen: %v", port)

		if err := s.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	},
}
