package seo_strategy

import (
	"fmt"
	"strings"
	"time"
)

type IntentHypothesisReadinessIssue struct {
	Code    string
	Message string
}

type IntentHypothesisReadiness struct {
	Ready    bool
	Blockers []IntentHypothesisReadinessIssue
	Warnings []IntentHypothesisReadinessIssue
}

func EvaluateIntentHypothesisReadiness(
	need *SeoUserNeed,
	intent *SeoIntentHypothesis,
	linkedEvidence []LinkedIntentQueryEvidence,
	now time.Time,
) IntentHypothesisReadiness {
	result := IntentHypothesisReadiness{
		Blockers: []IntentHypothesisReadinessIssue{},
		Warnings: []IntentHypothesisReadinessIssue{},
	}
	block := func(code, message string) {
		result.Blockers = append(result.Blockers, IntentHypothesisReadinessIssue{Code: code, Message: message})
	}
	warn := func(code, message string) {
		result.Warnings = append(result.Warnings, IntentHypothesisReadinessIssue{Code: code, Message: message})
	}

	if need == nil {
		block("user_need_not_found", "User Need không tồn tại")
		return result
	}
	if need.Status != UserNeedStatusActive {
		block("user_need_not_active", "User Need phải active trước khi activate Intent Hypothesis")
	}
	if intent == nil {
		block("intent_not_found", "Intent Hypothesis không tồn tại")
		return result
	}
	if intent.UserNeedID != need.ID || intent.ObjectiveID != need.ObjectiveID {
		block("intent_user_need_mismatch", "Intent Hypothesis không thuộc đúng User Need và Objective")
	}

	require := func(value, code, message string) {
		if strings.TrimSpace(value) == "" {
			block(code, message)
		}
	}
	require(intent.GoalStatement, "missing_goal_statement", "Chưa mô tả mục tiêu tìm kiếm")
	require(intent.PrimaryQuestion, "missing_primary_question", "Chưa xác định câu hỏi chính")
	require(intent.ExpectedResult, "missing_expected_result", "Chưa xác định kết quả người dùng mong đợi")
	require(intent.OwnerCode, "missing_owner", "Chưa có owner chịu trách nhiệm")

	if intent.SearchAction == IntentSearchActionUnknown {
		block("unknown_search_action", "Chưa xác định hành động tìm kiếm")
	}
	if intent.SubjectType == IntentSubjectUnknown {
		block("unknown_subject_type", "Chưa xác định loại chủ thể tìm kiếm")
	}
	if intent.ScopeType == IntentScopeUnknown {
		block("unknown_scope_type", "Chưa xác định phạm vi tìm kiếm")
	}

	supportCount := 0
	observedSupportCount := 0
	validatedObservedSupportCount := 0
	contradictionCount := 0
	contextCount := 0

	for _, linked := range linkedEvidence {
		evidence := linked.Evidence
		if !IsQueryEvidenceCurrent(evidence, now) {
			continue
		}
		switch linked.Link.RelationType {
		case IntentEvidenceSupports:
			supportCount++
			if IsQueryEvidenceObserved(evidence) {
				observedSupportCount++
				if evidence.Confidence == ObjectiveConfidenceValidated {
					validatedObservedSupportCount++
				}
			}
		case IntentEvidenceContradicts:
			contradictionCount++
		case IntentEvidenceContext:
			contextCount++
		}
	}

	if supportCount == 0 {
		block("missing_supporting_query_evidence", "Chưa có Query Evidence active hỗ trợ Intent Hypothesis")
	} else if observedSupportCount == 0 {
		warn("hypothesis_only_support", "Intent hiện chỉ được hỗ trợ bởi suggested hoặc assumed query evidence")
	}
	if contradictionCount > 0 {
		warn("contradicting_query_evidence", "Có Query Evidence đang mâu thuẫn với Intent Hypothesis")
	}
	if contextCount == 0 {
		warn("missing_context_evidence", "Chưa có Query Evidence context để mô tả ranh giới của intent")
	}
	if intent.Confidence == ObjectiveConfidenceSupported && observedSupportCount == 0 {
		block("supported_without_observed_evidence", "Intent supported phải có ít nhất một observed hoặc reported Query Evidence")
	}
	if intent.Confidence == ObjectiveConfidenceValidated && validatedObservedSupportCount == 0 {
		block("validated_without_validated_evidence", "Intent validated phải có observed hoặc reported Query Evidence đã validated")
	}
	if intent.ReviewDueAt == nil {
		warn("missing_review_due_at", "Chưa đặt thời điểm review Intent Hypothesis")
	} else if intent.ReviewDueAt.Before(now) {
		warn("review_overdue", "Intent Hypothesis đã quá hạn review")
	}

	result.Ready = len(result.Blockers) == 0
	return result
}

func ValidateIntentHypothesisStatusTransition(
	current IntentHypothesisStatus,
	target IntentHypothesisStatus,
) error {
	if !IsValidIntentHypothesisStatus(current) {
		return fmt.Errorf("intent hypothesis current status không hợp lệ: %s", current)
	}
	if !IsValidIntentHypothesisStatus(target) {
		return fmt.Errorf("intent hypothesis target status không hợp lệ: %s", target)
	}
	switch current {
	case IntentHypothesisStatusDraft:
		if target == IntentHypothesisStatusActive || target == IntentHypothesisStatusArchived {
			return nil
		}
	case IntentHypothesisStatusActive:
		if target == IntentHypothesisStatusPaused || target == IntentHypothesisStatusArchived {
			return nil
		}
	case IntentHypothesisStatusPaused:
		if target == IntentHypothesisStatusActive || target == IntentHypothesisStatusArchived {
			return nil
		}
	case IntentHypothesisStatusArchived:
		// Archived là terminal trong Module 3.
	}
	return fmt.Errorf("không thể chuyển intent hypothesis từ %s sang %s", current, target)
}
