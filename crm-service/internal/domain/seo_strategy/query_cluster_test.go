package seo_strategy

import (
	"testing"
	"time"
)

func TestCreateSeoQueryClusterDraft(t *testing.T) {
	now := time.Date(2026, 7, 14, 10, 0, 0, 0, time.UTC)
	cluster, err := CreateSeoQueryClusterDraft(10, 20, 30, " Parcel Planning Verify ", " Parcel planning verification ", now)
	if err != nil {
		t.Fatalf("CreateSeoQueryClusterDraft() error = %v", err)
	}
	if cluster.Code != "parcel_planning_verify" {
		t.Fatalf("Code = %q", cluster.Code)
	}
	if cluster.Status != QueryClusterStatusDraft || cluster.LockVersion != 1 {
		t.Fatalf("unexpected lifecycle: status=%s version=%d", cluster.Status, cluster.LockVersion)
	}
	if cluster.ObjectiveID != 10 || cluster.UserNeedID != 20 || cluster.IntentHypothesisID != 30 {
		t.Fatalf("unexpected parent identity: %#v", cluster)
	}
}

func TestValidateQueryClusterMemberRequiresRationaleForIncludedAndExcluded(t *testing.T) {
	for _, membershipType := range []QueryClusterMembershipType{
		QueryClusterMembershipIncluded,
		QueryClusterMembershipExcluded,
	} {
		member := &SeoQueryClusterMember{
			ObjectiveID:     1,
			UserNeedID:      2,
			ClusterID:       3,
			QueryEvidenceID: 4,
			MembershipType:  membershipType,
			Status:          QueryClusterMemberStatusActive,
		}
		if err := ValidateQueryClusterMember(member); err == nil {
			t.Fatalf("expected rationale error for %s", membershipType)
		}
	}
}

func TestValidateQueryClusterMembershipAgainstIntent(t *testing.T) {
	if err := ValidateQueryClusterMembershipAgainstIntent(
		QueryClusterMembershipPrimary,
		IntentEvidenceContradicts,
	); err == nil {
		t.Fatal("expected primary/contradicts to be rejected")
	}
	if err := ValidateQueryClusterMembershipAgainstIntent(
		QueryClusterMembershipExcluded,
		IntentEvidenceContradicts,
	); err != nil {
		t.Fatalf("excluded/contradicts should be valid: %v", err)
	}
}

func TestValidateQueryClusterStatusTransition(t *testing.T) {
	valid := [][2]QueryClusterStatus{
		{QueryClusterStatusDraft, QueryClusterStatusActive},
		{QueryClusterStatusDraft, QueryClusterStatusArchived},
		{QueryClusterStatusActive, QueryClusterStatusPaused},
		{QueryClusterStatusPaused, QueryClusterStatusActive},
	}
	for _, item := range valid {
		if err := ValidateQueryClusterStatusTransition(item[0], item[1]); err != nil {
			t.Fatalf("transition %s -> %s error = %v", item[0], item[1], err)
		}
	}
	if err := ValidateQueryClusterStatusTransition(QueryClusterStatusArchived, QueryClusterStatusActive); err == nil {
		t.Fatal("archived must be terminal")
	}
}
