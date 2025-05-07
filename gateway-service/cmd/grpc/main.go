package cmd

import (
	_jwt "common/jwt"
	"context"
	"fmt"
	"gateway/config"
	"log"
	"net/http"
	pb_noti "pb/types/notification"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func ProfileIDInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req) // Không có metadata thì tiếp tục
	}

	profileIDList := md.Get("profileid")
	if len(profileIDList) > 0 {
		profileIDStr := profileIDList[0]

		profileID, err := strconv.ParseUint(profileIDStr, 10, 64)
		if err == nil {
			log.Printf("Gán profileID", profileID)
			// Gán profileID vào context mới
			ctx = context.WithValue(ctx, "profileId", profileID)
		}
	}

	// Gọi handler tiếp theo (controller thực sự)
	return handler(ctx, req)
}

// @deprecated
var GrpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run the GRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadConfig()
		port := fmt.Sprintf(":%s", viper.GetString("server.tcp_port"))

		url := viper.GetString("service.domain")
		notification_endpoint := fmt.Sprintf("%s:%s", url, viper.GetString("service.port.notification"))

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
		// 	ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL)(c)
		// })

		r := gin.Default()

		// 1. gRPC Gateway xử lý mọi route
		mux := runtime.NewServeMux(
			runtime.WithMetadata(_jwt.ExtractMetadataFromRequest),
			runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
				log.Printf("Err", err)
				_jwt.ErrorToGatewayResponse(ctx, w, err)
			}),
		)
		ctx := context.Background()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		_ = pb_noti.RegisterNotificationServiceHandlerFromEndpoint(ctx, mux, notification_endpoint, opts)

		// gRPC Gateway qua Gin
		r.Any("/v2/*any", gin.WrapH(mux))
		r.Static("/swagger-doc", "./cmd/grpc/docs")
		// 3. Swagger UI dùng ginSwagger
		swaggerURL := ginSwagger.URL("/swagger-doc/swagger.json")
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))

		r.Run(port)
	},
}
