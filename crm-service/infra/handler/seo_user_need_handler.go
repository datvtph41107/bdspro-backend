package handler

import (
	_errors "common/errors"
	"context"

	"crm/internal/dto"
	crmpb "pb/types/crm"
)

func (h *SeoStrategyHandler) CreateSeoUserNeed(
	ctx context.Context,
	req *crmpb.CreateSeoUserNeedRequest,
) (*crmpb.SeoUserNeedResponse, error) {
	result, err := h.userNeedUsecase.Create(
		ctx,
		&dto.CreateSeoUserNeedRequest{
			ObjectiveID:       req.GetObjectiveId(),
			AudienceSegmentID: req.GetAudienceSegmentId(),
			Code:              req.GetCode(),
			Name:              req.GetName(),
		},
	)
	if err != nil {
		return nil, err
	}

	return h.userNeedMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) GetSeoUserNeed(
	ctx context.Context,
	req *crmpb.GetSeoUserNeedRequest,
) (*crmpb.SeoUserNeedResponse, error) {
	result, err := h.userNeedUsecase.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return h.userNeedMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) GetSeoUserNeedList(
	ctx context.Context,
	req *crmpb.GetSeoUserNeedListRequest,
) (*crmpb.GetSeoUserNeedListResponse, error) {
	items, total, err := h.userNeedUsecase.List(
		ctx,
		&dto.SeoUserNeedListRequest{
			ObjectiveID: req.GetObjectiveId(),
			AudienceID:  req.GetAudienceId(),
			Keyword:     req.GetKeyword(),
			Status:      req.GetStatus(),
			Confidence:  req.GetConfidence(),
			Priority:    req.GetPriority(),
			Page:        req.GetPage(),
			Size:        req.GetSize(),
		},
	)
	if err != nil {
		return nil, err
	}

	response := &crmpb.GetSeoUserNeedListResponse{
		Total: total,
	}

	for i := range items {
		response.Data = append(
			response.Data,
			h.userNeedMapper.SummaryToPb(&items[i]),
		)
	}

	return response, nil
}

func (h *SeoStrategyHandler) UpdateSeoUserNeed(
	ctx context.Context,
	req *crmpb.UpdateSeoUserNeedRequest,
) (*crmpb.SeoUserNeedResponse, error) {
	reviewDueAt, err := parseOptionalStrategyTime(req.ReviewDueAt)
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}

	result, err := h.userNeedUsecase.Update(
		ctx,
		&dto.UpdateSeoUserNeedRequest{
			ID:                 req.GetId(),
			ExpectedVersion:    req.GetExpectedVersion(),
			Name:               req.GetName(),
			DecisionStage:      req.GetDecisionStage(),
			Situation:          req.GetSituation(),
			TriggerStatement:   req.GetTriggerStatement(),
			JobToBeDone:        req.GetJobToBeDone(),
			BarrierStatement:   req.GetBarrierStatement(),
			CurrentAlternative: req.GetCurrentAlternative(),
			DesiredOutcome:     req.GetDesiredOutcome(),
			FailureImpact:      req.GetFailureImpact(),
			RiskLevel:          req.GetRiskLevel(),
			Priority:           req.GetPriority(),
			Confidence:         req.GetConfidence(),
			OwnerCode:          req.GetOwnerCode(),
			ReviewDueAt:        reviewDueAt,
		},
	)
	if err != nil {
		return nil, err
	}

	return h.userNeedMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) LinkSeoUserNeedEvidence(
	ctx context.Context,
	req *crmpb.LinkSeoUserNeedEvidenceRequest,
) (*crmpb.SeoUserNeedResponse, error) {
	result, err := h.userNeedUsecase.LinkEvidence(
		ctx,
		&dto.LinkSeoUserNeedEvidenceRequest{
			UserNeedID:    req.GetUserNeedId(),
			EvidenceID:    req.GetEvidenceId(),
			RelationType:  req.GetRelationType(),
			RelevanceNote: req.GetRelevanceNote(),
		},
	)
	if err != nil {
		return nil, err
	}

	return h.userNeedMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) UnlinkSeoUserNeedEvidence(
	ctx context.Context,
	req *crmpb.UnlinkSeoUserNeedEvidenceRequest,
) (*crmpb.SeoUserNeedResponse, error) {
	result, err := h.userNeedUsecase.UnlinkEvidence(
		ctx,
		&dto.UnlinkSeoUserNeedEvidenceRequest{
			UserNeedID: req.GetUserNeedId(),
			EvidenceID: req.GetEvidenceId(),
		},
	)
	if err != nil {
		return nil, err
	}

	return h.userNeedMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) EvaluateSeoUserNeedReadiness(
	ctx context.Context,
	req *crmpb.GetSeoUserNeedRequest,
) (*crmpb.SeoUserNeedReadinessResponse, error) {
	result, err := h.userNeedUsecase.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return h.userNeedMapper.ReadinessToPb(result.Readiness), nil
}

func (h *SeoStrategyHandler) ActivateSeoUserNeed(
	ctx context.Context,
	req *crmpb.SeoUserNeedCommandRequest,
) (*crmpb.SeoUserNeedResponse, error) {
	result, err := h.userNeedUsecase.Activate(
		ctx,
		userNeedCommandToDTO(req),
	)
	if err != nil {
		return nil, err
	}

	return h.userNeedMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) PauseSeoUserNeed(
	ctx context.Context,
	req *crmpb.SeoUserNeedCommandRequest,
) (*crmpb.SeoUserNeedResponse, error) {
	result, err := h.userNeedUsecase.Pause(
		ctx,
		userNeedCommandToDTO(req),
	)
	if err != nil {
		return nil, err
	}

	return h.userNeedMapper.AggregateToPb(result), nil
}

func (h *SeoStrategyHandler) ArchiveSeoUserNeed(
	ctx context.Context,
	req *crmpb.SeoUserNeedCommandRequest,
) (*crmpb.SeoUserNeedResponse, error) {
	result, err := h.userNeedUsecase.Archive(
		ctx,
		userNeedCommandToDTO(req),
	)
	if err != nil {
		return nil, err
	}

	return h.userNeedMapper.AggregateToPb(result), nil
}

func userNeedCommandToDTO(
	req *crmpb.SeoUserNeedCommandRequest,
) *dto.SeoUserNeedCommandRequest {
	return &dto.SeoUserNeedCommandRequest{
		ID:              req.GetId(),
		ExpectedVersion: req.GetExpectedVersion(),
		Reason:          req.GetReason(),
	}
}
