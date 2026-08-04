package subscriptionshadow

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"common/productcatalog"
	"user/internal/catalogmigration"
	"user/internal/models"
	"user/internal/subscriptiondomain"
)

type fakeRepository struct {
	membershipMappings map[uint64]Mapping
	planMappings       map[uint64]PlanMapping
	subscriptions      map[uint64]subscriptiondomain.Aggregate
	err                error
	calls              int
}

func (f *fakeRepository) GetLegacyMembershipMapping(
	ctx context.Context,
	legacyMembershipID uint64,
) (Mapping, error) {
	f.calls++
	if f.err != nil {
		return Mapping{}, f.err
	}
	mapping, ok := f.membershipMappings[legacyMembershipID]
	if !ok {
		return Mapping{}, ErrNotFound
	}
	return mapping, nil
}

func (f *fakeRepository) GetLegacyPlanMappingForSubscription(
	ctx context.Context,
	legacyPlanID uint64,
) (PlanMapping, error) {
	f.calls++
	if f.err != nil {
		return PlanMapping{}, f.err
	}
	mapping, ok := f.planMappings[legacyPlanID]
	if !ok {
		return PlanMapping{}, ErrNotFound
	}
	return mapping, nil
}

func (f *fakeRepository) GetSubscription(
	ctx context.Context,
	subscriptionID uint64,
) (subscriptiondomain.Aggregate, error) {
	f.calls++
	if f.err != nil {
		return subscriptiondomain.Aggregate{}, f.err
	}
	aggregate, ok := f.subscriptions[subscriptionID]
	if !ok {
		return subscriptiondomain.Aggregate{}, ErrNotFound
	}
	return aggregate.Copy(), nil
}

type captureLogger struct {
	messages []string
}

func (l *captureLogger) Printf(format string, values ...any) {
	l.messages = append(l.messages, fmt.Sprintf(format, values...))
}

func legacyMembership() models.Membership {
	paidAt := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	finishedAt := paidAt.AddDate(0, 1, 0)
	membership := models.Membership{
		ProfileID:    42,
		PlanID:       7,
		FinishedDate: &finishedAt,
		OrderID:      88,
		PaidAt:       &paidAt,
	}
	membership.ID = 9
	createdAt := paidAt.Add(-24 * time.Hour)
	updatedAt := paidAt
	membership.CreatedAt = &createdAt
	membership.UpdatedAt = &updatedAt
	return membership
}

func activeSubscription(membership models.Membership) subscriptiondomain.Aggregate {
	return subscriptiondomain.Aggregate{
		ID:                 19,
		SubscriptionKey:    "sub_legacy_membership_9",
		SubjectKind:        models.SubscriptionSubjectProfile,
		SubjectID:          "42",
		ProductID:          3,
		ProductCode:        "qhpro",
		PlanVersionID:      11,
		PlanCode:           "qhpro.pro",
		PlanVersion:        "1.0.0",
		PlanTermsChecksum:  "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		PlanStatus:         productcatalog.StatusActive,
		Status:             models.SubscriptionActive,
		StartedAt:          *membership.PaidAt,
		CurrentPeriodStart: *membership.PaidAt,
		CurrentPeriodEnd:   *membership.FinishedDate,
		AccessUntil:        copyTime(membership.FinishedDate),
		OrderReference:     "88",
		CreatedAt:          *membership.CreatedAt,
		UpdatedAt:          *membership.UpdatedAt,
	}
}

func matchingRepository(t *testing.T, membership models.Membership) *fakeRepository {
	t.Helper()
	checksum, err := catalogmigration.MembershipSourceChecksum(toLegacySnapshot(membership))
	if err != nil {
		t.Fatal(err)
	}
	subscriptionID := uint64(19)
	planVersionID := uint64(11)
	return &fakeRepository{
		membershipMappings: map[uint64]Mapping{
			9: {
				LegacyMembershipID: 9,
				SubscriptionID:     &subscriptionID,
				SourceChecksum:     checksum,
				Status:             models.CatalogMigrationShadowRead,
			},
		},
		planMappings: map[uint64]PlanMapping{
			7: {
				LegacyPlanID:  7,
				PlanVersionID: &planVersionID,
				Status:        models.CatalogMigrationShadowRead,
			},
		},
		subscriptions: map[uint64]subscriptiondomain.Aggregate{
			19: activeSubscription(membership),
		},
	}
}

func TestCompareMembershipMatchesApprovedSubscription(t *testing.T) {
	membership := legacyMembership()
	reader := NewReaderWithConfig(matchingRepository(t, membership), Config{Mode: ModeLog})
	entry := reader.CompareMembership(
		context.Background(),
		membership,
		membership.PaidAt.Add(time.Hour),
	)
	if entry.Status != StatusMatched {
		t.Fatalf("entry = %+v", entry)
	}
	if entry.MappingStatus != models.CatalogMigrationShadowRead ||
		entry.PlanMappingStatus != models.CatalogMigrationShadowRead {
		t.Fatalf("mapping evidence = %+v", entry)
	}
}

func TestCompareMembershipDetectsLegacyChange(t *testing.T) {
	membership := legacyMembership()
	repo := matchingRepository(t, membership)
	membership.OrderID++
	reader := NewReaderWithConfig(repo, Config{Mode: ModeLog})
	entry := reader.CompareMembership(context.Background(), membership, time.Now())
	if entry.Status != StatusLegacyChanged {
		t.Fatalf("status = %q", entry.Status)
	}
}

func TestCompareMembershipDetectsPlanMappingMismatch(t *testing.T) {
	membership := legacyMembership()
	repo := matchingRepository(t, membership)
	aggregate := repo.subscriptions[19]
	aggregate.PlanVersionID = 99
	repo.subscriptions[19] = aggregate
	reader := NewReaderWithConfig(repo, Config{Mode: ModeLog})
	entry := reader.CompareMembership(context.Background(), membership, membership.PaidAt.Add(time.Hour))
	if entry.Status != StatusPlanMappingMismatch {
		t.Fatalf("status = %q", entry.Status)
	}
}

func TestCompareMembershipDetectsStatusAndPeriodMismatch(t *testing.T) {
	membership := legacyMembership()
	repo := matchingRepository(t, membership)
	aggregate := repo.subscriptions[19]
	aggregate.Status = models.SubscriptionPending
	repo.subscriptions[19] = aggregate
	reader := NewReaderWithConfig(repo, Config{Mode: ModeLog})
	entry := reader.CompareMembership(context.Background(), membership, membership.PaidAt.Add(time.Hour))
	if entry.Status != StatusStatusMismatch {
		t.Fatalf("status = %q", entry.Status)
	}

	aggregate.Status = models.SubscriptionActive
	aggregate.CurrentPeriodEnd = aggregate.CurrentPeriodEnd.Add(time.Hour)
	repo.subscriptions[19] = aggregate
	entry = reader.CompareMembership(context.Background(), membership, membership.PaidAt.Add(time.Hour))
	if entry.Status != StatusPeriodMismatch {
		t.Fatalf("status = %q", entry.Status)
	}
}

func TestCompareMembershipTreatsFinishedDateAsExclusive(t *testing.T) {
	membership := legacyMembership()
	repo := matchingRepository(t, membership)
	aggregate := repo.subscriptions[19]
	aggregate.Status = models.SubscriptionExpired
	repo.subscriptions[19] = aggregate
	reader := NewReaderWithConfig(repo, Config{Mode: ModeLog})
	entry := reader.CompareMembership(context.Background(), membership, *membership.FinishedDate)
	if entry.Status != StatusMatched {
		t.Fatalf("entry = %+v", entry)
	}
}

func TestObserveMembershipOffDoesNotRead(t *testing.T) {
	repo := &fakeRepository{}
	reader := NewReaderWithConfig(repo, Config{Mode: ModeOff})
	membership := legacyMembership()
	reader.ObserveMembership(context.Background(), &membership)
	if repo.calls != 0 {
		t.Fatalf("repository calls = %d", repo.calls)
	}
}

func TestObserveMembershipLogsReadFailureWithoutPanic(t *testing.T) {
	logger := &captureLogger{}
	repo := &fakeRepository{err: errors.New("database unavailable")}
	reader := NewReaderWithConfig(repo, Config{
		Mode:    ModeLog,
		Timeout: time.Second,
		Now: func() time.Time {
			return time.Date(2026, 8, 4, 3, 0, 0, 0, time.UTC)
		},
		Logger: logger,
	})
	membership := legacyMembership()
	reader.ObserveMembership(context.Background(), &membership)
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
