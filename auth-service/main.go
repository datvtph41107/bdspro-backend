package main

import (
	"auth/config"
	"auth/infra/handler"
	authrpc "auth/infra/rpc"
	"auth/infra/services"
	_middleware "common/middleware"
	_ "common/models"
	"common/processlifecycle"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	authpb "pb/types/auth"

	"github.com/spf13/viper"
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
	if err := run(); err != nil {
		log.Printf("auth service failed: %v", err)
		os.Exit(1)
	}
}

func run() error {
	if err := config.LoadProperties(); err != nil {
		return fmt.Errorf("load auth config: %w", err)
	}

	userRPC, userCleanup, err := authrpc.NewUserRPCClient()
	if err != nil {
		return fmt.Errorf("configure auth user RPC: %w", err)
	}
	defer userCleanup()

	permissionService := services.NewPermissionService(userRPC)
	authInternalHandler := handler.NewAuthInternalHandler(permissionService)
	if err := os.Setenv("TZ", "Europe/London"); err != nil {
		return fmt.Errorf("set auth timezone: %w", err)
	}
	time.Local = time.FixedZone("Europe/London", 7*60*60)

	port := fmt.Sprintf(":%s", viper.GetString("server.tcp_port"))
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("listen auth gRPC on %s: %w", port, err)
	}
	defer func() { _ = lis.Close() }()

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			_middleware.ParseGrpcMetadataContextMiddleware,
		),
		grpc.ChainStreamInterceptor(
			_middleware.ParseGrpcMetadataContextStreamMiddleware,
		),
	)
	authpb.RegisterAuthInternalServiceServer(s, authInternalHandler)

	log.Printf("Listen: %v", port)
	signalCtx, stopSignals := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSignals()
	processCtx, processCancel := context.WithCancel(signalCtx)
	defer processCancel()

	return processlifecycle.Run(
		processCtx,
		processCancel,
		s,
		lis,
		[]processlifecycle.Actor{processlifecycle.ActorFunc(permissionService.Run)},
		processlifecycle.Config{GracefulStopTimeout: 10 * time.Second},
	)
}
