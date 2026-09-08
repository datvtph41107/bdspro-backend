package processing_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"common/operation"

	reportfilehttp "tqd/infra/client/file/generatedreport"
	"tqd/infra/worker/generatedreport/generator"
	"tqd/internal/usecase/generatedreport/memoryjob"
	"tqd/internal/usecase/generatedreport/processing"

	postgresdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// blockingFirstGenerator lets worker A finish the real external File effect
// and then stall before Worker.RunOnce can call CompleteJob.
//
// Worker B can therefore reclaim the same durable Job while A is stale.
type blockingFirstGenerator struct {
	underlying processing.Generator

	calls atomic.Int32

	firstDone    chan processing.Output
	releaseFirst chan struct{}
}

func (g *blockingFirstGenerator) GenerateReport(
	ctx context.Context,
	job processing.Job,
) (processing.Output, error) {
	output, err :=
		g.underlying.GenerateReport(
			ctx,
			job,
		)
	if err != nil {
		return processing.Output{}, err
	}

	call :=
		g.calls.Add(1)

	if call == 1 {
		g.firstDone <- output

		select {
		case <-g.releaseFirst:
			return output, nil

		case <-ctx.Done():
			return processing.Output{},
				ctx.Err()
		}
	}

	return output, nil
}

func TestReportLabStaleWorkerCannotFinalizeAfterFileEffect(
	t *testing.T,
) {
	if os.Getenv(
		"QHPRO_REPORT_LAB_INTEGRATION",
	) != "1" {
		t.Skip(
			"report-lab integration test",
		)
	}

	fileDSN :=
		strings.TrimSpace(
			os.Getenv(
				"QHPRO_TEST_FILE_DSN",
			),
		)
	if fileDSN == "" {
		t.Fatal(
			"QHPRO_TEST_FILE_DSN is required",
		)
	}

	fileBaseURL :=
		strings.TrimRight(
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

	serviceAuthKey :=
		strings.TrimSpace(
			os.Getenv(
				"QHPRO_FILE_HTTP_SERVICE_AUTH_KEY",
			),
		)
	if serviceAuthKey == "" {
		t.Fatal(
			"QHPRO_FILE_HTTP_SERVICE_AUTH_KEY is required",
		)
	}

	fileStorageRoot :=
		strings.TrimSpace(
			os.Getenv(
				"QHPRO_TEST_FILE_STORAGE_ROOT",
			),
		)
	if fileStorageRoot == "" {
		t.Fatal(
			"QHPRO_TEST_FILE_STORAGE_ROOT is required",
		)
	}

	fileDB, err :=
		gorm.Open(
			postgresdriver.Open(
				fileDSN,
			),
			&gorm.Config{},
		)
	if err != nil {
		t.Fatalf(
			"connect File PostgreSQL: %v",
			err,
		)
	}

	uploader, err :=
		reportfilehttp.New(
			reportfilehttp.Config{
				BaseURL: fileBaseURL,

				ServiceName: "tqd-service",

				ServiceAuthKey: serviceAuthKey,

				Timeout: 10 * time.Second,
			},
		)
	if err != nil {
		t.Fatalf(
			"create File HTTP client: %v",
			err,
		)
	}
	defer uploader.Close()

	smokeGenerator, err :=
		generator.NewSmokeGenerator(
			uploader,
			uploader,
		)
	if err != nil {
		t.Fatalf(
			"create smoke generator: %v",
			err,
		)
	}

	controlledGenerator :=
		&blockingFirstGenerator{
			underlying: smokeGenerator,

			firstDone: make(
				chan processing.Output,
				1,
			),

			releaseFirst: make(chan struct{}),
		}

	store :=
		memoryjob.NewStore()

	now :=
		time.Now().UTC()

	commandKey :=
		"idem_stale_worker_" +
			strconv.FormatInt(
				now.UnixNano(),
				10,
			)

	job :=
		processing.Job{
			ID: processing.BuildID(
				42,
				commandKey,
			),

			ReportID: 42,
			UserID:   42,

			Operation: operation.Code(
				"workspace.report.generate",
			),

			OperationID: "op_stale_worker_" +
				strconv.FormatInt(
					now.UnixNano(),
					10,
				),

			CommandKey: commandKey,

			Status: processing.StatusPending,

			AvailableAt: now.Add(
				-time.Second,
			),

			CreatedAt: now,
			UpdatedAt: now,
		}

	created, err :=
		store.SaveJob(
			context.Background(),
			job,
		)
	if err != nil {
		t.Fatalf(
			"SaveJob(): %v",
			err,
		)
	}
	if !created {
		t.Fatal(
			"stale-worker Job was not created",
		)
	}

	logger :=
		slog.New(
			slog.NewTextHandler(
				io.Discard,
				nil,
			),
		)

	workerA :=
		processing.NewWorker(
			store,
			controlledGenerator,
			logger,
			"worker-a",
		)

	workerB :=
		processing.NewWorker(
			store,
			controlledGenerator,
			logger,
			"worker-b",
		)

	// --------------------------------------------------------
	// WORKER A
	//
	// Claim v1 and create the real File effect.
	// Generator then stalls before A can CompleteJob.
	// --------------------------------------------------------

	type runResult struct {
		processed bool
		err       error
	}

	aResult :=
		make(chan runResult, 1)

	go func() {
		processed, err :=
			workerA.RunOnce(
				context.Background(),
			)

		aResult <- runResult{
			processed: processed,
			err:       err,
		}
	}()

	firstOutput :=
		<-controlledGenerator.firstDone

	if firstOutput.PDFURL == "" {
		t.Fatal(
			"worker A produced empty PDF URL",
		)
	}

	firstEffect :=
		loadOwnedFileEffect(
			t,
			fileDB,
			job.ID,
		)

	if firstEffect.ID == 0 {
		t.Fatal(
			"worker A did not create File effect",
		)
	}

	claimedA, ok :=
		store.Get(job.ID)
	if !ok {
		t.Fatal(
			"worker A Job disappeared",
		)
	}

	if claimedA.Status !=
		processing.StatusRunning {
		t.Fatalf(
			"worker A Job status=%q, want running",
			claimedA.Status,
		)
	}

	if claimedA.LockedBy !=
		"worker-a" {
		t.Fatalf(
			"worker A lock owner=%q",
			claimedA.LockedBy,
		)
	}

	if claimedA.ClaimVersion != 1 {
		t.Fatalf(
			"worker A claim version=%d, want 1",
			claimedA.ClaimVersion,
		)
	}

	t.Logf(
		"worker A created File effect before becoming stale job=%s claim=%d file_id=%d",
		job.ID,
		claimedA.ClaimVersion,
		firstEffect.ID,
	)

	// --------------------------------------------------------
	// STALE RECOVERY
	//
	// A remains blocked. A tiny stale threshold is deliberate:
	// it proves Worker.ReleaseStale semantics without sleeping 5 minutes.
	// --------------------------------------------------------

	time.Sleep(
		2 * time.Millisecond,
	)

	released, err :=
		workerB.ReleaseStale(
			context.Background(),
			time.Nanosecond,
			100,
		)
	if err != nil {
		t.Fatalf(
			"ReleaseStale(): %v",
			err,
		)
	}

	if released != 1 {
		t.Fatalf(
			"released stale Jobs=%d, want 1",
			released,
		)
	}

	releasedJob, ok :=
		store.Get(job.ID)
	if !ok {
		t.Fatal(
			"released Job disappeared",
		)
	}

	if releasedJob.Status !=
		processing.StatusPending {
		t.Fatalf(
			"released Job status=%q, want pending",
			releasedJob.Status,
		)
	}

	// --------------------------------------------------------
	// WORKER B
	//
	// Reclaims same Job as claim v2.
	// Same owner_key reaches File, so File returns the same effect.
	// B remains the current claim owner and completes.
	// --------------------------------------------------------

	processed, err :=
		workerB.RunOnce(
			context.Background(),
		)
	if err != nil {
		t.Fatalf(
			"worker B RunOnce(): %v",
			err,
		)
	}
	if !processed {
		t.Fatal(
			"worker B did not reclaim Job",
		)
	}

	afterB, ok :=
		store.Get(job.ID)
	if !ok {
		t.Fatal(
			"Job disappeared after worker B",
		)
	}

	if afterB.Status !=
		processing.StatusCompleted {
		t.Fatalf(
			"Job status after B=%q, want completed",
			afterB.Status,
		)
	}

	if afterB.Attempts != 2 {
		t.Fatalf(
			"Job attempts=%d, want 2",
			afterB.Attempts,
		)
	}

	if afterB.ClaimVersion != 2 {
		t.Fatalf(
			"current claim version=%d, want 2",
			afterB.ClaimVersion,
		)
	}

	secondEffect :=
		loadOwnedFileEffect(
			t,
			fileDB,
			job.ID,
		)

	if secondEffect.ID !=
		firstEffect.ID {
		t.Fatalf(
			"File ID changed across stale reclaim: first=%d second=%d",
			firstEffect.ID,
			secondEffect.ID,
		)
	}

	if secondEffect.Path !=
		firstEffect.Path {
		t.Fatalf(
			"File path changed across stale reclaim: first=%q second=%q",
			firstEffect.Path,
			secondEffect.Path,
		)
	}

	if secondEffect.Hash !=
		firstEffect.Hash {
		t.Fatalf(
			"File hash changed across stale reclaim",
		)
	}

	// --------------------------------------------------------
	// STALE WORKER A WAKES
	//
	// A still owns its old in-memory Job snapshot:
	// LockedBy=worker-a, ClaimVersion=1.
	//
	// CompleteJob must be fenced because durable authority moved to
	// worker-b / ClaimVersion=2.
	// --------------------------------------------------------

	close(
		controlledGenerator.releaseFirst,
	)

	resultA :=
		<-aResult

	if !resultA.processed {
		t.Fatal(
			"worker A did not process its claimed Job",
		)
	}

	if !errors.Is(
		resultA.err,
		processing.ErrClaimLost,
	) {
		t.Fatalf(
			"stale worker A error=%v, want ErrClaimLost",
			resultA.err,
		)
	}

	finalJob, ok :=
		store.Get(job.ID)
	if !ok {
		t.Fatal(
			"final Job disappeared",
		)
	}

	if finalJob.Status !=
		processing.StatusCompleted {
		t.Fatalf(
			"stale worker changed completed Job status=%q",
			finalJob.Status,
		)
	}

	if finalJob.ClaimVersion != 2 {
		t.Fatalf(
			"stale worker changed claim version=%d",
			finalJob.ClaimVersion,
		)
	}

	finalEffect :=
		loadOwnedFileEffect(
			t,
			fileDB,
			job.ID,
		)

	if finalEffect.ID !=
		firstEffect.ID {
		t.Fatalf(
			"stale worker created another File effect: first=%d final=%d",
			firstEffect.ID,
			finalEffect.ID,
		)
	}

	// physicalPath := requirePhysicalFileEffect(
	// 	t,
	// 	fileStorageRoot,
	// 	finalEffect.ID,
	// )

	if got :=
		controlledGenerator.calls.Load(); got != 2 {
		t.Fatalf(
			"GenerateReport calls=%d, want 2",
			got,
		)
	}

	t.Logf(
		"PASS stale-worker fence job=%s claimA=1 claimB=2 file_id=%d stale_error=%v",
		job.ID,
		finalEffect.ID,
		resultA.err,
	)
}
