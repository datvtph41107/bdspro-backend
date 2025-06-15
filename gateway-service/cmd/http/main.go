package cmd

import (
	_jwt "common/jwt"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gateway/config"
	"gateway/pkg"

	"gateway/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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

// Custom metadata injector
func CustomMetadataMethod(ctx context.Context, req *http.Request) metadata.MD {
	md := metadata.Pairs(
		"x-http-method", req.Method,
	)
	log.Println("Method: ", req.Method)
	return md
}

var publicRoutes = []string{
	"/v1/bdspro/v1/user/post",
	"/v1/bdspro/v1/user/post/global",
	"/v1/bdspro/v1/user/list/property-type",
	"/v1/bdspro/v1/user/post/global",
	"/v1/bdspro/v1/user/product/market",
	"/v1/bdspro/v2/user/post",
	"/v1/bdspro/v2/user/post/global",
	"/v1/bdspro/v2/user/list/property-type",
	"/v1/bdspro/v2/user/post/global",
	"/v1/bdspro/v2/user/product/market",
	"/v1/bdspro/v2/user/posts/publish",
	"/v1/user/rate/list",
	"/v1/public",
	"/v1/file/swagger",
	"/v1/file/upload",
	"/v1/file/load",
	"/v1/map/swagger",
	"/v1/map/locations/nearby",
	"/v1/map/locations/find-by-polygon",
	"/v1/membership/user/verify",
	"/v1/notification/swagger",
	"/v1/notification/auth/otp",
	"/v1/payment/swagger",
	"/v1/payment/payos/webhook/verify",
	"/v1/payment/wallet/account/new",
	"/v1/user/swagger",
	"/v1/user/auth/otp",
	"/v1/user/oauth/zalo/callback",
	"/v1/user/oauth/facebook/callback",
	"/v1/user/oauth/google/callback",
	"/v1/user/auth/token/refresh",
	"/v1/user/profile/info/",
	"/v1/user/profile/info/",
	"/v1/user/profile/find-by-phone",
	"/v1/membership/plan/list",
	"/v1/membership/user/plans",
	"/v1/search",
	"/v2/social/news-feed/public",
	"/v2/social/report/reasons",
	"/v2/social/comment/list",
	"/v2/social/news-feed/user",
	"/v2/social/news-feed/detail",
	"/v2/bdspro/v2/product/market",
	"/v2/bdspro/v2/list/property-type",
	"/v2/bdspro/v2/list/doc-type",
	"/v2/bdspro/v2/list/region",
	"/v2/bdspro/v2/list/project",
	"/v2/bdspro/v2/list/amenity",
	// "/v1/search/plan/list",
}

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 				"error": "Authorization header is required",
// 			})
// 			return
// 		}
// 		c.Next()
// 	}
// }

// var public

// httpCmd represents the http command
// func ExtractMetadataFromRequest(ctx context.Context, r *http.Request) metadata.MD {
// 	authHeader := r.Header.Get("Authorization")
// 	if authHeader == "" {
// 		// Gán lỗi vào context để ErrorHandler xử lý sau
// 		ctx = context.WithValue(ctx, "token_error", errors.New("Missing Authorization token"))
// 		return metadata.MD{}
// 	}

// 	return metadata.Pairs("authorization", authHeader)
// }

var HttpCmd = &cobra.Command{
	Use:   "http",
	Short: "Run the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadConfig()
		port := viper.GetString("server.port")
		gateway := pkg.NewGateway("url")
		// gateway.FetchServices()

		r := gin.Default()
		env := os.Getenv("APP_ENV")

		mux := runtime.NewServeMux(
			runtime.WithMetadata(_jwt.ExtractMetadataFromRequest),
			// // runtime.WithMetadata(CustomMetadataMethod),
			// runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
			// 	_jwt.ErrorToGatewayResponse(ctx, w, err)
			// }),
			// runtime.WithMetadata(ExtractMetadataFromRequest),
			runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
				if tokenErr, ok := ctx.Value("token_error").(error); ok {
					http.Error(w, tokenErr.Error(), http.StatusUnauthorized)
					return
				}
				_jwt.ErrorToGatewayResponse(ctx, w, err)
			}),
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &CustomMarshaler{}),
		)
		ctx := context.Background()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		service.RegisterNotificationService(ctx, mux, env, opts)
		service.RegisterPaymentService(ctx, mux, env, opts)
		service.RegisterOrganizationService(ctx, mux, env, opts)
		service.RegisterBdsproService(ctx, mux, env, opts)
		service.RegisterCrmService(ctx, mux, env, opts)
		service.RegisterSocialService(ctx, mux, env, opts)
		service.RegisterAppointmentService(ctx, mux, env, opts)

		r.Any("/v2/*any", _jwt.JWTAuthMiddleware(publicRoutes), gin.WrapH(mux))


		r.Static("/swagger-doc", "docs")
		r.GET("/swagger/:service/*any", func(c *gin.Context) {
			service := c.Param("service")

			// Ví dụ: /swagger/gateway/index.html => trỏ tới /swagger-doc/gateway/swagger.json
			swaggerURL := ginSwagger.URL(fmt.Sprintf("/swagger-doc/%s/swagger.json", service))

			// Tạo handler mới với URL tương ứng
			handler := ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL)

			handler(c)
		})

		r.Use(_jwt.JWTAuthMiddleware(publicRoutes))
		r.Use(cors.New(cors.Config{ // Cho phép gọi từ frontend
			AllowOrigins: []string{
				"http://localhost:3000",
				"http://localhost:8000",
				"http://14.225.210.29:8000",
				"http://14.225.205.222:8000",
			},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))

		Forwarding := pkg.ForwardingRequest(gateway, env)
		r.Any("/v1/:service/*path", Forwarding)

		log.Println("Start gateway server at port: ", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("Gin server failed: %v", err)
		}
	},
}
