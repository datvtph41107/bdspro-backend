package backfill

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"common/productcatalog"
	"user/internal/models"
	"user/internal/modules/catalog/migration"
	"user/internal/modules/subscription/domain"
)

const testChecksum = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

type fakeRepository struct {
	candidate     Candidate
	subscriptions map[string]domain.Aggregate
	nextID        uint64
	events        map[uint64]Event
	blockedReason string
	reviewed      bool
	failEvent     error
}

func newFakeRepository(t *testing.T) *fakeRepository {
	t.Helper()
	paidAt := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	finishedAt := paidAt.AddDate(0, 1, 0)
	membership := models.Membership{
		ProfileID:    42,
		PlanID:       7,
		PaidAt:       &paidAt,
		FinishedDate: &finishedAt,
		OrderID:      88,
	}
	membership.ID = 9
	createdAt := paidAt.Add(-24 * time.Hour)
	updatedAt := paidAt
	membership.CreatedAt = &createdAt
	membership.UpdatedAt = &updatedAt
	checksum, err := migration.MembershipSourceChecksum(
		toLegacySnapshot(membership),
	)
	if err != nil {
		t.Fatal(err)
	}
	reviewer := uint64(3)
	reviewedAt := paidAt.Add(-time.Hour)
	decision, _ := json.Marshal(validDecision())
	publishedAt := paidAt.AddDate(0, -1, 0)
	effectiveFrom := paidAt.AddDate(-1, 0, 0)
	return &fakeRepository{
		candidate: Candidate{
			Membership: membership,
			Mapping: Mapping{
				LegacyMembershipID: 9,
				SourceChecksum:     checksum,
				Status:             models.CatalogMigrationReviewed,
				Decision:           decision,
				ReviewedBy:         &reviewer,
				ReviewedAt:         &reviewedAt,
			},
			Target: PlanTarget{
				LegacyPlanID:      7,
				ProductID:         1,
				ProductCode:       "qhpro",
				PlanVersionID:     11,
				PlanCode:          "qhpro.pro",
				PlanVersion:       "1.0.0",
				PlanTermsChecksum: testChecksum,
				PlanStatus:        productcatalog.StatusActive,
				PlanPublishedAt:   &publishedAt,
				SubjectScope:      productcatalog.SubjectScopeProfile,
				EffectiveFrom:     &effectiveFrom,
				MappingStatus:     models.CatalogMigrationShadowRead,
			},
		},
		subscriptions: make(map[string]domain.Aggregate),
		nextID:        20,
		events:        make(map[uint64]Event),
	}
}

func validDecision() ReviewDecision {
	return ReviewDecision{
		PeriodHandling:          PeriodPaidAtToFinishedDate,
		SubscriptionKeyStrategy: KeyStrategyLegacyMembershipID,
		AutoRenew:               false,
		Reason:                  "Legacy paid period approved by product owner",
	}
}

func validRun() Run {
	return Run{
		RunID:       "run_20260804_001",
		ActorID:     "catalog-migration",
		RequestID:   "req_task45",
		OperationID: "op_task45",
		AsOf:        time.Date(2026, 8, 4, 3, 0, 0, 0, time.UTC),
		ExecutedAt:  time.Date(2026, 8, 4, 3, 1, 0, 0, time.UTC),
	}
}

func (r *fakeRepository) clone() *fakeRepository {
	copied := *r
	copied.candidate = r.candidate
	copied.candidate.Mapping.Decision = append(json.RawMessage(nil), r.candidate.Mapping.Decision...)
	copied.subscriptions = make(map[string]domain.Aggregate, len(r.subscriptions))
	for key, value := range r.subscriptions {
		copied.subscriptions[key] = value.Copy()
	}
	copied.events = make(map[uint64]Event, len(r.events))
	for key, value := range r.events {
		copied.events[key] = value
	}
	return &copied
}

func (r *fakeRepository) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	working := r.clone()
	if err := fn(context.WithValue(ctx, fakeContextKey{}, working)); err != nil {
		return err
	}
	*r = *working
	return nil
}

type fakeContextKey struct{}

func fakeFromContext(ctx context.Context, fallback *fakeRepository) *fakeRepository {
	if repo, ok := ctx.Value(fakeContextKey{}).(*fakeRepository); ok {
		return repo
	}
	return fallback
}

func (r *fakeRepository) LoadCandidate(ctx context.Context, id uint64, _ bool) (Candidate, error) {
	repo := fakeFromContext(ctx, r)
	if repo.candidate.Membership.ID != id {
		return Candidate{}, ErrNotFound
	}
	candidate := repo.candidate
	if candidate.Mapping.SubscriptionID != nil {
		if subscription, ok := repo.subscriptions["sub_legacy_membership_9"]; ok {
			candidate.ExistingSubscription = &subscription
		}
	}
	return candidate, nil
}

func (r *fakeRepository) SaveReview(
	ctx context.Context,
	request ReviewRequest,
	checksum string,
	decision json.RawMessage,
) error {
	repo := fakeFromContext(ctx, r)
	repo.reviewed = true
	repo.candidate.Mapping.SourceChecksum = checksum
	repo.candidate.Mapping.Status = models.CatalogMigrationReviewed
	repo.candidate.Mapping.Decision = append(json.RawMessage(nil), decision...)
	repo.candidate.Mapping.ReviewedBy = &request.ReviewedBy
	repo.candidate.Mapping.ReviewedAt = &request.ReviewedAt
	return nil
}

func (r *fakeRepository) ListReviewedMembershipIDs(context.Context, uint64, int) ([]uint64, error) {
	return []uint64{9}, nil
}

func (r *fakeRepository) FindSubscriptionByKey(
	ctx context.Context,
	key string,
) (domain.Aggregate, error) {
	repo := fakeFromContext(ctx, r)
	value, ok := repo.subscriptions[key]
	if !ok {
		return domain.Aggregate{}, ErrNotFound
	}
	return value.Copy(), nil
}

func (r *fakeRepository) FindLiveSubscription(
	ctx context.Context,
	_ models.SubscriptionSubjectKind,
	_ string,
	productID uint64,
) (domain.Aggregate, error) {
	repo := fakeFromContext(ctx, r)
	for _, value := range repo.subscriptions {
		live := value.Status == models.SubscriptionPending ||
			value.Status == models.SubscriptionActive ||
			value.Status == models.SubscriptionPastDue
		if value.ProductID == productID && live {
			return value.Copy(), nil
		}
	}
	return domain.Aggregate{}, ErrNotFound
}

func (r *fakeRepository) CreateSubscription(
	ctx context.Context,
	draft migration.SubscriptionDraft,
	target PlanTarget,
	decision ReviewDecision,
	_ Run,
) (uint64, error) {
	repo := fakeFromContext(ctx, r)
	id := repo.nextID
	repo.nextID++
	accessUntil := draft.AccessUntil
	repo.subscriptions[draft.SubscriptionKey] = domain.Aggregate{
		ID:                 id,
		SubscriptionKey:    draft.SubscriptionKey,
		SubjectKind:        models.SubscriptionSubjectProfile,
		SubjectID:          draft.SubjectID,
		ProductID:          target.ProductID,
		ProductCode:        target.ProductCode,
		PlanVersionID:      target.PlanVersionID,
		PlanCode:           target.PlanCode,
		PlanVersion:        target.PlanVersion,
		PlanTermsChecksum:  target.PlanTermsChecksum,
		PlanStatus:         target.PlanStatus,
		Status:             models.SubscriptionStatus(draft.Status),
		StartedAt:          draft.StartedAt,
		CurrentPeriodStart: draft.CurrentPeriodStart,
		CurrentPeriodEnd:   draft.CurrentPeriodEnd,
		AccessUntil:        &accessUntil,
		AutoRenew:          decision.AutoRenew,
		OrderReference:     draft.OrderReference,
		CreatedAt:          draft.CreatedAt,
		UpdatedAt:          draft.UpdatedAt,
	}
	return id, nil
}

func (r *fakeRepository) EnsureMigrationEvent(ctx context.Context, event Event) error {
	repo := fakeFromContext(ctx, r)
	if repo.failEvent != nil {
		return repo.failEvent
	}
	if _, exists := repo.events[event.SubscriptionID]; !exists {
		repo.events[event.SubscriptionID] = event
	}
	return nil
}

func (r *fakeRepository) MarkShadowRead(
	ctx context.Context,
	legacyMembershipID uint64,
	subscriptionID uint64,
	run Run,
) error {
	repo := fakeFromContext(ctx, r)
	if repo.candidate.Membership.ID != legacyMembershipID {
		return ErrNotFound
	}
	repo.candidate.Mapping.SubscriptionID = &subscriptionID
	repo.candidate.Mapping.Status = models.CatalogMigrationShadowRead
	repo.candidate.Mapping.LastRunID = run.RunID
	repo.candidate.Mapping.MigratedAt = &run.ExecutedAt
	return nil
}

func (r *fakeRepository) MarkBlocked(
	_ context.Context,
	_ uint64,
	_ Run,
	reason string,
) error {
	r.blockedReason = reason
	r.candidate.Mapping.Status = models.CatalogMigrationBlocked
	return nil
}

func TestReviewStoresExplicitEvidence(t *testing.T) {
	repo := newFakeRepository(t)
	repo.candidate.Mapping = Mapping{
		LegacyMembershipID: 9,
		Status:             models.CatalogMigrationPendingReview,
	}
	service := NewService(repo)
	err := service.Review(context.Background(), ReviewRequest{
		LegacyMembershipID: 9,
		ReviewedBy:         4,
		ReviewedAt:         time.Now().UTC(),
		Decision:           validDecision(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !repo.reviewed || repo.candidate.Mapping.Status != models.CatalogMigrationReviewed {
		t.Fatalf("mapping = %+v", repo.candidate.Mapping)
	}
}

func TestExecuteCreatesSubscriptionEventAndShadowMapping(t *testing.T) {
	repo := newFakeRepository(t)
	result, err := NewService(repo).Execute(context.Background(), 9, validRun())
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != ResultMigrated || result.SubscriptionID == 0 {
		t.Fatalf("result = %+v", result)
	}
	if repo.candidate.Mapping.Status != models.CatalogMigrationShadowRead {
		t.Fatalf("mapping status = %q", repo.candidate.Mapping.Status)
	}
	if len(repo.events) != 1 {
		t.Fatalf("events = %d", len(repo.events))
	}
}

func TestExecuteIsIdempotentAfterSuccessfulRun(t *testing.T) {
	repo := newFakeRepository(t)
	service := NewService(repo)
	first, err := service.Execute(context.Background(), 9, validRun())
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.Execute(context.Background(), 9, validRun())
	if err != nil {
		t.Fatal(err)
	}
	if second.Status != ResultIdempotent || second.SubscriptionID != first.SubscriptionID {
		t.Fatalf("second = %+v", second)
	}
	if len(repo.subscriptions) != 1 || len(repo.events) != 1 {
		t.Fatalf("subscriptions=%d events=%d", len(repo.subscriptions), len(repo.events))
	}
}

func TestExecuteBlocksChangedLegacyMembership(t *testing.T) {
	repo := newFakeRepository(t)
	repo.candidate.Membership.OrderID++
	result, err := NewService(repo).Execute(context.Background(), 9, validRun())
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != ResultBlocked || repo.blockedReason == "" {
		t.Fatalf("result=%+v reason=%q", result, repo.blockedReason)
	}
	if len(repo.subscriptions) != 0 || len(repo.events) != 0 {
		t.Fatal("blocked candidate wrote migration state")
	}
}

func TestExecuteRollsBackWhenEventWriteFails(t *testing.T) {
	repo := newFakeRepository(t)
	repo.failEvent = errors.New("event storage failed")
	result, err := NewService(repo).Execute(context.Background(), 9, validRun())
	if err == nil || result.Status != ResultFailed {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	if len(repo.subscriptions) != 0 || len(repo.events) != 0 {
		t.Fatal("failed transaction retained partial state")
	}
	if repo.candidate.Mapping.Status != models.CatalogMigrationReviewed {
		t.Fatalf("mapping status = %q", repo.candidate.Mapping.Status)
	}
}

func TestExecuteBlocksConflictingLiveSubscription(t *testing.T) {
	repo := newFakeRepository(t)
	draft, err := migration.BuildSubscriptionDraft(
		toLegacySnapshot(repo.candidate.Membership),
		migration.MembershipDecision{
			ProductCode:       repo.candidate.Target.ProductCode,
			PlanCode:          repo.candidate.Target.PlanCode,
			PlanVersion:       repo.candidate.Target.PlanVersion,
			PlanTermsChecksum: repo.candidate.Target.PlanTermsChecksum,
			PeriodHandling:    PeriodPaidAtToFinishedDate,
		},
		validRun().AsOf,
	)
	if err != nil {
		t.Fatal(err)
	}
	accessUntil := draft.AccessUntil
	repo.subscriptions["sub_existing"] = domain.Aggregate{
		ID:                 18,
		SubscriptionKey:    "sub_existing",
		SubjectKind:        models.SubscriptionSubjectProfile,
		SubjectID:          draft.SubjectID,
		ProductID:          repo.candidate.Target.ProductID,
		ProductCode:        repo.candidate.Target.ProductCode,
		PlanVersionID:      repo.candidate.Target.PlanVersionID,
		PlanCode:           repo.candidate.Target.PlanCode,
		PlanVersion:        repo.candidate.Target.PlanVersion,
		PlanTermsChecksum:  repo.candidate.Target.PlanTermsChecksum,
		PlanStatus:         repo.candidate.Target.PlanStatus,
		Status:             models.SubscriptionActive,
		StartedAt:          draft.StartedAt,
		CurrentPeriodStart: draft.CurrentPeriodStart,
		CurrentPeriodEnd:   draft.CurrentPeriodEnd,
		AccessUntil:        &accessUntil,
		OrderReference:     draft.OrderReference,
		CreatedAt:          draft.CreatedAt,
		UpdatedAt:          draft.UpdatedAt,
	}

	result, err := NewService(repo).Execute(context.Background(), 9, validRun())
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != ResultBlocked || !strings.Contains(result.Detail, "live subscription") {
		t.Fatalf("result=%+v", result)
	}
	if repo.candidate.Mapping.Status != models.CatalogMigrationBlocked {
		t.Fatalf("mapping status=%q", repo.candidate.Mapping.Status)
	}
	if len(repo.subscriptions) != 1 || len(repo.events) != 0 {
		t.Fatalf("subscriptions=%d events=%d", len(repo.subscriptions), len(repo.events))
	}
}

func TestExecuteRerunBlocksDriftAfterMigration(t *testing.T) {
	repo := newFakeRepository(t)
	service := NewService(repo)
	if _, err := service.Execute(context.Background(), 9, validRun()); err != nil {
		t.Fatal(err)
	}
	repo.candidate.Membership.OrderID++
	result, err := service.Execute(context.Background(), 9, validRun())
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != ResultBlocked || repo.candidate.Mapping.Status != models.CatalogMigrationBlocked {
		t.Fatalf("result=%+v mapping=%+v", result, repo.candidate.Mapping)
	}
	if len(repo.subscriptions) != 1 || len(repo.events) != 1 {
		t.Fatal("drift rerun changed durable migration evidence")
	}
}
