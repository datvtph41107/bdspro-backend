package usecase

import (
	"context"
	"testing"
	"time"

	seo_strategy "crm/internal/domain/seo_strategy"
	"crm/internal/dto"
	"crm/internal/repo"
)

type fakeIntentHypothesisRepo struct {
	items  map[uint64]*seo_strategy.SeoIntentHypothesis
	links  map[uint64][]seo_strategy.LinkedIntentQueryEvidence
	nextID uint64
}

func newFakeIntentHypothesisRepo() *fakeIntentHypothesisRepo {
	return &fakeIntentHypothesisRepo{
		items:  map[uint64]*seo_strategy.SeoIntentHypothesis{},
		links:  map[uint64][]seo_strategy.LinkedIntentQueryEvidence{},
		nextID: 1,
	}
}

func (r *fakeIntentHypothesisRepo) Create(_ context.Context, item *seo_strategy.SeoIntentHypothesis) error {
	for _, existing := range r.items {
		if existing.UserNeedID == item.UserNeedID && existing.Code == item.Code {
			return repo.ErrSeoIntentHypothesisConflict
		}
	}
	item.ID = r.nextID
	r.nextID++
	copyItem := *item
	r.items[item.ID] = &copyItem
	return nil
}

func (r *fakeIntentHypothesisRepo) GetByID(_ context.Context, id uint64) (*seo_strategy.SeoIntentHypothesis, error) {
	item, ok := r.items[id]
	if !ok {
		return nil, repo.ErrSeoIntentHypothesisNotFound
	}
	copyItem := *item
	return &copyItem, nil
}

func (r *fakeIntentHypothesisRepo) GetByCode(_ context.Context, userNeedID uint64, code string) (*seo_strategy.SeoIntentHypothesis, error) {
	for _, item := range r.items {
		if item.UserNeedID == userNeedID && item.Code == seo_strategy.NormalizeStrategyCode(code) {
			copyItem := *item
			return &copyItem, nil
		}
	}
	return nil, repo.ErrSeoIntentHypothesisNotFound
}

func (r *fakeIntentHypothesisRepo) List(_ context.Context, filter repo.SeoIntentHypothesisFilter) ([]seo_strategy.SeoIntentHypothesis, int64, error) {
	items := make([]seo_strategy.SeoIntentHypothesis, 0)
	for _, item := range r.items {
		if filter.UserNeedID > 0 && item.UserNeedID != filter.UserNeedID {
			continue
		}
		items = append(items, *item)
	}
	return items, int64(len(items)), nil
}

func (r *fakeIntentHypothesisRepo) UpdateDefinition(_ context.Context, item *seo_strategy.SeoIntentHypothesis, expectedVersion uint32) error {
	current, ok := r.items[item.ID]
	if !ok {
		return repo.ErrSeoIntentHypothesisNotFound
	}
	if current.LockVersion != expectedVersion {
		return repo.ErrSeoIntentHypothesisConflict
	}
	copyItem := *item
	copyItem.LockVersion = expectedVersion + 1
	r.items[item.ID] = &copyItem
	return nil
}

func (r *fakeIntentHypothesisRepo) UpdateStatus(_ context.Context, id uint64, status seo_strategy.IntentHypothesisStatus, expectedVersion uint32) error {
	current, ok := r.items[id]
	if !ok {
		return repo.ErrSeoIntentHypothesisNotFound
	}
	if current.LockVersion != expectedVersion {
		return repo.ErrSeoIntentHypothesisConflict
	}
	current.Status = status
	current.LockVersion++
	return nil
}

func (r *fakeIntentHypothesisRepo) LinkEvidence(_ context.Context, link *seo_strategy.SeoIntentEvidenceLink) error {
	for _, existing := range r.links[link.IntentHypothesisID] {
		if existing.Evidence.ID == link.QueryEvidenceID {
			return repo.ErrSeoIntentEvidenceAlreadyLinked
		}
	}
	return nil
}

func (r *fakeIntentHypothesisRepo) UnlinkEvidence(_ context.Context, intentID uint64, evidenceID uint64) error {
	items := r.links[intentID]
	for i := range items {
		if items[i].Evidence.ID == evidenceID {
			r.links[intentID] = append(items[:i], items[i+1:]...)
			return nil
		}
	}
	return repo.ErrSeoIntentEvidenceNotLinked
}

func (r *fakeIntentHypothesisRepo) ListLinkedEvidence(_ context.Context, intentID uint64) ([]seo_strategy.LinkedIntentQueryEvidence, error) {
	return append([]seo_strategy.LinkedIntentQueryEvidence(nil), r.links[intentID]...), nil
}

type fakeIntentUserNeedReader struct {
	item *seo_strategy.SeoUserNeed
	err  error
}

func (r *fakeIntentUserNeedReader) GetByID(_ context.Context, id uint64) (*seo_strategy.SeoUserNeed, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.item == nil || r.item.ID != id {
		return nil, repo.ErrSeoUserNeedNotFound
	}
	copyItem := *r.item
	return &copyItem, nil
}

type fakeIntentQueryEvidenceReader struct {
	items map[uint64]*seo_strategy.SeoQueryEvidence
}

func (r *fakeIntentQueryEvidenceReader) GetByID(_ context.Context, id uint64) (*seo_strategy.SeoQueryEvidence, error) {
	item, ok := r.items[id]
	if !ok {
		return nil, repo.ErrSeoQueryEvidenceNotFound
	}
	copyItem := *item
	return &copyItem, nil
}

func TestSeoIntentHypothesisUsecaseCreateRequiresActiveUserNeed(t *testing.T) {
	repository := newFakeIntentHypothesisRepo()
	u := &SeoIntentHypothesisUsecase{
		repository: repository,
		userNeedReader: &fakeIntentUserNeedReader{item: &seo_strategy.SeoUserNeed{
			ID: 20, ObjectiveID: 10, Status: seo_strategy.UserNeedStatusDraft,
		}},
		queryEvidenceReader: &fakeIntentQueryEvidenceReader{items: map[uint64]*seo_strategy.SeoQueryEvidence{}},
		now:                 func() time.Time { return time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC) },
	}
	_, err := u.Create(context.Background(), &dto.CreateSeoIntentHypothesisRequest{
		UserNeedID: 20, Code: "verify_parcel", Name: "Verify parcel",
	})
	if err == nil {
		t.Fatal("expected inactive user need error")
	}
}

func TestSeoIntentHypothesisUsecaseCreate(t *testing.T) {
	repository := newFakeIntentHypothesisRepo()
	u := &SeoIntentHypothesisUsecase{
		repository: repository,
		userNeedReader: &fakeIntentUserNeedReader{item: &seo_strategy.SeoUserNeed{
			ID: 20, ObjectiveID: 10, Status: seo_strategy.UserNeedStatusActive,
		}},
		queryEvidenceReader: &fakeIntentQueryEvidenceReader{items: map[uint64]*seo_strategy.SeoQueryEvidence{}},
		now:                 func() time.Time { return time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC) },
	}
	result, err := u.Create(context.Background(), &dto.CreateSeoIntentHypothesisRequest{
		UserNeedID: 20, Code: "verify parcel", Name: "Verify parcel",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Intent.ID == 0 || result.Intent.ObjectiveID != 10 || result.Intent.Code != "verify_parcel" {
		t.Fatalf("unexpected result: %+v", result.Intent)
	}
}

func TestSeoIntentHypothesisUsecaseLinkEvidenceRejectsMismatch(t *testing.T) {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	repository := newFakeIntentHypothesisRepo()
	repository.items[1] = &seo_strategy.SeoIntentHypothesis{
		ID: 1, ObjectiveID: 10, UserNeedID: 20, Code: "verify_parcel", Name: "Verify parcel",
		SearchAction: seo_strategy.IntentSearchActionUnknown, SubjectType: seo_strategy.IntentSubjectUnknown,
		ScopeType: seo_strategy.IntentScopeUnknown, Status: seo_strategy.IntentHypothesisStatusDraft,
		Confidence: seo_strategy.ObjectiveConfidenceExploratory, LockVersion: 1,
	}
	u := &SeoIntentHypothesisUsecase{
		repository: repository,
		userNeedReader: &fakeIntentUserNeedReader{item: &seo_strategy.SeoUserNeed{
			ID: 20, ObjectiveID: 10, Status: seo_strategy.UserNeedStatusActive,
		}},
		queryEvidenceReader: &fakeIntentQueryEvidenceReader{items: map[uint64]*seo_strategy.SeoQueryEvidence{
			40: {ID: 40, ObjectiveID: 11, UserNeedID: 21, Status: seo_strategy.QueryEvidenceStatusActive},
		}},
		now: func() time.Time { return now },
	}
	_, err := u.LinkEvidence(context.Background(), &dto.LinkSeoIntentEvidenceRequest{
		IntentHypothesisID: 1,
		QueryEvidenceID:    40,
		RelationType:       "supports",
	})
	if err == nil {
		t.Fatal("expected mismatch error")
	}
}

func TestSeoIntentHypothesisUsecaseActivateUsesReadinessAndVersion(t *testing.T) {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	review := now.Add(24 * time.Hour)
	repository := newFakeIntentHypothesisRepo()
	repository.items[1] = &seo_strategy.SeoIntentHypothesis{
		ID: 1, ObjectiveID: 10, UserNeedID: 20, Code: "verify_parcel", Name: "Verify parcel",
		GoalStatement: "Verify planning impact", PrimaryQuestion: "Is this parcel planned?",
		ExpectedResult: "Source-backed answer", SearchAction: seo_strategy.IntentSearchActionVerify,
		SubjectType: seo_strategy.IntentSubjectParcel, ScopeType: seo_strategy.IntentScopeSpecificEntity,
		Status: seo_strategy.IntentHypothesisStatusDraft, Confidence: seo_strategy.ObjectiveConfidenceExploratory,
		OwnerCode: "seo", ReviewDueAt: &review, LockVersion: 1,
	}
	evidence := seo_strategy.SeoQueryEvidence{
		ID: 40, ObjectiveID: 10, UserNeedID: 20, RawQuery: "parcel planning",
		NormalizedQuery: "parcel planning", QueryHash: seo_strategy.BuildQueryEvidenceHash("parcel planning", "vi-VN"),
		SourceSystem: "search_console", SourceKind: seo_strategy.QueryEvidenceSourceFirstParty,
		ObservationType: seo_strategy.QueryObservationObserved, SourceRecordKey: "40",
		Locale: "vi-VN", CountryCode: "VN", DeviceCategory: seo_strategy.QueryDeviceAll,
		Confidence: seo_strategy.ObjectiveConfidenceValidated, Status: seo_strategy.QueryEvidenceStatusActive,
	}
	repository.links[1] = []seo_strategy.LinkedIntentQueryEvidence{
		{Link: seo_strategy.SeoIntentEvidenceLink{RelationType: seo_strategy.IntentEvidenceSupports}, Evidence: evidence},
		{Link: seo_strategy.SeoIntentEvidenceLink{RelationType: seo_strategy.IntentEvidenceContext}, Evidence: evidence},
	}
	u := &SeoIntentHypothesisUsecase{
		repository: repository,
		userNeedReader: &fakeIntentUserNeedReader{item: &seo_strategy.SeoUserNeed{
			ID: 20, ObjectiveID: 10, Status: seo_strategy.UserNeedStatusActive,
		}},
		queryEvidenceReader: &fakeIntentQueryEvidenceReader{items: map[uint64]*seo_strategy.SeoQueryEvidence{40: &evidence}},
		now:                 func() time.Time { return now },
	}
	result, err := u.Activate(context.Background(), &dto.SeoIntentHypothesisCommandRequest{
		ID: 1, ExpectedVersion: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Intent.Status != seo_strategy.IntentHypothesisStatusActive || result.Intent.LockVersion != 2 {
		t.Fatalf("unexpected lifecycle result: %+v", result.Intent)
	}
	_, err = u.Pause(context.Background(), &dto.SeoIntentHypothesisCommandRequest{ID: 1, ExpectedVersion: 1})
	if err == nil {
		t.Fatal("expected stale version error")
	}
}
