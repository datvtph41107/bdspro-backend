package application

import (
	"common/identity"
	commonoperation "common/operation"
	"common/request"
	"context"
	"errors"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"
	"tqd/internal/access"
	"tqd/internal/usecase/generatedreport/processing"
	"tqd/internal/usecase/quota"
	"tqd/internal/usecase/quota/memory"
	"tqd/internal/usecase/usage"
)

type fakeAccessClient struct {
	result access.Result
	err    error
}

func (f fakeAccessClient) GetAccess(context.Context, commonoperation.Code) (access.Result, error) {
	return f.result, f.err
}

type fakeSourceStore struct {
	parcel Parcel
	region Region
}

func (f fakeSourceStore) FindParcel(context.Context, uint64) (Parcel, bool, error) {
	if f.parcel.ID == 0 {
		return Parcel{}, false, nil
	}
	return f.parcel, true, nil
}

func (f fakeSourceStore) FindRegion(context.Context, uint64) (Region, bool, error) {
	if f.region.ID == 0 {
		return Region{}, false, nil
	}
	return f.region, true, nil
}

type fakeSaveStore struct {
	mu sync.Mutex

	fail error

	nextID          uint64
	byKey           map[string]Report
	usage           map[string]usage.Event
	jobs            map[string]processing.Job
	forceFindMisses int
}

func newFakeSaveStore() *fakeSaveStore {
	return &fakeSaveStore{
		nextID: 1,
		byKey:  make(map[string]Report),
		usage:  make(map[string]usage.Event),
		jobs:   make(map[string]processing.Job),
	}
}

func (s *fakeSaveStore) FindReportByCommand(
	ctx context.Context,
	userID uint64,
	commandKey string,
) (Report, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.forceFindMisses > 0 {
		s.forceFindMisses--
		return Report{}, false, nil
	}
	report, ok := s.byKey[commandKey]
	return report, ok, nil
}

type failingQuotaStore struct {
	quota.Store
	commitErr error
	cancelErr error
}

func (s *failingQuotaStore) CommitQuota(
	ctx context.Context,
	reservationID string,
) (quota.Reservation, error) {
	if s.commitErr != nil {
		return quota.Reservation{}, s.commitErr
	}
	return s.Store.CommitQuota(ctx, reservationID)
}

func (s *failingQuotaStore) CancelQuota(
	ctx context.Context,
	reservationID string,
) (quota.Reservation, error) {
	if s.cancelErr != nil {
		return quota.Reservation{}, s.cancelErr
	}
	return s.Store.CancelQuota(ctx, reservationID)
}

func (s *fakeSaveStore) AcceptReport(
	ctx context.Context,
	input Acceptance,
) (AcceptanceResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.fail != nil {
		return AcceptanceResult{}, s.fail
	}

	report := input.Report
	if existing, ok := s.byKey[report.CommandKey]; ok {
		var durableUsage *usage.Event
		if input.Usage != nil {
			if savedUsage, found := s.usage[input.Usage.UsageKey]; found {
				usage := savedUsage
				durableUsage = &usage
			}
		}
		return AcceptanceResult{
			Report:       existing,
			Created:      false,
			DurableUsage: durableUsage,
		}, nil
	}

	report.ID = s.nextID
	s.nextID++
	report.CreatedAt = input.Job.CreatedAt
	report.UpdatedAt = input.Job.UpdatedAt
	s.byKey[report.CommandKey] = report

	var durableUsage *usage.Event
	if input.Usage != nil {
		usage := *input.Usage
		s.usage[usage.UsageKey] = usage
		durableUsage = &usage
	}

	job := input.Job
	job.ReportID = report.ID
	s.jobs[job.ID] = job

	return AcceptanceResult{
		Report:       report,
		Created:      true,
		DurableUsage: durableUsage,
	}, nil
}

func TestCreateGeneratedReportAcceptsAndChargesOnce(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 11, 5, 0, 0, 0, time.UTC)
	periodStart := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	subject := access.Subject{Type: access.SubjectProfile, ID: "42"}

	quotaStore := memoryquota.NewStore()
	quotaService := quota.NewService(quotaStore)
	saver := newFakeSaveStore()
	service := NewService(
		fakeSourceStore{parcel: Parcel{
			ID:         99,
			MapNumber:  "12",
			LandNumber: "44",
			AreaSqm:    125.5,
			Location:   Location{Address: "Test address"},
			Spatial:    Spatial{Centroid: Point{Lat: 10, Lon: 106}},
		}},
		saver,
		fakeAccessClient{result: access.Result{
			Subject:     subject,
			Operation:   commonoperation.Code("workspace.report.generate"),
			Allowed:     true,
			Metering:    access.Metering{FeatureCode: "workspace.report.generate", MeterCode: "workspace.report_generation.accepted", UnitsPerAction: 1, PolicyVersion: "1.0.0"},
			Limit:       100,
			Period:      access.PeriodSubscriptionCycle,
			PeriodStart: periodStart,
			PeriodEnd:   periodEnd,
		}},
		quotaService,
		true,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	service.now = func() time.Time { return now }

	ctx := testContext(t, "op_report_500", "idem_report_700")
	first, err := service.CreateGeneratedReport(ctx, Input{
		ReportType: 1,
		EntityType: 1,
		ParcelID:   99,
	})
	if err != nil {
		t.Fatalf("first CreateGeneratedReport() error = %v", err)
	}
	if first.Outcome != OutcomeAccepted || first.Report.ID == 0 {
		t.Fatalf("first result = %+v", first)
	}

	second, err := service.CreateGeneratedReport(ctx, Input{
		ReportType: 1,
		EntityType: 1,
		ParcelID:   99,
	})
	if err != nil {
		t.Fatalf("retry CreateGeneratedReport() error = %v", err)
	}
	if second.Outcome != OutcomeExisting || second.Report.ID != first.Report.ID {
		t.Fatalf("retry result = %+v, first = %+v", second, first)
	}

	current, err := quotaStore.GetUsage(
		context.Background(),
		subject,
		"workspace.report_generation.accepted",
		periodStart,
		periodEnd,
	)
	if err != nil {
		t.Fatalf("GetUsage() error = %v", err)
	}
	if current.Used != 1 || current.Reserved != 0 {
		t.Fatalf("usage = %+v, want used=1 reserved=0", current)
	}
	if len(saver.usage) != 1 {
		t.Fatalf("usage events = %d, want 1", len(saver.usage))
	}
	if len(saver.jobs) != 1 {
		t.Fatalf("report jobs = %d, want 1", len(saver.jobs))
	}
	job := saver.jobs[first.Report.JobID]
	if job.ReportID != first.Report.ID || job.UserID != 42 ||
		job.CommandKey != "idem_report_700" || job.Status != processing.StatusPending {
		t.Fatalf("accepted job = %+v, report = %+v", job, first.Report)
	}
	if !job.AvailableAt.Equal(now) || !job.CreatedAt.Equal(now) || !job.UpdatedAt.Equal(now) {
		t.Fatalf("accepted job timestamps = %+v, want %v", job, now)
	}
}

func TestCreateGeneratedReportKeepsAcceptedOutcomeWhenQuotaCommitFails(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 11, 5, 0, 0, 0, time.UTC)
	periodStart := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	subject := access.Subject{Type: access.SubjectProfile, ID: "42"}

	baseQuotaStore := memoryquota.NewStore()
	quotaStore := &failingQuotaStore{
		Store:     baseQuotaStore,
		commitErr: errors.New("redis unavailable"),
	}
	quotaService := quota.NewService(quotaStore)
	saver := newFakeSaveStore()
	service := NewService(
		fakeSourceStore{parcel: Parcel{ID: 99, LandNumber: "44"}},
		saver,
		fakeAccessClient{result: access.Result{
			Subject:     subject,
			Operation:   commonoperation.Code("workspace.report.generate"),
			Allowed:     true,
			Metering:    access.Metering{FeatureCode: "workspace.report.generate", MeterCode: "workspace.report_generation.accepted", UnitsPerAction: 1, PolicyVersion: "1.0.0"},
			Limit:       100,
			Period:      access.PeriodSubscriptionCycle,
			PeriodStart: periodStart,
			PeriodEnd:   periodEnd,
		}},
		quotaService,
		true,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	service.now = func() time.Time { return now }

	ctx := testContext(t, "op_report_502", "idem_report_702")
	first, err := service.CreateGeneratedReport(
		ctx,
		Input{ReportType: 1, EntityType: 1, ParcelID: 99},
	)
	if err != nil {
		t.Fatalf("CreateGeneratedReport() error = %v", err)
	}
	if first.Outcome != OutcomeAccepted || first.Report.ID == 0 {
		t.Fatalf("result = %+v, want accepted durable report", first)
	}
	if len(saver.usage) != 1 {
		t.Fatalf("usage events = %d, want 1 durable event", len(saver.usage))
	}

	current, err := baseQuotaStore.GetUsage(
		context.Background(),
		subject,
		"workspace.report_generation.accepted",
		periodStart,
		periodEnd,
	)
	if err != nil {
		t.Fatalf("GetUsage() error = %v", err)
	}
	if current.Used != 0 || current.Reserved != 1 {
		t.Fatalf("usage = %+v, want used=0 reserved=1 pending repair", current)
	}

	replay, err := service.CreateGeneratedReport(
		ctx,
		Input{ReportType: 1, EntityType: 1, ParcelID: 99},
	)
	if err != nil {
		t.Fatalf("replay CreateGeneratedReport() error = %v", err)
	}
	if replay.Outcome != OutcomeExisting || replay.Report.ID != first.Report.ID {
		t.Fatalf("replay = %+v, first = %+v", replay, first)
	}

	current, err = baseQuotaStore.GetUsage(
		context.Background(),
		subject,
		"workspace.report_generation.accepted",
		periodStart,
		periodEnd,
	)
	if err != nil {
		t.Fatalf("GetUsage() after replay error = %v", err)
	}
	if current.Used != 0 || current.Reserved != 1 {
		t.Fatalf("usage after replay = %+v, want same pending reservation", current)
	}
}

func TestCreateGeneratedReportKeepsExistingOutcomeWhenReservationCancelFails(t *testing.T) {
	t.Parallel()

	firstStart := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	firstEnd := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	secondStart := firstEnd
	secondEnd := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	subject := access.Subject{Type: access.SubjectProfile, ID: "42"}

	baseQuotaStore := memoryquota.NewStore()
	quotaStore := &failingQuotaStore{Store: baseQuotaStore}
	quotaService := quota.NewService(quotaStore)
	saver := newFakeSaveStore()
	source := fakeSourceStore{parcel: Parcel{ID: 99, LandNumber: "44"}}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	firstService := NewService(
		source,
		saver,
		fakeAccessClient{result: access.Result{
			Subject:     subject,
			Operation:   commonoperation.Code("workspace.report.generate"),
			Allowed:     true,
			Metering:    access.Metering{FeatureCode: "workspace.report.generate", MeterCode: "workspace.report_generation.accepted", UnitsPerAction: 1, PolicyVersion: "1.0.0"},
			Limit:       100,
			Period:      access.PeriodSubscriptionCycle,
			PeriodStart: firstStart,
			PeriodEnd:   firstEnd,
		}},
		quotaService,
		true,
		logger,
	)

	ctx := testContext(t, "op_report_562", "idem_report_762")
	first, err := firstService.CreateGeneratedReport(
		ctx,
		Input{ReportType: 1, EntityType: 1, ParcelID: 99},
	)
	if err != nil {
		t.Fatalf("first CreateGeneratedReport() error = %v", err)
	}

	quotaStore.cancelErr = errors.New("redis unavailable")
	secondService := NewService(
		source,
		saver,
		fakeAccessClient{result: access.Result{
			Subject:     subject,
			Operation:   commonoperation.Code("workspace.report.generate"),
			Allowed:     true,
			Metering:    access.Metering{FeatureCode: "workspace.report.generate", MeterCode: "workspace.report_generation.accepted", UnitsPerAction: 1, PolicyVersion: "1.0.0"},
			Limit:       100,
			Period:      access.PeriodSubscriptionCycle,
			PeriodStart: secondStart,
			PeriodEnd:   secondEnd,
		}},
		quotaService,
		true,
		logger,
	)

	// Simulate both requests missing the durable precheck before the unique
	// report constraint decides which transaction won.
	saver.forceFindMisses = 1
	second, err := secondService.CreateGeneratedReport(
		ctx,
		Input{ReportType: 1, EntityType: 1, ParcelID: 99},
	)
	if err != nil {
		t.Fatalf("race CreateGeneratedReport() error = %v", err)
	}
	if second.Outcome != OutcomeExisting || second.Report.ID != first.Report.ID {
		t.Fatalf("race result = %+v, first = %+v", second, first)
	}

	secondUsage, err := baseQuotaStore.GetUsage(
		context.Background(),
		subject,
		"workspace.report_generation.accepted",
		secondStart,
		secondEnd,
	)
	if err != nil {
		t.Fatalf("GetUsage() error = %v", err)
	}
	if secondUsage.Used != 0 || secondUsage.Reserved != 1 {
		t.Fatalf("second usage = %+v, want redundant reservation pending repair", secondUsage)
	}
}

func TestCreateGeneratedReportDoesNotChargeAgainWhenRaceCrossesPeriod(t *testing.T) {
	t.Parallel()

	firstStart := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	firstEnd := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	secondStart := firstEnd
	secondEnd := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	subject := access.Subject{Type: access.SubjectProfile, ID: "42"}
	quotaStore := memoryquota.NewStore()
	quotaService := quota.NewService(quotaStore)
	saver := newFakeSaveStore()
	source := fakeSourceStore{parcel: Parcel{ID: 99, LandNumber: "44"}}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	firstService := NewService(
		source,
		saver,
		fakeAccessClient{result: access.Result{
			Subject:     subject,
			Operation:   commonoperation.Code("workspace.report.generate"),
			Allowed:     true,
			Metering:    access.Metering{FeatureCode: "workspace.report.generate", MeterCode: "workspace.report_generation.accepted", UnitsPerAction: 1, PolicyVersion: "1.0.0"},
			Limit:       100,
			Period:      access.PeriodSubscriptionCycle,
			PeriodStart: firstStart,
			PeriodEnd:   firstEnd,
		}},
		quotaService,
		true,
		logger,
	)

	ctx := testContext(t, "op_report_560", "idem_report_760")
	first, err := firstService.CreateGeneratedReport(ctx, Input{ReportType: 1, EntityType: 1, ParcelID: 99})
	if err != nil {
		t.Fatalf("first CreateGeneratedReport() error = %v", err)
	}

	secondService := NewService(
		source,
		saver,
		fakeAccessClient{result: access.Result{
			Subject:     subject,
			Operation:   commonoperation.Code("workspace.report.generate"),
			Allowed:     true,
			Metering:    access.Metering{FeatureCode: "workspace.report.generate", MeterCode: "workspace.report_generation.accepted", UnitsPerAction: 1, PolicyVersion: "1.0.0"},
			Limit:       100,
			Period:      access.PeriodSubscriptionCycle,
			PeriodStart: secondStart,
			PeriodEnd:   secondEnd,
		}},
		quotaService,
		true,
		logger,
	)
	// Simulate both requests missing the durable precheck before the unique
	// report constraint decides which transaction won.
	saver.forceFindMisses = 1

	second, err := secondService.CreateGeneratedReport(ctx, Input{ReportType: 1, EntityType: 1, ParcelID: 99})
	if err != nil {
		t.Fatalf("retry in next period error = %v", err)
	}
	if second.Outcome != OutcomeExisting || second.Report.ID != first.Report.ID {
		t.Fatalf("retry result = %+v, first = %+v", second, first)
	}

	firstUsage, err := quotaStore.GetUsage(context.Background(), subject, "workspace.report_generation.accepted", firstStart, firstEnd)
	if err != nil {
		t.Fatal(err)
	}
	secondUsage, err := quotaStore.GetUsage(context.Background(), subject, "workspace.report_generation.accepted", secondStart, secondEnd)
	if err != nil {
		t.Fatal(err)
	}
	if firstUsage.Used != 1 || secondUsage.Used != 0 || secondUsage.Reserved != 0 {
		t.Fatalf("first usage=%+v second usage=%+v, retry must not charge next period", firstUsage, secondUsage)
	}
}

func TestCreateGeneratedReportRejectsReusedKeyWithDifferentInput(t *testing.T) {
	t.Parallel()

	periodStart := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	subject := access.Subject{Type: access.SubjectProfile, ID: "42"}
	quotaStore := memoryquota.NewStore()
	quotaService := quota.NewService(quotaStore)
	saver := newFakeSaveStore()

	service := NewService(
		fakeSourceStore{parcel: Parcel{ID: 99, LandNumber: "44"}},
		saver,
		fakeAccessClient{result: access.Result{
			Subject:     subject,
			Operation:   commonoperation.Code("workspace.report.generate"),
			Allowed:     true,
			Metering:    access.Metering{FeatureCode: "workspace.report.generate", MeterCode: "workspace.report_generation.accepted", UnitsPerAction: 1, PolicyVersion: "1.0.0"},
			Limit:       100,
			Period:      access.PeriodSubscriptionCycle,
			PeriodStart: periodStart,
			PeriodEnd:   periodEnd,
		}},
		quotaService,
		true,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	ctx := testContext(t, "op_report_550", "idem_report_750")
	_, err := service.CreateGeneratedReport(ctx, Input{ReportType: 1, EntityType: 1, ParcelID: 99, Title: "A"})
	if err != nil {
		t.Fatalf("first CreateGeneratedReport() error = %v", err)
	}

	_, err = service.CreateGeneratedReport(ctx, Input{ReportType: 1, EntityType: 1, ParcelID: 99, Title: "B"})
	if !errors.Is(err, ErrCommandConflict) {
		t.Fatalf("retry error = %v, want ErrCommandConflict", err)
	}

	current, err := quotaStore.GetUsage(context.Background(), subject, "workspace.report_generation.accepted", periodStart, periodEnd)
	if err != nil {
		t.Fatal(err)
	}
	if current.Used != 1 || current.Reserved != 0 {
		t.Fatalf("usage = %+v, want one charge", current)
	}
}

func TestCreateGeneratedReportRetainsQuotaLeaseWhenDatabaseFails(t *testing.T) {
	t.Parallel()

	periodStart := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	subject := access.Subject{Type: access.SubjectProfile, ID: "42"}
	quotaStore := memoryquota.NewStore()
	quotaService := quota.NewService(quotaStore)
	saver := newFakeSaveStore()
	saver.fail = errors.New("database unavailable")

	service := NewService(
		fakeSourceStore{parcel: Parcel{ID: 99, LandNumber: "44"}},
		saver,
		fakeAccessClient{result: access.Result{
			Subject:     subject,
			Operation:   commonoperation.Code("workspace.report.generate"),
			Allowed:     true,
			Metering:    access.Metering{FeatureCode: "workspace.report.generate", MeterCode: "workspace.report_generation.accepted", UnitsPerAction: 1, PolicyVersion: "1.0.0"},
			Limit:       100,
			Period:      access.PeriodSubscriptionCycle,
			PeriodStart: periodStart,
			PeriodEnd:   periodEnd,
		}},
		quotaService,
		true,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	_, err := service.CreateGeneratedReport(
		testContext(t, "op_report_501", "idem_report_701"),
		Input{ReportType: 1, EntityType: 1, ParcelID: 99},
	)
	if err == nil {
		t.Fatal("CreateGeneratedReport() error = nil")
	}

	current, err := quotaStore.GetUsage(
		context.Background(),
		subject,
		"workspace.report_generation.accepted",
		periodStart,
		periodEnd,
	)
	if err != nil {
		t.Fatalf("GetUsage() error = %v", err)
	}
	if current.Used != 0 || current.Reserved != 1 {
		t.Fatalf("usage = %+v, want reservation retained for reconciliation", current)
	}
}

func TestCreateGeneratedReportRequiresCommandKey(t *testing.T) {
	t.Parallel()

	service := NewService(
		fakeSourceStore{},
		newFakeSaveStore(),
		fakeAccessClient{},
		quota.NewService(memoryquota.NewStore()),
		true,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	ctx := context.Background()
	var err error
	ctx, err = identity.BindActor(ctx, identity.Actor{ProfileID: 42})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = commonoperation.Bind(ctx, commonoperation.Code("workspace.report.generate"))
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = request.BindOperationID(ctx, "op_without_command_key")
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.CreateGeneratedReport(ctx, Input{ReportType: 1, EntityType: 1, ParcelID: 99})
	if !errors.Is(err, ErrCommandKeyMissing) {
		t.Fatalf("CreateGeneratedReport() error = %v, want ErrCommandKeyMissing", err)
	}
}

func TestCreateGeneratedReportReplaysBeforeCurrentAccess(t *testing.T) {
	t.Parallel()

	saver := newFakeSaveStore()
	input := Input{ReportType: 1, EntityType: 1, ParcelID: 99}
	hash, err := BuildRequestHash(input)
	if err != nil {
		t.Fatal(err)
	}
	saver.byKey["idem_replay"] = Report{
		ID:          7,
		UserID:      42,
		CommandKey:  "idem_replay",
		RequestHash: hash,
	}

	service := NewService(
		fakeSourceStore{},
		saver,
		fakeAccessClient{err: errors.New("entitlement dependency unavailable")},
		quota.NewService(memoryquota.NewStore()),
		false,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	result, err := service.CreateGeneratedReport(testContext(t, "op_replay", "idem_replay"), input)
	if err != nil {
		t.Fatalf("CreateGeneratedReport() replay error = %v", err)
	}
	if result.Outcome != OutcomeExisting || result.Report.ID != 7 {
		t.Fatalf("replay result = %+v", result)
	}
}

func TestCreateGeneratedReportReplaysWhenProcessingUnavailable(t *testing.T) {
	t.Parallel()

	operation := commonoperation.Code("workspace.report.generate")
	input := Input{ReportType: 1, EntityType: 1, ParcelID: 99}
	hash, err := BuildRequestHash(input)
	if err != nil {
		t.Fatal(err)
	}
	saver := newFakeSaveStore()
	saver.byKey["idem_existing"] = Report{
		ID:          8,
		UserID:      42,
		Operation:   operation,
		OperationID: "op_original",
		CommandKey:  "idem_existing",
		RequestHash: hash,
	}

	service := NewService(
		fakeSourceStore{},
		saver,
		fakeAccessClient{result: access.Result{
			Operation: operation,
			Allowed:   true,
		}},
		quota.NewService(memoryquota.NewStore()),
		false,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	result, err := service.CreateGeneratedReport(testContext(t, "op_retry", "idem_existing"), input)
	if err != nil {
		t.Fatalf("CreateGeneratedReport() error = %v", err)
	}
	if result.Outcome != OutcomeExisting || result.Report.ID != 8 {
		t.Fatalf("replay result = %+v", result)
	}
}

func TestCreateGeneratedReportBlocksNewWorkWhenProcessingUnavailable(t *testing.T) {
	t.Parallel()

	operation := commonoperation.Code("workspace.report.generate")
	saver := newFakeSaveStore()
	service := NewService(
		fakeSourceStore{parcel: Parcel{ID: 99, LandNumber: "44"}},
		saver,
		fakeAccessClient{result: access.Result{
			Operation: operation,
			Allowed:   true,
		}},
		quota.NewService(memoryquota.NewStore()),
		false,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	_, err := service.CreateGeneratedReport(
		testContext(t, "op_processing_down", "idem_processing_down"),
		Input{ReportType: 1, EntityType: 1, ParcelID: 99},
	)
	if !errors.Is(err, ErrProcessingUnavailable) {
		t.Fatalf("CreateGeneratedReport() error = %v, want ErrProcessingUnavailable", err)
	}
	if len(saver.byKey) != 0 || len(saver.jobs) != 0 || len(saver.usage) != 0 {
		t.Fatalf("processing-unavailable command persisted state: reports=%d jobs=%d usage=%d", len(saver.byKey), len(saver.jobs), len(saver.usage))
	}
}

func testContext(t *testing.T, operationID, idempotencyKey string) context.Context {
	t.Helper()

	ctx, err := identity.BindActor(context.Background(), identity.Actor{ProfileID: 42})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = commonoperation.Bind(ctx, commonoperation.Code("workspace.report.generate"))
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = request.BindOperationID(ctx, operationID)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = request.BindIdempotencyKey(ctx, idempotencyKey)
	if err != nil {
		t.Fatal(err)
	}
	return ctx
}

func TestCreateGeneratedReportNormalizesAccessDependencyFailure(t *testing.T) {
	providerErr := errors.New("user access dependency failed")

	service := NewService(
		fakeSourceStore{},
		newFakeSaveStore(),
		fakeAccessClient{err: providerErr},
		quota.NewService(memoryquota.NewStore()),
		true,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	_, err := service.CreateGeneratedReport(
		testContext(t, "op_access_failure", "idem_access_failure"),
		Input{ReportType: 1, EntityType: 1, ParcelID: 99},
	)

	if !errors.Is(err, ErrAccessUnavailable) {
		t.Fatalf("error = %v, want ErrAccessUnavailable", err)
	}

	if !errors.Is(err, providerErr) {
		t.Fatalf("provider cause was lost: %v", err)
	}
}
