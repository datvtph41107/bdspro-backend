package cmd

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"
	"tqd/infra/client/user/generatedreport"
	"tqd/infra/handler/grpc/generatedreport"
	legacypostgres "tqd/infra/postgres"
	"tqd/infra/postgres/generatedreport/job"
	"tqd/infra/postgres/generatedreport/report"
	"tqd/infra/postgres/usage"
	"tqd/infra/redis/quota"
	"tqd/internal/usecase/generatedreport/application"
	"tqd/internal/usecase/generatedreport/processing"
	"tqd/internal/usecase/quota"
	"tqd/internal/usecase/quota/reconciliation/reservationcleanup"
	"tqd/internal/usecase/quota/reconciliation/usageprojection"

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
	adapter       *grpcadapter.Adapter
	jobRunner     *processing.Runner
	cleanupRunner *reservationcleanup.Runner
	usageRunner   *usageprojection.Runner

	wg sync.WaitGroup
}

/**
 * Prepare phục hồi Redis projection từ durable ledger trước khi nhận traffic.
 */
func (r *reportTarget) Prepare(ctx context.Context) error {
	if r == nil || r.usageRunner == nil {
		return errors.New("generated report usage projection is not configured")
	}
	return r.usageRunner.ReconcileAll(ctx)
}

/**
 * newReportTarget lắp generated-report flow từ concrete DB/Redis/gRPC implementations.
 *
 * generator là phần render PDF/image thật. Nếu nil, durable replay vẫn đọc được
 * Report đã accepted trước đó nhưng command mới fail-closed vì processing unavailable.
 */
func newReportTarget(
	database *gorm.DB,
	redisClient *redis.Client,
	accessGRPC userpb.InternalAccessServiceClient,
	generator processing.Generator,
	logger *slog.Logger,
	workerID string,
) (*reportTarget, error) {
	if database == nil || redisClient == nil || accessGRPC == nil {
		return nil, errors.New("generated report target dependencies are incomplete")
	}
	if logger == nil {
		logger = slog.Default()
	}

	workspaceStore := &legacypostgres.MapWorkspacePostgres{DB: database}
	sourceStore := postgresreport.NewSourceStore(workspaceStore)
	reportStore := postgresreport.NewReportStore(database)
	accessClient := useraccess.NewClient(accessGRPC)

	quotaStore := redisquota.NewStore(redisClient)
	quotaService := quota.NewService(quotaStore)

	reportService := application.NewService(
		sourceStore,
		reportStore,
		accessClient,
		quotaService,
		generator != nil,
		logger,
	)

	var jobRunner *processing.Runner
	if generator != nil {
		jobStore := postgresjob.NewStore(database)
		if workerID == "" {
			workerID = "tqd-report-worker"
		}
		jobWorker := processing.NewWorker(
			jobStore,
			generator,
			logger,
			workerID,
		)
		jobRunner = processing.NewRunner(jobWorker, time.Second, logger)
	}

	usageStore := postgresusage.NewStore(database)
	cleanupService := reservationcleanup.NewService(quotaStore, usageStore)
	usageService := usageprojection.NewService(quotaStore, usageStore)

	return &reportTarget{
		adapter:       grpcadapter.NewAdapter(reportService),
		jobRunner:     jobRunner,
		cleanupRunner: reservationcleanup.NewRunner(cleanupService, time.Minute, 100, logger),
		usageRunner:   usageprojection.NewRunner(usageService, usageStore, 5*time.Minute, 200, logger),
	}, nil
}

// Start launches the Report-owned actors but keeps their lifecycle observable
// by the composition root. Dependencies must not be closed until Wait returns.
func (r *reportTarget) Start(ctx context.Context) {
	if r == nil {
		return
	}
	if r.jobRunner != nil {
		r.start(ctx, r.jobRunner.Run)
	}
	if r.cleanupRunner != nil {
		r.start(ctx, r.cleanupRunner.Run)
	}
	if r.usageRunner != nil {
		r.start(ctx, r.usageRunner.Run)
	}
}

func (r *reportTarget) start(ctx context.Context, run func(context.Context)) {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		run(ctx)
	}()
}

// Wait blocks until every Report-owned actor has observed cancellation and
// returned. This is the dependency-close barrier used by process shutdown.
func (r *reportTarget) Wait() {
	if r == nil {
		return
	}
	r.wg.Wait()
}
