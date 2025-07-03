package cmd_grpc

import (
	"common"
	_db "common/db"
	_middleware "common/middleware"
	"context"
	"crm/config"
	"crm/infra/db"
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
		common.LoadConfig(config.ConfigFile)
		_db.NewDB()

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

		log.Printf("Listen: %v", port)

		startScheduler(app)

		if err := s.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	},
}

func startScheduler(app *config.ConfigApp) {
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
