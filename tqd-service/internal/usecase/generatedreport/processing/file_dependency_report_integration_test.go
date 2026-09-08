package processing_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
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

type firstFailureMode string

const (
	firstFailureUnavailable firstFailureMode = "unavailable"
	firstFailureTimeout     firstFailureMode = "timeout"
)

type ownedFileUploader interface {
	PutOwnedFile(
		ctx context.Context,
		ownerNamespace string,
		ownerKey string,
		file io.Reader,
		filename string,
		contentType string,
	) (string, error)
}

// recoveringOwnedUploader injects exactly one dependency failure.
//
// Call #1 uses the fault endpoint.
// Call #2+ use the real File HTTP endpoint.
//
// Both sides are the production reportfilehttp.Client implementation.
// The wrapper owns no HTTP semantics.
type recoveringOwnedUploader struct {
	fault   ownedFileUploader
	healthy ownedFileUploader

	calls          atomic.Int32
	faultCalls     atomic.Int32
	healthyCalls   atomic.Int32
	healthySuccess atomic.Int32
}

func (u *recoveringOwnedUploader) PutOwnedFile(
	ctx context.Context,
	ownerNamespace string,
	ownerKey string,
	file io.Reader,
	filename string,
	contentType string,
) (string, error) {
	call :=
		u.calls.Add(1)

	if call == 1 {
		u.faultCalls.Add(1)

		return u.fault.PutOwnedFile(
			ctx,
			ownerNamespace,
			ownerKey,
			file,
			filename,
			contentType,
		)
	}

	u.healthyCalls.Add(1)

	path, err :=
		u.healthy.PutOwnedFile(
			ctx,
			ownerNamespace,
			ownerKey,
			file,
			filename,
			contentType,
		)

	if err == nil {
		u.healthySuccess.Add(1)
	}

	return path, err
}

type faultFileServer struct {
	mode firstFailureMode

	requests atomic.Int32
}

func (s *faultFileServer) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	s.requests.Add(1)

	// Always consume the request body.
	//
	// This matters because the production HTTP client streams multipart
	// through io.Pipe. The fault fixture must not accidentally introduce
	// a blocked request-body writer.
	_, _ =
		io.Copy(
			io.Discard,
			r.Body,
		)

	switch s.mode {
	case firstFailureUnavailable:
		http.Error(
			w,
			"report-lab File unavailable",
			http.StatusServiceUnavailable,
		)

	case firstFailureTimeout:
		// The production client timeout in this test is 75ms.
		// Bound the fixture itself as well so httptest shutdown can never
		// wait indefinitely even if cancellation propagation changes.
		timer :=
			time.NewTimer(
				250 * time.Millisecond,
			)
		defer timer.Stop()

		select {
		case <-r.Context().Done():
			return

		case <-timer.C:
			http.Error(
				w,
				"report-lab delayed File response",
				http.StatusGatewayTimeout,
			)
		}

	default:
		http.Error(
			w,
			"unsupported fault mode",
			http.StatusInternalServerError,
		)
	}
}

func countOwnedFileEffects(
	t *testing.T,
	db *gorm.DB,
	ownerKey string,
) int64 {
	t.Helper()

	var count int64

	err :=
		db.Raw(
			`
SELECT count(*)
FROM file
WHERE owner_namespace = ?
  AND owner_key = ?
`,
			processing.FileOwnerNamespace,
			ownerKey,
		).
			Scan(&count).
			Error
	if err != nil {
		t.Fatalf(
			"count owned File effects: %v",
			err,
		)
	}

	return count
}

func runRecoveringFileDependencyCase(
	t *testing.T,
	mode firstFailureMode,
	faultTimeout time.Duration,
) {
	t.Helper()

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

	// ========================================================
	// FAILURE ENDPOINT
	// ========================================================

	faultHandler :=
		&faultFileServer{
			mode: mode,
		}

	faultServer :=
		httptest.NewServer(
			faultHandler,
		)
	defer faultServer.Close()

	faultClient, err :=
		reportfilehttp.New(
			reportfilehttp.Config{
				BaseURL: faultServer.URL,

				ServiceName: "tqd-service",

				ServiceAuthKey: serviceAuthKey,

				Timeout: faultTimeout,
			},
		)
	if err != nil {
		t.Fatalf(
			"create fault File client: %v",
			err,
		)
	}
	defer faultClient.Close()

	// ========================================================
	// HEALTHY ENDPOINT
	//
	// Direct production HTTP adapter → real File :8002.
	// No forwarding proxy sits in this path.
	// ========================================================

	healthyClient, err :=
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
			"create healthy File client: %v",
			err,
		)
	}
	defer healthyClient.Close()

	uploader :=
		&recoveringOwnedUploader{
			fault: faultClient,

			healthy: healthyClient,
		}

	reportGenerator, err :=
		generator.NewSmokeGenerator(
			uploader,
			healthyClient,
		)
	if err != nil {
		t.Fatalf(
			"create smoke generator: %v",
			err,
		)
	}

	store :=
		&immediateRetryStore{
			Store: memoryjob.NewStore(),
		}

	now :=
		time.Now().UTC()

	commandKey :=
		"idem_file_" +
			string(mode) +
			"_" +
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

			OperationID: "op_file_" +
				string(mode) +
				"_" +
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
			"dependency-failure Job was not created",
		)
	}

	worker :=
		processing.NewWorker(
			store,
			reportGenerator,
			slog.New(
				slog.NewTextHandler(
					io.Discard,
					nil,
				),
			),
			"g02-"+string(mode)+"-worker",
		)

	// ========================================================
	// ATTEMPT 1 — dependency failure.
	// ========================================================

	start :=
		time.Now()

	processed, err :=
		worker.RunOnce(
			context.Background(),
		)

	elapsed :=
		time.Since(start)

	if err != nil {
		t.Fatalf(
			"first Worker.RunOnce(): %v",
			err,
		)
	}
	if !processed {
		t.Fatal(
			"first Worker.RunOnce() did not claim Job",
		)
	}

	if got :=
		uploader.calls.Load(); got != 1 {
		t.Fatalf(
			"upload calls after attempt1=%d, want 1",
			got,
		)
	}

	if got :=
		uploader.faultCalls.Load(); got != 1 {
		t.Fatalf(
			"fault File calls=%d, want 1",
			got,
		)
	}

	if got :=
		uploader.healthyCalls.Load(); got != 0 {
		t.Fatalf(
			"healthy File calls after failure=%d, want 0",
			got,
		)
	}

	if got :=
		faultHandler.requests.Load(); got != 1 {
		t.Fatalf(
			"fault HTTP requests=%d, want 1",
			got,
		)
	}

	if store.retryCalls != 1 {
		t.Fatalf(
			"RetryJob calls=%d, want 1",
			store.retryCalls,
		)
	}

	if store.completeCalls != 0 {
		t.Fatalf(
			"CompleteJob calls=%d after File failure, want 0",
			store.completeCalls,
		)
	}

	afterFailure, ok :=
		store.Get(job.ID)
	if !ok {
		t.Fatal(
			"Job disappeared after File failure",
		)
	}

	if afterFailure.Status !=
		processing.StatusPending {
		t.Fatalf(
			"Job status after failure=%q, want pending",
			afterFailure.Status,
		)
	}

	if afterFailure.Attempts != 1 {
		t.Fatalf(
			"Job attempts after failure=%d, want 1",
			afterFailure.Attempts,
		)
	}

	if strings.TrimSpace(
		afterFailure.LastError,
	) == "" {
		t.Fatal(
			"Job LastError is empty after File failure",
		)
	}

	if got :=
		countOwnedFileEffects(
			t,
			fileDB,
			job.ID,
		); got != 0 {
		t.Fatalf(
			"File effects after failed attempt=%d, want 0",
			got,
		)
	}

	if mode ==
		firstFailureTimeout {
		if elapsed <
			faultTimeout/2 {
			t.Fatalf(
				"timeout returned too early: elapsed=%s timeout=%s",
				elapsed,
				faultTimeout,
			)
		}
	}

	t.Logf(
		"attempt1 mode=%s remained retryable elapsed=%s error=%q",
		mode,
		elapsed,
		afterFailure.LastError,
	)

	// ========================================================
	// ATTEMPT 2 — dependency recovered.
	// ========================================================

	processed, err =
		worker.RunOnce(
			context.Background(),
		)
	if err != nil {
		t.Fatalf(
			"second Worker.RunOnce(): %v",
			err,
		)
	}
	if !processed {
		t.Fatal(
			"second Worker.RunOnce() did not reclaim Job",
		)
	}

	if got :=
		uploader.calls.Load(); got != 2 {
		t.Fatalf(
			"total upload calls=%d, want 2",
			got,
		)
	}

	if got :=
		uploader.faultCalls.Load(); got != 1 {
		t.Fatalf(
			"fault File calls=%d, want exactly 1",
			got,
		)
	}

	if got :=
		uploader.healthyCalls.Load(); got != 1 {
		t.Fatalf(
			"healthy File calls=%d, want 1",
			got,
		)
	}

	if got :=
		uploader.healthySuccess.Load(); got != 1 {
		t.Fatalf(
			"successful healthy File calls=%d, want 1",
			got,
		)
	}

	if store.retryCalls != 1 {
		t.Fatalf(
			"RetryJob calls=%d, want exactly 1",
			store.retryCalls,
		)
	}

	if store.completeCalls != 1 {
		t.Fatalf(
			"CompleteJob calls=%d, want exactly 1",
			store.completeCalls,
		)
	}

	finalJob, ok :=
		store.Get(job.ID)
	if !ok {
		t.Fatal(
			"Job disappeared after recovery",
		)
	}

	if finalJob.Status !=
		processing.StatusCompleted {
		t.Fatalf(
			"final Job status=%q, want completed",
			finalJob.Status,
		)
	}

	if finalJob.Attempts != 2 {
		t.Fatalf(
			"final Job attempts=%d, want 2",
			finalJob.Attempts,
		)
	}

	if finalJob.LastError != "" {
		t.Fatalf(
			"completed Job retained LastError=%q",
			finalJob.LastError,
		)
	}

	if got :=
		countOwnedFileEffects(
			t,
			fileDB,
			job.ID,
		); got != 1 {
		t.Fatalf(
			"File effects after recovery=%d, want 1",
			got,
		)
	}

	effect :=
		loadOwnedFileEffect(
			t,
			fileDB,
			job.ID,
		)

	t.Logf(
		"PASS File dependency recovery mode=%s job=%s attempts=2 file_id=%d",
		mode,
		job.ID,
		effect.ID,
	)
}

func TestReportLabWorkerRetriesWhenFileUnavailable(
	t *testing.T,
) {
	runRecoveringFileDependencyCase(
		t,
		firstFailureUnavailable,
		5*time.Second,
	)
}

func TestReportLabWorkerRetriesWhenFileTimesOut(
	t *testing.T,
) {
	runRecoveringFileDependencyCase(
		t,
		firstFailureTimeout,
		75*time.Millisecond,
	)
}
