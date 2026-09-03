package publish

import (
	"context"
	"errors"
	"testing"
	"time"

	catalogdomain "user/internal/domain/plan"
)

type fakeTransaction struct {
	calls     int
	committed bool
}

func (f *fakeTransaction) WithTransaction(
	ctx context.Context,
	fn func(context.Context) error,
) error {
	f.calls++
	if err := fn(context.WithValue(ctx, "tx-test", true)); err != nil {
		return err
	}
	f.committed = true
	return nil
}

type fakeRepository struct {
	aggregate  catalogdomain.PlanVersionAggregate
	loadErr    error
	publishErr error
	published  *PublishRecord
}

func (f *fakeRepository) LoadPlanVersionForUpdate(
	ctx context.Context,
	planVersionID uint64,
) (catalogdomain.PlanVersionAggregate, error) {
	if f.loadErr != nil {
		return catalogdomain.PlanVersionAggregate{}, f.loadErr
	}
	return f.aggregate.Copy(), nil
}

func (f *fakeRepository) MarkPlanVersionPublished(
	ctx context.Context,
	input PublishRecord,
) error {
	if f.publishErr != nil {
		return f.publishErr
	}
	copied := input
	f.published = &copied
	return nil
}

func validAggregate() catalogdomain.PlanVersionAggregate {
	return catalogdomain.PlanVersionAggregate{
		ID:                   9,
		ProductID:            1,
		PlanID:               2,
		ProductCode:          "qhpro",
		ProductDisplayName:   "QHPro",
		ProductStatus:        catalogdomain.StatusActive,
		PlanCode:             "qhpro.pro",
		PlanStatus:           catalogdomain.StatusActive,
		Version:              "1.0.0",
		DisplayName:          "Pro",
		Status:               catalogdomain.StatusDraft,
		SubjectScope:         catalogdomain.SubjectScopeProfile,
		SubscriptionTermDays: 30,
		Entitlements: []catalogdomain.Entitlement{
			{
				Code:        "workspace.report.feature",
				Kind:        catalogdomain.EntitlementFeatureAccess,
				FeatureCode: "workspace.report.generate",
				Period:      catalogdomain.PeriodNone,
			},
		},
		Operations: []catalogdomain.OperationBinding{
			{
				Code:        "workspace.report.generate",
				FeatureCode: "workspace.report.generate",
			},
		},
		Prices: []catalogdomain.PriceItem{
			{
				Code:        "qhpro.pro.recurring",
				Kind:        catalogdomain.PriceRecurring,
				Currency:    "VND",
				AmountMinor: 9900000,
				BillingUnit: "subscription.month",
				Quantity:    1,
			},
		},
	}
}

func TestPublishCommitsValidatedDraft(t *testing.T) {
	tx := &fakeTransaction{}
	repo := &fakeRepository{aggregate: validAggregate()}
	service := NewService(tx, repo)
	publishedAt := time.Date(2026, 8, 4, 2, 0, 0, 0, time.UTC)
	effectiveFrom := publishedAt.Add(time.Hour)

	result, err := service.Publish(context.Background(), Command{
		PlanVersionID: 9,
		ActorID:       42,
		PublishedAt:   publishedAt,
		EffectiveFrom: effectiveFrom,
	})
	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if tx.calls != 1 || !tx.committed {
		t.Fatal("publication transaction was not committed")
	}
	if repo.published == nil {
		t.Fatal("publish update was not written")
	}
	if repo.published.PublishedBy != 42 ||
		repo.published.TermsChecksum == "" {
		t.Fatal("publish record is incomplete")
	}
	if result.TermsChecksum != repo.published.TermsChecksum {
		t.Fatal("result checksum does not match persisted checksum")
	}
}

func TestPublishRejectsNonDraft(t *testing.T) {
	aggregate := validAggregate()
	aggregate.Status = catalogdomain.StatusActive
	repo := &fakeRepository{aggregate: aggregate}
	service := NewService(&fakeTransaction{}, repo)

	_, err := service.Publish(context.Background(), Command{
		PlanVersionID: 9,
		ActorID:       42,
		PublishedAt:   time.Date(2026, 8, 4, 2, 0, 0, 0, time.UTC),
		EffectiveFrom: time.Date(2026, 8, 4, 3, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, ErrPlanVersionNotDraft) {
		t.Fatalf("Publish() error = %v", err)
	}
	if repo.published != nil {
		t.Fatal("non-draft version was published")
	}
}

func TestPublishRejectsRetiredParent(t *testing.T) {
	aggregate := validAggregate()
	aggregate.PlanStatus = catalogdomain.StatusRetired
	repo := &fakeRepository{aggregate: aggregate}
	service := NewService(&fakeTransaction{}, repo)

	_, err := service.Publish(context.Background(), Command{
		PlanVersionID: 9,
		ActorID:       42,
		PublishedAt:   time.Date(2026, 8, 4, 2, 0, 0, 0, time.UTC),
		EffectiveFrom: time.Date(2026, 8, 4, 3, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, ErrParentRetired) {
		t.Fatalf("Publish() error = %v", err)
	}
	if repo.published != nil {
		t.Fatal("retired parent was published")
	}
}

func TestPublishRollsBackRepositoryFailure(t *testing.T) {
	tx := &fakeTransaction{}
	repo := &fakeRepository{
		aggregate:  validAggregate(),
		publishErr: ErrConcurrentPublish,
	}
	service := NewService(tx, repo)

	_, err := service.Publish(context.Background(), Command{
		PlanVersionID: 9,
		ActorID:       42,
		PublishedAt:   time.Date(2026, 8, 4, 2, 0, 0, 0, time.UTC),
		EffectiveFrom: time.Date(2026, 8, 4, 3, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, ErrConcurrentPublish) {
		t.Fatalf("Publish() error = %v", err)
	}
	if tx.committed {
		t.Fatal("failed publication transaction was committed")
	}
}

func TestPublishRejectsInvalidWindowBeforeTransaction(t *testing.T) {
	tx := &fakeTransaction{}
	service := NewService(tx, &fakeRepository{aggregate: validAggregate()})
	from := time.Date(2026, 8, 4, 3, 0, 0, 0, time.UTC)
	until := from

	_, err := service.Publish(context.Background(), Command{
		PlanVersionID:  9,
		ActorID:        42,
		PublishedAt:    from.Add(-time.Hour),
		EffectiveFrom:  from,
		EffectiveUntil: &until,
	})
	if err == nil {
		t.Fatal("expected invalid effective window error")
	}
	if tx.calls != 0 {
		t.Fatal("invalid command opened a transaction")
	}
}

func TestPublishRejectsMissingSubscriptionTerm(t *testing.T) {
	tx := &fakeTransaction{}
	aggregate := validAggregate()
	aggregate.SubscriptionTermDays = 0
	repo := &fakeRepository{aggregate: aggregate}
	service := NewService(tx, repo)

	_, err := service.Publish(context.Background(), Command{
		PlanVersionID: 9,
		ActorID:       42,
		PublishedAt:   time.Date(2026, 8, 4, 2, 0, 0, 0, time.UTC),
		EffectiveFrom: time.Date(2026, 8, 4, 3, 0, 0, 0, time.UTC),
	})
	if err == nil {
		t.Fatal("Publish() accepted missing subscription_term_days")
	}
	if repo.published != nil {
		t.Fatal("missing commercial terms reached persistence")
	}
}
