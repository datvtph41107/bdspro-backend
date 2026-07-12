package mapper

import (
	"encoding/json"
	"time"

	seo_strategy "crm/internal/domain/seo_strategy"
	crmpb "pb/types/crm"
)

type SeoObjectiveMapper struct{}

func NewSeoObjectiveMapper() *SeoObjectiveMapper {
	return &SeoObjectiveMapper{}
}

func (m *SeoObjectiveMapper) AggregateToPb(
	aggregate *seo_strategy.SeoObjectiveAggregate,
) *crmpb.SeoObjectiveResponse {
	if aggregate == nil || aggregate.Objective == nil {
		return nil
	}

	objective := aggregate.Objective

	response := &crmpb.SeoObjectiveResponse{
		Id:               objective.ID,
		Code:             objective.Code,
		Name:             objective.Name,
		AudienceSummary:  objective.AudienceSummary,
		Situation:        objective.Situation,
		ProblemStatement: objective.ProblemStatement,
		UserOutcome:      objective.UserOutcome,
		ProductOutcome:   objective.ProductOutcome,
		BusinessOutcome:  objective.BusinessOutcome,
		NonGoals:         objective.NonGoals,
		Constraints:      objective.Constraints,
		Status:           string(objective.Status),
		Confidence:       string(objective.Confidence),
		OwnerCode:        objective.OwnerCode,
		LockVersion:      objective.LockVersion,
		Readiness:        m.ReadinessToPb(aggregate.Readiness),
		CreatedAt:        formatObjectiveTime(objective.CreatedAt),
		UpdatedAt:        formatObjectiveTime(objective.UpdatedAt),
	}

	if objective.ReviewDueAt != nil {
		value := formatObjectiveTime(*objective.ReviewDueAt)
		response.ReviewDueAt = &value
	}

	for i := range aggregate.Evidence {
		response.Evidence = append(
			response.Evidence,
			m.EvidenceToPb(&aggregate.Evidence[i]),
		)
	}

	for i := range aggregate.Metrics {
		response.Metrics = append(
			response.Metrics,
			m.MetricToPb(&aggregate.Metrics[i]),
		)
	}

	return response
}

func (m *SeoObjectiveMapper) SummaryToPb(
	objective *seo_strategy.SeoObjective,
) *crmpb.SeoObjectiveSummaryResponse {
	if objective == nil {
		return nil
	}

	return &crmpb.SeoObjectiveSummaryResponse{
		Id:          objective.ID,
		Code:        objective.Code,
		Name:        objective.Name,
		Status:      string(objective.Status),
		Confidence:  string(objective.Confidence),
		OwnerCode:   objective.OwnerCode,
		UserOutcome: objective.UserOutcome,
		LockVersion: objective.LockVersion,
		CreatedAt:   formatObjectiveTime(objective.CreatedAt),
		UpdatedAt:   formatObjectiveTime(objective.UpdatedAt),
	}
}

func (m *SeoObjectiveMapper) EvidenceToPb(
	item *seo_strategy.SeoObjectiveEvidence,
) *crmpb.SeoObjectiveEvidenceResponse {
	if item == nil {
		return nil
	}

	response := &crmpb.SeoObjectiveEvidenceResponse{
		Id:              item.ID,
		ObjectiveId:     item.ObjectiveID,
		EvidenceType:    string(item.Type),
		Stance:          string(item.Stance),
		SourceSystem:    item.SourceSystem,
		SourceReference: item.SourceReference,
		Finding:         item.Finding,
		Confidence:      string(item.Confidence),
		Status:          string(item.Status),
		MetadataJson:    string(item.Metadata),
		CreatedAt:       formatObjectiveTime(item.CreatedAt),
		UpdatedAt:       formatObjectiveTime(item.UpdatedAt),
	}

	if item.ObservedAt != nil {
		value := formatObjectiveTime(*item.ObservedAt)
		response.ObservedAt = &value
	}

	if item.ValidUntil != nil {
		value := formatObjectiveTime(*item.ValidUntil)
		response.ValidUntil = &value
	}

	if !json.Valid(item.Metadata) {
		response.MetadataJson = "{}"
	}

	return response
}

func (m *SeoObjectiveMapper) MetricToPb(
	item *seo_strategy.SeoObjectiveMetric,
) *crmpb.SeoObjectiveMetricResponse {
	if item == nil {
		return nil
	}

	response := &crmpb.SeoObjectiveMetricResponse{
		Id:                    item.ID,
		ObjectiveId:           item.ObjectiveID,
		MetricCode:            item.Code,
		MetricName:            item.Name,
		MetricLayer:           string(item.Layer),
		Description:           item.Description,
		SourceSystem:          item.SourceSystem,
		CalculationMethod:     item.CalculationMethod,
		Unit:                  item.Unit,
		TargetDirection:       string(item.TargetDirection),
		IsPrimary:             item.IsPrimary,
		InstrumentationStatus: string(item.InstrumentationStatus),
		Status:                string(item.Status),
		CreatedAt:             formatObjectiveTime(item.CreatedAt),
		UpdatedAt:             formatObjectiveTime(item.UpdatedAt),
	}

	response.BaselineValue = item.BaselineValue
	response.TargetValue = item.TargetValue

	return response
}

func (m *SeoObjectiveMapper) ReadinessToPb(
	readiness seo_strategy.ObjectiveReadiness,
) *crmpb.SeoObjectiveReadinessResponse {
	response := &crmpb.SeoObjectiveReadinessResponse{
		Ready: readiness.Ready,
	}

	for _, item := range readiness.Blockers {
		response.Blockers = append(
			response.Blockers,
			&crmpb.SeoObjectiveReadinessIssueResponse{
				Code:    item.Code,
				Message: item.Message,
			},
		)
	}

	for _, item := range readiness.Warnings {
		response.Warnings = append(
			response.Warnings,
			&crmpb.SeoObjectiveReadinessIssueResponse{
				Code:    item.Code,
				Message: item.Message,
			},
		)
	}

	return response
}

func formatObjectiveTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}
