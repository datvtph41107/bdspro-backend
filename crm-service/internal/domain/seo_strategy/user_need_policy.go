package seo_strategy

import (
	"strings"
	"time"
)

type UserNeedReadinessIssue struct {
	Code    string
	Message string
}

type UserNeedReadiness struct {
	Ready    bool
	Blockers []UserNeedReadinessIssue
	Warnings []UserNeedReadinessIssue
}

func EvaluateUserNeedReadiness(
	objective *SeoObjective,
	audience *SeoAudienceSegment,
	need *SeoUserNeed,
	linkedEvidence []LinkedUserNeedEvidence,
	now time.Time,
) UserNeedReadiness {
	result := UserNeedReadiness{Blockers: []UserNeedReadinessIssue{}, Warnings: []UserNeedReadinessIssue{}}
	block := func(code, message string) {
		result.Blockers = append(result.Blockers, UserNeedReadinessIssue{Code: code, Message: message})
	}
	warn := func(code, message string) {
		result.Warnings = append(result.Warnings, UserNeedReadinessIssue{Code: code, Message: message})
	}

	if objective == nil {
		block("objective_not_found", "Objective không tồn tại")
		return result
	}
	if objective.Status != ObjectiveStatusActive {
		block("objective_not_active", "Objective phải active trước khi activate user need")
	}
	if audience == nil {
		block("audience_not_found", "Audience không tồn tại")
	} else if audience.Status != AudienceStatusActive {
		block("audience_not_active", "Audience phải active trước khi activate user need")
	}
	if need == nil {
		block("user_need_not_found", "User need không tồn tại")
		return result
	}

	require := func(value, code, message string) {
		if strings.TrimSpace(value) == "" {
			block(code, message)
		}
	}
	require(need.Situation, "missing_situation", "Chưa xác định hoàn cảnh phát sinh nhu cầu")
	require(need.TriggerStatement, "missing_trigger", "Chưa xác định sự kiện kích hoạt nhu cầu")
	require(need.JobToBeDone, "missing_job_to_be_done", "Chưa xác định việc người dùng cần hoàn thành")
	require(need.BarrierStatement, "missing_barrier", "Chưa xác định rào cản")
	require(need.DesiredOutcome, "missing_desired_outcome", "Chưa xác định kết quả mong muốn")
	require(need.FailureImpact, "missing_failure_impact", "Chưa xác định hậu quả khi thất bại")
	require(need.OwnerCode, "missing_owner", "Chưa có owner chịu trách nhiệm")

	if need.DecisionStage == UserNeedStageUnknown {
		block("unknown_decision_stage", "Chưa xác định giai đoạn quyết định")
	}
	if need.RiskLevel == UserNeedRiskUnknown {
		block("unknown_risk_level", "Chưa đánh giá mức rủi ro")
	}
	if need.Priority == UserNeedPriorityUnassessed {
		block("priority_unassessed", "Chưa đánh giá priority")
	}

	hasRelevantEvidence := false
	hasNonAssumptionSupport := false
	hasValidatedSupport := false
	for _, linked := range linkedEvidence {
		evidence := linked.Evidence
		if evidence.Status != EvidenceStatusActive {
			continue
		}
		if linked.Link.RelationType == UserNeedEvidenceSupports || linked.Link.RelationType == UserNeedEvidenceContext {
			hasRelevantEvidence = true
		}
		if linked.Link.RelationType == UserNeedEvidenceSupports && evidence.Type != EvidenceTypeAssumption {
			hasNonAssumptionSupport = true
		}
		if linked.Link.RelationType == UserNeedEvidenceSupports &&
			evidence.Type != EvidenceTypeAssumption &&
			evidence.Confidence == ObjectiveConfidenceValidated {
			hasValidatedSupport = true
		}
	}
	if !hasRelevantEvidence {
		block("missing_relevant_evidence", "Chưa có evidence supports hoặc context")
	} else if !hasNonAssumptionSupport {
		warn("assumption_only", "User need hiện chỉ dựa trên assumption")
	}
	if need.Confidence == ObjectiveConfidenceValidated && !hasValidatedSupport {
		block("validated_without_validated_evidence", "Không thể validated khi chưa có validated evidence")
	}
	if strings.TrimSpace(need.CurrentAlternative) == "" {
		warn("missing_current_alternative", "Chưa biết người dùng đang dùng giải pháp nào")
	}
	if need.RiskLevel == UserNeedRiskHigh || need.RiskLevel == UserNeedRiskCritical {
		warn("high_risk_requires_trust_review", "Nhu cầu rủi ro cao cần trust/expert review ở module dữ liệu")
	}
	if need.ReviewDueAt == nil {
		warn("missing_review_date", "Chưa có ngày đánh giá lại user need")
	} else if need.ReviewDueAt.Before(now) {
		warn("review_overdue", "User need đã quá hạn đánh giá lại")
	}
	result.Ready = len(result.Blockers) == 0
	return result
}
