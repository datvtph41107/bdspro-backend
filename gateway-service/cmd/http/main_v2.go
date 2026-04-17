package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gateway/internal/config"
	"gateway/pkg"
	"gateway/service"

	_jwt "common/jwt"
	_middleware "common/middleware"
	sharepb "pb/types/shared"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	sharedcfg "common/pkg/config"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var HttpV2Cmd = &cobra.Command{
	Use:   "http",
	Short: "Chạy HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Load config (fail-fast)
		cfg := config.MustLoad()
		configHash := sharedcfg.ComputeHash(cfg)
		log.Printf("Gateway config loaded:\n%s", sharedcfg.MaskSensitive(cfg))
		log.Printf("Config hash: %s", configHash)

		// 2. Timezone
		os.Setenv("TZ", "Europe/London")
		time.Local = time.FixedZone("Europe/London", 7*60*60)

		// 3. Log file
		os.MkdirAll("./app/logs", os.ModePerm)
		f, err := os.OpenFile("./app/logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

		// 4. Gin engine
		r := gin.Default()
		r.RedirectTrailingSlash = false
		r.RedirectFixedPath = false

		// CORS từ config
		r.Use(cors.New(cors.Config{
			AllowOrigins:     cfg.Gateway.AllowedOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Lang"},
			ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))

		// 5. gRPC Gateway mux (giữ nguyên custom marshaler, error handler)
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
				if cr, ok := resp.(*sharepb.CommonResponse); ok {
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
					if cr.RedirectUrl != "" {
						http.Redirect(w, nil, cr.RedirectUrl, http.StatusFound)
						return nil
					}
				}
				return nil
			}),
		)

		ctx := context.Background()
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}

		// 6. Đăng ký gRPC services dùng config
		registerServices(ctx, mux, cfg, opts)

		// 7. Routes
		r.Any("/v2/*any", _jwt.JWTAuthMiddleware(cfg.Gateway.PublicRoutes, cfg.Gateway.TempRoutes), gin.WrapH(mux))
		r.Any("/v2.1/*any", _jwt.JWTAuthMiddleware(cfg.Gateway.PublicRoutes, cfg.Gateway.TempRoutes), gin.WrapH(mux))
		r.Any("/v3/*any", _jwt.JWTAuthMiddleware(cfg.Gateway.PublicRoutes, cfg.Gateway.TempRoutes), gin.WrapH(mux))

		// 8. Swagger (giữ nguyên)
		r.Static("/swagger-doc", "docs")
		r.GET("/swagger/merged/*any", func(c *gin.Context) {
			swaggerURL := ginSwagger.URL("/swagger-doc/merged_swagger.json")
			handler := ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL)
			handler(c)
		})
		r.GET("/swagger/:service/*any", func(c *gin.Context) {
			svc := c.Param("service")
			swaggerURL := ginSwagger.URL(fmt.Sprintf("/swagger-doc/%s/swagger.json", svc))
			handler := ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL)
			handler(c)
		})

		// 9. Legacy proxy v1 dùng config
		proxyHandler := pkg.NewProxyHandler(cfg)
		r.Any("/v1/:service/*path", _jwt.JWTAuthMiddleware(cfg.Gateway.PublicRoutes, cfg.Gateway.TempRoutes), proxyHandler)

		// 10. Debug endpoint (bảo vệ)
		debugToken := os.Getenv("DEBUG_TOKEN")
		if debugToken == "" {
			debugToken = "internal"
		}
		r.GET("/debug/config", func(c *gin.Context) {
			if c.GetHeader("X-Debug-Token") != debugToken {
				c.Status(http.StatusUnauthorized)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"hash":   configHash,
				"config": sharedcfg.MaskSensitive(cfg),
			})
		})

		// 11. Health check
		r.GET("/health/live", func(c *gin.Context) { c.Status(http.StatusOK) })
		r.GET("/health/ready", func(c *gin.Context) { c.Status(http.StatusOK) })

		// 12. Start server với graceful shutdown
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		srv := &http.Server{
			Addr:         addr,
			Handler:      r,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
		}

		go func() {
			log.Printf("Gateway listening on %s", addr)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server failed: %v", err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down gateway...")
		ctxShutdown, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(ctxShutdown); err != nil {
			log.Fatalf("Forced shutdown: %v", err)
		}
		log.Println("Gateway stopped")
	},
}

// registerServices kết nối đến các service và đăng ký handler
func registerServices(ctx context.Context, mux *runtime.ServeMux, cfg *config.GatewayConfig, opts []grpc.DialOption) {
	dialWithRetry := func(ep sharedcfg.ServiceEndpoint) (*grpc.ClientConn, error) {
		var lastErr error
		for i := 0; i <= ep.Retries; i++ {
			for _, addr := range ep.Addresses {
				dialCtx, cancel := context.WithTimeout(ctx, ep.Timeout)
				conn, err := grpc.DialContext(dialCtx, addr, opts...)
				cancel()
				if err == nil {
					return conn, nil
				}
				lastErr = err
			}
			if i < ep.Retries {
				time.Sleep(ep.Backoff)
			}
		}
		return nil, fmt.Errorf("failed after %d retries: %w", ep.Retries, lastErr)
	}

	for name, ep := range cfg.Services.Endpoints {
		conn, err := dialWithRetry(ep)
		if err != nil {
			log.Printf("Warning: cannot connect to %s: %v", name, err)
			continue
		}
		switch name {
		case "auth":
			service.RegisterAuthV2Service(ctx, mux, conn)
		// case "user":
		// 	service.RegisterUserService(ctx, mux, conn)
		// case "bdspro":
		// 	service.RegisterBdsproService(ctx, mux, conn)
		// case "crm":
		// 	service.RegisterCrmService(ctx, mux, conn)
		// case "chat":
		// 	service.RegisterChatService(ctx, mux, conn)
		// case "payment":
		// 	service.RegisterPaymentService(ctx, mux, conn)
		// case "notification":
		// 	service.RegisterNotificationService(ctx, mux, conn)
		// case "organization":
		// 	service.RegisterOrganizationService(ctx, mux, conn)
		// case "social":
		// 	service.RegisterSocialService(ctx, mux, conn)
		// case "appointment":
		// 	service.RegisterAppointmentService(ctx, mux, conn)
		// case "marketing":
		// 	service.RegisterMarketingService(ctx, mux, conn)
		// case "transaction":
		// 	service.RegisterTransactionService(ctx, mux, conn)
		// case "assistant":
		// 	service.RegisterAssistantService(ctx, mux, conn)
		// case "hub":
		// 	service.RegisterHubService(ctx, mux, conn)
		// case "tqd":
		// 	service.RegisterTqdService(ctx, mux, conn)
		// ... các service còn lại
		default:
			log.Printf("Unknown service %s, skip registration", name)
		}
	}
}
