package postgresjob_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"common/operation"

	reportfilehttp "tqd/infra/client/file/generatedreport"
	"tqd/infra/postgres/generatedreport/job"
	"tqd/infra/worker/generatedreport/generator"
	"tqd/internal/enums"
	"tqd/internal/usecase/generatedreport/processing"

	postgresdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type durableReportState struct {
	Status       int
	PDFURL       string
	FileSize     uint64
	ErrorMessage string
}

type durableFileEffect struct {
	ID               uint64
	OwnerNamespace   string
	OwnerKey         string
	OwnerRequestHash string
	Path             string
	Hash             string
	Size             int64
}

func openReportLabPostgres(
	t *testing.T,
	dsn string,
) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(
		postgresdriver.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		t.Fatalf(
			"open report-lab PostgreSQL: %v",
			err,
		)
	}

	return db
}

func insertProcessingReport(
	t *testing.T,
	db *gorm.DB,
	userID uint64,
	commandKey string,
	jobID string,
) uint64 {
	t.Helper()

	reportType :=
		enums.GeneratedReportType(1)

	if !reportType.IsValid() {
		t.Fatal(
			"GeneratedReportType(1) is no longer valid",
		)
	}

	var reportID uint64

	err := db.Raw(
		`
INSERT INTO user_reported (
    user_id,
    report_type,
    status,
    title,
    command_key,
    job_id,
    request_hash
)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING id
`,
		userID,
		reportType,
		enums.GeneratedReportStatusProcessing,
		"G02 P7 durable recovery",
		commandKey,
		jobID,
		strings.Repeat("a", 64),
	).Scan(&reportID).Error
	if err != nil {
		t.Fatalf(
			"insert processing report: %v",
			err,
		)
	}

	if reportID == 0 {
		t.Fatal(
			"inserted report ID is zero",
		)
	}

	return reportID
}

func loadReportState(
	t *testing.T,
	db *gorm.DB,
	reportID uint64,
) durableReportState {
	t.Helper()

	var state durableReportState

	err := db.Raw(
		`
SELECT
    status,
    pdf_url,
    file_size,
    error_message
FROM user_reported
WHERE id = ?
  AND deleted_at IS NULL
`,
		reportID,
	).Scan(&state).Error
	if err != nil {
		t.Fatalf(
			"load generated report state: %v",
			err,
		)
	}

	return state
}

func countFileEffects(
	t *testing.T,
	db *gorm.DB,
	jobID string,
) int64 {
	t.Helper()

	var count int64

	err := db.Raw(
		`
SELECT count(*)
FROM file
WHERE owner_namespace = ?
  AND owner_key = ?
`,
		processing.FileOwnerNamespace,
		jobID,
	).Scan(&count).Error
	if err != nil {
		t.Fatalf(
			"count File effects: %v",
			err,
		)
	}

	return count
}

func loadFileEffect(
	t *testing.T,
	db *gorm.DB,
	jobID string,
) durableFileEffect {
	t.Helper()

	var rows []durableFileEffect

	err := db.Raw(
		`
SELECT
    id,
    owner_namespace,
    owner_key,
    owner_request_hash,
    path,
    hash,
    size
FROM file
WHERE owner_namespace = ?
  AND owner_key = ?
ORDER BY id
`,
		processing.FileOwnerNamespace,
		jobID,
	).Scan(&rows).Error
	if err != nil {
		t.Fatalf(
			"load File effect: %v",
			err,
		)
	}

	if len(rows) != 1 {
		t.Fatalf(
			"File effects=%d, want 1",
			len(rows),
		)
	}

	return rows[0]
}

func requireFileHTTP(
	t *testing.T,
	baseURL string,
) {
	t.Helper()

	req, err := http.NewRequest(
		http.MethodPost,
		strings.TrimRight(baseURL, "/")+
			"/internal/v1/file/owned",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{
		Timeout: time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf(
			"File HTTP runtime is required for P7: %v",
			err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode !=
		http.StatusUnauthorized {
		t.Fatalf(
			"File HTTP preflight status=%d, want 401",
			resp.StatusCode,
		)
	}
}

func TestReportLabPostgresJobRecoversAfterProcessDeath(
	t *testing.T,
) {
	if os.Getenv(
		"QHPRO_REPORT_LAB_INTEGRATION",
	) != "1" {
		t.Skip(
			"report-lab integration test",
		)
	}

	tqdDSN :=
		strings.TrimSpace(
			os.Getenv(
				"QHPRO_TEST_TQD_DSN",
			),
		)

	fileDSN :=
		strings.TrimSpace(
			os.Getenv(
				"QHPRO_TEST_FILE_DSN",
			),
		)

	fileBaseURL :=
		strings.TrimRight(
			strings.TrimSpace(
				os.Getenv(
					"QHPRO_FILE_HTTP_BASE_URL",
				),
			),
			"/",
		)

	serviceAuthKey :=
		strings.TrimSpace(
			os.Getenv(
				"QHPRO_FILE_HTTP_SERVICE_AUTH_KEY",
			),
		)

	if tqdDSN == "" {
		t.Fatal(
			"QHPRO_TEST_TQD_DSN is required",
		)
	}
	if fileDSN == "" {
		t.Fatal(
			"QHPRO_TEST_FILE_DSN is required",
		)
	}
	if fileBaseURL == "" {
		fileBaseURL =
			"http://127.0.0.1:8002"
	}
	if serviceAuthKey == "" {
		t.Fatal(
			"QHPRO_FILE_HTTP_SERVICE_AUTH_KEY is required",
		)
	}

	requireFileHTTP(
		t,
		fileBaseURL,
	)

	// --------------------------------------------------------
	// Two independent DB pools represent two service processes.
	// --------------------------------------------------------

	dbA :=
		openReportLabPostgres(
			t,
			tqdDSN,
		)

	dbB :=
		openReportLabPostgres(
			t,
			tqdDSN,
		)

	fileDB :=
		openReportLabPostgres(
			t,
			fileDSN,
		)

	sqlB, err := dbB.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlB.Close()

	sqlFile, err := fileDB.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlFile.Close()

	storeA :=
		postgresjob.NewStore(
			dbA,
		)

	storeB :=
		postgresjob.NewStore(
			dbB,
		)

	now :=
		time.Now().UTC()

	userID :=
		uint64(42)

	commandKey :=
		"idem_p7_restart_" +
			strconv.FormatInt(
				now.UnixNano(),
				10,
			)

	jobID :=
		processing.BuildID(
			userID,
			commandKey,
		)

	reportID :=
		insertProcessingReport(
			t,
			dbA,
			userID,
			commandKey,
			jobID,
		)

	// Use an old deterministic queue time so this Job wins queue
	// ordering even if unrelated lab rows exist.
	claimAt :=
		time.Date(
			2000,
			time.January,
			1,
			0,
			0,
			0,
			0,
			time.UTC,
		)

	job := processing.Job{
		ID: jobID,

		ReportID: reportID,

		UserID: userID,

		Operation: operation.Code(
			"workspace.report.generate",
		),

		OperationID: "op_p7_restart_" +
			strconv.FormatInt(
				now.UnixNano(),
				10,
			),

		CommandKey: commandKey,

		Status: processing.StatusPending,

		AvailableAt: claimAt.Add(
			-time.Second,
		),

		CreatedAt: now,

		UpdatedAt: now,
	}

	created, err :=
		storeA.SaveJob(
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
			"P7 durable Job was not created",
		)
	}

	if got :=
		countFileEffects(
			t,
			fileDB,
			jobID,
		); got != 0 {
		t.Fatalf(
			"File effects before execution=%d, want 0",
			got,
		)
	}

	// ========================================================
	// PROCESS A CLAIMS
	// ========================================================

	claimA, found, err :=
		storeA.ClaimNextJob(
			context.Background(),
			"p7-worker-a",
			claimAt,
		)
	if err != nil {
		t.Fatalf(
			"process A ClaimNextJob(): %v",
			err,
		)
	}
	if !found {
		t.Fatal(
			"process A did not claim Job",
		)
	}

	if claimA.ID != jobID {
		t.Fatalf(
			"process A claimed unexpected Job=%s, want %s",
			claimA.ID,
			jobID,
		)
	}

	if claimA.Status !=
		processing.StatusRunning ||
		claimA.Attempts != 1 ||
		claimA.ClaimVersion != 1 ||
		claimA.LockedBy !=
			"p7-worker-a" {
		t.Fatalf(
			"process A claim=%+v",
			claimA,
		)
	}

	t.Logf(
		"process A claimed durable Job job=%s report=%d claim=%d",
		jobID,
		reportID,
		claimA.ClaimVersion,
	)

	// ========================================================
	// PROCESS A DIES
	//
	// No RetryJob.
	// No CompleteJob.
	// No FailJob.
	//
	// Close the entire DB pool to make the process boundary real
	// from PostgreSQL's perspective.
	// ========================================================

	sqlA, err :=
		dbA.DB()
	if err != nil {
		t.Fatal(err)
	}

	if err :=
		sqlA.Close(); err != nil {
		t.Fatalf(
			"close process A DB pool: %v",
			err,
		)
	}

	// Process B is a separately-created Store backed by another
	// PostgreSQL connection pool.
	beforeRecovery, found, err :=
		storeB.FindJobByReportID(
			context.Background(),
			reportID,
		)
	if err != nil {
		t.Fatalf(
			"process B FindJobByReportID(): %v",
			err,
		)
	}
	if !found {
		t.Fatal(
			"process B cannot see durable Job after A died",
		)
	}

	if beforeRecovery.Status !=
		processing.StatusRunning ||
		beforeRecovery.ClaimVersion != 1 ||
		beforeRecovery.LockedBy !=
			"p7-worker-a" {
		t.Fatalf(
			"durable Job after process A death=%+v",
			beforeRecovery,
		)
	}

	// ========================================================
	// STALE RECOVERY BY PROCESS B
	// ========================================================

	lockedBefore :=
		claimAt.Add(
			time.Minute,
		)

	released, err :=
		storeB.ReleaseStaleJobs(
			context.Background(),
			lockedBefore,
			100,
		)
	if err != nil {
		t.Fatalf(
			"process B ReleaseStaleJobs(): %v",
			err,
		)
	}

	if released != 1 {
		t.Fatalf(
			"released stale Jobs=%d, want 1",
			released,
		)
	}

	recovered, found, err :=
		storeB.FindJobByReportID(
			context.Background(),
			reportID,
		)
	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Fatal(
			"recovered Job disappeared",
		)
	}

	if recovered.Status !=
		processing.StatusPending ||
		recovered.LockedBy != "" ||
		recovered.LockedAt != nil ||
		recovered.ClaimVersion != 1 {
		t.Fatalf(
			"recovered Job=%+v",
			recovered,
		)
	}

	t.Logf(
		"process B recovered stale Job job=%s claim=%d status=%s",
		jobID,
		recovered.ClaimVersion,
		recovered.Status,
	)

	// ========================================================
	// PROCESS B RESTART EXECUTION
	// ========================================================

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

	reportGenerator, err :=
		generator.NewSmokeGenerator(
			uploader,
			uploader,
		)
	if err != nil {
		t.Fatalf(
			"create Report generator: %v",
			err,
		)
	}

	workerB :=
		processing.NewWorker(
			storeB,
			reportGenerator,
			slog.New(
				slog.NewTextHandler(
					io.Discard,
					nil,
				),
			),
			"p7-worker-b",
		)

	processed, err :=
		workerB.RunOnce(
			context.Background(),
		)
	if err != nil {
		t.Fatalf(
			"process B Worker.RunOnce(): %v",
			err,
		)
	}
	if !processed {
		t.Fatal(
			"process B did not reclaim recovered Job",
		)
	}

	// ========================================================
	// DURABLE POSTCONDITIONS
	// ========================================================

	finalJob, found, err :=
		storeB.FindJobByReportID(
			context.Background(),
			reportID,
		)
	if err != nil {
		t.Fatalf(
			"load final durable Job: %v",
			err,
		)
	}
	if !found {
		t.Fatal(
			"final Job disappeared",
		)
	}

	if finalJob.Status !=
		processing.StatusCompleted {
		t.Fatalf(
			"final Job status=%s, want completed",
			finalJob.Status,
		)
	}

	if finalJob.Attempts != 2 {
		t.Fatalf(
			"final Job attempts=%d, want 2",
			finalJob.Attempts,
		)
	}

	if finalJob.ClaimVersion != 2 {
		t.Fatalf(
			"final Job claim_version=%d, want 2",
			finalJob.ClaimVersion,
		)
	}

	if finalJob.LockedAt != nil ||
		finalJob.LockedBy != "" {
		t.Fatalf(
			"completed Job retained lock: %+v",
			finalJob,
		)
	}

	reportState :=
		loadReportState(
			t,
			dbB,
			reportID,
		)

	if reportState.Status !=
		int(
			enums.GeneratedReportStatusReady,
		) {
		t.Fatalf(
			"Report status=%d, want ready=%d",
			reportState.Status,
			enums.GeneratedReportStatusReady,
		)
	}

	if strings.TrimSpace(
		reportState.PDFURL,
	) == "" {
		t.Fatal(
			"ready Report has empty pdf_url",
		)
	}

	if reportState.FileSize == 0 {
		t.Fatal(
			"ready Report has zero file_size",
		)
	}

	if reportState.ErrorMessage != "" {
		t.Fatalf(
			"ready Report retained error=%q",
			reportState.ErrorMessage,
		)
	}

	if got :=
		countFileEffects(
			t,
			fileDB,
			jobID,
		); got != 1 {
		t.Fatalf(
			"File effects=%d, want 1",
			got,
		)
	}

	fileEffect :=
		loadFileEffect(
			t,
			fileDB,
			jobID,
		)

	if fileEffect.ID == 0 ||
		fileEffect.OwnerRequestHash == "" ||
		fileEffect.Path == "" ||
		fileEffect.Hash == "" ||
		fileEffect.Size <= 0 {
		t.Fatalf(
			"incomplete durable File effect=%+v",
			fileEffect,
		)
	}

	t.Logf(
		"PASS PostgreSQL restart recovery job=%s report=%d attempts=%d claim_version=%d file_id=%d",
		jobID,
		reportID,
		finalJob.Attempts,
		finalJob.ClaimVersion,
		fileEffect.ID,
	)
}
