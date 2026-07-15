package seo_strategy

import (
	"fmt"
	"strings"
	"time"
)

type QueryClusterReadinessIssue struct {
	Code    string
	Message string
}

type QueryClusterReadiness struct {
	Ready                  bool
	AnswerEquivalenceReady bool
	PrimaryCount           uint32
	IncludedCount          uint32
	ExcludedCount          uint32
	Blockers               []QueryClusterReadinessIssue
	Warnings               []QueryClusterReadinessIssue
}

func EvaluateQueryClusterReadiness(
	need *SeoUserNeed,
	intent *SeoIntentHypothesis,
	cluster *SeoQueryCluster,
	members []LinkedQueryClusterMember,
	intentEvidence []LinkedIntentQueryEvidence,
	now time.Time,
) QueryClusterReadiness {
	result := QueryClusterReadiness{
		Blockers: []QueryClusterReadinessIssue{},
		Warnings: []QueryClusterReadinessIssue{},
	}
	block := func(code, message string) {
		result.Blockers = append(result.Blockers, QueryClusterReadinessIssue{Code: code, Message: message})
	}
	warn := func(code, message string) {
		result.Warnings = append(result.Warnings, QueryClusterReadinessIssue{Code: code, Message: message})
	}

	if need == nil {
		block("user_need_not_found", "User Need không tồn tại")
		return result
	}
	if need.Status != UserNeedStatusActive {
		block("user_need_not_active", "User Need phải active trước khi activate Query Cluster")
	}
	if intent == nil {
		block("intent_not_found", "Intent Hypothesis không tồn tại")
		return result
	}
	if intent.Status != IntentHypothesisStatusActive {
		block("intent_not_active", "Intent Hypothesis phải active trước khi activate Query Cluster")
	}
	if cluster == nil {
		block("cluster_not_found", "Query Cluster không tồn tại")
		return result
	}
	if cluster.UserNeedID != need.ID || cluster.ObjectiveID != need.ObjectiveID ||
		cluster.IntentHypothesisID != intent.ID || intent.UserNeedID != need.ID ||
		intent.ObjectiveID != need.ObjectiveID {
		block("cluster_parent_mismatch", "Query Cluster không thuộc đúng Intent Hypothesis, User Need và Objective")
	}

	require := func(value, code, message string) {
		if strings.TrimSpace(value) == "" {
			block(code, message)
		}
	}
	require(cluster.SearchPromise, "missing_search_promise", "Chưa mô tả Search Promise")
	require(cluster.InclusionRule, "missing_inclusion_rule", "Chưa định nghĩa quy tắc query được phục vụ cùng primary answer")
	require(cluster.ExclusionRule, "missing_exclusion_rule", "Chưa định nghĩa ranh giới query không thuộc cluster")
	require(cluster.OwnerCode, "missing_owner", "Chưa có owner chịu trách nhiệm")

	intentRelations := make(map[uint64]IntentEvidenceRelation, len(intentEvidence))
	for _, linked := range intentEvidence {
		if linked.Evidence.Status == QueryEvidenceStatusActive {
			intentRelations[linked.Evidence.ID] = linked.Link.RelationType
		}
	}

	observedIncludedCount := 0
	validatedObservedIncludedCount := 0
	invalidIncludedRelation := 0
	staleIncludedCount := 0

	for _, linked := range members {
		member := linked.Member
		if member.Status != QueryClusterMemberStatusActive {
			continue
		}
		switch member.MembershipType {
		case QueryClusterMembershipPrimary:
			result.PrimaryCount++
		case QueryClusterMembershipIncluded:
			result.IncludedCount++
		case QueryClusterMembershipExcluded:
			result.ExcludedCount++
		}

		relation, linkedToIntent := intentRelations[linked.Evidence.ID]
		if !linkedToIntent {
			block("member_not_linked_to_intent", "Có Query Evidence thành viên chưa được liên kết với Intent Hypothesis")
			continue
		}

		if member.MembershipType == QueryClusterMembershipPrimary ||
			member.MembershipType == QueryClusterMembershipIncluded {
			if relation != IntentEvidenceSupports {
				invalidIncludedRelation++
			}
			if !IsQueryEvidenceCurrent(linked.Evidence, now) {
				staleIncludedCount++
				continue
			}
			if IsQueryEvidenceObserved(linked.Evidence) {
				observedIncludedCount++
				if linked.Evidence.Confidence == ObjectiveConfidenceValidated {
					validatedObservedIncludedCount++
				}
			}
		}
	}

	if result.PrimaryCount == 0 {
		block("missing_primary_query", "Cluster phải có đúng một primary query")
	} else if result.PrimaryCount > 1 {
		block("multiple_primary_queries", "Cluster chỉ được có một primary query active")
	}
	if result.IncludedCount == 0 {
		block("missing_included_query", "Cluster phải có ít nhất một included query để chứng minh answer equivalence")
	}
	if invalidIncludedRelation > 0 {
		block("included_query_not_supporting_intent", "Primary và included query phải là supporting evidence của Intent Hypothesis")
	}
	if staleIncludedCount > 0 {
		block("stale_or_inactive_cluster_member", "Primary hoặc included query phải active và còn hiệu lực")
	}
	if result.ExcludedCount == 0 {
		warn("missing_excluded_boundary_example", "Chưa có excluded query minh họa ranh giới của cluster")
	}

	result.AnswerEquivalenceReady = result.PrimaryCount == 1 &&
		result.IncludedCount > 0 &&
		invalidIncludedRelation == 0 &&
		staleIncludedCount == 0

	if cluster.Confidence == ObjectiveConfidenceSupported && observedIncludedCount < 2 {
		block("supported_without_observed_query_set", "Cluster supported cần ít nhất hai observed hoặc reported query trong primary answer set")
	}
	if cluster.Confidence == ObjectiveConfidenceValidated && validatedObservedIncludedCount < 2 {
		block("validated_without_validated_query_set", "Cluster validated cần ít nhất hai observed hoặc reported query đã validated")
	}
	if observedIncludedCount == 0 && result.PrimaryCount+result.IncludedCount > 0 {
		warn("hypothesis_only_cluster", "Cluster hiện chỉ dựa trên suggested hoặc assumed Query Evidence")
	}
	if cluster.ReviewDueAt == nil {
		warn("missing_review_due_at", "Chưa đặt thời điểm review Query Cluster")
	} else if cluster.ReviewDueAt.Before(now) {
		warn("review_overdue", "Query Cluster đã quá hạn review")
	}

	result.Ready = len(result.Blockers) == 0 && result.AnswerEquivalenceReady
	return result
}

func ValidateQueryClusterStatusTransition(current QueryClusterStatus, target QueryClusterStatus) error {
	if !IsValidQueryClusterStatus(current) {
		return fmt.Errorf("query cluster current status không hợp lệ: %s", current)
	}
	if !IsValidQueryClusterStatus(target) {
		return fmt.Errorf("query cluster target status không hợp lệ: %s", target)
	}
	switch current {
	case QueryClusterStatusDraft:
		if target == QueryClusterStatusActive || target == QueryClusterStatusArchived {
			return nil
		}
	case QueryClusterStatusActive:
		if target == QueryClusterStatusPaused || target == QueryClusterStatusArchived {
			return nil
		}
	case QueryClusterStatusPaused:
		if target == QueryClusterStatusActive || target == QueryClusterStatusArchived {
			return nil
		}
	case QueryClusterStatusArchived:
		// Archived là terminal trong Module 3.
	}
	return fmt.Errorf("không thể chuyển query cluster từ %s sang %s", current, target)
}

func ValidateQueryClusterMembershipAgainstIntent(
	membershipType QueryClusterMembershipType,
	relation IntentEvidenceRelation,
) error {
	if !IsValidQueryClusterMembershipType(membershipType) {
		return fmt.Errorf("query cluster membership type không hợp lệ: %s", membershipType)
	}
	if !IsValidIntentEvidenceRelation(relation) {
		return fmt.Errorf("intent evidence relation không hợp lệ: %s", relation)
	}
	if membershipType == QueryClusterMembershipPrimary || membershipType == QueryClusterMembershipIncluded {
		if relation != IntentEvidenceSupports {
			return fmt.Errorf("primary hoặc included query phải có relation supports với Intent Hypothesis")
		}
	}
	return nil
}
