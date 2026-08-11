package report

import (
	commonaccess "common/access"
	"common/businessoperation"
	"common/identitycontext"
	"common/requestcontext"
	"context"
	"errors"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"
	"tqd/infra/memoryquota"
	"tqd/internal/quota"
	"tqd/internal/usage"
)

type fakeAccessClient struct {
	result commonaccess.Result
	err    error
}

func (f fakeAccessClient) GetReportAccess(context.Context) (commonaccess.Result, error) {
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
	forceFindMisses int
}

func newFakeSaveStore() *fakeSaveStore {
	return &fakeSaveStore{
		nextID: 1,
		byKey:  make(map[string]Report),
		usage:  make(map[string]usage.Event),
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

func (s *fakeSaveStore) AcceptReport(
	ctx context.Context,
	report Report,
	usageEvent *usage.Event,
) (Report, bool, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.fail != nil {
		return Report{}, false, false, s.fail
	}

	if existing, ok := s.byKey[report.CommandKey]; ok {
		finalizeReservation := usageEvent == nil
		if usageEvent != nil {
			if savedUsage, found := s.usage[usageEvent.UsageKey]; found {
				finalizeReservation = savedUsage.ReservationID == usageEvent.ReservationID
			}
		}
		return existing, false, finalizeReservation, nil
	}

	report.ID = s.nextID
	s.nextID++
	report.CreatedAt = time.Now().UTC()
	report.UpdatedAt = report.CreatedAt
	s.byKey[report.CommandKey] = report
	if usageEvent != nil {
		s.usage[usageEvent.UsageKey] = *usageEvent
	}
	return report, true, true, nil
}

func TestCreateGeneratedReportAcceptsAndChargesOnce(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 11, 5, 0, 0, 0, time.UTC)
	periodStart := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	subject := commonaccess.Subject{Type: commonaccess.SubjectProfile, ID: "42"}

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
		fakeAccessClient{result: commonaccess.Result{
			Subject:     subject,
			Operation:   businessoperation.WorkspaceReportGenerate,
			Allowed:     true,
			Limit:       100,
			Period:      commonaccess.PeriodSubscriptionCycle,
			PeriodStart: periodStart,
			PeriodEnd:   periodEnd,
		}},
		quotaService,
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
		businessoperation.WorkspaceReportGenerate,
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
}

func TestCreateGeneratedReportDoesNotChargeAgainWhenRaceCrossesPeriod(t *testing.T) {
	t.Parallel()

	firstStart := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	firstEnd := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	secondStart := firstEnd
	secondEnd := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	subject := commonaccess.Subject{Type: commonaccess.SubjectProfile, ID: "42"}
	quotaStore := memoryquota.NewStore()
	quotaService := quota.NewService(quotaStore)
	saver := newFakeSaveStore()
	source := fakeSourceStore{parcel: Parcel{ID: 99, LandNumber: "44"}}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	firstService := NewService(
		source,
		saver,
		fakeAccessClient{result: commonaccess.Result{
			Subject:     subject,
			Operation:   businessoperation.WorkspaceReportGenerate,
			Allowed:     true,
			Limit:       100,
			Period:      commonaccess.PeriodSubscriptionCycle,
			PeriodStart: firstStart,
			PeriodEnd:   firstEnd,
		}},
		quotaService,
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
		fakeAccessClient{result: commonaccess.Result{
			Subject:     subject,
			Operation:   businessoperation.WorkspaceReportGenerate,
			Allowed:     true,
			Limit:       100,
			Period:      commonaccess.PeriodSubscriptionCycle,
			PeriodStart: secondStart,
			PeriodEnd:   secondEnd,
		}},
		quotaService,
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

	firstUsage, err := quotaStore.GetUsage(context.Background(), subject, businessoperation.WorkspaceReportGenerate, firstStart, firstEnd)
	if err != nil {
		t.Fatal(err)
	}
	secondUsage, err := quotaStore.GetUsage(context.Background(), subject, businessoperation.WorkspaceReportGenerate, secondStart, secondEnd)
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
	subject := commonaccess.Subject{Type: commonaccess.SubjectProfile, ID: "42"}
	quotaStore := memoryquota.NewStore()
	quotaService := quota.NewService(quotaStore)
	saver := newFakeSaveStore()

	service := NewService(
		fakeSourceStore{parcel: Parcel{ID: 99, LandNumber: "44"}},
		saver,
		fakeAccessClient{result: commonaccess.Result{
			Subject:     subject,
			Operation:   businessoperation.WorkspaceReportGenerate,
			Allowed:     true,
			Limit:       100,
			Period:      commonaccess.PeriodSubscriptionCycle,
			PeriodStart: periodStart,
			PeriodEnd:   periodEnd,
		}},
		quotaService,
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

	current, err := quotaStore.GetUsage(context.Background(), subject, businessoperation.WorkspaceReportGenerate, periodStart, periodEnd)
	if err != nil {
		t.Fatal(err)
	}
	if current.Used != 1 || current.Reserved != 0 {
		t.Fatalf("usage = %+v, want one charge", current)
	}
}

func TestCreateGeneratedReportCancelsQuotaWhenDatabaseFails(t *testing.T) {
	t.Parallel()

	periodStart := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	subject := commonaccess.Subject{Type: commonaccess.SubjectProfile, ID: "42"}
	quotaStore := memoryquota.NewStore()
	quotaService := quota.NewService(quotaStore)
	saver := newFakeSaveStore()
	saver.fail = errors.New("database unavailable")

	service := NewService(
		fakeSourceStore{parcel: Parcel{ID: 99, LandNumber: "44"}},
		saver,
		fakeAccessClient{result: commonaccess.Result{
			Subject:     subject,
			Operation:   businessoperation.WorkspaceReportGenerate,
			Allowed:     true,
			Limit:       100,
			Period:      commonaccess.PeriodSubscriptionCycle,
			PeriodStart: periodStart,
			PeriodEnd:   periodEnd,
		}},
		quotaService,
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
		businessoperation.WorkspaceReportGenerate,
		periodStart,
		periodEnd,
	)
	if err != nil {
		t.Fatalf("GetUsage() error = %v", err)
	}
	if current.Used != 0 || current.Reserved != 0 {
		t.Fatalf("usage = %+v, want zero after database failure", current)
	}
}

func testContext(t *testing.T, operationID, idempotencyKey string) context.Context {
	t.Helper()

	ctx, err := identitycontext.BindActor(context.Background(), identitycontext.Actor{ProfileID: 42})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = businessoperation.Bind(ctx, businessoperation.WorkspaceReportGenerate)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = requestcontext.BindOperationID(ctx, operationID)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = requestcontext.BindIdempotencyKey(ctx, idempotencyKey)
	if err != nil {
		t.Fatal(err)
	}
	return ctx
}
