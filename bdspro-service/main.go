package main

import (
	"fmt"
	"os"

	cmd_grpc "bdspro/cmd/grpc"
	cmd_http "bdspro/cmd/http"
	_ "bdspro/docs"
	_ "common/models"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// type ConfigApp struct {
// 	UserRoute   *user.UserRoute
// 	AdminRouter *admin.AdminRoute
// 	ServerHttp  *server.ServerHttp
// }

// var publicRoutes = []string{"/public"}

// @title BDSPRO API
// @version 1.0
// @description API phần logic chính của BĐSPro
// @BasePath /v1/bdspro/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

var rootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "Một ứng dụng quản lý server gRPC và HTTP",
}

func init() {
	// Đăng ký các subcommand vào rootCmd
	rootCmd.AddCommand(cmd_http.HttpCmd)
	rootCmd.AddCommand(cmd_grpc.GrpcCmd)

	// Cấu hình cho Cobra để lấy thông tin từ các tệp cấu hình (nếu cần)
	viper.AutomaticEnv() // Để lấy thông tin từ biến môi trường (nếu cần)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// engine := gin.Default()
	// os.Setenv("TZ", "Europe/London")
	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("notemptyarray", NotEmptyArray)
	// }

	// // common.LoadConfig(config.ConfigFile)

	// parameters := common.SetupParameters{
	// 	ConfigFile:   config.ConfigFile,
	// 	Engine:       engine,
	// 	PublicRoutes: publicRoutes,
	// 	Prefix:       "/bdspro",
	// }

	// app, _ := InitializeApp()
	// common.SetupRestfulAPI(parameters, func() {
	// 	config.LoadProperties()
	// 	log.Printf("Data: ", config.AppProperties)

	// 	models.AutoMigrate()

	// 	app.adminRouter.RegisterAdminRouter(engine)
	// 	app.userRoute.RegisterUserRouter(engine)
	// })

	// common.SetupGRPCServer(parameters, func(s *grpc.Server) {
	// 	config.LoadProperties()
	// 	log.Printf("Data: ", config.AppProperties)
	// 	pb_bdspro.RegisterBdsproServiceServer(s, app.ServerHttp)
	// })
}
