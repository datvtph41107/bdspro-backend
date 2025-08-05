package cmd_grpc

import (
	"bdspro/config"
	"bdspro/db"
	"bdspro/wire"
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	pb_bdspro "pb/types/bdspro"
	"time"

	"github.com/go-co-op/gocron"
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
		db.MigrateDomain()

		port := fmt.Sprintf(":%s", viper.GetString("server.tcp_port"))
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer(
			// grpc.UnaryInterceptor(common.ProfileIDInterceptor),
			grpc.UnaryInterceptor(_middleware.ParseGrpcMetadataContextMiddleware),
		)
		pb_bdspro.RegisterPostServiceServer(s, app.PostService)
		pb_bdspro.RegisterProductServiceServer(s, app.ProductService)
		pb_bdspro.RegisterBdsproPublicServiceServer(s, app.BdsproPublicService)
		pb_bdspro.RegisterAssetServiceServer(s, app.AssetService)
		pb_bdspro.RegisterAssetExploitationServiceServer(s, app.AssetExploitationService)
		pb_bdspro.RegisterAssetCostServiceServer(s, app.AssetCostService)
		pb_bdspro.RegisterAssetCostTypeServiceServer(s, app.AssetCostTypeService)
		pb_bdspro.RegisterSharingAccessServiceServer(s, app.SharingAccessService)
		pb_bdspro.RegisterAssetLegalServiceServer(s, app.AssetLegalService)
		pb_bdspro.RegisterProjectBuildServiceServer(s, app.ProjectBuildService)
		pb_bdspro.RegisterTransactionServiceServer(s, app.TransactionService)

		startScheduler(app)
		// pb_generic.RegisterGenericServiceServer(s, app.GrpcGenericServer)

		// @bind: register_server
		// pb_bdspro.RegisterCostDocumentServiceServer(s, app.CostDocumentServer)
		// pb_bdspro.RegisterIncomeDocumentServiceServer(s, app.IncomeDocumentServer)
		// pb_bdspro.RegisterAssetIncomeTypeServiceServer(s, app.AssetIncomeTypeServer)
		// pb_bdspro.RegisterAssetExploitationServiceServer(s, app.AssetExploitationServer)
		// pb_bdspro.RegisterAssetCostTypeServiceServer(s, app.AssetCostTypeServer)
		// pb_bdspro.RegisterAssetCostServiceServer(s, app.AssetCostServer)
		// pb_bdspro.RegisterAssetLegalServiceServer(s, app.AssetLegalServer)
		// pb_bdspro.RegisterAssetSplitHistoryServiceServer(s, app.AssetSplitHistoryServer)

		log.Printf("Listen: %v", port)

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}

func startScheduler(app *config.ConfigApp) {
	log.Printf("startScheduler")
	now := time.Now()
	nextHour := now.Truncate(time.Minute).Add(30 * time.Minute) // 30p
	wait := time.Until(nextHour)
	app.ProductUsecase.SyncProduct(context.Background())

	time.AfterFunc(wait, func() {
		s := gocron.NewScheduler(time.Local)
		s.Every(1).Hour().Do(func() {
			app.ProductUsecase.SyncProduct(context.Background())
		})
		s.StartAsync()
	})
}
