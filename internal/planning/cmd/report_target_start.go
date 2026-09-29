package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	userpb "pb/types/user"
	"tqd/config"
	reportfilehttp "tqd/infra/client/file/generatedreport"
	handler_grpc "tqd/infra/handler/grpc"
	postgresrender "tqd/infra/postgres/generatedreport/render"
	"tqd/infra/worker/generatedreport/generator"
	"tqd/infra/worker/generatedreport/weasyprintpdf"
	"tqd/internal/usecase/generatedreport/processing"
	"tqd/internal/usecase/generatedreport/rendering"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// startReportTarget attaches the canonical Generated Report runtime using only
// process-owned resources. The returned cleanup waits for report actors and
// closes Report-owned adapters; DB/Redis/User RPC remain owned by cmd/grpc.
func startReportTarget(
	ctx context.Context,
	database *gorm.DB,
	redisClient *redis.Client,
	accessClient userpb.InternalAccessServiceClient,
	handler *handler_grpc.MapWorkspaceGrpcHandler,
	cfg config.ReportRuntimeConfig,
	logger *slog.Logger,
) (func(), error) {
	if database == nil || redisClient == nil || accessClient == nil || handler == nil {
		return nil, fmt.Errorf("report target process dependencies are incomplete")
	}

	reportGenerator, closeGenerator, err := reportGeneratorFromConfig(ctx, database, cfg)
	if err != nil {
		return nil, err
	}

	target, err := newReportTarget(
		database,
		redisClient,
		accessClient,
		reportGenerator,
		logger,
		cfg.WorkerID,
	)
	if err != nil {
		closeGenerator()
		return nil, err
	}
	if err := target.Prepare(ctx); err != nil {
		closeGenerator()
		return nil, fmt.Errorf("prepare report quota projection: %w", err)
	}

	handler.UseCreateReportAdapter(target.adapter)
	target.Start(ctx)

	return func() {
		target.Wait()
		closeGenerator()
	}, nil
}

// reportGeneratorFromConfig selects the rendering capability at the process
// composition boundary. Disabled mode keeps non-report TQD capabilities alive;
// canonical report creation then fails closed as unavailable.
func reportGeneratorFromConfig(
	ctx context.Context,
	database *gorm.DB,
	cfg config.ReportRuntimeConfig,
) (processing.Generator, func(), error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	mode := strings.ToLower(strings.TrimSpace(cfg.GeneratorMode))
	if mode == "" || mode == "disabled" || mode == "none" {
		return nil, func() {}, nil
	}
	if mode != "production" && mode != "smoke" {
		return nil, nil, fmt.Errorf("unsupported report generator mode %q", mode)
	}
	if database == nil {
		return nil, nil, fmt.Errorf("report renderer database is missing")
	}
	if strings.TrimSpace(cfg.File.ServiceAuthKey) == "" {
		return nil, nil, fmt.Errorf("report File service auth key is required")
	}

	uploader, err := reportfilehttp.New(reportfilehttp.Config{
		BaseURL:        cfg.File.BaseURL,
		PublicBaseURL:  cfg.File.PublicBaseURL,
		ServiceName:    cfg.File.ServiceName,
		ServiceAuthKey: cfg.File.ServiceAuthKey,
		Timeout:        cfg.File.Timeout,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("configure Report File HTTP client: %w", err)
	}
	closeUploader := func() { uploader.Close() }

	if mode == "smoke" {
		smokeGenerator, err := generator.NewSmokeGenerator(uploader, uploader)
		if err != nil {
			closeUploader()
			return nil, nil, err
		}
		return smokeGenerator, closeUploader, nil
	}

	engine, err := weasyprintpdf.New(cfg.WeasyPrintBinary, cfg.RenderTimeout)
	if err != nil {
		closeUploader()
		return nil, nil, fmt.Errorf("configure production report PDF engine: %w", err)
	}

	productionRenderer, err := rendering.NewService(
		postgresrender.NewSource(database),
		engine,
		uploader,
		uploader,
	)
	if err != nil {
		closeUploader()
		return nil, nil, err
	}
	return productionRenderer, closeUploader, nil
}
