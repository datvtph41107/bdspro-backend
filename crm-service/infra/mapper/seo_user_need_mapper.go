package mapper

import (
	seo_strategy "crm/internal/domain/seo_strategy"
	crmpb "pb/types/crm"
)

// SeoUserNeedMapper tạo projection dành cho Admin API.
// Nó không quyết định audience, evidence hay readiness.
type SeoUserNeedMapper struct {
	objectiveMapper *SeoObjectiveMapper
	audienceMapper  *SeoAudienceMapper
}

func NewSeoUserNeedMapper(
	objectiveMapper *SeoObjectiveMapper,
	audienceMapper *SeoAudienceMapper,
) *SeoUserNeedMapper {
	return &SeoUserNeedMapper{
		objectiveMapper: objectiveMapper,
		audienceMapper:  audienceMapper,
	}
}

func (m *SeoUserNeedMapper) AggregateToPb(
	aggregate *seo_strategy.SeoUserNeedAggregate,
) *crmpb.SeoUserNeedResponse {
	if aggregate == nil || aggregate.Need == nil {
		return nil
	}

	need := aggregate.Need
	response := &crmpb.SeoUserNeedResponse{
		Id:                 need.ID,
		ObjectiveId:        need.ObjectiveID,
		AudienceSegmentId:  need.AudienceSegmentID,
		Code:               need.Code,
		Name:               need.Name,
		DecisionStage:      string(need.DecisionStage),
		Situation:          need.Situation,
		TriggerStatement:   need.TriggerStatement,
		JobToBeDone:        need.JobToBeDone,
		BarrierStatement:   need.BarrierStatement,
		CurrentAlternative: need.CurrentAlternative,
		DesiredOutcome:     need.DesiredOutcome,
		FailureImpact:      need.FailureImpact,
		RiskLevel:          string(need.RiskLevel),
		Priority:           string(need.Priority),
		Status:             string(need.Status),
		Confidence:         string(need.Confidence),
		OwnerCode:          need.OwnerCode,
		LockVersion:        need.LockVersion,
		Readiness:          m.ReadinessToPb(aggregate.Readiness),
		CreatedAt:          formatStrategyTime(need.CreatedAt),
		UpdatedAt:          formatStrategyTime(need.UpdatedAt),
	}

	if need.ReviewDueAt != nil {
		value := formatStrategyTime(*need.ReviewDueAt)
		response.ReviewDueAt = &value
	}

	if aggregate.Objective != nil {
		response.Objective = m.objectiveMapper.SummaryToPb(
			aggregate.Objective,
		)
	}

	if aggregate.Audience != nil {
		response.Audience = m.audienceMapper.SummaryToPb(
			aggregate.Audience,
		)
	}

	for i := range aggregate.Evidence {
		response.Evidence = append(
			response.Evidence,
			m.LinkedEvidenceToPb(&aggregate.Evidence[i]),
		)
	}

	return response
}

func (m *SeoUserNeedMapper) SummaryToPb(
	need *seo_strategy.SeoUserNeed,
) *crmpb.SeoUserNeedSummaryResponse {
	if need == nil {
		return nil
	}

	return &crmpb.SeoUserNeedSummaryResponse{
		Id:                need.ID,
		ObjectiveId:       need.ObjectiveID,
		AudienceSegmentId: need.AudienceSegmentID,
		Code:              need.Code,
		Name:              need.Name,
		DecisionStage:     string(need.DecisionStage),
		RiskLevel:         string(need.RiskLevel),
		Priority:          string(need.Priority),
		Status:            string(need.Status),
		Confidence:        string(need.Confidence),
		OwnerCode:         need.OwnerCode,
		LockVersion:       need.LockVersion,
		CreatedAt:         formatStrategyTime(need.CreatedAt),
		UpdatedAt:         formatStrategyTime(need.UpdatedAt),
	}
}

func (m *SeoUserNeedMapper) LinkedEvidenceToPb(
	linked *seo_strategy.LinkedUserNeedEvidence,
) *crmpb.SeoUserNeedLinkedEvidenceResponse {
	if linked == nil {
		return nil
	}

	evidence := linked.Evidence
	response := &crmpb.SeoUserNeedLinkedEvidenceResponse{
		EvidenceId:      evidence.ID,
		RelationType:    string(linked.Link.RelationType),
		RelevanceNote:   linked.Link.RelevanceNote,
		EvidenceType:    string(evidence.Type),
		Stance:          string(evidence.Stance),
		SourceSystem:    evidence.SourceSystem,
		SourceReference: evidence.SourceReference,
		Finding:         evidence.Finding,
		Confidence:      string(evidence.Confidence),
		Status:          string(evidence.Status),
		CreatedAt:       formatStrategyTime(linked.Link.CreatedAt),
	}

	if evidence.ObservedAt != nil {
		value := formatStrategyTime(*evidence.ObservedAt)
		response.ObservedAt = &value
	}

	if evidence.ValidUntil != nil {
		value := formatStrategyTime(*evidence.ValidUntil)
		response.ValidUntil = &value
	}

	return response
}

func (m *SeoUserNeedMapper) ReadinessToPb(
	readiness seo_strategy.UserNeedReadiness,
) *crmpb.SeoUserNeedReadinessResponse {
	response := &crmpb.SeoUserNeedReadinessResponse{
		Ready: readiness.Ready,
	}

	for _, issue := range readiness.Blockers {
		response.Blockers = append(
			response.Blockers,
			&crmpb.SeoUserNeedReadinessIssueResponse{
				Code:    issue.Code,
				Message: issue.Message,
			},
		)
	}

	for _, issue := range readiness.Warnings {
		response.Warnings = append(
			response.Warnings,
			&crmpb.SeoUserNeedReadinessIssueResponse{
				Code:    issue.Code,
				Message: issue.Message,
			},
		)
	}

	return response
}
