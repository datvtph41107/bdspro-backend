package main

import (
	"auth/config"
	"auth/infra/handler"
	authrpc "auth/infra/rpc"
	"auth/infra/services"
	"common/logging"
	_middleware "common/middleware"
	_ "common/models"
	process "common/process"
	qhprorpc "common/rpc"
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	authpb "pb/types/auth"

	"google.golang.org/grpc"
)

// @title BDSPRO API
// @version 1.0
// @description API phần logic chính của BĐSPro
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// auth không giữ logic liên quan tới lưu trữ db
// tức auth như 1 service đứng giữa gateway và service con
// gateway sau chỉ làm 1 việc duy nhất là điều hướng
// tương lai middleware jwt của gateway sẽ bê vào auth và toàn bộ request chia 2 hướng
// nếu public -> đi luôn vào service con
// nếu private -> đi vào auth -> auth điều hướng tới service con
func main() {
	os.Exit(runProcess(run))
}

func runProcess(runFn func() error) int {
	closeLogger, err := logging.Configure("auth-service")
	if err != nil {
		fmt.Fprintf(os.Stderr, "configure auth logging: %v\n", err)
		return 1
	}
	defer func() { _ = closeLogger() }()

	if err := runFn(); err != nil {
		slog.Error("auth service failed", slog.Any("error", err))
		return 1
	}

	return 0
}

func run() error {
	if err := config.LoadProperties(); err != nil {
		return fmt.Errorf("load auth config: %w", err)
	}
	runtime, err := config.NewRuntime()
	if err != nil {
		return fmt.Errorf("materialize auth runtime: %w", err)
	}

	userRPC, userCleanup, err := authrpc.NewUserRPCClient(
		runtime.UserRPCTarget,
		runtime.Transport,
	)
	if err != nil {
		return fmt.Errorf("configure auth user RPC: %w", err)
	}
	defer userCleanup()

	permissionService := services.NewPermissionService(
		userRPC,
		services.PermissionConfig{
			RefreshInterval: runtime.Permission.RefreshInterval,
			MaxStaleness:    runtime.Permission.MaxStaleness,
			RoleCacheTTL:    runtime.Permission.RoleCacheTTL,
			RequestTimeout:  runtime.Permission.RequestTimeout,
		},
	)
	authInternalHandler := handler.NewAuthInternalHandler(permissionService, userRPC)
	if err := os.Setenv("TZ", "Europe/London"); err != nil {
		return fmt.Errorf("set auth timezone: %w", err)
	}
	time.Local = time.FixedZone("Europe/London", 7*60*60)

	port := fmt.Sprintf(":%d", runtime.GRPCPort)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("listen auth gRPC on %s: %w", port, err)
	}
	defer func() { _ = lis.Close() }()

	// Auth Service chỉ phục vụ AuthInternal. Vì vậy mọi RPC đi vào đây phải
	// mang service assertion đã ký; Auth không tự tin tưởng metadata từ client.
	transport := qhprorpc.ServerTransport(runtime.Transport)
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			transport.Unary,
			_middleware.ParseGrpcMetadataContextMiddleware,
			_middleware.UnaryErrorInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			_middleware.ParseGrpcMetadataContextStreamMiddleware,
			_middleware.StreamErrorInterceptor(),
		),
	)
	authpb.RegisterAuthInternalServiceServer(s, authInternalHandler)

	slog.Info("auth gRPC listening", slog.String("address", port))
	signalCtx, stopSignals := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSignals()
	processCtx, processCancel := context.WithCancel(signalCtx)
	defer processCancel()

	return process.Run(
		processCtx,
		processCancel,
		s,
		lis,
		[]process.Actor{process.ActorFunc(permissionService.Run)},
		process.Config{GracefulStopTimeout: 10 * time.Second},
	)
}
