package postgre

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	seo_strategy "crm/internal/domain/seo_strategy"
	"crm/internal/repo"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestSeoIntentHypothesisPostgres(t *testing.T) {
	dsn := os.Getenv("SEO_STRATEGY_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("SEO_STRATEGY_TEST_DATABASE_URL chưa được gán")
	}

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)},
	)
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatal(tx.Error)
	}
	t.Cleanup(func() { _ = tx.Rollback().Error })

	objectiveID, userNeedID := seedIntentHypothesisParents(t, tx)
	now := time.Now().UTC().Truncate(time.Microsecond)

	intentRepo := &SeoIntentHypothesisPostgres{db: tx}
	queryRepo := &SeoQueryEvidencePostgres{db: tx}

	intent, err := seo_strategy.CreateSeoIntentHypothesisDraft(
		objectiveID,
		userNeedID,
		"verify_parcel_planning",
		"Xác minh quy hoạch thửa đất",
		now,
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := intentRepo.Create(context.Background(), intent); err != nil {
		t.Fatal(err)
	}
	if intent.ID == 0 {
		t.Fatal("intent repository phải trả persisted ID")
	}

	intent.GoalStatement = "Xác minh tác động quy hoạch trước quyết định bất động sản"
	intent.PrimaryQuestion = "Thửa đất này có thuộc quy hoạch hay không?"
	intent.ExpectedResult = "Kết quả có nguồn, thời gian và giới hạn dữ liệu"
	intent.SearchAction = seo_strategy.IntentSearchActionVerify
	intent.SubjectType = seo_strategy.IntentSubjectParcel
	intent.ScopeType = seo_strategy.IntentScopeSpecificEntity
	intent.OwnerCode = "seo_strategy"
	intent.Confidence = seo_strategy.ObjectiveConfidenceSupported

	if err := intentRepo.UpdateDefinition(context.Background(), intent, 1); err != nil {
		t.Fatal(err)
	}
	got, err := intentRepo.GetByID(context.Background(), intent.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.LockVersion != 2 || got.SearchAction != seo_strategy.IntentSearchActionVerify {
		t.Fatalf("unexpected updated intent: %+v", got)
	}

	err = intentRepo.UpdateDefinition(context.Background(), intent, 1)
	if !errors.Is(err, repo.ErrSeoIntentHypothesisConflict) {
		t.Fatalf("expected optimistic locking conflict, got %v", err)
	}

	evidence, err := seo_strategy.CreateSeoQueryEvidence(
		objectiveID,
		userNeedID,
		"tra cứu quy hoạch thửa đất",
		"integration_test",
		seo_strategy.QueryEvidenceSourceFirstParty,
		seo_strategy.QueryObservationObserved,
		fmt.Sprintf("intent|%d", now.UnixNano()),
		now,
	)
	if err != nil {
		t.Fatal(err)
	}
	evidence.Confidence = seo_strategy.ObjectiveConfidenceValidated
	if err := queryRepo.Upsert(context.Background(), evidence); err != nil {
		t.Fatal(err)
	}

	link := &seo_strategy.SeoIntentEvidenceLink{
		ObjectiveID:        objectiveID,
		UserNeedID:         userNeedID,
		IntentHypothesisID: intent.ID,
		QueryEvidenceID:    evidence.ID,
		RelationType:       seo_strategy.IntentEvidenceSupports,
		RelevanceNote:      "Query thực tế xác nhận hành vi verify",
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := intentRepo.LinkEvidence(context.Background(), link); err != nil {
		t.Fatal(err)
	}
	if err := intentRepo.LinkEvidence(context.Background(), link); !errors.Is(err, repo.ErrSeoIntentEvidenceAlreadyLinked) {
		t.Fatalf("expected duplicate link error, got %v", err)
	}

	linked, err := intentRepo.ListLinkedEvidence(context.Background(), intent.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(linked) != 1 || linked[0].Evidence.ID != evidence.ID {
		t.Fatalf("unexpected linked evidence: %+v", linked)
	}

	items, total, err := intentRepo.List(
		context.Background(),
		repo.SeoIntentHypothesisFilter{
			ObjectiveID:  objectiveID,
			UserNeedID:   userNeedID,
			SearchAction: seo_strategy.IntentSearchActionVerify,
			Size:         20,
		},
	)
	if err != nil || total != 1 || len(items) != 1 {
		t.Fatalf("List total=%d len=%d err=%v", total, len(items), err)
	}

	if err := intentRepo.UpdateStatus(
		context.Background(),
		intent.ID,
		seo_strategy.IntentHypothesisStatusActive,
		2,
	); err != nil {
		t.Fatal(err)
	}
	if err := intentRepo.UpdateStatus(
		context.Background(),
		intent.ID,
		seo_strategy.IntentHypothesisStatusPaused,
		2,
	); !errors.Is(err, repo.ErrSeoIntentHypothesisConflict) {
		t.Fatalf("expected stale status conflict, got %v", err)
	}

	if err := intentRepo.UnlinkEvidence(context.Background(), intent.ID, evidence.ID); err != nil {
		t.Fatal(err)
	}
	if err := intentRepo.UnlinkEvidence(context.Background(), intent.ID, evidence.ID); !errors.Is(err, repo.ErrSeoIntentEvidenceNotLinked) {
		t.Fatalf("expected not linked error, got %v", err)
	}
}

func seedIntentHypothesisParents(
	t *testing.T,
	tx *gorm.DB,
) (objectiveID uint64, userNeedID uint64) {
	t.Helper()
	suffix := fmt.Sprintf("%d", time.Now().UnixNano())

	if err := tx.Raw(
		`INSERT INTO seo_objective (code, name, status, confidence)
		 VALUES (?, ?, 'active', 'exploratory') RETURNING id`,
		"intent_objective_"+suffix,
		"Intent Objective",
	).Row().Scan(&objectiveID); err != nil {
		t.Fatal(err)
	}

	var audienceID uint64
	if err := tx.Raw(
		`INSERT INTO seo_audience_segment (code, name, status)
		 VALUES (?, ?, 'active') RETURNING id`,
		"intent_audience_"+suffix,
		"Intent Audience",
	).Row().Scan(&audienceID); err != nil {
		t.Fatal(err)
	}

	if err := tx.Raw(
		`INSERT INTO seo_user_need (
			objective_id,
			audience_segment_id,
			code,
			name,
			status,
			confidence
		 ) VALUES (?, ?, ?, ?, 'active', 'exploratory') RETURNING id`,
		objectiveID,
		audienceID,
		"intent_need_"+suffix,
		"Intent User Need",
	).Row().Scan(&userNeedID); err != nil {
		t.Fatal(err)
	}

	return objectiveID, userNeedID
}
