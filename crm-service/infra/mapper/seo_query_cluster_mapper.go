package mapper

import (
	seo_strategy "crm/internal/domain/seo_strategy"
	crmpb "pb/types/crm"
)

type SeoQueryClusterMapper struct {
	intentHypothesisMapper *SeoIntentHypothesisMapper
	userNeedMapper         *SeoUserNeedMapper
	queryEvidenceMapper    *SeoQueryEvidenceMapper
}

func NewSeoQueryClusterMapper(
	intentHypothesisMapper *SeoIntentHypothesisMapper,
	userNeedMapper *SeoUserNeedMapper,
	queryEvidenceMapper *SeoQueryEvidenceMapper,
) *SeoQueryClusterMapper {
	return &SeoQueryClusterMapper{
		intentHypothesisMapper: intentHypothesisMapper,
		userNeedMapper:         userNeedMapper,
		queryEvidenceMapper:    queryEvidenceMapper,
	}
}

func (m *SeoQueryClusterMapper) AggregateToPb(
	aggregate *seo_strategy.SeoQueryClusterAggregate,
) *crmpb.SeoQueryClusterResponse {
	if aggregate == nil || aggregate.Cluster == nil {
		return nil
	}
	cluster := aggregate.Cluster
	response := &crmpb.SeoQueryClusterResponse{
		Id:                 cluster.ID,
		ObjectiveId:        cluster.ObjectiveID,
		UserNeedId:         cluster.UserNeedID,
		IntentHypothesisId: cluster.IntentHypothesisID,
		Code:               cluster.Code,
		Name:               cluster.Name,
		SearchPromise:      cluster.SearchPromise,
		InclusionRule:      cluster.InclusionRule,
		ExclusionRule:      cluster.ExclusionRule,
		Status:             string(cluster.Status),
		Confidence:         string(cluster.Confidence),
		OwnerCode:          cluster.OwnerCode,
		LockVersion:        cluster.LockVersion,
		Readiness:          m.ReadinessToPb(aggregate.Readiness),
		CreatedAt:          formatStrategyTime(cluster.CreatedAt),
		UpdatedAt:          formatStrategyTime(cluster.UpdatedAt),
	}
	if cluster.ReviewDueAt != nil {
		value := formatStrategyTime(*cluster.ReviewDueAt)
		response.ReviewDueAt = &value
	}
	if aggregate.Intent != nil {
		response.IntentHypothesis = m.intentHypothesisMapper.SummaryToPb(aggregate.Intent)
	}
	if aggregate.UserNeed != nil {
		response.UserNeed = m.userNeedMapper.SummaryToPb(aggregate.UserNeed)
	}
	for i := range aggregate.Members {
		response.Members = append(response.Members, m.MemberToPb(&aggregate.Members[i]))
	}
	return response
}

func (m *SeoQueryClusterMapper) SummaryToPb(
	cluster *seo_strategy.SeoQueryCluster,
) *crmpb.SeoQueryClusterSummaryResponse {
	if cluster == nil {
		return nil
	}
	return &crmpb.SeoQueryClusterSummaryResponse{
		Id:                 cluster.ID,
		ObjectiveId:        cluster.ObjectiveID,
		UserNeedId:         cluster.UserNeedID,
		IntentHypothesisId: cluster.IntentHypothesisID,
		Code:               cluster.Code,
		Name:               cluster.Name,
		SearchPromise:      cluster.SearchPromise,
		Status:             string(cluster.Status),
		Confidence:         string(cluster.Confidence),
		OwnerCode:          cluster.OwnerCode,
		LockVersion:        cluster.LockVersion,
		CreatedAt:          formatStrategyTime(cluster.CreatedAt),
		UpdatedAt:          formatStrategyTime(cluster.UpdatedAt),
	}
}

func (m *SeoQueryClusterMapper) MemberToPb(
	linked *seo_strategy.LinkedQueryClusterMember,
) *crmpb.SeoQueryClusterMemberResponse {
	if linked == nil {
		return nil
	}
	return &crmpb.SeoQueryClusterMemberResponse{
		MembershipType: string(linked.Member.MembershipType),
		Rationale:      linked.Member.Rationale,
		Status:         string(linked.Member.Status),
		CreatedAt:      formatStrategyTime(linked.Member.CreatedAt),
		UpdatedAt:      formatStrategyTime(linked.Member.UpdatedAt),
		QueryEvidence:  m.queryEvidenceMapper.ToPb(&linked.Evidence),
	}
}

func (m *SeoQueryClusterMapper) ReadinessToPb(
	readiness seo_strategy.QueryClusterReadiness,
) *crmpb.SeoQueryClusterReadinessResponse {
	response := &crmpb.SeoQueryClusterReadinessResponse{
		Ready:                  readiness.Ready,
		AnswerEquivalenceReady: readiness.AnswerEquivalenceReady,
		PrimaryCount:           readiness.PrimaryCount,
		IncludedCount:          readiness.IncludedCount,
		ExcludedCount:          readiness.ExcludedCount,
	}
	for _, issue := range readiness.Blockers {
		response.Blockers = append(response.Blockers, &crmpb.SeoQueryClusterReadinessIssueResponse{
			Code:    issue.Code,
			Message: issue.Message,
		})
	}
	for _, issue := range readiness.Warnings {
		response.Warnings = append(response.Warnings, &crmpb.SeoQueryClusterReadinessIssueResponse{
			Code:    issue.Code,
			Message: issue.Message,
		})
	}
	return response
}
