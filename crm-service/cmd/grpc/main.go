package cmd_grpc

import (
	_db "common/db"
	_middleware "common/middleware"
	"context"
	"crm/config"
	"crm/infra/client"
	"crm/infra/db"
	handler_grpc "crm/infra/handler/grpc"
	"crm/infra/postgre"
	catalog_postgres "crm/infra/postgre/contentcatalog"
	"crm/initial"
	catalog_application "crm/internal/modules/contentcatalog/application"
	"crm/internal/usecase"
	"crm/wire"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	crmpb "pb/types/crm"

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
		// _db.NewDB()

		// Đặt múi giờ mặc định (VD: Asia/Ho_Chi_Minh)
		os.Setenv("TZ", "Europe/London")
		time.Local = time.FixedZone("Europe/London", 7*60*60)
		// timezone := viper.GetString("seo.timezone")
		// if timezone == "" {
		// 	timezone = "Asia/Ho_Chi_Minh"
		// }

		// loc, err := time.LoadLocation(timezone)
		// if err != nil {
		// 	log.Printf("invalid timezone %s, fallback Asia/Ho_Chi_Minh", timezone)
		// 	loc, _ = time.LoadLocation("Asia/Ho_Chi_Minh")
		// }

		// time.Local = loc

		app, _ := wire.InitializeApp()
		db.AutoMigrate()

		// Public SEO/content remains CRM-owned but is exposed only over gRPC.
		publicDB, publicDBErr := _db.NewDB()
		if publicDBErr != nil {
			log.Fatalf("initialize CRM public gRPC database: %v", publicDBErr)
		}
		seoDomainRepository := postgre.NewSeoDomainPostgres(publicDB)
		seoPublicUsecase := usecase.NewSeoPublicUsecase(seoDomainRepository)
		seoSitemapUsecase := usecase.NewSeoSitemapUsecase(seoDomainRepository)
		tqdClient := client.NewTqdClient()
		seoGenerationLogRepository := postgre.NewSeoGenerationLogPostgres(publicDB)
		seoRenderUsecase := usecase.NewSeoRenderUsecase(seoDomainRepository, seoGenerationLogRepository, tqdClient)
		seoProvisioningUsecase := usecase.NewSeoAdministrativeUnitProvisioningUsecase(seoDomainRepository, tqdClient, seoRenderUsecase)
		catalogRepository := catalog_postgres.NewRepository(publicDB)
		catalogService := catalog_application.NewService(catalogRepository)
		runtimeEventRepository := postgre.NewSeoRuntimeEventPostgres(publicDB)
		runtimeEventUsecase := usecase.NewSeoRuntimeEventUsecase(runtimeEventRepository)
		seoCountRepository := postgre.NewSeoCountPostgres(publicDB)
		seoCountUsecase := usecase.NewSeoCountUsecase(seoCountRepository, seoDomainRepository)
		publicContentGrpcHandler := handler_grpc.NewPublicContentGrpcHandler(seoPublicUsecase, seoSitemapUsecase, seoProvisioningUsecase, catalogService, runtimeEventUsecase, seoCountUsecase, config.AppProperties.Seo.RuntimeEventIngestKey)

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
		s := grpc.NewServer(
			interceptors,
			grpc.ChainStreamInterceptor(
				_middleware.ParseGrpcMetadataContextStreamMiddleware,
			),
		)
		crmpb.RegisterPipelineServiceServer(s, app.PipelineService)
		crmpb.RegisterContactServiceServer(s, app.ContactService)
		crmpb.RegisterLeadServiceServer(s, app.LeadService)
		crmpb.RegisterStageServiceServer(s, app.StageService)
		// crmpb.RegisterHistoryServiceServer(s, app.HistoryService)
		crmpb.RegisterRuleServiceServer(s, app.RuleService)
		crmpb.RegisterFollowServiceServer(s, app.FollowService)
		crmpb.RegisterBlockServiceServer(s, app.BlockService)
		crmpb.RegisterFriendGroupServiceServer(s, app.FriendGroupService)
		crmpb.RegisterFriendServiceServer(s, app.FriendService)
		crmpb.RegisterSharingAccessServiceServer(s, app.SharingAccessService)
		crmpb.RegisterInvitationInstallServiceServer(s, app.InvitationInstallService)
		crmpb.RegisterCrmInternalServiceServer(s, app.CrmInternalService)
		crmpb.RegisterEnumServiceServer(s, app.EnumService)
		// Feedback services
		crmpb.RegisterRateServiceServer(s, app.RateHandler)
		crmpb.RegisterReportServiceServer(s, app.ReportHandler)
		crmpb.RegisterReportAdminServiceServer(s, app.ReportAdminHandler)
		crmpb.RegisterReportReasonServiceServer(s, app.ReportReasonHandler)
		crmpb.RegisterSupportTicketServiceServer(s, app.SupportTicketHandler)
		crmpb.RegisterAdminSalesServiceServer(s, app.AdminSalesHandler)
		crmpb.RegisterAppointmentServiceServer(s, app.AppointmentService)
		crmpb.RegisterSeoDomainServiceServer(s, app.SeoDomainHandler)
		crmpb.RegisterSeoMeasurementRuntimeServiceServer(s, app.SeoDomainHandler)
		crmpb.RegisterSeoStrategyServiceServer(s, app.SeoStrategyHander)
		crmpb.RegisterCRMPublicContentServiceServer(s, publicContentGrpcHandler)
		// crmpb.RegisterRegistryServiceServer(s, app.RegistryHandler)

		log.Printf("Listen: %v", port)

		startScheduler(app)
		startSeoScheduler(app)

		if err := s.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	},
}

func startScheduler(app *initial.InitialApp) {
	now := time.Now()
	nextHour := now.Truncate(time.Hour).Add(1 * time.Hour) // giờ tròn kế tiếp
	wait := time.Until(nextHour)
	app.RuleEventUsecase.TriggerJobAutoEvent(context.Background())

	time.AfterFunc(wait, func() {
		s := gocron.NewScheduler(time.Local)
		s.Every(1).Hour().Do(func() {
			app.RuleEventUsecase.TriggerJobAutoEvent(context.Background())
		})
		s.StartAsync()
	})
}
