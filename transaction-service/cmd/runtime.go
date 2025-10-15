package cmd_grpc

import (
	_middleware "common/middleware"
	"fmt"
	"log"
	"net"
	"os"
	transactionpb "pb/types/transaction"
	"time"
	"tx/config"
	"tx/db"
	"tx/wire"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var GrpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run the GRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadConfig()
		// _db.NewDB()

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
			// common.ProfileIDInterceptor,
			_middleware.ParseGrpcMetadataContextMiddleware,
			_middleware.UnaryRecoveryInterceptor(app.Logger),
		)
		s := grpc.NewServer(interceptors)
		transactionpb.RegisterInternalTransactionServiceServer(s, app.InternalTransactionHandler)
		transactionpb.RegisterTransactionServiceServer(s, app.TransactionHandler)
		transactionpb.RegisterDealContractServiceServer(s, app.ContractDealHandler)
		transactionpb.RegisterDealCostServiceServer(s, app.DealCostHandler)
		transactionpb.RegisterCostTypeServiceServer(s, app.CostTypeHandler)
		transactionpb.RegisterAdminTransactionServiceServer(s, app.AdminTransactionHandler)

		log.Printf("Listen: %v", port)

		if err := s.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	},
}
