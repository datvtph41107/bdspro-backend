package shadow

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"common/productcatalog"
	"user/internal/models"
	"user/internal/modules/catalog/domain"
	"user/internal/modules/catalog/migration"
)

type fakeRepository struct {
	mappings   map[uint64]Mapping
	aggregates map[uint64]domain.PlanVersionAggregate
	err        error
	calls      int
}

func (f *fakeRepository) GetLegacyPlanMapping(
	ctx context.Context,
	legacyPlanID uint64,
) (Mapping, error) {
	f.calls++
	if f.err != nil {
		return Mapping{}, f.err
	}
	mapping, ok := f.mappings[legacyPlanID]
	if !ok {
		return Mapping{}, ErrNotFound
	}
	return mapping, nil
}

func (f *fakeRepository) GetPlanVersion(
	ctx context.Context,
	planVersionID uint64,
) (domain.PlanVersionAggregate, error) {
	f.calls++
	if f.err != nil {
		return domain.PlanVersionAggregate{}, f.err
	}
	aggregate, ok := f.aggregates[planVersionID]
	if !ok {
		return domain.PlanVersionAggregate{}, ErrNotFound
	}
	return aggregate.Copy(), nil
}

type captureLogger struct {
	messages []string
}

func (l *captureLogger) Printf(format string, values ...any) {
	l.messages = append(l.messages, fmt.Sprintf(format, values...))
}

func legacyPlan() models.Plan {
	plan := models.Plan{
		PackageName: "Pro",
		Duration:    30,
		Price:       99000,
		Visibility:  true,
		Features:    []string{"Generated reports"},
	}
	plan.ID = 7
	return plan
}

func activeAggregate(at time.Time) domain.PlanVersionAggregate {
	from := at.Add(-time.Hour)
	published := at.Add(-2 * time.Hour)
	aggregate := domain.PlanVersionAggregate{
		ID:                 19,
		ProductCode:        "qhpro",
		ProductDisplayName: "QHPro",
		PlanCode:           "qhpro.pro",
		Version:            "1.0.0",
		DisplayName:        "Pro",
		Status:             productcatalog.StatusActive,
		SubjectScope:       productcatalog.SubjectScopeProfile,
		EffectiveFrom:      &from,
		PublishedAt:        &published,
		Entitlements: []productcatalog.Entitlement{
			{
				Code:        "workspace.report.feature",
				Kind:        productcatalog.EntitlementFeatureAccess,
				FeatureCode: "workspace.report.generate",
				Period:      productcatalog.PeriodNone,
			},
		},
	}
	checksum, err := productcatalog.TermsChecksum(aggregate.CurrentContract())
	if err != nil {
		panic(err)
	}
	aggregate.TermsChecksum = checksum
	return aggregate
}

func sourceChecksum(t *testing.T, plan models.Plan) string {
	t.Helper()
	checksum, err := migration.SourceChecksum(toLegacySnapshot(plan))
	if err != nil {
		t.Fatal(err)
	}
	return checksum
}

func TestCompareVisiblePlansMatchesApprovedMapping(t *testing.T) {
	at := time.Date(2026, 8, 4, 3, 0, 0, 0, time.UTC)
	plan := legacyPlan()
	planVersionID := uint64(19)
	repo := &fakeRepository{
		mappings: map[uint64]Mapping{
			7: {
				LegacyPlanID:   7,
				PlanVersionID:  &planVersionID,
				SourceChecksum: sourceChecksum(t, plan),
				Status:         models.CatalogMigrationShadowRead,
			},
		},
		aggregates: map[uint64]domain.PlanVersionAggregate{
			19: activeAggregate(at),
		},
	}
	reader := BuildReaderWithConfig(repo, Config{Mode: ModeLog})

	report := reader.CompareVisiblePlans(context.Background(), []models.Plan{plan}, at)
	if len(report.Entries) != 1 || report.Entries[0].Status != StatusMatched {
		t.Fatalf("report = %+v", report)
	}
}

func TestCompareVisiblePlansDetectsLegacyChange(t *testing.T) {
	at := time.Date(2026, 8, 4, 3, 0, 0, 0, time.UTC)
	plan := legacyPlan()
	oldPlan := plan
	oldPlan.Price = 1
	planVersionID := uint64(19)
	repo := &fakeRepository{
		mappings: map[uint64]Mapping{
			7: {
				LegacyPlanID:   7,
				PlanVersionID:  &planVersionID,
				SourceChecksum: sourceChecksum(t, oldPlan),
				Status:         models.CatalogMigrationShadowRead,
			},
		},
		aggregates: map[uint64]domain.PlanVersionAggregate{
			19: activeAggregate(at),
		},
	}
	reader := BuildReaderWithConfig(repo, Config{Mode: ModeLog})

	report := reader.CompareVisiblePlans(context.Background(), []models.Plan{plan}, at)
	if report.Entries[0].Status != StatusLegacyChanged {
		t.Fatalf("status = %q", report.Entries[0].Status)
	}
}

func TestCompareVisiblePlansDetectsCatalogDrift(t *testing.T) {
	at := time.Date(2026, 8, 4, 3, 0, 0, 0, time.UTC)
	plan := legacyPlan()
	planVersionID := uint64(19)
	aggregate := activeAggregate(at)
	aggregate.TermsChecksum = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	repo := &fakeRepository{
		mappings: map[uint64]Mapping{
			7: {
				LegacyPlanID:   7,
				PlanVersionID:  &planVersionID,
				SourceChecksum: sourceChecksum(t, plan),
				Status:         models.CatalogMigrationShadowRead,
			},
		},
		aggregates: map[uint64]domain.PlanVersionAggregate{19: aggregate},
	}
	reader := BuildReaderWithConfig(repo, Config{Mode: ModeLog})

	report := reader.CompareVisiblePlans(context.Background(), []models.Plan{plan}, at)
	if report.Entries[0].Status != StatusCatalogTermsDrift {
		t.Fatalf("status = %q", report.Entries[0].Status)
	}
}

func TestObserveVisiblePlansOffDoesNotRead(t *testing.T) {
	repo := &fakeRepository{}
	reader := BuildReaderWithConfig(repo, Config{Mode: ModeOff})
	reader.ObserveVisiblePlans(context.Background(), []models.Plan{legacyPlan()})
	if repo.calls != 0 {
		t.Fatalf("repository calls = %d", repo.calls)
	}
}

func TestObserveVisiblePlansLogsReadFailureWithoutPanic(t *testing.T) {
	logger := &captureLogger{}
	repo := &fakeRepository{err: errors.New("database unavailable")}
	reader := BuildReaderWithConfig(repo, Config{
		Mode:    ModeLog,
		Timeout: time.Second,
		Now: func() time.Time {
			return time.Date(2026, 8, 4, 3, 0, 0, 0, time.UTC)
		},
		Logger: logger,
	})

	reader.ObserveVisiblePlans(context.Background(), []models.Plan{legacyPlan()})
	if len(logger.messages) != 1 {
		t.Fatalf("log messages = %d", len(logger.messages))
	}
}

func TestValidateMode(t *testing.T) {
	for _, value := range []string{"", "off", "log", "LOG"} {
		if err := ValidateMode(value); err != nil {
			t.Fatalf("ValidateMode(%q) error = %v", value, err)
		}
	}
	if err := ValidateMode("enforce"); err == nil {
		t.Fatal("expected invalid mode error")
	}
}
