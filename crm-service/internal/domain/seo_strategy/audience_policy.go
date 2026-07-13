package seo_strategy

import (
	"strings"
	"time"
)

type AudienceReadinessIssue struct {
	Code    string
	Message string
}

type AudienceReadiness struct {
	Ready    bool
	Blockers []AudienceReadinessIssue
	Warnings []AudienceReadinessIssue
}

func EvaluateAudienceReadiness(audience *SeoAudienceSegment, now time.Time) AudienceReadiness {
	result := AudienceReadiness{Blockers: []AudienceReadinessIssue{}, Warnings: []AudienceReadinessIssue{}}
	block := func(code, message string) {
		result.Blockers = append(result.Blockers, AudienceReadinessIssue{Code: code, Message: message})
	}
	warn := func(code, message string) {
		result.Warnings = append(result.Warnings, AudienceReadinessIssue{Code: code, Message: message})
	}

	if audience == nil {
		block("audience_not_found", "Audience không tồn tại")
		return result
	}
	if strings.TrimSpace(audience.Description) == "" {
		block("missing_description", "Chưa mô tả ranh giới audience")
	}
	if audience.DecisionRole == AudienceDecisionRoleUnknown {
		block("unknown_decision_role", "Chưa xác định vai trò trong quyết định")
	}
	if audience.KnowledgeLevel == AudienceKnowledgeUnknown {
		block("unknown_knowledge_level", "Chưa xác định mức hiểu biết")
	}
	if strings.TrimSpace(audience.OwnerCode) == "" {
		block("missing_owner", "Chưa có owner chịu trách nhiệm")
	}
	if audience.ReviewDueAt == nil {
		warn("missing_review_date", "Chưa có ngày đánh giá lại audience")
	} else if audience.ReviewDueAt.Before(now) {
		warn("review_overdue", "Audience đã quá hạn đánh giá lại")
	}
	result.Ready = len(result.Blockers) == 0
	return result
}
