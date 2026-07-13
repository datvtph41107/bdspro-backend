package mapper

import (
	"time"

	seo_strategy "crm/internal/domain/seo_strategy"
	crmpb "pb/types/crm"
)

// SeoAudienceMapper chỉ chuyển đổi domain -> transport.
// Readiness vẫn được tính trong domain policy, không tính lại tại mapper.
type SeoAudienceMapper struct{}

func NewSeoAudienceMapper() *SeoAudienceMapper {
	return &SeoAudienceMapper{}
}

func (m *SeoAudienceMapper) AggregateToPb(
	aggregate *seo_strategy.SeoAudienceAggregate,
) *crmpb.SeoAudienceSegmentResponse {
	if aggregate == nil || aggregate.Audience == nil {
		return nil
	}

	audience := aggregate.Audience
	response := &crmpb.SeoAudienceSegmentResponse{
		Id:             audience.ID,
		Code:           audience.Code,
		Name:           audience.Name,
		Description:    audience.Description,
		DecisionRole:   string(audience.DecisionRole),
		KnowledgeLevel: string(audience.KnowledgeLevel),
		Status:         string(audience.Status),
		OwnerCode:      audience.OwnerCode,
		LockVersion:    audience.LockVersion,
		Readiness:      m.ReadinessToPb(aggregate.Readiness),
		CreatedAt:      formatStrategyTime(audience.CreatedAt),
		UpdatedAt:      formatStrategyTime(audience.UpdatedAt),
	}

	if audience.ReviewDueAt != nil {
		value := formatStrategyTime(*audience.ReviewDueAt)
		response.ReviewDueAt = &value
	}

	return response
}

func (m *SeoAudienceMapper) SummaryToPb(
	audience *seo_strategy.SeoAudienceSegment,
) *crmpb.SeoAudienceSegmentSummaryResponse {
	if audience == nil {
		return nil
	}

	return &crmpb.SeoAudienceSegmentSummaryResponse{
		Id:             audience.ID,
		Code:           audience.Code,
		Name:           audience.Name,
		DecisionRole:   string(audience.DecisionRole),
		KnowledgeLevel: string(audience.KnowledgeLevel),
		Status:         string(audience.Status),
		OwnerCode:      audience.OwnerCode,
		LockVersion:    audience.LockVersion,
		CreatedAt:      formatStrategyTime(audience.CreatedAt),
		UpdatedAt:      formatStrategyTime(audience.UpdatedAt),
	}
}

func (m *SeoAudienceMapper) ReadinessToPb(
	readiness seo_strategy.AudienceReadiness,
) *crmpb.SeoAudienceReadinessResponse {
	response := &crmpb.SeoAudienceReadinessResponse{
		Ready: readiness.Ready,
	}

	for _, issue := range readiness.Blockers {
		response.Blockers = append(
			response.Blockers,
			&crmpb.SeoAudienceReadinessIssueResponse{
				Code:    issue.Code,
				Message: issue.Message,
			},
		)
	}

	for _, issue := range readiness.Warnings {
		response.Warnings = append(
			response.Warnings,
			&crmpb.SeoAudienceReadinessIssueResponse{
				Code:    issue.Code,
				Message: issue.Message,
			},
		)
	}

	return response
}

func formatStrategyTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}
