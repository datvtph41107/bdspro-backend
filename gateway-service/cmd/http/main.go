package cmd

import (
	_jwt "common/jwt"
	_middleware "common/middleware"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	sharepb "pb/types/shared"
	"time"

	"gateway/config"
	"gateway/pkg"
	"gateway/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/protobuf/proto"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

// publicRoutes danh sách các route không cần authentication
var publicRoutes = []string{
	// BDS Pro public routes
	"/v1/bdspro/v1/user/post",
	"/v1/bdspro/v1/user/post/global",
	"/v1/bdspro/v1/user/list/property-type",
	"/v1/bdspro/v1/user/product/market",
	"/v1/bdspro/v2/user/post",
	"/v1/bdspro/v2/user/post/global",
	"/v1/bdspro/v2/user/list/property-type",
	"/v1/bdspro/v2/user/posts/publish",
	"/v1/bdspro/v2/user/product/market",
	"/v1/bdspro/v2/list/property-type",
	"/v1/bdspro/v2/list/doc-type",
	"/v1/bdspro/v2/list/region",
	"/v1/bdspro/v2/list/project",
	"/v1/bdspro/v2/list/amenity",
	"/v1/bdspro/v2/product/market",

	// User public routes
	"/v1/user/oauth/zalo/login",
	"/v1/user/oauth/facebook/login",
	"/v1/user/oauth/google/login",
	"/v1/user/auth/otp",
	"/v1/user/oauth/zalo/callback",
	"/v1/user/oauth/facebook/callback",
	"/v1/user/oauth/google/callback",
	"/v1/user/auth/token/refresh",
	"/v1/user/profile/info/",
	"/v1/user/profile/find-by-phone",
	"/v2/user/profile/info/",
	"/v2/user/profile/search/public",
	"/v2/user/profile/find-by-phone",

	// Other services public routes
	"/v1/public",
	"/v1/file/swagger",
	"/v1/file/upload",
	"/v1/file/load",
	"/v1/map/swagger",
	"/v1/map/locations/nearby",
	"/v1/map/locations/find-by-polygon",
	"/v2/user/membership/verify",
	"/v2/user/plans",
	"/v1/notification/swagger",
	"/v1/notification/auth/otp",
	"/v1/payment/swagger",
	"/v1/payment/payos/webhook/verify",
	"/v1/payment/wallet/account/new",
	"/v1/payment/bank/all",
	"/v1/search",
	"/v1/crm/contact/relation-ship/",
	"/v2/crm/contact/relation-ship/",
	"/v2/payment/sepay/webhook",

	// Social public routes
	"/v2/social/news-feed/global",
	"/v2/social/report/reasons",
	"/v2/social/comment/list",
	"/v2/social/news-feed/user",
	"/v2/social/news-feed/detail",
	"/v2/social/news-feed/global/reel",
	"/v2/social/news-feed/reel/global",

	"/v2/org/organization/global/profile",

	// Auth public routes
	"/v2/auth/otp",
	"/v2/auth/login/key",
	"/v2/auth/qr",
	"/v2/auth/token/refresh",
	"/v2/auth/oauth",
	"/v2/auth/admin/login",
	"/v2/auth/phone/check",
	"/v2/bdspro/v2/product/market",
	"/v2/bdspro/v2/list/area-region",
	"/v2/bdspro/public/enums/",

	// Feedback public routes
	"/v2/feedback/rate/list",
	"/v2/feedback/rate",
	"/v2/feedback/rate/stats",
	"/v2/feedback/report/reasons",
	"/v2/feedback/report-reason",
	"/v2/feedback/report",
	"/v2/bdspro/v2/product/suggest/deepseek",
	"/v2/bdspro/v2/post/global",
	"/v2/hub/system-config/key",
	"/v2/hub/system-config/by-group",
	"/v2/hub/user-guides",
	"/v2/hub/app-bundle/version",
	"/v2/hub/app-bundle/update-check",

	"/v2/user/profile/search/public",
	"/v2/bdspro/v2/post/detail/",
	"/v2/social/news-feed/detail/",
	"/v2/social/public/news-feed/user/",
	"/v2/user/main-area",
	"/v2/user/purpose-use",
	"/v2/auth/device/check",
	"/v3/user/profile/p",
	"/v3/bdspro/post/personal",
	"/v2/assistant/national-card/extract",
	"/v2/pub/user/profile/tags",
}

var tempRoutes = []string{
	"/v2/auth/fullname",
	"/v2/auth/phone",
}

// HttpCmd chạy HTTP server
var HttpCmd = &cobra.Command{
	Use:   "http",
	Short: "Chạy HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadConfig()
		port := viper.GetString("server.port")
		gateway := pkg.NewGateway("url")

		// Set timezone
		os.Setenv("TZ", "Europe/London")
		time.Local = time.FixedZone("Europe/London", 7*60*60)

		// Tạo thư mục nếu chưa có
		os.MkdirAll("./app/logs", os.ModePerm)

		// Mở file log
		f, err := os.OpenFile("./app/logs/app.log",
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

		r := gin.Default()

		// CORS middleware
		r.Use(cors.New(cors.Config{
			AllowOrigins: []string{
				"http://localhost:3000",
				"http://localhost:3001",
				"http://localhost:8000",
				"http://localhost:5173",
				"http://localhost:5174",
				"http://14.225.210.29:8000",
				"http://14.225.205.222:8000",
				"http://14.225.205.222:3107",
				"http://14.225.205.222:3108",
				"http://103.145.63.185:3107",
				"http://localhost:5173",
				"https://admin.bdspro.com",
				"https://bdspro.com",
				"https://app.bdspro.com",
			},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Lang"},
			ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))

		env := os.Getenv("ENV_RUNTIME")

		// gRPC Gateway muxcs
		mux := runtime.NewServeMux(
			runtime.WithMetadata(_middleware.InjectGrpcMetadataContextMiddleware),
			runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
				if tokenErr, ok := ctx.Value("token_error").(error); ok {
					http.Error(w, tokenErr.Error(), http.StatusUnauthorized)
					return
				}
				_jwt.ErrorToGatewayResponse(ctx, w, err)
			}),
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &CustomMarshaler{}),
			runtime.WithForwardResponseOption(func(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
				// nếu resp implement interface nào đó, ví dụ CommonResponse
				if cr, ok := resp.(*sharepb.CommonResponse); ok {
					// 1. apply cookies
					for _, ck := range cr.Cookies {
						http.SetCookie(w, &http.Cookie{
							Name:     ck.Name,
							Value:    ck.Value,
							Path:     ck.Path,
							Domain:   ck.Domain,
							HttpOnly: ck.HttpOnly,
							Secure:   ck.Secure,
							MaxAge:   int(ck.MaxAge),
						})
					}

					// 2. handle redirect
					if cr.RedirectUrl != "" {
						http.Redirect(w, nil, cr.RedirectUrl, http.StatusFound) // 302
						return nil                                              // stop normal response
					}
				}
				return nil
			}),
		)

		ctx := context.Background()
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithConnectParams(grpc.ConnectParams{
				Backoff: backoff.Config{
					BaseDelay:  5 * time.Second,
					Multiplier: 1.0,
					Jitter:     0.0,
					MaxDelay:   5 * time.Second,
				},
				MinConnectTimeout: 5 * time.Second,
			}),
		}

		// Đăng ký các gRPC services
		service.RegisterNotificationService(ctx, mux, env, opts)
		service.RegisterPaymentService(ctx, mux, env, opts)
		service.RegisterOrganizationService(ctx, mux, env, opts)
		service.RegisterBdsproService(ctx, mux, env, opts)
		service.RegisterCrmService(ctx, mux, env, opts)
		service.RegisterSocialService(ctx, mux, env, opts)
		service.RegisterAppointmentService(ctx, mux, env, opts)
		service.RegisterUserService(ctx, mux, env, opts)
		service.RegisterChatService(ctx, mux, env, opts)
		service.RegisterMarketingService(ctx, mux, env, opts)
		service.RegisterTransactionService(ctx, mux, env, opts)
		service.RegisterAuthService(ctx, mux, env, opts)
		service.RegisterAssistantService(ctx, mux, env, opts)
		service.RegisterHubService(ctx, mux, env, opts)
		// service.RegisterTqdService(ctx, mux, env, opts)

		// v2 API routes (gRPC Gateway)
		r.Any("/v2/*any", _jwt.JWTAuthMiddleware(publicRoutes, tempRoutes), gin.WrapH(mux))
		r.Any("/v2.1/*any", _jwt.JWTAuthMiddleware(publicRoutes, tempRoutes), gin.WrapH(mux))
		r.Any("/v3/*any", _jwt.JWTAuthMiddleware(publicRoutes, tempRoutes), gin.WrapH(mux))

		// Swagger documentation
		r.Static("/swagger-doc", "docs")

		// Merged swagger endpoint
		r.GET("/swagger/merged/*any", func(c *gin.Context) {
			swaggerURL := ginSwagger.URL("/swagger-doc/merged_swagger.json")
			handler := ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL)
			handler(c)
		})

		// Individual service swagger endpoints
		r.GET("/swagger/:service/*any", func(c *gin.Context) {
			service := c.Param("service")
			swaggerURL := ginSwagger.URL(fmt.Sprintf("/swagger-doc/%s/swagger.json", service))
			handler := ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL)
			handler(c)
		})

		// v1 API routes (legacy proxy)
		forwarding := pkg.ForwardingRequest(gateway, env)
		r.Any("/v1/:service/*path", _jwt.JWTAuthMiddleware(publicRoutes, tempRoutes), forwarding)

		log.Printf("Gateway server starting on port: %s", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("Gateway server failed: %v", err)
		}
	},
}
