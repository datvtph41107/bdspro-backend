package seo_strategy

import (
	"testing"
	"time"
)

func TestEvaluateQueryClusterReadinessReady(t *testing.T) {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	need, intent, cluster := readyQueryClusterParents(now)
	members, intentEvidence := readyQueryClusterMembers(now)

	readiness := EvaluateQueryClusterReadiness(need, intent, cluster, members, intentEvidence, now)
	if !readiness.Ready {
		t.Fatalf("expected ready, blockers=%#v warnings=%#v", readiness.Blockers, readiness.Warnings)
	}
	if !readiness.AnswerEquivalenceReady {
		t.Fatal("expected answer equivalence ready")
	}
	if readiness.PrimaryCount != 1 || readiness.IncludedCount != 1 || readiness.ExcludedCount != 1 {
		t.Fatalf("unexpected counts: %#v", readiness)
	}
}

func TestEvaluateQueryClusterReadinessRequiresIncludedQuery(t *testing.T) {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	need, intent, cluster := readyQueryClusterParents(now)
	members, intentEvidence := readyQueryClusterMembers(now)
	members = []LinkedQueryClusterMember{members[0], members[2]}

	readiness := EvaluateQueryClusterReadiness(need, intent, cluster, members, intentEvidence, now)
	if readiness.Ready || readiness.AnswerEquivalenceReady {
		t.Fatalf("expected not ready: %#v", readiness)
	}
	assertQueryClusterBlocker(t, readiness, "missing_included_query")
}

func TestEvaluateQueryClusterReadinessRejectsNonSupportingIncludedQuery(t *testing.T) {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	need, intent, cluster := readyQueryClusterParents(now)
	members, intentEvidence := readyQueryClusterMembers(now)
	intentEvidence[1].Link.RelationType = IntentEvidenceContext

	readiness := EvaluateQueryClusterReadiness(need, intent, cluster, members, intentEvidence, now)
	if readiness.Ready || readiness.AnswerEquivalenceReady {
		t.Fatalf("expected not ready: %#v", readiness)
	}
	assertQueryClusterBlocker(t, readiness, "included_query_not_supporting_intent")
}

func TestEvaluateQueryClusterReadinessValidatedRequiresValidatedObservedSet(t *testing.T) {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	need, intent, cluster := readyQueryClusterParents(now)
	cluster.Confidence = ObjectiveConfidenceValidated
	members, intentEvidence := readyQueryClusterMembers(now)
	members[1].Evidence.Confidence = ObjectiveConfidenceSupported
	intentEvidence[1].Evidence.Confidence = ObjectiveConfidenceSupported

	readiness := EvaluateQueryClusterReadiness(need, intent, cluster, members, intentEvidence, now)
	if readiness.Ready {
		t.Fatalf("expected validated gate to block: %#v", readiness)
	}
	assertQueryClusterBlocker(t, readiness, "validated_without_validated_query_set")
}

func readyQueryClusterParents(now time.Time) (*SeoUserNeed, *SeoIntentHypothesis, *SeoQueryCluster) {
	review := now.Add(30 * 24 * time.Hour)
	need := &SeoUserNeed{ID: 20, ObjectiveID: 10, Status: UserNeedStatusActive}
	intent := &SeoIntentHypothesis{
		ID:          30,
		ObjectiveID: 10,
		UserNeedID:  20,
		Status:      IntentHypothesisStatusActive,
	}
	cluster := &SeoQueryCluster{
		ID:                 40,
		ObjectiveID:        10,
		UserNeedID:         20,
		IntentHypothesisID: 30,
		Code:               "parcel_planning_verify",
		Name:               "Parcel planning verification",
		SearchPromise:      "Trả lời liệu thửa đất có chịu tác động quy hoạch và chỉ rõ nguồn dữ liệu.",
		InclusionRule:      "Query có thể được trả lời bằng cùng kết quả xác minh thửa đất và quy hoạch.",
		ExclusionRule:      "Loại query chỉ hỏi giá, giao dịch hoặc pháp lý không liên quan quy hoạch.",
		Status:             QueryClusterStatusDraft,
		Confidence:         ObjectiveConfidenceValidated,
		OwnerCode:          "seo_strategy",
		ReviewDueAt:        &review,
		LockVersion:        3,
	}
	return need, intent, cluster
}

func readyQueryClusterMembers(now time.Time) ([]LinkedQueryClusterMember, []LinkedIntentQueryEvidence) {
	primary := queryClusterEvidence(101, "thửa đất có quy hoạch không", ObjectiveConfidenceValidated, QueryObservationObserved, now)
	included := queryClusterEvidence(102, "kiểm tra quy hoạch thửa đất", ObjectiveConfidenceValidated, QueryObservationReported, now)
	excluded := queryClusterEvidence(103, "giá thửa đất bao nhiêu", ObjectiveConfidenceSupported, QueryObservationObserved, now)

	members := []LinkedQueryClusterMember{
		{
			Member:   SeoQueryClusterMember{ClusterID: 40, QueryEvidenceID: 101, MembershipType: QueryClusterMembershipPrimary, Status: QueryClusterMemberStatusActive},
			Evidence: primary,
		},
		{
			Member:   SeoQueryClusterMember{ClusterID: 40, QueryEvidenceID: 102, MembershipType: QueryClusterMembershipIncluded, Rationale: "Cùng primary answer", Status: QueryClusterMemberStatusActive},
			Evidence: included,
		},
		{
			Member:   SeoQueryClusterMember{ClusterID: 40, QueryEvidenceID: 103, MembershipType: QueryClusterMembershipExcluded, Rationale: "Cần dữ liệu giá, không phải planning answer", Status: QueryClusterMemberStatusActive},
			Evidence: excluded,
		},
	}
	intentEvidence := []LinkedIntentQueryEvidence{
		{Link: SeoIntentEvidenceLink{RelationType: IntentEvidenceSupports}, Evidence: primary},
		{Link: SeoIntentEvidenceLink{RelationType: IntentEvidenceSupports}, Evidence: included},
		{Link: SeoIntentEvidenceLink{RelationType: IntentEvidenceContext}, Evidence: excluded},
	}
	return members, intentEvidence
}

func queryClusterEvidence(
	id uint64,
	raw string,
	confidence ObjectiveConfidence,
	observation QueryObservationType,
	now time.Time,
) SeoQueryEvidence {
	return SeoQueryEvidence{
		ID:              id,
		ObjectiveID:     10,
		UserNeedID:      20,
		RawQuery:        raw,
		NormalizedQuery: NormalizeSearchQuery(raw),
		ObservationType: observation,
		Confidence:      confidence,
		Status:          QueryEvidenceStatusActive,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func assertQueryClusterBlocker(t *testing.T, readiness QueryClusterReadiness, code string) {
	t.Helper()
	for _, item := range readiness.Blockers {
		if item.Code == code {
			return
		}
	}
	t.Fatalf("blocker %q not found in %#v", code, readiness.Blockers)
}
