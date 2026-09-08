package cmd_grpc

import (
	"bdspro/config"
	"bdspro/infra/db"
	"bdspro/initial"
	"bdspro/wire"
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	bdspropb "pb/types/bdspro"
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

		app, cleanup, err := wire.InitializeApp()
		if err != nil {
			log.Fatalf("failed to initialize app: %v", err)
		}
		defer cleanup()

		if app == nil {
			log.Fatal("app is nil after InitializeApp")
		}
		if app.RedisClient != nil {
			defer func() {
				if err := app.RedisClient.Close(); err != nil {
					log.Printf("failed to close Redis client: %v", err)
				}
			}()
		}
		db.MigrateDomain()

		// Chạy IdentifierProperty với lineageID=23 lúc start (tạo property_identify và đính vào lineage + các bảng thuộc tính)
		ids := []uint64{26, 27, 28}
		for _, id := range ids {
			if id, err := app.PropertyHandler.PropertyUsecase.IdentifierProperty(context.Background(), id); err != nil {
				log.Printf("IdentifierProperty(lineageID=%d) failed: %v", id, err)
			} else {
				log.Printf("IdentifierProperty(lineageID=%d) ok, identifyID=%d", id, id)
			}
		}

		port := fmt.Sprintf(":%s", viper.GetString("server.tcp_port"))
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer(
			// grpc.UnaryInterceptor(common.ProfileIDInterceptor),
			grpc.UnaryInterceptor(_middleware.ParseGrpcMetadataContextMiddleware),
			grpc.StreamInterceptor(_middleware.ParseGrpcMetadataContextStreamMiddleware),
		)
		bdspropb.RegisterPostServiceServer(s, app.PostService)
		bdspropb.RegisterProductServiceServer(s, app.ProductService)
		bdspropb.RegisterBdsproPublicServiceServer(s, app.BdsproPublicService)
		bdspropb.RegisterAssetServiceServer(s, app.AssetService)
		bdspropb.RegisterAssetExploitationServiceServer(s, app.AssetExploitationService)
		bdspropb.RegisterAssetCostServiceServer(s, app.AssetCostService)
		bdspropb.RegisterAssetCostTypeServiceServer(s, app.AssetCostTypeService)
		bdspropb.RegisterSharingAccessServiceServer(s, app.SharingAccessService)
		bdspropb.RegisterAssetLegalServiceServer(s, app.AssetLegalService)
		bdspropb.RegisterProjectBuildServiceServer(s, app.ProjectBuildService)
		bdspropb.RegisterTransactionServiceServer(s, app.TransactionService)
		bdspropb.RegisterAdminProductServiceServer(s, app.AdminProductHandler)
		bdspropb.RegisterAdminAssetServiceServer(s, app.AdminAssetHandler)
		bdspropb.RegisterAdminPostServiceServer(s, app.AdminPostHandler)
		bdspropb.RegisterAdminPropertyTypeServiceServer(s, app.AdminPropertyTypeHandler)
		bdspropb.RegisterAdminProjectServiceServer(s, app.AdminProjectHandler)
		bdspropb.RegisterAdminAreaRegionServiceServer(s, app.AdminAreaRegionHandler)
		bdspropb.RegisterAdminRegionServiceServer(s, app.AdminRegionHandler)
		bdspropb.RegisterBdsproInternalServiceServer(s, app.InternalHandler)
		bdspropb.RegisterBdsproDashboardServiceServer(s, app.BdsproDashboardHandler)
		bdspropb.RegisterBdsDomainServiceServer(s, app.BdsDomainHandler)

		bdspropb.RegisterAdminTxServiceServer(s, app.AdminTxHandler)
		bdspropb.RegisterTxServiceServer(s, app.TxHandler)
		bdspropb.RegisterDealContractServiceServer(s, app.TxContractDealHandler)
		bdspropb.RegisterCostTypeServiceServer(s, app.TxCostTypeHandler)
		bdspropb.RegisterDealCostServiceServer(s, app.TxDealCostHandler)
		bdspropb.RegisterInternalTransactionServiceServer(s, app.TxInternalHandler)
		bdspropb.RegisterDealInternalNoteServiceServer(s, app.DealInternalNoteHandler)

		bdspropb.RegisterDealServiceServer(s, app.DealHandler)
		bdspropb.RegisterDealCommissionServiceServer(s, app.DealCommissionHandler)
		bdspropb.RegisterInvestmentServiceServer(s, app.DealInvestmentHandler)
		bdspropb.RegisterDealMemberServiceServer(s, app.DealInvitationHandler)
		bdspropb.RegisterDealMilestoneServiceServer(s, app.DealMilestoneHandler)

		bdspropb.RegisterPropertyServiceServer(s, app.PropertyHandler)
		bdspropb.RegisterDistributionServiceServer(s, app.DistributeHandler)
		bdspropb.RegisterProductNoteServiceServer(s, app.ProductNoteHandler)
		bdspropb.RegisterIdentifierServiceServer(s, app.PropertyIdentifier)
		startScheduler(app)
		// pb_generic.RegisterGenericServiceServer(s, app.GrpcGenericServer)

		// @bind: register_server
		// bdspropb.RegisterCostDocumentServiceServer(s, app.CostDocumentServer)
		// bdspropb.RegisterIncomeDocumentServiceServer(s, app.IncomeDocumentServer)
		// bdspropb.RegisterAssetIncomeTypeServiceServer(s, app.AssetIncomeTypeServer)
		// bdspropb.RegisterAssetExploitationServiceServer(s, app.AssetExploitationServer)
		// bdspropb.RegisterAssetCostTypeServiceServer(s, app.AssetCostTypeServer)
		// bdspropb.RegisterAssetCostServiceServer(s, app.AssetCostServer)
		// bdspropb.RegisterAssetLegalServiceServer(s, app.AssetLegalServer)
		// bdspropb.RegisterAssetSplitHistoryServiceServer(s, app.AssetSplitHistoryServer)

		log.Printf("Listen: %v", port)

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}

func startScheduler(app *initial.InitialApp) {
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
	// go app.AssetStatsJob.Run(context.Background())
	// go app.ProductStatsJob.Run(context.Background())
	// go app.PostStatsJob.Run(context.Background())

}
