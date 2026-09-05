package postgresjob_test

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"common/operation"

	reportfilehttp "tqd/infra/client/file/generatedreport"
	"tqd/infra/postgres/generatedreport/job"
	"tqd/infra/worker/generatedreport/generator"
	"tqd/internal/usecase/generatedreport/processing"
)

func TestReportLabPostgresRejectsStaleClaimMutations(
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

	// Independent stores model two simultaneously alive processes.
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

	sqlA, err :=
		dbA.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlA.Close()

	sqlB, err :=
		dbB.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlB.Close()

	sqlFile, err :=
		fileDB.DB()
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
		"idem_p8_fence_" +
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

	// Keep this fixture at the front of the lab queue without
	// touching unrelated durable rows.
	claimAt :=
		time.Date(
			1900,
			time.January,
			1,
			0,
			0,
			0,
			0,
			time.UTC,
		)

	job :=
		processing.Job{
			ID: jobID,

			ReportID: reportID,

			UserID: userID,

			Operation: operation.Code(
				"workspace.report.generate",
			),

			OperationID: "op_p8_fence_" +
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
			"P8 durable Job was not created",
		)
	}

	// ========================================================
	// CLAIM A — VERSION 1
	// ========================================================

	claimA, found, err :=
		storeA.ClaimNextJob(
			context.Background(),
			"p8-worker-a",
			claimAt,
		)
	if err != nil {
		t.Fatalf(
			"worker A ClaimNextJob(): %v",
			err,
		)
	}
	if !found {
		t.Fatal(
			"worker A did not claim Job",
		)
	}

	if claimA.ID != jobID {
		t.Fatalf(
			"worker A claimed unexpected Job=%s, want %s",
			claimA.ID,
			jobID,
		)
	}

	if claimA.Status !=
		processing.StatusRunning ||
		claimA.Attempts != 1 ||
		claimA.ClaimVersion != 1 ||
		claimA.LockedBy !=
			"p8-worker-a" {
		t.Fatalf(
			"worker A claim=%+v",
			claimA,
		)
	}

	t.Logf(
		"worker A owns stale snapshot job=%s claim=%d",
		jobID,
		claimA.ClaimVersion,
	)

	// ========================================================
	// A BECOMES STALE
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
			"ReleaseStaleJobs(): %v",
			err,
		)
	}

	if released != 1 {
		t.Fatalf(
			"released Jobs=%d, want 1",
			released,
		)
	}

	// ========================================================
	// CLAIM B — VERSION 2
	// ========================================================

	claimBAt :=
		claimAt.Add(
			2 * time.Minute,
		)

	claimB, found, err :=
		storeB.ClaimNextJob(
			context.Background(),
			"p8-worker-b",
			claimBAt,
		)
	if err != nil {
		t.Fatalf(
			"worker B ClaimNextJob(): %v",
			err,
		)
	}
	if !found {
		t.Fatal(
			"worker B did not reclaim Job",
		)
	}

	if claimB.ID != jobID {
		t.Fatalf(
			"worker B claimed unexpected Job=%s, want %s",
			claimB.ID,
			jobID,
		)
	}

	if claimB.Status !=
		processing.StatusRunning ||
		claimB.Attempts != 2 ||
		claimB.ClaimVersion != 2 ||
		claimB.LockedBy !=
			"p8-worker-b" {
		t.Fatalf(
			"worker B claim=%+v",
			claimB,
		)
	}

	t.Logf(
		"worker B owns current claim job=%s claim=%d",
		jobID,
		claimB.ClaimVersion,
	)

	// ========================================================
	// B PERFORMS REAL EXTERNAL EFFECT
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

	winningOutput, err :=
		reportGenerator.GenerateReport(
			context.Background(),
			claimB,
		)
	if err != nil {
		t.Fatalf(
			"worker B GenerateReport(): %v",
			err,
		)
	}

	if strings.TrimSpace(
		winningOutput.PDFURL,
	) == "" {
		t.Fatal(
			"worker B generated empty PDF URL",
		)
	}

	if winningOutput.FileSize == 0 {
		t.Fatal(
			"worker B generated zero-size output",
		)
	}

	if got :=
		countFileEffects(
			t,
			fileDB,
			jobID,
		); got != 1 {
		t.Fatalf(
			"File effects before B finalization=%d, want 1",
			got,
		)
	}

	winningFile :=
		loadFileEffect(
			t,
			fileDB,
			jobID,
		)

	// ========================================================
	// CURRENT CLAIM B FINALIZES
	// ========================================================

	completeAt :=
		claimAt.Add(
			3 * time.Minute,
		)

	if err :=
		storeB.CompleteJob(
			context.Background(),
			claimB,
			winningOutput,
			completeAt,
		); err != nil {
		t.Fatalf(
			"worker B CompleteJob(): %v",
			err,
		)
	}

	finalBeforeAttack, found, err :=
		storeB.FindJobByReportID(
			context.Background(),
			reportID,
		)
	if err != nil {
		t.Fatalf(
			"load completed Job: %v",
			err,
		)
	}
	if !found {
		t.Fatal(
			"completed Job disappeared",
		)
	}

	if finalBeforeAttack.Status !=
		processing.StatusCompleted ||
		finalBeforeAttack.Attempts != 2 ||
		finalBeforeAttack.ClaimVersion != 2 ||
		finalBeforeAttack.LockedBy != "" ||
		finalBeforeAttack.LockedAt != nil ||
		finalBeforeAttack.LastError != "" {
		t.Fatalf(
			"completed Job=%+v",
			finalBeforeAttack,
		)
	}

	reportBeforeAttack :=
		loadReportState(
			t,
			dbB,
			reportID,
		)

	if reportBeforeAttack.PDFURL !=
		winningOutput.PDFURL {
		t.Fatalf(
			"Report PDF URL=%q, want winner=%q",
			reportBeforeAttack.PDFURL,
			winningOutput.PDFURL,
		)
	}

	if reportBeforeAttack.FileSize !=
		winningOutput.FileSize {
		t.Fatalf(
			"Report file size=%d, want winner=%d",
			reportBeforeAttack.FileSize,
			winningOutput.FileSize,
		)
	}

	if reportBeforeAttack.ErrorMessage != "" {
		t.Fatalf(
			"ready Report error=%q",
			reportBeforeAttack.ErrorMessage,
		)
	}

	// ========================================================
	// STALE CLAIM A ATTACKS DURABLE AUTHORITY
	//
	// All three mutation paths must reject the stale fencing
	// token. CompleteJob must also be unable to overwrite the
	// winner's Report output.
	// ========================================================

	staleOutput :=
		processing.Output{
			ThumbnailURL: "stale://thumbnail",

			ImageURL: "stale://image",

			PDFURL: "stale://pdf-must-not-persist",

			ShareURL: "stale://share",

			FileSize: 999999,
		}

	err =
		storeA.CompleteJob(
			context.Background(),
			claimA,
			staleOutput,
			completeAt.Add(
				time.Minute,
			),
		)
	if !errors.Is(
		err,
		processing.ErrClaimLost,
	) {
		t.Fatalf(
			"stale CompleteJob error=%v, want ErrClaimLost",
			err,
		)
	}

	t.Log(
		"stale claim CompleteJob rejected",
	)

	err =
		storeA.RetryJob(
			context.Background(),
			claimA,
			"stale retry must not persist",
			completeAt.Add(
				2*time.Minute,
			),
		)
	if !errors.Is(
		err,
		processing.ErrClaimLost,
	) {
		t.Fatalf(
			"stale RetryJob error=%v, want ErrClaimLost",
			err,
		)
	}

	t.Log(
		"stale claim RetryJob rejected",
	)

	err =
		storeA.FailJob(
			context.Background(),
			claimA,
			"stale failure must not persist",
			completeAt.Add(
				3*time.Minute,
			),
		)
	if !errors.Is(
		err,
		processing.ErrClaimLost,
	) {
		t.Fatalf(
			"stale FailJob error=%v, want ErrClaimLost",
			err,
		)
	}

	t.Log(
		"stale claim FailJob rejected",
	)

	// ========================================================
	// PROVE WINNER DURABLE STATE WAS NOT MUTATED
	// ========================================================

	finalAfterAttack, found, err :=
		storeB.FindJobByReportID(
			context.Background(),
			reportID,
		)
	if err != nil {
		t.Fatalf(
			"load Job after stale mutations: %v",
			err,
		)
	}
	if !found {
		t.Fatal(
			"Job disappeared after stale mutations",
		)
	}

	if finalAfterAttack.Status !=
		processing.StatusCompleted {
		t.Fatalf(
			"stale claim changed Job status=%s",
			finalAfterAttack.Status,
		)
	}

	if finalAfterAttack.Attempts !=
		finalBeforeAttack.Attempts {
		t.Fatalf(
			"stale claim changed attempts: before=%d after=%d",
			finalBeforeAttack.Attempts,
			finalAfterAttack.Attempts,
		)
	}

	if finalAfterAttack.ClaimVersion !=
		finalBeforeAttack.ClaimVersion {
		t.Fatalf(
			"stale claim changed claim_version: before=%d after=%d",
			finalBeforeAttack.ClaimVersion,
			finalAfterAttack.ClaimVersion,
		)
	}

	if finalAfterAttack.LockedBy != "" ||
		finalAfterAttack.LockedAt != nil {
		t.Fatalf(
			"stale claim restored Job lock=%+v",
			finalAfterAttack,
		)
	}

	if finalAfterAttack.LastError != "" {
		t.Fatalf(
			"stale claim changed LastError=%q",
			finalAfterAttack.LastError,
		)
	}

	reportAfterAttack :=
		loadReportState(
			t,
			dbB,
			reportID,
		)

	if reportAfterAttack.Status !=
		reportBeforeAttack.Status {
		t.Fatalf(
			"stale claim changed Report status: before=%d after=%d",
			reportBeforeAttack.Status,
			reportAfterAttack.Status,
		)
	}

	if reportAfterAttack.PDFURL !=
		reportBeforeAttack.PDFURL {
		t.Fatalf(
			"stale claim changed Report PDF: before=%q after=%q",
			reportBeforeAttack.PDFURL,
			reportAfterAttack.PDFURL,
		)
	}

	if reportAfterAttack.FileSize !=
		reportBeforeAttack.FileSize {
		t.Fatalf(
			"stale claim changed Report file size: before=%d after=%d",
			reportBeforeAttack.FileSize,
			reportAfterAttack.FileSize,
		)
	}

	if reportAfterAttack.ErrorMessage !=
		reportBeforeAttack.ErrorMessage {
		t.Fatalf(
			"stale claim changed Report error: before=%q after=%q",
			reportBeforeAttack.ErrorMessage,
			reportAfterAttack.ErrorMessage,
		)
	}

	if reportAfterAttack.PDFURL ==
		staleOutput.PDFURL {
		t.Fatal(
			"stale CompleteJob overwrote winner PDF",
		)
	}

	if got :=
		countFileEffects(
			t,
			fileDB,
			jobID,
		); got != 1 {
		t.Fatalf(
			"File effects after stale mutations=%d, want 1",
			got,
		)
	}

	finalFile :=
		loadFileEffect(
			t,
			fileDB,
			jobID,
		)

	if finalFile.ID !=
		winningFile.ID ||
		finalFile.Path !=
			winningFile.Path ||
		finalFile.Hash !=
			winningFile.Hash {
		t.Fatalf(
			"File effect changed after stale mutations: before=%+v after=%+v",
			winningFile,
			finalFile,
		)
	}

	t.Logf(
		"PASS PostgreSQL fencing job=%s report=%d stale_claim=%d winner_claim=%d file_id=%d",
		jobID,
		reportID,
		claimA.ClaimVersion,
		finalAfterAttack.ClaimVersion,
		finalFile.ID,
	)
}
