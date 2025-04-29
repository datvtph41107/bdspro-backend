package main

import (
	"gateway/config"
	"gateway/pkg"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Middleware CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("1GET_Request", c.Request.Header, c.Request.Method)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Cho phép tất cả domain (không nên dùng * trong production)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Lang") // Thêm 'Lang' vào đây
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")                          // Nếu dùng credentials (cookies, auth)
		log.Println("WRITE4")

		log.Println("2GET_Request", c.Request.Method == "OPTIONS")

		// Xử lý preflight request (OPTIONS)
		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Writer.WriteHeader(http.StatusOK)
			return
		}

		log.Println("GET_NEXT", c.Request.Method)
		c.Next()
	}
}

func main() {
	config.LoadConfig()

	// apiKey := viper.GetString("api.key")
	port := viper.GetString("server.port")

	// url := fmt.Sprintf("http://%s:%s")
	gateway := pkg.NewGateway("url")
	// gateway.FetchServices()

	r := gin.Default()

	// r.GET("/swagger/:service/*any", func(c *gin.Context) {
	// 	service := c.Param("service")

	// 	var swaggerPath string

	// 	switch service {
	// 	case "notification":
	// 		swaggerPath = "/swagger-doc/notification/notification.swagger.json"
	// 	case "bdspro":
	// 		swaggerPath = "/swagger-doc/bdspro/bdspro.swagger.json"
	// 	case "order":
	// 		swaggerPath = "/swagger-doc/order/order.swagger.json"
	// 	default:
	// 		c.JSON(404, gin.H{"error": "Service not found"})
	// 		return
	// 	}

	// 	swaggerURL := ginSwagger.URL(swaggerPath)
	// 	ginSwagger.WrapHandler(swaggerFile.Handler, swaggerURL)(c)
	// })
	// r.Static("/swagger-doc", "../protobuf/docs")
	//
	// ctx := context.Background()
	// muxGateway := runtime.NewServeMux(
	// 	runtime.WithMetadata(_jwt.ExtractMetadataFromRequest),
	// )

	// Đăng ký gRPC service vào HTTP proxy
	// err := pb_noti.RegisterNotificationServiceHandlerFromEndpoint(ctx, muxGateway, "localhost:8004", []grpc.DialOption{grpc.WithInsecure()})
	// if err != nil {
	// 	log.Fatalf("Failed to register gRPC service: %v", err)
	// }
	// err = pb_bdspro.RegisterBdsproServiceHandlerFromEndpoint(ctx, muxGateway, "localhost:8202", []grpc.DialOption{grpc.WithInsecure()})
	// if err != nil {
	// 	log.Fatalf("Failed to register gRPC service: %v", err)
	// }

	// Danh sách các path public
	// publicPaths := []string{
	// 	"/api/v1/public/login",
	// 	"/api/v1/public/register",
	// }
	// Định tuyến các request gRPC đến HTTP proxy
	// r.Any("/v2/*any", gin.WrapH(muxGateway))

	// go func() {
	// 	if err := r.Run(":8080"); err != nil {
	// 		log.Fatalf("Gin server failed: %v", err)
	// 	}
	// }()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Cho phép gọi từ frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Any("/v1/:service/*path", gateway.ProxyRequest)

	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{"message": "pong"})
	// })

	// log.Println("API Gateway started on port " + port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gin server failed: %v", err)
	}
}
