package usecase

import (
	"context"
	"testing"
	"time"

	seo_strategy "crm/internal/domain/seo_strategy"
	"crm/internal/dto"
	"crm/internal/repo"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/status"
)

type fakeQueryEvidenceUserNeedReader struct {
	need *seo_strategy.SeoUserNeed
	err  error
}

func (f *fakeQueryEvidenceUserNeedReader) GetByID(
	_ context.Context,
	_ uint64,
) (*seo_strategy.SeoUserNeed, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.need == nil {
		return nil, repo.ErrSeoUserNeedNotFound
	}
	copyValue := *f.need
	return &copyValue, nil
}

type fakeQueryEvidenceRepo struct {
	item       *seo_strategy.SeoQueryEvidence
	listItems  []seo_strategy.SeoQueryEvidence
	listTotal  int64
	lastFilter repo.SeoQueryEvidenceFilter
	upsertErr  error
	getErr     error
	listErr    error
	updateErr  error
}

func (f *fakeQueryEvidenceRepo) Upsert(
	_ context.Context,
	item *seo_strategy.SeoQueryEvidence,
) error {
	if f.upsertErr != nil {
		return f.upsertErr
	}
	copyValue := *item
	if copyValue.ID == 0 {
		copyValue.ID = 100
	}
	f.item = &copyValue
	*item = copyValue
	return nil
}

func (f *fakeQueryEvidenceRepo) GetByID(
	_ context.Context,
	_ uint64,
) (*seo_strategy.SeoQueryEvidence, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	if f.item == nil {
		return nil, repo.ErrSeoQueryEvidenceNotFound
	}
	copyValue := *f.item
	return &copyValue, nil
}

func (f *fakeQueryEvidenceRepo) GetBySourceRecordKey(
	_ context.Context,
	_ uint64,
	_ string,
	_ string,
) (*seo_strategy.SeoQueryEvidence, error) {
	if f.item == nil {
		return nil, repo.ErrSeoQueryEvidenceNotFound
	}
	copyValue := *f.item
	return &copyValue, nil
}

func (f *fakeQueryEvidenceRepo) List(
	_ context.Context,
	filter repo.SeoQueryEvidenceFilter,
) ([]seo_strategy.SeoQueryEvidence, int64, error) {
	f.lastFilter = filter
	if f.listErr != nil {
		return nil, 0, f.listErr
	}
	return f.listItems, f.listTotal, nil
}

func (f *fakeQueryEvidenceRepo) UpdateStatus(
	_ context.Context,
	_ uint64,
	expected seo_strategy.QueryEvidenceStatus,
	target seo_strategy.QueryEvidenceStatus,
) error {
	if f.updateErr != nil {
		return f.updateErr
	}
	if f.item == nil {
		return repo.ErrSeoQueryEvidenceNotFound
	}
	if f.item.Status != expected {
		return repo.ErrSeoQueryEvidenceConflict
	}
	f.item.Status = target
	return nil
}

func queryEvidenceApplicationCode(t *testing.T, err error) int32 {
	t.Helper()
	if err == nil {
		t.Fatal("expected error")
	}

	grpcStatus, ok := status.FromError(err)
	if !ok {
		t.Fatalf("error không phải grpc status: %v", err)
	}
	for _, detail := range grpcStatus.Details() {
		if errorResponse, ok := detail.(*sharepb.ErrorResponse); ok {
			return errorResponse.GetCode()
		}
	}
	t.Fatalf("không tìm thấy ErrorResponse: %v", err)
	return 0
}

func TestSeoQueryEvidenceUsecaseUpsert(t *testing.T) {
	fixedNow := time.Date(2026, time.July, 14, 10, 0, 0, 0, time.UTC)
	queryRepo := &fakeQueryEvidenceRepo{}
	uc := &SeoQueryEvidenceUsecase{
		repository: queryRepo,
		userNeedReader: &fakeQueryEvidenceUserNeedReader{
			need: &seo_strategy.SeoUserNeed{
				ID:          20,
				ObjectiveID: 10,
				Status:      seo_strategy.UserNeedStatusActive,
			},
		},
		now: func() time.Time { return fixedNow },
	}

	impressions := int64(100)
	clicks := int64(10)
	ctr := 0.1
	result, err := uc.Upsert(
		context.Background(),
		&dto.UpsertSeoQueryEvidenceRequest{
			UserNeedID:      20,
			RawQuery:        "  Tra CỨU quy hoạch  ",
			SourceSystem:    "google_search_console",
			SourceKind:      "first_party",
			ObservationType: "observed",
			SourceRecordKey: "gsc|2026-07|query-1",
			Locale:          "vi-VN",
			CountryCode:     "VN",
			DeviceCategory:  "mobile",
			Impressions:     &impressions,
			Clicks:          &clicks,
			CTR:             &ctr,
			Confidence:      "supported",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != 100 || result.ObjectiveID != 10 || result.UserNeedID != 20 {
		t.Fatalf("result=%+v", result)
	}
	if result.NormalizedQuery != "tra cứu quy hoạch" {
		t.Fatalf("normalizedQuery=%q", result.NormalizedQuery)
	}
	if result.Status != seo_strategy.QueryEvidenceStatusActive {
		t.Fatalf("status=%q", result.Status)
	}
}

func TestSeoQueryEvidenceUsecaseRejectsInactiveUserNeed(t *testing.T) {
	statuses := []seo_strategy.UserNeedStatus{
		seo_strategy.UserNeedStatusDraft,
		seo_strategy.UserNeedStatusPaused,
		seo_strategy.UserNeedStatusArchived,
	}

	for _, userNeedStatus := range statuses {
		userNeedStatus := userNeedStatus
		t.Run(string(userNeedStatus), func(t *testing.T) {
			uc := &SeoQueryEvidenceUsecase{
				repository: &fakeQueryEvidenceRepo{},
				userNeedReader: &fakeQueryEvidenceUserNeedReader{
					need: &seo_strategy.SeoUserNeed{
						ID:          20,
						ObjectiveID: 10,
						Status:      userNeedStatus,
					},
				},
				now: time.Now,
			}
			_, err := uc.Upsert(
				context.Background(),
				&dto.UpsertSeoQueryEvidenceRequest{
					UserNeedID:      20,
					RawQuery:        "tra cứu quy hoạch",
					SourceSystem:    "manual_research",
					SourceKind:      "manual",
					ObservationType: "assumed",
					SourceRecordKey: "manual|1",
				},
			)
			if code := queryEvidenceApplicationCode(t, err); code != 409 {
				t.Fatalf("code=%d, want=409", code)
			}
		})
	}
}

func TestSeoQueryEvidenceUsecaseRejectsInvalidEnum(t *testing.T) {
	uc := &SeoQueryEvidenceUsecase{
		repository: &fakeQueryEvidenceRepo{},
		userNeedReader: &fakeQueryEvidenceUserNeedReader{
			need: &seo_strategy.SeoUserNeed{
				ID:          20,
				ObjectiveID: 10,
				Status:      seo_strategy.UserNeedStatusActive,
			},
		},
		now: time.Now,
	}

	_, err := uc.Upsert(
		context.Background(),
		&dto.UpsertSeoQueryEvidenceRequest{
			UserNeedID:      20,
			RawQuery:        "tra cứu quy hoạch",
			SourceSystem:    "test_source",
			SourceKind:      "firstparty",
			ObservationType: "observed",
			SourceRecordKey: "test|1",
		},
	)
	if code := queryEvidenceApplicationCode(t, err); code != 400 {
		t.Fatalf("code=%d, want=400", code)
	}
}

func TestSeoQueryEvidenceUsecaseListBuildsTypedFilter(t *testing.T) {
	fixedNow := time.Date(2026, time.July, 14, 10, 0, 0, 0, time.UTC)
	queryRepo := &fakeQueryEvidenceRepo{}
	uc := &SeoQueryEvidenceUsecase{
		repository:     queryRepo,
		userNeedReader: &fakeQueryEvidenceUserNeedReader{},
		now:            func() time.Time { return fixedNow },
	}

	_, _, err := uc.List(
		context.Background(),
		&dto.SeoQueryEvidenceListRequest{
			ObjectiveID:     10,
			UserNeedID:      20,
			SourceKind:      "first_party",
			ObservationType: "observed",
			DeviceCategory:  "mobile",
			Confidence:      "supported",
			Status:          "active",
			CurrentOnly:     true,
			Page:            1,
			Size:            25,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	filter := queryRepo.lastFilter
	if filter.SourceKind != seo_strategy.QueryEvidenceSourceFirstParty ||
		filter.ObservationType != seo_strategy.QueryObservationObserved ||
		filter.DeviceCategory != seo_strategy.QueryDeviceMobile ||
		filter.Status != seo_strategy.QueryEvidenceStatusActive {
		t.Fatalf("filter=%+v", filter)
	}
	if filter.CurrentAt == nil || !filter.CurrentAt.Equal(fixedNow) {
		t.Fatalf("currentAt=%v", filter.CurrentAt)
	}
}

func TestSeoQueryEvidenceUsecaseStatusLifecycle(t *testing.T) {
	now := time.Now().UTC()
	evidence, err := seo_strategy.CreateSeoQueryEvidence(
		10,
		20,
		"tra cứu quy hoạch",
		"test_source",
		seo_strategy.QueryEvidenceSourceManual,
		seo_strategy.QueryObservationAssumed,
		"test|status",
		now,
	)
	if err != nil {
		t.Fatal(err)
	}
	evidence.ID = 100

	queryRepo := &fakeQueryEvidenceRepo{item: evidence}
	uc := &SeoQueryEvidenceUsecase{
		repository:     queryRepo,
		userNeedReader: &fakeQueryEvidenceUserNeedReader{},
		now:            func() time.Time { return now },
	}

	superseded, err := uc.Supersede(
		context.Background(),
		&dto.SeoQueryEvidenceStatusCommandRequest{ID: 100, ExpectedStatus: "active"},
	)
	if err != nil || superseded.Status != seo_strategy.QueryEvidenceStatusSuperseded {
		t.Fatalf("supersede=%+v err=%v", superseded, err)
	}

	rejected, err := uc.Reject(
		context.Background(),
		&dto.SeoQueryEvidenceStatusCommandRequest{ID: 100, ExpectedStatus: "superseded"},
	)
	if err != nil || rejected.Status != seo_strategy.QueryEvidenceStatusRejected {
		t.Fatalf("reject=%+v err=%v", rejected, err)
	}

	_, err = uc.Supersede(
		context.Background(),
		&dto.SeoQueryEvidenceStatusCommandRequest{ID: 100, ExpectedStatus: "rejected"},
	)
	if code := queryEvidenceApplicationCode(t, err); code != 409 {
		t.Fatalf("code=%d, want=409", code)
	}
}
