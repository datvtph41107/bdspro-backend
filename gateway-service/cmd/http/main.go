package cmd

import (
	"common/jwtverify"
	_logging "common/logging"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gateway/config"
	_authroutes "gateway/internal/authroutes"
	_grpcgateway "gateway/internal/grpcgateway"
	_grpcmetadata "gateway/internal/grpcmetadata"
	_httpauth "gateway/internal/httpauth"
	_httperror "gateway/internal/httperror"
	_httpmiddleware "gateway/internal/httpmiddleware"
	_httpresponse "gateway/internal/httpresponse"
	_observability "gateway/internal/observability"
	_tqdmultipart "gateway/internal/tqdmultipart"
	_upstream "gateway/internal/upstream"
	_v1proxy "gateway/internal/v1proxy"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// HttpCmd chạy HTTP server
var HttpCmd = &cobra.Command{
	Use:   "http",
	Short: "Chạy HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		runtimeConfig, err := config.LoadRuntime()
		if err != nil {
			return fmt.Errorf("load gateway runtime config: %w", err)
		}
		port := strconv.Itoa(runtimeConfig.Server.Port)
		env := runtimeConfig.Environment

		loggingConfig := _logging.FromEnv("gateway-service")
		loggingConfig.Environment = env
		// Local development keeps console output plus service-scoped JSONL
		// projections. Runtime platforms own production retention and rotation,
		// so non-development environments default to stdout unless explicitly
		// overridden by QHPRO_LOG_OUTPUT.
		if strings.TrimSpace(os.Getenv("QHPRO_LOG_OUTPUT")) == "" &&
			!strings.EqualFold(env, "development") &&
			!strings.EqualFold(env, "dev") &&
			!strings.EqualFold(env, "local") {
			loggingConfig.Output = "stdout"
		}
		logger, closeLogger, err := _logging.New(loggingConfig)
		if err != nil {
			return fmt.Errorf("configure gateway logging: %w", err)
		}
		defer func() { _ = closeLogger() }()
		slog.SetDefault(logger)

		// Gin's own framework output remains stdout. Application/runtime events
		// are emitted by the canonical structured logger above.
		gin.DefaultWriter = os.Stdout

		configureGinMode(
			os.Getenv(
				gin.EnvGinMode,
			),
		)

		httpRecorder :=
			_observability.
				NewHTTPRecorder(
					logger,
				)

		r := gin.New() // custom logging, recovery with req_identity

		// Tránh 301/307 redirect tự thêm/xóa "/" cuối path (httprouter/Gin mặc định bật).
		r.RedirectTrailingSlash = false
		r.RedirectFixedPath = false
		r.HandleMethodNotAllowed =
			true

		r.Use(
			_httpmiddleware.
				RequestID(
					httpRecorder,
				),
		)

		r.Use(
			_httpmiddleware.
				RequestLifecycle(
					httpRecorder,
				),
		)

		r.Use(
			_httpmiddleware.
				Recovery(
					logger,
				),
		)

		// CORS middleware
		corsConfig :=
			_httpmiddleware.
				WithRequestIDCORS(
					cors.Config{
						AllowOrigins: runtimeConfig.HTTP.AllowedOrigins,
						AllowMethods: []string{
							http.MethodGet,
							http.MethodHead,
							http.MethodPost,
							http.MethodPut,
							http.MethodPatch,
							http.MethodDelete,
							http.MethodOptions,
						},
						AllowHeaders: []string{
							"Origin",
							"Content-Type",
							"Accept",
							"Authorization",
							_httpauth.APIKeyHeader,
							_httpauth.LegacyAPIKeyHeader,
							"Lang",
							"If-None-Match",
							"X-QHPro-Telemetry-Key",
							"Referer",
							"timezone_offset",
						},

						ExposeHeaders: []string{
							"Content-Length",
							"Content-Type",
							"Set-Cookie",
							"ETag",
							"Cache-Control",
							"X-Robots-Tag",
							"X-Content-Type-Options",
						},
						AllowCredentials: true,
						MaxAge:           12 * time.Hour,
					},
				)

		corsConfig =
			_httpmiddleware.
				WithOperationIDCORS(
					corsConfig,
				)

		corsConfig =
			_httpmiddleware.
				WithIdempotencyKeyCORS(
					corsConfig,
				)
		r.Use(
			cors.New(
				corsConfig,
			),
		)

		r.Use(
			_httpmiddleware.
				OperationID(),
		)

		r.Use(
			_httpmiddleware.
				IdempotencyKey(),
		)

		// Operational health is process-local; readiness flips before drain.
		readiness := newReadinessState()
		registerOperationalHealth(r, readiness)

		// gRPC Gateway mux. Transport adaptation is owned by dedicated Gateway packages.
		mux := runtime.NewServeMux(
			runtime.WithMetadata(_grpcmetadata.FromHTTPRequest),
			runtime.WithErrorHandler(_httperror.HandleGRPC),
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &_httpresponse.Marshaler{}),
			runtime.WithForwardResponseOption(_httpresponse.ForwardMetadata),
			runtime.WithForwardResponseOption(_httpresponse.ForwardCommon),
		)

		processCtx, stop := signal.NotifyContext(
			context.Background(),
			os.Interrupt,
			syscall.SIGTERM,
		)
		defer stop()

		connections, err := _upstream.Open(runtimeConfig)
		if err != nil {
			return err
		}

		defer func() {
			if closeErr := connections.Close(); closeErr != nil {
				logger.Error("close Gateway upstream connections", "error", closeErr)
			}
		}()

		apiKeyVerifier, verifierErr := _httpauth.NewHubAPIKeyVerifier(connections.Hub)
		if verifierErr != nil {
			logger.Warn("Gateway API key verifier unavailable", "error", verifierErr)
		}

		tokenVerifier, tokenVerifierErr := jwtverify.NewHMACVerifier(runtimeConfig.JWT.VerificationKey)
		if tokenVerifierErr != nil {
			return fmt.Errorf("construct Gateway JWT verifier: %w", tokenVerifierErr)
		}

		authRoutes := _authroutes.AuthenticationRoutes()
		authMiddleware := _httpauth.Middleware(
			authRoutes.Public,
			authRoutes.TemporaryToken,
			_httpauth.WithAPIKeyVerifier(apiKeyVerifier),
			_httpauth.WithTokenVerifier(tokenVerifier),
			_httpauth.WithProviderCredentialRoutes(authRoutes.ProviderCredential),
		)

		// Register generated handlers against process-owned connections.
		if err := _grpcgateway.RegisterNotificationService(processCtx, mux, connections.Notification); err != nil {
			return fmt.Errorf("register Notification gateway handlers: %w", err)
		}
		if err := _grpcgateway.RegisterPaymentService(processCtx, mux, connections.Payment); err != nil {
			return fmt.Errorf("register Payment gateway handlers: %w", err)
		}
		if err := _grpcgateway.RegisterBdsproService(processCtx, mux, connections.BDSPro); err != nil {
			return fmt.Errorf("register Bdspro gateway handlers: %w", err)
		}
		if err := _grpcgateway.RegisterCrmService(processCtx, mux, connections.CRM); err != nil {
			return fmt.Errorf("register Crm gateway handlers: %w", err)
		}
		if err := _grpcgateway.RegisterOrganizationService(processCtx, mux, connections.Organization); err != nil {
			return fmt.Errorf("register Organization gateway handlers: %w", err)
		}
		if err := _grpcgateway.RegisterSocialService(processCtx, mux, connections.Social); err != nil {
			return fmt.Errorf("register Social gateway handlers: %w", err)
		}
		if err := _grpcgateway.RegisterUserService(processCtx, mux, connections.User); err != nil {
			return fmt.Errorf("register User gateway handlers: %w", err)
		}
		if err := _grpcgateway.RegisterChatService(processCtx, mux, connections.Chat); err != nil {
			return fmt.Errorf("register Chat gateway handlers: %w", err)
		}
		if err := _grpcgateway.RegisterAuthService(processCtx, mux, connections.Auth); err != nil {
			return fmt.Errorf("register Auth gateway handlers: %w", err)
		}
		if err := _grpcgateway.RegisterAssistantService(processCtx, mux, connections.Assistant); err != nil {
			return fmt.Errorf("register Assistant gateway handlers: %w", err)
		}
		if err := _grpcgateway.RegisterHubService(processCtx, mux, connections.Hub); err != nil {
			return fmt.Errorf("register Hub gateway handlers: %w", err)
		}
		if err := _grpcgateway.RegisterTqdService(processCtx, mux, connections.TQD); err != nil {
			return fmt.Errorf("register Tqd gateway handlers: %w", err)
		}

		// Multipart compatibility routes borrow the same process-owned TQD connection.
		multipartHandler, err := _tqdmultipart.New(connections.TQD, _tqdmultipart.DefaultLimits())
		if err != nil {
			return fmt.Errorf("construct TQD multipart boundary: %w", err)
		}
		v2Import := r.Group("/form", authMiddleware)
		multipartHandler.Register(v2Import)

		// Gin không cho phép catch-all /v2/*any cùng tồn tại với các route cụ thể
		// như /v2/tqd/public/* và /v2/crm/public/{content,seo}/*. Vì vậy các route
		// cụ thể được đăng ký bình thường; phần /v2 còn lại đi qua NoRoute.
		// v2Fallback := gin.WrapH(mux)
		// r.NoRoute(
		// 	func(c *gin.Context) {
		// 		path := c.Request.URL.Path

		// 		// RC2 gRPC-only: route Discovery legacy bị loại bỏ hoàn toàn và
		// 		// phải trả 404 trước khi JWT middleware xử lý.
		// 		if path == "/v2/discovery" || strings.HasPrefix(path, "/v2/discovery/") {
		// 			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		// 				"code":    "ROUTE_NOT_FOUND",
		// 				"message": "Route not found",
		// 			})
		// 			return
		// 		}

		// 		// NoRoute còn nhận cả URL ngoài /v2. Không chuyển các URL đó
		// 		// sang grpc-gateway vì chúng không thuộc API v2.
		// 		if path != "/v2" && !strings.HasPrefix(path, "/v2/") {
		// 			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		// 				"code":    "ROUTE_NOT_FOUND",
		// 				"message": "Route not found",
		// 			})
		// 			return
		// 		}

		// 		// Gin khởi tạo NoRoute với status 404. Đặt lại status trước khi
		// 		// grpc-gateway ghi response thành công, nếu không body 2xx sẽ
		// 		// bị trả kèm HTTP 404.
		// 		c.Status(http.StatusOK)
		// 		c.Next()
		// 	},
		// 	authMiddleware,
		// 	v2Fallback,
		// )
		r.Any("/v2/*any", authMiddleware, gin.WrapH(mux))
		r.Any("/v2.1/*any", authMiddleware, gin.WrapH(mux))
		r.Any("/v3/*any", authMiddleware, gin.WrapH(mux))

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

		// v1 API routes remain an explicit bounded compatibility proxy.
		v1, err := _v1proxy.New(runtimeConfig, 30*time.Second)
		if err != nil {
			return fmt.Errorf("construct v1 proxy: %w", err)
		}
		// Only true HTTP owners are reverse-proxied. Generated v1 contracts
		// (currently TQD location) stay on the process-owned grpc-gateway mux.
		r.Any("/v1/:service/*path", authMiddleware, v1.Handler(mux))

		server := &http.Server{
			Addr:              ":" + port,
			Handler:           r,
			ReadHeaderTimeout: 10 * time.Second,
		}

		listener, err := net.Listen("tcp", server.Addr)
		if err != nil {
			return fmt.Errorf("listen gateway HTTP on %s: %w", server.Addr, err)
		}

		logger.Info("gateway HTTP server starting", "port", port)

		return serveHTTP(
			processCtx,
			server,
			listener,
			15*time.Second,
			readiness,
		)
	},
}
