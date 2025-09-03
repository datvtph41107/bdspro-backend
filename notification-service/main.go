package main

import (
	"common"
	"fmt"
	"log"
	"net"
	"os"
	notificationpb "pb/types/notification"
	"time"

	"noti/config"
	"noti/db"
	_ "noti/docs"
	"noti/wire"

	_middleware "common/middleware"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var publicRoutes = []string{
	"/notification/swagger",
	"/notification/auth/otp",
}

// @title Notification API
// @version 1.0
// @description API phần logic thông báo
// @BasePath /v2/notification
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	common.LoadConfig(config.ConfigFile)

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

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			_middleware.ParseGrpcMetadataContextMiddleware,
		),
	)
	notificationpb.RegisterGatewayNotificationServiceServer(s, app.ServerHttp)
	notificationpb.RegisterNotificationServiceServer(s, app.InternalNotiService)
	notificationpb.RegisterHistoryServiceServer(s, app.HistoryService)

	log.Printf("GRPC Listen: %v", port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	// r := gin.Default()
	// os.Setenv("TZ", "Europe/London")

	// common.LoadConfig(config.ConfigFile)
	// port := fmt.Sprintf(":%s", viper.GetString("server.port"))

	// lis, err := net.Listen("tcp", port)
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }

	// s := grpc.NewServer()
	// notificationpb.RegisterNotificationServiceServer(s, &server.ServerHttp{})

	// // log.Fatalf("Server TCP listen on::  %v", port)
	// log.Printf("Server TCP listen on::  %v", port)

	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }

	// log.Fatalf("Server TCP listen on %v", port)

	// parameters := common.SetupParameters{
	// 	ConfigFile:   config.ConfigFile,
	// 	Engine:       r,
	// 	PublicRoutes: publicRoutes,
	// 	Prefix:       "/notification",
	// }

	// r := gin.Default()
	// // os.Setenv("TZ", "Europe/London")

	// parameters := common.SetupParameters{
	// 	ConfigFile:   config.ConfigFile,
	// 	Engine:       r,
	// 	PublicRoutes: publicRoutes,
	// 	Prefix:       "/notification",
	// }

	// // common.SetupRestfulAPI(parameters, func() {
	// // 	config.LoadProperties()
	// // 	log.Printf("Data: ", property.AppProperties)

	// // 	app.Noti.RegisterRoutes(r)
	// // })
	// common.SetupGRPCServer(parameters, func(s *grpc.Server) {
	// 	app, _ := InitializeApp()
	// 	db.AutoMigrate()
	// 	notificationpb.RegisterGatewayNotificationServiceServer(s, app.ServerHttp)
	// })

	// common.SetupRestfulAPI(parameters, func() {
	// 	config.LoadProperties()
	// 	log.Printf("Data: ", property.AppProperties)

	// 	app, _ := InitializeApp()
	// 	models.AutoMigrate()

	// 	app.Noti.RegisterRoutes(r)
	// })

	// common.LoadConfig(config.ConfigFile)
	// _db.NewDB()
	// port := fmt.Sprintf(":%s", viper.GetString("server.port"))

	// lis, err := net.Listen("tcp", port)
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }

	// s := grpc.NewServer()
	// notificationpb.RegisterNotificationServiceServer(s, &server.ServerHttp{})

	// log.Fatalf("Server TCP listen on %v", port)

	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }
}
