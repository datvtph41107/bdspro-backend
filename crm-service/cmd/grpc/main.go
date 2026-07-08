package cmd_grpc

import (
	_middleware "common/middleware"
	"context"
	"crm/config"
	"crm/infra/db"
	"crm/initial"
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
		crmpb.RegisterAppointmentServiceServer(s, app.AppointmentService)
		crmpb.RegisterSeoDomainServiceServer(s, app.SeoDomainHandler)
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