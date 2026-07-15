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

func TestSeoQueryEvidencePostgres(t *testing.T) {
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

	objectiveID, secondObjectiveID, userNeedID := seedQueryEvidenceParents(t, tx)
	repository := &SeoQueryEvidencePostgres{db: tx}
	now := time.Now().UTC().Truncate(time.Microsecond)
	sourceKey := fmt.Sprintf("test|%d", now.UnixNano())

	evidence, err := seo_strategy.CreateSeoQueryEvidence(
		objectiveID,
		userNeedID,
		"Tra cứu quy hoạch số tờ số thửa",
		"integration_test",
		seo_strategy.QueryEvidenceSourceFirstParty,
		seo_strategy.QueryObservationObserved,
		sourceKey,
		now,
	)
	if err != nil {
		t.Fatal(err)
	}
	impressions := int64(100)
	clicks := int64(12)
	ctr := 0.12
	evidence.Impressions = &impressions
	evidence.Clicks = &clicks
	evidence.CTR = &ctr

	if err := repository.Upsert(context.Background(), evidence); err != nil {
		t.Fatal(err)
	}
	if evidence.ID == 0 {
		t.Fatal("repository phải trả persisted ID")
	}
	firstID := evidence.ID

	updatedImpressions := int64(150)
	evidence.Impressions = &updatedImpressions
	if err := repository.Upsert(context.Background(), evidence); err != nil {
		t.Fatal(err)
	}
	if evidence.ID != firstID {
		t.Fatalf("upsert tạo duplicate: first=%d second=%d", firstID, evidence.ID)
	}

	got, err := repository.GetByID(context.Background(), firstID)
	if err != nil || got.SourceRecordKey != sourceKey {
		t.Fatalf("GetByID item=%+v err=%v", got, err)
	}

	got, err = repository.GetBySourceRecordKey(
		context.Background(),
		userNeedID,
		"integration_test",
		sourceKey,
	)
	if err != nil || got.ID != firstID {
		t.Fatalf("GetBySourceRecordKey item=%+v err=%v", got, err)
	}

	items, total, err := repository.List(
		context.Background(),
		repo.SeoQueryEvidenceFilter{
			ObjectiveID:     objectiveID,
			UserNeedID:      userNeedID,
			SourceKind:      seo_strategy.QueryEvidenceSourceFirstParty,
			ObservationType: seo_strategy.QueryObservationObserved,
			Status:          seo_strategy.QueryEvidenceStatusActive,
			Size:            20,
		},
	)
	if err != nil || total != 1 || len(items) != 1 {
		t.Fatalf("List total=%d len=%d err=%v", total, len(items), err)
	}

	if err := repository.UpdateStatus(
		context.Background(),
		firstID,
		seo_strategy.QueryEvidenceStatusActive,
		seo_strategy.QueryEvidenceStatusSuperseded,
	); err != nil {
		t.Fatal(err)
	}

	err = repository.UpdateStatus(
		context.Background(),
		firstID,
		seo_strategy.QueryEvidenceStatusActive,
		seo_strategy.QueryEvidenceStatusRejected,
	)
	if !errors.Is(err, repo.ErrSeoQueryEvidenceConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}

	invalid, err := seo_strategy.CreateSeoQueryEvidence(
		secondObjectiveID,
		userNeedID,
		"Query sai objective",
		"integration_test",
		seo_strategy.QueryEvidenceSourceManual,
		seo_strategy.QueryObservationAssumed,
		sourceKey+"|mismatch",
		now,
	)
	if err != nil {
		t.Fatal(err)
	}
	err = repository.Upsert(context.Background(), invalid)
	if !errors.Is(err, repo.ErrSeoQueryEvidenceUserNeedMismatch) {
		t.Fatalf("expected mismatch, got %v", err)
	}
}

func seedQueryEvidenceParents(
	t *testing.T,
	tx *gorm.DB,
) (objectiveID uint64, secondObjectiveID uint64, userNeedID uint64) {
	t.Helper()

	suffix := fmt.Sprintf("%d", time.Now().UnixNano())

	if err := tx.Raw(
		`INSERT INTO seo_objective (code, name, status, confidence)
		 VALUES (?, ?, 'active', 'exploratory') RETURNING id`,
		"query_evidence_objective_"+suffix,
		"Query Evidence Objective",
	).Row().Scan(&objectiveID); err != nil {
		t.Fatal(err)
	}

	if err := tx.Raw(
		`INSERT INTO seo_objective (code, name, status, confidence)
		 VALUES (?, ?, 'active', 'exploratory') RETURNING id`,
		"query_evidence_other_"+suffix,
		"Other Objective",
	).Row().Scan(&secondObjectiveID); err != nil {
		t.Fatal(err)
	}

	var audienceID uint64
	if err := tx.Raw(
		`INSERT INTO seo_audience_segment (code, name, status)
		 VALUES (?, ?, 'active') RETURNING id`,
		"query_evidence_audience_"+suffix,
		"Query Evidence Audience",
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
		"query_evidence_need_"+suffix,
		"Query Evidence User Need",
	).Row().Scan(&userNeedID); err != nil {
		t.Fatal(err)
	}

	return objectiveID, secondObjectiveID, userNeedID
}
