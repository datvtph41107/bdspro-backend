package cmd_grpc

import (
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	pb_social "pb/types/social"
	"social/config"
	"social/db"
	"social/initial"
	"social/wire"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var GrpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "gRPC server",
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
		pb_social.RegisterNewsFeedServiceServer(s, app.NewsFeedService)
		pb_social.RegisterReportServiceServer(s, app.ReportService)
		pb_social.RegisterCommentServiceServer(s, app.CommentService)
		pb_social.RegisterLikeServiceServer(s, app.LikeService)

		log.Printf("GRPC Listen: %v", port)

		startScheduler(app)

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}

func startScheduler(app *initial.InitialApp) {
	now := time.Now()
	nextHour := now.Truncate(time.Minute).Add(30 * time.Minute) // 30p
	wait := time.Until(nextHour)
	app.NewsFeedUsecase.SyncNewsFeed(context.Background())

	time.AfterFunc(wait, func() {
		s := gocron.NewScheduler(time.Local)
		s.Every(1).Hour().Do(func() {
			app.NewsFeedUsecase.SyncNewsFeed(context.Background())
		})
		s.StartAsync()
	})
}
