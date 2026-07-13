package handler

import (
	"context"
	"crm/infra/mapper"
	"crm/internal/dto"
	"crm/internal/usecase"
	"encoding/json"
	"time"

	crmpb "pb/types/crm"
)

type SeoStrategyHandler struct {
	crmpb.UnimplementedSeoStrategyServiceServer

	usecase *usecase.SeoObjectiveUsecase
	mapper  *mapper.SeoObjectiveMapper

	audienceUsecase *usecase.SeoAudienceUsecase
	audienceMapper  *mapper.SeoAudienceMapper

	userNeedUsecase *usecase.SeoUserNeedUsecase
	userNeedMapper  *mapper.SeoUserNeedMapper
}

func NewSeoStrategyHandler(
	objectiveUsecase *usecase.SeoObjectiveUsecase,
	objectiveMapper *mapper.SeoObjectiveMapper,
	audienceUsecase *usecase.SeoAudienceUsecase,
	audienceMapper *mapper.SeoAudienceMapper,
	userNeedUsecase *usecase.SeoUserNeedUsecase,
	userNeedMapper *mapper.SeoUserNeedMapper,
) *SeoStrategyHandler {
	return &SeoStrategyHandler{
		usecase: objectiveUsecase,
		mapper:  objectiveMapper,

		audienceUsecase: audienceUsecase,
		audienceMapper:  audienceMapper,

		userNeedUsecase: userNeedUsecase,
		userNeedMapper:  userNeedMapper,
	}
}

func (h *SeoStrategyHandler) CreateSeoObjective(
	ctx context.Context,
	req *crmpb.CreateSeoObjectiveRequest,
) (*crmpb.SeoObjectiveResponse, error) {
	result, err := h.usecase.Create(
		ctx,
		&dto.CreateSeoObjectiveRequest{
			Code: req.GetCode(),
			Name: req.GetName(),
		},
	)
	if err != nil {
		return nil, err
	}

	return h.mapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) GetSeoObjective(
	ctx context.Context,
	req *crmpb.GetSeoObjectiveRequest,
) (*crmpb.SeoObjectiveResponse, error) {
	result, err := h.usecase.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return h.mapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) GetSeoObjectiveList(
	ctx context.Context,
	req *crmpb.GetSeoObjectiveListRequest,
) (*crmpb.GetSeoObjectiveListResponse, error) {
	items, total, err := h.usecase.List(
		ctx,
		&dto.SeoObjectiveListRequest{
			Keyword:    req.GetKeyword(),
			Status:     req.GetStatus(),
			Confidence: req.GetConfidence(),
			Page:       req.GetPage(),
			Size:       req.GetSize(),
		},
	)
	if err != nil {
		return nil, err
	}

	response := &crmpb.GetSeoObjectiveListResponse{
		Total: total,
	}

	for i := range items {
		response.Data = append(
			response.Data,
			h.mapper.SummaryToPb(&items[i]),
		)
	}

	return response, nil
}

func (h *SeoStrategyHandler) UpdateSeoObjectiveDefinition(
	ctx context.Context,
	req *crmpb.UpdateSeoObjectiveDefinitionRequest,
) (*crmpb.SeoObjectiveResponse, error) {
	result, err := h.usecase.UpdateDefinition(
		ctx,
		&dto.UpdateSeoObjectiveDefinitionRequest{
			ID:               req.GetId(),
			ExpectedVersion:  req.GetExpectedVersion(),
			Name:             req.GetName(),
			AudienceSummary:  req.GetAudienceSummary(),
			Situation:        req.GetSituation(),
			ProblemStatement: req.GetProblemStatement(),
			UserOutcome:      req.GetUserOutcome(),
			ProductOutcome:   req.GetProductOutcome(),
			BusinessOutcome:  req.GetBusinessOutcome(),
			NonGoals:         req.GetNonGoals(),
			Constraints:      req.GetConstraints(),
			Confidence:       req.GetConfidence(),
			OwnerCode:        req.GetOwnerCode(),
			ReviewDueAt:      parseOptionalObjectiveTime(req.ReviewDueAt),
		},
	)
	if err != nil {
		return nil, err
	}

	return h.mapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) AddSeoObjectiveEvidence(
	ctx context.Context,
	req *crmpb.AddSeoObjectiveEvidenceRequest,
) (*crmpb.SeoObjectiveResponse, error) {
	metadata := json.RawMessage(req.GetMetadataJson())
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}

	result, err := h.usecase.AddEvidence(
		ctx,
		&dto.AddSeoObjectiveEvidenceRequest{
			ObjectiveID:     req.GetObjectiveId(),
			EvidenceType:    req.GetEvidenceType(),
			Stance:          req.GetStance(),
			SourceSystem:    req.GetSourceSystem(),
			SourceReference: req.GetSourceReference(),
			Finding:         req.GetFinding(),
			Confidence:      req.GetConfidence(),
			ObservedAt:      parseOptionalObjectiveTime(req.ObservedAt),
			ValidUntil:      parseOptionalObjectiveTime(req.ValidUntil),
			Metadata:        metadata,
		},
	)
	if err != nil {
		return nil, err
	}

	return h.mapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) UpsertSeoObjectiveMetric(
	ctx context.Context,
	req *crmpb.UpsertSeoObjectiveMetricRequest,
) (*crmpb.SeoObjectiveResponse, error) {
	result, err := h.usecase.UpsertMetric(
		ctx,
		&dto.UpsertSeoObjectiveMetricRequest{
			ObjectiveID:           req.GetObjectiveId(),
			MetricCode:            req.GetMetricCode(),
			MetricName:            req.GetMetricName(),
			MetricLayer:           req.GetMetricLayer(),
			Description:           req.GetDescription(),
			SourceSystem:          req.GetSourceSystem(),
			CalculationMethod:     req.GetCalculationMethod(),
			Unit:                  req.GetUnit(),
			TargetDirection:       req.GetTargetDirection(),
			BaselineValue:         req.BaselineValue,
			TargetValue:           req.TargetValue,
			IsPrimary:             req.GetIsPrimary(),
			InstrumentationStatus: req.GetInstrumentationStatus(),
			Status:                req.GetStatus(),
		},
	)
	if err != nil {
		return nil, err
	}

	return h.mapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) EvaluateSeoObjectiveReadiness(
	ctx context.Context,
	req *crmpb.GetSeoObjectiveRequest,
) (*crmpb.SeoObjectiveReadinessResponse, error) {
	result, err := h.usecase.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return h.mapper.ReadinessToPb(result.Readiness), nil
}

func (h *SeoStrategyHandler) ActivateSeoObjective(
	ctx context.Context,
	req *crmpb.SeoObjectiveCommandRequest,
) (*crmpb.SeoObjectiveResponse, error) {
	result, err := h.usecase.Activate(
		ctx,
		objectiveCommandToDTO(req),
	)
	if err != nil {
		return nil, err
	}

	return h.mapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) PauseSeoObjective(
	ctx context.Context,
	req *crmpb.SeoObjectiveCommandRequest,
) (*crmpb.SeoObjectiveResponse, error) {
	result, err := h.usecase.Pause(
		ctx,
		objectiveCommandToDTO(req),
	)
	if err != nil {
		return nil, err
	}

	return h.mapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) ArchiveSeoObjective(
	ctx context.Context,
	req *crmpb.SeoObjectiveCommandRequest,
) (*crmpb.SeoObjectiveResponse, error) {
	result, err := h.usecase.Archive(
		ctx,
		objectiveCommandToDTO(req),
	)
	if err != nil {
		return nil, err
	}

	return h.mapper.AggregateToPb(result), nil
}

func objectiveCommandToDTO(
	req *crmpb.SeoObjectiveCommandRequest,
) *dto.SeoObjectiveCommandRequest {
	return &dto.SeoObjectiveCommandRequest{
		ID:              req.GetId(),
		ExpectedVersion: req.GetExpectedVersion(),
		Reason:          req.GetReason(),
	}
}

func parseOptionalObjectiveTime(
	value *string,
) *time.Time {
	if value == nil || *value == "" {
		return nil
	}

	parsed, err := time.Parse(time.RFC3339, *value)
	if err != nil {
		return nil
	}

	return &parsed
}
