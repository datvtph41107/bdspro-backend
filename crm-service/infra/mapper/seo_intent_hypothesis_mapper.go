package mapper

import (
	seo_strategy "crm/internal/domain/seo_strategy"
	crmpb "pb/types/crm"
)

type SeoIntentHypothesisMapper struct {
	userNeedMapper      *SeoUserNeedMapper
	queryEvidenceMapper *SeoQueryEvidenceMapper
}

func NewSeoIntentHypothesisMapper(
	userNeedMapper *SeoUserNeedMapper,
	queryEvidenceMapper *SeoQueryEvidenceMapper,
) *SeoIntentHypothesisMapper {
	return &SeoIntentHypothesisMapper{
		userNeedMapper:      userNeedMapper,
		queryEvidenceMapper: queryEvidenceMapper,
	}
}

func (m *SeoIntentHypothesisMapper) AggregateToPb(
	aggregate *seo_strategy.SeoIntentHypothesisAggregate,
) *crmpb.SeoIntentHypothesisResponse {
	if aggregate == nil || aggregate.Intent == nil {
		return nil
	}

	intent := aggregate.Intent
	response := &crmpb.SeoIntentHypothesisResponse{
		Id:              intent.ID,
		ObjectiveId:     intent.ObjectiveID,
		UserNeedId:      intent.UserNeedID,
		Code:            intent.Code,
		Name:            intent.Name,
		GoalStatement:   intent.GoalStatement,
		PrimaryQuestion: intent.PrimaryQuestion,
		ExpectedResult:  intent.ExpectedResult,
		SearchAction:    string(intent.SearchAction),
		SubjectType:     string(intent.SubjectType),
		ScopeType:       string(intent.ScopeType),
		Status:          string(intent.Status),
		Confidence:      string(intent.Confidence),
		OwnerCode:       intent.OwnerCode,
		LockVersion:     intent.LockVersion,
		Readiness:       m.ReadinessToPb(aggregate.Readiness),
		CreatedAt:       formatStrategyTime(intent.CreatedAt),
		UpdatedAt:       formatStrategyTime(intent.UpdatedAt),
	}

	if intent.ReviewDueAt != nil {
		value := formatStrategyTime(*intent.ReviewDueAt)
		response.ReviewDueAt = &value
	}
	if aggregate.UserNeed != nil {
		response.UserNeed = m.userNeedMapper.SummaryToPb(aggregate.UserNeed)
	}
	for i := range aggregate.Evidence {
		response.Evidence = append(
			response.Evidence,
			m.LinkedEvidenceToPb(&aggregate.Evidence[i]),
		)
	}
	return response
}

func (m *SeoIntentHypothesisMapper) SummaryToPb(
	intent *seo_strategy.SeoIntentHypothesis,
) *crmpb.SeoIntentHypothesisSummaryResponse {
	if intent == nil {
		return nil
	}
	return &crmpb.SeoIntentHypothesisSummaryResponse{
		Id:           intent.ID,
		ObjectiveId:  intent.ObjectiveID,
		UserNeedId:   intent.UserNeedID,
		Code:         intent.Code,
		Name:         intent.Name,
		SearchAction: string(intent.SearchAction),
		SubjectType:  string(intent.SubjectType),
		ScopeType:    string(intent.ScopeType),
		Status:       string(intent.Status),
		Confidence:   string(intent.Confidence),
		OwnerCode:    intent.OwnerCode,
		LockVersion:  intent.LockVersion,
		CreatedAt:    formatStrategyTime(intent.CreatedAt),
		UpdatedAt:    formatStrategyTime(intent.UpdatedAt),
	}
}

func (m *SeoIntentHypothesisMapper) LinkedEvidenceToPb(
	linked *seo_strategy.LinkedIntentQueryEvidence,
) *crmpb.SeoIntentLinkedEvidenceResponse {
	if linked == nil {
		return nil
	}
	return &crmpb.SeoIntentLinkedEvidenceResponse{
		RelationType:  string(linked.Link.RelationType),
		RelevanceNote: linked.Link.RelevanceNote,
		LinkedAt:      formatStrategyTime(linked.Link.CreatedAt),
		QueryEvidence: m.queryEvidenceMapper.ToPb(&linked.Evidence),
	}
}

func (m *SeoIntentHypothesisMapper) ReadinessToPb(
	readiness seo_strategy.IntentHypothesisReadiness,
) *crmpb.SeoIntentHypothesisReadinessResponse {
	response := &crmpb.SeoIntentHypothesisReadinessResponse{
		Ready: readiness.Ready,
	}
	for _, issue := range readiness.Blockers {
		response.Blockers = append(
			response.Blockers,
			&crmpb.SeoIntentHypothesisReadinessIssueResponse{
				Code:    issue.Code,
				Message: issue.Message,
			},
		)
	}
	for _, issue := range readiness.Warnings {
		response.Warnings = append(
			response.Warnings,
			&crmpb.SeoIntentHypothesisReadinessIssueResponse{
				Code:    issue.Code,
				Message: issue.Message,
			},
		)
	}
	return response
}
