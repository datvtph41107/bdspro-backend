package main

import (
	"common/logging"
	logfabric "common/logging/fabric"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	_db "common/db"
	_jwt "common/jwt"
	"common/requestlog"
	"mdn4/config"
	userdb "mdn4/db"
	"mdn4/internal/requeststats"
	chatmodule "mdn4/modules/chat"
	"mdn4/wire"
)

func main() {
	defer logging.MustStart("mdn4-http")()
	logfabric.Bridge()
	if err := runHTTP(); err != nil {
		log.Fatal(err)
	}
}

func chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				if rec == http.ErrAbortHandler {
					panic(rec)
				}
				ctx := r.Context()
				requestlog.FromContext(ctx, nil).ErrorContext(ctx, "mdn4 HTTP panic",
					"category", "application",
					"error", fmt.Sprint(rec),
					"http_method", r.Method,
					"http_path", r.URL.Path,
					"stack_trace", string(debug.Stack()),
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// logUnmatchedRoutes warns when no registered pattern matches, since ServeMux
// answers 404/405 itself without reaching any handler.
func logUnmatchedRoutes(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, pattern := mux.Handler(r); pattern == "" {
			ctx := r.Context()
			requestlog.FromContext(ctx, nil).WarnContext(ctx, "HTTP route not registered",
				"category", "application",
				"http_method", r.Method,
				"http_path", r.URL.Path,
			)
		}
		mux.ServeHTTP(w, r)
	})
}

// runHTTP starts a real HTTP server for the v4 modular surface.
// Login and permission/me routes are owned by modules/auth; cmd/grpc remains the legacy gRPC entry.
func runHTTP() error {
	runtimeConfig, err := config.LoadRuntime()
	if err != nil {
		return fmt.Errorf("load mdn4 runtime config: %w", err)
	}
	if runtimeConfig.HTTPAddress == "" {
		return fmt.Errorf("http address is required (set HTTP_ADDRESS or server.port)")
	}
	log.Printf(
		"mdn4 database target: %s",
		config.DescribeDatabaseDSN(runtimeConfig.Database.DSN),
	)

	runtimePolicy, err := loadRuntimePolicy()
	if err != nil {
		return fmt.Errorf("load mdn4 runtime policy: %w", err)
	}
	log.Printf(
		"mdn4 http runtime policy: module_warmup=%t dashboard_stats=%t zns_scheduler=%t",
		runtimePolicy.ModuleWarmup,
		runtimePolicy.DashboardStats,
		runtimePolicy.ZNSScheduler,
	)

	app, cleanup, err := wire.InitializeApp()
	if err != nil {
		return fmt.Errorf("initialize mdn4 app: %w", err)
	}
	defer cleanup()
	if app == nil || app.DB == nil {
		return fmt.Errorf("initialize mdn4 app: app or database is nil")
	}
	if app.AuthHandler == nil || app.AuthHandler.LoginController == nil {
		return fmt.Errorf("initialize mdn4 app: login controller is nil")
	}
	if app.PermissionHandler == nil || app.PermissionHandler.PermissionController == nil {
		return fmt.Errorf("initialize mdn4 app: permission controller is nil")
	}
	if app.BdsproHandler == nil {
		return fmt.Errorf("initialize mdn4 app: bdspro handler is nil")
	}
	if app.HubHandler == nil {
		return fmt.Errorf("initialize mdn4 app: hub handler is nil")
	}
	if app.NotificationHandler == nil {
		return fmt.Errorf("initialize mdn4 app: notification handler is nil")
	}
	if app.TqdHandler == nil {
		return fmt.Errorf("initialize mdn4 app: tqd handler is nil")
	}
	if app.UserHandler == nil {
		return fmt.Errorf("initialize mdn4 app: user handler is nil")
	}
	if app.SeoHandler == nil {
		return fmt.Errorf("initialize mdn4 app: seo handler is nil")
	}
	if app.PaymentHandler == nil {
		return fmt.Errorf("initialize mdn4 app: payment handler is nil")
	}
	if app.AssistantHandler == nil {
		return fmt.Errorf("initialize mdn4 app: assistant handler is nil")
	}

	_ = userdb.RunPendingMigrations(app.DB)
	if err := _db.EnsureSchema(app.DB, "user", userdb.AutoMigrate); err != nil {
		return err
	}

	signalCtx, stopSignals := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stopSignals()

	if runtimePolicy.ModuleWarmup {
		if err := app.LoadRoleGroupsOnStartup(signalCtx); err != nil {
			log.Printf("mdn4 module cache warmup failed; starting degraded: %v", err)
		}
	}

	mux := http.NewServeMux()
	chatHandler, chatCleanup, err := chatmodule.NewHTTPHandler(app.DB, app.ProfileUsecase)
	if err != nil {
		return fmt.Errorf("initialize chat module: %w", err)
	}
	defer chatCleanup()
	mux.Handle("/v2/chat/", chatHandler)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	// Controllers register via httpapi.GET/Bind: path+query decode into DTO,
	// handlers are func(ctx, req DTO) (*httpapi.Result, error).
	app.AuthHandler.LoginController.RegisterHTTP(mux)
	app.OAuthHandler.RegisterHTTP(mux)
	app.PermissionHandler.PermissionController.RegisterHTTP(mux)
	app.UserHandler.RegisterHTTP(mux)
	app.BdsproHandler.RegisterHTTP(mux)
	app.HubHandler.RegisterHTTP(mux)
	app.NotificationHandler.RegisterHTTP(mux)
	app.TqdHandler.RegisterHTTP(mux)
	app.CrmHandler.RegisterHTTP(mux)
	app.SeoHandler.RegisterHTTP(mux)
	app.SeoHandler.RegisterInternalHTTP(mux, runtimeConfig.Seo.InternalSecretKey)
	app.PaymentHandler.RegisterHTTP(mux)
	app.AssistantHandler.RegisterHTTP(mux)

	handler := chain(
		logUnmatchedRoutes(mux),
		recoveryMiddleware,
		logging.AccessLogMiddleware,
		corsMiddleware(runtimeConfig.HTTP.AllowedOrigins),
		_jwt.JWTAuthHTTPMiddleware(publicHTTPRoutes, nil),
		bridgeTrustedHTTPIdentity(),
		requeststats.Middleware(requeststats.New(
			app.HubHandler.RequestStatsCollector(),
			runtimeConfig.RequestStats.SpecURLs,
			runtimeConfig.RequestStats.IgnorePrefixes,
		)),
	)

	server := &http.Server{
		Addr:              runtimeConfig.HTTPAddress,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return fmt.Errorf("listen mdn4 HTTP on %s: %w", server.Addr, err)
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("mdn4 HTTP listening on %s", runtimeConfig.HTTPAddress)
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-signalCtx.Done():
		log.Printf("mdn4 HTTP shutting down...")
	case err := <-errCh:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return server.Shutdown(shutdownCtx)
}
