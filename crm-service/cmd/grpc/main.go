package cmd_grpc

import (
	_middleware "common/middleware"
	process "common/process"
	"context"
	"crm/config"
	"crm/initial"
	"crm/wire"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	crmpb "pb/types/crm"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var GrpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run the GRPC server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGRPC()
	},
}

func runGRPC() error {
	if err := config.LoadConfig(); err != nil {
		return fmt.Errorf("load CRM config: %w", err)
	}
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

	app, cleanup, err := wire.InitializeApp()
	if err != nil {
		return fmt.Errorf("initialize CRM app: %w", err)
	}
	defer cleanup()

	port := fmt.Sprintf(":%s", viper.GetString("server.tcp_port"))
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("listen CRM gRPC %s: %w", port, err)
	}
	interceptors := grpc.ChainUnaryInterceptor(
		// common.ProfileIDInterceptor,
		_middleware.ParseGrpcMetadataContextMiddleware,
		_middleware.UnaryRecoveryInterceptor(),
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
	crmpb.RegisterCRMPublicContentServiceServer(s, app.PublicContentHandler)
	slog.
		// crmpb.RegisterRegistryServiceServer(s, app.RegistryHandler)
		Info(fmt.Sprintf("Listen: %v", port))

	signalCtx, stopSignals := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSignals()
	processCtx, processCancel := context.WithCancel(signalCtx)
	defer processCancel()

	actors := []process.Actor{
		process.ActorFunc(func(ctx context.Context) { runRuleEventScheduler(ctx, app) }),
		process.ActorFunc(func(ctx context.Context) { runSeoScheduler(ctx, app) }),
	}
	return process.Run(processCtx, processCancel, s, lis, actors, process.Config{GracefulStopTimeout: 10 * time.Second})
}

func runRuleEventScheduler(ctx context.Context, app *initial.InitialApp) {
	if app == nil || app.RuleEventUsecase == nil {
		<-ctx.Done()
		return
	}
	if err := app.RuleEventUsecase.TriggerJobAutoEvent(ctx); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("CRM rule event scheduler initial run failed: %v", err))
	}
	for {
		now := time.Now()
		nextHour := now.Truncate(time.Hour).Add(time.Hour)
		timer := time.NewTimer(time.Until(nextHour))
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			return
		case <-timer.C:
			if err := app.RuleEventUsecase.TriggerJobAutoEvent(ctx); err != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("CRM rule event scheduler run failed: %v", err))
			}
		}
	}
}
