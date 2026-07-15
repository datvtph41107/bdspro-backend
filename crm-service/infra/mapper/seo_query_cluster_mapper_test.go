package mapper

import (
	"testing"
	"time"

	seo_strategy "crm/internal/domain/seo_strategy"
)

func TestSeoQueryClusterMapperAggregateToPb(t *testing.T) {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	review := now.Add(30 * 24 * time.Hour)
	cluster := &seo_strategy.SeoQueryCluster{
		ID: 1, ObjectiveID: 10, UserNeedID: 20, IntentHypothesisID: 30,
		Code: "parcel_planning_verify", Name: "Parcel planning verify",
		SearchPromise: "Source-backed planning answer",
		InclusionRule: "Same primary answer", ExclusionRule: "Different answer",
		Status: seo_strategy.QueryClusterStatusDraft, Confidence: seo_strategy.ObjectiveConfidenceSupported,
		OwnerCode: "seo", ReviewDueAt: &review, LockVersion: 4,
		CreatedAt: now, UpdatedAt: now,
	}
	evidence := seo_strategy.SeoQueryEvidence{
		ID: 101, ObjectiveID: 10, UserNeedID: 20,
		RawQuery: "parcel planning", NormalizedQuery: "parcel planning",
		Status:    seo_strategy.QueryEvidenceStatusActive,
		CreatedAt: now, UpdatedAt: now,
	}
	intentMapper := NewSeoIntentHypothesisMapper(&SeoUserNeedMapper{}, NewSeoQueryEvidenceMapper())
	mapper := NewSeoQueryClusterMapper(intentMapper, &SeoUserNeedMapper{}, NewSeoQueryEvidenceMapper())
	response := mapper.AggregateToPb(&seo_strategy.SeoQueryClusterAggregate{
		Cluster: cluster,
		Intent: &seo_strategy.SeoIntentHypothesis{
			ID: 30, ObjectiveID: 10, UserNeedID: 20, Code: "verify", Name: "Verify",
			SearchAction: seo_strategy.IntentSearchActionVerify,
			SubjectType:  seo_strategy.IntentSubjectParcel,
			ScopeType:    seo_strategy.IntentScopeSpecificEntity,
			Status:       seo_strategy.IntentHypothesisStatusActive,
			Confidence:   seo_strategy.ObjectiveConfidenceSupported,
			LockVersion:  2, CreatedAt: now, UpdatedAt: now,
		},
		UserNeed: &seo_strategy.SeoUserNeed{ID: 20, ObjectiveID: 10},
		Members: []seo_strategy.LinkedQueryClusterMember{{
			Member: seo_strategy.SeoQueryClusterMember{
				ClusterID: 1, QueryEvidenceID: 101,
				MembershipType: seo_strategy.QueryClusterMembershipPrimary,
				Status:         seo_strategy.QueryClusterMemberStatusActive,
				CreatedAt:      now, UpdatedAt: now,
			},
			Evidence: evidence,
		}},
		Readiness: seo_strategy.QueryClusterReadiness{
			Ready: true, AnswerEquivalenceReady: true,
			PrimaryCount: 1, IncludedCount: 1, ExcludedCount: 1,
		},
	})
	if response == nil {
		t.Fatal("response must not be nil")
	}
	if response.Id != 1 || response.IntentHypothesisId != 30 || response.LockVersion != 4 {
		t.Fatalf("unexpected response: %+v", response)
	}
	if response.IntentHypothesis == nil || response.IntentHypothesis.Id != 30 {
		t.Fatalf("unexpected intent projection: %+v", response.IntentHypothesis)
	}
	if len(response.Members) != 1 || response.Members[0].QueryEvidence.GetId() != 101 {
		t.Fatalf("unexpected members: %+v", response.Members)
	}
	if response.Readiness == nil || !response.Readiness.AnswerEquivalenceReady {
		t.Fatalf("unexpected readiness: %+v", response.Readiness)
	}
}
