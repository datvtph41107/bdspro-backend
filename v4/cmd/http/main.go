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
	"syscall"
	"time"

	_db "common/db"
	_jwt "common/jwt"
	"userv4/config"
	userdb "userv4/db"
	"userv4/internal/requeststats"
	"userv4/wire"
)

func main() {
	defer logging.MustStart("userv4-http")()
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
				log.Printf("userv4 HTTP panic: %v", rec)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// runHTTP starts a real HTTP server for the v4 modular surface.
// Login and permission/me routes are owned by modules/auth; cmd/grpc remains the legacy gRPC entry.
func runHTTP() error {
	runtimeConfig, err := config.LoadRuntime()
	if err != nil {
		return fmt.Errorf("load userv4 runtime config: %w", err)
	}
	if runtimeConfig.HTTPAddress == "" {
		return fmt.Errorf("http address is required (set HTTP_ADDRESS or server.port)")
	}
	log.Printf(
		"userv4 database target: %s",
		config.DescribeDatabaseDSN(runtimeConfig.Database.DSN),
	)

	runtimePolicy, err := loadRuntimePolicy()
	if err != nil {
		return fmt.Errorf("load userv4 runtime policy: %w", err)
	}
	log.Printf(
		"userv4 http runtime policy: module_warmup=%t dashboard_stats=%t zns_scheduler=%t",
		runtimePolicy.ModuleWarmup,
		runtimePolicy.DashboardStats,
		runtimePolicy.ZNSScheduler,
	)

	app, cleanup, err := wire.InitializeApp()
	if err != nil {
		return fmt.Errorf("initialize userv4 app: %w", err)
	}
	defer cleanup()
	if app == nil || app.DB == nil {
		return fmt.Errorf("initialize userv4 app: app or database is nil")
	}
	if app.AuthHandler == nil || app.AuthHandler.LoginController == nil {
		return fmt.Errorf("initialize userv4 app: login controller is nil")
	}
	if app.PermissionHandler == nil || app.PermissionHandler.PermissionController == nil {
		return fmt.Errorf("initialize userv4 app: permission controller is nil")
	}
	if app.BdsproHandler == nil {
		return fmt.Errorf("initialize userv4 app: bdspro handler is nil")
	}
	if app.HubHandler == nil {
		return fmt.Errorf("initialize userv4 app: hub handler is nil")
	}
	if app.NotificationHandler == nil {
		return fmt.Errorf("initialize userv4 app: notification handler is nil")
	}
	if app.TqdHandler == nil {
		return fmt.Errorf("initialize userv4 app: tqd handler is nil")
	}
	if app.UserHandler == nil {
		return fmt.Errorf("initialize userv4 app: user handler is nil")
	}

	if err := userdb.RunPendingMigrations(app.DB); err != nil {
		return fmt.Errorf("userv4 schema migrate: %w", err)
	}
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
			log.Printf("userv4 module cache warmup failed; starting degraded: %v", err)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	app.AuthHandler.LoginController.RegisterHTTP(mux)
	app.PermissionHandler.PermissionController.RegisterHTTP(mux)
	app.UserHandler.RegisterHTTP(mux)
	app.BdsproHandler.RegisterHTTP(mux)
	app.HubHandler.RegisterHTTP(mux)
	app.NotificationHandler.RegisterHTTP(mux)
	app.TqdHandler.RegisterHTTP(mux)

	handler := chain(
		mux,
		recoveryMiddleware,
		loggingMiddleware,
		corsMiddleware(runtimeConfig.HTTP.AllowedOrigins),
		_jwt.JWTAuthHTTPMiddleware(publicHTTPRoutes, nil),
		bridgeActorProfileID(),
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
		return fmt.Errorf("listen userv4 HTTP on %s: %w", server.Addr, err)
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("userv4 HTTP listening on %s", runtimeConfig.HTTPAddress)
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-signalCtx.Done():
		log.Printf("userv4 HTTP shutting down...")
	case err := <-errCh:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return server.Shutdown(shutdownCtx)
}
