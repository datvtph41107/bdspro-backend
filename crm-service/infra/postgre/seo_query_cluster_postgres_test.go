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

func TestSeoQueryClusterPostgres(t *testing.T) {
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
	ctx := context.Background()

	intentRepo := &SeoIntentHypothesisPostgres{db: tx}
	queryRepo := &SeoQueryEvidencePostgres{db: tx}
	clusterRepo := &SeoQueryClusterPostgres{db: tx}

	intent, err := seo_strategy.CreateSeoIntentHypothesisDraft(
		objectiveID,
		userNeedID,
		"cluster_parent_intent",
		"Cluster parent intent",
		now,
	)
	if err != nil {
		t.Fatal(err)
	}
	intent.GoalStatement = "Verify parcel planning impact"
	intent.PrimaryQuestion = "Does this parcel have planning impact?"
	intent.ExpectedResult = "A source-backed planning answer"
	intent.SearchAction = seo_strategy.IntentSearchActionVerify
	intent.SubjectType = seo_strategy.IntentSubjectParcel
	intent.ScopeType = seo_strategy.IntentScopeSpecificEntity
	intent.OwnerCode = "seo_strategy"
	if err := intentRepo.Create(ctx, intent); err != nil {
		t.Fatal(err)
	}
	if err := intentRepo.UpdateStatus(ctx, intent.ID, seo_strategy.IntentHypothesisStatusActive, 1); err != nil {
		t.Fatal(err)
	}

	primary := seedQueryClusterEvidence(t, ctx, queryRepo, intentRepo, objectiveID, userNeedID, intent.ID, "primary", now)
	included := seedQueryClusterEvidence(t, ctx, queryRepo, intentRepo, objectiveID, userNeedID, intent.ID, "included", now)
	excluded := seedQueryClusterEvidence(t, ctx, queryRepo, intentRepo, objectiveID, userNeedID, intent.ID, "excluded", now)

	cluster, err := seo_strategy.CreateSeoQueryClusterDraft(
		objectiveID,
		userNeedID,
		intent.ID,
		"parcel_planning_answer_set",
		"Parcel planning answer set",
		now,
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := clusterRepo.Create(ctx, cluster); err != nil {
		t.Fatal(err)
	}
	if cluster.ID == 0 {
		t.Fatal("query cluster repository phải trả persisted ID")
	}

	cluster.SearchPromise = "Trả lời tác động quy hoạch của đúng thửa đất, có nguồn và giới hạn dữ liệu"
	cluster.InclusionRule = "Query có thể được trả lời bằng cùng primary answer về thửa đất và quy hoạch"
	cluster.ExclusionRule = "Query cần một answer khác như giá, giao dịch hoặc thủ tục pháp lý"
	cluster.OwnerCode = "seo_strategy"
	cluster.Confidence = seo_strategy.ObjectiveConfidenceSupported
	if err := clusterRepo.UpdateDefinition(ctx, cluster, 1); err != nil {
		t.Fatal(err)
	}
	got, err := clusterRepo.GetByID(ctx, cluster.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.LockVersion != 2 || got.SearchPromise == "" {
		t.Fatalf("unexpected updated cluster: %+v", got)
	}
	if err := clusterRepo.UpdateDefinition(ctx, cluster, 1); !errors.Is(err, repo.ErrSeoQueryClusterConflict) {
		t.Fatalf("expected optimistic locking conflict, got %v", err)
	}

	primaryMember := &seo_strategy.SeoQueryClusterMember{
		ObjectiveID: objectiveID, UserNeedID: userNeedID,
		ClusterID: cluster.ID, QueryEvidenceID: primary.ID,
		MembershipType: seo_strategy.QueryClusterMembershipPrimary,
		Status:         seo_strategy.QueryClusterMemberStatusActive,
		CreatedAt:      now, UpdatedAt: now,
	}
	if err := clusterRepo.UpsertMember(ctx, primaryMember, 2); err != nil {
		t.Fatal(err)
	}

	secondPrimary := &seo_strategy.SeoQueryClusterMember{
		ObjectiveID: objectiveID, UserNeedID: userNeedID,
		ClusterID: cluster.ID, QueryEvidenceID: included.ID,
		MembershipType: seo_strategy.QueryClusterMembershipPrimary,
		Status:         seo_strategy.QueryClusterMemberStatusActive,
		CreatedAt:      now, UpdatedAt: now,
	}
	if err := clusterRepo.UpsertMember(ctx, secondPrimary, 3); !errors.Is(err, repo.ErrSeoQueryClusterMembershipConflict) {
		t.Fatalf("expected one-primary conflict, got %v", err)
	}
	got, err = clusterRepo.GetByID(ctx, cluster.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.LockVersion != 3 {
		t.Fatalf("membership conflict phải rollback version, got %d", got.LockVersion)
	}

	includedMember := &seo_strategy.SeoQueryClusterMember{
		ObjectiveID: objectiveID, UserNeedID: userNeedID,
		ClusterID: cluster.ID, QueryEvidenceID: included.ID,
		MembershipType: seo_strategy.QueryClusterMembershipIncluded,
		Rationale:      "Cùng primary answer về planning impact",
		Status:         seo_strategy.QueryClusterMemberStatusActive,
		CreatedAt:      now, UpdatedAt: now,
	}
	if err := clusterRepo.UpsertMember(ctx, includedMember, 3); err != nil {
		t.Fatal(err)
	}

	excludedMember := &seo_strategy.SeoQueryClusterMember{
		ObjectiveID: objectiveID, UserNeedID: userNeedID,
		ClusterID: cluster.ID, QueryEvidenceID: excluded.ID,
		MembershipType: seo_strategy.QueryClusterMembershipExcluded,
		Rationale:      "Cần price answer, không phải planning answer",
		Status:         seo_strategy.QueryClusterMemberStatusActive,
		CreatedAt:      now, UpdatedAt: now,
	}
	if err := clusterRepo.UpsertMember(ctx, excludedMember, 4); err != nil {
		t.Fatal(err)
	}

	members, err := clusterRepo.ListMembers(ctx, cluster.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(members) != 3 || members[0].Member.MembershipType != seo_strategy.QueryClusterMembershipPrimary {
		t.Fatalf("unexpected members: %+v", members)
	}

	secondCluster, err := seo_strategy.CreateSeoQueryClusterDraft(
		objectiveID,
		userNeedID,
		intent.ID,
		"second_answer_set",
		"Second answer set",
		now,
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := clusterRepo.Create(ctx, secondCluster); err != nil {
		t.Fatal(err)
	}
	conflictingMember := &seo_strategy.SeoQueryClusterMember{
		ObjectiveID: objectiveID, UserNeedID: userNeedID,
		ClusterID: secondCluster.ID, QueryEvidenceID: included.ID,
		MembershipType: seo_strategy.QueryClusterMembershipIncluded,
		Rationale:      "Attempt duplicate active ownership",
		Status:         seo_strategy.QueryClusterMemberStatusActive,
		CreatedAt:      now, UpdatedAt: now,
	}
	if err := clusterRepo.UpsertMember(ctx, conflictingMember, 1); !errors.Is(err, repo.ErrSeoQueryClusterMembershipConflict) {
		t.Fatalf("expected query ownership conflict, got %v", err)
	}

	if err := clusterRepo.RemoveMember(ctx, cluster.ID, excluded.ID, 5); err != nil {
		t.Fatal(err)
	}
	if err := clusterRepo.RemoveMember(ctx, cluster.ID, excluded.ID, 5); !errors.Is(err, repo.ErrSeoQueryClusterConflict) {
		t.Fatalf("expected stale version conflict, got %v", err)
	}
	if err := clusterRepo.UpdateStatus(ctx, cluster.ID, seo_strategy.QueryClusterStatusActive, 6); err != nil {
		t.Fatal(err)
	}

	items, total, err := clusterRepo.List(ctx, repo.SeoQueryClusterFilter{
		ObjectiveID:        objectiveID,
		UserNeedID:         userNeedID,
		IntentHypothesisID: intent.ID,
		Status:             seo_strategy.QueryClusterStatusActive,
		Size:               20,
	})
	if err != nil || total != 1 || len(items) != 1 {
		t.Fatalf("List total=%d len=%d err=%v", total, len(items), err)
	}
}

func seedQueryClusterEvidence(
	t *testing.T,
	ctx context.Context,
	queryRepo *SeoQueryEvidencePostgres,
	intentRepo *SeoIntentHypothesisPostgres,
	objectiveID uint64,
	userNeedID uint64,
	intentID uint64,
	key string,
	now time.Time,
) *seo_strategy.SeoQueryEvidence {
	t.Helper()
	evidence, err := seo_strategy.CreateSeoQueryEvidence(
		objectiveID,
		userNeedID,
		"query cluster "+key,
		"integration_test",
		seo_strategy.QueryEvidenceSourceFirstParty,
		seo_strategy.QueryObservationObserved,
		fmt.Sprintf("cluster|%s|%d", key, now.UnixNano()),
		now,
	)
	if err != nil {
		t.Fatal(err)
	}
	evidence.Confidence = seo_strategy.ObjectiveConfidenceValidated
	if err := queryRepo.Upsert(ctx, evidence); err != nil {
		t.Fatal(err)
	}
	relation := seo_strategy.IntentEvidenceSupports
	if key == "excluded" {
		relation = seo_strategy.IntentEvidenceContext
	}
	if err := intentRepo.LinkEvidence(ctx, &seo_strategy.SeoIntentEvidenceLink{
		ObjectiveID:        objectiveID,
		UserNeedID:         userNeedID,
		IntentHypothesisID: intentID,
		QueryEvidenceID:    evidence.ID,
		RelationType:       relation,
		RelevanceNote:      "cluster integration evidence",
		CreatedAt:          now,
		UpdatedAt:          now,
	}); err != nil {
		t.Fatal(err)
	}
	return evidence
}
