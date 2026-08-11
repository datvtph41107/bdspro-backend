package cmd

import (
	"common"
	"common/grpcclient"
	_redis "common/redis"
	"common/transporttrust"
	"context"
	"fmt"
	"log/slog"
	"os"
	filepb "pb/types/file"
	userpb "pb/types/user"
	"strings"
	"tqd/env"
	"tqd/infra/filegrpc"
	handler_grpc "tqd/infra/handler/grpc"
	"tqd/infra/reportgenerator"
	"tqd/internal/reportjob"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

/**
 * startReportTarget bật generated-report flow mới khi migration flag được bật.
 *
 * Default là tắt. Việc bật flag chỉ được làm sau khi proto đã generate và
 * database migrations của target đã được chạy.
 */
func startReportTarget(
	ctx context.Context,
	database *gorm.DB,
	handler *handler_grpc.MapWorkspaceGrpcHandler,
	logger *slog.Logger,
	targetEnabled bool,
) (func(), error) {
	if !targetEnabled {
		return func() {}, nil
	}
	if database == nil || handler == nil {
		return nil, fmt.Errorf("report target database/handler is missing")
	}

	accessConn, accessClient, err := dialAccessService(ctx)
	if err != nil {
		return nil, err
	}

	redisService := _redis.NewRedisService()
	if redisService == nil || redisService.Client == nil {
		_ = accessConn.Close()
		return nil, fmt.Errorf("report target redis client is unavailable")
	}

	generator, closeGenerator, err := reportGeneratorFromEnv(ctx)
	if err != nil {
		_ = accessConn.Close()
		return nil, err
	}

	target, err := newReportTarget(
		database,
		redisService.Client,
		accessClient,
		generator,
		logger,
	)
	if err != nil {
		closeGenerator()
		_ = accessConn.Close()
		return nil, err
	}

	handler.UseCreateReportAdapter(target.adapter)
	target.Start(ctx)

	return func() {
		closeGenerator()
		_ = accessConn.Close()
	}, nil
}

// reportGeneratorFromEnv keeps production behavior unchanged by default.
// Smoke mode exists only to prove the complete local queue/file lifecycle.
func reportGeneratorFromEnv(
	ctx context.Context,
) (reportjob.Generator, func(), error) {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("QHPRO_REPORT_GENERATOR_MODE")))
	if mode == "" || mode == "disabled" || mode == "none" {
		return nil, func() {}, nil
	}
	if mode != "smoke" {
		return nil, nil, fmt.Errorf("unsupported report generator mode %q", mode)
	}

	addr := strings.TrimSpace(os.Getenv("QHPRO_FILE_GRPC_ADDR"))
	if addr == "" {
		addr = "127.0.0.1:8203"
	}
	publicURL := strings.TrimSpace(os.Getenv("QHPRO_FILE_PUBLIC_BASE_URL"))
	if publicURL == "" {
		publicURL = "http://127.0.0.1:8002"
	}

	trust := transporttrust.ConfigFromEnv()
	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(grpcclient.UnaryClientInterceptor(trust)),
		grpc.WithChainStreamInterceptor(grpcclient.StreamClientInterceptor(trust)),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("dial file service for smoke report generator: %w", err)
	}

	uploader := filegrpc.NewClient(filepb.NewInternalFileServiceClient(conn))
	generator, err := reportgenerator.NewSmokeGenerator(uploader, publicURL)
	if err != nil {
		_ = conn.Close()
		return nil, nil, err
	}
	return generator, func() { _ = conn.Close() }, nil
}

func reportTargetEnabled() bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv("QHPRO_REPORT_TARGET_ENABLED")))
	return value == "1" || value == "true" || value == "yes"
}

func dialAccessService(
	ctx context.Context,
) (*grpc.ClientConn, userpb.InternalAccessServiceClient, error) {
	addr := strings.TrimSpace(os.Getenv("QHPRO_USER_GRPC_ADDR"))
	if addr == "" {
		addr = fmt.Sprintf("%s:8201", common.GetPrefixProtobufUrl(env.ENV_RUNTIME, "user"))
	}

	trust := transporttrust.ConfigFromEnv()
	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(grpcclient.UnaryClientInterceptor(trust)),
		grpc.WithChainStreamInterceptor(grpcclient.StreamClientInterceptor(trust)),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("dial user access service: %w", err)
	}
	return conn, userpb.NewInternalAccessServiceClient(conn), nil
}
