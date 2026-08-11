package cmd

import (
	"context"
	"errors"
	"log/slog"
	"time"
	legacypostgres "tqd/infra/postgres"
	reportpostgres "tqd/infra/postgres/report"
	reportjobpostgres "tqd/infra/postgres/reportjob"
	usageeventpostgres "tqd/infra/postgres/usageevent"
	"tqd/infra/redisquota"
	"tqd/infra/reportgrpc"
	"tqd/infra/useraccess"
	"tqd/internal/quota"
	"tqd/internal/quotacleanup"
	"tqd/internal/report"
	"tqd/internal/reportjob"
	"tqd/internal/usagecheck"

	userpb "pb/types/user"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

/**
 * reportTarget gom các thành phần của generated-report flow mới.
 *
 * Struct này chỉ dùng ở composition root, business code không import cmd package.
 */
type reportTarget struct {
	adapter       *reportgrpc.Adapter
	jobRunner     *reportjob.Runner
	cleanupRunner *quotacleanup.Runner
	usageRunner   *usagecheck.Runner
}

/**
 * newReportTarget lắp generated-report flow từ concrete DB/Redis/gRPC implementations.
 *
 * generator là phần render PDF/image thật. Có thể để nil trong giai đoạn đầu:
 * request vẫn accepted và report_job được lưu pending, chỉ worker render chưa chạy.
 */
func newReportTarget(
	database *gorm.DB,
	redisClient *redis.Client,
	accessGRPC userpb.InternalAccessServiceClient,
	generator reportjob.Generator,
	logger *slog.Logger,
) (*reportTarget, error) {
	if database == nil || redisClient == nil || accessGRPC == nil {
		return nil, errors.New("generated report target dependencies are incomplete")
	}
	if logger == nil {
		logger = slog.Default()
	}

	workspaceStore := &legacypostgres.MapWorkspacePostgres{DB: database}
	sourceStore := reportpostgres.NewSourceStore(workspaceStore)
	reportStore := reportpostgres.NewReportStore(database)
	accessClient := useraccess.NewClient(accessGRPC)

	quotaStore := redisquota.NewStore(redisClient)
	quotaService := quota.NewService(quotaStore)

	reportService := report.NewService(
		sourceStore,
		reportStore,
		accessClient,
		quotaService,
		logger,
	)

	var jobRunner *reportjob.Runner
	if generator != nil {
		jobStore := reportjobpostgres.NewStore(database)
		jobWorker := reportjob.NewWorker(
			jobStore,
			generator,
			logger,
			"tqd-report-worker",
		)
		jobRunner = reportjob.NewRunner(jobWorker, time.Second, logger)
	}

	usageStore := usageeventpostgres.NewStore(database)
	cleanupService := quotacleanup.NewService(quotaStore, usageStore)
	usageService := usagecheck.NewService(quotaStore, usageStore)

	return &reportTarget{
		adapter:       reportgrpc.NewAdapter(reportService),
		jobRunner:     jobRunner,
		cleanupRunner: quotacleanup.NewRunner(cleanupService, time.Minute, 100, logger),
		usageRunner:   usagecheck.NewRunner(usageService, usageStore, 5*time.Minute, 200, logger),
	}, nil
}

/**
 * Start chạy các background loop đi cùng report flow.
 */
func (r *reportTarget) Start(ctx context.Context) {
	if r == nil {
		return
	}
	if r.jobRunner != nil {
		go r.jobRunner.Run(ctx)
	}
	go r.cleanupRunner.Run(ctx)
	go r.usageRunner.Run(ctx)
}
