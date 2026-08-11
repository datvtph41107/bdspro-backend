package cmd

import (
	_jwt "common/jwt"
	_middleware "common/middleware"
	"context"
	"fmt"
	"gateway/config"
	_grpcmetadata "gateway/internal/grpcmetadata"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ProfileIDInterceptor xử lý profileID từ metadata
func ProfileIDInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	profileIDList := md.Get("profileid")
	if len(profileIDList) > 0 {
		profileIDStr := profileIDList[0]
		profileID, err := strconv.ParseUint(profileIDStr, 10, 64)
		if err == nil {
			log.Printf("Gán profileID: %d", profileID)
			ctx = context.WithValue(ctx, "profileId", profileID)
		}
	}

	return handler(ctx, req)
}

// GrpcCmd chạy gRPC server
var GrpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Chạy gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadConfig()
		port := fmt.Sprintf(":%s", viper.GetString("server.tcp_port"))
		// url := viper.GetString("service.domain")
		// notificationEndpoint := fmt.Sprintf("%s:%s", url, viper.GetString("service.port.notification"))

		r := gin.Default()
		r.Use(_jwt.JWTAuthMiddleware(nil, nil))

		// gRPC Gateway
		mux := runtime.NewServeMux(
			runtime.WithMetadata(_grpcmetadata.DuringReportMigration(
				_grpcmetadata.ReportTargetEnabled(),
				_middleware.InjectGrpcMetadataContextMiddleware,
			)),
			runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler,
				w http.ResponseWriter, r *http.Request, err error) {
				log.Printf("gRPC Error: %v", err)
				_jwt.ErrorToGatewayResponse(ctx, w, err)
			}),
		)

		// ctx := context.Background()
		// opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		// Đăng ký notification service
		// if err := notificationpb.RegisterNotificationServiceHandlerFromEndpoint(ctx, mux, notificationEndpoint, opts); err != nil {
		// 	log.Printf("Failed to register notification service: %v", err)
		// }

		// Routes
		r.Any("/v2/*any", gin.WrapH(mux))
		r.Static("/swagger-doc", "./cmd/grpc/docs")

		// Swagger UI
		swaggerURL := ginSwagger.URL("/swagger-doc/swagger.json")
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))

		log.Printf("gRPC Gateway server starting on port %s", port)
		if err := r.Run(port); err != nil {
			log.Fatalf("gRPC Gateway server failed: %v", err)
		}
	},
}
