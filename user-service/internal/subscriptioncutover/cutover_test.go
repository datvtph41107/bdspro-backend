package subscriptioncutover

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"testing"
	"time"

	"common/productcatalog"
	"user/internal/models"
	"user/internal/subscriptionbackfill"
	"user/internal/subscriptiondomain"
	"user/internal/subscriptionshadow"
)

type fakeRepository struct {
	candidate        subscriptionbackfill.Candidate
	cutoverEvents    []Event
	rollbackEvents   []Event
	cutoverChecksum  string
	rollbackChecksum string
	failEvent        error
}

func (f *fakeRepository) clone() *fakeRepository {
	copied := *f
	copied.candidate = f.candidate
	copied.candidate.Mapping.Decision = append(json.RawMessage(nil), f.candidate.Mapping.Decision...)
	if f.candidate.ExistingSubscription != nil {
		value := f.candidate.ExistingSubscription.Copy()
		copied.candidate.ExistingSubscription = &value
	}
	copied.cutoverEvents = append([]Event(nil), f.cutoverEvents...)
	copied.rollbackEvents = append([]Event(nil), f.rollbackEvents...)
	return &copied
}

func (f *fakeRepository) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	working := f.clone()
	if err := fn(context.WithValue(ctx, fakeKey{}, working)); err != nil {
		return err
	}
	*f = *working
	return nil
}

type fakeKey struct{}

func fromContext(ctx context.Context, fallback *fakeRepository) *fakeRepository {
	if value, ok := ctx.Value(fakeKey{}).(*fakeRepository); ok {
		return value
	}
	return fallback
}

func (f *fakeRepository) LoadCandidate(ctx context.Context, id uint64, _ bool) (subscriptionbackfill.Candidate, error) {
	repo := fromContext(ctx, f)
	if repo.candidate.Membership.ID != id {
		return subscriptionbackfill.Candidate{}, subscriptionbackfill.ErrNotFound
	}
	return repo.candidate, nil
}
func (f *fakeRepository) EnsureCutoverEvent(ctx context.Context, event Event) error {
	repo := fromContext(ctx, f)
	if repo.failEvent != nil {
		return repo.failEvent
	}
	repo.cutoverEvents = append(repo.cutoverEvents, event)
	return nil
}
func (f *fakeRepository) MarkCutover(ctx context.Context, entry Entry, checksum string, _ Run) error {
	repo := fromContext(ctx, f)
	if repo.candidate.Mapping.Status != models.CatalogMigrationShadowRead {
		return ErrConcurrentChange
	}
	repo.candidate.Mapping.Status = models.CatalogMigrationCutover
	repo.candidate.Mapping.CutoverEvidenceChecksum = checksum
	repo.cutoverChecksum = checksum
	return nil
}
func (f *fakeRepository) EnsureRollbackEvent(ctx context.Context, event Event) error {
	repo := fromContext(ctx, f)
	if repo.failEvent != nil {
		return repo.failEvent
	}
	repo.rollbackEvents = append(repo.rollbackEvents, event)
	return nil
}
func (f *fakeRepository) MarkRolledBack(ctx context.Context, entry Entry, checksum string, _ Run) error {
	repo := fromContext(ctx, f)
	if repo.candidate.Mapping.Status != models.CatalogMigrationCutover {
		return ErrConcurrentChange
	}
	repo.candidate.Mapping.Status = models.CatalogMigrationRolledBack
	repo.candidate.Mapping.RollbackEvidenceChecksum = checksum
	repo.rollbackChecksum = checksum
	return nil
}

type populationRepository struct {
	candidates map[uint64]subscriptionbackfill.Candidate
}

func (r *populationRepository) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}
func (r *populationRepository) LoadCandidate(_ context.Context, id uint64, _ bool) (subscriptionbackfill.Candidate, error) {
	candidate, ok := r.candidates[id]
	if !ok {
		return subscriptionbackfill.Candidate{}, subscriptionbackfill.ErrNotFound
	}
	return candidate, nil
}
func (*populationRepository) EnsureCutoverEvent(context.Context, Event) error          { return nil }
func (*populationRepository) MarkCutover(context.Context, Entry, string, Run) error    { return nil }
func (*populationRepository) EnsureRollbackEvent(context.Context, Event) error         { return nil }
func (*populationRepository) MarkRolledBack(context.Context, Entry, string, Run) error { return nil }

func populationFixture(count int) (*populationRepository, subscriptionbackfill.ReadinessResult) {
	base := fixtureRepository().candidate
	repo := &populationRepository{candidates: make(map[uint64]subscriptionbackfill.Candidate, count)}
	readiness := strictReadiness()
	readiness.Report.TotalLegacy = count
	readiness.Report.EligibleLegacy = count
	readiness.Report.Observed = count
	readiness.Report.Matched = count
	readiness.Report.RollbackEvidence = count
	readiness.Report.Entries = make([]subscriptionshadow.Entry, 0, count)
	readiness.Report.Statuses = map[string]int{string(subscriptionshadow.StatusMatched): count}
	for index := 1; index <= count; index++ {
		candidate := base
		id := uint64(index)
		profileID := uint64(1000 + index)
		subscriptionID := uint64(10000 + index)
		candidate.Membership.ID = id
		candidate.Membership.ProfileID = profileID
		candidate.Mapping.LegacyMembershipID = id
		candidate.Mapping.SubscriptionID = &subscriptionID
		aggregate := candidate.ExistingSubscription.Copy()
		aggregate.ID = subscriptionID
		aggregate.SubjectID = strconv.FormatUint(profileID, 10)
		aggregate.SubscriptionKey = "sub_legacy_membership_" + strconv.Itoa(index)
		candidate.ExistingSubscription = &aggregate
		repo.candidates[id] = candidate
		readiness.Report.Entries = append(readiness.Report.Entries, subscriptionshadow.Entry{
			LegacyMembershipID: id,
			SubscriptionID:     subscriptionID,
			MappingStatus:      models.CatalogMigrationShadowRead,
			PlanMappingStatus:  models.CatalogMigrationShadowRead,
			Status:             subscriptionshadow.StatusMatched,
		})
	}
	return repo, readiness
}

func fixtureRepository() *fakeRepository {
	paidAt := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	finishedAt := paidAt.AddDate(0, 1, 0)
	createdAt := paidAt.Add(-time.Hour)
	updatedAt := paidAt
	membership := models.Membership{ProfileID: 42, PlanID: 7, PaidAt: &paidAt, FinishedDate: &finishedAt, OrderID: 88}
	membership.ID = 9
	membership.CreatedAt = &createdAt
	membership.UpdatedAt = &updatedAt
	subscriptionID := uint64(19)
	reviewer := uint64(3)
	reviewedAt := paidAt.Add(-2 * time.Hour)
	aggregate := subscriptiondomain.Aggregate{
		ID: 19, SubscriptionKey: "sub_legacy_membership_9",
		SubjectKind: models.SubscriptionSubjectProfile, SubjectID: "42",
		ProductID: 3, ProductCode: "qhpro", PlanVersionID: 11,
		PlanCode: "qhpro.pro", PlanVersion: "1.0.0",
		PlanTermsChecksum: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		PlanStatus:        productcatalog.StatusActive, Status: models.SubscriptionActive,
		StartedAt: paidAt, CurrentPeriodStart: paidAt, CurrentPeriodEnd: finishedAt,
		AccessUntil: &finishedAt, OrderReference: "88", CreatedAt: createdAt, UpdatedAt: updatedAt,
	}
	return &fakeRepository{candidate: subscriptionbackfill.Candidate{
		Membership: membership,
		Mapping: subscriptionbackfill.Mapping{
			LegacyMembershipID: 9, SubscriptionID: &subscriptionID,
			SourceChecksum: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			Status:         models.CatalogMigrationShadowRead, ReviewedBy: &reviewer, ReviewedAt: &reviewedAt,
		},
		Target: subscriptionbackfill.PlanTarget{
			LegacyPlanID: 7, ProductID: 3, ProductCode: "qhpro", PlanVersionID: 11,
			PlanCode: "qhpro.pro", PlanVersion: "1.0.0",
			PlanTermsChecksum: aggregate.PlanTermsChecksum, PlanStatus: productcatalog.StatusActive,
			SubjectScope:  productcatalog.SubjectScopeProfile,
			MappingStatus: models.CatalogMigrationShadowRead,
		},
		ExistingSubscription: &aggregate,
	}}
}

func strictReadiness() subscriptionbackfill.ReadinessResult {
	return subscriptionbackfill.ReadinessResult{
		Ready:  true,
		Policy: subscriptionbackfill.StrictReadinessPolicy(),
		Report: subscriptionbackfill.CoverageReport{
			AsOf:        time.Date(2026, 8, 4, 0, 0, 0, 0, time.UTC),
			TotalLegacy: 1, EligibleLegacy: 1, Observed: 1, Matched: 1,
			RollbackEvidence: 1,
			Statuses:         map[string]int{string(subscriptionshadow.StatusMatched): 1},
			Entries: []subscriptionshadow.Entry{{
				LegacyMembershipID: 9, SubscriptionID: 19,
				MappingStatus:     models.CatalogMigrationShadowRead,
				PlanMappingStatus: models.CatalogMigrationShadowRead,
				Status:            subscriptionshadow.StatusMatched,
			}},
			Checksum: "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
		},
	}
}

func validRun(reason string) Run {
	return Run{RunID: "cutover_20260804", ActorID: "catalog-operator", RequestID: "req_cutover", OperationID: "op_cutover", ExecutedAt: time.Date(2026, 8, 4, 1, 0, 0, 0, time.UTC), Reason: reason}
}

func TestPrepareBindsConcreteSortedEvidence(t *testing.T) {
	repo := fixtureRepository()
	manifest, err := NewService(repo).Prepare(context.Background(), strictReadiness(), time.Date(2026, 8, 4, 0, 30, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if len(manifest.Entries) != 1 || manifest.Entries[0].LegacyMembershipID != 9 || !validChecksum(manifest.EvidenceChecksum) {
		t.Fatalf("manifest=%+v", manifest)
	}
	changed := manifest
	changed.Entries = append([]Entry(nil), manifest.Entries...)
	changed.Entries[0].SubscriptionID++
	if changed.Validate() == nil {
		t.Fatal("tampered manifest passed validation")
	}
}

func TestActivateIsAtomicAndMarksCutover(t *testing.T) {
	repo := fixtureRepository()
	service := NewService(repo)
	manifest, err := service.Prepare(context.Background(), strictReadiness(), time.Date(2026, 8, 4, 0, 30, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if err := service.Activate(context.Background(), manifest, validRun("Approved canary evidence")); err != nil {
		t.Fatal(err)
	}
	if repo.candidate.Mapping.Status != models.CatalogMigrationCutover || repo.cutoverChecksum != manifest.EvidenceChecksum || len(repo.cutoverEvents) != 1 {
		t.Fatalf("repo=%+v", repo)
	}
}

func TestActivateRollsBackWhenEventFails(t *testing.T) {
	repo := fixtureRepository()
	service := NewService(repo)
	manifest, _ := service.Prepare(context.Background(), strictReadiness(), time.Now().UTC())
	repo.failEvent = errors.New("event unavailable")
	if err := service.Activate(context.Background(), manifest, validRun("Approved")); err == nil {
		t.Fatal("expected activation failure")
	}
	if repo.candidate.Mapping.Status != models.CatalogMigrationShadowRead || repo.cutoverChecksum != "" {
		t.Fatal("partial cutover escaped transaction")
	}
}

func TestRollbackRequiresOriginalCutoverEvidence(t *testing.T) {
	repo := fixtureRepository()
	service := NewService(repo)
	manifest, _ := service.Prepare(context.Background(), strictReadiness(), time.Now().UTC())
	if err := service.Activate(context.Background(), manifest, validRun("Approved")); err != nil {
		t.Fatal(err)
	}
	rollback, err := BuildRollbackManifest(manifest, "Rollback drill after catalog read failure", time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if err := service.Rollback(context.Background(), rollback, validRun("Rollback drill after catalog read failure")); err != nil {
		t.Fatal(err)
	}
	if repo.candidate.Mapping.Status != models.CatalogMigrationRolledBack || repo.rollbackChecksum != rollback.EvidenceChecksum || len(repo.rollbackEvents) != 1 {
		t.Fatalf("repo=%+v", repo)
	}
}

func TestPrepareRejectsAggregateOnlyReadiness(t *testing.T) {
	readiness := strictReadiness()
	readiness.Report.Entries = nil
	if _, err := NewService(fixtureRepository()).Prepare(context.Background(), readiness, time.Now().UTC()); !errors.Is(err, ErrEvidenceInvalid) {
		t.Fatalf("err=%v", err)
	}
}

func TestPreparePlanCreatesBoundedCompleteBatches(t *testing.T) {
	repo, readiness := populationFixture(1001)
	plan, err := NewService(repo).PreparePlan(
		context.Background(), readiness,
		time.Date(2026, 8, 4, 0, 30, 0, 0, time.UTC),
		250,
	)
	if err != nil {
		t.Fatal(err)
	}
	if plan.TotalEntries != 1001 || len(plan.Batches) != 5 ||
		len(plan.Batches[0].Entries) != 250 || len(plan.Batches[4].Entries) != 1 {
		t.Fatalf("plan shape = total:%d batches:%d first:%d last:%d", plan.TotalEntries, len(plan.Batches), len(plan.Batches[0].Entries), len(plan.Batches[4].Entries))
	}
	if err := plan.Validate(); err != nil {
		t.Fatal(err)
	}
	tampered := plan
	tampered.Batches = append([]Manifest(nil), plan.Batches...)
	tampered.Batches[1] = plan.Batches[1]
	tampered.Batches[1].PopulationChecksum = "dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
	if tampered.Validate() == nil {
		t.Fatal("tampered batch plan passed validation")
	}
}

func TestActivateCanResumeSameEvidence(t *testing.T) {
	repo := fixtureRepository()
	service := NewService(repo)
	manifest, err := service.Prepare(context.Background(), strictReadiness(), time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	run := validRun("Approved")
	if err := service.Activate(context.Background(), manifest, run); err != nil {
		t.Fatal(err)
	}
	if err := service.Activate(context.Background(), manifest, run); err != nil {
		t.Fatalf("idempotent resume failed: %v", err)
	}
	if len(repo.cutoverEvents) != 1 {
		t.Fatalf("cutover events = %d", len(repo.cutoverEvents))
	}
}
