package read

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"common/productcatalog"
	"user/internal/models"
	"user/internal/modules/catalog/migration"
	"user/internal/modules/subscription/domain"
	"user/internal/modules/subscription/shadow"
)

type fakeRepository struct {
	membershipMapping  shadow.Mapping
	planMapping        shadow.PlanMapping
	subscription       domain.Aggregate
	err                error
	calls              int
	blockUntilCanceled bool
}

func (f *fakeRepository) GetLegacyMembershipMapping(ctx context.Context, _ uint64) (shadow.Mapping, error) {
	f.calls++
	if f.blockUntilCanceled {
		<-ctx.Done()
		return shadow.Mapping{}, ctx.Err()
	}
	if f.err != nil {
		return shadow.Mapping{}, f.err
	}
	return f.membershipMapping, nil
}

func (f *fakeRepository) GetLegacyPlanMappingForSubscription(context.Context, uint64) (shadow.PlanMapping, error) {
	f.calls++
	if f.err != nil {
		return shadow.PlanMapping{}, f.err
	}
	return f.planMapping, nil
}

func (f *fakeRepository) GetSubscription(context.Context, uint64) (domain.Aggregate, error) {
	f.calls++
	if f.err != nil {
		return domain.Aggregate{}, f.err
	}
	return f.subscription.Copy(), nil
}

type captureLogger struct{ messages []string }

func (l *captureLogger) Printf(format string, values ...any) {
	l.messages = append(l.messages, fmt.Sprintf(format, values...))
}

func legacyFixture(t *testing.T) models.Membership {
	t.Helper()
	createdAt := time.Date(2026, 7, 1, 8, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	paidAt := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	finishedAt := paidAt.AddDate(0, 1, 0)
	membership := models.Membership{
		ProfileID:    42,
		PlanID:       7,
		PaidAt:       &paidAt,
		FinishedDate: &finishedAt,
		OrderID:      88,
	}
	membership.ID = 9
	membership.CreatedAt = &createdAt
	membership.UpdatedAt = &updatedAt
	return membership
}

func matchingRepo(t *testing.T, status models.CatalogMigrationStatus) *fakeRepository {
	t.Helper()
	legacy := legacyFixture(t)
	checksum, err := migration.MembershipSourceChecksum(migration.LegacyMembershipSnapshot{
		ID: legacy.ID, ProfileID: legacy.ProfileID, PlanID: legacy.PlanID,
		FinishedDate: legacy.FinishedDate, OrderID: legacy.OrderID, PaidAt: legacy.PaidAt,
		CreatedAt: legacy.CreatedAt, UpdatedAt: legacy.UpdatedAt,
	})
	if err != nil {
		t.Fatal(err)
	}
	subscriptionID := uint64(19)
	planVersionID := uint64(11)
	return &fakeRepository{
		membershipMapping: shadow.Mapping{
			LegacyMembershipID: legacy.ID,
			SubscriptionID:     &subscriptionID,
			SourceChecksum:     checksum,
			Status:             status,
		},
		planMapping: shadow.PlanMapping{
			LegacyPlanID:  legacy.PlanID,
			PlanVersionID: &planVersionID,
			Status:        models.CatalogMigrationShadowRead,
		},
		subscription: domain.Aggregate{
			ID:                 subscriptionID,
			SubscriptionKey:    "sub_legacy_membership_9",
			SubjectKind:        models.SubscriptionSubjectProfile,
			SubjectID:          "42",
			ProductID:          3,
			ProductCode:        "qhpro",
			PlanVersionID:      planVersionID,
			PlanCode:           "qhpro.pro",
			PlanVersion:        "1.0.0",
			PlanTermsChecksum:  "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			PlanStatus:         productcatalog.StatusActive,
			Status:             models.SubscriptionActive,
			StartedAt:          *legacy.PaidAt,
			CurrentPeriodStart: *legacy.PaidAt,
			CurrentPeriodEnd:   *legacy.FinishedDate,
			AccessUntil:        legacy.FinishedDate,
			OrderReference:     "88",
			CreatedAt:          *legacy.CreatedAt,
			UpdatedAt:          *legacy.UpdatedAt,
		},
	}
}

func TestLegacyModeDoesNotReadCatalog(t *testing.T) {
	repo := matchingRepo(t, models.CatalogMigrationShadowRead)
	router := BuildRouterWithConfig(repo, Config{Mode: ModeLegacy})
	legacy := legacyFixture(t)
	selected, err := router.SelectMembership(context.Background(), &legacy)
	if err != nil {
		t.Fatal(err)
	}
	if repo.calls != 0 || selected == &legacy {
		t.Fatalf("calls=%d selected aliases legacy=%t", repo.calls, selected == &legacy)
	}
}

func TestCanaryFallbackDoesNotChangeLegacyResponse(t *testing.T) {
	repo := matchingRepo(t, models.CatalogMigrationShadowRead)
	repo.err = errors.New("database unavailable")
	logger := &captureLogger{}
	router := BuildRouterWithConfig(repo, Config{
		Mode: ModeCanary, CanaryPercent: 100, CanarySalt: "stable", Logger: logger,
	})
	legacy := legacyFixture(t)
	selected, err := router.SelectMembership(context.Background(), &legacy)
	if err != nil {
		t.Fatal(err)
	}
	if selected.ID != legacy.ID || selected.PlanID != legacy.PlanID || len(logger.messages) == 0 {
		t.Fatalf("selected=%+v logs=%v", selected, logger.messages)
	}
}

func TestCanaryReturnsCompatibleCatalogProjection(t *testing.T) {
	repo := matchingRepo(t, models.CatalogMigrationShadowRead)
	router := BuildRouterWithConfig(repo, Config{
		Mode: ModeCanary, CanaryPercent: 100, CanarySalt: "stable",
		Now: func() time.Time { return time.Date(2026, 8, 4, 0, 0, 0, 0, time.UTC) },
	})
	legacy := legacyFixture(t)
	selected, err := router.SelectMembership(context.Background(), &legacy)
	if err != nil {
		t.Fatal(err)
	}
	if selected.ID != legacy.ID || selected.ProfileID != legacy.ProfileID || selected.PlanID != legacy.PlanID {
		t.Fatalf("immutable API identity changed: %+v", selected)
	}
	if !selected.CreatedAt.Equal(*legacy.CreatedAt) || !selected.UpdatedAt.Equal(*legacy.UpdatedAt) {
		t.Fatal("legacy audit timestamps changed")
	}
}

func TestCatalogModeUsesLegacyUntilMappingIsCutover(t *testing.T) {
	repo := matchingRepo(t, models.CatalogMigrationShadowRead)
	router := BuildRouterWithConfig(repo, Config{Mode: ModeCatalog})
	legacy := legacyFixture(t)
	selected, err := router.SelectMembership(context.Background(), &legacy)
	if err != nil {
		t.Fatal(err)
	}
	if selected.ID != legacy.ID || repo.calls != 1 {
		t.Fatalf("selected=%+v calls=%d", selected, repo.calls)
	}
}

func TestCatalogModeFailsClosedForInvalidCutover(t *testing.T) {
	repo := matchingRepo(t, models.CatalogMigrationCutover)
	repo.subscription.OrderReference = "wrong"
	router := BuildRouterWithConfig(repo, Config{
		Mode: ModeCatalog,
		Now:  func() time.Time { return time.Date(2026, 8, 4, 0, 0, 0, 0, time.UTC) },
	})
	legacy := legacyFixture(t)
	selected, err := router.SelectMembership(context.Background(), &legacy)
	if selected != nil || !errors.Is(err, ErrCatalogCutoverInvalid) {
		t.Fatalf("selected=%+v err=%v", selected, err)
	}
}

func TestCanaryReadPreservesRequestCancellation(t *testing.T) {
	repo := matchingRepo(t, models.CatalogMigrationShadowRead)
	repo.blockUntilCanceled = true
	router := BuildRouterWithConfig(repo, Config{
		Mode: ModeCanary, CanaryPercent: 100, CanarySalt: "stable",
		Timeout: time.Second,
	})
	legacy := legacyFixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	started := time.Now()
	selected, err := router.SelectMembership(ctx, &legacy)
	if err != nil {
		t.Fatal(err)
	}
	if selected.ID != legacy.ID || time.Since(started) > 100*time.Millisecond {
		t.Fatalf("cancellation was detached: selected=%+v elapsed=%s", selected, time.Since(started))
	}
}

func TestCanaryCohortIsStable(t *testing.T) {
	router := BuildRouterWithConfig(nil, Config{
		Mode: ModeCanary, CanaryPercent: 50, CanarySalt: "deployment-20260804",
	})
	first := router.inCanary(42)
	for index := 0; index < 100; index++ {
		if router.inCanary(42) != first {
			t.Fatal("canary cohort changed for the same profile")
		}
	}
}

func TestValidateConfig(t *testing.T) {
	if err := ValidateConfig(Config{Mode: ModeCanary, CanaryPercent: 10}); err == nil {
		t.Fatal("expected missing salt error")
	}
	if err := ValidateConfig(Config{Mode: ModeCatalog}); err != nil {
		t.Fatal(err)
	}
}
