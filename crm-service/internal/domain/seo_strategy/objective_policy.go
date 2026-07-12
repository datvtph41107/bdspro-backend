package seo_strategy

import (
	"strings"
	"time"
)

type ObjectiveReadinessIssue struct {
	Code    string
	Message string
}

type ObjectiveReadiness struct {
	Ready    bool
	Blockers []ObjectiveReadinessIssue
	Warnings []ObjectiveReadinessIssue
}

func EvaluateObjectiveReadiness(
	objective *SeoObjective,
	evidence []SeoObjectiveEvidence,
	metrics []SeoObjectiveMetric,
	now time.Time,
) ObjectiveReadiness {
	result := ObjectiveReadiness{
		Blockers: []ObjectiveReadinessIssue{},
		Warnings: []ObjectiveReadinessIssue{},
	}

	if objective == nil {
		result.Blockers = append(
			result.Blockers,
			ObjectiveReadinessIssue{
				Code:    "objective_not_found",
				Message: "Objective không tồn tại",
			},
		)
		result.Ready = false
		return result
	}

	requireText := func(
		value string,
		code string,
		message string,
	) {
		if strings.TrimSpace(value) == "" {
			result.Blockers = append(
				result.Blockers,
				ObjectiveReadinessIssue{
					Code:    code,
					Message: message,
				},
			)
		}
	}

	requireText(
		objective.AudienceSummary,
		"missing_audience_summary",
		"Chưa xác định nhóm người dùng trọng tâm",
	)

	requireText(
		objective.Situation,
		"missing_situation",
		"Chưa xác định hoàn cảnh phát sinh nhu cầu",
	)

	requireText(
		objective.ProblemStatement,
		"missing_problem_statement",
		"Chưa mô tả vấn đề thực tế",
	)

	requireText(
		objective.UserOutcome,
		"missing_user_outcome",
		"Chưa xác định kết quả dành cho người dùng",
	)

	requireText(
		objective.ProductOutcome,
		"missing_product_outcome",
		"Chưa xác định kết quả sản phẩm",
	)

	requireText(
		objective.BusinessOutcome,
		"missing_business_outcome",
		"Chưa xác định giá trị kinh doanh",
	)

	requireText(
		objective.OwnerCode,
		"missing_owner",
		"Chưa có owner chịu trách nhiệm",
	)

	activeEvidence := make([]SeoObjectiveEvidence, 0, len(evidence))
	hasNonAssumptionEvidence := false
	hasValidatedEvidence := false

	for _, item := range evidence {
		if item.Status != EvidenceStatusActive {
			continue
		}

		activeEvidence = append(activeEvidence, item)

		if item.Type != EvidenceTypeAssumption {
			hasNonAssumptionEvidence = true
		}

		if item.Type != EvidenceTypeAssumption &&
			item.Confidence == ObjectiveConfidenceValidated {
			hasValidatedEvidence = true
		}
	}

	if len(activeEvidence) == 0 {
		result.Blockers = append(
			result.Blockers,
			ObjectiveReadinessIssue{
				Code:    "missing_evidence",
				Message: "Chưa có evidence hoặc assumption được ghi nhận",
			},
		)
	} else if !hasNonAssumptionEvidence {
		result.Warnings = append(
			result.Warnings,
			ObjectiveReadinessIssue{
				Code:    "assumption_only",
				Message: "Objective hiện chỉ dựa trên assumption",
			},
		)
	}

	primaryMetrics := make([]SeoObjectiveMetric, 0, 1)

	for _, metric := range metrics {
		if metric.Status == MetricStatusActive && metric.IsPrimary {
			primaryMetrics = append(primaryMetrics, metric)
		}
	}

	if len(primaryMetrics) == 0 {
		result.Blockers = append(
			result.Blockers,
			ObjectiveReadinessIssue{
				Code:    "missing_primary_metric",
				Message: "Chưa có primary metric",
			},
		)
	}

	if len(primaryMetrics) > 1 {
		result.Blockers = append(
			result.Blockers,
			ObjectiveReadinessIssue{
				Code:    "multiple_primary_metrics",
				Message: "Objective có nhiều hơn một primary metric",
			},
		)
	}

	if len(primaryMetrics) == 1 {
		primary := primaryMetrics[0]

		if strings.TrimSpace(primary.SourceSystem) == "" {
			result.Blockers = append(
				result.Blockers,
				ObjectiveReadinessIssue{
					Code:    "primary_metric_missing_source",
					Message: "Primary metric chưa có source system",
				},
			)
		}

		if strings.TrimSpace(primary.CalculationMethod) == "" {
			result.Blockers = append(
				result.Blockers,
				ObjectiveReadinessIssue{
					Code:    "primary_metric_missing_calculation",
					Message: "Primary metric chưa có cách tính",
				},
			)
		}

		if primary.InstrumentationStatus ==
			InstrumentationStatusUnknown {
			result.Blockers = append(
				result.Blockers,
				ObjectiveReadinessIssue{
					Code:    "primary_metric_instrumentation_unknown",
					Message: "Chưa biết primary metric sẽ được đo bằng cách nào",
				},
			)
		}

		if primary.InstrumentationStatus ==
			InstrumentationStatusPlanned {
			result.Warnings = append(
				result.Warnings,
				ObjectiveReadinessIssue{
					Code:    "primary_metric_not_instrumented",
					Message: "Primary metric mới ở trạng thái planned",
				},
			)
		}
	}

	if objective.Confidence == ObjectiveConfidenceValidated &&
		!hasValidatedEvidence {
		result.Blockers = append(
			result.Blockers,
			ObjectiveReadinessIssue{
				Code:    "validated_without_validated_evidence",
				Message: "Không thể đánh dấu validated khi chưa có evidence phù hợp",
			},
		)
	}

	if len(objective.NonGoals) == 0 {
		result.Warnings = append(
			result.Warnings,
			ObjectiveReadinessIssue{
				Code:    "missing_non_goals",
				Message: "Chưa xác định các mục tiêu không thuộc phạm vi",
			},
		)
	}

	if objective.ReviewDueAt == nil {
		result.Warnings = append(
			result.Warnings,
			ObjectiveReadinessIssue{
				Code:    "missing_review_date",
				Message: "Chưa có ngày đánh giá lại objective",
			},
		)
	} else if objective.ReviewDueAt.Before(now) {
		result.Warnings = append(
			result.Warnings,
			ObjectiveReadinessIssue{
				Code:    "review_overdue",
				Message: "Objective đã quá hạn đánh giá lại",
			},
		)
	}

	result.Ready = len(result.Blockers) == 0
	return result
}
