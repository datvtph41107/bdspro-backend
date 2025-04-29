package main

import (
	"common"
	"noti/config"
	_ "noti/docs"
	"noti/routes"
	"noti/server"
	pb_noti "pb/types/notification"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var publicRoutes = []string{
	"/notification/swagger",
	"/notification/auth/otp",
}

type ConfigApp struct {
	Noti       *routes.NotificationRouter
	ServerHttp *server.ServerHttp
}

func NewConfigApp(router *routes.NotificationRouter, ServerHttp *server.ServerHttp) *ConfigApp {
	return &ConfigApp{
		Noti:       router,
		ServerHttp: ServerHttp,
	}
}

// @title User API
// @version 1.0
// @description API quản lý người dùng
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// r := gin.Default()
	// os.Setenv("TZ", "Europe/London")

	// common.LoadConfig(config.ConfigFile)
	// port := fmt.Sprintf(":%s", viper.GetString("server.port"))

	// lis, err := net.Listen("tcp", port)
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }

	// s := grpc.NewServer()
	// pb_noti.RegisterNotificationServiceServer(s, &server.ServerHttp{})

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

	r := gin.Default()
	// os.Setenv("TZ", "Europe/London")

	parameters := common.SetupParameters{
		ConfigFile:   config.ConfigFile,
		Engine:       r,
		PublicRoutes: publicRoutes,
		Prefix:       "/notification",
	}

	// common.SetupRestfulAPI(parameters, func() {
	// 	config.LoadProperties()
	// 	log.Printf("Data: ", config.AppProperties)

	// 	models.AutoMigrate()

	// 	app.Noti.RegisterRoutes(r)
	// })
	common.SetupGRPCServer(parameters, func(s *grpc.Server) {
		app, _ := InitializeApp()
		pb_noti.RegisterNotificationServiceServer(s, app.ServerHttp)
	})

	// common.SetupRestfulAPI(parameters, func() {
	// 	config.LoadProperties()
	// 	log.Printf("Data: ", config.AppProperties)

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
	// pb_noti.RegisterNotificationServiceServer(s, &server.ServerHttp{})

	// log.Fatalf("Server TCP listen on %v", port)

	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }
}
