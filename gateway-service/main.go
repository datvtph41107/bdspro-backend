package main

import (
	"fmt"
	cmd_grpc "gateway/cmd/grpc"
	cmd_http "gateway/cmd/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
}

// func main() {
// 	config.LoadConfig()

// 	// apiKey := viper.GetString("api.key")
// 	port := viper.GetString("server.port")

// 	// url := fmt.Sprintf("http://%s:%s")
// 	gateway := pkg.NewGateway("url")
// 	// gateway.FetchServices()

// 	r := gin.Default()

// 	// r.GET("/swagger/:service/*any", func(c *gin.Context) {
// 	// 	service := c.Param("service")

// 	// 	var swaggerPath string

// 	// 	switch service {
// 	// 	case "notification":
// 	// 		swaggerPath = "/swagger-doc/notification/notification.swagger.json"
// 	// 	case "bdspro":
// 	// 		swaggerPath = "/swagger-doc/bdspro/bdspro.swagger.json"
// 	// 	case "order":
// 	// 		swaggerPath = "/swagger-doc/order/order.swagger.json"
// 	// 	default:
// 	// 		c.JSON(404, gin.H{"error": "Service not found"})
// 	// 		return
// 	// 	}

// 	// 	swaggerURL := ginSwagger.URL(swaggerPath)
// 	// 	ginSwagger.WrapHandler(swaggerFile.Handler, swaggerURL)(c)
// 	// })
// 	// r.Static("/swagger-doc", "../protobuf/docs")
// 	//
// 	ctx := context.Background()
// 	muxGateway := runtime.NewServeMux(
// 		runtime.WithMetadata(_jwt.ExtractMetadataFromRequest),
// 	)

// 	// Đăng ký gRPC service vào HTTP proxy
// 	err := pb_noti.RegisterNotificationServiceHandlerFromEndpoint(ctx, muxGateway, "localhost:8004", []grpc.DialOption{grpc.WithInsecure()})
// 	if err != nil {
// 		log.Fatalf("Failed to register gRPC service: %v", err)
// 	}
// 	// err = pb_bdspro.RegisterBdsproServiceHandlerFromEndpoint(ctx, muxGateway, "localhost:8202", []grpc.DialOption{grpc.WithInsecure()})
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to register gRPC service: %v", err)
// 	// }

// 	// Danh sách các path public
// 	// publicPaths := []string{
// 	// 	"/api/v1/public/login",
// 	// 	"/api/v1/public/register",
// 	// }
// 	// Định tuyến các request gRPC đến HTTP proxy
// 	// r.Any("/v2/*any", gin.WrapH(muxGateway))

// 	// go func() {
// 	// 	if err := r.Run(":8080"); err != nil {
// 	// 		log.Fatalf("Gin server failed: %v", err)
// 	// 	}
// 	// }()

// 	r.Use(cors.New(cors.Config{
// 		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8000"}, // Cho phép gọi từ frontend
// 		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
// 		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
// 		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
// 		AllowCredentials: true,
// 		MaxAge:           12 * time.Hour,
// 	}))
// 	r.Any("/v1/:service/*path", gateway.ProxyRequest)

// 	// r.GET("/ping", func(c *gin.Context) {
// 	// 	c.JSON(200, gin.H{"message": "pong"})
// 	// })

// 	// log.Println("API Gateway started on port " + port)

// 	if err := r.Run(":" + port); err != nil {
// 		log.Fatalf("Gin server failed: %v", err)
// 	}
// }
