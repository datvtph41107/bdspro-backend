package cmd

import (
	"common"
	_redis "common/redis"
	qhprorpc "common/rpc"
	"common/transporttrust"
	"context"
	"fmt"
	"log/slog"
	"os"
	userpb "pb/types/user"
	"strings"
	"time"
	"tqd/env"
	handler_grpc "tqd/infra/handler/grpc"
	reportfilehttp "tqd/internal/workspace/report/adapters/filehttp"
	"tqd/internal/workspace/report/adapters/generator"
	"tqd/internal/workspace/report/adapters/postgresrender"
	"tqd/internal/workspace/report/adapters/weasyprintpdf"
	"tqd/internal/workspace/report/processing"
	"tqd/internal/workspace/report/rendering"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

// startReportTarget attempts to attach the canonical Generated Report runtime.
// Failure is capability-local: the caller may keep the rest of TQD running,
// while CreateGeneratedReport fails closed with UNAVAILABLE.
func startReportTarget(
	ctx context.Context,
	database *gorm.DB,
	handler *handler_grpc.MapWorkspaceGrpcHandler,
	logger *slog.Logger,
) (func(), error) {
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
	closeRedis := func() { _ = redisService.Client.Close() }

	generator, closeGenerator, err := reportGeneratorFromEnv(ctx, database)
	if err != nil {
		closeRedis()
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
		closeRedis()
		_ = accessConn.Close()
		return nil, err
	}
	if err := target.Prepare(ctx); err != nil {
		closeGenerator()
		closeRedis()
		_ = accessConn.Close()
		return nil, fmt.Errorf("prepare report quota projection: %w", err)
	}

	handler.UseCreateReportAdapter(target.adapter)
	target.Start(ctx)

	return func() {
		// The runners may still be using File/User clients after cancellation.
		// Wait is therefore the dependency-close barrier.
		target.Wait()
		closeGenerator()
		closeRedis()
		_ = accessConn.Close()
	}, nil
}

// reportGeneratorFromEnv selects an explicit rendering capability at the
// composition boundary. Empty/disabled preserves capability-local UNAVAILABLE;
// production uses the accepted PostgreSQL Report snapshot + WeasyPrint + File gRPC.
// Smoke remains a lab-only implementation for local end-to-end tests.
func reportGeneratorFromEnv(
	ctx context.Context,
	database *gorm.DB,
) (processing.Generator, func(), error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	mode := strings.ToLower(
		strings.TrimSpace(
			os.Getenv("QHPRO_REPORT_GENERATOR_MODE"),
		),
	)

	if mode == "" ||
		mode == "disabled" ||
		mode == "none" {
		return nil, func() {}, nil
	}

	if mode != "production" &&
		mode != "smoke" {
		return nil, nil, fmt.Errorf(
			"unsupported report generator mode %q",
			mode,
		)
	}

	if mode == "production" &&
		database == nil {
		return nil, nil, fmt.Errorf(
			"production report renderer database is missing",
		)
	}

	fileBaseURL := strings.TrimRight(
		strings.TrimSpace(
			os.Getenv(
				"QHPRO_FILE_HTTP_BASE_URL",
			),
		),
		"/",
	)

	if fileBaseURL == "" {
		fileBaseURL =
			"http://127.0.0.1:8002"
	}

	publicURL := strings.TrimRight(
		strings.TrimSpace(
			os.Getenv(
				"QHPRO_FILE_PUBLIC_BASE_URL",
			),
		),
		"/",
	)

	if publicURL == "" {
		publicURL = fileBaseURL
	}

	serviceAuthKey := strings.TrimSpace(
		os.Getenv(
			"QHPRO_FILE_HTTP_SERVICE_AUTH_KEY",
		),
	)

	if serviceAuthKey == "" {
		return nil, nil, fmt.Errorf(
			"QHPRO_FILE_HTTP_SERVICE_AUTH_KEY is required",
		)
	}

	serviceName := strings.TrimSpace(
		os.Getenv("QHPRO_SERVICE_ID"),
	)

	if serviceName == "" {
		serviceName = "tqd-service"
	}

	uploadTimeout := 30 * time.Second

	if raw := strings.TrimSpace(
		os.Getenv("QHPRO_FILE_HTTP_TIMEOUT"),
	); raw != "" {
		parsed, err :=
			time.ParseDuration(raw)

		if err != nil ||
			parsed <= 0 {
			return nil, nil, fmt.Errorf(
				"invalid QHPRO_FILE_HTTP_TIMEOUT %q",
				raw,
			)
		}

		uploadTimeout = parsed
	}

	uploader, err :=
		reportfilehttp.New(
			reportfilehttp.Config{
				BaseURL: fileBaseURL,

				ServiceName: serviceName,

				ServiceAuthKey: serviceAuthKey,

				Timeout: uploadTimeout,
			},
		)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"configure Report File HTTP client: %w",
			err,
		)
	}

	closeUploader := func() {
		uploader.Close()
	}

	if mode == "smoke" {
		smokeGenerator, err :=
			generator.NewSmokeGenerator(
				uploader,
				publicURL,
			)
		if err != nil {
			closeUploader()
			return nil, nil, err
		}

		return smokeGenerator,
			closeUploader,
			nil
	}

	binary := strings.TrimSpace(
		os.Getenv(
			"QHPRO_REPORT_WEASYPRINT_BINARY",
		),
	)

	if binary == "" {
		binary = "weasyprint"
	}

	renderTimeout :=
		45 * time.Second

	if raw := strings.TrimSpace(
		os.Getenv(
			"QHPRO_REPORT_RENDER_TIMEOUT",
		),
	); raw != "" {
		parsed, err :=
			time.ParseDuration(raw)

		if err != nil ||
			parsed <= 0 {
			closeUploader()

			return nil, nil, fmt.Errorf(
				"invalid QHPRO_REPORT_RENDER_TIMEOUT %q",
				raw,
			)
		}

		renderTimeout = parsed
	}

	engine, err :=
		weasyprintpdf.New(
			binary,
			renderTimeout,
		)
	if err != nil {
		closeUploader()

		return nil, nil, fmt.Errorf(
			"configure production report PDF engine: %w",
			err,
		)
	}

	productionRenderer, err :=
		rendering.NewService(
			postgresrender.NewSource(
				database,
			),
			engine,
			uploader,
			publicURL,
		)
	if err != nil {
		closeUploader()
		return nil, nil, err
	}

	return productionRenderer,
		closeUploader,
		nil
}

func dialAccessService(
	ctx context.Context,
) (*grpc.ClientConn, userpb.InternalAccessServiceClient, error) {
	addr := strings.TrimSpace(os.Getenv("QHPRO_USER_GRPC_ADDR"))
	if addr == "" {
		addr = fmt.Sprintf("%s:8201", common.GetPrefixProtobufUrl(env.ENV_RUNTIME, "user"))
	}

	trust := transporttrust.ConfigFromEnv()
	conn, err := qhprorpc.NewClient(qhprorpc.ClientConfig{
		Target:          addr,
		BackoffMaxDelay: 5 * time.Second,
		Credentials:     insecure.NewCredentials(),
		Transport:       qhprorpc.TransportConfig{Trust: trust},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("dial user access service: %w", err)
	}
	return conn, userpb.NewInternalAccessServiceClient(conn), nil
}
