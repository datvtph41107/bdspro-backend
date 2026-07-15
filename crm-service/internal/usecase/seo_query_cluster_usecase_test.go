package usecase

import (
	"context"
	"testing"
	"time"

	seo_strategy "crm/internal/domain/seo_strategy"
	"crm/internal/dto"
	"crm/internal/repo"
)

type fakeQueryClusterRepo struct {
	items   map[uint64]*seo_strategy.SeoQueryCluster
	members map[uint64][]seo_strategy.LinkedQueryClusterMember
	nextID  uint64
}

func newFakeQueryClusterRepo() *fakeQueryClusterRepo {
	return &fakeQueryClusterRepo{
		items:   map[uint64]*seo_strategy.SeoQueryCluster{},
		members: map[uint64][]seo_strategy.LinkedQueryClusterMember{},
		nextID:  1,
	}
}

func (r *fakeQueryClusterRepo) Create(_ context.Context, item *seo_strategy.SeoQueryCluster) error {
	for _, existing := range r.items {
		if existing.UserNeedID == item.UserNeedID && existing.Code == item.Code {
			return repo.ErrSeoQueryClusterConflict
		}
	}
	item.ID = r.nextID
	r.nextID++
	copyItem := *item
	r.items[item.ID] = &copyItem
	return nil
}

func (r *fakeQueryClusterRepo) GetByID(_ context.Context, id uint64) (*seo_strategy.SeoQueryCluster, error) {
	item, ok := r.items[id]
	if !ok {
		return nil, repo.ErrSeoQueryClusterNotFound
	}
	copyItem := *item
	return &copyItem, nil
}

func (r *fakeQueryClusterRepo) GetByCode(_ context.Context, userNeedID uint64, code string) (*seo_strategy.SeoQueryCluster, error) {
	for _, item := range r.items {
		if item.UserNeedID == userNeedID && item.Code == seo_strategy.NormalizeStrategyCode(code) {
			copyItem := *item
			return &copyItem, nil
		}
	}
	return nil, repo.ErrSeoQueryClusterNotFound
}

func (r *fakeQueryClusterRepo) List(_ context.Context, filter repo.SeoQueryClusterFilter) ([]seo_strategy.SeoQueryCluster, int64, error) {
	items := make([]seo_strategy.SeoQueryCluster, 0)
	for _, item := range r.items {
		if filter.IntentHypothesisID > 0 && item.IntentHypothesisID != filter.IntentHypothesisID {
			continue
		}
		items = append(items, *item)
	}
	return items, int64(len(items)), nil
}

func (r *fakeQueryClusterRepo) UpdateDefinition(_ context.Context, item *seo_strategy.SeoQueryCluster, expectedVersion uint32) error {
	current, ok := r.items[item.ID]
	if !ok {
		return repo.ErrSeoQueryClusterNotFound
	}
	if current.LockVersion != expectedVersion {
		return repo.ErrSeoQueryClusterConflict
	}
	copyItem := *item
	copyItem.LockVersion = expectedVersion + 1
	r.items[item.ID] = &copyItem
	return nil
}

func (r *fakeQueryClusterRepo) UpdateStatus(_ context.Context, id uint64, status seo_strategy.QueryClusterStatus, expectedVersion uint32) error {
	current, ok := r.items[id]
	if !ok {
		return repo.ErrSeoQueryClusterNotFound
	}
	if current.LockVersion != expectedVersion {
		return repo.ErrSeoQueryClusterConflict
	}
	current.Status = status
	current.LockVersion++
	return nil
}

func (r *fakeQueryClusterRepo) UpsertMember(_ context.Context, member *seo_strategy.SeoQueryClusterMember, expectedVersion uint32) error {
	cluster, ok := r.items[member.ClusterID]
	if !ok {
		return repo.ErrSeoQueryClusterNotFound
	}
	if cluster.LockVersion != expectedVersion {
		return repo.ErrSeoQueryClusterConflict
	}
	items := r.members[member.ClusterID]
	for i := range items {
		if items[i].Evidence.ID == member.QueryEvidenceID {
			items[i].Member = *member
			r.members[member.ClusterID] = items
			cluster.LockVersion++
			return nil
		}
	}
	r.members[member.ClusterID] = append(items, seo_strategy.LinkedQueryClusterMember{Member: *member})
	cluster.LockVersion++
	return nil
}

func (r *fakeQueryClusterRepo) RemoveMember(_ context.Context, clusterID uint64, queryEvidenceID uint64, expectedVersion uint32) error {
	cluster, ok := r.items[clusterID]
	if !ok {
		return repo.ErrSeoQueryClusterNotFound
	}
	if cluster.LockVersion != expectedVersion {
		return repo.ErrSeoQueryClusterConflict
	}
	items := r.members[clusterID]
	for i := range items {
		if items[i].Evidence.ID == queryEvidenceID && items[i].Member.Status == seo_strategy.QueryClusterMemberStatusActive {
			items[i].Member.Status = seo_strategy.QueryClusterMemberStatusRemoved
			r.members[clusterID] = items
			cluster.LockVersion++
			return nil
		}
	}
	return repo.ErrSeoQueryClusterMemberNotFound
}

func (r *fakeQueryClusterRepo) ListMembers(_ context.Context, clusterID uint64) ([]seo_strategy.LinkedQueryClusterMember, error) {
	items := r.members[clusterID]
	result := make([]seo_strategy.LinkedQueryClusterMember, len(items))
	copy(result, items)
	return result, nil
}

type fakeQueryClusterIntentReader struct {
	item  *seo_strategy.SeoIntentHypothesis
	links []seo_strategy.LinkedIntentQueryEvidence
	err   error
}

func (r *fakeQueryClusterIntentReader) GetByID(_ context.Context, id uint64) (*seo_strategy.SeoIntentHypothesis, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.item == nil || r.item.ID != id {
		return nil, repo.ErrSeoIntentHypothesisNotFound
	}
	copyItem := *r.item
	return &copyItem, nil
}

func (r *fakeQueryClusterIntentReader) ListLinkedEvidence(_ context.Context, intentID uint64) ([]seo_strategy.LinkedIntentQueryEvidence, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.item == nil || r.item.ID != intentID {
		return nil, repo.ErrSeoIntentHypothesisNotFound
	}
	result := make([]seo_strategy.LinkedIntentQueryEvidence, len(r.links))
	copy(result, r.links)
	return result, nil
}

type fakeQueryClusterUserNeedReader struct {
	item *seo_strategy.SeoUserNeed
	err  error
}

func (r *fakeQueryClusterUserNeedReader) GetByID(_ context.Context, id uint64) (*seo_strategy.SeoUserNeed, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.item == nil || r.item.ID != id {
		return nil, repo.ErrSeoUserNeedNotFound
	}
	copyItem := *r.item
	return &copyItem, nil
}

type fakeQueryClusterEvidenceReader struct {
	items map[uint64]*seo_strategy.SeoQueryEvidence
}

func (r *fakeQueryClusterEvidenceReader) GetByID(_ context.Context, id uint64) (*seo_strategy.SeoQueryEvidence, error) {
	item, ok := r.items[id]
	if !ok {
		return nil, repo.ErrSeoQueryEvidenceNotFound
	}
	copyItem := *item
	return &copyItem, nil
}

func TestSeoQueryClusterUsecaseCreateRequiresActiveIntent(t *testing.T) {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	u := &SeoQueryClusterUsecase{
		repository: newFakeQueryClusterRepo(),
		intentReader: &fakeQueryClusterIntentReader{item: &seo_strategy.SeoIntentHypothesis{
			ID: 30, ObjectiveID: 10, UserNeedID: 20, Status: seo_strategy.IntentHypothesisStatusDraft,
		}},
		userNeedReader: &fakeQueryClusterUserNeedReader{item: &seo_strategy.SeoUserNeed{
			ID: 20, ObjectiveID: 10, Status: seo_strategy.UserNeedStatusActive,
		}},
		queryEvidenceReader: &fakeQueryClusterEvidenceReader{items: map[uint64]*seo_strategy.SeoQueryEvidence{}},
		now:                 func() time.Time { return now },
	}
	_, err := u.Create(context.Background(), &dto.CreateSeoQueryClusterRequest{
		IntentHypothesisID: 30, Code: "parcel verify", Name: "Parcel verify",
	})
	if err == nil {
		t.Fatal("expected inactive intent error")
	}
}

func TestSeoQueryClusterUsecaseCreate(t *testing.T) {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	repository := newFakeQueryClusterRepo()
	u := queryClusterUsecaseFixture(now, repository)
	result, err := u.Create(context.Background(), &dto.CreateSeoQueryClusterRequest{
		IntentHypothesisID: 30, Code: "parcel verify", Name: "Parcel verify",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Cluster.ID == 0 || result.Cluster.Code != "parcel_verify" || result.Cluster.ObjectiveID != 10 {
		t.Fatalf("unexpected cluster: %+v", result.Cluster)
	}
}

func TestSeoQueryClusterUsecaseUpsertMemberRequiresSupportingIntentRelation(t *testing.T) {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	repository := newFakeQueryClusterRepo()
	repository.items[1] = queryClusterDraftForUsecase(now)
	evidence := queryClusterEvidenceForUsecase(101, now)
	u := queryClusterUsecaseFixture(now, repository)
	u.queryEvidenceReader = &fakeQueryClusterEvidenceReader{items: map[uint64]*seo_strategy.SeoQueryEvidence{101: &evidence}}
	u.intentReader = &fakeQueryClusterIntentReader{
		item: &seo_strategy.SeoIntentHypothesis{ID: 30, ObjectiveID: 10, UserNeedID: 20, Status: seo_strategy.IntentHypothesisStatusActive},
		links: []seo_strategy.LinkedIntentQueryEvidence{{
			Link: seo_strategy.SeoIntentEvidenceLink{RelationType: seo_strategy.IntentEvidenceContext}, Evidence: evidence,
		}},
	}
	_, err := u.UpsertMember(context.Background(), &dto.UpsertSeoQueryClusterMemberRequest{
		ClusterID: 1, QueryEvidenceID: 101, ExpectedVersion: 1,
		MembershipType: "included", Rationale: "same answer",
	})
	if err == nil {
		t.Fatal("expected non-supporting relation error")
	}
}

func TestSeoQueryClusterUsecaseUpsertMemberIncrementsVersion(t *testing.T) {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	repository := newFakeQueryClusterRepo()
	repository.items[1] = queryClusterDraftForUsecase(now)
	evidence := queryClusterEvidenceForUsecase(101, now)
	u := queryClusterUsecaseFixture(now, repository)
	u.queryEvidenceReader = &fakeQueryClusterEvidenceReader{items: map[uint64]*seo_strategy.SeoQueryEvidence{101: &evidence}}
	u.intentReader = &fakeQueryClusterIntentReader{
		item: &seo_strategy.SeoIntentHypothesis{ID: 30, ObjectiveID: 10, UserNeedID: 20, Status: seo_strategy.IntentHypothesisStatusActive},
		links: []seo_strategy.LinkedIntentQueryEvidence{{
			Link: seo_strategy.SeoIntentEvidenceLink{RelationType: seo_strategy.IntentEvidenceSupports}, Evidence: evidence,
		}},
	}
	result, err := u.UpsertMember(context.Background(), &dto.UpsertSeoQueryClusterMemberRequest{
		ClusterID: 1, QueryEvidenceID: 101, ExpectedVersion: 1,
		MembershipType: "primary",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Cluster.LockVersion != 2 || len(result.Members) != 1 {
		t.Fatalf("unexpected result: cluster=%+v members=%+v", result.Cluster, result.Members)
	}
}

func TestSeoQueryClusterUsecaseActivateUsesAnswerEquivalenceReadiness(t *testing.T) {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	review := now.Add(30 * 24 * time.Hour)
	repository := newFakeQueryClusterRepo()
	cluster := queryClusterDraftForUsecase(now)
	cluster.SearchPromise = "Source-backed parcel planning answer"
	cluster.InclusionRule = "Same parcel planning answer"
	cluster.ExclusionRule = "Different answer or non-planning task"
	cluster.OwnerCode = "seo"
	cluster.ReviewDueAt = &review
	cluster.Confidence = seo_strategy.ObjectiveConfidenceSupported
	repository.items[1] = cluster

	primary := queryClusterEvidenceForUsecase(101, now)
	included := queryClusterEvidenceForUsecase(102, now)
	repository.members[1] = []seo_strategy.LinkedQueryClusterMember{
		{Member: seo_strategy.SeoQueryClusterMember{ClusterID: 1, QueryEvidenceID: 101, MembershipType: seo_strategy.QueryClusterMembershipPrimary, Status: seo_strategy.QueryClusterMemberStatusActive}, Evidence: primary},
		{Member: seo_strategy.SeoQueryClusterMember{ClusterID: 1, QueryEvidenceID: 102, MembershipType: seo_strategy.QueryClusterMembershipIncluded, Rationale: "same answer", Status: seo_strategy.QueryClusterMemberStatusActive}, Evidence: included},
	}
	u := queryClusterUsecaseFixture(now, repository)
	u.intentReader = &fakeQueryClusterIntentReader{
		item: &seo_strategy.SeoIntentHypothesis{ID: 30, ObjectiveID: 10, UserNeedID: 20, Status: seo_strategy.IntentHypothesisStatusActive},
		links: []seo_strategy.LinkedIntentQueryEvidence{
			{Link: seo_strategy.SeoIntentEvidenceLink{RelationType: seo_strategy.IntentEvidenceSupports}, Evidence: primary},
			{Link: seo_strategy.SeoIntentEvidenceLink{RelationType: seo_strategy.IntentEvidenceSupports}, Evidence: included},
		},
	}

	result, err := u.Activate(context.Background(), &dto.SeoQueryClusterCommandRequest{ID: 1, ExpectedVersion: 1})
	if err != nil {
		t.Fatal(err)
	}
	if result.Cluster.Status != seo_strategy.QueryClusterStatusActive || result.Cluster.LockVersion != 2 {
		t.Fatalf("unexpected lifecycle result: %+v", result.Cluster)
	}
	_, err = u.Pause(context.Background(), &dto.SeoQueryClusterCommandRequest{ID: 1, ExpectedVersion: 1})
	if err == nil {
		t.Fatal("expected stale version error")
	}
}

func queryClusterUsecaseFixture(now time.Time, repository *fakeQueryClusterRepo) *SeoQueryClusterUsecase {
	return &SeoQueryClusterUsecase{
		repository: repository,
		intentReader: &fakeQueryClusterIntentReader{item: &seo_strategy.SeoIntentHypothesis{
			ID: 30, ObjectiveID: 10, UserNeedID: 20, Status: seo_strategy.IntentHypothesisStatusActive,
		}},
		userNeedReader: &fakeQueryClusterUserNeedReader{item: &seo_strategy.SeoUserNeed{
			ID: 20, ObjectiveID: 10, Status: seo_strategy.UserNeedStatusActive,
		}},
		queryEvidenceReader: &fakeQueryClusterEvidenceReader{items: map[uint64]*seo_strategy.SeoQueryEvidence{}},
		now:                 func() time.Time { return now },
	}
}

func queryClusterDraftForUsecase(now time.Time) *seo_strategy.SeoQueryCluster {
	return &seo_strategy.SeoQueryCluster{
		ID: 1, ObjectiveID: 10, UserNeedID: 20, IntentHypothesisID: 30,
		Code: "parcel_verify", Name: "Parcel verify", Status: seo_strategy.QueryClusterStatusDraft,
		Confidence: seo_strategy.ObjectiveConfidenceExploratory, LockVersion: 1,
		CreatedAt: now, UpdatedAt: now,
	}
}

func queryClusterEvidenceForUsecase(id uint64, now time.Time) seo_strategy.SeoQueryEvidence {
	return seo_strategy.SeoQueryEvidence{
		ID: id, ObjectiveID: 10, UserNeedID: 20,
		RawQuery: "parcel planning", NormalizedQuery: "parcel planning",
		ObservationType: seo_strategy.QueryObservationObserved,
		Confidence:      seo_strategy.ObjectiveConfidenceValidated,
		Status:          seo_strategy.QueryEvidenceStatusActive,
		CreatedAt:       now, UpdatedAt: now,
	}
}
