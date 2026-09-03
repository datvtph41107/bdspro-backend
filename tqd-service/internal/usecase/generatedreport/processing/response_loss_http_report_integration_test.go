package processing_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

type immediateRetryStore struct {
	*memoryjob.Store

	retryCalls    int
	completeCalls int
}

func (s *immediateRetryStore) RetryJob(
	ctx context.Context,
	job processing.Job,
	lastError string,
	_ time.Time,
) error {
	s.retryCalls++

	// Production Worker keeps its normal retry policy.
	//
	// Only the lab Store projection makes this already-retried Job
	// claimable immediately so the proof does not sleep 30 seconds.
	return s.Store.RetryJob(
		ctx,
		job,
		lastError,
		time.Now().UTC().Add(-time.Second),
	)
}

func (s *immediateRetryStore) CompleteJob(
	ctx context.Context,
	job processing.Job,
	output processing.Output,
	now time.Time,
) error {
	if err := s.Store.CompleteJob(
		ctx,
		job,
		output,
		now,
	); err != nil {
		return err
	}

	s.completeCalls++
	return nil
}

// responseLossProxy forwards the request to the real File HTTP service.
//
// The first successful backend request deliberately loses the usable
// response body. The File effect is already durable, but the TQD caller
// cannot decode the success result and must treat the outcome as unknown.
type responseLossProxy struct {
	target string
	client *http.Client

	attempts        atomic.Int32
	backendSuccess  atomic.Int32
	droppedResponse atomic.Int32
}

func (p *responseLossProxy) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	attempt :=
		p.attempts.Add(1)

	target :=
		strings.TrimRight(
			p.target,
			"/",
		) +
			r.URL.RequestURI()

	req, err :=
		http.NewRequestWithContext(
			r.Context(),
			r.Method,
			target,
			r.Body,
		)
	if err != nil {
		http.Error(
			w,
			"cannot build upstream request",
			http.StatusBadGateway,
		)
		return
	}

	req.Header =
		r.Header.Clone()

	resp, err :=
		p.client.Do(req)
	if err != nil {
		http.Error(
			w,
			"File upstream unavailable",
			http.StatusBadGateway,
		)
		return
	}
	defer resp.Body.Close()

	body, err :=
		io.ReadAll(resp.Body)
	if err != nil {
		http.Error(
			w,
			"cannot read File upstream response",
			http.StatusBadGateway,
		)
		return
	}

	if resp.StatusCode == http.StatusOK {
		p.backendSuccess.Add(1)
	}

	if attempt == 1 &&
		resp.StatusCode == http.StatusOK {
		p.droppedResponse.Add(1)

		// Returning without writing preserves an HTTP 200 response with
		// an empty body. The downstream client therefore cannot decode
		// the successful File result even though File already committed.
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(
				key,
				value,
			)
		}
	}

	w.WriteHeader(
		resp.StatusCode,
	)

	_, _ =
		w.Write(body)
}

type ownedFileRow struct {
	ID uint64

	OwnerNamespace   string
	OwnerKey         string
	OwnerRequestHash string

	Name string
	Path string
	Hash string
	Size int64
}

func loadOwnedFileEffect(
	t *testing.T,
	db *gorm.DB,
	ownerKey string,
) ownedFileRow {
	t.Helper()

	var rows []ownedFileRow

	err :=
		db.Raw(
			`
SELECT
    id,
    owner_namespace,
    owner_key,
    owner_request_hash,
    name,
    path,
    hash,
    size
FROM file
WHERE owner_namespace = ?
  AND owner_key = ?
ORDER BY id
`,
			processing.FileOwnerNamespace,
			ownerKey,
		).
			Scan(&rows).
			Error
	if err != nil {
		t.Fatalf(
			"query owned File effect: %v",
			err,
		)
	}

	if len(rows) != 1 {
		t.Fatalf(
			"owned File rows=%d, want 1",
			len(rows),
		)
	}

	return rows[0]
}

// requirePhysicalFileEffect proves the File DB effect has one physical
// artifact without reading any deployment-specific path from durable metadata.
// File owns path resolution; this cross-service lab only observes the configured
// storage root and the durable File ID.
func requirePhysicalFileEffect(
	t *testing.T,
	root string,
	fileID uint64,
) string {
	t.Helper()

	wantBase := strconv.FormatUint(fileID, 10)
	var matches []string

	err := filepath.Walk(
		root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !info.Mode().IsRegular() {
				return nil
			}

			base := filepath.Base(path)
			if strings.TrimSuffix(base, filepath.Ext(base)) == wantBase {
				matches = append(matches, path)
			}
			return nil
		},
	)
	if err != nil {
		t.Fatalf("walk File storage root: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf(
			"physical File artifacts for file_id=%d: got %d, want 1 (%v)",
			fileID,
			len(matches),
			matches,
		)
	}

	return matches[0]
}

func TestReportLabWorkerConvergesAfterLostHTTPFileResponse(
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

	proxy :=
		&responseLossProxy{
			target: fileBaseURL,
			client: &http.Client{
				Timeout: 10 * time.Second,
			},
		}

	proxyServer :=
		httptest.NewServer(proxy)
	defer proxyServer.Close()

	uploader, err :=
		reportfilehttp.New(
			reportfilehttp.Config{
				BaseURL: proxyServer.URL,

				ServiceName: "tqd-service",

				ServiceAuthKey: serviceAuthKey,

				Timeout: 10 * time.Second,
			},
		)
	if err != nil {
		t.Fatalf(
			"create Report File HTTP client: %v",
			err,
		)
	}
	defer uploader.Close()

	reportGenerator, err :=
		generator.NewSmokeGenerator(
			uploader,
			uploader,
		)
	if err != nil {
		t.Fatalf(
			"create smoke report generator: %v",
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
		"idem_response_loss_" +
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

			OperationID: "op_response_loss_" +
				strconv.FormatInt(
					now.UnixNano(),
					10,
				),

			CommandKey: commandKey,

			Status: processing.StatusPending,

			AvailableAt: now.Add(-time.Second),

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
			"response-loss Job was not created",
		)
	}

	logger :=
		slog.New(
			slog.NewTextHandler(
				io.Discard,
				nil,
			),
		)

	worker :=
		processing.NewWorker(
			store,
			reportGenerator,
			logger,
			"g02-response-loss-worker",
		)

	// --------------------------------------------------------
	// ATTEMPT 1
	//
	// File succeeds durably.
	// Proxy loses the usable success response.
	// Worker must schedule retry.
	// --------------------------------------------------------

	processed, err :=
		worker.RunOnce(
			context.Background(),
		)
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
		proxy.attempts.Load(); got != 1 {
		t.Fatalf(
			"HTTP attempts after first run=%d, want 1",
			got,
		)
	}

	if got :=
		proxy.backendSuccess.Load(); got != 1 {
		t.Fatalf(
			"successful File backend calls=%d, want 1",
			got,
		)
	}

	if got :=
		proxy.droppedResponse.Load(); got != 1 {
		t.Fatalf(
			"dropped File responses=%d, want 1",
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
			"CompleteJob calls=%d after ambiguous outcome, want 0",
			store.completeCalls,
		)
	}

	firstEffect :=
		loadOwnedFileEffect(
			t,
			fileDB,
			job.ID,
		)

	if firstEffect.ID == 0 ||
		firstEffect.Path == "" ||
		firstEffect.Hash == "" ||
		firstEffect.OwnerRequestHash == "" {
		t.Fatalf(
			"incomplete File effect after lost response: %+v",
			firstEffect,
		)
	}

	t.Logf(
		"attempt1: File committed before response loss job=%s file_id=%d",
		job.ID,
		firstEffect.ID,
	)

	// --------------------------------------------------------
	// ATTEMPT 2
	//
	// Same durable Job identity retries.
	// File must return the already-owned effect.
	// Worker can now complete.
	// --------------------------------------------------------

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
		proxy.attempts.Load(); got != 2 {
		t.Fatalf(
			"HTTP attempts=%d, want 2",
			got,
		)
	}

	if got :=
		proxy.backendSuccess.Load(); got != 2 {
		t.Fatalf(
			"successful File backend calls=%d, want 2",
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

	savedJob, ok :=
		store.Get(job.ID)
	if !ok {
		t.Fatal(
			"completed Job disappeared",
		)
	}

	if savedJob.Status !=
		processing.StatusCompleted {
		t.Fatalf(
			"Job status=%q, want completed",
			savedJob.Status,
		)
	}

	if savedJob.Attempts != 2 {
		t.Fatalf(
			"Job attempts=%d, want 2",
			savedJob.Attempts,
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
			"File ID changed across retry: first=%d second=%d",
			firstEffect.ID,
			secondEffect.ID,
		)
	}

	if secondEffect.Path !=
		firstEffect.Path {
		t.Fatalf(
			"File path changed across retry: first=%q second=%q",
			firstEffect.Path,
			secondEffect.Path,
		)
	}

	if secondEffect.Hash !=
		firstEffect.Hash {
		t.Fatalf(
			"File content hash changed across retry: first=%q second=%q",
			firstEffect.Hash,
			secondEffect.Hash,
		)
	}

	physicalPath := requirePhysicalFileEffect(
		t,
		fileStorageRoot,
		secondEffect.ID,
	)

	pattern :=
		filepath.Join(
			filepath.Dir(
				physicalPath,
			),
			strconv.FormatUint(
				secondEffect.ID,
				10,
			)+".*",
		)

	matches, err :=
		filepath.Glob(
			pattern,
		)
	if err != nil {
		t.Fatalf(
			"inspect physical File effects: %v",
			err,
		)
	}

	if len(matches) != 1 {
		t.Fatalf(
			"physical File effects=%d, want 1: %#v",
			len(matches),
			matches,
		)
	}

	// Completed Job must not be claimable again.
	processed, err =
		worker.RunOnce(
			context.Background(),
		)
	if err != nil {
		t.Fatalf(
			"third Worker.RunOnce(): %v",
			err,
		)
	}
	if processed {
		t.Fatal(
			"completed Job became claimable again",
		)
	}

	t.Logf(
		"PASS response-loss convergence job=%s attempts=2 file_id=%d physical=%s",
		job.ID,
		secondEffect.ID,
		physicalPath,
	)
}
